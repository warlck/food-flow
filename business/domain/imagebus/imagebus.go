// Package imagebus provides business APIs for tracked image uploads.
package imagebus

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/warlck/food-flow/business/sdk/order"
	"github.com/warlck/food-flow/business/sdk/page"
	"github.com/warlck/food-flow/foundation/logger"
	"github.com/warlck/food-flow/foundation/storage"
)

// DefaultMaxSizeBytes bounds an upload when no explicit limit is configured.
const DefaultMaxSizeBytes int64 = 5 << 20 // 5 MiB

// allowedContentTypes maps permitted MIME types to file extensions.
var allowedContentTypes = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/webp": ".webp",
}

// entityDirs maps entity types to storage object path prefixes.
var entityDirs = map[string]string{
	EntityTypeRestaurant: "restaurants",
	EntityTypeMenuItem:   "menu-items",
}

// Storer declares the persistence needed by the image domain.
type Storer interface {
	Create(ctx context.Context, img Image) error
	Update(ctx context.Context, img Image) error
	Delete(ctx context.Context, img Image) error
	Query(ctx context.Context, filter QueryFilter, orderBy order.By, page page.Page) ([]Image, error)
	Count(ctx context.Context, filter QueryFilter) (int, error)
	QueryByID(ctx context.Context, imageID uuid.UUID) (Image, error)
}

// Business manages the set of APIs for image upload access.
type Business struct {
	log          *logger.Logger
	storer       Storer
	signer       storage.Signer
	maxSizeBytes int64
}

// NewBusiness constructs an imagebus business API for use.
func NewBusiness(log *logger.Logger, storer Storer, signer storage.Signer, maxSizeBytes int64) *Business {
	if maxSizeBytes <= 0 {
		maxSizeBytes = DefaultMaxSizeBytes
	}

	return &Business{
		log:          log,
		storer:       storer,
		signer:       signer,
		maxSizeBytes: maxSizeBytes,
	}
}

// MaxSizeBytes returns the configured upload size limit.
func (b *Business) MaxSizeBytes() int64 {
	return b.maxSizeBytes
}

// UploadRequest describes a desired image upload.
type UploadRequest struct {
	RestaurantID uuid.UUID
	EntityType   string
	ContentType  string
	SizeBytes    int64
	UploadedBy   *uuid.UUID
}

// UploadGrant carries the signed upload target plus the tracked image row.
type UploadGrant struct {
	Image     Image
	UploadURL string
	Method    string
	Headers   map[string]string
	ExpiresAt time.Time
}

// CreateUpload validates the request, mints a signed upload URL, and records
// a pending image row.
func (b *Business) CreateUpload(ctx context.Context, req UploadRequest) (UploadGrant, error) {
	ext, err := b.validateUploadRequest(req)
	if err != nil {
		return UploadGrant{}, err
	}

	objectPath := fmt.Sprintf("%s/%s/%s%s", entityDirs[req.EntityType], req.RestaurantID, uuid.New(), ext)

	signed, err := b.signer.SignUpload(ctx, storage.SignRequest{
		ObjectPath:  objectPath,
		ContentType: req.ContentType,
		SizeBytes:   req.SizeBytes,
	})
	if err != nil {
		return UploadGrant{}, fmt.Errorf("sign upload: %w", err)
	}

	now := time.Now()
	img := Image{
		ID:           uuid.New(),
		RestaurantID: req.RestaurantID,
		EntityType:   req.EntityType,
		ObjectPath:   objectPath,
		PublicURL:    signed.PublicURL,
		ContentType:  req.ContentType,
		SizeBytes:    req.SizeBytes,
		Status:       StatusPending,
		UploadedBy:   req.UploadedBy,
		DateCreated:  now,
		DateUpdated:  now,
	}

	if err := b.storer.Create(ctx, img); err != nil {
		return UploadGrant{}, fmt.Errorf("create: %w", err)
	}

	return UploadGrant{
		Image:     img,
		UploadURL: signed.UploadURL,
		Method:    signed.Method,
		Headers:   signed.Headers,
		ExpiresAt: signed.ExpiresAt,
	}, nil
}

// ConfirmUpload verifies the object landed in storage within the size limits
// and marks the image confirmed. Confirming an already-confirmed image is a
// no-op so clients can safely retry.
func (b *Business) ConfirmUpload(ctx context.Context, imageID uuid.UUID) (Image, error) {
	img, err := b.storer.QueryByID(ctx, imageID)
	if err != nil {
		return Image{}, fmt.Errorf("query by id: %w", err)
	}

	if img.Status == StatusConfirmed {
		return img, nil
	}

	info, err := b.signer.Stat(ctx, img.ObjectPath)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return Image{}, fmt.Errorf("%w: upload has not completed", ErrInvalid)
		}
		return Image{}, fmt.Errorf("stat object: %w", err)
	}

	if info.SizeBytes <= 0 || info.SizeBytes > b.maxSizeBytes {
		// Remove all trace of the out-of-bounds upload.
		if err := b.signer.Delete(ctx, img.ObjectPath); err != nil {
			b.log.Error(ctx, "confirm upload: cleanup object", "image_id", img.ID, "err", err)
		}
		if err := b.storer.Delete(ctx, img); err != nil {
			b.log.Error(ctx, "confirm upload: cleanup row", "image_id", img.ID, "err", err)
		}
		return Image{}, fmt.Errorf("%w: uploaded object size %d exceeds the limit of %d bytes", ErrInvalid, info.SizeBytes, b.maxSizeBytes)
	}

	img.Status = StatusConfirmed
	img.SizeBytes = info.SizeBytes
	img.DateUpdated = time.Now()

	if err := b.storer.Update(ctx, img); err != nil {
		return Image{}, fmt.Errorf("update: %w", err)
	}

	return img, nil
}

// Delete removes the stored object and its tracking row.
func (b *Business) Delete(ctx context.Context, imageID uuid.UUID) error {
	img, err := b.storer.QueryByID(ctx, imageID)
	if err != nil {
		return fmt.Errorf("query by id: %w", err)
	}

	if err := b.signer.Delete(ctx, img.ObjectPath); err != nil {
		return fmt.Errorf("delete object: %w", err)
	}

	if err := b.storer.Delete(ctx, img); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}

// Query retrieves a list of images based on the filter.
func (b *Business) Query(ctx context.Context, filter QueryFilter, orderBy order.By, page page.Page) ([]Image, error) {
	images, err := b.storer.Query(ctx, filter, orderBy, page)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	return images, nil
}

// Count returns the number of images matching the filter.
func (b *Business) Count(ctx context.Context, filter QueryFilter) (int, error) {
	total, err := b.storer.Count(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("count: %w", err)
	}
	return total, nil
}

// QueryByID retrieves an image by its ID.
func (b *Business) QueryByID(ctx context.Context, imageID uuid.UUID) (Image, error) {
	img, err := b.storer.QueryByID(ctx, imageID)
	if err != nil {
		return Image{}, fmt.Errorf("query by id: %w", err)
	}
	return img, nil
}

func (b *Business) validateUploadRequest(req UploadRequest) (string, error) {
	if req.RestaurantID == uuid.Nil {
		return "", fmt.Errorf("%w: restaurant_id is required", ErrInvalid)
	}

	if _, ok := entityDirs[req.EntityType]; !ok {
		return "", fmt.Errorf("%w: invalid entity_type %q", ErrInvalid, req.EntityType)
	}

	ext, ok := allowedContentTypes[req.ContentType]
	if !ok {
		return "", fmt.Errorf("%w: unsupported content_type %q", ErrInvalid, req.ContentType)
	}

	if req.SizeBytes <= 0 {
		return "", fmt.Errorf("%w: size_bytes must be greater than 0", ErrInvalid)
	}

	if req.SizeBytes > b.maxSizeBytes {
		return "", fmt.Errorf("%w: size_bytes %d exceeds the limit of %d bytes", ErrInvalid, req.SizeBytes, b.maxSizeBytes)
	}

	return ext, nil
}
