package modifiergroupdb

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/warlck/food-flow/business/domain/modifiergroupbus"
	"github.com/warlck/food-flow/business/types/name"
)

type modifierGroup struct {
	ID            uuid.UUID `db:"modifier_group_id"`
	MenuItemID    uuid.UUID `db:"menu_item_id"`
	RestaurantID  uuid.UUID `db:"restaurant_id"`
	Name          string    `db:"name"`
	Description   string    `db:"description"`
	MinSelections int       `db:"min_selections"`
	MaxSelections int       `db:"max_selections"`
	Available     bool      `db:"available"`
	Rank          *int      `db:"rank"`
	DateCreated   time.Time `db:"date_created"`
	DateUpdated   time.Time `db:"date_updated"`
}

func toDBModifierGroup(bus modifiergroupbus.ModifierGroup) modifierGroup {
	return modifierGroup{
		ID:            bus.ID,
		MenuItemID:    bus.MenuItemID,
		RestaurantID:  bus.RestaurantID,
		Name:          bus.Name.String(),
		Description:   bus.Description,
		MinSelections: bus.MinSelections,
		MaxSelections: bus.MaxSelections,
		Available:     bus.Available,
		Rank:          bus.Rank,
		DateCreated:   bus.DateCreated.UTC(),
		DateUpdated:   bus.DateUpdated.UTC(),
	}
}

func toBusModifierGroup(db modifierGroup) (modifiergroupbus.ModifierGroup, error) {
	nme, err := name.Parse(db.Name)
	if err != nil {
		return modifiergroupbus.ModifierGroup{}, fmt.Errorf("parse name: %w", err)
	}

	bus := modifiergroupbus.ModifierGroup{
		ID:            db.ID,
		MenuItemID:    db.MenuItemID,
		RestaurantID:  db.RestaurantID,
		Name:          nme,
		Description:   db.Description,
		MinSelections: db.MinSelections,
		MaxSelections: db.MaxSelections,
		Available:     db.Available,
		Rank:          db.Rank,
		DateCreated:   db.DateCreated.In(time.Local),
		DateUpdated:   db.DateUpdated.In(time.Local),
	}

	return bus, nil
}

func toBusModifierGroups(dbs []modifierGroup) ([]modifiergroupbus.ModifierGroup, error) {
	bus := make([]modifiergroupbus.ModifierGroup, len(dbs))

	for i, db := range dbs {
		var err error
		bus[i], err = toBusModifierGroup(db)
		if err != nil {
			return nil, err
		}
	}

	return bus, nil
}
