package restaurantapi_test

import (
	"net/http"

	"github.com/google/go-cmp/cmp"
	restaurantapi "github.com/warlck/food-flow/app/domain/restaurantapp"
	"github.com/warlck/food-flow/app/sdk/apitest"
	"github.com/warlck/food-flow/app/sdk/errs"
)

func create201(sd apitest.SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:       "basic",
			URL:        "/v1/restaurants",
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPost,
			StatusCode: http.StatusCreated,
			Input: &restaurantapi.NewRestaurant{
				OrganizationID: sd.Restaurants[0].OrganizationID.String(),
				Name:           "New Restaurant",
				Description:    "A wonderful new restaurant with great food",
				Address:        "123 Main Street",
				Phone:          "+1-555-0100",
				Email:          "info@newrestaurant.com",
				ImageURL:       "restaurant.jpg",
			},
			GotResp: &restaurantapi.Restaurant{},
			ExpResp: &restaurantapi.Restaurant{
				Name:        "New Restaurant",
				Description: "A wonderful new restaurant with great food",
				Address:     "123 Main Street",
				Phone:       "+1-555-0100",
				Email:       "info@newrestaurant.com",
				ImageURL:    "restaurant.jpg",
				Enabled:     true,
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.(*restaurantapi.Restaurant)
				if !exists {
					return "error occurred"
				}

				expResp := exp.(*restaurantapi.Restaurant)

				expResp.ID = gotResp.ID
				expResp.DateCreated = gotResp.DateCreated
				expResp.DateUpdated = gotResp.DateUpdated

				return cmp.Diff(got, exp)
			},
		},
	}

	return table
}

func create400(sd apitest.SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:       "missing-input",
			URL:        "/v1/restaurants",
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPost,
			StatusCode: http.StatusBadRequest,
			Input:      &restaurantapi.NewRestaurant{},
			GotResp:    &errs.Error{},
			ExpResp:    errs.Newf(errs.InvalidArgument, "validate: [{\"field\":\"organizationId\",\"error\":\"organizationId is a required field\"},{\"field\":\"name\",\"error\":\"name is a required field\"},{\"field\":\"address\",\"error\":\"address is a required field\"},{\"field\":\"phone\",\"error\":\"phone is a required field\"},{\"field\":\"email\",\"error\":\"email is a required field\"}]"),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
		{
			Name:       "invalid-email",
			URL:        "/v1/restaurants",
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPost,
			StatusCode: http.StatusBadRequest,
			Input: &restaurantapi.NewRestaurant{
				OrganizationID: sd.Restaurants[0].OrganizationID.String(),
				Name:           "Test Restaurant",
				Address:        "123 Main St",
				Phone:          "+1-555-0100",
				Email:          "invalid-email",
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

func create401(sd apitest.SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:       "emptytoken",
			URL:        "/v1/restaurants",
			Token:      "&nbsp;",
			Method:     http.MethodPost,
			StatusCode: http.StatusUnauthorized,
			GotResp:    &errs.Error{},
			ExpResp:    errs.Newf(errs.Unauthenticated, "error parsing token: token is malformed: token contains an invalid number of segments"),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
		{
			Name:       "wronguser",
			URL:        "/v1/restaurants",
			Token:      sd.Users[0].Token,
			Method:     http.MethodPost,
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
