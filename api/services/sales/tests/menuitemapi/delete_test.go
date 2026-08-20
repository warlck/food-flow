package menuitemapi_test

import (
	"net/http"

	"github.com/warlck/food-flow/app/sdk/apitest"
	"github.com/warlck/food-flow/app/sdk/errs"
)

func delete200(sd apitest.SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:       "basic",
			URL:        "/v1/menuitems/" + sd.MenuItems[3].ID.String(),
			Token:      sd.Admins[0].Token,
			Method:     http.MethodDelete,
			StatusCode: http.StatusNoContent,
		},
	}

	return table
}

func delete401(sd apitest.SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:       "emptytoken",
			URL:        "/v1/menuitems/" + sd.MenuItems[3].ID.String(),
			Token:      "",
			Method:     http.MethodDelete,
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
			URL:        "/v1/menuitems/" + sd.MenuItems[2].ID.String(),
			Token:      sd.Users[0].Token,
			Method:     http.MethodDelete,
			StatusCode: http.StatusForbidden,
			GotResp:    &errs.Error{},
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
