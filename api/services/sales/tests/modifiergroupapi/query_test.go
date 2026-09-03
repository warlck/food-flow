package modifiergroupapi_test

import (
	"fmt"
	"net/http"
	"sort"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/warlck/food-flow/app/domain/modifiergroupapp"
	"github.com/warlck/food-flow/app/sdk/apitest"
	"github.com/warlck/food-flow/app/sdk/errs"
	"github.com/warlck/food-flow/app/sdk/query"
)

func query200(sd apitest.SeedData) []apitest.Table {
	groups := make([]modifiergroupapp.ModifierGroup, 0, 2)
	for _, g := range sd.ModifierGroups[:2] {
		groups = append(groups, toAppModifierGroup(g.ModifierGroup))
	}

	sort.Slice(groups, func(i, j int) bool {
		return groups[i].ID <= groups[j].ID
	})

	table := []apitest.Table{
		{
			Name:       "by-menu-item",
			URL:        fmt.Sprintf("/v1/modifiergroups?page=1&rows=10&orderBy=id,ASC&menu_item_id=%s", sd.MenuItems[0].ID),
			Token:      sd.Admins[0].Token,
			StatusCode: http.StatusOK,
			Method:     http.MethodGet,
			GotResp:    &query.Result[modifiergroupapp.ModifierGroup]{},
			ExpResp: &query.Result[modifiergroupapp.ModifierGroup]{
				Page:        1,
				RowsPerPage: 10,
				Total:       2,
				Items:       groups,
			},
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
		{
			Name:       "by-restaurant",
			URL:        fmt.Sprintf("/v1/modifiergroups?page=1&rows=10&orderBy=id,ASC&restaurant_id=%s", sd.Restaurants[0].ID),
			Token:      sd.Admins[0].Token,
			StatusCode: http.StatusOK,
			Method:     http.MethodGet,
			GotResp:    &query.Result[modifiergroupapp.ModifierGroup]{},
			ExpResp: &query.Result[modifiergroupapp.ModifierGroup]{
				Page:        1,
				RowsPerPage: 10,
				Total:       2,
				Items:       groups,
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
			URL:        "/v1/modifiergroups?page=1&rows=10",
			Token:      sd.Admins[0].Token,
			StatusCode: http.StatusBadRequest,
			Method:     http.MethodGet,
			GotResp:    &errs.Error{},
			ExpResp: &errs.Error{
				Code:    errs.InvalidArgument,
				Message: "query requires a restaurant_id or menu_item_id filter",
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
			Name:       "cross-org-menu-item",
			URL:        fmt.Sprintf("/v1/modifiergroups?page=1&rows=10&menu_item_id=%s", sd.MenuItems[2].ID),
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
		{
			Name:       "cross-org-restaurant",
			URL:        fmt.Sprintf("/v1/modifiergroups?page=1&rows=10&restaurant_id=%s", sd.Restaurants[2].ID),
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

func queryByID403(sd apitest.SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:       "cross-org",
			URL:        fmt.Sprintf("/v1/modifiergroups/%s", sd.ModifierGroups[2].ID),
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
