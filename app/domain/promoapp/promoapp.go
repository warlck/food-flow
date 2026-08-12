package promoapp

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/warlck/food-flow/app/sdk/errs"
	"github.com/warlck/food-flow/app/sdk/query"
	"github.com/warlck/food-flow/business/domain/promobus"
	"github.com/warlck/food-flow/business/sdk/order"
	"github.com/warlck/food-flow/business/sdk/page"
	"github.com/warlck/food-flow/foundation/web"
)

type app struct {
	promoBus *promobus.Business
}

func newApp(promoBus *promobus.Business) *app {
	return &app{
		promoBus: promoBus,
	}
}

// validate checks a promo code against restaurant and order subtotal.
func (a *app) validate(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var req ValidateRequest
	if err := web.Decode(r, &req); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	var restaurantID *uuid.UUID
	if req.RestaurantID != nil && *req.RestaurantID != "" {
		id, err := uuid.Parse(*req.RestaurantID)
		if err != nil {
			return errs.NewFieldErrors("restaurantId", err)
		}
		restaurantID = &id
	}

	res, err := a.promoBus.ValidatePromoCode(ctx, req.PromoCode, restaurantID, req.Subtotal)
	if err != nil {
		return errs.Newf(errs.Internal, "validate promo code: %s", err)
	}

	resp := ValidateResponse{
		Valid:          res.Valid,
		Reason:         res.Reason,
		Code:           res.Code,
		DiscountType:   res.DiscountType,
		DiscountValue:  res.DiscountValue,
		DiscountAmount: res.DiscountAmount,
		FinalSubtotal:  res.FinalSubtotal,
	}

	return web.Respond(ctx, w, resp, http.StatusOK)
}

// create adds a new promotion campaign.
func (a *app) create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var app NewPromotion
	if err := web.Decode(r, &app); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	np, err := toBusNewPromotion(app)
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	promo, err := a.promoBus.Create(ctx, np)
	if err != nil {
		return fmt.Errorf("create: promo[%+v]: %w", promo, err)
	}

	return web.Respond(ctx, w, ToAppPromotion(promo), http.StatusCreated)
}

// query retrieves a list of promotions.
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

	promos, err := a.promoBus.Query(ctx, filter, orderBy, pg)
	if err != nil {
		return errs.Newf(errs.Internal, "query: %s", err)
	}

	total, err := a.promoBus.Count(ctx, filter)
	if err != nil {
		return errs.Newf(errs.Internal, "count: %s", err)
	}

	result := query.NewResult(ToAppPromotions(promos), total, pg)
	return web.Respond(ctx, w, result, http.StatusOK)
}

// queryByID retrieves a promotion by its ID.
func (a *app) queryByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	promotionIDStr := web.Param(r, "promotion_id")

	promotionID, err := uuid.Parse(promotionIDStr)
	if err != nil {
		return errs.NewFieldErrors("promotion_id", err)
	}

	promo, err := a.promoBus.QueryByID(ctx, promotionID)
	if err != nil {
		if errors.Is(err, promobus.ErrNotFound) {
			return errs.New(errs.NotFound, err)
		}
		return fmt.Errorf("querybyid: promotionID[%s]: %w", promotionID, err)
	}

	return web.Respond(ctx, w, ToAppPromotion(promo), http.StatusOK)
}

// update modifies an existing promotion.
func (a *app) update(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	promotionIDStr := web.Param(r, "promotion_id")

	promotionID, err := uuid.Parse(promotionIDStr)
	if err != nil {
		return errs.NewFieldErrors("promotion_id", err)
	}

	var app UpdatePromotion
	if err := web.Decode(r, &app); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	up, err := toBusUpdatePromotion(app)
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	promo, err := a.promoBus.QueryByID(ctx, promotionID)
	if err != nil {
		if errors.Is(err, promobus.ErrNotFound) {
			return errs.New(errs.NotFound, err)
		}
		return fmt.Errorf("querybyid: promotionID[%s]: %w", promotionID, err)
	}

	updPromo, err := a.promoBus.Update(ctx, promo, up)
	if err != nil {
		return errs.Newf(errs.Internal, "update: promotionID[%s] up[%+v]: %s", promotionID, up, err)
	}

	return web.Respond(ctx, w, ToAppPromotion(updPromo), http.StatusOK)
}

// delete removes a promotion.
func (a *app) delete(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	promotionIDStr := web.Param(r, "promotion_id")

	promotionID, err := uuid.Parse(promotionIDStr)
	if err != nil {
		return errs.NewFieldErrors("promotion_id", err)
	}

	promo, err := a.promoBus.QueryByID(ctx, promotionID)
	if err != nil {
		if errors.Is(err, promobus.ErrNotFound) {
			return errs.New(errs.NotFound, err)
		}
		return fmt.Errorf("querybyid: promotionID[%s]: %w", promotionID, err)
	}

	if err := a.promoBus.Delete(ctx, promo); err != nil {
		return errs.Newf(errs.Internal, "delete: promotionID[%s]: %s", promotionID, err)
	}

	return web.Respond(ctx, w, struct{}{}, http.StatusNoContent)
}
