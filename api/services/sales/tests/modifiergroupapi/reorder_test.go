package modifiergroupapi_test

import (
	"fmt"
	"net/http"

	"github.com/warlck/food-flow/app/domain/modifiergroupapp"
	"github.com/warlck/food-flow/app/sdk/apitest"
	"github.com/warlck/food-flow/app/sdk/errs"
)

func reorder200(sd apitest.SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:       "swap-groups",
			URL:        "/v1/modifiergroups/order",
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPut,
			StatusCode: http.StatusOK,
			Input: &modifiergroupapp.ReorderModifierGroups{
				MenuItemID: sd.MenuItems[0].ID.String(),
				OrderedIDs: []string{
					sd.ModifierGroups[1].ID.String(),
					sd.ModifierGroups[0].ID.String(),
				},
			},
			GotResp: &[]modifiergroupapp.ModifierGroup{},
			ExpResp: &[]modifiergroupapp.ModifierGroup{},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.(*[]modifiergroupapp.ModifierGroup)
				if !exists {
					return "got is not *[]modifiergroupapp.ModifierGroup"
				}

				groups := *gotResp
				if len(groups) != 2 {
					return fmt.Sprintf("expected 2 groups, got %d", len(groups))
				}

				if groups[0].ID != sd.ModifierGroups[1].ID.String() {
					return fmt.Sprintf("pos 0: expected %s, got %s", sd.ModifierGroups[1].ID, groups[0].ID)
				}
				if groups[1].ID != sd.ModifierGroups[0].ID.String() {
					return fmt.Sprintf("pos 1: expected %s, got %s", sd.ModifierGroups[0].ID, groups[1].ID)
				}

				if groups[0].Rank == nil || *groups[0].Rank != 10 {
					return "pos 0 rank is not 10"
				}
				if groups[1].Rank == nil || *groups[1].Rank != 20 {
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
			URL:        "/v1/modifiergroups/order",
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPut,
			StatusCode: http.StatusBadRequest,
			Input: &modifiergroupapp.ReorderModifierGroups{
				MenuItemID: sd.MenuItems[0].ID.String(),
				OrderedIDs: []string{
					sd.ModifierGroups[0].ID.String(),
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
			URL:        "/v1/modifiergroups/order",
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
	}

	return table
}
