package addonapi_test

import (
	"fmt"
	"net/http"
	"sort"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/warlck/food-flow/app/domain/addonapp"
	"github.com/warlck/food-flow/app/sdk/apitest"
	"github.com/warlck/food-flow/app/sdk/errs"
	"github.com/warlck/food-flow/app/sdk/query"
)

func query200(sd apitest.SeedData) []apitest.Table {
	itemAddons := make([]addonapp.Addon, 0, len(sd.Addons))
	for _, a := range sd.Addons {
		if a.MenuItemID == sd.MenuItems[0].ID {
			itemAddons = append(itemAddons, toAppAddon(a.Addon))
		}
	}

	sort.Slice(itemAddons, func(i, j int) bool {
		return itemAddons[i].ID <= itemAddons[j].ID
	})

	restAddons := make([]addonapp.Addon, 0, len(sd.Addons))
	for _, a := range sd.Addons {
		if a.RestaurantID == sd.Restaurants[0].ID {
			restAddons = append(restAddons, toAppAddon(a.Addon))
		}
	}

	sort.Slice(restAddons, func(i, j int) bool {
		return restAddons[i].ID <= restAddons[j].ID
	})

	table := []apitest.Table{
		{
			Name:       "by-restaurant",
			URL:        fmt.Sprintf("/v1/addons?page=1&rows=10&orderBy=addon_id,ASC&restaurant_id=%s", sd.Restaurants[0].ID),
			Token:      sd.Admins[0].Token,
			StatusCode: http.StatusOK,
			Method:     http.MethodGet,
			GotResp:    &query.Result[addonapp.Addon]{},
			ExpResp: &query.Result[addonapp.Addon]{
				Page:        1,
				RowsPerPage: 10,
				Total:       len(restAddons),
				Items:       restAddons,
			},
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
		{
			Name:       "by-menu-item",
			URL:        fmt.Sprintf("/v1/addons?page=1&rows=10&orderBy=addon_id,ASC&menu_item_id=%s", sd.MenuItems[0].ID),
			Token:      sd.Admins[0].Token,
			StatusCode: http.StatusOK,
			Method:     http.MethodGet,
			GotResp:    &query.Result[addonapp.Addon]{},
			ExpResp: &query.Result[addonapp.Addon]{
				Page:        1,
				RowsPerPage: 10,
				Total:       len(itemAddons),
				Items:       itemAddons,
			},
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
		{
			Name:       "by-menu-item-empty",
			URL:        fmt.Sprintf("/v1/addons?page=1&rows=10&orderBy=addon_id,ASC&menu_item_id=%s", sd.MenuItems[1].ID),
			Token:      sd.Admins[0].Token,
			StatusCode: http.StatusOK,
			Method:     http.MethodGet,
			GotResp:    &query.Result[addonapp.Addon]{},
			ExpResp: &query.Result[addonapp.Addon]{
				Page:        1,
				RowsPerPage: 10,
				Total:       0,
				Items:       []addonapp.Addon{},
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
			URL:        "/v1/addons?page=1&rows=10",
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
			URL:        fmt.Sprintf("/v1/addons?page=1&rows=10&menu_item_id=%s", sd.MenuItems[2].ID),
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
			URL:        fmt.Sprintf("/v1/addons?page=1&rows=10&restaurant_id=%s", sd.Restaurants[2].ID),
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
			URL:        fmt.Sprintf("/v1/addons/%s", sd.Addons[0].ID),
			Token:      sd.Admins[0].Token,
			StatusCode: http.StatusOK,
			Method:     http.MethodGet,
			GotResp:    &addonapp.Addon{},
			ExpResp:    toAppAddonPtr(sd.Addons[0].Addon),
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
			URL:        fmt.Sprintf("/v1/addons/%s", sd.Addons[4].ID),
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
			URL:        fmt.Sprintf("/v1/addons/%s", unknownID),
			Token:      sd.Admins[0].Token,
			StatusCode: http.StatusNotFound,
			Method:     http.MethodGet,
			GotResp:    &errs.Error{},
			ExpResp: &errs.Error{
				Code:    errs.NotFound,
				Message: fmt.Sprintf("query: addonID[%s]: addon not found", unknownID),
			},
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
	}

	return table
}
