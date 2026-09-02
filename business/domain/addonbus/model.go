package addonbus

import (
	"time"

	"github.com/google/uuid"
	"github.com/warlck/food-flow/business/types/money"
	"github.com/warlck/food-flow/business/types/name"
)

// Addon represents a restaurant-level addon definition in the system.
type Addon struct {
	ID           uuid.UUID
	RestaurantID uuid.UUID
	Name         name.Name
	Description  string
	Price        money.Money
	Available    bool
	MaxQuantity  int
	DateCreated  time.Time
	DateUpdated  time.Time
}

// NewAddon contains information needed to create a new addon definition.
type NewAddon struct {
	RestaurantID uuid.UUID
	Name         name.Name
	Description  string
	Price        money.Money
	Available    *bool
	MaxQuantity  int
}

// UpdateAddon contains information needed to update an addon definition.
type UpdateAddon struct {
	Name        *name.Name
	Description *string
	Price       *money.Money
	Available   *bool
	MaxQuantity *int
}

// ItemAddonAssignment represents an item-to-addon assignment with an optional rank.
type ItemAddonAssignment struct {
	AddonID uuid.UUID
	Rank    *int
}

// MenuItemAddonInfo represents an assigned add-on on a specific menu item.
type MenuItemAddonInfo struct {
	Addon Addon
	Rank  *int
}
