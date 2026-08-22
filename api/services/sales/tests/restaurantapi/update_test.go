package restaurantapi_test

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/go-cmp/cmp"
	restaurantapi "github.com/warlck/food-flow/app/domain/restaurantapp"
	"github.com/warlck/food-flow/app/sdk/apitest"
	"github.com/warlck/food-flow/app/sdk/errs"
	"github.com/warlck/food-flow/business/sdk/dbtest"
)

func update200(sd apitest.SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:       "basic",
			URL:        fmt.Sprintf("/v1/restaurants/%s", sd.Restaurants[0].ID),
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPut,
			StatusCode: http.StatusOK,
			Input: &restaurantapi.UpdateRestaurant{
				Name:        dbtest.StringPointer("Updated Restaurant"),
				Description: dbtest.StringPointer("Updated description for this restaurant"),
				Address:     dbtest.StringPointer("456 New Street"),
				Phone:       dbtest.StringPointer("+1-555-9999"),
				Email:       dbtest.StringPointer("updated@restaurant.com"),
				ImageURL:    dbtest.StringPointer("updated.jpg"),
				LogoURL:     dbtest.StringPointer("updated_logo.jpg"),
				Enabled:     dbtest.BoolPointer(false),
			},
			GotResp: &restaurantapi.Restaurant{},
			ExpResp: &restaurantapi.Restaurant{
				ID:                    sd.Restaurants[0].ID.String(),
				Name:                  "Updated Restaurant",
				Description:           "Updated description for this restaurant",
				Address:               "456 New Street",
				Phone:                 "+1-555-9999",
				Email:                 "updated@restaurant.com",
				ImageURL:              "updated.jpg",
				LogoURL:               "updated_logo.jpg",
				Enabled:               false,
				Latitude:              sd.Restaurants[0].Latitude,
				Longitude:             sd.Restaurants[0].Longitude,
				MaxDeliveryDistanceKm: sd.Restaurants[0].MaxDeliveryDistanceKm,
				TaxRate:               sd.Restaurants[0].TaxRate,
				DateCreated:           sd.Restaurants[0].DateCreated.Format(time.RFC3339),
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.(*restaurantapi.Restaurant)
				if !exists {
					return "error occurred"
				}

				expResp := exp.(*restaurantapi.Restaurant)
				expResp.DateUpdated = gotResp.DateUpdated

				return cmp.Diff(got, exp)
			},
		},
		{
			Name:       "logo-url-update",
			URL:        fmt.Sprintf("/v1/restaurants/%s", sd.Restaurants[1].ID),
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPut,
			StatusCode: http.StatusOK,
			Input: &restaurantapi.UpdateRestaurant{
				LogoURL: dbtest.StringPointer("brand_logo_new.png"),
			},
			GotResp: &restaurantapi.Restaurant{},
			ExpResp: &restaurantapi.Restaurant{
				ID:                    sd.Restaurants[1].ID.String(),
				Name:                  sd.Restaurants[1].Name.String(),
				Description:           sd.Restaurants[1].Description,
				Address:               sd.Restaurants[1].Address,
				Phone:                 sd.Restaurants[1].Phone,
				Email:                 sd.Restaurants[1].Email,
				ImageURL:              sd.Restaurants[1].ImageURL,
				LogoURL:               "brand_logo_new.png",
				Enabled:               sd.Restaurants[1].Enabled,
				Latitude:              sd.Restaurants[1].Latitude,
				Longitude:             sd.Restaurants[1].Longitude,
				MaxDeliveryDistanceKm: sd.Restaurants[1].MaxDeliveryDistanceKm,
				TaxRate:               sd.Restaurants[1].TaxRate,
				DateCreated:           sd.Restaurants[1].DateCreated.Format(time.RFC3339),
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.(*restaurantapi.Restaurant)
				if !exists {
					return "error occurred"
				}

				expResp := exp.(*restaurantapi.Restaurant)
				expResp.DateUpdated = gotResp.DateUpdated

				return cmp.Diff(got, exp)
			},
		},
	}

	return table
}

func update400(sd apitest.SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:       "invalid-email",
			URL:        fmt.Sprintf("/v1/restaurants/%s", sd.Restaurants[0].ID),
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPut,
			StatusCode: http.StatusBadRequest,
			Input: &restaurantapi.UpdateRestaurant{
				Email: dbtest.StringPointer("invalid-email"),
			},
			GotResp: &errs.Error{},
			ExpResp: errs.Newf(errs.InvalidArgument, "validate: [{\"field\":\"email\",\"error\":\"email must be a valid email address\"}]"),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
	}

	return table
}

func update401(sd apitest.SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:       "emptytoken",
			URL:        fmt.Sprintf("/v1/restaurants/%s", sd.Restaurants[0].ID),
			Token:      "&nbsp;",
			Method:     http.MethodPut,
			StatusCode: http.StatusUnauthorized,
			GotResp:    &errs.Error{},
			ExpResp:    errs.Newf(errs.Unauthenticated, "error parsing token: token is malformed: token contains an invalid number of segments"),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
		{
			Name:       "wronguser",
			URL:        fmt.Sprintf("/v1/restaurants/%s", sd.Restaurants[0].ID),
			Token:      sd.Users[0].Token,
			Method:     http.MethodPut,
			StatusCode: http.StatusForbidden,
			GotResp:    &errs.Error{},
			ExpResp:    errs.Newf(errs.PermissionDenied, "authorize: you are not authorized for that action, claims[[USER]] rule[rule_admin_only]"),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
	}

	return table
}
