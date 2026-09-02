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

// create adds a new addon definition to the system.
func (a *app) create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var app NewAddon
	if err := web.Decode(r, &app); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	na, err := toBusNewAddon(app)
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
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

	return web.Respond(ctx, w, ToAppAddon(addon), http.StatusOK)
}

// update modifies an existing addon definition.
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
		return errs.Newf(errs.Internal, "update: addonID[%s] ua[%+v]: %s", addonID, ua, err)
	}

	return web.Respond(ctx, w, ToAppAddon(updAddon), http.StatusOK)
}

// delete removes an addon definition from the system.
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

// queryMenuItemAddons retrieves assigned addons for a menu item.
func (a *app) queryMenuItemAddons(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	menuItemIDStr := web.Param(r, "menu_item_id")
	menuItemID, err := uuid.Parse(menuItemIDStr)
	if err != nil {
		return errs.NewFieldErrors("menu_item_id", err)
	}

	infos, err := a.addonBus.QueryMenuItemAddons(ctx, menuItemID)
	if err != nil {
		return errs.Newf(errs.Internal, "queryMenuItemAddons: %s", err)
	}

	return web.Respond(ctx, w, ToAppMenuItemAddons(infos), http.StatusOK)
}

// replaceMenuItemAddons replaces the list of addons assigned to a menu item.
func (a *app) replaceMenuItemAddons(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	menuItemIDStr := web.Param(r, "menu_item_id")
	menuItemID, err := uuid.Parse(menuItemIDStr)
	if err != nil {
		return errs.NewFieldErrors("menu_item_id", err)
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

	var app ReplaceMenuItemAddons
	if err := web.Decode(r, &app); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	assignments := make([]addonbus.ItemAddonAssignment, len(app.Addons))
	for i, in := range app.Addons {
		addonID, err := uuid.Parse(in.AddonID)
		if err != nil {
			return errs.NewFieldErrors(fmt.Sprintf("addons[%d].addonId", i), err)
		}
		assignments[i] = addonbus.ItemAddonAssignment{
			AddonID: addonID,
			Rank:    in.Rank,
		}
	}

	assigned, err := a.addonBus.ReplaceMenuItemAddons(ctx, menuItemID, rest.ID, assignments)
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	return web.Respond(ctx, w, ToAppMenuItemAddons(assigned), http.StatusOK)
}

// reorderMenuItemAddons updates the rank of assigned addons on a menu item.
func (a *app) reorderMenuItemAddons(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	menuItemIDStr := web.Param(r, "menu_item_id")
	menuItemID, err := uuid.Parse(menuItemIDStr)
	if err != nil {
		return errs.NewFieldErrors("menu_item_id", err)
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

	var app ReorderMenuItemAddons
	if err := web.Decode(r, &app); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	orderedIDs := make([]uuid.UUID, len(app.OrderedIDs))
	for i, idStr := range app.OrderedIDs {
		id, err := uuid.Parse(idStr)
		if err != nil {
			return errs.NewFieldErrors("orderedIds", err)
		}
		orderedIDs[i] = id
	}

	reordered, err := a.addonBus.ReorderMenuItemAddons(ctx, menuItemID, orderedIDs)
	if err != nil {
		if errors.Is(err, addonbus.ErrInvalidReorder) {
			return errs.New(errs.InvalidArgument, err)
		}
		return errs.Newf(errs.Internal, "reorder: %s", err)
	}

	return web.Respond(ctx, w, ToAppMenuItemAddons(reordered), http.StatusOK)
}
