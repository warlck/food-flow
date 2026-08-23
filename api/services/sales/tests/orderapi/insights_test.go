package orderapi_test

import (
	"net/http"

	"github.com/google/go-cmp/cmp"
	"github.com/warlck/food-flow/app/domain/orderapp"
	"github.com/warlck/food-flow/app/sdk/apitest"
	"github.com/warlck/food-flow/app/sdk/errs"
)

func insights200(sd apitest.SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:       "insights-all-restaurants",
			URL:        "/v1/insights",
			Token:      sd.Admins[0].Token,
			StatusCode: http.StatusOK,
			Method:     http.MethodGet,
			GotResp:    &orderapp.AppInsights{},
			ExpResp:    &orderapp.AppInsights{},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.(*orderapp.AppInsights)
				if !exists {
					return "error occurred"
				}

				if gotResp.Summary.TotalOrders < 0 {
					return "expected non-negative total orders"
				}

				return ""
			},
		},
		{
			Name:       "insights-with-restaurant-filter",
			URL:        "/v1/insights?restaurant_id=" + sd.Restaurants[0].ID.String(),
			Token:      sd.Admins[0].Token,
			StatusCode: http.StatusOK,
			Method:     http.MethodGet,
			GotResp:    &orderapp.AppInsights{},
			ExpResp:    &orderapp.AppInsights{},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.(*orderapp.AppInsights)
				if !exists {
					return "error occurred"
				}

				if gotResp.Summary.TotalOrders < 0 {
					return "expected non-negative total orders"
				}

				return ""
			},
		},
	}

	return table
}

func insights401(sd apitest.SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:       "non-admin-user",
			URL:        "/v1/insights",
			Token:      sd.Users[0].Token,
			Method:     http.MethodGet,
			StatusCode: http.StatusForbidden,
			GotResp:    &errs.Error{},
			ExpResp:    errs.Newf(errs.PermissionDenied, "authorize: you are not authorized for that action, claims[[USER]] rule[rule_admin_only]"),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
		{
			Name:       "empty-token",
			URL:        "/v1/insights",
			Token:      "&nbsp;",
			Method:     http.MethodGet,
			StatusCode: http.StatusUnauthorized,
			GotResp:    &errs.Error{},
			ExpResp:    errs.Newf(errs.Unauthenticated, "error parsing token: token is malformed: token contains an invalid number of segments"),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
	}

	return table
}
