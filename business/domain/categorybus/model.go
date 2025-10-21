package categorybus

import (
	"time"

	"github.com/google/uuid"
	"github.com/warlck/food-flow/business/types/name"
)

// Category represents a food category in the system.
type Category struct {
	ID           uuid.UUID
	Name         name.Name
	Description  name.Null
	RestaurantID uuid.UUID
	Enabled      bool
	DateCreated  time.Time
	DateUpdated  time.Time
}

// NewCategory contains information needed to create a new category.
type NewCategory struct {
	Name         name.Name
	Description  name.Null
	RestaurantID uuid.UUID
}

// UpdateCategory contains information needed to update a category.
type UpdateCategory struct {
	Name        *name.Name
	Description *name.Null
	Enabled     *bool
}
