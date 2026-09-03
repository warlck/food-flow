package addonapi_test

import (
	"fmt"
	"net/http"

	"github.com/google/go-cmp/cmp"
	"github.com/warlck/food-flow/app/domain/addonapp"
	"github.com/warlck/food-flow/app/sdk/apitest"
	"github.com/warlck/food-flow/app/sdk/errs"
	"github.com/warlck/food-flow/business/sdk/dbtest"
	"github.com/warlck/food-flow/business/types/opt"
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
				Rank:        opt.NewNullInt(*sd.Addons[0].Rank),
			},

			GotResp: &addonapp.Addon{},
			ExpResp: &addonapp.Addon{
				ID:           sd.Addons[0].ID.String(),
				MenuItemID:   sd.Addons[0].MenuItemID.String(),
				RestaurantID: sd.Addons[0].RestaurantID.String(),
				Name:         "Updated Extra Cheese",
				Description:  "Extra melted cheese portion",
				Price:        3.00,
				Available:    sd.Addons[0].Available,
				MaxQuantity:  sd.Addons[0].MaxQuantity,
				Rank:         sd.Addons[0].Rank,
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
		{
			Name:       "rank-set",
			URL:        fmt.Sprintf("/v1/addons/%s", sd.Addons[0].ID),
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPut,
			StatusCode: http.StatusOK,
			Input: &addonapp.UpdateAddon{
				Rank: opt.NewNullInt(25),
			},
			GotResp: &addonapp.Addon{},
			ExpResp: &addonapp.Addon{
				ID:           sd.Addons[0].ID.String(),
				MenuItemID:   sd.Addons[0].MenuItemID.String(),
				RestaurantID: sd.Addons[0].RestaurantID.String(),
				Name:         "Updated Extra Cheese",
				Description:  "Extra melted cheese portion",
				Price:        3.00,
				Available:    sd.Addons[0].Available,
				MaxQuantity:  sd.Addons[0].MaxQuantity,
				Rank:         dbtest.IntPointer(25),
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
		{
			Name:       "rank-clear",
			URL:        fmt.Sprintf("/v1/addons/%s", sd.Addons[0].ID),
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPut,
			StatusCode: http.StatusOK,
			Input: &addonapp.UpdateAddon{
				Rank: opt.NullInt{Present: true, Value: nil},
			},
			GotResp: &addonapp.Addon{},
			ExpResp: &addonapp.Addon{
				ID:           sd.Addons[0].ID.String(),
				MenuItemID:   sd.Addons[0].MenuItemID.String(),
				RestaurantID: sd.Addons[0].RestaurantID.String(),
				Name:         "Updated Extra Cheese",
				Description:  "Extra melted cheese portion",
				Price:        3.00,
				Available:    sd.Addons[0].Available,
				MaxQuantity:  sd.Addons[0].MaxQuantity,
				Rank:         nil,
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
		{
			Name:       "invalid-rank",
			URL:        fmt.Sprintf("/v1/addons/%s", sd.Addons[0].ID),
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPut,
			StatusCode: http.StatusBadRequest,
			Input: &addonapp.UpdateAddon{
				Rank: opt.NewNullInt(0),
			},
			GotResp: &errs.Error{},
			ExpResp: &errs.Error{
				Code:    errs.InvalidArgument,
				Message: `[{"field":"rank","error":"rank must be \u003e= 1"}]`,
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
