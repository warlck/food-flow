package menuitembus

import (
	"time"

	"github.com/google/uuid"
	"github.com/warlck/food-flow/business/types/money"
	"github.com/warlck/food-flow/business/types/name"
)

// MenuItem represents a menu item in the system.
type MenuItem struct {
	ID           uuid.UUID
	Name         name.Name
	Description  name.Null
	Price        money.Money
	CategoryID   uuid.UUID
	RestaurantID uuid.UUID
	ImageURL     name.Null
	Available    bool
	DateCreated  time.Time
	DateUpdated  time.Time
}

// NewMenuItem contains information needed to create a new menu item.
type NewMenuItem struct {
	Name         name.Name
	Description  name.Null
	Price        money.Money
	CategoryID   uuid.UUID
	RestaurantID uuid.UUID
	ImageURL     name.Null
}

// UpdateMenuItem contains information needed to update a menu item.
type UpdateMenuItem struct {
	Name        *name.Name
	Description *name.Null
	Price       *money.Money
	CategoryID  *uuid.UUID
	ImageURL    *name.Null
	Available   *bool
}
