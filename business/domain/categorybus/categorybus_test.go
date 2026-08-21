package categorybus_test

import (
	"context"
	"fmt"
	"github.com/warlck/food-flow/business/domain/organizationbus"
	"sort"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/warlck/food-flow/business/domain/categorybus"
	"github.com/warlck/food-flow/business/domain/restaurantbus"
	"github.com/warlck/food-flow/business/sdk/dbtest"
	"github.com/warlck/food-flow/business/sdk/page"
	"github.com/warlck/food-flow/business/sdk/unittest"
	"github.com/warlck/food-flow/business/types/name"
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
	unittest.Run(t, create(db.BusDomain, sd), "create")
	unittest.Run(t, update(db.BusDomain, sd), "update")
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
		return cats[i].ID.String() <= cats[j].ID.String()
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

func create(busDomain dbtest.BusDomain, sd unittest.SeedData) []unittest.Table {
	table := []unittest.Table{
		{
			Name: "basic",
			ExpResp: categorybus.Category{
				Name:         name.MustParse("Appetizers"),
				Description:  "Delicious starters and appetizers",
				RestaurantID: sd.Restaurants[0].ID,
				Enabled:      true,
			},
			ExcFunc: func(ctx context.Context) any {
				nc := categorybus.NewCategory{
					Name:         name.MustParse("Appetizers"),
					Description:  "Delicious starters and appetizers",
					RestaurantID: sd.Restaurants[0].ID,
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
	table := []unittest.Table{
		{
			Name: "basic",
			ExpResp: categorybus.Category{
				ID:           sd.Categories[0].ID,
				Name:         name.MustParse("Updated Category"),
				Description:  "Updated description for this category",
				RestaurantID: sd.Categories[0].RestaurantID,
				Enabled:      false,
				DateCreated:  sd.Categories[0].DateCreated,
			},
			ExcFunc: func(ctx context.Context) any {
				uc := categorybus.UpdateCategory{
					Name:        dbtest.NamePointer("Updated Category"),
					Description: dbtest.StringPointer("Updated description for this category"),
					Enabled:     dbtest.BoolPointer(false),
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
