package orderapi_test

import (
	"net/http"

	"github.com/google/go-cmp/cmp"
	"github.com/warlck/food-flow/app/domain/orderapp"
	"github.com/warlck/food-flow/app/sdk/apitest"
	"github.com/warlck/food-flow/app/sdk/errs"
)

func create201(sd apitest.SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:       "pickup-order",
			URL:        "/v1/orders",
			Token:      sd.Users[0].Token,
			Method:     http.MethodPost,
			StatusCode: http.StatusCreated,
			Input: &orderapp.NewOrder{
				RestaurantID:  sd.Restaurants[0].ID.String(),
				CustomerName:  "John Doe",
				CustomerEmail: "john@example.com",
				CustomerPhone: "555-1234",
				OrderType:     "pickup",
				PaymentMethod: "creditCard",
				Items: []orderapp.NewOrderItem{
					{
						MenuItemID: sd.MenuItems[0].ID.String(),
						Quantity:   2,
					},
				},
			},
			GotResp: &orderapp.Order{},
			ExpResp: &orderapp.Order{
				RestaurantID:  sd.Restaurants[0].ID.String(),
				CustomerName:  "John Doe",
				CustomerEmail: "john@example.com",
				CustomerPhone: "555-1234",
				OrderType:     "pickup",
				OrderStatus:   "pending",
				PaymentStatus: "pending",
				PaymentMethod: "creditCard",
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.(*orderapp.Order)
				if !exists {
					return "error occurred"
				}

				expResp := exp.(*orderapp.Order)

				// Copy dynamic fields from got to exp
				expResp.ID = gotResp.ID
				expResp.Subtotal = gotResp.Subtotal
				expResp.DeliveryFee = gotResp.DeliveryFee
				expResp.Tax = gotResp.Tax
				expResp.Total = gotResp.Total
				expResp.Items = gotResp.Items
				expResp.DateCreated = gotResp.DateCreated
				expResp.DateUpdated = gotResp.DateUpdated

				return cmp.Diff(gotResp, expResp)
			},
		},
		{
			Name:       "delivery-order",
			URL:        "/v1/orders",
			Token:      sd.Users[0].Token,
			Method:     http.MethodPost,
			StatusCode: http.StatusCreated,
			Input: &orderapp.NewOrder{
				RestaurantID:        sd.Restaurants[0].ID.String(),
				CustomerName:        "Jane Smith",
				CustomerEmail:       "jane@example.com",
				CustomerPhone:       "555-5678",
				OrderType:           "delivery",
				PaymentMethod:       "creditCard",
				SpecialInstructions: "Ring the doorbell",
				Items: []orderapp.NewOrderItem{
					{
						MenuItemID:          sd.MenuItems[0].ID.String(),
						Quantity:            1,
						SpecialInstructions: "Extra sauce",
					},
					{
						MenuItemID: sd.MenuItems[1].ID.String(),
						Quantity:   2,
					},
				},
				DeliveryAddress: &orderapp.NewDeliveryAddress{
					Street:               "123 Main St",
					City:                 "Anytown",
					State:                "CA",
					PostalCode:           "12345",
					DeliveryInstructions: "Leave at door",
					Latitude:             ptr(1.30719),
					Longitude:            ptr(103.87434),
				},
			},
			GotResp: &orderapp.Order{},
			ExpResp: &orderapp.Order{
				RestaurantID:        sd.Restaurants[0].ID.String(),
				CustomerName:        "Jane Smith",
				CustomerEmail:       "jane@example.com",
				CustomerPhone:       "555-5678",
				OrderType:           "delivery",
				OrderStatus:         "pending",
				PaymentStatus:       "pending",
				PaymentMethod:       "creditCard",
				SpecialInstructions: "Ring the doorbell",
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.(*orderapp.Order)
				if !exists {
					return "error occurred"
				}

				expResp := exp.(*orderapp.Order)

				// Copy dynamic fields from got to exp
				expResp.ID = gotResp.ID
				expResp.Subtotal = gotResp.Subtotal
				expResp.DeliveryFee = gotResp.DeliveryFee
				expResp.Tax = gotResp.Tax
				expResp.Total = gotResp.Total
				expResp.Items = gotResp.Items
				expResp.DeliveryAddress = gotResp.DeliveryAddress
				expResp.DateCreated = gotResp.DateCreated
				expResp.DateUpdated = gotResp.DateUpdated

				return cmp.Diff(gotResp, expResp)
			},
		},
	}

	return table
}

func create400(sd apitest.SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:       "missing-input",
			URL:        "/v1/orders",
			Token:      sd.Users[0].Token,
			Method:     http.MethodPost,
			StatusCode: http.StatusBadRequest,
			Input:      &orderapp.NewOrder{},
			GotResp:    &errs.Error{},
			ExpResp:    errs.Newf(errs.InvalidArgument, "validate: [{\"field\":\"restaurantId\",\"error\":\"restaurantId is a required field\"},{\"field\":\"customerName\",\"error\":\"customerName is a required field\"},{\"field\":\"customerEmail\",\"error\":\"customerEmail is a required field\"},{\"field\":\"customerPhone\",\"error\":\"customerPhone is a required field\"},{\"field\":\"orderType\",\"error\":\"orderType is a required field\"},{\"field\":\"paymentMethod\",\"error\":\"paymentMethod is a required field\"},{\"field\":\"items\",\"error\":\"items is a required field\"}]"),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
		{
			Name:       "invalid-order-type",
			URL:        "/v1/orders",
			Token:      sd.Users[0].Token,
			Method:     http.MethodPost,
			StatusCode: http.StatusBadRequest,
			Input: &orderapp.NewOrder{
				RestaurantID:  sd.Restaurants[0].ID.String(),
				CustomerName:  "John Doe",
				CustomerEmail: "john@example.com",
				CustomerPhone: "555-1234",
				OrderType:     "invalid",
				PaymentMethod: "creditCard",
				Items: []orderapp.NewOrderItem{
					{
						MenuItemID: sd.MenuItems[0].ID.String(),
						Quantity:   1,
					},
				},
			},
			GotResp: &errs.Error{},
			ExpResp: errs.Newf(errs.InvalidArgument, "validate: [{\"field\":\"orderType\",\"error\":\"orderType must be one of [pickup delivery]\"}]"),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
		{
			Name:       "invalid-email",
			URL:        "/v1/orders",
			Token:      sd.Users[0].Token,
			Method:     http.MethodPost,
			StatusCode: http.StatusBadRequest,
			Input: &orderapp.NewOrder{
				RestaurantID:  sd.Restaurants[0].ID.String(),
				CustomerName:  "John Doe",
				CustomerEmail: "invalid-email",
				CustomerPhone: "555-1234",
				OrderType:     "pickup",
				PaymentMethod: "creditCard",
				Items: []orderapp.NewOrderItem{
					{
						MenuItemID: sd.MenuItems[0].ID.String(),
						Quantity:   1,
					},
				},
			},
			GotResp: &errs.Error{},
			ExpResp: errs.Newf(errs.InvalidArgument, "validate: [{\"field\":\"customerEmail\",\"error\":\"customerEmail must be a valid email address\"}]"),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
		{
			Name:       "empty-items",
			URL:        "/v1/orders",
			Token:      sd.Users[0].Token,
			Method:     http.MethodPost,
			StatusCode: http.StatusBadRequest,
			Input: &orderapp.NewOrder{
				RestaurantID:  sd.Restaurants[0].ID.String(),
				CustomerName:  "John Doe",
				CustomerEmail: "john@example.com",
				CustomerPhone: "555-1234",
				OrderType:     "pickup",
				PaymentMethod: "creditCard",
				Items:         []orderapp.NewOrderItem{},
			},
			GotResp: &errs.Error{},
			ExpResp: errs.Newf(errs.InvalidArgument, "validate: [{\"field\":\"items\",\"error\":\"items must contain at least 1 item\"}]"),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
	}

	return table
}

func create401(sd apitest.SeedData) []apitest.Table {
	table := []apitest.Table{}
	return table
}

// ptr is a helper function for getting the address of a value.
func ptr[T any](v T) *T {
	return &v
}
