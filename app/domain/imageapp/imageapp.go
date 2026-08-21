// Package imageapp provides HTTP handlers for tracked image uploads.
package imageapp

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/warlck/food-flow/app/sdk/errs"
	"github.com/warlck/food-flow/app/sdk/mid"
	"github.com/warlck/food-flow/app/sdk/query"
	"github.com/warlck/food-flow/business/domain/imagebus"
	"github.com/warlck/food-flow/business/domain/restaurantbus"
	"github.com/warlck/food-flow/business/sdk/order"
	"github.com/warlck/food-flow/business/sdk/page"
	"github.com/warlck/food-flow/foundation/storage"
	"github.com/warlck/food-flow/foundation/web"
)

type app struct {
	imageBus      *imagebus.Business
	restaurantBus *restaurantbus.Business
	localStore    storage.LocalStore
}

func newApp(imageBus *imagebus.Business, restaurantBus *restaurantbus.Business, localStore storage.LocalStore) *app {
	return &app{
		imageBus:      imageBus,
		restaurantBus: restaurantBus,
		localStore:    localStore,
	}
}

// createUpload mints a signed upload URL and records a pending image.
func (a *app) createUpload(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var req NewUpload
	if err := web.Decode(r, &req); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	rest, err := a.restaurantBus.QueryByID(ctx, req.RestaurantID)
	if err != nil {
		return errs.New(errs.InvalidArgument, fmt.Errorf("restaurant lookup: %w", err))
	}

	claims := mid.GetClaims(ctx)
	if !claims.IsOrgAuthorized(rest.OrganizationID) {
		return errs.Newf(errs.PermissionDenied, "user not in organization %s", rest.OrganizationID)
	}

	var uploadedBy *uuid.UUID
	if userID, err := mid.GetUserID(ctx); err == nil {
		uploadedBy = &userID
	}

	grant, err := a.imageBus.CreateUpload(ctx, imagebus.UploadRequest{
		RestaurantID: req.RestaurantID,
		EntityType:   req.EntityType,
		ContentType:  req.ContentType,
		SizeBytes:    req.SizeBytes,
		UploadedBy:   uploadedBy,
	})
	if err != nil {
		if errors.Is(err, imagebus.ErrInvalid) {
			return errs.New(errs.InvalidArgument, err)
		}
		return fmt.Errorf("create upload: %w", err)
	}

	return web.Respond(ctx, w, ToAppUploadGrant(grant), http.StatusCreated)
}

// complete marks an upload confirmed after verifying the object in storage.
func (a *app) complete(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	imageID, err := uuid.Parse(web.Param(r, "image_id"))
	if err != nil {
		return errs.NewFieldErrors("image_id", err)
	}

	img, err := a.imageBus.ConfirmUpload(ctx, imageID)
	if err != nil {
		switch {
		case errors.Is(err, imagebus.ErrNotFound):
			return errs.New(errs.NotFound, err)
		case errors.Is(err, imagebus.ErrInvalid):
			return errs.New(errs.InvalidArgument, err)
		default:
			return fmt.Errorf("complete: %w", err)
		}
	}

	return web.Respond(ctx, w, ToAppImage(img), http.StatusOK)
}

// delete removes the stored object and its tracking row.
func (a *app) delete(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	imageID, err := uuid.Parse(web.Param(r, "image_id"))
	if err != nil {
		return errs.NewFieldErrors("image_id", err)
	}

	img, err := a.imageBus.QueryByID(ctx, imageID)
	if err != nil {
		if errors.Is(err, imagebus.ErrNotFound) {
			return errs.New(errs.NotFound, err)
		}
		return fmt.Errorf("query by id: %w", err)
	}

	rest, err := a.restaurantBus.QueryByID(ctx, img.RestaurantID)
	if err != nil {
		return errs.New(errs.InvalidArgument, fmt.Errorf("restaurant lookup: %w", err))
	}

	claims := mid.GetClaims(ctx)
	if !claims.IsOrgAuthorized(rest.OrganizationID) {
		return errs.Newf(errs.PermissionDenied, "user not in organization %s", rest.OrganizationID)
	}

	if err := a.imageBus.Delete(ctx, imageID); err != nil {
		if errors.Is(err, imagebus.ErrNotFound) {
			return errs.New(errs.NotFound, err)
		}
		return fmt.Errorf("delete: %w", err)
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}

// query retrieves a list of tracked images based on the filter.
func (a *app) query(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	qp, err := parseQueryParams(r)
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	pg, err := page.Parse(qp.Page, qp.Rows)
	if err != nil {
		return errs.NewFieldErrors("page", err)
	}

	filter, err := parseFilter(qp)
	if err != nil {
		// TODO: Apply this safer error assertion pattern to other app domains (userapp, restaurantapp, etc.)
		return errs.NewError(err)
	}

	orderBy, err := order.Parse(orderByFields, qp.OrderBy, defaultOrderBy)
	if err != nil {
		return errs.NewFieldErrors("order", err)
	}

	images, err := a.imageBus.Query(ctx, filter, orderBy, pg)
	if err != nil {
		return errs.Newf(errs.Internal, "query: %s", err)
	}

	total, err := a.imageBus.Count(ctx, filter)
	if err != nil {
		return errs.Newf(errs.Internal, "count: %s", err)
	}

	result := query.NewResult(ToAppImages(images), total, pg)
	return web.Respond(ctx, w, result, http.StatusOK)
}

// localUpload receives upload bytes for the local development backend. The
// signed URL produced by the local signer points here.
func (a *app) localUpload(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	objectPath := web.Param(r, "path")
	if objectPath == "" {
		return errs.Newf(errs.InvalidArgument, "missing object path")
	}

	if err := a.localStore.Put(ctx, objectPath, r.Body, a.imageBus.MaxSizeBytes()); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	return web.Respond(ctx, w, map[string]string{"status": "uploaded"}, http.StatusOK)
}

// localDownload serves objects stored by the local development backend.
func (a *app) localDownload(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	objectPath := web.Param(r, "path")
	if objectPath == "" {
		return errs.Newf(errs.InvalidArgument, "missing object path")
	}

	fullPath, err := a.localStore.LocalPath(objectPath)
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	if _, err := os.Stat(fullPath); err != nil {
		return errs.New(errs.NotFound, err)
	}

	http.ServeFile(w, r, fullPath)
	return nil
}
