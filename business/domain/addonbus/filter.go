package addonbus

import (
	"github.com/google/uuid"
	"github.com/warlck/food-flow/business/types/name"
)

// QueryFilter holds the available fields for filtering addons.
type QueryFilter struct {
	ID           *uuid.UUID
	MenuItemID   *uuid.UUID
	RestaurantID *uuid.UUID
	Name         *name.Name
	Available    *bool
}
