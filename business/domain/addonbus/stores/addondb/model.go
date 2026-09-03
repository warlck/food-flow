package addondb

import (
	"time"

	"github.com/google/uuid"
	"github.com/warlck/food-flow/business/domain/addonbus"
	"github.com/warlck/food-flow/business/types/money"
	"github.com/warlck/food-flow/business/types/name"
)

// dbAddon represents the database row for an addon.
type dbAddon struct {
	ID           uuid.UUID `db:"addon_id"`
	MenuItemID   uuid.UUID `db:"menu_item_id"`
	RestaurantID uuid.UUID `db:"restaurant_id"`
	Name         string    `db:"name"`
	Description  string    `db:"description"`
	Price        float64   `db:"price"`
	Available    bool      `db:"available"`
	MaxQuantity  int       `db:"max_quantity"`
	Rank         *int      `db:"rank"`
	DateCreated  time.Time `db:"date_created"`
	DateUpdated  time.Time `db:"date_updated"`
}

func toDBAddon(bus addonbus.Addon) dbAddon {
	return dbAddon{
		ID:           bus.ID,
		MenuItemID:   bus.MenuItemID,
		RestaurantID: bus.RestaurantID,
		Name:         bus.Name.String(),
		Description:  bus.Description,
		Price:        bus.Price.Value(),
		Available:    bus.Available,
		MaxQuantity:  bus.MaxQuantity,
		Rank:         bus.Rank,
		DateCreated:  bus.DateCreated.UTC(),
		DateUpdated:  bus.DateUpdated.UTC(),
	}
}

func toBusAddon(dbo dbAddon) (addonbus.Addon, error) {
	n, err := name.Parse(dbo.Name)
	if err != nil {
		return addonbus.Addon{}, err
	}

	price, err := money.Parse(dbo.Price)
	if err != nil {
		return addonbus.Addon{}, err
	}

	return addonbus.Addon{
		ID:           dbo.ID,
		MenuItemID:   dbo.MenuItemID,
		RestaurantID: dbo.RestaurantID,
		Name:         n,
		Description:  dbo.Description,
		Price:        price,
		Available:    dbo.Available,
		MaxQuantity:  dbo.MaxQuantity,
		Rank:         dbo.Rank,
		DateCreated:  dbo.DateCreated.In(time.Local),
		DateUpdated:  dbo.DateUpdated.In(time.Local),
	}, nil
}

func toBusAddons(dbos []dbAddon) ([]addonbus.Addon, error) {
	addons := make([]addonbus.Addon, len(dbos))
	for i, dbo := range dbos {
		addon, err := toBusAddon(dbo)
		if err != nil {
			return nil, err
		}
		addons[i] = addon
	}
	return addons, nil
}
