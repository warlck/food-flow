package restaurantapi_test

import (
	"fmt"
	"net/http"
	"sort"

	"github.com/google/go-cmp/cmp"
	restaurantapi "github.com/warlck/food-flow/app/domain/restaurantapp"
	"github.com/warlck/food-flow/app/sdk/apitest"
)

func queryByIDWithDetails200(sd apitest.SeedData) []apitest.Table {
	// Build expected categories with menu items
	expectedCategories := make([]restaurantapi.Category, 2)

	// Category 0 with 3 menu items
	expectedCategories[0] = restaurantapi.Category{
		ID:          sd.Categories[0].ID.String(),
		Name:        sd.Categories[0].Name.String(),
		Description: sd.Categories[0].Description,
		Enabled:     sd.Categories[0].Enabled,
		MenuItems: []restaurantapi.MenuItem{
			{
				ID:          sd.MenuItems[0].ID.String(),
				Name:        sd.MenuItems[0].Name.String(),
				Description: sd.MenuItems[0].Description,
				Price:       sd.MenuItems[0].Price.Value(),
				ImageURL:    sd.MenuItems[0].ImageURL,
				Available:   sd.MenuItems[0].Available,
				Addons:      []restaurantapi.Addon{},
			},
			{
				ID:          sd.MenuItems[1].ID.String(),
				Name:        sd.MenuItems[1].Name.String(),
				Description: sd.MenuItems[1].Description,
				Price:       sd.MenuItems[1].Price.Value(),
				ImageURL:    sd.MenuItems[1].ImageURL,
				Available:   sd.MenuItems[1].Available,
				Addons:      []restaurantapi.Addon{},
			},
			{
				ID:          sd.MenuItems[2].ID.String(),
				Name:        sd.MenuItems[2].Name.String(),
				Description: sd.MenuItems[2].Description,
				Price:       sd.MenuItems[2].Price.Value(),
				ImageURL:    sd.MenuItems[2].ImageURL,
				Available:   sd.MenuItems[2].Available,
				Addons:      []restaurantapi.Addon{},
			},
		},
	}

	// Category 1 with 2 menu items
	expectedCategories[1] = restaurantapi.Category{
		ID:          sd.Categories[1].ID.String(),
		Name:        sd.Categories[1].Name.String(),
		Description: sd.Categories[1].Description,
		Enabled:     sd.Categories[1].Enabled,
		MenuItems: []restaurantapi.MenuItem{
			{
				ID:          sd.MenuItems[3].ID.String(),
				Name:        sd.MenuItems[3].Name.String(),
				Description: sd.MenuItems[3].Description,
				Price:       sd.MenuItems[3].Price.Value(),
				ImageURL:    sd.MenuItems[3].ImageURL,
				Available:   sd.MenuItems[3].Available,
				Addons:      []restaurantapi.Addon{},
			},
			{
				ID:          sd.MenuItems[4].ID.String(),
				Name:        sd.MenuItems[4].Name.String(),
				Description: sd.MenuItems[4].Description,
				Price:       sd.MenuItems[4].Price.Value(),
				ImageURL:    sd.MenuItems[4].ImageURL,
				Available:   sd.MenuItems[4].Available,
				Addons:      []restaurantapi.Addon{},
			},
		},
	}

	// Sort menu items within each category by ID for consistent comparison
	for i := range expectedCategories {
		sort.Slice(expectedCategories[i].MenuItems, func(a, b int) bool {
			return expectedCategories[i].MenuItems[a].ID < expectedCategories[i].MenuItems[b].ID
		})
	}

	// Sort categories by ID
	sort.Slice(expectedCategories, func(i, j int) bool {
		return expectedCategories[i].ID < expectedCategories[j].ID
	})

	expected := &restaurantapi.RestaurantWithMenuCategories{
		ID:          sd.Restaurants[0].ID.String(),
		Name:        sd.Restaurants[0].Name.String(),
		Description: sd.Restaurants[0].Description,
		Address:     sd.Restaurants[0].Address,
		Phone:       sd.Restaurants[0].Phone,
		Email:       sd.Restaurants[0].Email,
		ImageURL:    sd.Restaurants[0].ImageURL,
		Enabled:     sd.Restaurants[0].Enabled,
		Categories:  expectedCategories,
		DateCreated: sd.Restaurants[0].DateCreated.Format("2006-01-02T15:04:05Z07:00"),
		DateUpdated: sd.Restaurants[0].DateUpdated.Format("2006-01-02T15:04:05Z07:00"),
	}

	table := []apitest.Table{
		{
			Name:       "basic",
			URL:        fmt.Sprintf("/v1/restaurants/%s/details", sd.Restaurants[0].ID),
			Token:      sd.Admins[0].Token,
			StatusCode: http.StatusOK,
			Method:     http.MethodGet,
			GotResp:    &restaurantapi.RestaurantWithMenuCategories{},
			ExpResp:    expected,
			CmpFunc: func(got any, exp any) string {
				gotResp := got.(*restaurantapi.RestaurantWithMenuCategories)

				// Assert menu items are returned sorted by price (cheapest first) within each category.
				for _, cat := range gotResp.Categories {
					for i := 1; i < len(cat.MenuItems); i++ {
						if cat.MenuItems[i-1].Price > cat.MenuItems[i].Price {
							return fmt.Sprintf("menu items not sorted by price in category %s", cat.ID)
						}
					}
				}

				// Sort categories and menu items in the response for consistent comparison.
				sort.Slice(gotResp.Categories, func(i, j int) bool {
					return gotResp.Categories[i].ID < gotResp.Categories[j].ID
				})

				for i := range gotResp.Categories {
					sort.Slice(gotResp.Categories[i].MenuItems, func(a, b int) bool {
						return gotResp.Categories[i].MenuItems[a].ID < gotResp.Categories[i].MenuItems[b].ID
					})
				}

				return cmp.Diff(got, exp)
			},
		},
	}

	return table
}

func queryByIDWithDetails400(sd apitest.SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:       "bad-restaurant-id",
			URL:        "/v1/restaurants/invalid-uuid/details",
			Token:      sd.Admins[0].Token,
			StatusCode: http.StatusBadRequest,
			Method:     http.MethodGet,
			GotResp:    &restaurantapi.RestaurantWithMenuCategories{},
			ExpResp:    &restaurantapi.RestaurantWithMenuCategories{},
			CmpFunc: func(got any, exp any) string {
				return "" // We just check status code for error cases
			},
		},
	}

	return table
}

func queryByIDWithDetails401(sd apitest.SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:       "missing-token",
			URL:        fmt.Sprintf("/v1/restaurants/%s/details", sd.Restaurants[0].ID),
			Token:      "",
			StatusCode: http.StatusUnauthorized,
			Method:     http.MethodGet,
			GotResp:    &restaurantapi.RestaurantWithMenuCategories{},
			ExpResp:    &restaurantapi.RestaurantWithMenuCategories{},
			CmpFunc: func(got any, exp any) string {
				return "" // We just check status code for auth errors
			},
		},
	}

	return table
}
