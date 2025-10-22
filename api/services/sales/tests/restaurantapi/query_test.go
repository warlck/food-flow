package restaurantapi_test

import (
	"fmt"
	"net/http"
	"sort"

	"github.com/google/go-cmp/cmp"
	restaurantapi "github.com/warlck/food-flow/app/domain/restaurantapp"
	"github.com/warlck/food-flow/app/sdk/apitest"
	"github.com/warlck/food-flow/app/sdk/query"
)

func query200(sd apitest.SeedData) []apitest.Table {
	rests := make([]restaurantapi.Restaurant, 0, len(sd.Restaurants))

	for _, rest := range sd.Restaurants {
		rests = append(rests, toAppRestaurant(rest.Restaurant))
	}

	sort.Slice(rests, func(i, j int) bool {
		return rests[i].ID <= rests[j].ID
	})

	table := []apitest.Table{
		{
			Name:       "basic",
			URL:        "/v1/restaurants?page=1&rows=10&orderBy=restaurant_id,ASC",
			Token:      sd.Admins[0].Token,
			StatusCode: http.StatusOK,
			Method:     http.MethodGet,
			GotResp:    &query.Result[restaurantapi.Restaurant]{},
			ExpResp: &query.Result[restaurantapi.Restaurant]{
				Page:        1,
				RowsPerPage: 10,
				Total:       len(rests),
				Items:       rests,
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
			URL:        fmt.Sprintf("/v1/restaurants/%s", sd.Restaurants[0].ID),
			Token:      sd.Admins[0].Token,
			StatusCode: http.StatusOK,
			Method:     http.MethodGet,
			GotResp:    &restaurantapi.Restaurant{},
			ExpResp:    toAppRestaurantPtr(sd.Restaurants[0].Restaurant),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
	}

	return table
}
