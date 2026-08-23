package restaurantdb

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/warlck/food-flow/business/domain/restaurantbus"
	"github.com/warlck/food-flow/business/types/name"
)

type restaurant struct {
	ID                    uuid.UUID `db:"restaurant_id"`
	OrganizationID        uuid.UUID `db:"organization_id"`
	Name                  string    `db:"name"`
	Description           string    `db:"description"`
	Address               string    `db:"address"`
	Phone                 string    `db:"phone"`
	Email                 string    `db:"email"`
	ImageURL              string    `db:"image_url"`
	LogoURL               string    `db:"logo_url"`
	OperatingHours        string    `db:"operating_hours"`
	Enabled               bool      `db:"enabled"`
	Latitude              *float64  `db:"latitude"`
	Longitude             *float64  `db:"longitude"`
	MaxDeliveryDistanceKm float64   `db:"max_delivery_distance_km"`
	MinSpend              float64   `db:"min_spend"`
	TaxRate               float64   `db:"tax_rate"`
	DateCreated           time.Time `db:"date_created"`
	DateUpdated           time.Time `db:"date_updated"`
}

func toDBRestaurant(bus restaurantbus.Restaurant) restaurant {
	hours := bus.OperatingHours
	if len(hours) == 0 {
		hours = restaurantbus.DefaultOperatingHours()
	}
	hoursJSON, _ := json.Marshal(hours)

	return restaurant{
		ID:                    bus.ID,
		OrganizationID:        bus.OrganizationID,
		Name:                  bus.Name.String(),
		Description:           bus.Description,
		Address:               bus.Address,
		Phone:                 bus.Phone,
		Email:                 bus.Email,
		ImageURL:              bus.ImageURL,
		LogoURL:               bus.LogoURL,
		OperatingHours:        string(hoursJSON),
		Enabled:               bus.Enabled,
		Latitude:              bus.Latitude,
		Longitude:             bus.Longitude,
		MaxDeliveryDistanceKm: bus.MaxDeliveryDistanceKm,
		MinSpend:              bus.MinSpend,
		TaxRate:               bus.TaxRate,
		DateCreated:           bus.DateCreated.UTC(),
		DateUpdated:           bus.DateUpdated.UTC(),
	}
}

func toBusRestaurant(db restaurant) (restaurantbus.Restaurant, error) {
	nme, err := name.Parse(db.Name)
	if err != nil {
		return restaurantbus.Restaurant{}, fmt.Errorf("parse name: %w", err)
	}

	var hours restaurantbus.OperatingHours
	if db.OperatingHours != "" {
		_ = json.Unmarshal([]byte(db.OperatingHours), &hours)
	}
	if len(hours) == 0 {
		hours = restaurantbus.DefaultOperatingHours()
	}

	bus := restaurantbus.Restaurant{
		ID:                    db.ID,
		OrganizationID:        db.OrganizationID,
		Name:                  nme,
		Description:           db.Description,
		Address:               db.Address,
		Phone:                 db.Phone,
		Email:                 db.Email,
		ImageURL:              db.ImageURL,
		LogoURL:               db.LogoURL,
		OperatingHours:        hours,
		Enabled:               db.Enabled,
		Latitude:              db.Latitude,
		Longitude:             db.Longitude,
		MaxDeliveryDistanceKm: db.MaxDeliveryDistanceKm,
		MinSpend:              db.MinSpend,
		TaxRate:               db.TaxRate,
		DateCreated:           db.DateCreated.In(time.Local),
		DateUpdated:           db.DateUpdated.In(time.Local),
	}

	return bus, nil
}

func toBusRestaurants(dbs []restaurant) ([]restaurantbus.Restaurant, error) {
	bus := make([]restaurantbus.Restaurant, len(dbs))

	for i, db := range dbs {
		var err error
		bus[i], err = toBusRestaurant(db)
		if err != nil {
			return nil, err
		}
	}

	return bus, nil
}
