package modifieroptionapi_test

import (
	"fmt"
	"net/http"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/warlck/food-flow/app/domain/modifieroptionapp"
	"github.com/warlck/food-flow/app/sdk/apitest"
	"github.com/warlck/food-flow/app/sdk/errs"
	"github.com/warlck/food-flow/business/types/opt"
)

func update200(sd apitest.SeedData) []apitest.Table {
	disable := false

	exp := toAppModifierOptionPtr(sd.ModifierOptions[1].ModifierOption)
	exp.Available = false

	table := []apitest.Table{
		{
			Name:       "disable-one-of-two-available",
			URL:        fmt.Sprintf("/v1/modifieroptions/%s", sd.ModifierOptions[1].ID),
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPut,
			StatusCode: http.StatusOK,
			Input: &modifieroptionapp.UpdateModifierOption{
				Available: &disable,
				// A zero NullInt marshals as explicit null, so pass the
				// current rank to keep it untouched.
				Rank: opt.NewNullInt(*sd.ModifierOptions[1].Rank),
			},
			GotResp: &modifieroptionapp.ModifierOption{},
			ExpResp: exp,
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.(*modifieroptionapp.ModifierOption)
				if !exists {
					return "got is not *modifieroptionapp.ModifierOption"
				}
				expResp := exp.(*modifieroptionapp.ModifierOption)
				expResp.DateUpdated = gotResp.DateUpdated
				return cmp.Diff(got, exp)
			},
		},
	}

	return table
}

func update400(sd apitest.SeedData) []apitest.Table {
	disable := false

	table := []apitest.Table{
		{
			Name:       "disable-last-available-of-active-required-group",
			URL:        fmt.Sprintf("/v1/modifieroptions/%s", sd.ModifierOptions[3].ID),
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPut,
			StatusCode: http.StatusBadRequest,
			Input: &modifieroptionapp.UpdateModifierOption{
				Available: &disable,
			},
			GotResp: &errs.Error{},
			ExpResp: &errs.Error{
				Code:    errs.InvalidArgument,
				Message: "cannot disable or delete the last available option of an active required group",
			},
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
	}

	return table
}

func update404(sd apitest.SeedData) []apitest.Table {
	unknownID := uuid.New()

	table := []apitest.Table{
		{
			Name:       "not-found",
			URL:        fmt.Sprintf("/v1/modifieroptions/%s", unknownID),
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPut,
			StatusCode: http.StatusNotFound,
			Input:      &modifieroptionapp.UpdateModifierOption{},
			GotResp:    &errs.Error{},
			ExpResp: &errs.Error{
				Code:    errs.NotFound,
				Message: fmt.Sprintf("query: optionID[%s]: db: modifier option not found", unknownID),
			},
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
	}

	return table
}
