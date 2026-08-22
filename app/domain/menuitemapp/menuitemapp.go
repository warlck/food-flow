package menuitemapp

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
	"github.com/warlck/food-flow/business/domain/menuitembus"
	"github.com/warlck/food-flow/business/domain/restaurantbus"
	"github.com/warlck/food-flow/business/sdk/order"
	"github.com/warlck/food-flow/business/sdk/page"
	"github.com/warlck/food-flow/foundation/web"
)

// app manages the set of menu item endpoints.
type app struct {
	menuItemBus   *menuitembus.Business
	restaurantBus *restaurantbus.Business
	categoryBus   *categorybus.Business
}

// newApp constructs a handlers for route access.
func newApp(menuItemBus *menuitembus.Business, restaurantBus *restaurantbus.Business, categoryBus *categorybus.Business) *app {
	return &app{
		menuItemBus:   menuItemBus,
		restaurantBus: restaurantBus,
		categoryBus:   categoryBus,
	}
}

// create adds a new menu item to the system.
func (a *app) create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var app NewMenuItem
	if err := web.Decode(r, &app); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	nb, err := toBusNewMenuItem(app)
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	rest, err := a.restaurantBus.QueryByID(ctx, nb.RestaurantID)
	if err != nil {
		return errs.New(errs.InvalidArgument, fmt.Errorf("restaurant lookup: %w", err))
	}

	claims := mid.GetClaims(ctx)
	if !claims.IsOrgAuthorized(rest.OrganizationID) {
		return errs.Newf(errs.PermissionDenied, "user not in organization %s", rest.OrganizationID)
	}

	item, err := a.menuItemBus.Create(ctx, nb)
	if err != nil {
		return fmt.Errorf("create: menuItem[%+v]: %w", item, err)
	}

	return web.Respond(ctx, w, ToAppMenuItem(item), http.StatusCreated)
}

// query retrieves a list of menu items based on query parameters.
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

	items, err := a.menuItemBus.Query(ctx, filter, orderBy, pg)
	if err != nil {
		return errs.Newf(errs.Internal, "query: %s", err)
	}

	total, err := a.menuItemBus.Count(ctx, filter)
	if err != nil {
		return errs.Newf(errs.Internal, "count: %s", err)
	}

	result := query.NewResult(ToAppMenuItems(items), total, pg)
	return web.Respond(ctx, w, result, http.StatusOK)
}

// queryByID retrieves a menu item by its ID.
func (a *app) queryByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	menuItemIDStr := web.Param(r, "menuitem_id")

	menuItemID, err := uuid.Parse(menuItemIDStr)
	if err != nil {
		return errs.NewFieldErrors("menuitem_id", err)
	}

	item, err := a.menuItemBus.QueryByID(ctx, menuItemID)
	if err != nil {
		return fmt.Errorf("querybyid: menuItemID[%s]: %w", menuItemID, err)
	}

	return web.Respond(ctx, w, ToAppMenuItem(item), http.StatusOK)
}

// update modifies an existing menu item.
func (a *app) update(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var app UpdateMenuItem
	if err := web.Decode(r, &app); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	menuItemIDStr := web.Param(r, "menuitem_id")
	menuItemID, err := uuid.Parse(menuItemIDStr)
	if err != nil {
		return errs.NewFieldErrors("menuitem_id", err)
	}

	item, err := a.menuItemBus.QueryByID(ctx, menuItemID)
	if err != nil {
		return fmt.Errorf("querybyid: menuItemID[%s]: %w", menuItemID, err)
	}

	ub, err := toBusUpdateMenuItem(app)
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	rest, err := a.restaurantBus.QueryByID(ctx, item.RestaurantID)
	if err != nil {
		return errs.New(errs.InvalidArgument, fmt.Errorf("restaurant lookup: %w", err))
	}

	claims := mid.GetClaims(ctx)
	if !claims.IsOrgAuthorized(rest.OrganizationID) {
		return errs.Newf(errs.PermissionDenied, "user not in organization %s", rest.OrganizationID)
	}

	updItem, err := a.menuItemBus.Update(ctx, item, ub)
	if err != nil {
		return errs.Newf(errs.Internal, "update: menuItemID[%s] ub[%+v]: %s", menuItemID, ub, err)
	}

	return web.Respond(ctx, w, ToAppMenuItem(updItem), http.StatusOK)
}

// delete removes a menu item from the system.
func (a *app) delete(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	menuItemIDStr := web.Param(r, "menuitem_id")
	menuItemID, err := uuid.Parse(menuItemIDStr)
	if err != nil {
		return errs.NewFieldErrors("menuitem_id", err)
	}

	item, err := a.menuItemBus.QueryByID(ctx, menuItemID)
	if err != nil {
		return fmt.Errorf("querybyid: menuItemID[%s]: %w", menuItemID, err)
	}

	rest, err := a.restaurantBus.QueryByID(ctx, item.RestaurantID)
	if err != nil {
		return errs.New(errs.InvalidArgument, fmt.Errorf("restaurant lookup: %w", err))
	}

	claims := mid.GetClaims(ctx)
	if !claims.IsOrgAuthorized(rest.OrganizationID) {
		return errs.Newf(errs.PermissionDenied, "user not in organization %s", rest.OrganizationID)
	}

	if err := a.menuItemBus.Delete(ctx, item); err != nil {
		return errs.Newf(errs.Internal, "delete: menuItemID[%s]: %s", menuItemID, err)
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}

// reorder updates the display rank of menu items in a category.
func (a *app) reorder(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var app ReorderMenuItems
	if err := web.Decode(r, &app); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	categoryID, err := uuid.Parse(app.CategoryID)
	if err != nil {
		return errs.NewFieldErrors("categoryId", err)
	}

	cat, err := a.categoryBus.QueryByID(ctx, categoryID)
	if err != nil {
		return errs.New(errs.InvalidArgument, fmt.Errorf("category lookup: %w", err))
	}

	rest, err := a.restaurantBus.QueryByID(ctx, cat.RestaurantID)
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

	if err := a.menuItemBus.Reorder(ctx, categoryID, orderedIDs); err != nil {
		if errors.Is(err, menuitembus.ErrInvalidOrder) {
			return errs.New(errs.InvalidArgument, err)
		}
		return errs.Newf(errs.Internal, "reorder: %s", err)
	}

	items, err := a.menuItemBus.QueryByCategoryID(ctx, categoryID)
	if err != nil {
		return errs.Newf(errs.Internal, "query reordered menu items: categoryID[%s]: %s", categoryID, err)
	}

	return web.Respond(ctx, w, ToAppMenuItems(items), http.StatusOK)
}
