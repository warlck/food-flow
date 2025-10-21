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
	DateCreated  string  `json:"dateCreated"`
	DateUpdated  string  `json:"dateUpdated"`
}

// Encode implements the encoder interface.
func (app MenuItem) Encode() ([]byte, string, error) {
	data, err := json.Marshal(app)
	return data, "application/json", err
}

func toAppMenuItem(bus menuitembus.MenuItem) MenuItem {
	return MenuItem{
		ID:           bus.ID.String(),
		Name:         bus.Name.String(),
		Description:  bus.Description.String(),
		Price:        bus.Price.Value(),
		CategoryID:   bus.CategoryID.String(),
		RestaurantID: bus.RestaurantID.String(),
		ImageURL:     bus.ImageURL.String(),
		Available:    bus.Available,
		DateCreated:  bus.DateCreated.Format(time.RFC3339),
		DateUpdated:  bus.DateUpdated.Format(time.RFC3339),
	}
}

func toAppMenuItems(items []menuitembus.MenuItem) []MenuItem {
	app := make([]MenuItem, len(items))
	for i, item := range items {
		app[i] = toAppMenuItem(item)
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

	description, err := name.ParseNull(app.Description)
	if err != nil {
		return menuitembus.NewMenuItem{}, fmt.Errorf("parse description: %w", err)
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

	imageURL, err := name.ParseNull(app.ImageURL)
	if err != nil {
		return menuitembus.NewMenuItem{}, fmt.Errorf("parse imageURL: %w", err)
	}

	bus := menuitembus.NewMenuItem{
		Name:         nme,
		Description:  description,
		Price:        price,
		CategoryID:   categoryID,
		RestaurantID: restaurantID,
		ImageURL:     imageURL,
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

	var description *name.Null
	if app.Description != nil {
		desc, err := name.ParseNull(*app.Description)
		if err != nil {
			return menuitembus.UpdateMenuItem{}, fmt.Errorf("parse description: %w", err)
		}
		description = &desc
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

	var imageURL *name.Null
	if app.ImageURL != nil {
		img, err := name.ParseNull(*app.ImageURL)
		if err != nil {
			return menuitembus.UpdateMenuItem{}, fmt.Errorf("parse imageURL: %w", err)
		}
		imageURL = &img
	}

	bus := menuitembus.UpdateMenuItem{
		Name:        nme,
		Description: description,
		Price:       price,
		CategoryID:  categoryID,
		ImageURL:    imageURL,
		Available:   app.Available,
	}

	return bus, nil
}
