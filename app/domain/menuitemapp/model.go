package menuitemapp

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/warlck/food-flow/app/sdk/errs"
	"github.com/warlck/food-flow/business/domain/menuitembus"
	"github.com/warlck/food-flow/business/types/money"
	"github.com/warlck/food-flow/business/types/name"
)

// MenuItem represents information about a menu item for API responses.
type MenuItem struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	Price        float64 `json:"price"`
	CategoryID   string  `json:"categoryId"`
	RestaurantID string  `json:"restaurantId"`
	ImageURL     string  `json:"imageUrl"`
	Available    bool    `json:"available"`
	Rank         *int    `json:"rank"`
	DateCreated  string  `json:"dateCreated"`
	DateUpdated  string  `json:"dateUpdated"`
}

// Encode implements the encoder interface.
func (app MenuItem) Encode() ([]byte, string, error) {
	data, err := json.Marshal(app)
	return data, "application/json", err
}

// ToAppMenuItem converts a business layer menu item to an app layer menu item.
func ToAppMenuItem(bus menuitembus.MenuItem) MenuItem {
	return MenuItem{
		ID:           bus.ID.String(),
		Name:         bus.Name.String(),
		Description:  bus.Description,
		Price:        bus.Price.Value(),
		CategoryID:   bus.CategoryID.String(),
		RestaurantID: bus.RestaurantID.String(),
		ImageURL:     bus.ImageURL,
		Available:    bus.Available,
		Rank:         bus.Rank,
		DateCreated:  bus.DateCreated.Format(time.RFC3339),
		DateUpdated:  bus.DateUpdated.Format(time.RFC3339),
	}
}

// ToAppMenuItems converts a slice of business layer menu items to app layer menu items.
func ToAppMenuItems(items []menuitembus.MenuItem) []MenuItem {
	app := make([]MenuItem, len(items))
	for i, item := range items {
		app[i] = ToAppMenuItem(item)
	}

	return app
}

// =============================================================================

// NewMenuItem defines the data needed to add a new menu item.
type NewMenuItem struct {
	Name         string  `json:"name" validate:"required"`
	Description  string  `json:"description"`
	Price        float64 `json:"price" validate:"required,gt=0"`
	CategoryID   string  `json:"categoryId" validate:"required"`
	RestaurantID string  `json:"restaurantId" validate:"required"`
	ImageURL     string  `json:"imageUrl"`
	Rank         *int    `json:"rank"`
}

// Decode implements the decoder interface.
func (app *NewMenuItem) Decode(data []byte) error {
	return json.Unmarshal(data, app)
}

// Validate checks the data in the model is considered clean.
func (app NewMenuItem) Validate() error {
	if err := errs.Check(app); err != nil {
		return fmt.Errorf("validate: %w", err)
	}

	return nil
}

func toBusNewMenuItem(app NewMenuItem) (menuitembus.NewMenuItem, error) {
	nme, err := name.Parse(app.Name)
	if err != nil {
		return menuitembus.NewMenuItem{}, fmt.Errorf("parse name: %w", err)
	}

	price, err := money.Parse(app.Price)
	if err != nil {
		return menuitembus.NewMenuItem{}, fmt.Errorf("parse price: %w", err)
	}

	categoryID, err := uuid.Parse(app.CategoryID)
	if err != nil {
		return menuitembus.NewMenuItem{}, fmt.Errorf("parse categoryID: %w", err)
	}

	restaurantID, err := uuid.Parse(app.RestaurantID)
	if err != nil {
		return menuitembus.NewMenuItem{}, fmt.Errorf("parse restaurantID: %w", err)
	}

	bus := menuitembus.NewMenuItem{
		Name:         nme,
		Description:  app.Description,
		Price:        price,
		CategoryID:   categoryID,
		RestaurantID: restaurantID,
		ImageURL:     app.ImageURL,
		Rank:         app.Rank,
	}

	return bus, nil
}

// =============================================================================

// UpdateMenuItem defines the data needed to update a menu item.
type UpdateMenuItem struct {
	Name        *string  `json:"name"`
	Description *string  `json:"description"`
	Price       *float64 `json:"price" validate:"omitempty,gt=0"`
	CategoryID  *string  `json:"categoryId"`
	ImageURL    *string  `json:"imageUrl"`
	Available   *bool    `json:"available"`
	Rank        *int     `json:"rank"`
}

// Decode implements the decoder interface.
func (app *UpdateMenuItem) Decode(data []byte) error {
	return json.Unmarshal(data, app)
}

// Validate checks the data in the model is considered clean.
func (app UpdateMenuItem) Validate() error {
	if err := errs.Check(app); err != nil {
		return fmt.Errorf("validate: %w", err)
	}

	return nil
}

func toBusUpdateMenuItem(app UpdateMenuItem) (menuitembus.UpdateMenuItem, error) {
	var nme *name.Name
	if app.Name != nil {
		nm, err := name.Parse(*app.Name)
		if err != nil {
			return menuitembus.UpdateMenuItem{}, fmt.Errorf("parse name: %w", err)
		}
		nme = &nm
	}

	var price *money.Money
	if app.Price != nil {
		p, err := money.Parse(*app.Price)
		if err != nil {
			return menuitembus.UpdateMenuItem{}, fmt.Errorf("parse price: %w", err)
		}
		price = &p
	}

	var categoryID *uuid.UUID
	if app.CategoryID != nil {
		catID, err := uuid.Parse(*app.CategoryID)
		if err != nil {
			return menuitembus.UpdateMenuItem{}, fmt.Errorf("parse categoryID: %w", err)
		}
		categoryID = &catID
	}

	bus := menuitembus.UpdateMenuItem{
		Name:        nme,
		Description: app.Description,
		Price:       price,
		CategoryID:  categoryID,
		ImageURL:    app.ImageURL,
		Available:   app.Available,
		Rank:        app.Rank,
	}

	return bus, nil
}

// =============================================================================

// ReorderMenuItems defines the data needed to update menu items order.
type ReorderMenuItems struct {
	CategoryID string   `json:"categoryId" validate:"required,uuid"`
	OrderedIDs []string `json:"orderedIds" validate:"required,min=1,dive,uuid"`
}

// Decode implements the decoder interface.
func (app *ReorderMenuItems) Decode(data []byte) error {
	return json.Unmarshal(data, app)
}

// Validate checks the data in the model is considered clean.
func (app ReorderMenuItems) Validate() error {
	if err := errs.Check(app); err != nil {
		return fmt.Errorf("validate: %w", err)
	}

	return nil
}
