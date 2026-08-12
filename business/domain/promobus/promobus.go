package promobus

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/warlck/food-flow/business/sdk/order"
	"github.com/warlck/food-flow/business/sdk/page"
	"github.com/warlck/food-flow/business/types/name"
	"github.com/warlck/food-flow/foundation/logger"
)

var (
	ErrNotFound     = errors.New("promotion not found")
	ErrInvalid      = errors.New("invalid promotion data")
	ErrLimitReached = errors.New("promo code usage limit reached")
)

type Storer interface {
	Create(ctx context.Context, promo Promotion) error
	Update(ctx context.Context, promo Promotion) error
	Delete(ctx context.Context, promo Promotion) error
	Query(ctx context.Context, filter QueryFilter, orderBy order.By, page page.Page) ([]Promotion, error)
	Count(ctx context.Context, filter QueryFilter) (int, error)
	QueryByID(ctx context.Context, promotionID uuid.UUID) (Promotion, error)
	QueryByCode(ctx context.Context, code string) (Promotion, error)
	IncrementUsage(ctx context.Context, promotionID uuid.UUID) error
}

type Business struct {
	log    *logger.Logger
	storer Storer
}

func NewBusiness(log *logger.Logger, storer Storer) *Business {
	return &Business{
		log:    log,
		storer: storer,
	}
}

func (b *Business) Create(ctx context.Context, np NewPromotion) (Promotion, error) {
	code := strings.ToUpper(strings.TrimSpace(np.Code))
	if code == "" {
		return Promotion{}, fmt.Errorf("%w: code is required", ErrInvalid)
	}

	if np.DiscountType != DiscountTypePercentage && np.DiscountType != DiscountTypeFixedAmount {
		return Promotion{}, fmt.Errorf("%w: invalid discount_type %q", ErrInvalid, np.DiscountType)
	}

	if np.DiscountValue <= 0 {
		return Promotion{}, fmt.Errorf("%w: discount_value must be greater than 0", ErrInvalid)
	}

	if np.DiscountType == DiscountTypePercentage && np.DiscountValue > 100 {
		return Promotion{}, fmt.Errorf("%w: percentage discount value cannot exceed 100", ErrInvalid)
	}

	now := time.Now()
	promo := Promotion{
		ID:                uuid.New(),
		RestaurantID:      np.RestaurantID,
		Code:              code,
		Name:              np.Name,
		Description:       np.Description,
		DiscountType:      np.DiscountType,
		DiscountValue:     np.DiscountValue,
		MinOrderAmount:    np.MinOrderAmount,
		MaxDiscountAmount: np.MaxDiscountAmount,
		UsageLimit:        np.UsageLimit,
		UsageCount:        0,
		StartDate:         np.StartDate,
		EndDate:           np.EndDate,
		Enabled:           np.Enabled,
		DateCreated:       now,
		DateUpdated:       now,
	}

	if err := b.storer.Create(ctx, promo); err != nil {
		return Promotion{}, fmt.Errorf("create: %w", err)
	}

	return promo, nil
}

func (b *Business) Update(ctx context.Context, promo Promotion, up UpdatePromotion) (Promotion, error) {
	if up.Code != nil {
		code := strings.ToUpper(strings.TrimSpace(*up.Code))
		if code == "" {
			return Promotion{}, fmt.Errorf("%w: code cannot be empty", ErrInvalid)
		}
		promo.Code = code
	}
	if up.Name != nil {
		promo.Name = *up.Name
	}
	if up.Description != nil {
		promo.Description = *up.Description
	}
	if up.DiscountType != nil {
		if *up.DiscountType != DiscountTypePercentage && *up.DiscountType != DiscountTypeFixedAmount {
			return Promotion{}, fmt.Errorf("%w: invalid discount_type %q", ErrInvalid, *up.DiscountType)
		}
		promo.DiscountType = *up.DiscountType
	}
	if up.DiscountValue != nil {
		if *up.DiscountValue <= 0 {
			return Promotion{}, fmt.Errorf("%w: discount_value must be greater than 0", ErrInvalid)
		}
		promo.DiscountValue = *up.DiscountValue
	}

	if promo.DiscountType == DiscountTypePercentage && promo.DiscountValue > 100 {
		return Promotion{}, fmt.Errorf("%w: percentage discount value cannot exceed 100", ErrInvalid)
	}
	if up.MinOrderAmount != nil {
		promo.MinOrderAmount = *up.MinOrderAmount
	}
	if up.Enabled != nil {
		promo.Enabled = *up.Enabled
	}
	if up.RestaurantID != nil {
		promo.RestaurantID = *up.RestaurantID
	}
	if up.MaxDiscountAmount != nil {
		promo.MaxDiscountAmount = *up.MaxDiscountAmount
	}
	if up.UsageLimit != nil {
		promo.UsageLimit = *up.UsageLimit
	}
	if up.StartDate != nil {
		promo.StartDate = *up.StartDate
	}
	if up.EndDate != nil {
		promo.EndDate = *up.EndDate
	}

	promo.DateUpdated = time.Now()

	if err := b.storer.Update(ctx, promo); err != nil {
		return Promotion{}, fmt.Errorf("update: %w", err)
	}

	return promo, nil
}

func (b *Business) Delete(ctx context.Context, promo Promotion) error {
	if err := b.storer.Delete(ctx, promo); err != nil {
		return fmt.Errorf("delete: %w", err)
	}
	return nil
}

func (b *Business) Query(ctx context.Context, filter QueryFilter, orderBy order.By, page page.Page) ([]Promotion, error) {
	promos, err := b.storer.Query(ctx, filter, orderBy, page)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	return promos, nil
}

func (b *Business) Count(ctx context.Context, filter QueryFilter) (int, error) {
	total, err := b.storer.Count(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("count: %w", err)
	}
	return total, nil
}

func (b *Business) QueryByID(ctx context.Context, promotionID uuid.UUID) (Promotion, error) {
	promo, err := b.storer.QueryByID(ctx, promotionID)
	if err != nil {
		return Promotion{}, fmt.Errorf("query by id: %w", err)
	}
	return promo, nil
}

func (b *Business) QueryByCode(ctx context.Context, code string) (Promotion, error) {
	cleanCode := strings.ToUpper(strings.TrimSpace(code))
	promo, err := b.storer.QueryByCode(ctx, cleanCode)
	if err != nil {
		return Promotion{}, fmt.Errorf("query by code: %w", err)
	}
	return promo, nil
}

func (b *Business) IncrementUsage(ctx context.Context, promotionID uuid.UUID) error {
	if err := b.storer.IncrementUsage(ctx, promotionID); err != nil {
		return fmt.Errorf("increment usage: %w", err)
	}
	return nil
}

func (b *Business) ValidatePromoCode(ctx context.Context, code string, restaurantID *uuid.UUID, subtotal float64) (ValidateResult, error) {
	cleanCode := strings.ToUpper(strings.TrimSpace(code))
	if cleanCode == "" {
		return ValidateResult{Valid: false, Reason: "Promo code is required"}, nil
	}

	promo, err := b.storer.QueryByCode(ctx, cleanCode)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return ValidateResult{Valid: false, Reason: "Invalid promo code"}, nil
		}
		return ValidateResult{}, fmt.Errorf("query promo by code: %w", err)
	}

	if !promo.Enabled {
		return ValidateResult{Valid: false, Reason: "Promo code is inactive"}, nil
	}

	now := time.Now()
	if promo.StartDate != nil && now.Before(*promo.StartDate) {
		return ValidateResult{Valid: false, Reason: "Promo code campaign has not started yet"}, nil
	}
	if promo.EndDate != nil && now.After(*promo.EndDate) {
		return ValidateResult{Valid: false, Reason: "Promo code has expired"}, nil
	}

	if promo.UsageLimit != nil && promo.UsageCount >= *promo.UsageLimit {
		return ValidateResult{Valid: false, Reason: "Promo code usage limit reached"}, nil
	}

	if promo.RestaurantID != nil {
		if restaurantID == nil || *restaurantID != *promo.RestaurantID {
			return ValidateResult{Valid: false, Reason: "Promo code is not applicable to this restaurant"}, nil
		}
	}

	if subtotal < promo.MinOrderAmount {
		return ValidateResult{
			Valid:  false,
			Reason: fmt.Sprintf("Minimum order subtotal of $%.2f required for this promo code", promo.MinOrderAmount),
		}, nil
	}

	discountAmount := b.CalculateDiscount(promo, subtotal)
	finalSubtotal := math.Max(0, subtotal-discountAmount)

	return ValidateResult{
		Valid:          true,
		Code:           promo.Code,
		DiscountType:   promo.DiscountType,
		DiscountValue:  promo.DiscountValue,
		DiscountAmount: roundToTwoDecimals(discountAmount),
		FinalSubtotal:  roundToTwoDecimals(finalSubtotal),
		Promotion:      &promo,
	}, nil
}

func (b *Business) CalculateDiscount(promo Promotion, subtotal float64) float64 {
	if subtotal <= 0 {
		return 0
	}

	var discount float64
	if promo.DiscountType == DiscountTypePercentage {
		discount = subtotal * (promo.DiscountValue / 100.0)
		if promo.MaxDiscountAmount != nil && discount > *promo.MaxDiscountAmount {
			discount = *promo.MaxDiscountAmount
		}
	} else if promo.DiscountType == DiscountTypeFixedAmount {
		discount = promo.DiscountValue
	}

	if discount > subtotal {
		discount = subtotal
	}

	return roundToTwoDecimals(discount)
}

func roundToTwoDecimals(val float64) float64 {
	return math.Round(val*100) / 100
}

// TestSeedPromotions is a helper for seeding promotions in unit/integration tests.
func TestSeedPromotions(ctx context.Context, n int, bus *Business) ([]Promotion, error) {
	promos := make([]Promotion, n)
	for i := 0; i < n; i++ {
		np := NewPromotion{
			Code:           fmt.Sprintf("TESTPROMO%d", i+1),
			Name:           name.MustParse(fmt.Sprintf("Test Promo %d", i+1)),
			Description:    fmt.Sprintf("Description for test promo %d", i+1),
			DiscountType:   DiscountTypePercentage,
			DiscountValue:  10 + float64(i*5),
			MinOrderAmount: 15.0,
			Enabled:        true,
		}

		promo, err := bus.Create(ctx, np)
		if err != nil {
			return nil, fmt.Errorf("seeding promo %d: %w", i, err)
		}

		promos[i] = promo
	}

	return promos, nil
}
