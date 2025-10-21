package categoryapi_test

import (
	"fmt"
	"net/http"

	"github.com/google/go-cmp/cmp"
	categoryapi "github.com/warlck/food-flow/app/domain/categoryapp"
	"github.com/warlck/food-flow/app/sdk/apitest"
	"github.com/warlck/food-flow/app/sdk/errs"
	"github.com/warlck/food-flow/business/sdk/dbtest"
)

func update200(sd apitest.SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:       "basic",
			URL:        fmt.Sprintf("/v1/categories/%s", sd.Categories[0].ID),
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPut,
			StatusCode: http.StatusOK,
			Input: &categoryapi.UpdateCategory{
				Name:        dbtest.StringPointer("Updated Category"),
				Description: dbtest.StringPointer("Updated description"),
				Enabled:     dbtest.BoolPointer(false),
			},
			GotResp: &categoryapi.Category{},
			ExpResp: &categoryapi.Category{
				ID:           sd.Categories[0].ID.String(),
				Name:         "Updated Category",
				Description:  "Updated description",
				RestaurantID: sd.Categories[0].RestaurantID.String(),
				Enabled:      false,
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.(*categoryapi.Category)
				if !exists {
					return "got is not *categoryapi.Category"
				}

				expResp := exp.(*categoryapi.Category)
				expResp.DateCreated = gotResp.DateCreated
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
			Name:       "invalid-name",
			URL:        fmt.Sprintf("/v1/categories/%s", sd.Categories[0].ID),
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPut,
			StatusCode: http.StatusBadRequest,
			Input: &categoryapi.UpdateCategory{
				Name: dbtest.StringPointer(""),
			},
			GotResp: &errs.Error{},
			ExpResp: errs.Newf(errs.InvalidArgument, `parse name: invalid name ""`),
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
			URL:        fmt.Sprintf("/v1/categories/%s", sd.Categories[0].ID),
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
			URL:        fmt.Sprintf("/v1/categories/%s", sd.Categories[0].ID),
			Token:      sd.Users[0].Token,
			Method:     http.MethodPut,
			StatusCode: http.StatusUnauthorized,
			GotResp:    &errs.Error{},
			ExpResp:    errs.Newf(errs.Unauthenticated, "authorize: you are not authorized for that action, claims[[USER]] rule[rule_admin_only]"),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
	}

	return table
}
