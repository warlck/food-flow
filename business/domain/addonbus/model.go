package addonbus

import (
	"time"

	"github.com/google/uuid"
	"github.com/warlck/food-flow/business/types/money"
	"github.com/warlck/food-flow/business/types/name"
)

// Addon represents an addon for a category in the system.
type Addon struct {
	ID           uuid.UUID
	CategoryID   uuid.UUID
	RestaurantID uuid.UUID
	Name         name.Name
	Description  string
	Price        money.Money
	Available    bool
	MaxQuantity  int
	DateCreated  time.Time
	DateUpdated  time.Time
}

// NewAddon contains information needed to create a new addon.
type NewAddon struct {
	CategoryID   uuid.UUID
	RestaurantID uuid.UUID
	Name         name.Name
	Description  string
	Price        money.Money
	MaxQuantity  int
}

// UpdateAddon contains information needed to update an addon.
type UpdateAddon struct {
	Name        *name.Name
	Description *string
	Price       *money.Money
	Available   *bool
	MaxQuantity *int
}
