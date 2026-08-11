package promobus

import (
	"github.com/google/uuid"
	"github.com/warlck/food-flow/business/types/name"
)

// QueryFilter holds fields for filtering Promotion queries.
type QueryFilter struct {
	ID           *uuid.UUID
	Code         *string
	Name         *name.Name
	RestaurantID *uuid.UUID
	Enabled      *bool
}
