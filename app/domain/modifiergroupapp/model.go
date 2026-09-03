package modifiergroupapp

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/warlck/food-flow/app/sdk/errs"
	"github.com/warlck/food-flow/business/domain/modifiergroupbus"
	"github.com/warlck/food-flow/business/types/name"
	"github.com/warlck/food-flow/business/types/opt"
)

// ModifierGroup represents a modifier group for API responses.
type ModifierGroup struct {
	ID            string `json:"id"`
	MenuItemID    string `json:"menuItemId"`
	RestaurantID  string `json:"restaurantId"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	MinSelections int    `json:"minSelections"`
	MaxSelections int    `json:"maxSelections"`
	Available     bool   `json:"available"`
	Rank          *int   `json:"rank"`
	DateCreated   string `json:"dateCreated"`
	DateUpdated   string `json:"dateUpdated"`
}

// Encode implements the encoder interface.
func (app ModifierGroup) Encode() ([]byte, string, error) {
	data, err := json.Marshal(app)
	return data, "application/json", err
}

// ToAppModifierGroup converts a business layer modifier group to an app layer modifier group.
func ToAppModifierGroup(bus modifiergroupbus.ModifierGroup) ModifierGroup {
	return ModifierGroup{
		ID:            bus.ID.String(),
		MenuItemID:    bus.MenuItemID.String(),
		RestaurantID:  bus.RestaurantID.String(),
		Name:          bus.Name.String(),
		Description:   bus.Description,
		MinSelections: bus.MinSelections,
		MaxSelections: bus.MaxSelections,
		Available:     bus.Available,
		Rank:          bus.Rank,
		DateCreated:   bus.DateCreated.Format(time.RFC3339),
		DateUpdated:   bus.DateUpdated.Format(time.RFC3339),
	}
}

// ToAppModifierGroups converts a slice of business layer modifier groups.
func ToAppModifierGroups(groups []modifiergroupbus.ModifierGroup) []ModifierGroup {
	app := make([]ModifierGroup, len(groups))
	for i, group := range groups {
		app[i] = ToAppModifierGroup(group)
	}
	return app
}

// =============================================================================

// NewModifierGroup defines the data needed to add a new modifier group.
type NewModifierGroup struct {
	MenuItemID    string `json:"menuItemId" validate:"required,uuid"`
	RestaurantID  string `json:"restaurantId" validate:"required,uuid"`
	Name          string `json:"name" validate:"required"`
	Description   string `json:"description"`
	MinSelections int    `json:"minSelections" validate:"gte=0,lte=1"`
	MaxSelections int    `json:"maxSelections" validate:"required,eq=1"`
	Available     bool   `json:"available"`
	Rank          *int   `json:"rank" validate:"omitempty,gte=1"`
}

// Decode implements the decoder interface.
func (app *NewModifierGroup) Decode(data []byte) error {
	return json.Unmarshal(data, app)
}

// Validate checks the data in the model is clean.
func (app NewModifierGroup) Validate() error {
	if err := errs.Check(app); err != nil {
		return fmt.Errorf("validate: %w", err)
	}

	if app.Available && app.MinSelections >= 1 {
		return errs.NewFieldErrors("available", fmt.Errorf("required modifier group cannot be active before options are added"))
	}

	return nil
}

func toBusNewModifierGroup(app NewModifierGroup) (modifiergroupbus.NewModifierGroup, error) {
	nme, err := name.Parse(app.Name)
	if err != nil {
		return modifiergroupbus.NewModifierGroup{}, fmt.Errorf("parse name: %w", err)
	}

	menuItemID, err := uuid.Parse(app.MenuItemID)
	if err != nil {
		return modifiergroupbus.NewModifierGroup{}, fmt.Errorf("parse menuItemId: %w", err)
	}

	restaurantID, err := uuid.Parse(app.RestaurantID)
	if err != nil {
		return modifiergroupbus.NewModifierGroup{}, fmt.Errorf("parse restaurantId: %w", err)
	}

	bus := modifiergroupbus.NewModifierGroup{
		MenuItemID:    menuItemID,
		RestaurantID:  restaurantID,
		Name:          nme,
		Description:   app.Description,
		MinSelections: app.MinSelections,
		MaxSelections: app.MaxSelections,
		Available:     app.Available,
		Rank:          app.Rank,
	}

	return bus, nil
}

// =============================================================================

// UpdateModifierGroup defines the data needed to update a modifier group.
type UpdateModifierGroup struct {
	Name          *string     `json:"name"`
	Description   *string     `json:"description"`
	MinSelections *int        `json:"minSelections" validate:"omitempty,gte=0,lte=1"`
	MaxSelections *int        `json:"maxSelections" validate:"omitempty,eq=1"`
	Available     *bool       `json:"available"`
	Rank          opt.NullInt `json:"rank"`
}

// Decode implements the decoder interface.
func (app *UpdateModifierGroup) Decode(data []byte) error {
	return json.Unmarshal(data, app)
}

// Validate checks the data in the model is clean.
func (app UpdateModifierGroup) Validate() error {
	if err := errs.Check(app); err != nil {
		return fmt.Errorf("validate: %w", err)
	}
	if app.Rank.Present && app.Rank.Value != nil && *app.Rank.Value < 1 {
		return errs.NewFieldErrors("rank", fmt.Errorf("rank must be >= 1"))
	}
	return nil
}

func toBusUpdateModifierGroup(app UpdateModifierGroup) (modifiergroupbus.UpdateModifierGroup, error) {
	var nme *name.Name
	if app.Name != nil {
		nm, err := name.Parse(*app.Name)
		if err != nil {
			return modifiergroupbus.UpdateModifierGroup{}, fmt.Errorf("parse name: %w", err)
		}
		nme = &nm
	}

	bus := modifiergroupbus.UpdateModifierGroup{
		Name:          nme,
		Description:   app.Description,
		MinSelections: app.MinSelections,
		MaxSelections: app.MaxSelections,
		Available:     app.Available,
		Rank:          app.Rank,
	}

	return bus, nil
}

// =============================================================================

// ReorderModifierGroups defines the payload for reordering modifier groups within a menu item.
type ReorderModifierGroups struct {
	MenuItemID string   `json:"menuItemId" validate:"required,uuid"`
	OrderedIDs []string `json:"orderedIds" validate:"required,min=1,dive,uuid"`
}

// Decode implements the decoder interface.
func (app *ReorderModifierGroups) Decode(data []byte) error {
	return json.Unmarshal(data, app)
}

// Validate checks the data in the model is clean.
func (app ReorderModifierGroups) Validate() error {
	if err := errs.Check(app); err != nil {
		return fmt.Errorf("validate: %w", err)
	}
	return nil
}
