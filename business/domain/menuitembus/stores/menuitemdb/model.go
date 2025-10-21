package menuitemdb

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/warlck/food-flow/business/domain/menuitembus"
	"github.com/warlck/food-flow/business/types/money"
	"github.com/warlck/food-flow/business/types/name"
)

type menuItem struct {
	ID           uuid.UUID `db:"menu_item_id"`
	Name         string    `db:"name"`
	Description  string    `db:"description"`
	Price        float64   `db:"price"`
	CategoryID   uuid.UUID `db:"category_id"`
	RestaurantID uuid.UUID `db:"restaurant_id"`
	ImageURL     string    `db:"image_url"`
	Available    bool      `db:"available"`
	DateCreated  time.Time `db:"date_created"`
	DateUpdated  time.Time `db:"date_updated"`
}

func toDBMenuItem(bus menuitembus.MenuItem) menuItem {
	return menuItem{
		ID:           bus.ID,
		Name:         bus.Name.String(),
		Description:  bus.Description,
		Price:        bus.Price.Value(),
		CategoryID:   bus.CategoryID,
		RestaurantID: bus.RestaurantID,
		ImageURL:     bus.ImageURL,
		Available:    bus.Available,
		DateCreated:  bus.DateCreated.UTC(),
		DateUpdated:  bus.DateUpdated.UTC(),
	}
}

func toBusMenuItem(db menuItem) (menuitembus.MenuItem, error) {
	nme, err := name.Parse(db.Name)
	if err != nil {
		return menuitembus.MenuItem{}, fmt.Errorf("parse name: %w", err)
	}

	price, err := money.Parse(db.Price)
	if err != nil {
		return menuitembus.MenuItem{}, fmt.Errorf("parse price: %w", err)
	}

	bus := menuitembus.MenuItem{
		ID:           db.ID,
		Name:         nme,
		Description:  db.Description,
		Price:        price,
		CategoryID:   db.CategoryID,
		RestaurantID: db.RestaurantID,
		ImageURL:     db.ImageURL,
		Available:    db.Available,
		DateCreated:  db.DateCreated.In(time.Local),
		DateUpdated:  db.DateUpdated.In(time.Local),
	}

	return bus, nil
}

func toBusMenuItems(dbs []menuItem) ([]menuitembus.MenuItem, error) {
	bus := make([]menuitembus.MenuItem, len(dbs))

	for i, db := range dbs {
		var err error
		bus[i], err = toBusMenuItem(db)
		if err != nil {
			return nil, err
		}
	}

	return bus, nil
}
