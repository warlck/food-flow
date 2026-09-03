package addonbus

import (
	"time"

	"github.com/google/uuid"
	"github.com/warlck/food-flow/business/types/money"
	"github.com/warlck/food-flow/business/types/name"
	"github.com/warlck/food-flow/business/types/opt"
)

// Addon represents a menu-item level addon in the system.
type Addon struct {
	ID           uuid.UUID
	MenuItemID   uuid.UUID
	RestaurantID uuid.UUID
	Name         name.Name
	Description  string
	Price        money.Money
	Available    bool
	MaxQuantity  int
	Rank         *int
	DateCreated  time.Time
	DateUpdated  time.Time
}

// NewAddon contains information needed to create a new addon.
type NewAddon struct {
	MenuItemID   uuid.UUID
	RestaurantID uuid.UUID
	Name         name.Name
	Description  string
	Price        money.Money
	Available    *bool
	MaxQuantity  int
	Rank         *int
}

// UpdateAddon contains information needed to update an addon.
type UpdateAddon struct {
	Name        *name.Name
	Description *string
	Price       *money.Money
	Available   *bool
	MaxQuantity *int
	Rank        opt.NullInt
}
