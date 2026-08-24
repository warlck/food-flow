package orderapi_test

import (
	"fmt"
	"math"
	"net/http"

	"github.com/google/go-cmp/cmp"
	"github.com/warlck/food-flow/app/domain/orderapp"
	"github.com/warlck/food-flow/app/sdk/apitest"
	"github.com/warlck/food-flow/app/sdk/errs"
)

func insights200(sd apitest.SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:       "insights-all-restaurants",
			URL:        "/v1/insights",
			Token:      sd.Admins[0].Token,
			StatusCode: http.StatusOK,
			Method:     http.MethodGet,
			GotResp:    &orderapp.AppInsights{},
			ExpResp:    &orderapp.AppInsights{},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.(*orderapp.AppInsights)
				if !exists {
					return "error occurred: could not cast response"
				}

				if gotResp.Summary.TotalOrders <= 0 {
					return fmt.Sprintf("expected positive total orders, got %d", gotResp.Summary.TotalOrders)
				}

				// Accounting Identity 1: GrossSales - Discounts == NetSales
				grossMinusDisc := gotResp.Summary.GrossSales - gotResp.Summary.TotalDiscounts
				if math.Abs(grossMinusDisc-gotResp.Summary.NetSales) > 0.01 {
					return fmt.Sprintf("identity 1 failed: GrossSales (%.2f) - Discounts (%.2f) != NetSales (%.2f)",
						gotResp.Summary.GrossSales, gotResp.Summary.TotalDiscounts, gotResp.Summary.NetSales)
				}

				// Accounting Identity 2: NetSales + DeliveryFees + Tax == TotalCollected
				netPlusFeesTax := gotResp.Summary.NetSales + gotResp.Summary.TotalDeliveryFees + gotResp.Summary.TotalTax
				if math.Abs(netPlusFeesTax-gotResp.Summary.TotalCollected) > 0.01 {
					return fmt.Sprintf("identity 2 failed: NetSales (%.2f) + Deliv (%.2f) + Tax (%.2f) != TotalCollected (%.2f)",
						gotResp.Summary.NetSales, gotResp.Summary.TotalDeliveryFees, gotResp.Summary.TotalTax, gotResp.Summary.TotalCollected)
				}

				// Accounting Identity 3: AOV == TotalCollected / CompletedOrders
				if gotResp.Summary.CompletedOrders > 0 {
					expectedAOV := gotResp.Summary.TotalCollected / float64(gotResp.Summary.CompletedOrders)
					if math.Abs(expectedAOV-gotResp.Summary.AverageOrderValue) > 0.01 {
						return fmt.Sprintf("identity 3 failed: AOV (%.2f) != TotalCollected (%.2f) / CompletedOrders (%d) = %.2f",
							gotResp.Summary.AverageOrderValue, gotResp.Summary.TotalCollected, gotResp.Summary.CompletedOrders, expectedAOV)
					}
				}

				// Dataset slices verification
				if len(gotResp.SalesOverTime) == 0 {
					return "expected non-empty SalesOverTime"
				}
				if len(gotResp.TopItems) == 0 {
					return "expected non-empty TopItems"
				}
				if len(gotResp.TopCategories) == 0 {
					return "expected non-empty TopCategories"
				}
				if len(gotResp.TopAddons) == 0 {
					return "expected non-empty TopAddons"
				}
				if len(gotResp.OrderTypes) == 0 {
					return "expected non-empty OrderTypes"
				}
				if len(gotResp.PeakHours) == 0 {
					return "expected non-empty PeakHours"
				}

				return ""
			},
		},
		{
			Name:       "insights-with-restaurant-filter",
			URL:        "/v1/insights?restaurant_id=" + sd.Restaurants[0].ID.String(),
			Token:      sd.Admins[0].Token,
			StatusCode: http.StatusOK,
			Method:     http.MethodGet,
			GotResp:    &orderapp.AppInsights{},
			ExpResp:    &orderapp.AppInsights{},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.(*orderapp.AppInsights)
				if !exists {
					return "error occurred: could not cast response"
				}

				if gotResp.Summary.TotalOrders <= 0 {
					return fmt.Sprintf("expected positive total orders, got %d", gotResp.Summary.TotalOrders)
				}

				// Accounting Identity 1: GrossSales - Discounts == NetSales
				grossMinusDisc := gotResp.Summary.GrossSales - gotResp.Summary.TotalDiscounts
				if math.Abs(grossMinusDisc-gotResp.Summary.NetSales) > 0.01 {
					return fmt.Sprintf("identity 1 failed: GrossSales (%.2f) - Discounts (%.2f) != NetSales (%.2f)",
						gotResp.Summary.GrossSales, gotResp.Summary.TotalDiscounts, gotResp.Summary.NetSales)
				}

				// Accounting Identity 2: NetSales + DeliveryFees + Tax == TotalCollected
				netPlusFeesTax := gotResp.Summary.NetSales + gotResp.Summary.TotalDeliveryFees + gotResp.Summary.TotalTax
				if math.Abs(netPlusFeesTax-gotResp.Summary.TotalCollected) > 0.01 {
					return fmt.Sprintf("identity 2 failed: NetSales (%.2f) + Deliv (%.2f) + Tax (%.2f) != TotalCollected (%.2f)",
						gotResp.Summary.NetSales, gotResp.Summary.TotalDeliveryFees, gotResp.Summary.TotalTax, gotResp.Summary.TotalCollected)
				}

				return ""
			},
		},
	}

	return table
}

func insights401(sd apitest.SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:       "non-admin-user",
			URL:        "/v1/insights",
			Token:      sd.Users[0].Token,
			Method:     http.MethodGet,
			StatusCode: http.StatusForbidden,
			GotResp:    &errs.Error{},
			ExpResp:    errs.Newf(errs.PermissionDenied, "authorize: you are not authorized for that action, claims[[USER]] rule[rule_admin_only]"),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
		{
			Name:       "empty-token",
			URL:        "/v1/insights",
			Token:      "&nbsp;",
			Method:     http.MethodGet,
			StatusCode: http.StatusUnauthorized,
			GotResp:    &errs.Error{},
			ExpResp:    errs.Newf(errs.Unauthenticated, "error parsing token: token is malformed: token contains an invalid number of segments"),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
	}

	return table
}

func insights404(sd apitest.SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:       "restaurant-not-found",
			URL:        "/v1/insights?restaurant_id=00000000-0000-0000-0000-000000000000",
			Token:      sd.Admins[0].Token,
			Method:     http.MethodGet,
			StatusCode: http.StatusNotFound,
			GotResp:    &errs.Error{},
			ExpResp:    errs.Newf(errs.NotFound, "restaurant not found"),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
	}

	return table
}

