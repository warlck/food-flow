package categoryapi_test

import (
	"fmt"
	"net/http"

	"github.com/warlck/food-flow/app/domain/categoryapp"
	"github.com/warlck/food-flow/app/sdk/apitest"
	"github.com/warlck/food-flow/app/sdk/errs"
)

func reorder200(sd apitest.SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:       "swap-categories",
			URL:        "/v1/categories/order",
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPut,
			StatusCode: http.StatusOK,
			Input: &categoryapp.ReorderCategories{
				RestaurantID: sd.Restaurants[1].ID.String(),
				OrderedIDs: []string{
					sd.Categories[3].ID.String(),
					sd.Categories[2].ID.String(),
				},
			},
			GotResp: &[]categoryapp.Category{},
			ExpResp: &[]categoryapp.Category{},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.(*[]categoryapp.Category)
				if !exists {
					return "got is not *[]categoryapp.Category"
				}

				cats := *gotResp
				if len(cats) != 2 {
					return fmt.Sprintf("expected 2 categories, got %d", len(cats))
				}

				if cats[0].ID != sd.Categories[3].ID.String() {
					return fmt.Sprintf("pos 0: expected %s, got %s", sd.Categories[3].ID, cats[0].ID)
				}
				if cats[1].ID != sd.Categories[2].ID.String() {
					return fmt.Sprintf("pos 1: expected %s, got %s", sd.Categories[2].ID, cats[1].ID)
				}

				if cats[0].Rank == nil || *cats[0].Rank != 10 {
					return "pos 0 rank is not 10"
				}
				if cats[1].Rank == nil || *cats[1].Rank != 20 {
					return "pos 1 rank is not 20"
				}

				return ""
			},
		},
	}

	return table
}

func reorder400(sd apitest.SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:       "subset-mismatch",
			URL:        "/v1/categories/order",
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPut,
			StatusCode: http.StatusBadRequest,
			Input: &categoryapp.ReorderCategories{
				RestaurantID: sd.Restaurants[1].ID.String(),
				OrderedIDs: []string{
					sd.Categories[2].ID.String(),
				},
			},
			GotResp: &errs.Error{},
			ExpResp: &errs.Error{
				Code:    errs.InvalidArgument,
				Message: "invalid reorder set: exact set mismatch: expected 2 categories, got 1",
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
			Name:       "duplicate-id",
			URL:        "/v1/categories/order",
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPut,
			StatusCode: http.StatusBadRequest,
			Input: &categoryapp.ReorderCategories{
				RestaurantID: sd.Restaurants[0].ID.String(),
				OrderedIDs: []string{
					sd.Categories[0].ID.String(),
					sd.Categories[0].ID.String(),
				},
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

func reorder401(sd apitest.SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:       "empty-token",
			URL:        "/v1/categories/order",
			Token:      "",
			Method:     http.MethodPut,
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
		{
			Name:       "user-not-admin",
			URL:        "/v1/categories/order",
			Token:      sd.Users[0].Token,
			Method:     http.MethodPut,
			StatusCode: http.StatusForbidden,
			GotResp:    &errs.Error{},
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
