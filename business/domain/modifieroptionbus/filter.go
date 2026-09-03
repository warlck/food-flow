package modifieroptionbus

import (
	"github.com/google/uuid"
	"github.com/warlck/food-flow/business/types/name"
)

// QueryFilter holds the available fields a query can be filtered on.
type QueryFilter struct {
	ID              *uuid.UUID
	ModifierGroupID *uuid.UUID
	RestaurantID    *uuid.UUID
	Name            *name.Name
	Available       *bool
}
