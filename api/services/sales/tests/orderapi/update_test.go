package orderapi_test

import (
	"fmt"
	"net/http"

	"github.com/google/go-cmp/cmp"
	"github.com/warlck/food-flow/app/domain/orderapp"
	"github.com/warlck/food-flow/app/sdk/apitest"
	"github.com/warlck/food-flow/app/sdk/errs"
)

func updateStatus200(sd apitest.SeedData) []apitest.Table {
	orderStatus := "confirmed"
	paymentStatus := "paid"
	outForDelivery := "out_for_delivery"

	table := []apitest.Table{
		{
			Name:       "update-order-status",
			URL:        fmt.Sprintf("/v1/orders/%s/status", sd.Orders[0].ID),
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPatch,
			StatusCode: http.StatusOK,
			Input: &orderapp.UpdateOrderStatus{
				OrderStatus: &orderStatus,
			},
			GotResp: &orderapp.Order{},
			ExpResp: &orderapp.Order{
				ID:            sd.Orders[0].ID.String(),
				RestaurantID:  sd.Orders[0].RestaurantID.String(),
				CustomerName:  sd.Orders[0].CustomerName,
				CustomerEmail: sd.Orders[0].CustomerEmail,
				CustomerPhone: sd.Orders[0].CustomerPhone,
				OrderType:     sd.Orders[0].OrderType,
				OrderStatus:   "confirmed",
				PaymentStatus: sd.Orders[0].PaymentStatus,
				PaymentMethod: sd.Orders[0].PaymentMethod,
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.(*orderapp.Order)
				if !exists {
					return "error occurred"
				}

				expResp := exp.(*orderapp.Order)

				// Copy dynamic fields from got to exp
				expResp.Subtotal = gotResp.Subtotal
				expResp.DeliveryFee = gotResp.DeliveryFee
				expResp.Tax = gotResp.Tax
				expResp.Total = gotResp.Total
				expResp.Items = gotResp.Items
				expResp.DeliveryAddress = gotResp.DeliveryAddress
				expResp.SpecialInstructions = gotResp.SpecialInstructions
				expResp.DateCreated = gotResp.DateCreated
				expResp.DateUpdated = gotResp.DateUpdated

				return cmp.Diff(gotResp, expResp)
			},
		},
		{
			Name:       "update-payment-status",
			URL:        fmt.Sprintf("/v1/orders/%s/status", sd.Orders[1].ID),
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPatch,
			StatusCode: http.StatusOK,
			Input: &orderapp.UpdateOrderStatus{
				PaymentStatus: &paymentStatus,
			},
			GotResp: &orderapp.Order{},
			ExpResp: &orderapp.Order{
				ID:            sd.Orders[1].ID.String(),
				RestaurantID:  sd.Orders[1].RestaurantID.String(),
				CustomerName:  sd.Orders[1].CustomerName,
				CustomerEmail: sd.Orders[1].CustomerEmail,
				CustomerPhone: sd.Orders[1].CustomerPhone,
				OrderType:     sd.Orders[1].OrderType,
				OrderStatus:   sd.Orders[1].OrderStatus,
				PaymentStatus: "paid",
				PaymentMethod: sd.Orders[1].PaymentMethod,
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.(*orderapp.Order)
				if !exists {
					return "error occurred"
				}

				expResp := exp.(*orderapp.Order)

				// Copy dynamic fields from got to exp
				expResp.Subtotal = gotResp.Subtotal
				expResp.DeliveryFee = gotResp.DeliveryFee
				expResp.Tax = gotResp.Tax
				expResp.Total = gotResp.Total
				expResp.Items = gotResp.Items
				expResp.DeliveryAddress = gotResp.DeliveryAddress
				expResp.SpecialInstructions = gotResp.SpecialInstructions
				expResp.DateCreated = gotResp.DateCreated
				expResp.DateUpdated = gotResp.DateUpdated

				return cmp.Diff(gotResp, expResp)
			},
		},
		{
			Name:       "update-order-status-out-for-delivery",
			URL:        fmt.Sprintf("/v1/orders/%s/status", sd.Orders[0].ID),
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPatch,
			StatusCode: http.StatusOK,
			Input: &orderapp.UpdateOrderStatus{
				OrderStatus: &outForDelivery,
			},
			GotResp: &orderapp.Order{},
			ExpResp: &orderapp.Order{
				ID:            sd.Orders[0].ID.String(),
				RestaurantID:  sd.Orders[0].RestaurantID.String(),
				CustomerName:  sd.Orders[0].CustomerName,
				CustomerEmail: sd.Orders[0].CustomerEmail,
				CustomerPhone: sd.Orders[0].CustomerPhone,
				OrderType:     sd.Orders[0].OrderType,
				OrderStatus:   "out_for_delivery",
				PaymentStatus: sd.Orders[0].PaymentStatus,
				PaymentMethod: sd.Orders[0].PaymentMethod,
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.(*orderapp.Order)
				if !exists {
					return "error occurred"
				}

				expResp := exp.(*orderapp.Order)

				// Copy dynamic fields from got to exp
				expResp.Subtotal = gotResp.Subtotal
				expResp.DeliveryFee = gotResp.DeliveryFee
				expResp.Tax = gotResp.Tax
				expResp.Total = gotResp.Total
				expResp.Items = gotResp.Items
				expResp.DeliveryAddress = gotResp.DeliveryAddress
				expResp.SpecialInstructions = gotResp.SpecialInstructions
				expResp.DateCreated = gotResp.DateCreated
				expResp.DateUpdated = gotResp.DateUpdated

				return cmp.Diff(gotResp, expResp)
			},
		},
	}

	return table
}

func updateStatus400(sd apitest.SeedData) []apitest.Table {
	outForDelivery := "out_for_delivery"

	table := []apitest.Table{
		{
			Name:       "out-for-delivery-on-pickup-order",
			URL:        fmt.Sprintf("/v1/orders/%s/status", sd.Orders[1].ID),
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPatch,
			StatusCode: http.StatusBadRequest,
			Input: &orderapp.UpdateOrderStatus{
				OrderStatus: &outForDelivery,
			},
			GotResp: &errs.Error{},
			ExpResp: errs.Newf(errs.InvalidArgument, "out for delivery status requires a delivery order"),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
	}

	return table
}

func updateStatus401(sd apitest.SeedData) []apitest.Table {
	orderStatus := "confirmed"

	table := []apitest.Table{
		{
			Name:       "non-admin-user",
			URL:        fmt.Sprintf("/v1/orders/%s/status", sd.Orders[0].ID),
			Token:      sd.Users[0].Token,
			Method:     http.MethodPatch,
			StatusCode: http.StatusForbidden,
			Input: &orderapp.UpdateOrderStatus{
				OrderStatus: &orderStatus,
			},
			GotResp: &errs.Error{},
			ExpResp: errs.Newf(errs.PermissionDenied, "authorize: you are not authorized for that action, claims[[USER]] rule[rule_admin_only]"),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
		{
			Name:       "empty-token",
			URL:        fmt.Sprintf("/v1/orders/%s/status", sd.Orders[0].ID),
			Token:      "&nbsp;",
			Method:     http.MethodPatch,
			StatusCode: http.StatusUnauthorized,
			Input: &orderapp.UpdateOrderStatus{
				OrderStatus: &orderStatus,
			},
			GotResp: &errs.Error{},
			ExpResp: errs.Newf(errs.Unauthenticated, "error parsing token: token is malformed: token contains an invalid number of segments"),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
	}

	return table
}
