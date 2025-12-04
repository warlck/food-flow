package addonbus

import "github.com/google/uuid"

// QueryFilter holds the available fields for filtering addons.
type QueryFilter struct {
	ID           *uuid.UUID
	MenuItemID   *uuid.UUID
	RestaurantID *uuid.UUID
	Name         *string
	Available    *bool
}
