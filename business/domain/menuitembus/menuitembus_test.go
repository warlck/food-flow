package menuitembus_test

import (
	"context"
	"errors"
	"fmt"
	"github.com/warlck/food-flow/business/domain/organizationbus"
	"sort"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
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
	unittest.Run(t, reorder(db.BusDomain, sd), "reorder")
}

// =============================================================================

func insertSeedData(busDomain dbtest.BusDomain) (unittest.SeedData, error) {
	ctx := context.Background()

	// Seed restaurants first
	orgs, err := organizationbus.TestSeedOrganizations(ctx, 1, busDomain.Organization)
	if err != nil {
		return unittest.SeedData{}, fmt.Errorf("seeding organizations: %w", err)
	}
	rests, err := restaurantbus.TestSeedRestaurants(ctx, 2, busDomain.Restaurant, orgs[0].ID)
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
		{
			Name: "with_rank",
			ExpResp: menuitembus.MenuItem{
				Name:         name.MustParse("Pepperoni Pizza"),
				Description:  "Spicy pepperoni with mozzarella",
				Price:        money.MustParse(14.99),
				CategoryID:   sd.Categories[0].ID,
				RestaurantID: sd.Restaurants[0].ID,
				ImageURL:     "pepperoni.jpg",
				Available:    true,
				Rank:         dbtest.IntPointer(10),
			},
			ExcFunc: func(ctx context.Context) any {
				ni := menuitembus.NewMenuItem{
					Name:         name.MustParse("Pepperoni Pizza"),
					Description:  "Spicy pepperoni with mozzarella",
					Price:        money.MustParse(14.99),
					CategoryID:   sd.Categories[0].ID,
					RestaurantID: sd.Restaurants[0].ID,
					ImageURL:     "pepperoni.jpg",
					Rank:         dbtest.IntPointer(10),
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
		{
			Name: "update_rank",
			ExpResp: menuitembus.MenuItem{
				ID:           sd.MenuItems[0].ID,
				Name:         sd.MenuItems[0].Name,
				Description:  sd.MenuItems[0].Description,
				Price:        sd.MenuItems[0].Price,
				CategoryID:   sd.MenuItems[0].CategoryID,
				RestaurantID: sd.MenuItems[0].RestaurantID,
				ImageURL:     sd.MenuItems[0].ImageURL,
				Available:    sd.MenuItems[0].Available,
				Rank:         dbtest.IntPointer(42),
				DateCreated:  sd.MenuItems[0].DateCreated,
			},
			ExcFunc: func(ctx context.Context) any {
				ui := menuitembus.UpdateMenuItem{
					Rank: dbtest.IntPointer(42),
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

func reorder(busDomain dbtest.BusDomain, sd unittest.SeedData) []unittest.Table {
	categoryID := sd.Categories[1].ID
	// In Category 1, items are items2[0] and items2[1], which correspond to sd.MenuItems[2] and sd.MenuItems[3]
	item1 := sd.MenuItems[2]
	item2 := sd.MenuItems[3]

	table := []unittest.Table{
		{
			Name:    "mismatch_length",
			ExpResp: "invalid menu item order: orderedIds must contain all menu items in the category exactly once",
			ExcFunc: func(ctx context.Context) any {
				err := busDomain.MenuItem.Reorder(ctx, categoryID, []uuid.UUID{item1.ID})
				if err != nil {
					return err.Error()
				}
				return nil
			},
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
		{
			Name:    "invalid_id",
			ExpResp: "invalid menu item order: orderedIds contains invalid or duplicate menu item id",
			ExcFunc: func(ctx context.Context) any {
				err := busDomain.MenuItem.Reorder(ctx, categoryID, []uuid.UUID{item1.ID, uuid.New()})
				if err != nil {
					return err.Error()
				}
				return nil
			},
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
		{
			Name:    "invalid_order_sentinel",
			ExpResp: true,
			ExcFunc: func(ctx context.Context) any {
				// Validation failures must wrap ErrInvalidOrder so the API
				// layer can map them to 400 instead of 500.
				mismatchErr := busDomain.MenuItem.Reorder(ctx, categoryID, []uuid.UUID{item1.ID})
				if !errors.Is(mismatchErr, menuitembus.ErrInvalidOrder) {
					return fmt.Sprintf("mismatch length: error %v does not match ErrInvalidOrder", mismatchErr)
				}

				unknownErr := busDomain.MenuItem.Reorder(ctx, categoryID, []uuid.UUID{item1.ID, uuid.New()})
				if !errors.Is(unknownErr, menuitembus.ErrInvalidOrder) {
					return fmt.Sprintf("invalid id: error %v does not match ErrInvalidOrder", unknownErr)
				}

				return true
			},
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
		{
			Name:    "success",
			ExpResp: []int{10, 20},
			ExcFunc: func(ctx context.Context) any {
				// Reverse order: item2 first, then item1
				err := busDomain.MenuItem.Reorder(ctx, categoryID, []uuid.UUID{item2.ID, item1.ID})
				if err != nil {
					return err
				}

				itm2Updated, err := busDomain.MenuItem.QueryByID(ctx, item2.ID)
				if err != nil {
					return err
				}
				itm1Updated, err := busDomain.MenuItem.QueryByID(ctx, item1.ID)
				if err != nil {
					return err
				}

				if itm2Updated.Rank == nil || itm1Updated.Rank == nil {
					return "ranks are nil"
				}

				return []int{*itm2Updated.Rank, *itm1Updated.Rank}
			},
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
		{
			// Regression: the validation fetch must not be page-capped, or
			// categories with more than 100 items could never be reordered.
			Name:    "more_than_100_items",
			ExpResp: true,
			ExcFunc: func(ctx context.Context) any {
				cats, err := categorybus.TestSeedCategories(ctx, 1, sd.Restaurants[0].ID, busDomain.Category)
				if err != nil {
					return err
				}

				items, err := menuitembus.TestSeedMenuItems(ctx, 101, cats[0].ID, sd.Restaurants[0].ID, busDomain.MenuItem)
				if err != nil {
					return err
				}

				// Reverse the seeded order.
				orderedIDs := make([]uuid.UUID, len(items))
				for i, itm := range items {
					orderedIDs[len(items)-1-i] = itm.ID
				}

				if err := busDomain.MenuItem.Reorder(ctx, cats[0].ID, orderedIDs); err != nil {
					return err
				}

				got, err := busDomain.MenuItem.QueryByCategoryID(ctx, cats[0].ID)
				if err != nil {
					return err
				}

				if len(got) != len(orderedIDs) {
					return fmt.Sprintf("expected %d items, got %d", len(orderedIDs), len(got))
				}

				for i, itm := range got {
					if itm.ID != orderedIDs[i] {
						return fmt.Sprintf("position %d: expected item %s, got %s", i, orderedIDs[i], itm.ID)
					}
				}

				return true
			},
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
	}

	return table
}
