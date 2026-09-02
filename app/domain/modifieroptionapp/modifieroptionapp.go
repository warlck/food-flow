package modifieroptionapp

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/warlck/food-flow/app/sdk/errs"
	"github.com/warlck/food-flow/app/sdk/mid"
	"github.com/warlck/food-flow/app/sdk/query"
	"github.com/warlck/food-flow/business/domain/modifiergroupbus"
	"github.com/warlck/food-flow/business/domain/modifieroptionbus"
	"github.com/warlck/food-flow/business/domain/restaurantbus"
	"github.com/warlck/food-flow/business/sdk/order"
	"github.com/warlck/food-flow/business/sdk/page"
	"github.com/warlck/food-flow/foundation/web"
)

// Handlers manages the set of modifier option endpoints.
type app struct {
	modifierOptionBus *modifieroptionbus.Business
	modifierGroupBus  *modifiergroupbus.Business
	restaurantBus     *restaurantbus.Business
}

// newApp constructs a handlers for route access.
func newApp(modifierOptionBus *modifieroptionbus.Business, modifierGroupBus *modifiergroupbus.Business, restaurantBus *restaurantbus.Business) *app {
	return &app{
		modifierOptionBus: modifierOptionBus,
		modifierGroupBus:  modifierGroupBus,
		restaurantBus:     restaurantBus,
	}
}

// Create adds a new modifier option to the system.
func (a *app) create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var app NewModifierOption
	if err := web.Decode(r, &app); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	no, err := toBusNewModifierOption(app)
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	rest, err := a.restaurantBus.QueryByID(ctx, no.RestaurantID)
	if err != nil {
		return errs.New(errs.InvalidArgument, fmt.Errorf("restaurant lookup: %w", err))
	}

	claims := mid.GetClaims(ctx)
	if !claims.IsOrgAuthorized(rest.OrganizationID) {
		return errs.Newf(errs.PermissionDenied, "user not in organization %s", rest.OrganizationID)
	}

	group, err := a.modifierGroupBus.QueryByID(ctx, no.ModifierGroupID)
	if err != nil {
		return errs.New(errs.InvalidArgument, fmt.Errorf("modifier group lookup: %w", err))
	}
	if group.RestaurantID != rest.ID {
		return errs.Newf(errs.InvalidArgument, "modifier group %s does not belong to restaurant %s", group.ID, rest.ID)
	}

	opt, err := a.modifierOptionBus.Create(ctx, no)
	if err != nil {
		return fmt.Errorf("create: modifierOption[%+v]: %w", opt, err)
	}

	return web.Respond(ctx, w, ToAppModifierOption(opt), http.StatusCreated)
}

// Reorder updates the rank of all modifier options within a group.
func (a *app) reorder(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var app ReorderModifierOptions
	if err := web.Decode(r, &app); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	groupID, err := uuid.Parse(app.ModifierGroupID)
	if err != nil {
		return errs.NewFieldErrors("modifierGroupId", err)
	}

	group, err := a.modifierGroupBus.QueryByID(ctx, groupID)
	if err != nil {
		return errs.New(errs.InvalidArgument, fmt.Errorf("modifier group lookup: %w", err))
	}

	rest, err := a.restaurantBus.QueryByID(ctx, group.RestaurantID)
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

	reordered, err := a.modifierOptionBus.Reorder(ctx, groupID, orderedIDs)
	if err != nil {
		if errors.Is(err, modifieroptionbus.ErrInvalidReorder) {
			return errs.New(errs.InvalidArgument, err)
		}
		return errs.Newf(errs.Internal, "reorder: %s", err)
	}

	return web.Respond(ctx, w, ToAppModifierOptions(reordered), http.StatusOK)
}

// Query retrieves a list of modifier options based on query parameters.
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

	options, err := a.modifierOptionBus.Query(ctx, filter, orderBy, pg)
	if err != nil {
		return errs.Newf(errs.Internal, "query: %s", err)
	}

	total, err := a.modifierOptionBus.Count(ctx, filter)
	if err != nil {
		return errs.Newf(errs.Internal, "count: %s", err)
	}

	result := query.NewResult(ToAppModifierOptions(options), total, pg)
	return web.Respond(ctx, w, result, http.StatusOK)
}

// QueryByID retrieves a modifier option by its ID.
func (a *app) queryByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	optionIDStr := web.Param(r, "modifier_option_id")

	optionID, err := uuid.Parse(optionIDStr)
	if err != nil {
		return errs.NewFieldErrors("modifier_option_id", err)
	}

	opt, err := a.modifierOptionBus.QueryByID(ctx, optionID)
	if err != nil {
		if errors.Is(err, modifieroptionbus.ErrNotFound) {
			return errs.New(errs.NotFound, err)
		}
		return fmt.Errorf("querybyid: optionID[%s]: %w", optionID, err)
	}

	return web.Respond(ctx, w, ToAppModifierOption(opt), http.StatusOK)
}

// Update modifies a modifier option.
func (a *app) update(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	optionIDStr := web.Param(r, "modifier_option_id")

	optionID, err := uuid.Parse(optionIDStr)
	if err != nil {
		return errs.NewFieldErrors("modifier_option_id", err)
	}

	var app UpdateModifierOption
	if err := web.Decode(r, &app); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	uo, err := toBusUpdateModifierOption(app)
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	opt, err := a.modifierOptionBus.QueryByID(ctx, optionID)
	if err != nil {
		if errors.Is(err, modifieroptionbus.ErrNotFound) {
			return errs.New(errs.NotFound, err)
		}
		return fmt.Errorf("querybyid: optionID[%s]: %w", optionID, err)
	}

	rest, err := a.restaurantBus.QueryByID(ctx, opt.RestaurantID)
	if err != nil {
		return errs.New(errs.InvalidArgument, fmt.Errorf("restaurant lookup: %w", err))
	}

	claims := mid.GetClaims(ctx)
	if !claims.IsOrgAuthorized(rest.OrganizationID) {
		return errs.Newf(errs.PermissionDenied, "user not in organization %s", rest.OrganizationID)
	}

	updOpt, err := a.modifierOptionBus.Update(ctx, opt, uo)
	if err != nil {
		return errs.Newf(errs.Internal, "update: optionID[%s] uo[%+v]: %s", optionID, uo, err)
	}

	return web.Respond(ctx, w, ToAppModifierOption(updOpt), http.StatusOK)
}

// Delete removes a modifier option.
func (a *app) delete(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	optionIDStr := web.Param(r, "modifier_option_id")

	optionID, err := uuid.Parse(optionIDStr)
	if err != nil {
		return errs.NewFieldErrors("modifier_option_id", err)
	}

	opt, err := a.modifierOptionBus.QueryByID(ctx, optionID)
	if err != nil {
		if errors.Is(err, modifieroptionbus.ErrNotFound) {
			return errs.New(errs.NotFound, err)
		}
		return fmt.Errorf("querybyid: optionID[%s]: %w", optionID, err)
	}

	rest, err := a.restaurantBus.QueryByID(ctx, opt.RestaurantID)
	if err != nil {
		return errs.New(errs.InvalidArgument, fmt.Errorf("restaurant lookup: %w", err))
	}

	claims := mid.GetClaims(ctx)
	if !claims.IsOrgAuthorized(rest.OrganizationID) {
		return errs.Newf(errs.PermissionDenied, "user not in organization %s", rest.OrganizationID)
	}

	if err := a.modifierOptionBus.Delete(ctx, opt); err != nil {
		return errs.Newf(errs.Internal, "delete: optionID[%s]: %s", optionID, err)
	}

	return web.Respond(ctx, w, struct{}{}, http.StatusNoContent)
}
