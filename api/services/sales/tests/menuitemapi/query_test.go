package menuitemapi_test

import (
	"net/http"
	"sort"

	"github.com/google/go-cmp/cmp"
	"github.com/warlck/food-flow/app/domain/menuitemapp"
	"github.com/warlck/food-flow/app/sdk/apitest"
	"github.com/warlck/food-flow/app/sdk/query"
	"github.com/warlck/food-flow/business/domain/menuitembus"
)

func query200(sd apitest.SeedData) []apitest.Table {
	// Restaurant 0 holds MenuItems[0..3]; indices 4..5 belong to restaurant 1
	// and 6..7 to the second organization (cross-org test fixtures).
	items := toAppMenuItems([]menuitembus.MenuItem{
		sd.MenuItems[0].MenuItem,
		sd.MenuItems[1].MenuItem,
		sd.MenuItems[2].MenuItem,
		sd.MenuItems[3].MenuItem,
	})

	sort.Slice(items, func(i, j int) bool {
		return items[i].ID < items[j].ID
	})

	table := []apitest.Table{
		{
			Name: "basic",
			// Scoped to the first restaurant: the seed also creates menu items
			// in a second organization (for cross-org authorization tests),
			// and the list endpoint is intentionally not org-scoped.
			URL:        "/v1/menuitems?page=1&rows=20&orderBy=menu_item_id,ASC&restaurant_id=" + sd.Restaurants[0].ID.String(),
			Token:      sd.Admins[0].Token,
			StatusCode: http.StatusOK,
			Method:     http.MethodGet,
			GotResp:    &query.Result[menuitemapp.MenuItem]{},
			ExpResp: &query.Result[menuitemapp.MenuItem]{
				Page:        1,
				RowsPerPage: 20,
				Total:       len(items),
				Items:       items,
			},
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
	}

	return table
}

func queryByID200(sd apitest.SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:       "basic",
			URL:        "/v1/menuitems/" + sd.MenuItems[0].ID.String(),
			Token:      sd.Admins[0].Token,
			StatusCode: http.StatusOK,
			Method:     http.MethodGet,
			GotResp:    &menuitemapp.MenuItem{},
			ExpResp:    toAppMenuItemPtr(sd.MenuItems[0].MenuItem),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
	}

	return table
}
