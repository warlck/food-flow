package imageapp

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/warlck/food-flow/app/sdk/errs"
	"github.com/warlck/food-flow/business/domain/imagebus"
)

// Image represents information about a tracked uploaded image.
type Image struct {
	ID           string `json:"imageId"`
	RestaurantID string `json:"restaurantId"`
	EntityType   string `json:"entityType"`
	ObjectPath   string `json:"objectPath"`
	PublicURL    string `json:"publicUrl"`
	ContentType  string `json:"contentType"`
	SizeBytes    int64  `json:"sizeBytes"`
	Status       string `json:"status"`
	DateCreated  string `json:"dateCreated"`
	DateUpdated  string `json:"dateUpdated"`
}

// Encode implements encoder interface.
func (app Image) Encode() ([]byte, string, error) {
	data, err := json.Marshal(app)
	return data, "application/json", err
}

// ToAppImage converts a business layer image to an app layer image.
func ToAppImage(bus imagebus.Image) Image {
	return Image{
		ID:           bus.ID.String(),
		RestaurantID: bus.RestaurantID.String(),
		EntityType:   bus.EntityType,
		ObjectPath:   bus.ObjectPath,
		PublicURL:    bus.PublicURL,
		ContentType:  bus.ContentType,
		SizeBytes:    bus.SizeBytes,
		Status:       bus.Status,
		DateCreated:  bus.DateCreated.Format(time.RFC3339),
		DateUpdated:  bus.DateUpdated.Format(time.RFC3339),
	}
}

// ToAppImages converts a list of business layer images to app layer images.
func ToAppImages(bus []imagebus.Image) []Image {
	items := make([]Image, len(bus))
	for i, img := range bus {
		items[i] = ToAppImage(img)
	}
	return items
}

// =============================================================================

// NewUpload defines data needed to request a signed image upload URL.
type NewUpload struct {
	// TODO: Refactor DTOs in other app domains to use uuid.UUID directly instead of string to avoid manual parsing.
	RestaurantID uuid.UUID `json:"restaurantId" validate:"required"`
	EntityType   string    `json:"entityType" validate:"required,oneof=restaurant menu_item"`
	ContentType  string    `json:"contentType" validate:"required,oneof=image/jpeg image/png image/webp"`
	SizeBytes    int64     `json:"sizeBytes" validate:"required,gt=0"`
}

// Decode implements decoder interface.
func (app *NewUpload) Decode(data []byte) error {
	return json.Unmarshal(data, app)
}

// Validate checks request validity.
func (app NewUpload) Validate() error {
	if err := errs.Check(app); err != nil {
		return fmt.Errorf("validate: %w", err)
	}
	return nil
}

// =============================================================================

// UploadGrant defines the signed upload target plus the tracked image.
type UploadGrant struct {
	Image     Image             `json:"image"`
	UploadURL string            `json:"uploadUrl"`
	Method    string            `json:"method"`
	Headers   map[string]string `json:"headers"`
	ExpiresAt string            `json:"expiresAt"`
}

// Encode implements encoder interface.
func (app UploadGrant) Encode() ([]byte, string, error) {
	data, err := json.Marshal(app)
	return data, "application/json", err
}

// ToAppUploadGrant converts a business layer upload grant to an app layer grant.
func ToAppUploadGrant(bus imagebus.UploadGrant) UploadGrant {
	return UploadGrant{
		Image:     ToAppImage(bus.Image),
		UploadURL: bus.UploadURL,
		Method:    bus.Method,
		Headers:   bus.Headers,
		ExpiresAt: bus.ExpiresAt.Format(time.RFC3339),
	}
}
