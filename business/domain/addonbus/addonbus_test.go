package addonbus_test

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/warlck/food-flow/business/domain/addonbus"
	"github.com/warlck/food-flow/business/domain/categorybus"
	"github.com/warlck/food-flow/business/domain/menuitembus"
	"github.com/warlck/food-flow/business/domain/organizationbus"
	"github.com/warlck/food-flow/business/domain/restaurantbus"
	"github.com/warlck/food-flow/business/sdk/dbtest"
	"github.com/warlck/food-flow/business/sdk/page"
	"github.com/warlck/food-flow/business/sdk/unittest"
	"github.com/warlck/food-flow/business/types/money"
	"github.com/warlck/food-flow/business/types/name"
	"github.com/warlck/food-flow/business/types/opt"
)

type seedData struct {
	RestaurantID uuid.UUID
	MenuItemID   uuid.UUID
	Addons       []addonbus.Addon
}

func Test_Addon(t *testing.T) {
	t.Parallel()

	db := dbtest.New(t, "Test_Addon")

	sd, err := insertSeedData(db.BusDomain)
	if err != nil {
		t.Fatalf("Seeding error: %s", err)
	}

	unittest.Run(t, query(db.BusDomain, sd), "query")
	unittest.Run(t, create(db.BusDomain, sd), "create")
	unittest.Run(t, update(db.BusDomain, sd), "update")
	unittest.Run(t, reorder(db.BusDomain, sd), "reorder")
	unittest.Run(t, delete(db.BusDomain, sd), "delete")
}

func insertSeedData(busDomain dbtest.BusDomain) (seedData, error) {
	ctx := context.Background()

	orgs, err := organizationbus.TestSeedOrganizations(ctx, 1, busDomain.Organization)
	if err != nil {
		return seedData{}, fmt.Errorf("seeding organizations: %w", err)
	}
	rests, err := restaurantbus.TestSeedRestaurants(ctx, 1, busDomain.Restaurant, orgs[0].ID)
	if err != nil {
		return seedData{}, fmt.Errorf("seeding restaurants: %w", err)
	}
	cats, err := categorybus.TestSeedCategories(ctx, 1, rests[0].ID, busDomain.Category)
	if err != nil {
		return seedData{}, fmt.Errorf("seeding categories: %w", err)
	}
	items, err := menuitembus.TestSeedMenuItems(ctx, 1, cats[0].ID, rests[0].ID, busDomain.MenuItem)
	if err != nil {
		return seedData{}, fmt.Errorf("seeding menu items: %w", err)
	}
	addons, err := addonbus.TestSeedAddons(ctx, 4, items[0].ID, rests[0].ID, busDomain.Addon)
	if err != nil {
		return seedData{}, fmt.Errorf("seeding addons: %w", err)
	}

	return seedData{
		RestaurantID: rests[0].ID,
		MenuItemID:   items[0].ID,
		Addons:       addons,
	}, nil
}

func query(busDomain dbtest.BusDomain, sd seedData) []unittest.Table {
	addons := make([]addonbus.Addon, len(sd.Addons))
	copy(addons, sd.Addons)

	sort.Slice(addons, func(i, j int) bool {
		return addons[i].ID.String() <= addons[j].ID.String()
	})

	table := []unittest.Table{
		{
			Name:    "all",
			ExpResp: addons,
			ExcFunc: func(ctx context.Context) any {
				filter := addonbus.QueryFilter{
					MenuItemID:   &sd.MenuItemID,
					RestaurantID: &sd.RestaurantID,
				}
				resp, err := busDomain.Addon.Query(ctx, filter, addonbus.DefaultOrderBy, page.MustParse("1", "10"))
				if err != nil {
					return err
				}
				return resp
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, ok := got.([]addonbus.Addon)
				if !ok {
					return "expected []addonbus.Addon"
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
			Name:    "by-id",
			ExpResp: sd.Addons[0],
			ExcFunc: func(ctx context.Context) any {
				resp, err := busDomain.Addon.QueryByID(ctx, sd.Addons[0].ID)
				if err != nil {
					return err
				}
				return resp
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, ok := got.(addonbus.Addon)
				if !ok {
					return "expected addonbus.Addon"
				}
				expResp := exp.(addonbus.Addon)
				expResp.DateCreated = gotResp.DateCreated
				expResp.DateUpdated = gotResp.DateUpdated
				return cmp.Diff(gotResp, expResp)
			},
		},
	}

	return table
}

func create(busDomain dbtest.BusDomain, sd seedData) []unittest.Table {
	avail := true
	r10 := 10
	table := []unittest.Table{
		{
			Name: "basic",
			ExpResp: addonbus.Addon{
				MenuItemID:   sd.MenuItemID,
				RestaurantID: sd.RestaurantID,
				Name:         name.MustParse("French Fries"),
				Description:  "Crispy golden fries",
				Price:        money.MustParse(4.50),
				Available:    true,
				MaxQuantity:  5,
				Rank:         &r10,
			},
			ExcFunc: func(ctx context.Context) any {
				na := addonbus.NewAddon{
					MenuItemID:   sd.MenuItemID,
					RestaurantID: sd.RestaurantID,
					Name:         name.MustParse("French Fries"),
					Description:  "Crispy golden fries",
					Price:        money.MustParse(4.50),
					Available:    &avail,
					MaxQuantity:  5,
					Rank:         &r10,
				}
				resp, err := busDomain.Addon.Create(ctx, na)
				if err != nil {
					return err
				}
				return resp
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, ok := got.(addonbus.Addon)
				if !ok {
					return "expected addonbus.Addon"
				}
				expResp := exp.(addonbus.Addon)
				expResp.ID = gotResp.ID
				expResp.DateCreated = gotResp.DateCreated
				expResp.DateUpdated = gotResp.DateUpdated
				return cmp.Diff(gotResp, expResp)
			},
		},
		{
			Name:    "duplicate-name",
			ExpResp: addonbus.ErrDuplicateName,
			ExcFunc: func(ctx context.Context) any {
				na := addonbus.NewAddon{
					MenuItemID:   sd.MenuItemID,
					RestaurantID: sd.RestaurantID,
					Name:         name.MustParse("French Fries"),
					Description:  "Another fries",
					Price:        money.MustParse(5.00),
					Available:    &avail,
					MaxQuantity:  5,
				}
				_, err := busDomain.Addon.Create(ctx, na)
				return err
			},
			CmpFunc: func(got any, exp any) string {
				gotErr, ok := got.(error)
				if !ok {
					return "expected error response"
				}
				if !errors.Is(gotErr, addonbus.ErrDuplicateName) {
					return fmt.Sprintf("expected ErrDuplicateName, got %v", gotErr)
				}
				return ""
			},
		},
	}

	return table
}

func update(busDomain dbtest.BusDomain, sd seedData) []unittest.Table {
	newPrice := money.MustParse(6.00)
	newMax := 8
	table := []unittest.Table{
		{
			Name: "basic",
			ExpResp: addonbus.Addon{
				ID:           sd.Addons[0].ID,
				MenuItemID:   sd.Addons[0].MenuItemID,
				RestaurantID: sd.Addons[0].RestaurantID,
				Name:         name.MustParse("Updated Addon Name"),
				Description:  "Updated description",
				Price:        newPrice,
				Available:    sd.Addons[0].Available,
				MaxQuantity:  newMax,
				Rank:         nil,
				DateCreated:  sd.Addons[0].DateCreated,
			},
			ExcFunc: func(ctx context.Context) any {
				ua := addonbus.UpdateAddon{
					Name:        dbtest.NamePointer("Updated Addon Name"),
					Description: dbtest.StringPointer("Updated description"),
					Price:       &newPrice,
					MaxQuantity: &newMax,
					Rank:        opt.NewNullIntNull(), // test unsetting rank
				}
				current, err := busDomain.Addon.QueryByID(ctx, sd.Addons[0].ID)
				if err != nil {
					return err
				}
				resp, err := busDomain.Addon.Update(ctx, current, ua)
				if err != nil {
					return err
				}
				return resp
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, ok := got.(addonbus.Addon)
				if !ok {
					return "expected addonbus.Addon"
				}
				expResp := exp.(addonbus.Addon)
				expResp.DateCreated = gotResp.DateCreated
				expResp.DateUpdated = gotResp.DateUpdated
				return cmp.Diff(gotResp, expResp)
			},
		},
	}

	return table
}

func reorder(busDomain dbtest.BusDomain, sd seedData) []unittest.Table {
	table := []unittest.Table{
		{
			Name:    "reorder-success",
			ExpResp: 5,
			ExcFunc: func(ctx context.Context) any {
				current, err := busDomain.Addon.QueryAll(ctx, addonbus.QueryFilter{
					MenuItemID: &sd.MenuItemID,
				}, addonbus.DefaultOrderBy)
				if err != nil {
					return err
				}

				orderIDs := make([]uuid.UUID, len(current))
				for i := range current {
					orderIDs[i] = current[len(current)-1-i].ID
				}

				reordered, err := busDomain.Addon.Reorder(ctx, sd.MenuItemID, orderIDs)
				if err != nil {
					return err
				}
				return len(reordered)
			},
			CmpFunc: func(got any, exp any) string {
				gotLen, ok := got.(int)
				if !ok {
					return fmt.Sprintf("expected int, got %T: %v", got, got)
				}
				return cmp.Diff(gotLen, exp.(int))
			},
		},
		{
			Name:    "reorder-mismatch-error",
			ExpResp: addonbus.ErrInvalidReorder,
			ExcFunc: func(ctx context.Context) any {
				// Pass fewer addons than exist
				_, err := busDomain.Addon.Reorder(ctx, sd.MenuItemID, []uuid.UUID{sd.Addons[0].ID})
				return err
			},
			CmpFunc: func(got any, exp any) string {
				gotErr, ok := got.(error)
				if !ok {
					return "expected error response"
				}
				if !errors.Is(gotErr, addonbus.ErrInvalidReorder) {
					return fmt.Sprintf("expected ErrInvalidReorder, got %v", gotErr)
				}
				return ""
			},
		},
	}

	return table
}

func delete(busDomain dbtest.BusDomain, sd seedData) []unittest.Table {
	table := []unittest.Table{
		{
			Name:    "basic",
			ExpResp: nil,
			ExcFunc: func(ctx context.Context) any {
				current, err := busDomain.Addon.QueryByID(ctx, sd.Addons[3].ID)
				if err != nil {
					return err
				}
				return busDomain.Addon.Delete(ctx, current)
			},
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
	}

	return table
}
