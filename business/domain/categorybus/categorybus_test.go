package categorybus_test

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
	"github.com/warlck/food-flow/business/domain/organizationbus"
	"github.com/warlck/food-flow/business/domain/restaurantbus"
	"github.com/warlck/food-flow/business/sdk/dbtest"
	"github.com/warlck/food-flow/business/sdk/page"
	"github.com/warlck/food-flow/business/sdk/unittest"
	"github.com/warlck/food-flow/business/types/name"
	"github.com/warlck/food-flow/business/types/opt"
)

func Test_Category(t *testing.T) {
	t.Parallel()

	db := dbtest.New(t, "Test_Category")

	sd, err := insertSeedData(db.BusDomain)
	if err != nil {
		t.Fatalf("Seeding error: %s", err)
	}

	// -------------------------------------------------------------------------

	unittest.Run(t, query(db.BusDomain, sd), "query")
	unittest.Run(t, queryAll(db.BusDomain, sd), "query-all")
	unittest.Run(t, create(db.BusDomain, sd), "create")
	unittest.Run(t, update(db.BusDomain, sd), "update")
	unittest.Run(t, reorder(db.BusDomain, sd), "reorder")
	unittest.Run(t, delete(db.BusDomain, sd), "delete")
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

	// Seed categories for restaurant 1
	cats1, err := categorybus.TestSeedCategories(ctx, 2, rests[0].ID, busDomain.Category)
	if err != nil {
		return unittest.SeedData{}, fmt.Errorf("seeding categories : %w", err)
	}

	// Seed categories for restaurant 2
	cats2, err := categorybus.TestSeedCategories(ctx, 2, rests[1].ID, busDomain.Category)
	if err != nil {
		return unittest.SeedData{}, fmt.Errorf("seeding categories : %w", err)
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
	}

	return sd, nil
}

// =============================================================================

func query(busDomain dbtest.BusDomain, sd unittest.SeedData) []unittest.Table {
	cats := make([]categorybus.Category, 0, len(sd.Categories))

	for _, cat := range sd.Categories {
		cats = append(cats, cat.Category)
	}

	sort.Slice(cats, func(i, j int) bool {
		r1 := cats[i].Rank
		r2 := cats[j].Rank
		if r1 != nil && r2 != nil && *r1 != *r2 {
			return *r1 < *r2
		}
		if r1 != nil && r2 == nil {
			return true
		}
		if r1 == nil && r2 != nil {
			return false
		}
		if cats[i].Name.String() != cats[j].Name.String() {
			return cats[i].Name.String() < cats[j].Name.String()
		}
		return cats[i].ID.String() < cats[j].ID.String()
	})

	table := []unittest.Table{
		{
			Name:    "all",
			ExpResp: cats,
			ExcFunc: func(ctx context.Context) any {
				filter := categorybus.QueryFilter{
					Name: dbtest.NamePointer("Category"),
				}

				resp, err := busDomain.Category.Query(ctx, filter, categorybus.DefaultOrderBy, page.MustParse("1", "10"))
				if err != nil {
					return err
				}

				return resp
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.([]categorybus.Category)
				if !exists {
					return "error occurred"
				}

				expResp := exp.([]categorybus.Category)

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
			ExpResp: sd.Categories[0].Category,
			ExcFunc: func(ctx context.Context) any {
				resp, err := busDomain.Category.QueryByID(ctx, sd.Categories[0].ID)
				if err != nil {
					return err
				}

				return resp
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.(categorybus.Category)
				if !exists {
					return "error occurred"
				}

				expResp := exp.(categorybus.Category)

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

func queryAll(busDomain dbtest.BusDomain, sd unittest.SeedData) []unittest.Table {
	table := []unittest.Table{
		{
			Name:    "all-unfiltered",
			ExpResp: len(sd.Categories),
			ExcFunc: func(ctx context.Context) any {
				resp, err := busDomain.Category.QueryAll(ctx, categorybus.QueryFilter{}, categorybus.DefaultOrderBy)
				if err != nil {
					return err
				}
				return len(resp)
			},
			CmpFunc: func(got any, exp any) string {
				gotLen, ok := got.(int)
				if !ok {
					return fmt.Sprintf("expected int, got %T", got)
				}
				expLen := exp.(int)
				if gotLen != expLen {
					return fmt.Sprintf("expected %d categories, got %d", expLen, gotLen)
				}
				return ""
			},
		},
		{
			Name:    "by-restaurant",
			ExpResp: 2,
			ExcFunc: func(ctx context.Context) any {
				resp, err := busDomain.Category.QueryAll(ctx, categorybus.QueryFilter{
					RestaurantID: &sd.Restaurants[0].ID,
				}, categorybus.DefaultOrderBy)
				if err != nil {
					return err
				}
				return len(resp)
			},
			CmpFunc: func(got any, exp any) string {
				gotLen, ok := got.(int)
				if !ok {
					return fmt.Sprintf("expected int, got %T", got)
				}
				expLen := exp.(int)
				if gotLen != expLen {
					return fmt.Sprintf("expected %d categories, got %d", expLen, gotLen)
				}
				return ""
			},
		},
	}

	return table
}

func create(busDomain dbtest.BusDomain, sd unittest.SeedData) []unittest.Table {
	rankVal := 10
	table := []unittest.Table{
		{
			Name: "basic",
			ExpResp: categorybus.Category{
				Name:         name.MustParse("Appetizers"),
				Description:  "Delicious starters and appetizers",
				RestaurantID: sd.Restaurants[1].ID, // use restaurant 1 to keep restaurant 0 clean for reorder
				Enabled:      true,
				Rank:         &rankVal,
			},
			ExcFunc: func(ctx context.Context) any {
				nc := categorybus.NewCategory{
					Name:         name.MustParse("Appetizers"),
					Description:  "Delicious starters and appetizers",
					RestaurantID: sd.Restaurants[1].ID,
					Rank:         &rankVal,
				}

				resp, err := busDomain.Category.Create(ctx, nc)
				if err != nil {
					return err
				}

				return resp
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.(categorybus.Category)
				if !exists {
					return "error occurred"
				}

				expResp := exp.(categorybus.Category)

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
	rankVal := 50
	table := []unittest.Table{
		{
			Name: "set-rank",
			ExpResp: categorybus.Category{
				ID:           sd.Categories[0].ID,
				Name:         name.MustParse("Updated Category"),
				Description:  "Updated description for this category",
				RestaurantID: sd.Categories[0].RestaurantID,
				Enabled:      false,
				Rank:         &rankVal,
				DateCreated:  sd.Categories[0].DateCreated,
			},
			ExcFunc: func(ctx context.Context) any {
				uc := categorybus.UpdateCategory{
					Name:        dbtest.NamePointer("Updated Category"),
					Description: dbtest.StringPointer("Updated description for this category"),
					Enabled:     dbtest.BoolPointer(false),
					Rank:        opt.NewNullInt(50),
				}

				resp, err := busDomain.Category.Update(ctx, sd.Categories[0].Category, uc)
				if err != nil {
					return err
				}

				return resp
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.(categorybus.Category)
				if !exists {
					return "error occurred"
				}

				expResp := exp.(categorybus.Category)
				expResp.DateCreated = gotResp.DateCreated
				expResp.DateUpdated = gotResp.DateUpdated

				return cmp.Diff(gotResp, expResp)
			},
		},
		{
			Name: "clear-rank",
			ExpResp: categorybus.Category{
				ID:           sd.Categories[0].ID,
				Name:         name.MustParse("Updated Category"),
				Description:  "Updated description for this category",
				RestaurantID: sd.Categories[0].RestaurantID,
				Enabled:      false,
				Rank:         nil,
				DateCreated:  sd.Categories[0].DateCreated,
			},
			ExcFunc: func(ctx context.Context) any {
				uc := categorybus.UpdateCategory{
					Rank: opt.NewNullIntNull(),
				}

				current, err := busDomain.Category.QueryByID(ctx, sd.Categories[0].ID)
				if err != nil {
					return err
				}

				resp, err := busDomain.Category.Update(ctx, current, uc)
				if err != nil {
					return err
				}

				return resp
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.(categorybus.Category)
				if !exists {
					return "error occurred"
				}

				expResp := exp.(categorybus.Category)
				expResp.DateCreated = gotResp.DateCreated
				expResp.DateUpdated = gotResp.DateUpdated

				return cmp.Diff(gotResp, expResp)
			},
		},
	}

	return table
}

func reorder(busDomain dbtest.BusDomain, sd unittest.SeedData) []unittest.Table {
	restID := sd.Restaurants[0].ID
	cat0 := sd.Categories[0].ID
	cat1 := sd.Categories[1].ID

	table := []unittest.Table{
		{
			Name:    "success-exact-set",
			ExpResp: []int{10, 20},
			ExcFunc: func(ctx context.Context) any {
				// Reverse order: cat1 then cat0
				reordered, err := busDomain.Category.Reorder(ctx, restID, []uuid.UUID{cat1, cat0})
				if err != nil {
					return err
				}

				ranks := make([]int, len(reordered))
				for i, cat := range reordered {
					if cat.Rank != nil {
						ranks[i] = *cat.Rank
					}
				}
				return ranks
			},
			CmpFunc: func(got any, exp any) string {
				gotRanks, ok := got.([]int)
				if !ok {
					return fmt.Sprintf("expected []int, got %T", got)
				}
				return cmp.Diff(gotRanks, exp.([]int))
			},
		},
		{
			Name:    "duplicate-id-error",
			ExpResp: categorybus.ErrInvalidReorder,
			ExcFunc: func(ctx context.Context) any {
				_, err := busDomain.Category.Reorder(ctx, restID, []uuid.UUID{cat1, cat1})
				return err
			},
			CmpFunc: func(got any, exp any) string {
				gotErr, ok := got.(error)
				if !ok {
					return "expected error response"
				}
				if !errors.Is(gotErr, categorybus.ErrInvalidReorder) {
					return fmt.Sprintf("expected ErrInvalidReorder, got %v", gotErr)
				}
				return ""
			},
		},
		{
			Name:    "foreign-restaurant-id-error",
			ExpResp: categorybus.ErrInvalidReorder,
			ExcFunc: func(ctx context.Context) any {
				foreignCat := sd.Categories[2].ID
				_, err := busDomain.Category.Reorder(ctx, restID, []uuid.UUID{cat0, foreignCat})
				return err
			},
			CmpFunc: func(got any, exp any) string {
				gotErr, ok := got.(error)
				if !ok {
					return "expected error response"
				}
				if !errors.Is(gotErr, categorybus.ErrInvalidReorder) {
					return fmt.Sprintf("expected ErrInvalidReorder, got %v", gotErr)
				}
				return ""
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
				if err := busDomain.Category.Delete(ctx, sd.Categories[1].Category); err != nil {
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
