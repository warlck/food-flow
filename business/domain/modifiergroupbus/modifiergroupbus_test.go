package modifiergroupbus_test

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/warlck/food-flow/business/domain/categorybus"
	"github.com/warlck/food-flow/business/domain/menuitembus"
	"github.com/warlck/food-flow/business/domain/modifiergroupbus"
	"github.com/warlck/food-flow/business/domain/organizationbus"
	"github.com/warlck/food-flow/business/domain/restaurantbus"
	"github.com/warlck/food-flow/business/sdk/dbtest"
	"github.com/warlck/food-flow/business/sdk/page"
	"github.com/warlck/food-flow/business/sdk/unittest"
	"github.com/warlck/food-flow/business/types/name"
	"github.com/warlck/food-flow/business/types/opt"
)

type seedData struct {
	RestaurantID uuid.UUID
	MenuItemID   uuid.UUID
	Groups       []modifiergroupbus.ModifierGroup
}

func Test_ModifierGroup(t *testing.T) {
	t.Parallel()

	db := dbtest.New(t, "Test_ModifierGroup")

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
	items, err := menuitembus.TestSeedMenuItems(ctx, 2, cats[0].ID, rests[0].ID, busDomain.MenuItem)
	if err != nil {
		return seedData{}, fmt.Errorf("seeding menu items: %w", err)
	}
	groups, err := modifiergroupbus.TestSeedModifierGroups(ctx, 2, items[0].ID, rests[0].ID, busDomain.ModifierGroup)
	if err != nil {
		return seedData{}, fmt.Errorf("seeding modifier groups: %w", err)
	}

	return seedData{
		RestaurantID: rests[0].ID,
		MenuItemID:   items[0].ID,
		Groups:       groups,
	}, nil
}

func query(busDomain dbtest.BusDomain, sd seedData) []unittest.Table {
	groups := make([]modifiergroupbus.ModifierGroup, len(sd.Groups))
	copy(groups, sd.Groups)

	sort.Slice(groups, func(i, j int) bool {
		r1 := groups[i].Rank
		r2 := groups[j].Rank
		if r1 != nil && r2 != nil && *r1 != *r2 {
			return *r1 < *r2
		}
		if r1 != nil && r2 == nil {
			return true
		}
		if r1 == nil && r2 != nil {
			return false
		}
		if groups[i].Name.String() != groups[j].Name.String() {
			return groups[i].Name.String() < groups[j].Name.String()
		}
		return groups[i].ID.String() < groups[j].ID.String()
	})

	table := []unittest.Table{
		{
			Name:    "by-menu-item",
			ExpResp: groups,
			ExcFunc: func(ctx context.Context) any {
				filter := modifiergroupbus.QueryFilter{
					MenuItemID: &sd.MenuItemID,
				}
				resp, err := busDomain.ModifierGroup.Query(ctx, filter, modifiergroupbus.DefaultOrderBy, page.MustParse("1", "10"))
				if err != nil {
					return err
				}
				return resp
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, ok := got.([]modifiergroupbus.ModifierGroup)
				if !ok {
					return "expected []modifiergroupbus.ModifierGroup"
				}
				expResp := exp.([]modifiergroupbus.ModifierGroup)
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
			ExpResp: sd.Groups[0],
			ExcFunc: func(ctx context.Context) any {
				resp, err := busDomain.ModifierGroup.QueryByID(ctx, sd.Groups[0].ID)
				if err != nil {
					return err
				}
				return resp
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, ok := got.(modifiergroupbus.ModifierGroup)
				if !ok {
					return "expected modifiergroupbus.ModifierGroup"
				}
				expResp := exp.(modifiergroupbus.ModifierGroup)
				expResp.DateCreated = gotResp.DateCreated
				expResp.DateUpdated = gotResp.DateUpdated
				return cmp.Diff(gotResp, expResp)
			},
		},
	}

	return table
}

func create(busDomain dbtest.BusDomain, sd seedData) []unittest.Table {
	rankVal := 30
	table := []unittest.Table{
		{
			Name: "basic",
			ExpResp: modifiergroupbus.ModifierGroup{
				MenuItemID:    sd.MenuItemID,
				RestaurantID:  sd.RestaurantID,
				Name:          name.MustParse("Choice of Sauce"),
				Description:   "Select your primary sauce",
				MinSelections: 1,
				MaxSelections: 1,
				Available:     false,
				Rank:          &rankVal,
			},
			ExcFunc: func(ctx context.Context) any {
				ng := modifiergroupbus.NewModifierGroup{
					MenuItemID:    sd.MenuItemID,
					RestaurantID:  sd.RestaurantID,
					Name:          name.MustParse("Choice of Sauce"),
					Description:   "Select your primary sauce",
					MinSelections: 1,
					MaxSelections: 1,
					Available:     false,
					Rank:          &rankVal,
				}
				resp, err := busDomain.ModifierGroup.Create(ctx, ng)
				if err != nil {
					return err
				}
				return resp
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, ok := got.(modifiergroupbus.ModifierGroup)
				if !ok {
					return "expected modifiergroupbus.ModifierGroup"
				}
				expResp := exp.(modifiergroupbus.ModifierGroup)
				expResp.ID = gotResp.ID
				expResp.DateCreated = gotResp.DateCreated
				expResp.DateUpdated = gotResp.DateUpdated
				return cmp.Diff(gotResp, expResp)
			},
		},
		{
			Name:    "invalid-min-selections",
			ExpResp: "min_selections must be 0 or 1",
			ExcFunc: func(ctx context.Context) any {
				ng := modifiergroupbus.NewModifierGroup{
					MenuItemID:    sd.MenuItemID,
					RestaurantID:  sd.RestaurantID,
					Name:          name.MustParse("Invalid Group"),
					MinSelections: 2,
					MaxSelections: 1,
				}
				_, err := busDomain.ModifierGroup.Create(ctx, ng)
				if err != nil {
					return err.Error()
				}
				return "expected error"
			},
			CmpFunc: func(got any, exp any) string {
				if got != exp {
					return fmt.Sprintf("expected %v, got %v", exp, got)
				}
				return ""
			},
		},
	}

	return table
}

func update(busDomain dbtest.BusDomain, sd seedData) []unittest.Table {
	newRank := 50
	table := []unittest.Table{
		{
			Name: "set-rank-and-name",
			ExpResp: modifiergroupbus.ModifierGroup{
				ID:            sd.Groups[0].ID,
				MenuItemID:    sd.Groups[0].MenuItemID,
				RestaurantID:  sd.Groups[0].RestaurantID,
				Name:          name.MustParse("Updated Group Name"),
				Description:   sd.Groups[0].Description,
				MinSelections: sd.Groups[0].MinSelections,
				MaxSelections: sd.Groups[0].MaxSelections,
				Available:     sd.Groups[0].Available,
				Rank:          &newRank,
				DateCreated:   sd.Groups[0].DateCreated,
			},
			ExcFunc: func(ctx context.Context) any {
				ug := modifiergroupbus.UpdateModifierGroup{
					Name: dbtest.NamePointer("Updated Group Name"),
					Rank: opt.NewNullInt(50),
				}
				current, err := busDomain.ModifierGroup.QueryByID(ctx, sd.Groups[0].ID)
				if err != nil {
					return err
				}
				resp, err := busDomain.ModifierGroup.Update(ctx, current, ug)
				if err != nil {
					return err
				}
				return resp
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, ok := got.(modifiergroupbus.ModifierGroup)
				if !ok {
					return "expected modifiergroupbus.ModifierGroup"
				}
				expResp := exp.(modifiergroupbus.ModifierGroup)
				expResp.DateCreated = gotResp.DateCreated
				expResp.DateUpdated = gotResp.DateUpdated
				return cmp.Diff(gotResp, expResp)
			},
		},
		{
			Name: "clear-rank",
			ExpResp: modifiergroupbus.ModifierGroup{
				ID:            sd.Groups[0].ID,
				MenuItemID:    sd.Groups[0].MenuItemID,
				RestaurantID:  sd.Groups[0].RestaurantID,
				Name:          name.MustParse("Updated Group Name"),
				Description:   sd.Groups[0].Description,
				MinSelections: sd.Groups[0].MinSelections,
				MaxSelections: sd.Groups[0].MaxSelections,
				Available:     sd.Groups[0].Available,
				Rank:          nil,
				DateCreated:   sd.Groups[0].DateCreated,
			},
			ExcFunc: func(ctx context.Context) any {
				ug := modifiergroupbus.UpdateModifierGroup{
					Rank: opt.NewNullIntNull(),
				}
				current, err := busDomain.ModifierGroup.QueryByID(ctx, sd.Groups[0].ID)
				if err != nil {
					return err
				}
				resp, err := busDomain.ModifierGroup.Update(ctx, current, ug)
				if err != nil {
					return err
				}
				return resp
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, ok := got.(modifiergroupbus.ModifierGroup)
				if !ok {
					return "expected modifiergroupbus.ModifierGroup"
				}
				expResp := exp.(modifiergroupbus.ModifierGroup)
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
			Name:    "success-exact-set",
			ExpResp: []int{10, 20, 30},
			ExcFunc: func(ctx context.Context) any {
				curr, err := busDomain.ModifierGroup.QueryAll(ctx, modifiergroupbus.QueryFilter{MenuItemID: &sd.MenuItemID}, modifiergroupbus.DefaultOrderBy)
				if err != nil {
					return err
				}
				if len(curr) != 3 {
					return fmt.Errorf("expected 3 groups, got %d", len(curr))
				}
				// Reorder all 3: reverse
				reordered, err := busDomain.ModifierGroup.Reorder(ctx, sd.MenuItemID, []uuid.UUID{curr[2].ID, curr[0].ID, curr[1].ID})
				if err != nil {
					return err
				}
				ranks := make([]int, len(reordered))
				for i, g := range reordered {
					if g.Rank != nil {
						ranks[i] = *g.Rank
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
			Name:    "duplicate-id-error",
			ExpResp: modifiergroupbus.ErrInvalidReorder,
			ExcFunc: func(ctx context.Context) any {
				_, err := busDomain.ModifierGroup.Reorder(ctx, sd.MenuItemID, []uuid.UUID{sd.Groups[0].ID, sd.Groups[0].ID})
				return err
			},
			CmpFunc: func(got any, exp any) string {
				gotErr, ok := got.(error)
				if !ok {
					return "expected error response"
				}
				if !errors.Is(gotErr, modifiergroupbus.ErrInvalidReorder) {
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
				current, err := busDomain.ModifierGroup.QueryByID(ctx, sd.Groups[1].ID)
				if err != nil {
					return err
				}
				return busDomain.ModifierGroup.Delete(ctx, current)
			},
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
	}

	return table
}
