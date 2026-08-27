package addonapi_test

import (
	"net/http"

	"github.com/google/go-cmp/cmp"
	"github.com/warlck/food-flow/app/domain/addonapp"
	"github.com/warlck/food-flow/app/sdk/apitest"
	"github.com/warlck/food-flow/app/sdk/errs"
	"github.com/warlck/food-flow/business/sdk/dbtest"
)

func create201(sd apitest.SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:       "basic",
			URL:        "/v1/addons",
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPost,
			StatusCode: http.StatusCreated,
			Input: &addonapp.NewAddon{
				CategoryID:   sd.Categories[0].ID.String(),
				RestaurantID: sd.Restaurants[0].ID.String(),
				Name:         "Extra Special Sauce",
				Description:  "Special homemade sauce",
				Price:        2.50,
				MaxQuantity:  5,
			},
			GotResp: &addonapp.Addon{},
			ExpResp: &addonapp.Addon{
				CategoryID:   sd.Categories[0].ID.String(),
				RestaurantID: sd.Restaurants[0].ID.String(),
				Name:         "Extra Special Sauce",
				Description:  "Special homemade sauce",
				Price:        2.50,
				Available:    true,
				MaxQuantity:  5,
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.(*addonapp.Addon)
				if !exists {
					return "got is not *addonapp.Addon"
				}

				expResp := exp.(*addonapp.Addon)

				expResp.ID = gotResp.ID
				expResp.DateCreated = gotResp.DateCreated
				expResp.DateUpdated = gotResp.DateUpdated

				return cmp.Diff(got, exp)
			},
		},
		{
			Name:       "with-rank",
			URL:        "/v1/addons",
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPost,
			StatusCode: http.StatusCreated,
			Input: &addonapp.NewAddon{
				CategoryID:   sd.Categories[0].ID.String(),
				RestaurantID: sd.Restaurants[0].ID.String(),
				Name:         "Extra Pickles",
				Description:  "Pickle slices",
				Price:        1.00,
				MaxQuantity:  3,
				Rank:         func(i int) *int { return &i }(15),
			},
			GotResp: &addonapp.Addon{},
			ExpResp: &addonapp.Addon{},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.(*addonapp.Addon)
				if !exists {
					return "got is not *addonapp.Addon"
				}
				if gotResp.Rank == nil || *gotResp.Rank != 15 {
					return "rank mismatch"
				}
				return ""
			},
		},
		{
			Name:       "with-brackets-and-punctuation",
			URL:        "/v1/addons",
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPost,
			StatusCode: http.StatusCreated,
			Input: &addonapp.NewAddon{
				CategoryID:   sd.Categories[0].ID.String(),
				RestaurantID: sd.Restaurants[0].ID.String(),
				Name:         "Extra Sauce {Special} + Garlic & Herb, 0.5L (Cold)",
				Description:  "House blend",
				Price:        3.50,
				MaxQuantity:  5,
			},
			GotResp: &addonapp.Addon{},
			ExpResp: &addonapp.Addon{},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.(*addonapp.Addon)
				if !exists {
					return "got is not *addonapp.Addon"
				}
				if gotResp.Name != "Extra Sauce {Special} + Garlic & Herb, 0.5L (Cold)" {
					return "name mismatch"
				}
				if gotResp.Price != 3.50 {
					return "price mismatch"
				}
				return ""
			},
		},
	}

	return table
}

func create400(sd apitest.SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:       "missing-input",
			URL:        "/v1/addons",
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPost,
			StatusCode: http.StatusBadRequest,
			Input:      &addonapp.NewAddon{},
			GotResp:    &errs.Error{},
			ExpResp: &errs.Error{
				Code:    errs.InvalidArgument,
				Message: `validate: [{"field":"categoryId","error":"categoryId is a required field"},{"field":"restaurantId","error":"restaurantId is a required field"},{"field":"name","error":"name is a required field"},{"field":"price","error":"price is a required field"}]`,
			},
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
		{
			Name:       "rank-zero",
			URL:        "/v1/addons",
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPost,
			StatusCode: http.StatusBadRequest,
			Input: &addonapp.NewAddon{
				CategoryID:   sd.Categories[0].ID.String(),
				RestaurantID: sd.Restaurants[0].ID.String(),
				Name:         "Extra Sauce",
				Price:        2.50,
				Rank:         dbtest.IntPointer(0),
			},
			GotResp: &errs.Error{},
			ExpResp: &errs.Error{
				Code:    errs.InvalidArgument,
				Message: `validate: [{"field":"rank","error":"rank must be 1 or greater"}]`,
			},
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
		{
			Name:       "rank-negative",
			URL:        "/v1/addons",
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPost,
			StatusCode: http.StatusBadRequest,
			Input: &addonapp.NewAddon{
				CategoryID:   sd.Categories[0].ID.String(),
				RestaurantID: sd.Restaurants[0].ID.String(),
				Name:         "Extra Sauce",
				Price:        2.50,
				Rank:         dbtest.IntPointer(-1),
			},
			GotResp: &errs.Error{},
			ExpResp: &errs.Error{
				Code:    errs.InvalidArgument,
				Message: `validate: [{"field":"rank","error":"rank must be 1 or greater"}]`,
			},
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
		{
			Name:       "invalid-name-disallowed-char",
			URL:        "/v1/addons",
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPost,
			StatusCode: http.StatusBadRequest,
			Input: &addonapp.NewAddon{
				CategoryID:   sd.Categories[0].ID.String(),
				RestaurantID: sd.Restaurants[0].ID.String(),
				Name:         "Extra Dip #1",
				Price:        2.50,
			},
			GotResp: &errs.Error{},
			ExpResp: &errs.Error{
				Code:    errs.InvalidArgument,
				Message: `parse name: invalid name "Extra Dip #1"`,
			},
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
		{
			Name:       "invalid-name-angle-brackets",
			URL:        "/v1/addons",
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPost,
			StatusCode: http.StatusBadRequest,
			Input: &addonapp.NewAddon{
				CategoryID:   sd.Categories[0].ID.String(),
				RestaurantID: sd.Restaurants[0].ID.String(),
				Name:         "Dip <Special>",
				Price:        2.50,
			},
			GotResp: &errs.Error{},
			ExpResp: &errs.Error{
				Code:    errs.InvalidArgument,
				Message: `parse name: invalid name "Dip <Special>"`,
			},
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
			URL:        "/v1/addons",
			Token:      "",
			Method:     http.MethodPost,
			StatusCode: http.StatusUnauthorized,
			GotResp:    &errs.Error{},
			ExpResp: &errs.Error{
				Code:    errs.Unauthenticated,
				Message: "expected authorization header format: Bearer <token>",
			},
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
		{
			Name:       "user-role-unauthorized",
			URL:        "/v1/addons",
			Token:      sd.Users[0].Token,
			Method:     http.MethodPost,
			StatusCode: http.StatusForbidden,
			GotResp:    &errs.Error{},
			ExpResp: &errs.Error{
				Code:    errs.PermissionDenied,
				Message: "authorize: you are not authorized for that action, claims[[USER]] rule[rule_admin_only]",
			},
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
	}

	return table
}
