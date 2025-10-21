package categoryapi_test

import (
	"fmt"
	"net/http"
	"sort"

	"github.com/google/go-cmp/cmp"
	categoryapi "github.com/warlck/food-flow/app/domain/categoryapp"
	"github.com/warlck/food-flow/app/sdk/apitest"
	"github.com/warlck/food-flow/app/sdk/query"
)

func query200(sd apitest.SeedData) []apitest.Table {
	cats := make([]categoryapi.Category, 0, len(sd.Categories))

	for _, cat := range sd.Categories {
		cats = append(cats, *toAppCategoryPtr(cat.Category))
	}

	sort.Slice(cats, func(i, j int) bool {
		return cats[i].ID <= cats[j].ID
	})

	table := []apitest.Table{
		{
			Name:       "basic",
			URL:        "/v1/categories?page=1&rows=10&orderBy=category_id,ASC",
			Token:      sd.Admins[0].Token,
			StatusCode: http.StatusOK,
			Method:     http.MethodGet,
			GotResp:    &query.Result[categoryapi.Category]{},
			ExpResp: &query.Result[categoryapi.Category]{
				Page:        1,
				RowsPerPage: 10,
				Total:       len(cats),
				Items:       cats,
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
			URL:        fmt.Sprintf("/v1/categories/%s", sd.Categories[0].ID),
			Token:      sd.Admins[0].Token,
			StatusCode: http.StatusOK,
			Method:     http.MethodGet,
			GotResp:    &categoryapi.Category{},
			ExpResp:    toAppCategoryPtr(sd.Categories[0].Category),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
	}

	return table
}
