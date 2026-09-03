package modifieroptiondb

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/warlck/food-flow/business/domain/modifieroptionbus"
	"github.com/warlck/food-flow/business/types/money"
	"github.com/warlck/food-flow/business/types/name"
)

type modifierOption struct {
	ID              uuid.UUID `db:"modifier_option_id"`
	ModifierGroupID uuid.UUID `db:"modifier_group_id"`
	RestaurantID    uuid.UUID `db:"restaurant_id"`
	Name            string    `db:"name"`
	Description     string    `db:"description"`
	PriceDelta      float64   `db:"price_delta"`
	Available       bool      `db:"available"`
	Rank            *int      `db:"rank"`
	DateCreated     time.Time `db:"date_created"`
	DateUpdated     time.Time `db:"date_updated"`
}

func toDBModifierOption(bus modifieroptionbus.ModifierOption) modifierOption {
	return modifierOption{
		ID:              bus.ID,
		ModifierGroupID: bus.ModifierGroupID,
		RestaurantID:    bus.RestaurantID,
		Name:            bus.Name.String(),
		Description:     bus.Description,
		PriceDelta:      bus.PriceDelta.Value(),
		Available:       bus.Available,
		Rank:            bus.Rank,
		DateCreated:     bus.DateCreated.UTC(),
		DateUpdated:     bus.DateUpdated.UTC(),
	}
}

func toBusModifierOption(db modifierOption) (modifieroptionbus.ModifierOption, error) {
	nme, err := name.Parse(db.Name)
	if err != nil {
		return modifieroptionbus.ModifierOption{}, fmt.Errorf("parse name: %w", err)
	}

	priceDelta, err := money.Parse(db.PriceDelta)
	if err != nil {
		return modifieroptionbus.ModifierOption{}, fmt.Errorf("parse priceDelta: %w", err)
	}

	bus := modifieroptionbus.ModifierOption{
		ID:              db.ID,
		ModifierGroupID: db.ModifierGroupID,
		RestaurantID:    db.RestaurantID,
		Name:            nme,
		Description:     db.Description,
		PriceDelta:      priceDelta,
		Available:       db.Available,
		Rank:            db.Rank,
		DateCreated:     db.DateCreated.In(time.Local),
		DateUpdated:     db.DateUpdated.In(time.Local),
	}

	return bus, nil
}

func toBusModifierOptions(dbs []modifierOption) ([]modifieroptionbus.ModifierOption, error) {
	bus := make([]modifieroptionbus.ModifierOption, len(dbs))

	for i, db := range dbs {
		var err error
		bus[i], err = toBusModifierOption(db)
		if err != nil {
			return nil, err
		}
	}

	return bus, nil
}
