package categoryapp

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/warlck/food-flow/app/sdk/errs"
	"github.com/warlck/food-flow/app/sdk/mid"
	"github.com/warlck/food-flow/app/sdk/query"
	"github.com/warlck/food-flow/business/domain/categorybus"
	"github.com/warlck/food-flow/business/domain/restaurantbus"
	"github.com/warlck/food-flow/business/sdk/order"
	"github.com/warlck/food-flow/business/sdk/page"
	"github.com/warlck/food-flow/foundation/web"
)

// Handlers manages the set of category endpoints.
type app struct {
	categoryBus   *categorybus.Business
	restaurantBus *restaurantbus.Business
}

// newApp constructs a handlers for route access.
func newApp(categoryBus *categorybus.Business, restaurantBus *restaurantbus.Business) *app {
	return &app{
		categoryBus:   categoryBus,
		restaurantBus: restaurantBus,
	}
}

// Create adds a new category to the system.
func (a *app) create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var app NewCategory
	if err := web.Decode(r, &app); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	nc, err := toBusNewCategory(app)
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	rest, err := a.restaurantBus.QueryByID(ctx, nc.RestaurantID)
	if err != nil {
		return errs.New(errs.InvalidArgument, fmt.Errorf("restaurant lookup: %w", err))
	}

	claims := mid.GetClaims(ctx)
	if !claims.IsOrgAuthorized(rest.OrganizationID) {
		return errs.Newf(errs.PermissionDenied, "user not in organization %s", rest.OrganizationID)
	}

	cat, err := a.categoryBus.Create(ctx, nc)
	if err != nil {
		return fmt.Errorf("create: category[%+v]: %w", cat, err)
	}

	return web.Respond(ctx, w, ToAppCategory(cat), http.StatusCreated)
}

// Reorder updates the rank of all categories within a restaurant.
func (a *app) reorder(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var app ReorderCategories
	if err := web.Decode(r, &app); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	restID, err := uuid.Parse(app.RestaurantID)
	if err != nil {
		return errs.NewFieldErrors("restaurantId", err)
	}

	rest, err := a.restaurantBus.QueryByID(ctx, restID)
	if err != nil {
		return errs.New(errs.InvalidArgument, fmt.Errorf("restaurant lookup: %w", err))
	}

	claims := mid.GetClaims(ctx)
	if !claims.IsOrgAuthorized(rest.OrganizationID) {
		return errs.Newf(errs.PermissionDenied, "user not in organization %s", rest.OrganizationID)
	}

	orderedIDs := make([]uuid.UUID, len(app.OrderedIDs))
	for i, idStr := range app.OrderedIDs {
		id, err := uuid.Parse(idStr)
		if err != nil {
			return errs.NewFieldErrors("orderedIds", err)
		}
		orderedIDs[i] = id
	}

	reordered, err := a.categoryBus.Reorder(ctx, restID, orderedIDs)
	if err != nil {
		if errors.Is(err, categorybus.ErrInvalidReorder) {
			return errs.New(errs.InvalidArgument, err)
		}
		return errs.Newf(errs.Internal, "reorder: %s", err)
	}

	return web.Respond(ctx, w, ToAppCategories(reordered), http.StatusOK)
}

// Query retrieves a list of categories based on query parameters.
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
		return err.(*errs.Error)
	}

	orderBy, err := order.Parse(orderByFields, qp.OrderBy, defaultOrderBy)
	if err != nil {
		return errs.NewFieldErrors("order", err)
	}

	categories, err := a.categoryBus.Query(ctx, filter, orderBy, pg)
	if err != nil {
		return errs.Newf(errs.Internal, "query: %s", err)
	}

	total, err := a.categoryBus.Count(ctx, filter)
	if err != nil {
		return errs.Newf(errs.Internal, "count: %s", err)
	}

	result := query.NewResult(ToAppCategories(categories), total, pg)
	return web.Respond(ctx, w, result, http.StatusOK)
}

// QueryByID retrieves a category by its ID.
func (a *app) queryByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	categoryIDStr := web.Param(r, "category_id")

	categoryID, err := uuid.Parse(categoryIDStr)
	if err != nil {
		return errs.NewFieldErrors("category_id", err)
	}

	cat, err := a.categoryBus.QueryByID(ctx, categoryID)
	if err != nil {
		return fmt.Errorf("querybyid: categoryID[%s]: %w", categoryID, err)
	}

	return web.Respond(ctx, w, ToAppCategory(cat), http.StatusOK)
}

// Update modifies a category.
func (a *app) update(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	categoryIDStr := web.Param(r, "category_id")

	categoryID, err := uuid.Parse(categoryIDStr)
	if err != nil {
		return errs.NewFieldErrors("category_id", err)
	}

	var app UpdateCategory
	if err := web.Decode(r, &app); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	uc, err := toBusUpdateCategory(app)
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	cat, err := a.categoryBus.QueryByID(ctx, categoryID)
	if err != nil {
		return fmt.Errorf("querybyid: categoryID[%s]: %w", categoryID, err)
	}

	rest, err := a.restaurantBus.QueryByID(ctx, cat.RestaurantID)
	if err != nil {
		return errs.New(errs.InvalidArgument, fmt.Errorf("restaurant lookup: %w", err))
	}

	claims := mid.GetClaims(ctx)
	if !claims.IsOrgAuthorized(rest.OrganizationID) {
		return errs.Newf(errs.PermissionDenied, "user not in organization %s", rest.OrganizationID)
	}

	updCat, err := a.categoryBus.Update(ctx, cat, uc)
	if err != nil {
		return errs.Newf(errs.Internal, "update: categoryID[%s] uc[%+v]: %s", categoryID, uc, err)
	}

	return web.Respond(ctx, w, ToAppCategory(updCat), http.StatusOK)
}

// Delete removes a category.
func (a *app) delete(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	categoryIDStr := web.Param(r, "category_id")

	categoryID, err := uuid.Parse(categoryIDStr)
	if err != nil {
		return errs.NewFieldErrors("category_id", err)
	}

	cat, err := a.categoryBus.QueryByID(ctx, categoryID)
	if err != nil {
		return fmt.Errorf("querybyid: categoryID[%s]: %w", categoryID, err)
	}

	rest, err := a.restaurantBus.QueryByID(ctx, cat.RestaurantID)
	if err != nil {
		return errs.New(errs.InvalidArgument, fmt.Errorf("restaurant lookup: %w", err))
	}

	claims := mid.GetClaims(ctx)
	if !claims.IsOrgAuthorized(rest.OrganizationID) {
		return errs.Newf(errs.PermissionDenied, "user not in organization %s", rest.OrganizationID)
	}

	if err := a.categoryBus.Delete(ctx, cat); err != nil {
		return errs.Newf(errs.Internal, "delete: categoryID[%s]: %s", categoryID, err)
	}

	return web.Respond(ctx, w, struct{}{}, http.StatusNoContent)
}
