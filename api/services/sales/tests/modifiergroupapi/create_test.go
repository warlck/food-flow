package modifiergroupapi_test

import (
	"fmt"
	"net/http"

	"github.com/google/go-cmp/cmp"
	"github.com/warlck/food-flow/app/domain/modifiergroupapp"
	"github.com/warlck/food-flow/app/sdk/apitest"
	"github.com/warlck/food-flow/app/sdk/errs"
)

func create201(sd apitest.SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:       "optional-group",
			URL:        "/v1/modifiergroups",
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPost,
			StatusCode: http.StatusCreated,
			Input: &modifiergroupapp.NewModifierGroup{
				MenuItemID:    sd.MenuItems[1].ID.String(),
				RestaurantID:  sd.Restaurants[0].ID.String(),
				Name:          "Choice of Dressing",
				Description:   "Optional salad dressing",
				MinSelections: 0,
				MaxSelections: 1,
				Available:     true,
			},
			GotResp: &modifiergroupapp.ModifierGroup{},
			ExpResp: &modifiergroupapp.ModifierGroup{
				MenuItemID:    sd.MenuItems[1].ID.String(),
				RestaurantID:  sd.Restaurants[0].ID.String(),
				Name:          "Choice of Dressing",
				Description:   "Optional salad dressing",
				MinSelections: 0,
				MaxSelections: 1,
				Available:     true,
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.(*modifiergroupapp.ModifierGroup)
				if !exists {
					return "got is not *modifiergroupapp.ModifierGroup"
				}

				expResp := exp.(*modifiergroupapp.ModifierGroup)
				expResp.ID = gotResp.ID
				expResp.DateCreated = gotResp.DateCreated
				expResp.DateUpdated = gotResp.DateUpdated

				return cmp.Diff(got, exp)
			},
		},
		{
			Name:       "required-group-disabled",
			URL:        "/v1/modifiergroups",
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPost,
			StatusCode: http.StatusCreated,
			Input: &modifiergroupapp.NewModifierGroup{
				MenuItemID:    sd.MenuItems[1].ID.String(),
				RestaurantID:  sd.Restaurants[0].ID.String(),
				Name:          "Required Crust Type",
				Description:   "Crust style",
				MinSelections: 1,
				MaxSelections: 1,
				Available:     false,
			},
			GotResp: &modifiergroupapp.ModifierGroup{},
			ExpResp: &modifiergroupapp.ModifierGroup{
				MenuItemID:    sd.MenuItems[1].ID.String(),
				RestaurantID:  sd.Restaurants[0].ID.String(),
				Name:          "Required Crust Type",
				Description:   "Crust style",
				MinSelections: 1,
				MaxSelections: 1,
				Available:     false,
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.(*modifiergroupapp.ModifierGroup)
				if !exists {
					return "got is not *modifiergroupapp.ModifierGroup"
				}

				expResp := exp.(*modifiergroupapp.ModifierGroup)
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
			Name:       "active-required-rejected",
			URL:        "/v1/modifiergroups",
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPost,
			StatusCode: http.StatusBadRequest,
			Input: &modifiergroupapp.NewModifierGroup{
				MenuItemID:    sd.MenuItems[1].ID.String(),
				RestaurantID:  sd.Restaurants[0].ID.String(),
				Name:          "Invalid Active Required",
				MinSelections: 1,
				MaxSelections: 1,
				Available:     true,
			},
			GotResp: &errs.Error{},
			ExpResp: &errs.Error{
				Code: errs.InvalidArgument,
			},
			CmpFunc: func(got any, exp any) string {
				gotErr := got.(*errs.Error)
				expErr := exp.(*errs.Error)
				if gotErr.Code != expErr.Code {
					return fmt.Sprintf("code mismatch: got %v, exp %v", gotErr.Code, expErr.Code)
				}
				return ""
			},
		},
		{
			Name:       "invalid-name",
			URL:        "/v1/modifiergroups",
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPost,
			StatusCode: http.StatusBadRequest,
			Input: &modifiergroupapp.NewModifierGroup{
				MenuItemID:    sd.MenuItems[1].ID.String(),
				RestaurantID:  sd.Restaurants[0].ID.String(),
				Name:          "",
				MinSelections: 0,
				MaxSelections: 1,
			},
			GotResp: &errs.Error{},
			ExpResp: &errs.Error{
				Code: errs.InvalidArgument,
			},
			CmpFunc: func(got any, exp any) string {
				gotErr := got.(*errs.Error)
				expErr := exp.(*errs.Error)
				if gotErr.Code != expErr.Code {
					return fmt.Sprintf("code mismatch: got %v, exp %v", gotErr.Code, expErr.Code)
				}
				return ""
			},
		},
	}

	return table
}

func create401(sd apitest.SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:       "empty-token",
			URL:        "/v1/modifiergroups",
			Token:      "",
			Method:     http.MethodPost,
			StatusCode: http.StatusUnauthorized,
			GotResp:    &errs.Error{},
			ExpResp: &errs.Error{
				Code: errs.Unauthenticated,
			},
			CmpFunc: func(got any, exp any) string {
				gotErr := got.(*errs.Error)
				expErr := exp.(*errs.Error)
				if gotErr.Code != expErr.Code {
					return fmt.Sprintf("code mismatch: got %v, exp %v", gotErr.Code, expErr.Code)
				}
				return ""
			},
		},
	}

	return table
}

func create403(sd apitest.SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:       "cross-org",
			URL:        "/v1/modifiergroups",
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPost,
			StatusCode: http.StatusForbidden,
			Input: &modifiergroupapp.NewModifierGroup{
				MenuItemID:    sd.MenuItems[2].ID.String(),
				RestaurantID:  sd.Restaurants[2].ID.String(),
				Name:          "Cross Org Group",
				MinSelections: 0,
				MaxSelections: 1,
			},
			GotResp: &errs.Error{},
			ExpResp: &errs.Error{
				Code: errs.PermissionDenied,
			},
			CmpFunc: func(got any, exp any) string {
				gotErr := got.(*errs.Error)
				expErr := exp.(*errs.Error)
				if gotErr.Code != expErr.Code {
					return fmt.Sprintf("code mismatch: got %v, exp %v", gotErr.Code, expErr.Code)
				}
				return ""
			},
		},
	}

	return table
}
