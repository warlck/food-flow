package promoapp

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/warlck/food-flow/app/sdk/errs"
	"github.com/warlck/food-flow/business/domain/promobus"
	"github.com/warlck/food-flow/business/types/name"
)

// ValidateRequest defines data needed to validate a promo code.
type ValidateRequest struct {
	PromoCode    string  `json:"promoCode" validate:"required"`
	RestaurantID *string `json:"restaurantId"`
	Subtotal     float64 `json:"subtotal" validate:"gte=0"`
}

// Decode implements decoder interface.
func (app *ValidateRequest) Decode(data []byte) error {
	return json.Unmarshal(data, app)
}

// Validate checks request validity.
func (app ValidateRequest) Validate() error {
	if err := errs.Check(app); err != nil {
		return fmt.Errorf("validate: %w", err)
	}
	return nil
}

// ValidateResponse defines the response structure for promo validation.
type ValidateResponse struct {
	Valid          bool    `json:"valid"`
	Reason         string  `json:"reason,omitempty"`
	Code           string  `json:"code,omitempty"`
	DiscountType   string  `json:"discountType,omitempty"`
	DiscountValue  float64 `json:"discountValue"`
	DiscountAmount float64 `json:"discountAmount"`
	FinalSubtotal  float64 `json:"finalSubtotal"`
}

// Encode implements encoder interface.
func (app ValidateResponse) Encode() ([]byte, string, error) {
	data, err := json.Marshal(app)
	return data, "application/json", err
}

// =============================================================================

// Promotion represents information about a promotion campaign for API responses.
type Promotion struct {
	ID                string   `json:"id"`
	RestaurantID      *string  `json:"restaurantId,omitempty"`
	Code              string   `json:"code"`
	Name              string   `json:"name"`
	Description       string   `json:"description,omitempty"`
	DiscountType      string   `json:"discountType"`
	DiscountValue     float64  `json:"discountValue"`
	MinOrderAmount    float64  `json:"minOrderAmount"`
	MaxDiscountAmount *float64 `json:"maxDiscountAmount,omitempty"`
	UsageLimit        *int     `json:"usageLimit,omitempty"`
	UsageCount        int      `json:"usageCount"`
	StartDate         *string  `json:"startDate,omitempty"`
	EndDate           *string  `json:"endDate,omitempty"`
	Enabled           bool     `json:"enabled"`
	DateCreated       string   `json:"dateCreated"`
	DateUpdated       string   `json:"dateUpdated"`
}

// Encode implements encoder interface.
func (app Promotion) Encode() ([]byte, string, error) {
	data, err := json.Marshal(app)
	return data, "application/json", err
}

// ToAppPromotion converts a business layer promotion to an app layer promotion.
func ToAppPromotion(bus promobus.Promotion) Promotion {
	var restaurantID *string
	if bus.RestaurantID != nil {
		id := bus.RestaurantID.String()
		restaurantID = &id
	}

	var startDate *string
	if bus.StartDate != nil {
		d := bus.StartDate.Format(time.RFC3339)
		startDate = &d
	}

	var endDate *string
	if bus.EndDate != nil {
		d := bus.EndDate.Format(time.RFC3339)
		endDate = &d
	}

	return Promotion{
		ID:                bus.ID.String(),
		RestaurantID:      restaurantID,
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
		DateCreated:       bus.DateCreated.Format(time.RFC3339),
		DateUpdated:       bus.DateUpdated.Format(time.RFC3339),
	}
}

// ToAppPromotions converts a slice of business layer promotions to app layer promotions.
func ToAppPromotions(promos []promobus.Promotion) []Promotion {
	app := make([]Promotion, len(promos))
	for i, p := range promos {
		app[i] = ToAppPromotion(p)
	}
	return app
}

// =============================================================================

// NewPromotion defines data needed to create a promotion.
type NewPromotion struct {
	RestaurantID      *string  `json:"restaurantId"`
	Code              string   `json:"code" validate:"required"`
	Name              string   `json:"name" validate:"required"`
	Description       string   `json:"description"`
	DiscountType      string   `json:"discountType" validate:"required,oneof=percentage fixed_amount"`
	DiscountValue     float64  `json:"discountValue" validate:"required,gt=0"`
	MinOrderAmount    float64  `json:"minOrderAmount" validate:"gte=0"`
	MaxDiscountAmount *float64 `json:"maxDiscountAmount" validate:"omitempty,gt=0"`
	UsageLimit        *int     `json:"usageLimit" validate:"omitempty,gt=0"`
	StartDate         *string  `json:"startDate"`
	EndDate           *string  `json:"endDate"`
	Enabled           bool     `json:"enabled"`
}

// Decode implements decoder interface.
func (app *NewPromotion) Decode(data []byte) error {
	return json.Unmarshal(data, app)
}

// Validate checks request validity.
func (app NewPromotion) Validate() error {
	if err := errs.Check(app); err != nil {
		return fmt.Errorf("validate: %w", err)
	}

	if app.DiscountType == promobus.DiscountTypePercentage && app.DiscountValue > 100 {
		return fmt.Errorf("validate: percentage discount value cannot exceed 100")
	}

	return nil
}

func toBusNewPromotion(app NewPromotion) (promobus.NewPromotion, error) {
	nme, err := name.Parse(app.Name)
	if err != nil {
		return promobus.NewPromotion{}, fmt.Errorf("parse name: %w", err)
	}

	var restaurantID *uuid.UUID
	if app.RestaurantID != nil && *app.RestaurantID != "" {
		id, err := uuid.Parse(*app.RestaurantID)
		if err != nil {
			return promobus.NewPromotion{}, fmt.Errorf("parse restaurantID: %w", err)
		}
		restaurantID = &id
	}

	var startDate *time.Time
	if app.StartDate != nil && *app.StartDate != "" {
		t, err := time.Parse(time.RFC3339, *app.StartDate)
		if err != nil {
			return promobus.NewPromotion{}, fmt.Errorf("parse startDate: %w", err)
		}
		startDate = &t
	}

	var endDate *time.Time
	if app.EndDate != nil && *app.EndDate != "" {
		t, err := time.Parse(time.RFC3339, *app.EndDate)
		if err != nil {
			return promobus.NewPromotion{}, fmt.Errorf("parse endDate: %w", err)
		}
		endDate = &t
	}

	return promobus.NewPromotion{
		RestaurantID:      restaurantID,
		Code:              app.Code,
		Name:              nme,
		Description:       app.Description,
		DiscountType:      app.DiscountType,
		DiscountValue:     app.DiscountValue,
		MinOrderAmount:    app.MinOrderAmount,
		MaxDiscountAmount: app.MaxDiscountAmount,
		UsageLimit:        app.UsageLimit,
		StartDate:         startDate,
		EndDate:           endDate,
		Enabled:           app.Enabled,
	}, nil
}

// =============================================================================

// UpdatePromotion defines data needed to update a promotion.
type UpdatePromotion struct {
	RestaurantID      **string  `json:"restaurantId"`
	Code              *string   `json:"code"`
	Name              *string   `json:"name"`
	Description       *string   `json:"description"`
	DiscountType      *string   `json:"discountType" validate:"omitempty,oneof=percentage fixed_amount"`
	DiscountValue     *float64  `json:"discountValue" validate:"omitempty,gt=0"`
	MinOrderAmount    *float64  `json:"minOrderAmount" validate:"omitempty,gte=0"`
	MaxDiscountAmount **float64 `json:"maxDiscountAmount"`
	UsageLimit        **int     `json:"usageLimit"`
	StartDate         **string  `json:"startDate"`
	EndDate           **string  `json:"endDate"`
	Enabled           *bool     `json:"enabled"`
}

// Decode implements decoder interface.
func (app *UpdatePromotion) Decode(data []byte) error {
	return json.Unmarshal(data, app)
}

// Validate checks request validity.
func (app UpdatePromotion) Validate() error {
	if err := errs.Check(app); err != nil {
		return fmt.Errorf("validate: %w", err)
	}

	if app.DiscountType != nil && *app.DiscountType == promobus.DiscountTypePercentage {
		if app.DiscountValue != nil && *app.DiscountValue > 100 {
			return fmt.Errorf("validate: percentage discount value cannot exceed 100")
		}
	}

	return nil
}

func toBusUpdatePromotion(app UpdatePromotion) (promobus.UpdatePromotion, error) {
	var nme *name.Name
	if app.Name != nil {
		n, err := name.Parse(*app.Name)
		if err != nil {
			return promobus.UpdatePromotion{}, fmt.Errorf("parse name: %w", err)
		}
		nme = &n
	}

	var restaurantID **uuid.UUID
	if app.RestaurantID != nil {
		if *app.RestaurantID != nil && **app.RestaurantID != "" {
			id, err := uuid.Parse(**app.RestaurantID)
			if err != nil {
				return promobus.UpdatePromotion{}, fmt.Errorf("parse restaurantID: %w", err)
			}
			pID := &id
			restaurantID = &pID
		} else {
			var pID *uuid.UUID
			restaurantID = &pID
		}
	}

	var startDate **time.Time
	if app.StartDate != nil {
		if *app.StartDate != nil && **app.StartDate != "" {
			t, err := time.Parse(time.RFC3339, **app.StartDate)
			if err != nil {
				return promobus.UpdatePromotion{}, fmt.Errorf("parse startDate: %w", err)
			}
			pT := &t
			startDate = &pT
		} else {
			var pT *time.Time
			startDate = &pT
		}
	}

	var endDate **time.Time
	if app.EndDate != nil {
		if *app.EndDate != nil && **app.EndDate != "" {
			t, err := time.Parse(time.RFC3339, **app.EndDate)
			if err != nil {
				return promobus.UpdatePromotion{}, fmt.Errorf("parse endDate: %w", err)
			}
			pT := &t
			endDate = &pT
		} else {
			var pT *time.Time
			endDate = &pT
		}
	}

	return promobus.UpdatePromotion{
		RestaurantID:      restaurantID,
		Code:              app.Code,
		Name:              nme,
		Description:       app.Description,
		DiscountType:      app.DiscountType,
		DiscountValue:     app.DiscountValue,
		MinOrderAmount:    app.MinOrderAmount,
		MaxDiscountAmount: app.MaxDiscountAmount,
		UsageLimit:        app.UsageLimit,
		StartDate:         startDate,
		EndDate:           endDate,
		Enabled:           app.Enabled,
	}, nil
}
