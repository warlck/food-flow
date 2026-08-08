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

func create201(sd apitest.SeedData) []apitest.Table {
	// Expected totals for the "order-with-addons" case: menu item 0 (qty 2)
	// with addon 0 (qty 1) and addon 1 (qty 2). Addon cost is addon price *
	// addon qty * item qty.
	expAddonSubtotal := sd.MenuItems[0].Price.Value()*2 +
		sd.Addons[0].Price.Value()*1*2 +
		sd.Addons[1].Price.Value()*2*2
	expAddonSubtotal = math.Round(expAddonSubtotal*100) / 100
	expAddonTax := math.Round(expAddonSubtotal*sd.Restaurants[0].TaxRate*100) / 100
	expAddonTotal := math.Round((expAddonSubtotal+expAddonTax)*100) / 100

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
		{
			Name:       "order-with-addons",
			URL:        "/v1/orders",
			Token:      sd.Users[0].Token,
			Method:     http.MethodPost,
			StatusCode: http.StatusCreated,
			Input: &orderapp.NewOrder{
				RestaurantID:  sd.Restaurants[0].ID.String(),
				CustomerName:  "Addon Customer",
				CustomerEmail: "addon@example.com",
				CustomerPhone: "555-4321",
				OrderType:     "pickup",
				PaymentMethod: "creditCard",
				Items: []orderapp.NewOrderItem{
					{
						MenuItemID: sd.MenuItems[0].ID.String(),
						Quantity:   2,
						Addons: []orderapp.NewOrderItemAddon{
							{AddonID: sd.Addons[0].ID.String(), Quantity: 1},
							{AddonID: sd.Addons[1].ID.String(), Quantity: 2},
						},
					},
				},
			},
			GotResp: &orderapp.Order{},
			ExpResp: &orderapp.Order{
				RestaurantID:  sd.Restaurants[0].ID.String(),
				CustomerName:  "Addon Customer",
				CustomerEmail: "addon@example.com",
				CustomerPhone: "555-4321",
				OrderType:     "pickup",
				OrderStatus:   "pending",
				PaymentStatus: "pending",
				PaymentMethod: "creditCard",
				Subtotal:      expAddonSubtotal,
				Tax:           expAddonTax,
				Total:         expAddonTotal,
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.(*orderapp.Order)
				if !exists {
					return "error occurred"
				}

				expResp := exp.(*orderapp.Order)

				// Copy dynamic fields from got to exp
				expResp.ID = gotResp.ID
				expResp.DeliveryFee = gotResp.DeliveryFee
				expResp.DateCreated = gotResp.DateCreated
				expResp.DateUpdated = gotResp.DateUpdated

				if len(gotResp.Items) != 1 {
					return fmt.Sprintf("expected 1 item, got %d", len(gotResp.Items))
				}

				item := gotResp.Items[0]
				if len(item.Addons) != 2 {
					return fmt.Sprintf("expected 2 addons, got %d", len(item.Addons))
				}

				expAddons := map[string]struct {
					name     string
					price    float64
					quantity int
				}{
					sd.Addons[0].ID.String(): {sd.Addons[0].Name.String(), sd.Addons[0].Price.Value(), 1},
					sd.Addons[1].ID.String(): {sd.Addons[1].Name.String(), sd.Addons[1].Price.Value(), 2},
				}

				for _, a := range item.Addons {
					expA, ok := expAddons[a.AddonID]
					if !ok {
						return fmt.Sprintf("unexpected addon %s", a.AddonID)
					}

					if a.AddonName != expA.name {
						return fmt.Sprintf("addon name: got %q, want %q", a.AddonName, expA.name)
					}

					if a.AddonPrice != expA.price {
						return fmt.Sprintf("addon price: got %.2f, want %.2f", a.AddonPrice, expA.price)
					}

					if a.Quantity != expA.quantity {
						return fmt.Sprintf("addon quantity: got %d, want %d", a.Quantity, expA.quantity)
					}
				}

				// Verify the item addons are included in the response while
				// ignoring dynamic item fields.
				expResp.Items = gotResp.Items

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
		{
			Name:       "addon-quantity-zero",
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
				Items: []orderapp.NewOrderItem{
					{
						MenuItemID: sd.MenuItems[0].ID.String(),
						Quantity:   1,
						Addons: []orderapp.NewOrderItemAddon{
							{AddonID: sd.Addons[0].ID.String(), Quantity: 0},
						},
					},
				},
			},
			GotResp: &errs.Error{},
			ExpResp: errs.Newf(errs.InvalidArgument, "validate: [{\"field\":\"quantity\",\"error\":\"quantity is a required field\"}]"),
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
