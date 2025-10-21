package menuitembus_test

import (
	"context"
	"fmt"
	"sort"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/warlck/food-flow/business/domain/categorybus"
	"github.com/warlck/food-flow/business/domain/menuitembus"
	"github.com/warlck/food-flow/business/domain/restaurantbus"
	"github.com/warlck/food-flow/business/sdk/dbtest"
	"github.com/warlck/food-flow/business/sdk/page"
	"github.com/warlck/food-flow/business/sdk/unittest"
	"github.com/warlck/food-flow/business/types/money"
	"github.com/warlck/food-flow/business/types/name"
)

func Test_MenuItem(t *testing.T) {
	t.Parallel()

	db := dbtest.New(t, "Test_MenuItem")

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
	rests, err := restaurantbus.TestSeedRestaurants(ctx, 2, busDomain.Restaurant)
	if err != nil {
		return unittest.SeedData{}, fmt.Errorf("seeding restaurants : %w", err)
	}

	// Seed categories for restaurants
	cats1, err := categorybus.TestSeedCategories(ctx, 2, rests[0].ID, busDomain.Category)
	if err != nil {
		return unittest.SeedData{}, fmt.Errorf("seeding categories : %w", err)
	}

	cats2, err := categorybus.TestSeedCategories(ctx, 2, rests[1].ID, busDomain.Category)
	if err != nil {
		return unittest.SeedData{}, fmt.Errorf("seeding categories : %w", err)
	}

	// Seed menu items for categories
	items1, err := menuitembus.TestSeedMenuItems(ctx, 2, cats1[0].ID, rests[0].ID, busDomain.MenuItem)
	if err != nil {
		return unittest.SeedData{}, fmt.Errorf("seeding menu items : %w", err)
	}

	items2, err := menuitembus.TestSeedMenuItems(ctx, 2, cats1[1].ID, rests[0].ID, busDomain.MenuItem)
	if err != nil {
		return unittest.SeedData{}, fmt.Errorf("seeding menu items : %w", err)
	}

	// -------------------------------------------------------------------------

	sd := unittest.SeedData{
		Restaurants: []unittest.Restaurant{
			{Restaurant: rests[0]},
			{Restaurant: rests[1]},
		},
		Categories: []unittest.Category{
			{Category: cats1[0]},
			{Category: cats1[1]},
			{Category: cats2[0]},
			{Category: cats2[1]},
		},
		MenuItems: []unittest.MenuItem{
			{MenuItem: items1[0]},
			{MenuItem: items1[1]},
			{MenuItem: items2[0]},
			{MenuItem: items2[1]},
		},
	}

	return sd, nil
}

// =============================================================================

func query(busDomain dbtest.BusDomain, sd unittest.SeedData) []unittest.Table {
	items := make([]menuitembus.MenuItem, 0, len(sd.MenuItems))

	for _, item := range sd.MenuItems {
		items = append(items, item.MenuItem)
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].ID.String() <= items[j].ID.String()
	})

	table := []unittest.Table{
		{
			Name:    "all",
			ExpResp: items,
			ExcFunc: func(ctx context.Context) any {
				filter := menuitembus.QueryFilter{
					Name: dbtest.NamePointer("MenuItem"),
				}

				resp, err := busDomain.MenuItem.Query(ctx, filter, menuitembus.DefaultOrderBy, page.MustParse("1", "10"))
				if err != nil {
					return err
				}

				return resp
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.([]menuitembus.MenuItem)
				if !exists {
					return "error occurred"
				}

				expResp := exp.([]menuitembus.MenuItem)

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
			ExpResp: sd.MenuItems[0].MenuItem,
			ExcFunc: func(ctx context.Context) any {
				resp, err := busDomain.MenuItem.QueryByID(ctx, sd.MenuItems[0].ID)
				if err != nil {
					return err
				}

				return resp
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.(menuitembus.MenuItem)
				if !exists {
					return "error occurred"
				}

				expResp := exp.(menuitembus.MenuItem)

				if gotResp.DateCreated.Format(time.RFC3339) == expResp.DateCreated.Format(time.RFC3339) {
					expResp.DateCreated = gotResp.DateCreated
				}

				if gotResp.DateUpdated.Format(time.RFC3339) == expResp.DateUpdated.Format(time.RFC3339) {
					expResp.DateUpdated = gotResp.DateUpdated
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
			ExpResp: menuitembus.MenuItem{
				Name:         name.MustParse("Margherita Pizza"),
				Description:  "Classic Italian pizza with fresh mozzarella",
				Price:        money.MustParse(12.99),
				CategoryID:   sd.Categories[0].ID,
				RestaurantID: sd.Restaurants[0].ID,
				ImageURL:     "pizza.jpg",
				Available:    true,
			},
			ExcFunc: func(ctx context.Context) any {
				ni := menuitembus.NewMenuItem{
					Name:         name.MustParse("Margherita Pizza"),
					Description:  "Classic Italian pizza with fresh mozzarella",
					Price:        money.MustParse(12.99),
					CategoryID:   sd.Categories[0].ID,
					RestaurantID: sd.Restaurants[0].ID,
					ImageURL:     "pizza.jpg",
				}

				resp, err := busDomain.MenuItem.Create(ctx, ni)
				if err != nil {
					return err
				}

				return resp
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.(menuitembus.MenuItem)
				if !exists {
					return "error occurred"
				}

				expResp := exp.(menuitembus.MenuItem)

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
			ExpResp: menuitembus.MenuItem{
				ID:           sd.MenuItems[0].ID,
				Name:         name.MustParse("Updated Pizza"),
				Description:  "Updated description for this delicious pizza",
				Price:        money.MustParse(15.99),
				CategoryID:   sd.Categories[1].ID,
				RestaurantID: sd.MenuItems[0].RestaurantID,
				ImageURL:     "updated.jpg",
				Available:    false,
				DateCreated:  sd.MenuItems[0].DateCreated,
			},
			ExcFunc: func(ctx context.Context) any {
				newPrice := money.MustParse(15.99)
				ui := menuitembus.UpdateMenuItem{
					Name:        dbtest.NamePointer("Updated Pizza"),
					Description: dbtest.StringPointer("Updated description for this delicious pizza"),
					Price:       &newPrice,
					CategoryID:  &sd.Categories[1].ID,
					ImageURL:    dbtest.StringPointer("updated.jpg"),
					Available:   dbtest.BoolPointer(false),
				}

				resp, err := busDomain.MenuItem.Update(ctx, sd.MenuItems[0].MenuItem, ui)
				if err != nil {
					return err
				}

				return resp
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.(menuitembus.MenuItem)
				if !exists {
					return "error occurred"
				}

				expResp := exp.(menuitembus.MenuItem)

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
				if err := busDomain.MenuItem.Delete(ctx, sd.MenuItems[1].MenuItem); err != nil {
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
