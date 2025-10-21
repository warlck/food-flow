package restaurantapp

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/warlck/food-flow/app/sdk/errs"
	"github.com/warlck/food-flow/app/sdk/query"
	"github.com/warlck/food-flow/business/domain/restaurantbus"
	"github.com/warlck/food-flow/business/sdk/order"
	"github.com/warlck/food-flow/business/sdk/page"
	"github.com/warlck/food-flow/foundation/web"
)

// Handlers manages the set of restaurant endpoints.
type app struct {
	restaurantBus *restaurantbus.Business
}

// newApp constructs a handlers for route access.
func newApp(restaurantBus *restaurantbus.Business) *app {
	return &app{
		restaurantBus: restaurantBus,
	}
}

// Create adds a new restaurant to the system.
func (a *app) create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var app NewRestaurant
	if err := web.Decode(r, &app); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	nr, err := toBusNewRestaurant(app)
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	res, err := a.restaurantBus.Create(ctx, nr)
	if err != nil {
		return fmt.Errorf("create: restaurant[%+v]: %w", res, err)
	}

	return web.Respond(ctx, w, toAppRestaurant(res), http.StatusCreated)
}

// Query retrieves a list of restaurants based on query parameters.
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

	restaurants, err := a.restaurantBus.Query(ctx, filter, orderBy, pg)
	if err != nil {
		return errs.Newf(errs.Internal, "query: %s", err)
	}

	total, err := a.restaurantBus.Count(ctx, filter)
	if err != nil {
		return errs.Newf(errs.Internal, "count: %s", err)
	}

	result := query.NewResult(toAppRestaurants(restaurants), total, pg)
	return web.Respond(ctx, w, result, http.StatusOK)
}

// QueryByID retrieves a restaurant by its ID.
func (a *app) queryByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	restaurantIDStr := web.Param(r, "restaurant_id")

	restaurantID, err := uuid.Parse(restaurantIDStr)
	if err != nil {
		return errs.NewFieldErrors("restaurant_id", err)
	}

	res, err := a.restaurantBus.QueryByID(ctx, restaurantID)
	if err != nil {
		return fmt.Errorf("querybyid: restaurantID[%s]: %w", restaurantID, err)
	}

	return web.Respond(ctx, w, toAppRestaurant(res), http.StatusOK)
}

// Update modifies a restaurant.
func (a *app) update(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	restaurantIDStr := web.Param(r, "restaurant_id")

	restaurantID, err := uuid.Parse(restaurantIDStr)
	if err != nil {
		return errs.NewFieldErrors("restaurant_id", err)
	}

	var app UpdateRestaurant
	if err := web.Decode(r, &app); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	ur, err := toBusUpdateRestaurant(app)
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	res, err := a.restaurantBus.QueryByID(ctx, restaurantID)
	if err != nil {
		return fmt.Errorf("querybyid: restaurantID[%s]: %w", restaurantID, err)
	}

	updRes, err := a.restaurantBus.Update(ctx, res, ur)
	if err != nil {
		return errs.Newf(errs.Internal, "update: restaurantID[%s] ur[%+v]: %s", restaurantID, ur, err)
	}

	return web.Respond(ctx, w, toAppRestaurant(updRes), http.StatusOK)
}

// Delete removes a restaurant.
func (a *app) delete(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	restaurantIDStr := web.Param(r, "restaurant_id")

	restaurantID, err := uuid.Parse(restaurantIDStr)
	if err != nil {
		return errs.NewFieldErrors("restaurant_id", err)
	}

	res, err := a.restaurantBus.QueryByID(ctx, restaurantID)
	if err != nil {
		return fmt.Errorf("querybyid: restaurantID[%s]: %w", restaurantID, err)
	}

	if err := a.restaurantBus.Delete(ctx, res); err != nil {
		return errs.Newf(errs.Internal, "delete: restaurantID[%s]: %s", restaurantID, err)
	}

	return web.Respond(ctx, w, struct{}{}, http.StatusNoContent)
}
