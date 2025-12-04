package addonbus_test

import (
	"context"
	"fmt"
	"sort"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/warlck/food-flow/business/domain/addonbus"
	"github.com/warlck/food-flow/business/domain/categorybus"
	"github.com/warlck/food-flow/business/domain/menuitembus"
	"github.com/warlck/food-flow/business/domain/restaurantbus"
	"github.com/warlck/food-flow/business/sdk/dbtest"
	"github.com/warlck/food-flow/business/sdk/page"
	"github.com/warlck/food-flow/business/sdk/unittest"
	"github.com/warlck/food-flow/business/types/money"
	"github.com/warlck/food-flow/business/types/name"
)

func Test_Addon(t *testing.T) {
	t.Parallel()

	db := dbtest.New(t, "Test_Addon")

	sd, err := insertSeedData(db.BusDomain)
	if err != nil {
		t.Fatalf("Seeding error: %s", err)
	}

	// -------------------------------------------------------------------------

	unittest.Run(t, query(db.BusDomain, sd), "query")
	unittest.Run(t, create(db.BusDomain, sd), "create")
	unittest.Run(t, update(db.BusDomain, sd), "update")
	unittest.Run(t, delete(db.BusDomain, sd), "delete")
}

// =============================================================================

func insertSeedData(busDomain dbtest.BusDomain) (unittest.SeedData, error) {
	ctx := context.Background()

	// Seed restaurants first
	rests, err := restaurantbus.TestSeedRestaurants(ctx, 1, busDomain.Restaurant)
	if err != nil {
		return unittest.SeedData{}, fmt.Errorf("seeding restaurants : %w", err)
	}

	// Seed categories for restaurant
	cats, err := categorybus.TestSeedCategories(ctx, 1, rests[0].ID, busDomain.Category)
	if err != nil {
		return unittest.SeedData{}, fmt.Errorf("seeding categories : %w", err)
	}

	// Seed menu items for category
	items, err := menuitembus.TestSeedMenuItems(ctx, 2, cats[0].ID, rests[0].ID, busDomain.MenuItem)
	if err != nil {
		return unittest.SeedData{}, fmt.Errorf("seeding menu items : %w", err)
	}

	// Seed addons for menu items
	addons1, err := addonbus.TestSeedAddons(ctx, 2, items[0].ID, rests[0].ID, busDomain.Addon)
	if err != nil {
		return unittest.SeedData{}, fmt.Errorf("seeding addons for item 1 : %w", err)
	}

	addons2, err := addonbus.TestSeedAddons(ctx, 2, items[1].ID, rests[0].ID, busDomain.Addon)
	if err != nil {
		return unittest.SeedData{}, fmt.Errorf("seeding addons for item 2 : %w", err)
	}

	// -------------------------------------------------------------------------

	sd := unittest.SeedData{
		Restaurants: []unittest.Restaurant{
			{Restaurant: rests[0]},
		},
		Categories: []unittest.Category{
			{Category: cats[0]},
		},
		MenuItems: []unittest.MenuItem{
			{MenuItem: items[0]},
			{MenuItem: items[1]},
		},
		Addons: []unittest.Addon{
			{Addon: addons1[0]},
			{Addon: addons1[1]},
			{Addon: addons2[0]},
			{Addon: addons2[1]},
		},
	}

	return sd, nil
}

// =============================================================================

func query(busDomain dbtest.BusDomain, sd unittest.SeedData) []unittest.Table {
	addons := make([]addonbus.Addon, 0, len(sd.Addons))

	for _, addon := range sd.Addons {
		addons = append(addons, addon.Addon)
	}

	sort.Slice(addons, func(i, j int) bool {
		return addons[i].ID.String() <= addons[j].ID.String()
	})

	table := []unittest.Table{
		{
			Name:    "all",
			ExpResp: addons,
			ExcFunc: func(ctx context.Context) any {
				filter := addonbus.QueryFilter{
					RestaurantID: &sd.Restaurants[0].ID,
				}

				resp, err := busDomain.Addon.Query(ctx, filter, addonbus.DefaultOrderBy, page.MustParse("1", "10"))
				if err != nil {
					return err
				}

				return resp
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.([]addonbus.Addon)
				if !exists {
					return "error occurred"
				}

				expResp := exp.([]addonbus.Addon)

				sort.Slice(gotResp, func(i, j int) bool {
					return gotResp[i].ID.String() <= gotResp[j].ID.String()
				})

				for i := range gotResp {
					if gotResp[i].DateCreated.Format(time.RFC3339) == expResp[i].DateCreated.Format(time.RFC3339) {
						expResp[i].DateCreated = gotResp[i].DateCreated
					}

					if gotResp[i].DateUpdated.Format(time.RFC3339) == expResp[i].DateUpdated.Format(time.RFC3339) {
						expResp[i].DateUpdated = gotResp[i].DateUpdated
					}
				}

				return cmp.Diff(gotResp, expResp)
			},
		},
		{
			Name:    "byid",
			ExpResp: sd.Addons[0].Addon,
			ExcFunc: func(ctx context.Context) any {
				resp, err := busDomain.Addon.QueryByID(ctx, sd.Addons[0].ID)
				if err != nil {
					return err
				}

				return resp
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.(addonbus.Addon)
				if !exists {
					return "error occurred"
				}

				expResp := exp.(addonbus.Addon)

				if gotResp.DateCreated.Format(time.RFC3339) == expResp.DateCreated.Format(time.RFC3339) {
					expResp.DateCreated = gotResp.DateCreated
				}

				if gotResp.DateUpdated.Format(time.RFC3339) == expResp.DateUpdated.Format(time.RFC3339) {
					expResp.DateUpdated = gotResp.DateUpdated
				}

				return cmp.Diff(gotResp, expResp)
			},
		},
		{
			Name: "bymenuitemid",
			ExpResp: []addonbus.Addon{
				sd.Addons[0].Addon,
				sd.Addons[1].Addon,
			},
			ExcFunc: func(ctx context.Context) any {
				resp, err := busDomain.Addon.QueryByMenuItemID(ctx, sd.MenuItems[0].ID)
				if err != nil {
					return err
				}

				return resp
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.([]addonbus.Addon)
				if !exists {
					return "error occurred"
				}

				expResp := exp.([]addonbus.Addon)

				sort.Slice(gotResp, func(i, j int) bool {
					return gotResp[i].ID.String() <= gotResp[j].ID.String()
				})

				sort.Slice(expResp, func(i, j int) bool {
					return expResp[i].ID.String() <= expResp[j].ID.String()
				})

				for i := range gotResp {
					if gotResp[i].DateCreated.Format(time.RFC3339) == expResp[i].DateCreated.Format(time.RFC3339) {
						expResp[i].DateCreated = gotResp[i].DateCreated
					}

					if gotResp[i].DateUpdated.Format(time.RFC3339) == expResp[i].DateUpdated.Format(time.RFC3339) {
						expResp[i].DateUpdated = gotResp[i].DateUpdated
					}
				}

				return cmp.Diff(gotResp, expResp)
			},
		},
	}

	return table
}

func create(busDomain dbtest.BusDomain, sd unittest.SeedData) []unittest.Table {
	table := []unittest.Table{
		{
			Name: "basic",
			ExpResp: addonbus.Addon{
				MenuItemID:   sd.MenuItems[0].ID,
				RestaurantID: sd.Restaurants[0].ID,
				Name:         name.MustParse("Extra Bacon"),
				Description:  "Crispy bacon strips",
				Price:        money.MustParse(3.50),
				Available:    true,
				MaxQuantity:  2,
			},
			ExcFunc: func(ctx context.Context) any {
				na := addonbus.NewAddon{
					MenuItemID:   sd.MenuItems[0].ID,
					RestaurantID: sd.Restaurants[0].ID,
					Name:         name.MustParse("Extra Bacon"),
					Description:  "Crispy bacon strips",
					Price:        money.MustParse(3.50),
					MaxQuantity:  2,
				}

				resp, err := busDomain.Addon.Create(ctx, na)
				if err != nil {
					return err
				}

				return resp
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.(addonbus.Addon)
				if !exists {
					return "error occurred"
				}

				expResp := exp.(addonbus.Addon)

				expResp.ID = gotResp.ID
				expResp.DateCreated = gotResp.DateCreated
				expResp.DateUpdated = gotResp.DateUpdated

				return cmp.Diff(gotResp, expResp)
			},
		},
		{
			Name: "default_max_quantity",
			ExpResp: addonbus.Addon{
				MenuItemID:   sd.MenuItems[0].ID,
				RestaurantID: sd.Restaurants[0].ID,
				Name:         name.MustParse("Extra Sauce"),
				Description:  "Additional sauce portion",
				Price:        money.MustParse(1.00),
				Available:    true,
				MaxQuantity:  10,
			},
			ExcFunc: func(ctx context.Context) any {
				na := addonbus.NewAddon{
					MenuItemID:   sd.MenuItems[0].ID,
					RestaurantID: sd.Restaurants[0].ID,
					Name:         name.MustParse("Extra Sauce"),
					Description:  "Additional sauce portion",
					Price:        money.MustParse(1.00),
					MaxQuantity:  0,
				}

				resp, err := busDomain.Addon.Create(ctx, na)
				if err != nil {
					return err
				}

				return resp
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.(addonbus.Addon)
				if !exists {
					return "error occurred"
				}

				expResp := exp.(addonbus.Addon)

				expResp.ID = gotResp.ID
				expResp.DateCreated = gotResp.DateCreated
				expResp.DateUpdated = gotResp.DateUpdated

				return cmp.Diff(gotResp, expResp)
			},
		},
	}

	return table
}

func update(busDomain dbtest.BusDomain, sd unittest.SeedData) []unittest.Table {
	table := []unittest.Table{
		{
			Name: "basic",
			ExpResp: addonbus.Addon{
				ID:           sd.Addons[0].ID,
				MenuItemID:   sd.Addons[0].MenuItemID,
				RestaurantID: sd.Addons[0].RestaurantID,
				Name:         name.MustParse("Updated Addon"),
				Description:  "Updated description for this addon",
				Price:        money.MustParse(5.99),
				Available:    false,
				MaxQuantity:  5,
				DateCreated:  sd.Addons[0].DateCreated,
			},
			ExcFunc: func(ctx context.Context) any {
				newName := name.MustParse("Updated Addon")
				newPrice := money.MustParse(5.99)
				newDesc := "Updated description for this addon"
				newAvailable := false
				newMaxQty := 5

				ua := addonbus.UpdateAddon{
					Name:        &newName,
					Description: &newDesc,
					Price:       &newPrice,
					Available:   &newAvailable,
					MaxQuantity: &newMaxQty,
				}

				resp, err := busDomain.Addon.Update(ctx, sd.Addons[0].Addon, ua)
				if err != nil {
					return err
				}

				return resp
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.(addonbus.Addon)
				if !exists {
					return "error occurred"
				}

				expResp := exp.(addonbus.Addon)

				expResp.DateUpdated = gotResp.DateUpdated

				return cmp.Diff(gotResp, expResp)
			},
		},
		{
			Name: "partial_update",
			ExpResp: addonbus.Addon{
				ID:           sd.Addons[1].ID,
				MenuItemID:   sd.Addons[1].MenuItemID,
				RestaurantID: sd.Addons[1].RestaurantID,
				Name:         sd.Addons[1].Name,
				Description:  sd.Addons[1].Description,
				Price:        money.MustParse(7.99),
				Available:    sd.Addons[1].Available,
				MaxQuantity:  sd.Addons[1].MaxQuantity,
				DateCreated:  sd.Addons[1].DateCreated,
			},
			ExcFunc: func(ctx context.Context) any {
				newPrice := money.MustParse(7.99)

				ua := addonbus.UpdateAddon{
					Price: &newPrice,
				}

				resp, err := busDomain.Addon.Update(ctx, sd.Addons[1].Addon, ua)
				if err != nil {
					return err
				}

				return resp
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.(addonbus.Addon)
				if !exists {
					return "error occurred"
				}

				expResp := exp.(addonbus.Addon)

				expResp.DateUpdated = gotResp.DateUpdated

				return cmp.Diff(gotResp, expResp)
			},
		},
	}

	return table
}

func delete(busDomain dbtest.BusDomain, sd unittest.SeedData) []unittest.Table {
	table := []unittest.Table{
		{
			Name:    "basic",
			ExpResp: nil,
			ExcFunc: func(ctx context.Context) any {
				if err := busDomain.Addon.Delete(ctx, sd.Addons[2].Addon); err != nil {
					return err
				}

				return nil
			},
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
	}

	return table
}
