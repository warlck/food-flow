package categorybus

import (
	"time"

	"github.com/google/uuid"
	"github.com/warlck/food-flow/business/types/name"
	"github.com/warlck/food-flow/business/types/opt"
)

// Category represents a food category in the system.
type Category struct {
	ID           uuid.UUID
	Name         name.Name
	Description  string
	RestaurantID uuid.UUID
	Enabled      bool
	Rank         *int
	DateCreated  time.Time
	DateUpdated  time.Time
}

// NewCategory contains information needed to create a new category.
type NewCategory struct {
	Name         name.Name
	Description  string
	RestaurantID uuid.UUID
	Rank         *int
}

// UpdateCategory contains information needed to update a category.
type UpdateCategory struct {
	Name        *name.Name
	Description *string
	Enabled     *bool
	Rank        opt.NullInt
}
