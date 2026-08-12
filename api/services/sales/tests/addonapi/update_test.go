package addonapi_test

import (
	"fmt"
	"net/http"

	"github.com/google/go-cmp/cmp"
	"github.com/warlck/food-flow/app/domain/addonapp"
	"github.com/warlck/food-flow/app/sdk/apitest"
	"github.com/warlck/food-flow/app/sdk/errs"
	"github.com/warlck/food-flow/business/sdk/dbtest"
)

func update200(sd apitest.SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:       "basic",
			URL:        fmt.Sprintf("/v1/addons/%s", sd.Addons[0].ID),
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPut,
			StatusCode: http.StatusOK,
			Input: &addonapp.UpdateAddon{
				Name:        dbtest.StringPointer("Updated Extra Cheese"),
				Price:       dbtest.FloatPointer(3.00),
				Description: dbtest.StringPointer("Extra melted cheese portion"),
			},

			GotResp: &addonapp.Addon{},
			ExpResp: &addonapp.Addon{
				ID:           sd.Addons[0].ID.String(),
				CategoryID:   sd.Addons[0].CategoryID.String(),
				RestaurantID: sd.Addons[0].RestaurantID.String(),
				Name:         "Updated Extra Cheese",
				Description:  "Extra melted cheese portion",
				Price:        3.00,
				Available:    sd.Addons[0].Available,
				MaxQuantity:  sd.Addons[0].MaxQuantity,
				DateCreated:  sd.Addons[0].DateCreated.Format("2006-01-02T15:04:05Z07:00"),
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.(*addonapp.Addon)
				if !exists {
					return "got is not *addonapp.Addon"
				}

				expResp := exp.(*addonapp.Addon)
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
			URL:        fmt.Sprintf("/v1/addons/%s", sd.Addons[0].ID),
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPut,
			StatusCode: http.StatusBadRequest,
			Input: &addonapp.UpdateAddon{
				Name: dbtest.StringPointer(""),
			},
			GotResp: &errs.Error{},
			ExpResp: &errs.Error{
				Code:    errs.InvalidArgument,
				Message: `parse name: invalid name ""`,
			},
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
			URL:        fmt.Sprintf("/v1/addons/%s", sd.Addons[0].ID),
			Token:      "",
			Method:     http.MethodPut,
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
	}

	return table
}
