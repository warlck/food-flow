package orderapi_test

import (
	"fmt"
	"net/http"
	"sort"

	"github.com/google/go-cmp/cmp"
	"github.com/warlck/food-flow/app/domain/orderapp"
	"github.com/warlck/food-flow/app/sdk/apitest"
	"github.com/warlck/food-flow/app/sdk/query"
	"github.com/warlck/food-flow/business/domain/orderbus"
)

func query200(sd apitest.SeedData) []apitest.Table {
	orders := make([]orderbus.Order, 0, len(sd.Orders))

	for _, ord := range sd.Orders {
		orders = append(orders, ord.Order)
	}

	sort.Slice(orders, func(i, j int) bool {
		return orders[i].DateCreated.After(orders[j].DateCreated)
	})

	table := []apitest.Table{
		{
			Name:       "basic",
			URL:        "/v1/orders?page=1&rows=10&orderBy=date,DESC",
			Token:      sd.Admins[0].Token,
			StatusCode: http.StatusOK,
			Method:     http.MethodGet,
			GotResp:    &query.Result[orderapp.Order]{},
			ExpResp: &query.Result[orderapp.Order]{
				Page:        1,
				RowsPerPage: 10,
				Total:       len(orders),
				Items:       toAppOrders(orders),
			},
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
		{
			Name:       "by-restaurant",
			URL:        fmt.Sprintf("/v1/orders?page=1&rows=10&restaurant_id=%s", sd.Restaurants[0].ID),
			Token:      sd.Admins[0].Token,
			StatusCode: http.StatusOK,
			Method:     http.MethodGet,
			GotResp:    &query.Result[orderapp.Order]{},
			ExpResp: &query.Result[orderapp.Order]{
				Page:        1,
				RowsPerPage: 10,
				Total:       len(orders),
				Items:       toAppOrders(orders),
			},
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
		{
			Name:       "by-status",
			URL:        "/v1/orders?page=1&rows=10&order_status=pending",
			Token:      sd.Admins[0].Token,
			StatusCode: http.StatusOK,
			Method:     http.MethodGet,
			GotResp:    &query.Result[orderapp.Order]{},
			ExpResp: &query.Result[orderapp.Order]{
				Page:        1,
				RowsPerPage: 10,
				Total:       len(orders),
				Items:       toAppOrders(orders),
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
			URL:        fmt.Sprintf("/v1/orders/%s", sd.Orders[0].ID),
			Token:      sd.Users[0].Token,
			StatusCode: http.StatusOK,
			Method:     http.MethodGet,
			GotResp:    &orderapp.Order{},
			ExpResp:    toAppOrderPtr(sd.Orders[0].Order),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
	}

	return table
}
