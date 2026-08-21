package orderapp

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/warlck/food-flow/app/sdk/errs"
	"github.com/warlck/food-flow/app/sdk/mid"
	"github.com/warlck/food-flow/app/sdk/query"
	"github.com/warlck/food-flow/business/domain/orderbus"
	"github.com/warlck/food-flow/business/domain/restaurantbus"
	"github.com/warlck/food-flow/business/sdk/order"
	"github.com/warlck/food-flow/business/sdk/page"
	"github.com/warlck/food-flow/foundation/web"
)

// app manages the set of order endpoints.
type app struct {
	orderBus      *orderbus.Business
	restaurantBus *restaurantbus.Business
}

// newApp constructs a handlers for route access.
func newApp(orderBus *orderbus.Business, restaurantBus *restaurantbus.Business) *app {
	return &app{
		orderBus:      orderBus,
		restaurantBus: restaurantBus,
	}
}

// create adds a new order to the system.
func (a *app) create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var app NewOrder
	if err := web.Decode(r, &app); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	nb, err := toBusNewOrder(app)
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	ord, err := a.orderBus.Create(ctx, nb)
	if err != nil {
		if errors.Is(err, orderbus.ErrMinSpendNotMet) || strings.Contains(err.Error(), "invalid promo code") {
			return errs.New(errs.InvalidArgument, err)
		}
		return fmt.Errorf("create: order[%+v]: %w", ord, err)
	}

	return web.Respond(ctx, w, ToAppOrder(ord), http.StatusCreated)
}

// query retrieves a list of orders based on query parameters.
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

	orders, err := a.orderBus.Query(ctx, filter, orderBy, pg)
	if err != nil {
		return errs.Newf(errs.Internal, "query: %s", err)
	}

	total, err := a.orderBus.Count(ctx, filter)
	if err != nil {
		return errs.Newf(errs.Internal, "count: %s", err)
	}

	result := query.NewResult(ToAppOrders(orders), total, pg)
	return web.Respond(ctx, w, result, http.StatusOK)
}

// queryByID retrieves an order by its ID.
func (a *app) queryByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	orderIDStr := web.Param(r, "order_id")

	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		return errs.NewFieldErrors("order_id", err)
	}

	ord, err := a.orderBus.QueryByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("querybyid: orderID[%s]: %w", orderID, err)
	}

	return web.Respond(ctx, w, ToAppOrder(ord), http.StatusOK)
}

// updateStatus updates the order or payment status.
func (a *app) updateStatus(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var app UpdateOrderStatus
	if err := web.Decode(r, &app); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	orderIDStr := web.Param(r, "order_id")
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		return errs.NewFieldErrors("order_id", err)
	}

	ub := toBusUpdateOrderStatus(app)

	ord, err := a.orderBus.QueryByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("querybyid: orderID[%s]: %w", orderID, err)
	}

	rest, err := a.restaurantBus.QueryByID(ctx, ord.RestaurantID)
	if err != nil {
		return errs.New(errs.InvalidArgument, fmt.Errorf("restaurant lookup: %w", err))
	}

	claims := mid.GetClaims(ctx)
	if !claims.IsOrgAuthorized(rest.OrganizationID) {
		return errs.Newf(errs.PermissionDenied, "user not in organization %s", rest.OrganizationID)
	}

	if err := a.orderBus.UpdateStatus(ctx, orderID, ub); err != nil {
		if errors.Is(err, orderbus.ErrOutForDeliveryRequiresDelivery) {
			return errs.New(errs.InvalidArgument, err)
		}
		return errs.Newf(errs.Internal, "updateStatus: orderID[%s] ub[%+v]: %s", orderID, ub, err)
	}

	// Query the updated order to return it
	ord, err = a.orderBus.QueryByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("querybyid: orderID[%s]: %w", orderID, err)
	}

	return web.Respond(ctx, w, ToAppOrder(ord), http.StatusOK)
}

// deliveryQuote calculates the delivery fee and range check for a destination.
func (a *app) deliveryQuote(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	restLat, err := strconv.ParseFloat(r.URL.Query().Get("restLat"), 64)
	if err != nil {
		return errs.NewFieldErrors("restLat", err)
	}

	restLng, err := strconv.ParseFloat(r.URL.Query().Get("restLng"), 64)
	if err != nil {
		return errs.NewFieldErrors("restLng", err)
	}

	lat, err := strconv.ParseFloat(r.URL.Query().Get("lat"), 64)
	if err != nil {
		return errs.NewFieldErrors("lat", err)
	}

	lng, err := strconv.ParseFloat(r.URL.Query().Get("lng"), 64)
	if err != nil {
		return errs.NewFieldErrors("lng", err)
	}

	var maxDeliveryDistanceKm float64
	if maxStr := r.URL.Query().Get("maxDeliveryDistanceKm"); maxStr != "" {
		maxDeliveryDistanceKm, err = strconv.ParseFloat(maxStr, 64)
		if err != nil {
			return errs.NewFieldErrors("maxDeliveryDistanceKm", err)
		}
	}

	quote, err := a.orderBus.DeliveryQuote(restLat, restLng, lat, lng, maxDeliveryDistanceKm)
	if err != nil {
		return fmt.Errorf("deliveryquote: %w", err)
	}

	return web.Respond(ctx, w, ToAppDeliveryQuote(quote), http.StatusOK)
}

// cancel cancels an order.
func (a *app) cancel(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	orderIDStr := web.Param(r, "order_id")
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		return errs.NewFieldErrors("order_id", err)
	}

	ord, err := a.orderBus.QueryByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("querybyid: orderID[%s]: %w", orderID, err)
	}

	rest, err := a.restaurantBus.QueryByID(ctx, ord.RestaurantID)
	if err != nil {
		return errs.New(errs.InvalidArgument, fmt.Errorf("restaurant lookup: %w", err))
	}

	claims := mid.GetClaims(ctx)
	fmt.Printf("DEBUG: ordID=%s restID=%s rest.OrgID=%s claims.OrgIDs=%v\n", orderID, rest.ID, rest.OrganizationID, claims.OrganizationIDs)
	if !claims.IsOrgAuthorized(rest.OrganizationID) {
		return errs.Newf(errs.PermissionDenied, "user not in organization %s", rest.OrganizationID)
	}

	if err := a.orderBus.Cancel(ctx, orderID); err != nil {
		return errs.Newf(errs.Internal, "cancel: orderID[%s]: %s", orderID, err)
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}
