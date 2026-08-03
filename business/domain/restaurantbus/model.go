package restaurantbus

import (
	"time"

	"github.com/google/uuid"
	"github.com/warlck/food-flow/business/types/name"
)

// Restaurant represents information about a restaurant.
type Restaurant struct {
	ID                    uuid.UUID
	Name                  name.Name
	Description           string
	Address               string
	Phone                 string
	Email                 string
	ImageURL              string
	Enabled               bool
	Latitude              *float64
	Longitude             *float64
	MaxDeliveryDistanceKm float64
	DateCreated           time.Time
	DateUpdated           time.Time
}

// NewRestaurant contains information needed to create a new restaurant.
type NewRestaurant struct {
	Name                  name.Name
	Description           string
	Address               string
	Phone                 string
	Email                 string
	ImageURL              string
	Latitude              *float64
	Longitude             *float64
	MaxDeliveryDistanceKm float64
}

// UpdateRestaurant contains information needed to update a restaurant.
type UpdateRestaurant struct {
	Name                  *name.Name
	Description           *string
	Address               *string
	Phone                 *string
	Email                 *string
	ImageURL              *string
	Enabled               *bool
	Latitude              *float64
	Longitude             *float64
	MaxDeliveryDistanceKm *float64
}
