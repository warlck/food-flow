package modifieroptionbus_test

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
	"github.com/warlck/food-flow/business/domain/modifieroptionbus"
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
	RestaurantID    uuid.UUID
	MenuItemID      uuid.UUID
	ModifierGroupID uuid.UUID
	Options         []modifieroptionbus.ModifierOption
}

func Test_ModifierOption(t *testing.T) {
	t.Parallel()

	db := dbtest.New(t, "Test_ModifierOption")

	sd, err := insertSeedData(db.BusDomain)
	if err != nil {
		t.Fatalf("Seeding error: %s", err)
	}

	unittest.Run(t, query(db.BusDomain, sd), "query")
	unittest.Run(t, create(db.BusDomain, sd), "create")
	unittest.Run(t, update(db.BusDomain, sd), "update")
	unittest.Run(t, reorder(db.BusDomain, sd), "reorder")
	unittest.Run(t, lastAvailableOption(db.BusDomain, sd), "last-available-option")
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
	groups, err := modifiergroupbus.TestSeedModifierGroups(ctx, 1, items[0].ID, rests[0].ID, busDomain.ModifierGroup)
	if err != nil {
		return seedData{}, fmt.Errorf("seeding modifier groups: %w", err)
	}
	options, err := modifieroptionbus.TestSeedModifierOptions(ctx, 2, groups[0].ID, rests[0].ID, busDomain.ModifierOption)
	if err != nil {
		return seedData{}, fmt.Errorf("seeding modifier options: %w", err)
	}

	return seedData{
		RestaurantID:    rests[0].ID,
		MenuItemID:      items[0].ID,
		ModifierGroupID: groups[0].ID,
		Options:         options,
	}, nil
}

func query(busDomain dbtest.BusDomain, sd seedData) []unittest.Table {
	options := make([]modifieroptionbus.ModifierOption, len(sd.Options))
	copy(options, sd.Options)

	sort.Slice(options, func(i, j int) bool {
		r1 := options[i].Rank
		r2 := options[j].Rank
		if r1 != nil && r2 != nil && *r1 != *r2 {
			return *r1 < *r2
		}
		if r1 != nil && r2 == nil {
			return true
		}
		if r1 == nil && r2 != nil {
			return false
		}
		if options[i].Name.String() != options[j].Name.String() {
			return options[i].Name.String() < options[j].Name.String()
		}
		return options[i].ID.String() < options[j].ID.String()
	})

	table := []unittest.Table{
		{
			Name:    "by-group",
			ExpResp: options,
			ExcFunc: func(ctx context.Context) any {
				filter := modifieroptionbus.QueryFilter{
					ModifierGroupID: &sd.ModifierGroupID,
				}
				resp, err := busDomain.ModifierOption.Query(ctx, filter, modifieroptionbus.DefaultOrderBy, page.MustParse("1", "10"))
				if err != nil {
					return err
				}
				return resp
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, ok := got.([]modifieroptionbus.ModifierOption)
				if !ok {
					return "expected []modifieroptionbus.ModifierOption"
				}
				expResp := exp.([]modifieroptionbus.ModifierOption)
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
			ExpResp: sd.Options[0],
			ExcFunc: func(ctx context.Context) any {
				resp, err := busDomain.ModifierOption.QueryByID(ctx, sd.Options[0].ID)
				if err != nil {
					return err
				}
				return resp
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, ok := got.(modifieroptionbus.ModifierOption)
				if !ok {
					return "expected modifieroptionbus.ModifierOption"
				}
				expResp := exp.(modifieroptionbus.ModifierOption)
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
	avail := true
	table := []unittest.Table{
		{
			Name: "basic",
			ExpResp: modifieroptionbus.ModifierOption{
				ModifierGroupID: sd.ModifierGroupID,
				RestaurantID:    sd.RestaurantID,
				Name:            name.MustParse("Garlic Yogurt"),
				Description:     "Creamy garlic yogurt sauce",
				PriceDelta:      money.MustParse(1.50),
				Available:       true,
				Rank:            &rankVal,
			},
			ExcFunc: func(ctx context.Context) any {
				no := modifieroptionbus.NewModifierOption{
					ModifierGroupID: sd.ModifierGroupID,
					RestaurantID:    sd.RestaurantID,
					Name:            name.MustParse("Garlic Yogurt"),
					Description:     "Creamy garlic yogurt sauce",
					PriceDelta:      money.MustParse(1.50),
					Available:       &avail,
					Rank:            &rankVal,
				}
				resp, err := busDomain.ModifierOption.Create(ctx, no)
				if err != nil {
					return err
				}
				return resp
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, ok := got.(modifieroptionbus.ModifierOption)
				if !ok {
					return "expected modifieroptionbus.ModifierOption"
				}
				expResp := exp.(modifieroptionbus.ModifierOption)
				expResp.ID = gotResp.ID
				expResp.DateCreated = gotResp.DateCreated
				expResp.DateUpdated = gotResp.DateUpdated
				return cmp.Diff(gotResp, expResp)
			},
		},
		{
			Name: "zero-price-delta",
			ExpResp: modifieroptionbus.ModifierOption{
				ModifierGroupID: sd.ModifierGroupID,
				RestaurantID:    sd.RestaurantID,
				Name:            name.MustParse("No Sauce"),
				Description:     "Plain without sauce",
				PriceDelta:      money.MustParse(0.00),
				Available:       true,
				Rank:            nil,
			},
			ExcFunc: func(ctx context.Context) any {
				no := modifieroptionbus.NewModifierOption{
					ModifierGroupID: sd.ModifierGroupID,
					RestaurantID:    sd.RestaurantID,
					Name:            name.MustParse("No Sauce"),
					Description:     "Plain without sauce",
					PriceDelta:      money.MustParse(0.00),
					Available:       &avail,
					Rank:            nil,
				}
				resp, err := busDomain.ModifierOption.Create(ctx, no)
				if err != nil {
					return err
				}
				return resp
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, ok := got.(modifieroptionbus.ModifierOption)
				if !ok {
					return "expected modifieroptionbus.ModifierOption"
				}
				expResp := exp.(modifieroptionbus.ModifierOption)
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
	newRank := 50
	newPrice := money.MustParse(2.00)
	table := []unittest.Table{
		{
			Name: "set-price-and-rank",
			ExpResp: modifieroptionbus.ModifierOption{
				ID:              sd.Options[0].ID,
				ModifierGroupID: sd.Options[0].ModifierGroupID,
				RestaurantID:    sd.Options[0].RestaurantID,
				Name:            sd.Options[0].Name,
				Description:     sd.Options[0].Description,
				PriceDelta:      newPrice,
				Available:       sd.Options[0].Available,
				Rank:            &newRank,
				DateCreated:     sd.Options[0].DateCreated,
			},
			ExcFunc: func(ctx context.Context) any {
				uo := modifieroptionbus.UpdateModifierOption{
					PriceDelta: &newPrice,
					Rank:       opt.NewNullInt(50),
				}
				current, err := busDomain.ModifierOption.QueryByID(ctx, sd.Options[0].ID)
				if err != nil {
					return err
				}
				resp, err := busDomain.ModifierOption.Update(ctx, current, uo)
				if err != nil {
					return err
				}
				return resp
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, ok := got.(modifieroptionbus.ModifierOption)
				if !ok {
					return "expected modifieroptionbus.ModifierOption"
				}
				expResp := exp.(modifieroptionbus.ModifierOption)
				expResp.DateCreated = gotResp.DateCreated
				expResp.DateUpdated = gotResp.DateUpdated
				return cmp.Diff(gotResp, expResp)
			},
		},
		{
			Name: "clear-rank",
			ExpResp: modifieroptionbus.ModifierOption{
				ID:              sd.Options[0].ID,
				ModifierGroupID: sd.Options[0].ModifierGroupID,
				RestaurantID:    sd.Options[0].RestaurantID,
				Name:            sd.Options[0].Name,
				Description:     sd.Options[0].Description,
				PriceDelta:      newPrice,
				Available:       sd.Options[0].Available,
				Rank:            nil,
				DateCreated:     sd.Options[0].DateCreated,
			},
			ExcFunc: func(ctx context.Context) any {
				uo := modifieroptionbus.UpdateModifierOption{
					Rank: opt.NewNullIntNull(),
				}
				current, err := busDomain.ModifierOption.QueryByID(ctx, sd.Options[0].ID)
				if err != nil {
					return err
				}
				resp, err := busDomain.ModifierOption.Update(ctx, current, uo)
				if err != nil {
					return err
				}
				return resp
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, ok := got.(modifieroptionbus.ModifierOption)
				if !ok {
					return "expected modifieroptionbus.ModifierOption"
				}
				expResp := exp.(modifieroptionbus.ModifierOption)
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
			ExpResp: []int{10, 20, 30, 40},
			ExcFunc: func(ctx context.Context) any {
				curr, err := busDomain.ModifierOption.QueryAll(ctx, modifieroptionbus.QueryFilter{ModifierGroupID: &sd.ModifierGroupID}, modifieroptionbus.DefaultOrderBy)
				if err != nil {
					return err
				}
				if len(curr) != 4 {
					return fmt.Errorf("expected 4 options, got %d", len(curr))
				}
				// Reverse all 4
				reordered, err := busDomain.ModifierOption.Reorder(ctx, sd.ModifierGroupID, []uuid.UUID{curr[3].ID, curr[2].ID, curr[1].ID, curr[0].ID})
				if err != nil {
					return err
				}
				ranks := make([]int, len(reordered))
				for i, o := range reordered {
					if o.Rank != nil {
						ranks[i] = *o.Rank
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
			ExpResp: modifieroptionbus.ErrInvalidReorder,
			ExcFunc: func(ctx context.Context) any {
				_, err := busDomain.ModifierOption.Reorder(ctx, sd.ModifierGroupID, []uuid.UUID{sd.Options[0].ID, sd.Options[0].ID})
				return err
			},
			CmpFunc: func(got any, exp any) string {
				gotErr, ok := got.(error)
				if !ok {
					return "expected error response"
				}
				if !errors.Is(gotErr, modifieroptionbus.ErrInvalidReorder) {
					return fmt.Sprintf("expected ErrInvalidReorder, got %v", gotErr)
				}
				return ""
			},
		},
	}

	return table
}

// lastAvailableOption covers the invariant that the last available option of
// an active required group cannot be disabled or deleted. Each case builds
// its own group/options so the cases stay independent. Runs after reorder
// because it creates extra groups and options on the seeded menu item.
func lastAvailableOption(busDomain dbtest.BusDomain, sd seedData) []unittest.Table {
	disable := false

	newGroup := func(ctx context.Context, minSelections int, available bool) (modifiergroupbus.ModifierGroup, error) {
		return busDomain.ModifierGroup.Create(ctx, modifiergroupbus.NewModifierGroup{
			MenuItemID:    sd.MenuItemID,
			RestaurantID:  sd.RestaurantID,
			Name:          name.MustParse(fmt.Sprintf("Guard Group %s", uuid.NewString()[:8])),
			MinSelections: minSelections,
			MaxSelections: 1,
			Available:     available,
		})
	}

	newOption := func(ctx context.Context, groupID uuid.UUID, available bool) (modifieroptionbus.ModifierOption, error) {
		return busDomain.ModifierOption.Create(ctx, modifieroptionbus.NewModifierOption{
			ModifierGroupID: groupID,
			RestaurantID:    sd.RestaurantID,
			Name:            name.MustParse(fmt.Sprintf("Guard Option %s", uuid.NewString()[:8])),
			PriceDelta:      money.MustParse(0.50),
			Available:       &available,
		})
	}

	expectLastOptionErr := func(got any, exp any) string {
		gotErr, ok := got.(error)
		if !ok {
			return "expected error response"
		}
		if !errors.Is(gotErr, modifieroptionbus.ErrLastAvailableOption) {
			return fmt.Sprintf("expected ErrLastAvailableOption, got %v", gotErr)
		}
		return ""
	}

	table := []unittest.Table{
		{
			Name:    "disable-last-available-of-active-required-group",
			ExpResp: modifieroptionbus.ErrLastAvailableOption,
			ExcFunc: func(ctx context.Context) any {
				grp, err := newGroup(ctx, 1, true)
				if err != nil {
					return err
				}
				opt, err := newOption(ctx, grp.ID, true)
				if err != nil {
					return err
				}
				_, err = busDomain.ModifierOption.Update(ctx, opt, modifieroptionbus.UpdateModifierOption{
					Available: &disable,
				})
				return err
			},
			CmpFunc: expectLastOptionErr,
		},
		{
			Name:    "delete-last-available-of-active-required-group",
			ExpResp: modifieroptionbus.ErrLastAvailableOption,
			ExcFunc: func(ctx context.Context) any {
				grp, err := newGroup(ctx, 1, true)
				if err != nil {
					return err
				}
				opt, err := newOption(ctx, grp.ID, true)
				if err != nil {
					return err
				}
				return busDomain.ModifierOption.Delete(ctx, opt)
			},
			CmpFunc: expectLastOptionErr,
		},
		{
			Name:    "disable-one-of-two-available",
			ExpResp: false,
			ExcFunc: func(ctx context.Context) any {
				grp, err := newGroup(ctx, 1, true)
				if err != nil {
					return err
				}
				opt, err := newOption(ctx, grp.ID, true)
				if err != nil {
					return err
				}
				if _, err := newOption(ctx, grp.ID, true); err != nil {
					return err
				}
				resp, err := busDomain.ModifierOption.Update(ctx, opt, modifieroptionbus.UpdateModifierOption{
					Available: &disable,
				})
				if err != nil {
					return err
				}
				return resp.Available
			},
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
		{
			Name:    "delete-one-of-two-available",
			ExpResp: nil,
			ExcFunc: func(ctx context.Context) any {
				grp, err := newGroup(ctx, 1, true)
				if err != nil {
					return err
				}
				opt, err := newOption(ctx, grp.ID, true)
				if err != nil {
					return err
				}
				if _, err := newOption(ctx, grp.ID, true); err != nil {
					return err
				}
				return busDomain.ModifierOption.Delete(ctx, opt)
			},
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
		{
			Name:    "disable-only-option-of-optional-group",
			ExpResp: false,
			ExcFunc: func(ctx context.Context) any {
				grp, err := newGroup(ctx, 0, true)
				if err != nil {
					return err
				}
				opt, err := newOption(ctx, grp.ID, true)
				if err != nil {
					return err
				}
				resp, err := busDomain.ModifierOption.Update(ctx, opt, modifieroptionbus.UpdateModifierOption{
					Available: &disable,
				})
				if err != nil {
					return err
				}
				return resp.Available
			},
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
		{
			Name:    "disable-only-option-of-inactive-required-group",
			ExpResp: false,
			ExcFunc: func(ctx context.Context) any {
				grp, err := newGroup(ctx, 1, false)
				if err != nil {
					return err
				}
				opt, err := newOption(ctx, grp.ID, true)
				if err != nil {
					return err
				}
				resp, err := busDomain.ModifierOption.Update(ctx, opt, modifieroptionbus.UpdateModifierOption{
					Available: &disable,
				})
				if err != nil {
					return err
				}
				return resp.Available
			},
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
		{
			Name:    "delete-unavailable-option-of-active-required-group",
			ExpResp: nil,
			ExcFunc: func(ctx context.Context) any {
				grp, err := newGroup(ctx, 1, true)
				if err != nil {
					return err
				}
				if _, err := newOption(ctx, grp.ID, true); err != nil {
					return err
				}
				unavailable, err := newOption(ctx, grp.ID, false)
				if err != nil {
					return err
				}
				return busDomain.ModifierOption.Delete(ctx, unavailable)
			},
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
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
				current, err := busDomain.ModifierOption.QueryByID(ctx, sd.Options[1].ID)
				if err != nil {
					return err
				}
				return busDomain.ModifierOption.Delete(ctx, current)
			},
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
	}

	return table
}
