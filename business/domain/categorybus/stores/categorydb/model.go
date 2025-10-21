package categorydb

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/warlck/food-flow/business/domain/categorybus"
	"github.com/warlck/food-flow/business/types/name"
)

type category struct {
	ID           uuid.UUID `db:"category_id"`
	Name         string    `db:"name"`
	Description  string    `db:"description"`
	RestaurantID uuid.UUID `db:"restaurant_id"`
	Enabled      bool      `db:"enabled"`
	DateCreated  time.Time `db:"date_created"`
	DateUpdated  time.Time `db:"date_updated"`
}

func toDBCategory(bus categorybus.Category) category {
	return category{
		ID:           bus.ID,
		Name:         bus.Name.String(),
		Description:  bus.Description,
		RestaurantID: bus.RestaurantID,
		Enabled:      bus.Enabled,
		DateCreated:  bus.DateCreated.UTC(),
		DateUpdated:  bus.DateUpdated.UTC(),
	}
}

func toBusCategory(db category) (categorybus.Category, error) {
	nme, err := name.Parse(db.Name)
	if err != nil {
		return categorybus.Category{}, fmt.Errorf("parse name: %w", err)
	}

	bus := categorybus.Category{
		ID:           db.ID,
		Name:         nme,
		Description:  db.Description,
		RestaurantID: db.RestaurantID,
		Enabled:      db.Enabled,
		DateCreated:  db.DateCreated.In(time.Local),
		DateUpdated:  db.DateUpdated.In(time.Local),
	}

	return bus, nil
}

func toBusCategories(dbs []category) ([]categorybus.Category, error) {
	bus := make([]categorybus.Category, len(dbs))

	for i, db := range dbs {
		var err error
		bus[i], err = toBusCategory(db)
		if err != nil {
			return nil, err
		}
	}

	return bus, nil
}
