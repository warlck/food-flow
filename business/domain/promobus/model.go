package promobus

import (
	"time"

	"github.com/google/uuid"
	"github.com/warlck/food-flow/business/types/name"
)

// DiscountType constants.
const (
	DiscountTypePercentage  = "percentage"
	DiscountTypeFixedAmount = "fixed_amount"
)

// Promotion represents information about a promo code campaign.
type Promotion struct {
	ID                uuid.UUID
	RestaurantID      *uuid.UUID
	Code              string
	Name              name.Name
	Description       string
	DiscountType      string
	DiscountValue     float64
	MinOrderAmount    float64
	MaxDiscountAmount *float64
	UsageLimit        *int
	UsageCount        int
	StartDate         *time.Time
	EndDate           *time.Time
	Enabled           bool
	DateCreated       time.Time
	DateUpdated       time.Time
}

// NewPromotion contains information needed to create a promotion.
type NewPromotion struct {
	RestaurantID      *uuid.UUID
	Code              string
	Name              name.Name
	Description       string
	DiscountType      string
	DiscountValue     float64
	MinOrderAmount    float64
	MaxDiscountAmount *float64
	UsageLimit        *int
	StartDate         *time.Time
	EndDate           *time.Time
	Enabled           bool
}

// UpdatePromotion contains information needed to update a promotion.
type UpdatePromotion struct {
	RestaurantID      **uuid.UUID
	Code              *string
	Name              *name.Name
	Description       *string
	DiscountType      *string
	DiscountValue     *float64
	MinOrderAmount    *float64
	MaxDiscountAmount **float64
	UsageLimit        **int
	StartDate         **time.Time
	EndDate           **time.Time
	Enabled           *bool
}

// ValidateResult represents the outcome of validating a promo code.
type ValidateResult struct {
	Valid          bool       `json:"valid"`
	Reason         string     `json:"reason,omitempty"`
	Code           string     `json:"code,omitempty"`
	DiscountType   string     `json:"discountType,omitempty"`
	DiscountValue  float64    `json:"discountValue"`
	DiscountAmount float64    `json:"discountAmount"`
	FinalSubtotal  float64    `json:"finalSubtotal"`
	Promotion      *Promotion `json:"-"`
}
