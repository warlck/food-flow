package menuitemapi_test

import (
	"net/http"

	"github.com/google/go-cmp/cmp"
	"github.com/warlck/food-flow/app/domain/menuitemapp"
	"github.com/warlck/food-flow/app/sdk/apitest"
	"github.com/warlck/food-flow/app/sdk/errs"
	"github.com/warlck/food-flow/business/sdk/dbtest"
)

func create201(sd apitest.SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:       "basic",
			URL:        "/v1/menuitems",
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPost,
			StatusCode: http.StatusCreated,
			Input: &menuitemapp.NewMenuItem{
				CategoryID:   sd.Categories[0].ID.String(),
				RestaurantID: sd.Restaurants[0].ID.String(),
				Name:         "Test MenuItem",
				Description:  "Test Description",
				Price:        19.99,
			},
			GotResp: &menuitemapp.MenuItem{},
			ExpResp: &menuitemapp.MenuItem{},
			CmpFunc: func(got any, exp any) string {
				gotResp := got.(*menuitemapp.MenuItem)
				if gotResp.ID == "" {
					return "id should not be empty"
				}
				if gotResp.Name != "Test MenuItem" {
					return "name mismatch"
				}
				if gotResp.Price != 19.99 {
					return "price mismatch"
				}
				if gotResp.Rank != nil {
					return "rank should be nil when not provided"
				}
				return ""
			},
		},
		{
			Name:       "with-rank",
			URL:        "/v1/menuitems",
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPost,
			StatusCode: http.StatusCreated,
			Input: &menuitemapp.NewMenuItem{
				CategoryID:   sd.Categories[0].ID.String(),
				RestaurantID: sd.Restaurants[0].ID.String(),
				Name:         "Test MenuItem Ranked",
				Description:  "Test Description",
				Price:        24.99,
				Rank:         dbtest.IntPointer(5),
			},
			GotResp: &menuitemapp.MenuItem{},
			ExpResp: &menuitemapp.MenuItem{},
			CmpFunc: func(got any, exp any) string {
				gotResp := got.(*menuitemapp.MenuItem)
				if gotResp.Rank == nil || *gotResp.Rank != 5 {
					return "rank mismatch"
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
			URL:        "/v1/menuitems",
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPost,
			StatusCode: http.StatusBadRequest,
			Input: &menuitemapp.NewMenuItem{
				CategoryID: sd.Categories[0].ID.String(),
			},
			GotResp: &errs.Error{},
			ExpResp: &errs.Error{
				Code:    errs.InvalidArgument,
				Message: "validate:",
			},
			CmpFunc: func(got any, exp any) string {
				gotErr := got.(*errs.Error)
				expErr := exp.(*errs.Error)

				if gotErr.Code != expErr.Code {
					return "error code mismatch"
				}

				// Just check that the message starts with "validate:"
				if len(gotErr.Message) < len(expErr.Message) {
					return "error message too short"
				}
				if gotErr.Message[:len(expErr.Message)] != expErr.Message {
					return "error message doesn't start with 'validate:'"
				}

				return ""
			},
		},
		{
			Name:       "rank-zero",
			URL:        "/v1/menuitems",
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPost,
			StatusCode: http.StatusBadRequest,
			Input: &menuitemapp.NewMenuItem{
				CategoryID:   sd.Categories[0].ID.String(),
				RestaurantID: sd.Restaurants[0].ID.String(),
				Name:         "Test MenuItem",
				Price:        19.99,
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
			URL:        "/v1/menuitems",
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPost,
			StatusCode: http.StatusBadRequest,
			Input: &menuitemapp.NewMenuItem{
				CategoryID:   sd.Categories[0].ID.String(),
				RestaurantID: sd.Restaurants[0].ID.String(),
				Name:         "Test MenuItem",
				Price:        19.99,
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
			Name:       "category-from-different-restaurant",
			URL:        "/v1/menuitems",
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPost,
			StatusCode: http.StatusBadRequest,
			Input: &menuitemapp.NewMenuItem{
				CategoryID:   sd.Categories[2].ID.String(),
				RestaurantID: sd.Restaurants[0].ID.String(),
				Name:         "Cross Restaurant Item",
				Description:  "invalid ownership fixture",
				Price:        19.99,
			},
			GotResp: &errs.Error{},
			ExpResp: &errs.Error{
				Code: errs.InvalidArgument,
			},
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got.(*errs.Error).Code, exp.(*errs.Error).Code)
			},
		},
	}

	return table
}

func create401(sd apitest.SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:       "emptytoken",
			URL:        "/v1/menuitems",
			Token:      "",
			Method:     http.MethodPost,
			StatusCode: http.StatusUnauthorized,
			GotResp:    &errs.Error{},
			ExpResp: &errs.Error{
				Code:    errs.Unauthenticated,
				Message: "expected authorization header format: Bearer <token>",
			},
			CmpFunc: func(got any, exp any) string {
				gotErr := got.(*errs.Error)
				expErr := exp.(*errs.Error)

				if gotErr.Code != expErr.Code {
					return "error code mismatch"
				}

				if gotErr.Message != expErr.Message {
					return "error message mismatch"
				}

				return ""
			},
		},
		{
			Name:       "wronguser",
			URL:        "/v1/menuitems",
			Token:      sd.Users[0].Token,
			Method:     http.MethodPost,
			StatusCode: http.StatusForbidden,
			Input: &menuitemapp.NewMenuItem{
				CategoryID:   sd.Categories[0].ID.String(),
				RestaurantID: sd.Restaurants[0].ID.String(),
				Name:         "Test MenuItem",
				Description:  "Test Description",
				Price:        19.99,
			},
			GotResp: &errs.Error{},
			ExpResp: &errs.Error{
				Code:    errs.PermissionDenied,
				Message: "authorize: you are not authorized for that action, claims[[USER]] rule[rule_admin_only]",
			},
			CmpFunc: func(got any, exp any) string {
				gotErr := got.(*errs.Error)
				expErr := exp.(*errs.Error)

				if gotErr.Code != expErr.Code {
					return "error code mismatch"
				}

				if gotErr.Message != expErr.Message {
					return "error message mismatch"
				}

				return ""
			},
		},
	}

	return table
}
