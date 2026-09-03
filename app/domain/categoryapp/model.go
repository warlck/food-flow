package categoryapp

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/warlck/food-flow/app/sdk/errs"
	"github.com/warlck/food-flow/business/domain/categorybus"
	"github.com/warlck/food-flow/business/types/name"
	"github.com/warlck/food-flow/business/types/opt"
)

// Category represents information about a category for API responses.
type Category struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	RestaurantID string `json:"restaurantId"`
	Enabled      bool   `json:"enabled"`
	Rank         *int   `json:"rank"`
	DateCreated  string `json:"dateCreated"`
	DateUpdated  string `json:"dateUpdated"`
}

// Encode implements the encoder interface.
func (app Category) Encode() ([]byte, string, error) {
	data, err := json.Marshal(app)
	return data, "application/json", err
}

// ToAppCategory converts a business layer category to an app layer category.
func ToAppCategory(bus categorybus.Category) Category {
	return Category{
		ID:           bus.ID.String(),
		Name:         bus.Name.String(),
		Description:  bus.Description,
		RestaurantID: bus.RestaurantID.String(),
		Enabled:      bus.Enabled,
		Rank:         bus.Rank,
		DateCreated:  bus.DateCreated.Format(time.RFC3339),
		DateUpdated:  bus.DateUpdated.Format(time.RFC3339),
	}
}

// ToAppCategories converts a slice of business layer categories to app layer categories.
func ToAppCategories(categories []categorybus.Category) []Category {
	app := make([]Category, len(categories))
	for i, cat := range categories {
		app[i] = ToAppCategory(cat)
	}

	return app
}

// =============================================================================

// NewCategory defines the data needed to add a new category.
type NewCategory struct {
	Name         string `json:"name" validate:"required"`
	Description  string `json:"description"`
	RestaurantID string `json:"restaurantId" validate:"required,uuid"`
	Rank         *int   `json:"rank" validate:"omitempty,gte=1"`
}

// Decode implements the decoder interface.
func (app *NewCategory) Decode(data []byte) error {
	return json.Unmarshal(data, app)
}

// Validate checks the data in the model is considered clean.
func (app NewCategory) Validate() error {
	if err := errs.Check(app); err != nil {
		return fmt.Errorf("validate: %w", err)
	}

	return nil
}

func toBusNewCategory(app NewCategory) (categorybus.NewCategory, error) {
	nme, err := name.Parse(app.Name)
	if err != nil {
		return categorybus.NewCategory{}, fmt.Errorf("parse name: %w", err)
	}

	restaurantID, err := uuid.Parse(app.RestaurantID)
	if err != nil {
		return categorybus.NewCategory{}, fmt.Errorf("parse restaurantID: %w", err)
	}

	bus := categorybus.NewCategory{
		Name:         nme,
		Description:  app.Description,
		RestaurantID: restaurantID,
		Rank:         app.Rank,
	}

	return bus, nil
}

// =============================================================================

// UpdateCategory defines the data needed to update a category.
type UpdateCategory struct {
	Name        *string     `json:"name"`
	Description *string     `json:"description"`
	Enabled     *bool       `json:"enabled"`
	Rank        opt.NullInt `json:"rank"`
}

// Decode implements the decoder interface.
func (app *UpdateCategory) Decode(data []byte) error {
	return json.Unmarshal(data, app)
}

// Validate checks the data in the model is considered clean.
func (app UpdateCategory) Validate() error {
	if err := errs.Check(app); err != nil {
		return fmt.Errorf("validate: %w", err)
	}

	if app.Rank.Present && app.Rank.Value != nil && *app.Rank.Value < 1 {
		return errs.NewFieldErrors("rank", fmt.Errorf("rank must be 1 or greater"))
	}

	return nil
}

func toBusUpdateCategory(app UpdateCategory) (categorybus.UpdateCategory, error) {
	var nme *name.Name
	if app.Name != nil {
		nm, err := name.Parse(*app.Name)
		if err != nil {
			return categorybus.UpdateCategory{}, fmt.Errorf("parse name: %w", err)
		}
		nme = &nm
	}

	bus := categorybus.UpdateCategory{
		Name:        nme,
		Description: app.Description,
		Enabled:     app.Enabled,
		Rank:        app.Rank,
	}

	return bus, nil
}

// =============================================================================

// ReorderCategories defines the payload for reordering categories.
type ReorderCategories struct {
	RestaurantID string   `json:"restaurantId" validate:"required,uuid"`
	OrderedIDs   []string `json:"orderedIds" validate:"required,min=1,dive,uuid"`
}

// Decode implements the decoder interface.
func (app *ReorderCategories) Decode(data []byte) error {
	return json.Unmarshal(data, app)
}

// Validate checks the data in the model is considered clean.
func (app ReorderCategories) Validate() error {
	if err := errs.Check(app); err != nil {
		return fmt.Errorf("validate: %w", err)
	}

	return nil
}
