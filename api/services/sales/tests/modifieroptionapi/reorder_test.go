package modifieroptionapi_test

import (
	"fmt"
	"net/http"

	"github.com/warlck/food-flow/app/domain/modifieroptionapp"
	"github.com/warlck/food-flow/app/sdk/apitest"
	"github.com/warlck/food-flow/app/sdk/errs"
)

func reorder200(sd apitest.SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:       "swap-options",
			URL:        "/v1/modifieroptions/order",
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPut,
			StatusCode: http.StatusOK,
			Input: &modifieroptionapp.ReorderModifierOptions{
				ModifierGroupID: sd.ModifierGroups[0].ID.String(),
				OrderedIDs: []string{
					sd.ModifierOptions[1].ID.String(),
					sd.ModifierOptions[0].ID.String(),
				},
			},
			GotResp: &[]modifieroptionapp.ModifierOption{},
			ExpResp: &[]modifieroptionapp.ModifierOption{},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.(*[]modifieroptionapp.ModifierOption)
				if !exists {
					return "got is not *[]modifieroptionapp.ModifierOption"
				}

				opts := *gotResp
				if len(opts) != 2 {
					return fmt.Sprintf("expected 2 options, got %d", len(opts))
				}

				if opts[0].ID != sd.ModifierOptions[1].ID.String() {
					return fmt.Sprintf("pos 0: expected %s, got %s", sd.ModifierOptions[1].ID, opts[0].ID)
				}
				if opts[1].ID != sd.ModifierOptions[0].ID.String() {
					return fmt.Sprintf("pos 1: expected %s, got %s", sd.ModifierOptions[0].ID, opts[1].ID)
				}

				if opts[0].Rank == nil || *opts[0].Rank != 10 {
					return "pos 0 rank is not 10"
				}
				if opts[1].Rank == nil || *opts[1].Rank != 20 {
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
			URL:        "/v1/modifieroptions/order",
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPut,
			StatusCode: http.StatusBadRequest,
			Input: &modifieroptionapp.ReorderModifierOptions{
				ModifierGroupID: sd.ModifierGroups[0].ID.String(),
				OrderedIDs: []string{
					sd.ModifierOptions[0].ID.String(),
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
			URL:        "/v1/modifieroptions/order",
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
