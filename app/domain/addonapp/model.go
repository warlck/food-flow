package addonapp

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/warlck/food-flow/app/sdk/errs"
	"github.com/warlck/food-flow/business/domain/addonbus"
	"github.com/warlck/food-flow/business/types/money"
	"github.com/warlck/food-flow/business/types/name"
)

// Addon represents information about an addon for API responses.
type Addon struct {
	ID           string  `json:"id"`
	CategoryID   string  `json:"categoryId"`
	RestaurantID string  `json:"restaurantId"`
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	Price        float64 `json:"price"`
	Available    bool    `json:"available"`
	MaxQuantity  int     `json:"maxQuantity"`
	DateCreated  string  `json:"dateCreated"`
	DateUpdated  string  `json:"dateUpdated"`
}

// Encode implements the web.Encoder interface.
func (app Addon) Encode() ([]byte, string, error) {
	data, err := json.Marshal(app)
	return data, "application/json", err
}

// ToAppAddon converts a business layer addon to an app layer addon.
func ToAppAddon(bus addonbus.Addon) Addon {
	return Addon{
		ID:           bus.ID.String(),
		CategoryID:   bus.CategoryID.String(),
		RestaurantID: bus.RestaurantID.String(),
		Name:         bus.Name.String(),
		Description:  bus.Description,
		Price:        bus.Price.Value(),
		Available:    bus.Available,
		MaxQuantity:  bus.MaxQuantity,
		DateCreated:  bus.DateCreated.Format(time.RFC3339),
		DateUpdated:  bus.DateUpdated.Format(time.RFC3339),
	}
}

// ToAppAddons converts a slice of business layer addons to app layer addons.
func ToAppAddons(addons []addonbus.Addon) []Addon {
	app := make([]Addon, len(addons))
	for i, addon := range addons {
		app[i] = ToAppAddon(addon)
	}

	return app
}

// =============================================================================

// NewAddon defines the data needed to add a new addon.
type NewAddon struct {
	CategoryID   string  `json:"categoryId" validate:"required"`
	RestaurantID string  `json:"restaurantId" validate:"required"`
	Name         string  `json:"name" validate:"required"`
	Description  string  `json:"description"`
	Price        float64 `json:"price" validate:"required"`
	MaxQuantity  int     `json:"maxQuantity"`
}

// Decode implements the web.Decoder interface.
func (app *NewAddon) Decode(data []byte) error {
	return json.Unmarshal(data, app)
}

// Validate checks the data in the model is considered clean.
func (app NewAddon) Validate() error {
	if err := errs.Check(app); err != nil {
		return fmt.Errorf("validate: %w", err)
	}

	return nil
}

func toBusNewAddon(app NewAddon) (addonbus.NewAddon, error) {
	nme, err := name.Parse(app.Name)
	if err != nil {
		return addonbus.NewAddon{}, fmt.Errorf("parse name: %w", err)
	}

	categoryID, err := uuid.Parse(app.CategoryID)
	if err != nil {
		return addonbus.NewAddon{}, fmt.Errorf("parse categoryID: %w", err)
	}

	restaurantID, err := uuid.Parse(app.RestaurantID)
	if err != nil {
		return addonbus.NewAddon{}, fmt.Errorf("parse restaurantID: %w", err)
	}

	prc, err := money.Parse(app.Price)
	if err != nil {
		return addonbus.NewAddon{}, fmt.Errorf("parse price: %w", err)
	}

	bus := addonbus.NewAddon{
		CategoryID:   categoryID,
		RestaurantID: restaurantID,
		Name:         nme,
		Description:  app.Description,
		Price:        prc,
		MaxQuantity:  app.MaxQuantity,
	}

	return bus, nil
}

// =============================================================================

// UpdateAddon defines the data needed to update an addon.
type UpdateAddon struct {
	Name        *string  `json:"name"`
	Description *string  `json:"description"`
	Price       *float64 `json:"price"`
	Available   *bool    `json:"available"`
	MaxQuantity *int     `json:"maxQuantity"`
}

// Decode implements the web.Decoder interface.
func (app *UpdateAddon) Decode(data []byte) error {
	return json.Unmarshal(data, app)
}

// Validate checks the data in the model is considered clean.
func (app UpdateAddon) Validate() error {
	if err := errs.Check(app); err != nil {
		return fmt.Errorf("validate: %w", err)
	}

	return nil
}

func toBusUpdateAddon(app UpdateAddon) (addonbus.UpdateAddon, error) {
	var nme *name.Name
	if app.Name != nil {
		nm, err := name.Parse(*app.Name)
		if err != nil {
			return addonbus.UpdateAddon{}, fmt.Errorf("parse name: %w", err)
		}
		nme = &nm
	}

	var prc *money.Money
	if app.Price != nil {
		pr, err := money.Parse(*app.Price)
		if err != nil {
			return addonbus.UpdateAddon{}, fmt.Errorf("parse price: %w", err)
		}
		prc = &pr
	}

	bus := addonbus.UpdateAddon{
		Name:        nme,
		Description: app.Description,
		Price:       prc,
		Available:   app.Available,
		MaxQuantity: app.MaxQuantity,
	}

	return bus, nil
}
