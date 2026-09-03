package modifieroptionapi_test

import (
	"fmt"
	"net/http"

	"github.com/google/go-cmp/cmp"
	"github.com/warlck/food-flow/app/domain/modifieroptionapp"
	"github.com/warlck/food-flow/app/sdk/apitest"
	"github.com/warlck/food-flow/app/sdk/errs"
)

func create201(sd apitest.SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:       "basic-delta",
			URL:        "/v1/modifieroptions",
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPost,
			StatusCode: http.StatusCreated,
			Input: &modifieroptionapp.NewModifierOption{
				ModifierGroupID: sd.ModifierGroups[3].ID.String(),
				RestaurantID:    sd.Restaurants[0].ID.String(),
				Name:            "Gluten Free Crust",
				Description:     "Specialty crust",
				PriceDelta:      2.50,
			},
			GotResp: &modifieroptionapp.ModifierOption{},
			ExpResp: &modifieroptionapp.ModifierOption{
				ModifierGroupID: sd.ModifierGroups[3].ID.String(),
				RestaurantID:    sd.Restaurants[0].ID.String(),
				Name:            "Gluten Free Crust",
				Description:     "Specialty crust",
				PriceDelta:      2.50,
				Available:       true,
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.(*modifieroptionapp.ModifierOption)
				if !exists {
					return "got is not *modifieroptionapp.ModifierOption"
				}

				expResp := exp.(*modifieroptionapp.ModifierOption)
				expResp.ID = gotResp.ID
				expResp.DateCreated = gotResp.DateCreated
				expResp.DateUpdated = gotResp.DateUpdated

				return cmp.Diff(got, exp)
			},
		},
		{
			Name:       "zero-delta",
			URL:        "/v1/modifieroptions",
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPost,
			StatusCode: http.StatusCreated,
			Input: &modifieroptionapp.NewModifierOption{
				ModifierGroupID: sd.ModifierGroups[3].ID.String(),
				RestaurantID:    sd.Restaurants[0].ID.String(),
				Name:            "Regular Crust",
				Description:     "Standard dough",
				PriceDelta:      0.00,
			},
			GotResp: &modifieroptionapp.ModifierOption{},
			ExpResp: &modifieroptionapp.ModifierOption{
				ModifierGroupID: sd.ModifierGroups[3].ID.String(),
				RestaurantID:    sd.Restaurants[0].ID.String(),
				Name:            "Regular Crust",
				Description:     "Standard dough",
				PriceDelta:      0.00,
				Available:       true,
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.(*modifieroptionapp.ModifierOption)
				if !exists {
					return "got is not *modifieroptionapp.ModifierOption"
				}

				expResp := exp.(*modifieroptionapp.ModifierOption)
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
			Name:       "negative-delta",
			URL:        "/v1/modifieroptions",
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPost,
			StatusCode: http.StatusBadRequest,
			Input: &modifieroptionapp.NewModifierOption{
				ModifierGroupID: sd.ModifierGroups[3].ID.String(),
				RestaurantID:    sd.Restaurants[0].ID.String(),
				Name:            "Discount Crust",
				PriceDelta:      -1.00,
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
			URL:        "/v1/modifieroptions",
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
			URL:        "/v1/modifieroptions",
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPost,
			StatusCode: http.StatusForbidden,
			Input: &modifieroptionapp.NewModifierOption{
				ModifierGroupID: sd.ModifierGroups[1].ID.String(),
				RestaurantID:    sd.Restaurants[2].ID.String(),
				Name:            "Cross Org Option",
				PriceDelta:      1.00,
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
