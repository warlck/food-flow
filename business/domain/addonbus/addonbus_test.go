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
	unittest.Run(t, menuItemAddons(db.BusDomain, sd), "menu-item-addons")
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
	addons, err := addonbus.TestSeedAddons(ctx, 4, rests[0].ID, busDomain.Addon)
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
	table := []unittest.Table{
		{
			Name: "basic",
			ExpResp: addonbus.Addon{
				RestaurantID: sd.RestaurantID,
				Name:         name.MustParse("French Fries"),
				Description:  "Crispy golden fries",
				Price:        money.MustParse(4.50),
				Available:    true,
				MaxQuantity:  5,
			},
			ExcFunc: func(ctx context.Context) any {
				na := addonbus.NewAddon{
					RestaurantID: sd.RestaurantID,
					Name:         name.MustParse("French Fries"),
					Description:  "Crispy golden fries",
					Price:        money.MustParse(4.50),
					Available:    &avail,
					MaxQuantity:  5,
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
				RestaurantID: sd.Addons[0].RestaurantID,
				Name:         name.MustParse("Updated Addon Name"),
				Description:  "Updated description",
				Price:        newPrice,
				Available:    sd.Addons[0].Available,
				MaxQuantity:  newMax,
				DateCreated:  sd.Addons[0].DateCreated,
			},
			ExcFunc: func(ctx context.Context) any {
				ua := addonbus.UpdateAddon{
					Name:        dbtest.NamePointer("Updated Addon Name"),
					Description: dbtest.StringPointer("Updated description"),
					Price:       &newPrice,
					MaxQuantity: &newMax,
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

func menuItemAddons(busDomain dbtest.BusDomain, sd seedData) []unittest.Table {
	r10 := 10
	r20 := 20
	table := []unittest.Table{
		{
			Name:    "replace-assignments",
			ExpResp: []int{10, 20},
			ExcFunc: func(ctx context.Context) any {
				assignments := []addonbus.ItemAddonAssignment{
					{AddonID: sd.Addons[0].ID, Rank: &r10},
					{AddonID: sd.Addons[1].ID, Rank: &r20},
				}
				assigned, err := busDomain.Addon.ReplaceMenuItemAddons(ctx, sd.MenuItemID, sd.RestaurantID, assignments)
				if err != nil {
					return err
				}
				ranks := make([]int, len(assigned))
				for i, a := range assigned {
					if a.Rank != nil {
						ranks[i] = *a.Rank
					}
				}
				return ranks
			},
			CmpFunc: func(got any, exp any) string {
				gotRanks, ok := got.([]int)
				if !ok {
					return fmt.Sprintf("expected []int, got %T: %v", got, got)
				}
				return cmp.Diff(gotRanks, exp.([]int))
			},
		},
		{
			Name:    "reorder-assignments",
			ExpResp: []uuid.UUID{sd.Addons[1].ID, sd.Addons[0].ID},
			ExcFunc: func(ctx context.Context) any {
				reordered, err := busDomain.Addon.ReorderMenuItemAddons(ctx, sd.MenuItemID, []uuid.UUID{sd.Addons[1].ID, sd.Addons[0].ID})
				if err != nil {
					return err
				}
				ids := make([]uuid.UUID, len(reordered))
				for i, a := range reordered {
					ids[i] = a.Addon.ID
				}
				return ids
			},
			CmpFunc: func(got any, exp any) string {
				gotIDs, ok := got.([]uuid.UUID)
				if !ok {
					return fmt.Sprintf("expected []uuid.UUID, got %T: %v", got, got)
				}
				return cmp.Diff(gotIDs, exp.([]uuid.UUID))
			},
		},
		{
			Name:    "reorder-mismatch-error",
			ExpResp: addonbus.ErrInvalidReorder,
			ExcFunc: func(ctx context.Context) any {
				// Only pass 1 addon when 2 are assigned
				_, err := busDomain.Addon.ReorderMenuItemAddons(ctx, sd.MenuItemID, []uuid.UUID{sd.Addons[1].ID})
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
