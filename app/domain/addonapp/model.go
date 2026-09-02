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

// Addon represents an addon definition for API responses.
type Addon struct {
	ID           string  `json:"id"`
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

// MenuItemAddon represents an assigned addon on a menu item.
type MenuItemAddon struct {
	ID           string  `json:"id"`
	RestaurantID string  `json:"restaurantId"`
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	Price        float64 `json:"price"`
	Available    bool    `json:"available"`
	MaxQuantity  int     `json:"maxQuantity"`
	Rank         *int    `json:"rank"`
	DateCreated  string  `json:"dateCreated"`
	DateUpdated  string  `json:"dateUpdated"`
}

// Encode implements the web.Encoder interface.
func (app MenuItemAddon) Encode() ([]byte, string, error) {
	data, err := json.Marshal(app)
	return data, "application/json", err
}

// ToAppMenuItemAddon converts a MenuItemAddonInfo to an app MenuItemAddon.
func ToAppMenuItemAddon(bus addonbus.MenuItemAddonInfo) MenuItemAddon {
	return MenuItemAddon{
		ID:           bus.Addon.ID.String(),
		RestaurantID: bus.Addon.RestaurantID.String(),
		Name:         bus.Addon.Name.String(),
		Description:  bus.Addon.Description,
		Price:        bus.Addon.Price.Value(),
		Available:    bus.Addon.Available,
		MaxQuantity:  bus.Addon.MaxQuantity,
		Rank:         bus.Rank,
		DateCreated:  bus.Addon.DateCreated.Format(time.RFC3339),
		DateUpdated:  bus.Addon.DateUpdated.Format(time.RFC3339),
	}
}

// ToAppMenuItemAddons converts a slice of MenuItemAddonInfo to app MenuItemAddon.
func ToAppMenuItemAddons(infos []addonbus.MenuItemAddonInfo) []MenuItemAddon {
	app := make([]MenuItemAddon, len(infos))
	for i, info := range infos {
		app[i] = ToAppMenuItemAddon(info)
	}
	return app
}

// =============================================================================

// NewAddon defines the data needed to add a new addon definition.
type NewAddon struct {
	RestaurantID string  `json:"restaurantId" validate:"required,uuid"`
	Name         string  `json:"name" validate:"required"`
	Description  string  `json:"description"`
	Price        float64 `json:"price" validate:"gte=0"`
	Available    *bool   `json:"available"`
	MaxQuantity  int     `json:"maxQuantity" validate:"omitempty,gte=1"`
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

	restaurantID, err := uuid.Parse(app.RestaurantID)
	if err != nil {
		return addonbus.NewAddon{}, fmt.Errorf("parse restaurantID: %w", err)
	}

	maxQty := app.MaxQuantity
	if maxQty <= 0 {
		maxQty = 10
	}

	bus := addonbus.NewAddon{
		RestaurantID: restaurantID,
		Name:         nme,
		Description:  app.Description,
		Price:        price,
		Available:    app.Available,
		MaxQuantity:  maxQty,
	}

	return bus, nil
}

// =============================================================================

// UpdateAddon defines the data needed to update an addon definition.
type UpdateAddon struct {
	Name        *string  `json:"name"`
	Description *string  `json:"description"`
	Price       *float64 `json:"price" validate:"omitempty,gte=0"`
	Available   *bool    `json:"available"`
	MaxQuantity *int     `json:"maxQuantity" validate:"omitempty,gte=1"`
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

	bus := addonbus.UpdateAddon{
		Name:        nme,
		Description: app.Description,
		Price:       price,
		Available:   app.Available,
		MaxQuantity: app.MaxQuantity,
	}

	return bus, nil
}

// =============================================================================

// ItemAddonAssignmentInput represents an addon assigned to a menu item with an optional rank.
type ItemAddonAssignmentInput struct {
	AddonID string `json:"addonId" validate:"required,uuid"`
	Rank    *int   `json:"rank" validate:"omitempty,gte=1"`
}

// ReplaceMenuItemAddons defines the payload for replacing assigned addons for a menu item.
type ReplaceMenuItemAddons struct {
	Addons []ItemAddonAssignmentInput `json:"addons" validate:"required,dive"`
}

// Decode implements the web.Decoder interface.
func (app *ReplaceMenuItemAddons) Decode(data []byte) error {
	return json.Unmarshal(data, app)
}

// Validate checks the data in the model is clean.
func (app ReplaceMenuItemAddons) Validate() error {
	if err := errs.Check(app); err != nil {
		return fmt.Errorf("validate: %w", err)
	}
	return nil
}

// ReorderMenuItemAddons defines the payload for reordering assigned addons on a menu item.
type ReorderMenuItemAddons struct {
	OrderedIDs []string `json:"orderedIds" validate:"required,min=1,dive,uuid"`
}

// Decode implements the web.Decoder interface.
func (app *ReorderMenuItemAddons) Decode(data []byte) error {
	return json.Unmarshal(data, app)
}

// Validate checks the data in the model is clean.
func (app ReorderMenuItemAddons) Validate() error {
	if err := errs.Check(app); err != nil {
		return fmt.Errorf("validate: %w", err)
	}
	return nil
}
