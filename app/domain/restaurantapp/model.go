package restaurantapp

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/warlck/food-flow/app/sdk/errs"
	"github.com/warlck/food-flow/business/domain/restaurantbus"
	"github.com/warlck/food-flow/business/types/name"
)

// DaySchedule represents open/close timings and closed status for a single day.
type DaySchedule struct {
	Open     string `json:"open"`
	Close    string `json:"close"`
	IsClosed bool   `json:"isClosed"`
}

// OperatingHours represents a restaurant's weekly schedule.
type OperatingHours map[string]DaySchedule

// ToAppOperatingHours converts domain operating hours to app operating hours.
func ToAppOperatingHours(bus restaurantbus.OperatingHours) OperatingHours {
	if len(bus) == 0 {
		bus = restaurantbus.DefaultOperatingHours()
	}
	app := make(OperatingHours, len(bus))
	for k, v := range bus {
		app[k] = DaySchedule{
			Open:     v.Open,
			Close:    v.Close,
			IsClosed: v.IsClosed,
		}
	}
	return app
}

func toBusOperatingHours(app OperatingHours) restaurantbus.OperatingHours {
	if len(app) == 0 {
		return restaurantbus.DefaultOperatingHours()
	}
	bus := make(restaurantbus.OperatingHours, len(app))
	for k, v := range app {
		bus[k] = restaurantbus.DaySchedule{
			Open:     v.Open,
			Close:    v.Close,
			IsClosed: v.IsClosed,
		}
	}
	return bus
}

// Restaurant represents information about a restaurant for API responses.
type Restaurant struct {
	ID                    string         `json:"id"`
	Name                  string         `json:"name"`
	Description           string         `json:"description"`
	Address               string         `json:"address"`
	Phone                 string         `json:"phone"`
	Email                 string         `json:"email"`
	ImageURL              string         `json:"imageUrl"`
	LogoURL               string         `json:"logoUrl"`
	OperatingHours        OperatingHours `json:"operatingHours"`
	Enabled               bool           `json:"enabled"`
	Latitude              *float64       `json:"latitude,omitempty"`
	Longitude             *float64       `json:"longitude,omitempty"`
	MaxDeliveryDistanceKm float64        `json:"maxDeliveryDistanceKm"`
	MinSpend              float64        `json:"minSpend"`
	TaxRate               float64        `json:"taxRate"`
	DateCreated           string         `json:"dateCreated"`
	DateUpdated           string         `json:"dateUpdated"`
}

// Encode implements the encoder interface.
func (app Restaurant) Encode() ([]byte, string, error) {
	data, err := json.Marshal(app)
	return data, "application/json", err
}

// ToAppRestaurant converts a business layer restaurant to an app layer restaurant.
func ToAppRestaurant(bus restaurantbus.Restaurant) Restaurant {
	return Restaurant{
		ID:                    bus.ID.String(),
		Name:                  bus.Name.String(),
		Description:           bus.Description,
		Address:               bus.Address,
		Phone:                 bus.Phone,
		Email:                 bus.Email,
		ImageURL:              bus.ImageURL,
		LogoURL:               bus.LogoURL,
		OperatingHours:        ToAppOperatingHours(bus.OperatingHours),
		Enabled:               bus.Enabled,
		Latitude:              bus.Latitude,
		Longitude:             bus.Longitude,
		MaxDeliveryDistanceKm: bus.MaxDeliveryDistanceKm,
		MinSpend:              bus.MinSpend,
		TaxRate:               bus.TaxRate,
		DateCreated:           bus.DateCreated.Format(time.RFC3339),
		DateUpdated:           bus.DateUpdated.Format(time.RFC3339),
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
	OrganizationID        string         `json:"organizationId" validate:"required,uuid"`
	Name                  string         `json:"name" validate:"required"`
	Description           string         `json:"description"`
	Address               string         `json:"address" validate:"required"`
	Phone                 string         `json:"phone" validate:"required"`
	Email                 string         `json:"email" validate:"required,email"`
	ImageURL              string         `json:"imageUrl"`
	LogoURL               string         `json:"logoUrl"`
	OperatingHours        OperatingHours `json:"operatingHours,omitempty"`
	Latitude              *float64       `json:"latitude" validate:"omitempty,latitude"`
	Longitude             *float64       `json:"longitude" validate:"omitempty,longitude"`
	MaxDeliveryDistanceKm float64        `json:"maxDeliveryDistanceKm" validate:"gte=0"`
	MinSpend              float64        `json:"minSpend" validate:"gte=0"`
	TaxRate               float64        `json:"taxRate" validate:"gte=0"`
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

	orgID, err := uuid.Parse(app.OrganizationID)
	if err != nil {
		return restaurantbus.NewRestaurant{}, fmt.Errorf("parse organization id: %w", err)
	}

	bus := restaurantbus.NewRestaurant{
		OrganizationID:        orgID,
		Name:                  nme,
		Description:           app.Description,
		Address:               app.Address,
		Phone:                 app.Phone,
		Email:                 app.Email,
		ImageURL:              app.ImageURL,
		LogoURL:               app.LogoURL,
		OperatingHours:        toBusOperatingHours(app.OperatingHours),
		Latitude:              app.Latitude,
		Longitude:             app.Longitude,
		MaxDeliveryDistanceKm: app.MaxDeliveryDistanceKm,
		MinSpend:              app.MinSpend,
		TaxRate:               app.TaxRate,
	}

	return bus, nil
}

// =============================================================================

// UpdateRestaurant defines the data needed to update a restaurant.
type UpdateRestaurant struct {
	Name                  *string         `json:"name"`
	Description           *string         `json:"description"`
	Address               *string         `json:"address"`
	Phone                 *string         `json:"phone"`
	Email                 *string         `json:"email" validate:"omitempty,email"`
	ImageURL              *string         `json:"imageUrl"`
	LogoURL               *string         `json:"logoUrl"`
	OperatingHours        *OperatingHours `json:"operatingHours,omitempty"`
	Enabled               *bool           `json:"enabled"`
	Latitude              *float64        `json:"latitude" validate:"omitempty,latitude"`
	Longitude             *float64        `json:"longitude" validate:"omitempty,longitude"`
	MaxDeliveryDistanceKm *float64        `json:"maxDeliveryDistanceKm" validate:"omitempty,gte=0"`
	MinSpend              *float64        `json:"minSpend" validate:"omitempty,gte=0"`
	TaxRate               *float64        `json:"taxRate" validate:"omitempty,gte=0"`
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

	var hours *restaurantbus.OperatingHours
	if app.OperatingHours != nil {
		h := toBusOperatingHours(*app.OperatingHours)
		hours = &h
	}

	bus := restaurantbus.UpdateRestaurant{
		Name:                  nme,
		Description:           app.Description,
		Address:               app.Address,
		Phone:                 app.Phone,
		Email:                 app.Email,
		ImageURL:              app.ImageURL,
		LogoURL:               app.LogoURL,
		OperatingHours:        hours,
		Enabled:               app.Enabled,
		Latitude:              app.Latitude,
		Longitude:             app.Longitude,
		MaxDeliveryDistanceKm: app.MaxDeliveryDistanceKm,
		MinSpend:              app.MinSpend,
		TaxRate:               app.TaxRate,
	}

	return bus, nil
}

// =============================================================================

// Option represents information about a modifier option in the restaurant API response.
type Option struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	PriceDelta  float64 `json:"priceDelta"`
	Available   bool    `json:"available"`
	Rank        *int    `json:"rank,omitempty"`
}

// ModifierGroup represents information about a modifier group in the restaurant API response.
type ModifierGroup struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	MinSelections int      `json:"minSelections"`
	MaxSelections int      `json:"maxSelections"`
	Available     bool     `json:"available"`
	Rank          *int     `json:"rank,omitempty"`
	Options       []Option `json:"options"`
}

// Addon represents information about an addon that is included with menu items
// in the restaurant API response.
type Addon struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Available   bool    `json:"available"`
	MaxQuantity int     `json:"maxQuantity"`
	Rank        *int    `json:"rank,omitempty"`
}

// MenuItem represents information about a menu item embedded with restaurant categories
// in the restaurant API response.
type MenuItem struct {
	ID             string          `json:"id"`
	Name           string          `json:"name"`
	Description    string          `json:"description"`
	Price          float64         `json:"price"`
	ImageURL       string          `json:"imageUrl"`
	Available      bool            `json:"available"`
	Orderable      bool            `json:"orderable"`
	Rank           *int            `json:"rank,omitempty"`
	ModifierGroups []ModifierGroup `json:"modifierGroups"`
	Addons         []Addon         `json:"addons"`
}

// Category represents information about a category that is included with restaurant
// API response.
type Category struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Enabled     bool       `json:"enabled"`
	Rank        *int       `json:"rank,omitempty"`
	MenuItems   []MenuItem `json:"menuItems"`
}

// RestaurantWithMenuCategories represents information about restaurant including menu item categories
// of the restaurant. Each category object embeds all the menuItems in that category.
type RestaurantWithMenuCategories struct {
	ID                    string         `json:"id"`
	Name                  string         `json:"name"`
	Description           string         `json:"description"`
	Address               string         `json:"address"`
	Phone                 string         `json:"phone"`
	Email                 string         `json:"email"`
	ImageURL              string         `json:"imageUrl"`
	LogoURL               string         `json:"logoUrl"`
	OperatingHours        OperatingHours `json:"operatingHours"`
	Enabled               bool           `json:"enabled"`
	Latitude              *float64       `json:"latitude,omitempty"`
	Longitude             *float64       `json:"longitude,omitempty"`
	MaxDeliveryDistanceKm float64        `json:"maxDeliveryDistanceKm"`
	MinSpend              float64        `json:"minSpend"`
	TaxRate               float64        `json:"taxRate"`
	Categories            []Category     `json:"categories"`
	DateCreated           string         `json:"dateCreated"`
	DateUpdated           string         `json:"dateUpdated"`
}

// Encode implements the encoder interface.
func (app RestaurantWithMenuCategories) Encode() ([]byte, string, error) {
	data, err := json.Marshal(app)
	return data, "application/json", err
}
