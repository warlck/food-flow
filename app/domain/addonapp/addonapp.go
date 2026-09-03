package addonapp

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/warlck/food-flow/app/sdk/errs"
	"github.com/warlck/food-flow/app/sdk/mid"
	"github.com/warlck/food-flow/app/sdk/query"
	"github.com/warlck/food-flow/business/domain/addonbus"
	"github.com/warlck/food-flow/business/domain/menuitembus"
	"github.com/warlck/food-flow/business/domain/restaurantbus"
	"github.com/warlck/food-flow/business/sdk/order"
	"github.com/warlck/food-flow/business/sdk/page"
	"github.com/warlck/food-flow/foundation/web"
)

// app manages the set of addon endpoints.
type app struct {
	addonBus      *addonbus.Business
	menuItemBus   *menuitembus.Business
	restaurantBus *restaurantbus.Business
}

// newApp constructs an app handler for route access.
func newApp(addonBus *addonbus.Business, menuItemBus *menuitembus.Business, restaurantBus *restaurantbus.Business) *app {
	return &app{
		addonBus:      addonBus,
		menuItemBus:   menuItemBus,
		restaurantBus: restaurantBus,
	}
}

// create adds a new addon to a menu item.
func (a *app) create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var app NewAddon
	if err := web.Decode(r, &app); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	na, err := toBusNewAddon(app)
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	item, err := a.menuItemBus.QueryByID(ctx, na.MenuItemID)
	if err != nil {
		return errs.New(errs.InvalidArgument, fmt.Errorf("menu item lookup: %w", err))
	}

	if item.RestaurantID != na.RestaurantID {
		return errs.Newf(errs.InvalidArgument, "menu item %s does not belong to restaurant %s", na.MenuItemID, na.RestaurantID)
	}

	rest, err := a.restaurantBus.QueryByID(ctx, na.RestaurantID)
	if err != nil {
		return errs.New(errs.InvalidArgument, fmt.Errorf("restaurant lookup: %w", err))
	}

	claims := mid.GetClaims(ctx)
	if !claims.IsOrgAuthorized(rest.OrganizationID) {
		return errs.Newf(errs.PermissionDenied, "user not in organization %s", rest.OrganizationID)
	}

	addon, err := a.addonBus.Create(ctx, na)
	if err != nil {
		if errors.Is(err, addonbus.ErrDuplicateName) {
			return errs.New(errs.InvalidArgument, err)
		}
		return fmt.Errorf("create: addon[%+v]: %w", addon, err)
	}

	return web.Respond(ctx, w, ToAppAddon(addon), http.StatusCreated)
}

// query retrieves a list of addons based on query parameters.
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

	if err := a.authorizeScope(ctx, filter); err != nil {
		return err
	}

	orderBy, err := order.Parse(orderByFields, qp.OrderBy, defaultOrderBy)
	if err != nil {
		return errs.NewFieldErrors("order", err)
	}

	addons, err := a.addonBus.Query(ctx, filter, orderBy, pg)
	if err != nil {
		return errs.Newf(errs.Internal, "query: %s", err)
	}

	total, err := a.addonBus.Count(ctx, filter)
	if err != nil {
		return errs.Newf(errs.Internal, "count: %s", err)
	}

	result := query.NewResult(ToAppAddons(addons), total, pg)
	return web.Respond(ctx, w, result, http.StatusOK)
}

// authorizeScope verifies the caller's organization owns the restaurant the
// list query is scoped to. A scoping filter (restaurant_id or menu_item_id)
// is required so list reads can never enumerate another organization's
// catalog; a supplied filter is never treated as authorization by itself.
func (a *app) authorizeScope(ctx context.Context, filter addonbus.QueryFilter) error {
	var restaurantID uuid.UUID

	switch {
	case filter.RestaurantID != nil:
		restaurantID = *filter.RestaurantID
	case filter.MenuItemID != nil:
		item, err := a.menuItemBus.QueryByID(ctx, *filter.MenuItemID)
		if err != nil {
			return errs.New(errs.InvalidArgument, fmt.Errorf("menu item lookup: %w", err))
		}
		restaurantID = item.RestaurantID
	default:
		return errs.Newf(errs.InvalidArgument, "query requires a restaurant_id or menu_item_id filter")
	}

	rest, err := a.restaurantBus.QueryByID(ctx, restaurantID)
	if err != nil {
		return errs.New(errs.InvalidArgument, fmt.Errorf("restaurant lookup: %w", err))
	}

	claims := mid.GetClaims(ctx)
	if !claims.IsOrgAuthorized(rest.OrganizationID) {
		return errs.Newf(errs.PermissionDenied, "user not in organization %s", rest.OrganizationID)
	}

	return nil
}

// queryByID retrieves an addon by its ID.
func (a *app) queryByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	addonIDStr := web.Param(r, "addon_id")

	addonID, err := uuid.Parse(addonIDStr)
	if err != nil {
		return errs.NewFieldErrors("addon_id", err)
	}

	addon, err := a.addonBus.QueryByID(ctx, addonID)
	if err != nil {
		if errors.Is(err, addonbus.ErrNotFound) {
			return errs.New(errs.NotFound, err)
		}
		return fmt.Errorf("querybyid: addonID[%s]: %w", addonID, err)
	}

	rest, err := a.restaurantBus.QueryByID(ctx, addon.RestaurantID)
	if err != nil {
		return errs.New(errs.InvalidArgument, fmt.Errorf("restaurant lookup: %w", err))
	}

	claims := mid.GetClaims(ctx)
	if !claims.IsOrgAuthorized(rest.OrganizationID) {
		return errs.Newf(errs.PermissionDenied, "user not in organization %s", rest.OrganizationID)
	}

	return web.Respond(ctx, w, ToAppAddon(addon), http.StatusOK)
}

// update modifies an existing addon.
func (a *app) update(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var app UpdateAddon
	if err := web.Decode(r, &app); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	addonIDStr := web.Param(r, "addon_id")
	addonID, err := uuid.Parse(addonIDStr)
	if err != nil {
		return errs.NewFieldErrors("addon_id", err)
	}

	addon, err := a.addonBus.QueryByID(ctx, addonID)
	if err != nil {
		if errors.Is(err, addonbus.ErrNotFound) {
			return errs.New(errs.NotFound, err)
		}
		return fmt.Errorf("querybyid: addonID[%s]: %w", addonID, err)
	}

	ua, err := toBusUpdateAddon(app)
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	rest, err := a.restaurantBus.QueryByID(ctx, addon.RestaurantID)
	if err != nil {
		return errs.New(errs.InvalidArgument, fmt.Errorf("restaurant lookup: %w", err))
	}

	claims := mid.GetClaims(ctx)
	if !claims.IsOrgAuthorized(rest.OrganizationID) {
		return errs.Newf(errs.PermissionDenied, "user not in organization %s", rest.OrganizationID)
	}

	updAddon, err := a.addonBus.Update(ctx, addon, ua)
	if err != nil {
		if errors.Is(err, addonbus.ErrDuplicateName) {
			return errs.New(errs.InvalidArgument, err)
		}
		return errs.Newf(errs.Internal, "update: addonID[%s] ua[%+v]: %s", addonID, ua, err)
	}

	return web.Respond(ctx, w, ToAppAddon(updAddon), http.StatusOK)
}

// delete removes an addon from the system.
func (a *app) delete(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	addonIDStr := web.Param(r, "addon_id")
	addonID, err := uuid.Parse(addonIDStr)
	if err != nil {
		return errs.NewFieldErrors("addon_id", err)
	}

	addon, err := a.addonBus.QueryByID(ctx, addonID)
	if err != nil {
		if errors.Is(err, addonbus.ErrNotFound) {
			return errs.New(errs.NotFound, err)
		}
		return fmt.Errorf("querybyid: addonID[%s]: %w", addonID, err)
	}

	rest, err := a.restaurantBus.QueryByID(ctx, addon.RestaurantID)
	if err != nil {
		return errs.New(errs.InvalidArgument, fmt.Errorf("restaurant lookup: %w", err))
	}

	claims := mid.GetClaims(ctx)
	if !claims.IsOrgAuthorized(rest.OrganizationID) {
		return errs.Newf(errs.PermissionDenied, "user not in organization %s", rest.OrganizationID)
	}

	if err := a.addonBus.Delete(ctx, addon); err != nil {
		return errs.Newf(errs.Internal, "delete: addonID[%s]: %s", addonID, err)
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}

// reorder updates the display rank of addons on a menu item.
func (a *app) reorder(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var app ReorderAddons
	if err := web.Decode(r, &app); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	menuItemID, err := uuid.Parse(app.MenuItemID)
	if err != nil {
		return errs.NewFieldErrors("menuItemId", err)
	}

	item, err := a.menuItemBus.QueryByID(ctx, menuItemID)
	if err != nil {
		return errs.New(errs.InvalidArgument, fmt.Errorf("menu item lookup: %w", err))
	}

	rest, err := a.restaurantBus.QueryByID(ctx, item.RestaurantID)
	if err != nil {
		return errs.New(errs.InvalidArgument, fmt.Errorf("restaurant lookup: %w", err))
	}

	claims := mid.GetClaims(ctx)
	if !claims.IsOrgAuthorized(rest.OrganizationID) {
		return errs.Newf(errs.PermissionDenied, "user not in organization %s", rest.OrganizationID)
	}

	orderedIDs := make([]uuid.UUID, len(app.AddonIDs))
	for i, idStr := range app.AddonIDs {
		id, err := uuid.Parse(idStr)
		if err != nil {
			return errs.NewFieldErrors(fmt.Sprintf("addonIds[%d]", i), err)
		}
		orderedIDs[i] = id
	}

	reordered, err := a.addonBus.Reorder(ctx, menuItemID, orderedIDs)
	if err != nil {
		if errors.Is(err, addonbus.ErrInvalidReorder) {
			return errs.New(errs.InvalidArgument, err)
		}
		return errs.Newf(errs.Internal, "reorder: %s", err)
	}

	return web.Respond(ctx, w, ToAppAddons(reordered), http.StatusOK)
}
