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
	"github.com/warlck/food-flow/business/types/opt"
)

// Addon represents a menu-item scoped addon for API responses.
type Addon struct {
	ID           string  `json:"id"`
	MenuItemID   string  `json:"menuItemId"`
	RestaurantID string  `json:"restaurantId"`
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	Price        float64 `json:"price"`
	Available    bool    `json:"available"`
	MaxQuantity  int     `json:"maxQuantity"`
	Rank         *int    `json:"rank,omitempty"`
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
		MenuItemID:   bus.MenuItemID.String(),
		RestaurantID: bus.RestaurantID.String(),
		Name:         bus.Name.String(),
		Description:  bus.Description,
		Price:        bus.Price.Value(),
		Available:    bus.Available,
		MaxQuantity:  bus.MaxQuantity,
		Rank:         bus.Rank,
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
	MenuItemID   string  `json:"menuItemId" validate:"required,uuid"`
	RestaurantID string  `json:"restaurantId" validate:"required,uuid"`
	Name         string  `json:"name" validate:"required"`
	Description  string  `json:"description"`
	Price        float64 `json:"price" validate:"gte=0"`
	Available    *bool   `json:"available"`
	MaxQuantity  int     `json:"maxQuantity" validate:"omitempty,gte=1"`
	Rank         *int    `json:"rank" validate:"omitempty,gte=1"`
}

// Decode implements the web.Decoder interface.
func (app *NewAddon) Decode(data []byte) error {
	return json.Unmarshal(data, app)
}

// Validate checks the data in the model is clean.
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

	price, err := money.Parse(app.Price)
	if err != nil {
		return addonbus.NewAddon{}, fmt.Errorf("parse price: %w", err)
	}

	menuItemID, err := uuid.Parse(app.MenuItemID)
	if err != nil {
		return addonbus.NewAddon{}, fmt.Errorf("parse menuItemId: %w", err)
	}

	restaurantID, err := uuid.Parse(app.RestaurantID)
	if err != nil {
		return addonbus.NewAddon{}, fmt.Errorf("parse restaurantId: %w", err)
	}

	maxQty := app.MaxQuantity
	if maxQty <= 0 {
		maxQty = 10
	}

	bus := addonbus.NewAddon{
		MenuItemID:   menuItemID,
		RestaurantID: restaurantID,
		Name:         nme,
		Description:  app.Description,
		Price:        price,
		Available:    app.Available,
		MaxQuantity:  maxQty,
		Rank:         app.Rank,
	}

	return bus, nil
}

// =============================================================================

// UpdateAddon defines the data needed to update an addon.
type UpdateAddon struct {
	Name        *string      `json:"name"`
	Description *string      `json:"description"`
	Price       *float64     `json:"price" validate:"omitempty,gte=0"`
	Available   *bool        `json:"available"`
	MaxQuantity *int         `json:"maxQuantity" validate:"omitempty,gte=1"`
	Rank        *opt.NullInt `json:"rank"`
}

// Decode implements the web.Decoder interface.
func (app *UpdateAddon) Decode(data []byte) error {
	return json.Unmarshal(data, app)
}

// Validate checks the data in the model is clean.
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

	var price *money.Money
	if app.Price != nil {
		p, err := money.Parse(*app.Price)
		if err != nil {
			return addonbus.UpdateAddon{}, fmt.Errorf("parse price: %w", err)
		}
		price = &p
	}

	var rankOpt opt.NullInt
	if app.Rank != nil {
		rankOpt = *app.Rank
	}

	bus := addonbus.UpdateAddon{
		Name:        nme,
		Description: app.Description,
		Price:       price,
		Available:   app.Available,
		MaxQuantity: app.MaxQuantity,
		Rank:        rankOpt,
	}

	return bus, nil
}

// =============================================================================

// ReorderAddons defines the payload for reordering addons on a menu item.
type ReorderAddons struct {
	MenuItemID string   `json:"menuItemId" validate:"required,uuid"`
	AddonIDs   []string `json:"addonIds" validate:"required,min=1,dive,uuid"`
}

// Decode implements the web.Decoder interface.
func (app *ReorderAddons) Decode(data []byte) error {
	return json.Unmarshal(data, app)
}

// Validate checks the data in the model is clean.
func (app ReorderAddons) Validate() error {
	if err := errs.Check(app); err != nil {
		return fmt.Errorf("validate: %w", err)
	}
	return nil
}
