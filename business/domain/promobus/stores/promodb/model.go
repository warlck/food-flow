package promodb

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/warlck/food-flow/business/domain/promobus"
	"github.com/warlck/food-flow/business/types/name"
)

type promotion struct {
	ID                uuid.UUID  `db:"promotion_id"`
	RestaurantID      *uuid.UUID `db:"restaurant_id"`
	Code              string     `db:"code"`
	Name              string     `db:"name"`
	Description       string     `db:"description"`
	DiscountType      string     `db:"discount_type"`
	DiscountValue     float64    `db:"discount_value"`
	MinOrderAmount    float64    `db:"min_order_amount"`
	MaxDiscountAmount *float64   `db:"max_discount_amount"`
	UsageLimit        *int       `db:"usage_limit"`
	UsageCount        int        `db:"usage_count"`
	StartDate         *time.Time `db:"start_date"`
	EndDate           *time.Time `db:"end_date"`
	Enabled           bool       `db:"enabled"`
	DateCreated       time.Time  `db:"date_created"`
	DateUpdated       time.Time  `db:"date_updated"`
}

func toDBPromotion(bus promobus.Promotion) promotion {
	var startDate *time.Time
	if bus.StartDate != nil {
		t := bus.StartDate.UTC()
		startDate = &t
	}
	var endDate *time.Time
	if bus.EndDate != nil {
		t := bus.EndDate.UTC()
		endDate = &t
	}

	return promotion{
		ID:                bus.ID,
		RestaurantID:      bus.RestaurantID,
		Code:              bus.Code,
		Name:              bus.Name.String(),
		Description:       bus.Description,
		DiscountType:      bus.DiscountType,
		DiscountValue:     bus.DiscountValue,
		MinOrderAmount:    bus.MinOrderAmount,
		MaxDiscountAmount: bus.MaxDiscountAmount,
		UsageLimit:        bus.UsageLimit,
		UsageCount:        bus.UsageCount,
		StartDate:         startDate,
		EndDate:           endDate,
		Enabled:           bus.Enabled,
		DateCreated:       bus.DateCreated.UTC(),
		DateUpdated:       bus.DateUpdated.UTC(),
	}
}

func toBusPromotion(db promotion) (promobus.Promotion, error) {
	nme, err := name.Parse(db.Name)
	if err != nil {
		return promobus.Promotion{}, fmt.Errorf("parse name: %w", err)
	}

	var startDate *time.Time
	if db.StartDate != nil {
		t := db.StartDate.In(time.Local)
		startDate = &t
	}
	var endDate *time.Time
	if db.EndDate != nil {
		t := db.EndDate.In(time.Local)
		endDate = &t
	}

	bus := promobus.Promotion{
		ID:                db.ID,
		RestaurantID:      db.RestaurantID,
		Code:              db.Code,
		Name:              nme,
		Description:       db.Description,
		DiscountType:      db.DiscountType,
		DiscountValue:     db.DiscountValue,
		MinOrderAmount:    db.MinOrderAmount,
		MaxDiscountAmount: db.MaxDiscountAmount,
		UsageLimit:        db.UsageLimit,
		UsageCount:        db.UsageCount,
		StartDate:         startDate,
		EndDate:           endDate,
		Enabled:           db.Enabled,
		DateCreated:       db.DateCreated.In(time.Local),
		DateUpdated:       db.DateUpdated.In(time.Local),
	}

	return bus, nil
}

func toBusPromotions(dbs []promotion) ([]promobus.Promotion, error) {
	bus := make([]promobus.Promotion, len(dbs))
	for i, db := range dbs {
		var err error
		bus[i], err = toBusPromotion(db)
		if err != nil {
			return nil, err
		}
	}
	return bus, nil
}
