package modifiergroupbus

import (
	"time"

	"github.com/google/uuid"
	"github.com/warlck/food-flow/business/types/name"
	"github.com/warlck/food-flow/business/types/opt"
)

// ModifierGroup represents a modifier group attached to a menu item.
type ModifierGroup struct {
	ID            uuid.UUID
	MenuItemID    uuid.UUID
	RestaurantID  uuid.UUID
	Name          name.Name
	Description   string
	MinSelections int
	MaxSelections int
	Available     bool
	Rank          *int
	DateCreated   time.Time
	DateUpdated   time.Time
}

// NewModifierGroup contains information needed to create a new modifier group.
type NewModifierGroup struct {
	MenuItemID    uuid.UUID
	RestaurantID  uuid.UUID
	Name          name.Name
	Description   string
	MinSelections int
	MaxSelections int
	Available     bool
	Rank          *int
}

// UpdateModifierGroup contains information needed to update a modifier group.
type UpdateModifierGroup struct {
	Name          *name.Name
	Description   *string
	MinSelections *int
	MaxSelections *int
	Available     *bool
	Rank          opt.NullInt
}
