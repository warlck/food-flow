package promobus_test

import (
	"context"
	"testing"
	"time"

	"github.com/warlck/food-flow/business/domain/promobus"
	"github.com/warlck/food-flow/business/sdk/dbtest"
	"github.com/warlck/food-flow/business/types/name"
)

func Test_Promotion(t *testing.T) {
	t.Parallel()

	db := dbtest.New(t, "Test_Promotion")

	ctx := context.Background()

	// 1. Create a percentage promo
	np := promobus.NewPromotion{
		Code:              "SAVE15",
		Name:              name.MustParse("Save 15 Percent"),
		Description:       "15% off orders over $20",
		DiscountType:      promobus.DiscountTypePercentage,
		DiscountValue:     15.0,
		MinOrderAmount:    20.0,
		MaxDiscountAmount: ptr(10.0),
		UsageLimit:        ptr(100),
		Enabled:           true,
	}

	promo, err := db.BusDomain.Promo.Create(ctx, np)
	if err != nil {
		t.Fatalf("Create promotion: %s", err)
	}

	if promo.Code != "SAVE15" {
		t.Errorf("expected code SAVE15, got %s", promo.Code)
	}

	// 2. Query by ID and Code
	fetchedByID, err := db.BusDomain.Promo.QueryByID(ctx, promo.ID)
	if err != nil {
		t.Fatalf("QueryByID error: %s", err)
	}
	if fetchedByID.Name.String() != "Save 15 Percent" {
		t.Errorf("expected name Save 15 Percent, got %s", fetchedByID.Name.String())
	}

	fetchedByCode, err := db.BusDomain.Promo.QueryByCode(ctx, "save15")
	if err != nil {
		t.Fatalf("QueryByCode error: %s", err)
	}
	if fetchedByCode.ID != promo.ID {
		t.Errorf("expected ID %s, got %s", promo.ID, fetchedByCode.ID)
	}

	// 3. Test ValidatePromoCode
	// Subtotal below min order
	res, err := db.BusDomain.Promo.ValidatePromoCode(ctx, "SAVE15", nil, 10.0)
	if err != nil {
		t.Fatalf("ValidatePromoCode error: %s", err)
	}
	if res.Valid {
		t.Errorf("expected invalid for subtotal below min order")
	}

	// Valid subtotal: $50 subtotal * 15% = $7.50 discount
	res, err = db.BusDomain.Promo.ValidatePromoCode(ctx, "SAVE15", nil, 50.0)
	if err != nil {
		t.Fatalf("ValidatePromoCode error: %s", err)
	}
	if !res.Valid {
		t.Errorf("expected valid promo code, got reason: %s", res.Reason)
	}
	if res.DiscountAmount != 7.50 {
		t.Errorf("expected discount $7.50, got $%.2f", res.DiscountAmount)
	}
	if res.FinalSubtotal != 42.50 {
		t.Errorf("expected final subtotal $42.50, got $%.2f", res.FinalSubtotal)
	}

	// Cap test: $100 subtotal * 15% = $15.00, capped at $10.00
	res, err = db.BusDomain.Promo.ValidatePromoCode(ctx, "SAVE15", nil, 100.0)
	if err != nil {
		t.Fatalf("ValidatePromoCode error: %s", err)
	}
	if res.DiscountAmount != 10.00 {
		t.Errorf("expected capped discount $10.00, got $%.2f", res.DiscountAmount)
	}

	// 4. Update Promotion fields (Exhaustive Testing)
	up := promobus.UpdatePromotion{
		Name:           ptr(name.MustParse("Save 20 Percent")),
		DiscountValue:  ptr(20.0),
		MinOrderAmount: ptr(25.0),
		Enabled:        ptr(false),
	}

	updated, err := db.BusDomain.Promo.Update(ctx, promo, up)
	if err != nil {
		t.Fatalf("Update error: %s", err)
	}
	if updated.Name.String() != "Save 20 Percent" {
		t.Errorf("expected updated name Save 20 Percent, got %s", updated.Name.String())
	}
	if updated.DiscountValue != 20.0 {
		t.Errorf("expected updated value 20.0, got %f", updated.DiscountValue)
	}

	// Validation after disable
	res, err = db.BusDomain.Promo.ValidatePromoCode(ctx, "SAVE15", nil, 50.0)
	if err != nil {
		t.Fatalf("ValidatePromoCode error: %s", err)
	}
	if res.Valid {
		t.Errorf("expected invalid for disabled promo code")
	}

	// 5. Test Fixed Amount Promotion
	npFixed := promobus.NewPromotion{
		Code:           "FLAT10",
		Name:           name.MustParse("Flat 10 Off"),
		DiscountType:   promobus.DiscountTypeFixedAmount,
		DiscountValue:  10.0,
		MinOrderAmount: 0.0,
		Enabled:        true,
	}
	promoFixed, err := db.BusDomain.Promo.Create(ctx, npFixed)
	if err != nil {
		t.Fatalf("Create fixed promo: %s", err)
	}

	res, err = db.BusDomain.Promo.ValidatePromoCode(ctx, "FLAT10", nil, 25.0)
	if err != nil {
		t.Fatalf("Validate fixed promo error: %s", err)
	}
	if res.DiscountAmount != 10.0 {
		t.Errorf("expected discount $10.00, got $%.2f", res.DiscountAmount)
	}

	// 6. Delete Promotion
	if err := db.BusDomain.Promo.Delete(ctx, promoFixed); err != nil {
		t.Fatalf("Delete error: %s", err)
	}
	_, err = db.BusDomain.Promo.QueryByID(ctx, promoFixed.ID)
	if err == nil {
		t.Errorf("expected error querying deleted promo, got nil")
	}
}

func ptr[T any](v T) *T {
	return &v
}

func ptrTime(t time.Time) *time.Time {
	return &t
}
