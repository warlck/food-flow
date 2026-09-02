package modifieroptionapi_test

import (
	"fmt"
	"net/http"
	"sort"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/warlck/food-flow/app/domain/modifieroptionapp"
	"github.com/warlck/food-flow/app/sdk/apitest"
	"github.com/warlck/food-flow/app/sdk/errs"
	"github.com/warlck/food-flow/app/sdk/query"
)

func query200(sd apitest.SeedData) []apitest.Table {
	options := make([]modifieroptionapp.ModifierOption, 0, 2)
	for _, o := range sd.ModifierOptions[:2] {
		options = append(options, toAppModifierOption(o.ModifierOption))
	}

	sort.Slice(options, func(i, j int) bool {
		return options[i].ID <= options[j].ID
	})

	table := []apitest.Table{
		{
			Name:       "by-modifier-group",
			URL:        fmt.Sprintf("/v1/modifieroptions?page=1&rows=10&orderBy=id,ASC&modifier_group_id=%s", sd.ModifierGroups[0].ID),
			Token:      sd.Admins[0].Token,
			StatusCode: http.StatusOK,
			Method:     http.MethodGet,
			GotResp:    &query.Result[modifieroptionapp.ModifierOption]{},
			ExpResp: &query.Result[modifieroptionapp.ModifierOption]{
				Page:        1,
				RowsPerPage: 10,
				Total:       2,
				Items:       options,
			},
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
	}

	return table
}

func query400(sd apitest.SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:       "unscoped",
			URL:        "/v1/modifieroptions?page=1&rows=10",
			Token:      sd.Admins[0].Token,
			StatusCode: http.StatusBadRequest,
			Method:     http.MethodGet,
			GotResp:    &errs.Error{},
			ExpResp: &errs.Error{
				Code:    errs.InvalidArgument,
				Message: "query requires a restaurant_id or modifier_group_id filter",
			},
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
	}

	return table
}

func query403(sd apitest.SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:       "cross-org-modifier-group",
			URL:        fmt.Sprintf("/v1/modifieroptions?page=1&rows=10&modifier_group_id=%s", sd.ModifierGroups[1].ID),
			Token:      sd.Admins[0].Token,
			StatusCode: http.StatusForbidden,
			Method:     http.MethodGet,
			GotResp:    &errs.Error{},
			ExpResp: &errs.Error{
				Code:    errs.PermissionDenied,
				Message: fmt.Sprintf("user not in organization %s", sd.Organizations[1].ID),
			},
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
	}

	return table
}

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

func queryByID403(sd apitest.SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:       "cross-org",
			URL:        fmt.Sprintf("/v1/modifieroptions/%s", sd.ModifierOptions[2].ID),
			Token:      sd.Admins[0].Token,
			StatusCode: http.StatusForbidden,
			Method:     http.MethodGet,
			GotResp:    &errs.Error{},
			ExpResp: &errs.Error{
				Code:    errs.PermissionDenied,
				Message: fmt.Sprintf("user not in organization %s", sd.Organizations[1].ID),
			},
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
