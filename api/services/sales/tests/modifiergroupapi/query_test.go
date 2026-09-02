package modifiergroupapi_test

import (
	"fmt"
	"net/http"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/warlck/food-flow/app/domain/modifiergroupapp"
	"github.com/warlck/food-flow/app/sdk/apitest"
	"github.com/warlck/food-flow/app/sdk/errs"
)

func queryByID200(sd apitest.SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:       "basic",
			URL:        fmt.Sprintf("/v1/modifiergroups/%s", sd.ModifierGroups[0].ID),
			Token:      sd.Admins[0].Token,
			StatusCode: http.StatusOK,
			Method:     http.MethodGet,
			GotResp:    &modifiergroupapp.ModifierGroup{},
			ExpResp:    toAppModifierGroupPtr(sd.ModifierGroups[0].ModifierGroup),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
	}

	return table
}

func queryByID404(sd apitest.SeedData) []apitest.Table {
	unknownID := uuid.New()

	table := []apitest.Table{
		{
			Name:       "not-found",
			URL:        fmt.Sprintf("/v1/modifiergroups/%s", unknownID),
			Token:      sd.Admins[0].Token,
			StatusCode: http.StatusNotFound,
			Method:     http.MethodGet,
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
