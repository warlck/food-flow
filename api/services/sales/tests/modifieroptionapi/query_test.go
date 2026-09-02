package modifieroptionapi_test

import (
	"fmt"
	"net/http"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/warlck/food-flow/app/domain/modifieroptionapp"
	"github.com/warlck/food-flow/app/sdk/apitest"
	"github.com/warlck/food-flow/app/sdk/errs"
)

func queryByID200(sd apitest.SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:       "basic",
			URL:        fmt.Sprintf("/v1/modifieroptions/%s", sd.ModifierOptions[0].ID),
			Token:      sd.Admins[0].Token,
			StatusCode: http.StatusOK,
			Method:     http.MethodGet,
			GotResp:    &modifieroptionapp.ModifierOption{},
			ExpResp:    toAppModifierOptionPtr(sd.ModifierOptions[0].ModifierOption),
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
			URL:        fmt.Sprintf("/v1/modifieroptions/%s", unknownID),
			Token:      sd.Admins[0].Token,
			StatusCode: http.StatusNotFound,
			Method:     http.MethodGet,
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
