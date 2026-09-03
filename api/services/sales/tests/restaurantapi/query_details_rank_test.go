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

				// Addons: owned per item; ranked ascending first (10, 20), unranked last.
				rank10 := 10
				rank20 := 20
				wantRanks := []*int{&rank10, &rank20, nil}

				for _, gotItem := range cat.MenuItems {
					gotAddons := gotItem.Addons
					if len(gotAddons) != len(wantRanks) {
						return fmt.Sprintf("menu item %s: expected %d addons, got %d", gotItem.ID, len(wantRanks), len(gotAddons))
					}

					for i, wantRank := range wantRanks {
						gotAddon := gotAddons[i]
						if (gotAddon.Rank == nil) != (wantRank == nil) {
							return fmt.Sprintf("menu item %s addon position %d: rank nil mismatch (got %+v)", gotItem.ID, i, gotAddon.Rank)
						}
						if gotAddon.Rank != nil && *gotAddon.Rank != *wantRank {
							return fmt.Sprintf("menu item %s addon position %d: expected rank %d, got %d", gotItem.ID, i, *wantRank, *gotAddon.Rank)
						}
					}
				}

				return ""
			},
		},
	}

	return table
}
