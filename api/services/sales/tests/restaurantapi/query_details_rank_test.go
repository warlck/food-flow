package restaurantapi_test

import (
	"fmt"
	"net/http"

	restaurantapi "github.com/warlck/food-flow/app/domain/restaurantapp"
	"github.com/warlck/food-flow/app/sdk/apitest"
)

// queryByIDWithDetailsRanked200 verifies the details endpoint orders ranked
// entries first (ascending rank) and unranked entries last, for both menu
// items and addons. The ranked fixtures live on a dedicated restaurant,
// Restaurants[4] / Categories[2], so the unranked seeds on Restaurants[0]
// keep proving backward compatibility and no other test mutates them.
//
// Seed layout for the ranked category (see seed_test.go):
//   - MenuItems[5] -> rank 20, MenuItems[6] -> rank 10, MenuItems[7] -> unranked
//   - Addons[0] -> rank 20, Addons[1] -> rank 10, Addons[2] -> unranked
func queryByIDWithDetailsRanked200(sd apitest.SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:       "ranked-first-ordering",
			URL:        fmt.Sprintf("/v1/restaurants/%s/details", sd.Restaurants[4].ID),
			Token:      sd.Admins[0].Token,
			StatusCode: http.StatusOK,
			Method:     http.MethodGet,
			GotResp:    &restaurantapi.RestaurantWithMenuCategories{},
			ExpResp:    &restaurantapi.RestaurantWithMenuCategories{},
			CmpFunc: func(got any, exp any) string {
				gotResp := got.(*restaurantapi.RestaurantWithMenuCategories)

				if len(gotResp.Categories) != 1 {
					return fmt.Sprintf("expected 1 category, got %d", len(gotResp.Categories))
				}

				cat := gotResp.Categories[0]
				if cat.ID != sd.Categories[2].ID.String() {
					return fmt.Sprintf("expected ranked category %s, got %s", sd.Categories[2].ID, cat.ID)
				}

				// Menu items: ranked ascending first (10, 20), unranked last.
				wantItems := []struct {
					id   string
					rank *int
				}{
					{sd.MenuItems[6].ID.String(), sd.MenuItems[6].Rank},
					{sd.MenuItems[5].ID.String(), sd.MenuItems[5].Rank},
					{sd.MenuItems[7].ID.String(), nil},
				}

				if len(cat.MenuItems) != len(wantItems) {
					return fmt.Sprintf("expected %d menu items, got %d", len(wantItems), len(cat.MenuItems))
				}

				for i, want := range wantItems {
					gotItem := cat.MenuItems[i]
					if gotItem.ID != want.id {
						return fmt.Sprintf("menu item position %d: expected %s, got %s", i, want.id, gotItem.ID)
					}
					if (gotItem.Rank == nil) != (want.rank == nil) {
						return fmt.Sprintf("menu item position %d: rank nil mismatch (got %+v)", i, gotItem.Rank)
					}
					if gotItem.Rank != nil && *gotItem.Rank != *want.rank {
						return fmt.Sprintf("menu item position %d: expected rank %d, got %d", i, *want.rank, *gotItem.Rank)
					}
				}

				// Addons: shared per category; ranked ascending first, unranked last.
				wantAddons := []struct {
					id   string
					rank *int
				}{
					{sd.Addons[1].ID.String(), sd.Addons[1].Rank},
					{sd.Addons[0].ID.String(), sd.Addons[0].Rank},
					{sd.Addons[2].ID.String(), nil},
				}

				// The backend attaches the same shared addon list to every menu
				// item in the category, so assert the ordering on all of them —
				// a bug that ordered only the first item's addons would slip
				// past a MenuItems[0]-only check.
				for _, gotItem := range cat.MenuItems {
					gotAddons := gotItem.Addons
					if len(gotAddons) != len(wantAddons) {
						return fmt.Sprintf("menu item %s: expected %d addons, got %d", gotItem.ID, len(wantAddons), len(gotAddons))
					}

					for i, want := range wantAddons {
						gotAddon := gotAddons[i]
						if gotAddon.ID != want.id {
							return fmt.Sprintf("menu item %s addon position %d: expected %s, got %s", gotItem.ID, i, want.id, gotAddon.ID)
						}
						if (gotAddon.Rank == nil) != (want.rank == nil) {
							return fmt.Sprintf("menu item %s addon position %d: rank nil mismatch (got %+v)", gotItem.ID, i, gotAddon.Rank)
						}
						if gotAddon.Rank != nil && *gotAddon.Rank != *want.rank {
							return fmt.Sprintf("menu item %s addon position %d: expected rank %d, got %d", gotItem.ID, i, *want.rank, *gotAddon.Rank)
						}
					}
				}

				return ""
			},
		},
	}

	return table
}
