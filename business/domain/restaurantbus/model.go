package restaurantbus

import (
	"time"

	"github.com/google/uuid"
	"github.com/warlck/food-flow/business/types/name"
)

// Restaurant represents information about a restaurant.
type Restaurant struct {
	ID                    uuid.UUID
	OrganizationID        uuid.UUID
	Name                  name.Name
	Description           string
	Address               string
	Phone                 string
	Email                 string
	ImageURL              string
	LogoURL               string
	Enabled               bool
	Latitude              *float64
	Longitude             *float64
	MaxDeliveryDistanceKm float64
	MinSpend              float64
	TaxRate               float64
	DateCreated           time.Time
	DateUpdated           time.Time
}

// NewRestaurant contains information needed to create a new restaurant.
type NewRestaurant struct {
	OrganizationID        uuid.UUID
	Name                  name.Name
	Description           string
	Address               string
	Phone                 string
	Email                 string
	ImageURL              string
	LogoURL               string
	Latitude              *float64
	Longitude             *float64
	MaxDeliveryDistanceKm float64
	MinSpend              float64
	TaxRate               float64
}

// UpdateRestaurant contains information needed to update a restaurant.
type UpdateRestaurant struct {
	Name                  *name.Name
	Description           *string
	Address               *string
	Phone                 *string
	Email                 *string
	ImageURL              *string
	LogoURL               *string
	Enabled               *bool
	Latitude              *float64
	Longitude             *float64
	MaxDeliveryDistanceKm *float64
	MinSpend              *float64
	TaxRate               *float64
}
