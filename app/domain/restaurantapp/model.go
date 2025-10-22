package restaurantapp

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/warlck/food-flow/app/sdk/errs"
	"github.com/warlck/food-flow/business/domain/restaurantbus"
	"github.com/warlck/food-flow/business/types/name"
)

// Restaurant represents information about a restaurant for API responses.
type Restaurant struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Address     string `json:"address"`
	Phone       string `json:"phone"`
	Email       string `json:"email"`
	ImageURL    string `json:"imageUrl"`
	Enabled     bool   `json:"enabled"`
	DateCreated string `json:"dateCreated"`
	DateUpdated string `json:"dateUpdated"`
}

// Encode implements the encoder interface.
func (app Restaurant) Encode() ([]byte, string, error) {
	data, err := json.Marshal(app)
	return data, "application/json", err
}

// ToAppRestaurant converts a business layer restaurant to an app layer restaurant.
func ToAppRestaurant(bus restaurantbus.Restaurant) Restaurant {
	return Restaurant{
		ID:          bus.ID.String(),
		Name:        bus.Name.String(),
		Description: bus.Description,
		Address:     bus.Address,
		Phone:       bus.Phone,
		Email:       bus.Email,
		ImageURL:    bus.ImageURL,
		Enabled:     bus.Enabled,
		DateCreated: bus.DateCreated.Format(time.RFC3339),
		DateUpdated: bus.DateUpdated.Format(time.RFC3339),
	}
}

// ToAppRestaurants converts a slice of business layer restaurants to app layer restaurants.
func ToAppRestaurants(restaurants []restaurantbus.Restaurant) []Restaurant {
	app := make([]Restaurant, len(restaurants))
	for i, res := range restaurants {
		app[i] = ToAppRestaurant(res)
	}

	return app
}

// =============================================================================

// NewRestaurant defines the data needed to add a new restaurant.
type NewRestaurant struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	Address     string `json:"address" validate:"required"`
	Phone       string `json:"phone" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	ImageURL    string `json:"imageUrl"`
}

// Decode implements the decoder interface.
func (app *NewRestaurant) Decode(data []byte) error {
	return json.Unmarshal(data, app)
}

// Validate checks the data in the model is considered clean.
func (app NewRestaurant) Validate() error {
	if err := errs.Check(app); err != nil {
		return fmt.Errorf("validate: %w", err)
	}

	return nil
}

func toBusNewRestaurant(app NewRestaurant) (restaurantbus.NewRestaurant, error) {
	nme, err := name.Parse(app.Name)
	if err != nil {
		return restaurantbus.NewRestaurant{}, fmt.Errorf("parse name: %w", err)
	}

	bus := restaurantbus.NewRestaurant{
		Name:        nme,
		Description: app.Description,
		Address:     app.Address,
		Phone:       app.Phone,
		Email:       app.Email,
		ImageURL:    app.ImageURL,
	}

	return bus, nil
}

// =============================================================================

// UpdateRestaurant defines the data needed to update a restaurant.
type UpdateRestaurant struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Address     *string `json:"address"`
	Phone       *string `json:"phone"`
	Email       *string `json:"email" validate:"omitempty,email"`
	ImageURL    *string `json:"imageUrl"`
	Enabled     *bool   `json:"enabled"`
}

// Decode implements the decoder interface.
func (app *UpdateRestaurant) Decode(data []byte) error {
	return json.Unmarshal(data, app)
}

// Validate checks the data in the model is considered clean.
func (app UpdateRestaurant) Validate() error {
	if err := errs.Check(app); err != nil {
		return fmt.Errorf("validate: %w", err)
	}

	return nil
}

func toBusUpdateRestaurant(app UpdateRestaurant) (restaurantbus.UpdateRestaurant, error) {
	var nme *name.Name
	if app.Name != nil {
		nm, err := name.Parse(*app.Name)
		if err != nil {
			return restaurantbus.UpdateRestaurant{}, fmt.Errorf("parse name: %w", err)
		}
		nme = &nm
	}

	bus := restaurantbus.UpdateRestaurant{
		Name:        nme,
		Description: app.Description,
		Address:     app.Address,
		Phone:       app.Phone,
		Email:       app.Email,
		ImageURL:    app.ImageURL,
		Enabled:     app.Enabled,
	}

	return bus, nil
}

// =============================================================================

// MenuItem represents information about a menu items embedded with restaurant categories
// that are included in the restaurant API response
type MenuItem struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	ImageURL    string  `json:"imageUrl"`
	Available   bool    `json:"available"`
}

// Category represents information about a category that is included with restaurant
// API response
type Category struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Enabled     bool       `json:"enabled"`
	MenuItems   []MenuItem `json:"mentuItems"`
}

// RestaurantWithMenuCategories represents information about restaurant including menu item categories
// of the restaurant. Each category object embeds all the menuItems that in that category
type RestaurantWithMenuCategories struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Address     string     `json:"address"`
	Phone       string     `json:"phone"`
	Email       string     `json:"email"`
	ImageURL    string     `json:"imageUrl"`
	Enabled     bool       `json:"enabled"`
	Categories  []Category `json:"categories"`
	DateCreated string     `json:"dateCreated"`
	DateUpdated string     `json:"dateUpdated"`
}

// Encode implements the encoder interface.
func (app RestaurantWithMenuCategories) Encode() ([]byte, string, error) {
	data, err := json.Marshal(app)
	return data, "application/json", err
}
