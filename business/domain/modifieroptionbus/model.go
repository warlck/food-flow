package modifieroptionbus

import (
	"time"

	"github.com/google/uuid"
	"github.com/warlck/food-flow/business/types/money"
	"github.com/warlck/food-flow/business/types/name"
	"github.com/warlck/food-flow/business/types/opt"
)

// ModifierOption represents a selectable option within a modifier group.
type ModifierOption struct {
	ID              uuid.UUID
	ModifierGroupID uuid.UUID
	RestaurantID    uuid.UUID
	Name            name.Name
	Description     string
	PriceDelta      money.Money
	Available       bool
	Rank            *int
	DateCreated     time.Time
	DateUpdated     time.Time
}

// NewModifierOption contains information needed to create a new modifier option.
type NewModifierOption struct {
	ModifierGroupID uuid.UUID
	RestaurantID    uuid.UUID
	Name            name.Name
	Description     string
	PriceDelta      money.Money
	Available       *bool
	Rank            *int
}

// UpdateModifierOption contains information needed to update a modifier option.
type UpdateModifierOption struct {
	Name        *name.Name
	Description *string
	PriceDelta  *money.Money
	Available   *bool
	Rank        opt.NullInt
}
