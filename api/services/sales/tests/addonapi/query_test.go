package addonapi_test

import (
	"fmt"
	"net/http"
	"sort"

	"github.com/google/go-cmp/cmp"
	"github.com/warlck/food-flow/app/domain/addonapp"
	"github.com/warlck/food-flow/app/sdk/apitest"
	"github.com/warlck/food-flow/app/sdk/query"
)

func query200(sd apitest.SeedData) []apitest.Table {
	addons := make([]addonapp.Addon, 0, len(sd.Addons))

	for _, a := range sd.Addons {
		addons = append(addons, toAppAddon(a.Addon))
	}

	sort.Slice(addons, func(i, j int) bool {
		return addons[i].ID <= addons[j].ID
	})

	itemAddons := make([]addonapp.Addon, 0, len(sd.Addons))
	for _, a := range sd.Addons {
		if a.MenuItemID == sd.MenuItems[0].ID {
			itemAddons = append(itemAddons, toAppAddon(a.Addon))
		}
	}

	sort.Slice(itemAddons, func(i, j int) bool {
		return itemAddons[i].ID <= itemAddons[j].ID
	})

	table := []apitest.Table{
		{
			Name:       "basic",
			URL:        "/v1/addons?page=1&rows=10&orderBy=addon_id,ASC",
			Token:      sd.Admins[0].Token,
			StatusCode: http.StatusOK,
			Method:     http.MethodGet,
			GotResp:    &query.Result[addonapp.Addon]{},
			ExpResp: &query.Result[addonapp.Addon]{
				Page:        1,
				RowsPerPage: 10,
				Total:       len(addons),
				Items:       addons,
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
