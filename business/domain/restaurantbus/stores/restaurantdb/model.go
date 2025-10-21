package restaurantdb

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/warlck/food-flow/business/domain/restaurantbus"
	"github.com/warlck/food-flow/business/types/name"
)

type restaurant struct {
	ID          uuid.UUID `db:"restaurant_id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	Address     string    `db:"address"`
	Phone       string    `db:"phone"`
	Email       string    `db:"email"`
	ImageURL    string    `db:"image_url"`
	Enabled     bool      `db:"enabled"`
	DateCreated time.Time `db:"date_created"`
	DateUpdated time.Time `db:"date_updated"`
}

func toDBRestaurant(bus restaurantbus.Restaurant) restaurant {
	return restaurant{
		ID:          bus.ID,
		Name:        bus.Name.String(),
		Description: bus.Description,
		Address:     bus.Address,
		Phone:       bus.Phone,
		Email:       bus.Email,
		ImageURL:    bus.ImageURL,
		Enabled:     bus.Enabled,
		DateCreated: bus.DateCreated.UTC(),
		DateUpdated: bus.DateUpdated.UTC(),
	}
}

func toBusRestaurant(db restaurant) (restaurantbus.Restaurant, error) {
	nme, err := name.Parse(db.Name)
	if err != nil {
		return restaurantbus.Restaurant{}, fmt.Errorf("parse name: %w", err)
	}

	bus := restaurantbus.Restaurant{
		ID:          db.ID,
		Name:        nme,
		Description: db.Description,
		Address:     db.Address,
		Phone:       db.Phone,
		Email:       db.Email,
		ImageURL:    db.ImageURL,
		Enabled:     db.Enabled,
		DateCreated: db.DateCreated.In(time.Local),
		DateUpdated: db.DateUpdated.In(time.Local),
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
