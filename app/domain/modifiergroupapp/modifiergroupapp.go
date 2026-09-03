package modifiergroupapp

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/warlck/food-flow/app/sdk/errs"
	"github.com/warlck/food-flow/app/sdk/mid"
	"github.com/warlck/food-flow/app/sdk/query"
	"github.com/warlck/food-flow/business/domain/menuitembus"
	"github.com/warlck/food-flow/business/domain/modifiergroupbus"
	"github.com/warlck/food-flow/business/domain/restaurantbus"
	"github.com/warlck/food-flow/business/sdk/order"
	"github.com/warlck/food-flow/business/sdk/page"
	"github.com/warlck/food-flow/foundation/web"
)

// Handlers manages the set of modifier group endpoints.
type app struct {
	modifierGroupBus *modifiergroupbus.Business
	menuItemBus      *menuitembus.Business
	restaurantBus    *restaurantbus.Business
}

// newApp constructs a handlers for route access.
func newApp(modifierGroupBus *modifiergroupbus.Business, menuItemBus *menuitembus.Business, restaurantBus *restaurantbus.Business) *app {
	return &app{
		modifierGroupBus: modifierGroupBus,
		menuItemBus:      menuItemBus,
		restaurantBus:    restaurantBus,
	}
}

// Create adds a new modifier group to the system.
func (a *app) create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var app NewModifierGroup
	if err := web.Decode(r, &app); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	ng, err := toBusNewModifierGroup(app)
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	rest, err := a.restaurantBus.QueryByID(ctx, ng.RestaurantID)
	if err != nil {
		return errs.New(errs.InvalidArgument, fmt.Errorf("restaurant lookup: %w", err))
	}

	claims := mid.GetClaims(ctx)
	if !claims.IsOrgAuthorized(rest.OrganizationID) {
		return errs.Newf(errs.PermissionDenied, "user not in organization %s", rest.OrganizationID)
	}

	item, err := a.menuItemBus.QueryByID(ctx, ng.MenuItemID)
	if err != nil {
		return errs.New(errs.InvalidArgument, fmt.Errorf("menu item lookup: %w", err))
	}
	if item.RestaurantID != rest.ID {
		return errs.Newf(errs.InvalidArgument, "menu item %s does not belong to restaurant %s", item.ID, rest.ID)
	}

	group, err := a.modifierGroupBus.Create(ctx, ng)
	if err != nil {
		if errors.Is(err, modifiergroupbus.ErrRequiredNoOptions) {
			return errs.New(errs.InvalidArgument, err)
		}
		return fmt.Errorf("create: modifierGroup[%+v]: %w", group, err)
	}

	return web.Respond(ctx, w, ToAppModifierGroup(group), http.StatusCreated)
}

// Reorder updates the rank of all modifier groups within a menu item.
func (a *app) reorder(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var app ReorderModifierGroups
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

	orderedIDs := make([]uuid.UUID, len(app.OrderedIDs))
	for i, idStr := range app.OrderedIDs {
		id, err := uuid.Parse(idStr)
		if err != nil {
			return errs.NewFieldErrors("orderedIds", err)
		}
		orderedIDs[i] = id
	}

	reordered, err := a.modifierGroupBus.Reorder(ctx, menuItemID, orderedIDs)
	if err != nil {
		if errors.Is(err, modifiergroupbus.ErrInvalidReorder) {
			return errs.New(errs.InvalidArgument, err)
		}
		return errs.Newf(errs.Internal, "reorder: %s", err)
	}

	return web.Respond(ctx, w, ToAppModifierGroups(reordered), http.StatusOK)
}

// Query retrieves a list of modifier groups based on query parameters.
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

	groups, err := a.modifierGroupBus.Query(ctx, filter, orderBy, pg)
	if err != nil {
		return errs.Newf(errs.Internal, "query: %s", err)
	}

	total, err := a.modifierGroupBus.Count(ctx, filter)
	if err != nil {
		return errs.Newf(errs.Internal, "count: %s", err)
	}

	result := query.NewResult(ToAppModifierGroups(groups), total, pg)
	return web.Respond(ctx, w, result, http.StatusOK)
}

// authorizeScope verifies the caller's organization owns the restaurant the
// list query is scoped to. A scoping filter (restaurant_id or menu_item_id)
// is required so list reads can never enumerate another organization's
// catalog; a supplied filter is never treated as authorization by itself.
func (a *app) authorizeScope(ctx context.Context, filter modifiergroupbus.QueryFilter) error {
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

// QueryByID retrieves a modifier group by its ID.
func (a *app) queryByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	groupIDStr := web.Param(r, "modifier_group_id")

	groupID, err := uuid.Parse(groupIDStr)
	if err != nil {
		return errs.NewFieldErrors("modifier_group_id", err)
	}

	group, err := a.modifierGroupBus.QueryByID(ctx, groupID)
	if err != nil {
		if errors.Is(err, modifiergroupbus.ErrNotFound) {
			return errs.New(errs.NotFound, err)
		}
		return fmt.Errorf("querybyid: groupID[%s]: %w", groupID, err)
	}

	rest, err := a.restaurantBus.QueryByID(ctx, group.RestaurantID)
	if err != nil {
		return errs.New(errs.InvalidArgument, fmt.Errorf("restaurant lookup: %w", err))
	}

	claims := mid.GetClaims(ctx)
	if !claims.IsOrgAuthorized(rest.OrganizationID) {
		return errs.Newf(errs.PermissionDenied, "user not in organization %s", rest.OrganizationID)
	}

	return web.Respond(ctx, w, ToAppModifierGroup(group), http.StatusOK)
}

// Update modifies a modifier group.
func (a *app) update(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	groupIDStr := web.Param(r, "modifier_group_id")

	groupID, err := uuid.Parse(groupIDStr)
	if err != nil {
		return errs.NewFieldErrors("modifier_group_id", err)
	}

	var app UpdateModifierGroup
	if err := web.Decode(r, &app); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	ug, err := toBusUpdateModifierGroup(app)
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	group, err := a.modifierGroupBus.QueryByID(ctx, groupID)
	if err != nil {
		if errors.Is(err, modifiergroupbus.ErrNotFound) {
			return errs.New(errs.NotFound, err)
		}
		return fmt.Errorf("querybyid: groupID[%s]: %w", groupID, err)
	}

	rest, err := a.restaurantBus.QueryByID(ctx, group.RestaurantID)
	if err != nil {
		return errs.New(errs.InvalidArgument, fmt.Errorf("restaurant lookup: %w", err))
	}

	claims := mid.GetClaims(ctx)
	if !claims.IsOrgAuthorized(rest.OrganizationID) {
		return errs.Newf(errs.PermissionDenied, "user not in organization %s", rest.OrganizationID)
	}

	updGroup, err := a.modifierGroupBus.Update(ctx, group, ug)
	if err != nil {
		if errors.Is(err, modifiergroupbus.ErrRequiredNoOptions) {
			return errs.New(errs.InvalidArgument, err)
		}
		return errs.Newf(errs.Internal, "update: groupID[%s] ug[%+v]: %s", groupID, ug, err)
	}

	return web.Respond(ctx, w, ToAppModifierGroup(updGroup), http.StatusOK)
}

// Delete removes a modifier group.
func (a *app) delete(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	groupIDStr := web.Param(r, "modifier_group_id")

	groupID, err := uuid.Parse(groupIDStr)
	if err != nil {
		return errs.NewFieldErrors("modifier_group_id", err)
	}

	group, err := a.modifierGroupBus.QueryByID(ctx, groupID)
	if err != nil {
		if errors.Is(err, modifiergroupbus.ErrNotFound) {
			return errs.New(errs.NotFound, err)
		}
		return fmt.Errorf("querybyid: groupID[%s]: %w", groupID, err)
	}

	rest, err := a.restaurantBus.QueryByID(ctx, group.RestaurantID)
	if err != nil {
		return errs.New(errs.InvalidArgument, fmt.Errorf("restaurant lookup: %w", err))
	}

	claims := mid.GetClaims(ctx)
	if !claims.IsOrgAuthorized(rest.OrganizationID) {
		return errs.Newf(errs.PermissionDenied, "user not in organization %s", rest.OrganizationID)
	}

	if err := a.modifierGroupBus.Delete(ctx, group); err != nil {
		return errs.Newf(errs.Internal, "delete: groupID[%s]: %s", groupID, err)
	}

	return web.Respond(ctx, w, struct{}{}, http.StatusNoContent)
}
