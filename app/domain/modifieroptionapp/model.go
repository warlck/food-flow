package modifieroptionapp

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/warlck/food-flow/app/sdk/errs"
	"github.com/warlck/food-flow/business/domain/modifieroptionbus"
	"github.com/warlck/food-flow/business/types/money"
	"github.com/warlck/food-flow/business/types/name"
	"github.com/warlck/food-flow/business/types/opt"
)

// ModifierOption represents a modifier option for API responses.
type ModifierOption struct {
	ID              string  `json:"id"`
	ModifierGroupID string  `json:"modifierGroupId"`
	RestaurantID    string  `json:"restaurantId"`
	Name            string  `json:"name"`
	Description     string  `json:"description"`
	PriceDelta      float64 `json:"priceDelta"`
	Available       bool    `json:"available"`
	Rank            *int    `json:"rank"`
	DateCreated     string  `json:"dateCreated"`
	DateUpdated     string  `json:"dateUpdated"`
}

// Encode implements the encoder interface.
func (app ModifierOption) Encode() ([]byte, string, error) {
	data, err := json.Marshal(app)
	return data, "application/json", err
}

// ToAppModifierOption converts a business layer modifier option to an app layer modifier option.
func ToAppModifierOption(bus modifieroptionbus.ModifierOption) ModifierOption {
	return ModifierOption{
		ID:              bus.ID.String(),
		ModifierGroupID: bus.ModifierGroupID.String(),
		RestaurantID:    bus.RestaurantID.String(),
		Name:            bus.Name.String(),
		Description:     bus.Description,
		PriceDelta:      bus.PriceDelta.Value(),
		Available:       bus.Available,
		Rank:            bus.Rank,
		DateCreated:     bus.DateCreated.Format(time.RFC3339),
		DateUpdated:     bus.DateUpdated.Format(time.RFC3339),
	}
}

// ToAppModifierOptions converts a slice of business layer modifier options.
func ToAppModifierOptions(options []modifieroptionbus.ModifierOption) []ModifierOption {
	app := make([]ModifierOption, len(options))
	for i, option := range options {
		app[i] = ToAppModifierOption(option)
	}
	return app
}

// =============================================================================

// NewModifierOption defines the data needed to add a new modifier option.
type NewModifierOption struct {
	ModifierGroupID string  `json:"modifierGroupId" validate:"required,uuid"`
	RestaurantID    string  `json:"restaurantId" validate:"required,uuid"`
	Name            string  `json:"name" validate:"required"`
	Description     string  `json:"description"`
	PriceDelta      float64 `json:"priceDelta" validate:"gte=0"`
	Available       *bool   `json:"available"`
	Rank            *int    `json:"rank" validate:"omitempty,gte=1"`
}

// Decode implements the decoder interface.
func (app *NewModifierOption) Decode(data []byte) error {
	return json.Unmarshal(data, app)
}

// Validate checks the data in the model is clean.
func (app NewModifierOption) Validate() error {
	if err := errs.Check(app); err != nil {
		return fmt.Errorf("validate: %w", err)
	}
	return nil
}

func toBusNewModifierOption(app NewModifierOption) (modifieroptionbus.NewModifierOption, error) {
	nme, err := name.Parse(app.Name)
	if err != nil {
		return modifieroptionbus.NewModifierOption{}, fmt.Errorf("parse name: %w", err)
	}

	priceDelta, err := money.Parse(app.PriceDelta)
	if err != nil {
		return modifieroptionbus.NewModifierOption{}, fmt.Errorf("parse priceDelta: %w", err)
	}

	modifierGroupID, err := uuid.Parse(app.ModifierGroupID)
	if err != nil {
		return modifieroptionbus.NewModifierOption{}, fmt.Errorf("parse modifierGroupId: %w", err)
	}

	restaurantID, err := uuid.Parse(app.RestaurantID)
	if err != nil {
		return modifieroptionbus.NewModifierOption{}, fmt.Errorf("parse restaurantId: %w", err)
	}

	bus := modifieroptionbus.NewModifierOption{
		ModifierGroupID: modifierGroupID,
		RestaurantID:    restaurantID,
		Name:            nme,
		Description:     app.Description,
		PriceDelta:      priceDelta,
		Available:       app.Available,
		Rank:            app.Rank,
	}

	return bus, nil
}

// =============================================================================

// UpdateModifierOption defines the data needed to update a modifier option.
type UpdateModifierOption struct {
	Name        *string     `json:"name"`
	Description *string     `json:"description"`
	PriceDelta  *float64    `json:"priceDelta" validate:"omitempty,gte=0"`
	Available   *bool       `json:"available"`
	Rank        opt.NullInt `json:"rank"`
}

// Decode implements the decoder interface.
func (app *UpdateModifierOption) Decode(data []byte) error {
	return json.Unmarshal(data, app)
}

// Validate checks the data in the model is clean.
func (app UpdateModifierOption) Validate() error {
	if err := errs.Check(app); err != nil {
		return fmt.Errorf("validate: %w", err)
	}
	if app.Rank.Present && app.Rank.Value != nil && *app.Rank.Value < 1 {
		return errs.NewFieldErrors("rank", fmt.Errorf("rank must be >= 1"))
	}
	return nil
}

func toBusUpdateModifierOption(app UpdateModifierOption) (modifieroptionbus.UpdateModifierOption, error) {
	var nme *name.Name
	if app.Name != nil {
		nm, err := name.Parse(*app.Name)
		if err != nil {
			return modifieroptionbus.UpdateModifierOption{}, fmt.Errorf("parse name: %w", err)
		}
		nme = &nm
	}

	var priceDelta *money.Money
	if app.PriceDelta != nil {
		p, err := money.Parse(*app.PriceDelta)
		if err != nil {
			return modifieroptionbus.UpdateModifierOption{}, fmt.Errorf("parse priceDelta: %w", err)
		}
		priceDelta = &p
	}

	bus := modifieroptionbus.UpdateModifierOption{
		Name:        nme,
		Description: app.Description,
		PriceDelta:  priceDelta,
		Available:   app.Available,
		Rank:        app.Rank,
	}

	return bus, nil
}

// =============================================================================

// ReorderModifierOptions defines the payload for reordering modifier options within a group.
type ReorderModifierOptions struct {
	ModifierGroupID string   `json:"modifierGroupId" validate:"required,uuid"`
	OrderedIDs      []string `json:"orderedIds" validate:"required,min=1,dive,uuid"`
}

// Decode implements the decoder interface.
func (app *ReorderModifierOptions) Decode(data []byte) error {
	return json.Unmarshal(data, app)
}

// Validate checks the data in the model is clean.
func (app ReorderModifierOptions) Validate() error {
	if err := errs.Check(app); err != nil {
		return fmt.Errorf("validate: %w", err)
	}
	return nil
}
