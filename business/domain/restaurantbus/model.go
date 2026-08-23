package restaurantbus

import (
	"time"

	"github.com/google/uuid"
	"github.com/warlck/food-flow/business/types/name"
)

// DaySchedule represents open/close timings and closed status for a single day.
type DaySchedule struct {
	Open     string `json:"open"`
	Close    string `json:"close"`
	IsClosed bool   `json:"isClosed"`
}

// OperatingHours represents a restaurant's weekly schedule indexed by day of week.
type OperatingHours map[string]DaySchedule

// DefaultOperatingHours returns the standard default weekly schedule.
func DefaultOperatingHours() OperatingHours {
	return OperatingHours{
		"monday":    {Open: "10:00", Close: "22:00", IsClosed: false},
		"tuesday":   {Open: "10:00", Close: "22:00", IsClosed: false},
		"wednesday": {Open: "10:00", Close: "22:00", IsClosed: false},
		"thursday":  {Open: "10:00", Close: "22:00", IsClosed: false},
		"friday":    {Open: "10:00", Close: "23:00", IsClosed: false},
		"saturday":  {Open: "11:00", Close: "23:00", IsClosed: false},
		"sunday":    {Open: "11:00", Close: "22:00", IsClosed: false},
	}
}

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
	OperatingHours        OperatingHours
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
	OperatingHours        OperatingHours
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
	OperatingHours        *OperatingHours
	Enabled               *bool
	Latitude              *float64
	Longitude             *float64
	MaxDeliveryDistanceKm *float64
	MinSpend              *float64
	TaxRate               *float64
}

