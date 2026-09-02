package modifiergroupapi_test

import (
	"fmt"
	"net/http"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/warlck/food-flow/app/domain/modifiergroupapp"
	"github.com/warlck/food-flow/app/sdk/apitest"
	"github.com/warlck/food-flow/app/sdk/errs"
	"github.com/warlck/food-flow/business/types/opt"
)

// update200 exercises the disable/enable cycle on a required group that has
// available options: disabling is always allowed, enabling requires at least
// one available option.
func update200(sd apitest.SeedData) []apitest.Table {
	disable := false
	enable := true

	// A zero NullInt marshals as explicit null, so pass the current rank to
	// keep it untouched.
	rank := opt.NewNullInt(*sd.ModifierGroups[1].Rank)

	groupResp := func(available bool) *modifiergroupapp.ModifierGroup {
		exp := toAppModifierGroupPtr(sd.ModifierGroups[1].ModifierGroup)
		exp.Available = available
		return exp
	}

	cmpUpdatedGroup := func(got any, exp any) string {
		gotResp, exists := got.(*modifiergroupapp.ModifierGroup)
		if !exists {
			return "got is not *modifiergroupapp.ModifierGroup"
		}
		expResp := exp.(*modifiergroupapp.ModifierGroup)
		expResp.DateUpdated = gotResp.DateUpdated
		return cmp.Diff(got, exp)
	}

	table := []apitest.Table{
		{
			Name:       "disable-required-with-options",
			URL:        fmt.Sprintf("/v1/modifiergroups/%s", sd.ModifierGroups[1].ID),
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPut,
			StatusCode: http.StatusOK,
			Input: &modifiergroupapp.UpdateModifierGroup{
				Available: &disable,
				Rank:      rank,
			},
			GotResp: &modifiergroupapp.ModifierGroup{},
			ExpResp: groupResp(false),
			CmpFunc: cmpUpdatedGroup,
		},
		{
			Name:       "enable-required-with-options",
			URL:        fmt.Sprintf("/v1/modifiergroups/%s", sd.ModifierGroups[1].ID),
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPut,
			StatusCode: http.StatusOK,
			Input: &modifiergroupapp.UpdateModifierGroup{
				Available: &enable,
				Rank:      rank,
			},
			GotResp: &modifiergroupapp.ModifierGroup{},
			ExpResp: groupResp(true),
			CmpFunc: cmpUpdatedGroup,
		},
	}

	return table
}

func update400(sd apitest.SeedData) []apitest.Table {
	enable := true

	table := []apitest.Table{
		{
			Name:       "enable-required-without-options",
			URL:        fmt.Sprintf("/v1/modifiergroups/%s", sd.ModifierGroups[0].ID),
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPut,
			StatusCode: http.StatusBadRequest,
			Input: &modifiergroupapp.UpdateModifierGroup{
				Available: &enable,
			},
			GotResp: &errs.Error{},
			ExpResp: &errs.Error{
				Code:    errs.InvalidArgument,
				Message: "required modifier group must have at least one available option",
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
			URL:        fmt.Sprintf("/v1/modifiergroups/%s", unknownID),
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPut,
			StatusCode: http.StatusNotFound,
			Input:      &modifiergroupapp.UpdateModifierGroup{},
			GotResp:    &errs.Error{},
			ExpResp: &errs.Error{
				Code:    errs.NotFound,
				Message: fmt.Sprintf("query: groupID[%s]: db: modifier group not found", unknownID),
			},
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
	}

	return table
}
