package promoapi_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/warlck/food-flow/app/domain/promoapp"
	"github.com/warlck/food-flow/app/sdk/apitest"
	"github.com/warlck/food-flow/app/sdk/errs"
	"github.com/warlck/food-flow/app/sdk/query"
)

func Test_Promo(t *testing.T) {
	t.Parallel()

	test := apitest.New(t, "Test_Promo")

	sd, err := insertSeedData(test.DB, test.Auth)
	if err != nil {
		t.Fatalf("Seeding error: %s", err)
	}

	test.Run(t, validate200(sd), "validate-200")
	test.Run(t, validate400(sd), "validate-400")

	test.Run(t, query200(sd), "query-200")
	test.Run(t, queryByID404(sd), "querybyid-404")
	test.Run(t, create201(sd), "create-201")
	test.Run(t, create400(sd), "create-400")
	test.Run(t, create403(sd), "create-403")
	test.Run(t, update200(sd), "update-200")
	test.Run(t, update400(sd), "update-400")
	test.Run(t, update403(sd), "update-403")
	test.Run(t, update404(sd), "update-404")
	test.Run(t, delete200(sd), "delete-200")
	test.Run(t, delete403(sd), "delete-403")
	test.Run(t, delete404(sd), "delete-404")
}

func validate200(sd SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:       "valid-promo-code",
			URL:        "/v1/promotions/validate",
			Token:      "",
			Method:     http.MethodPost,
			StatusCode: http.StatusOK,
			Input: &promoapp.ValidateRequest{
				PromoCode: sd.Promotions[0].Code,
				Subtotal:  50.0,
			},
			GotResp: &promoapp.ValidateResponse{},
			ExpResp: &promoapp.ValidateResponse{
				Valid:         true,
				Code:          sd.Promotions[0].Code,
				DiscountType:  "percentage",
				DiscountValue: sd.Promotions[0].DiscountValue,

				DiscountAmount: 5.0,
				FinalSubtotal:  45.0,
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, ok := got.(*promoapp.ValidateResponse)
				if !ok {
					return "error occurred"
				}
				expResp := exp.(*promoapp.ValidateResponse)
				return cmp.Diff(gotResp, expResp)
			},
		},
		{
			Name:       "invalid-promo-code",
			URL:        "/v1/promotions/validate",
			Token:      "",
			Method:     http.MethodPost,
			StatusCode: http.StatusOK,
			Input: &promoapp.ValidateRequest{
				PromoCode: "NOTREALCODE",
				Subtotal:  50.0,
			},
			GotResp: &promoapp.ValidateResponse{},
			ExpResp: &promoapp.ValidateResponse{
				Valid:  false,
				Reason: "Invalid promo code",
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, ok := got.(*promoapp.ValidateResponse)
				if !ok {
					return "error occurred"
				}
				expResp := exp.(*promoapp.ValidateResponse)
				return cmp.Diff(gotResp, expResp)
			},
		},
		{
			Name:       "subtotal-below-minimum",
			URL:        "/v1/promotions/validate",
			Token:      "",
			Method:     http.MethodPost,
			StatusCode: http.StatusOK,
			Input: &promoapp.ValidateRequest{
				PromoCode: sd.Promotions[0].Code,
				Subtotal:  5.0,
			},
			GotResp: &promoapp.ValidateResponse{},
			ExpResp: &promoapp.ValidateResponse{
				Valid:  false,
				Reason: fmt.Sprintf("Minimum order subtotal of $%.2f required for this promo code", sd.Promotions[0].MinOrderAmount),
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, ok := got.(*promoapp.ValidateResponse)
				if !ok {
					return "error occurred"
				}
				expResp := exp.(*promoapp.ValidateResponse)
				return cmp.Diff(gotResp, expResp)
			},
		},
		{
			Name:       "disabled-promo-code",
			URL:        "/v1/promotions/validate",
			Token:      "",
			Method:     http.MethodPost,
			StatusCode: http.StatusOK,
			Input: &promoapp.ValidateRequest{
				PromoCode: "DISABLED10",
				Subtotal:  50.0,
			},
			GotResp: &promoapp.ValidateResponse{},
			ExpResp: &promoapp.ValidateResponse{
				Valid:  false,
				Reason: "Promo code is inactive",
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, ok := got.(*promoapp.ValidateResponse)
				if !ok {
					return "error occurred"
				}
				expResp := exp.(*promoapp.ValidateResponse)
				return cmp.Diff(gotResp, expResp)
			},
		},
		{
			Name:       "restaurant-mismatch",
			URL:        "/v1/promotions/validate",
			Token:      "",
			Method:     http.MethodPost,
			StatusCode: http.StatusOK,
			Input: &promoapp.ValidateRequest{
				PromoCode:    "REST2ONLY",
				RestaurantID: ptr(sd.Restaurants[0].ID.String()),
				Subtotal:     50.0,
			},
			GotResp: &promoapp.ValidateResponse{},
			ExpResp: &promoapp.ValidateResponse{
				Valid:  false,
				Reason: "Promo code is not applicable to this restaurant",
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, ok := got.(*promoapp.ValidateResponse)
				if !ok {
					return "error occurred"
				}
				expResp := exp.(*promoapp.ValidateResponse)
				return cmp.Diff(gotResp, expResp)
			},
		},
	}

	return table
}

func validate400(sd SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:       "missing-promo-code",
			URL:        "/v1/promotions/validate",
			Token:      "",
			Method:     http.MethodPost,
			StatusCode: http.StatusBadRequest,
			Input:      &promoapp.ValidateRequest{},
			GotResp:    &errs.Error{},
			ExpResp:    errs.Newf(errs.InvalidArgument, "validate: [{\"field\":\"promoCode\",\"error\":\"promoCode is a required field\"}]"),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
	}

	return table
}

func query200(sd SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:       "admin-list-promotions",
			URL:        "/v1/promotions?page=1&rows=10",
			Token:      sd.Admins[0].Token,
			Method:     http.MethodGet,
			StatusCode: http.StatusOK,
			GotResp:    &query.Result[promoapp.Promotion]{},
			ExpResp:    &query.Result[promoapp.Promotion]{},
			CmpFunc: func(got any, exp any) string {
				gotResp, ok := got.(*query.Result[promoapp.Promotion])
				if !ok {
					return "error occurred"
				}
				if len(gotResp.Items) < 3 {
					return fmt.Sprintf("expected at least 3 promotions, got %d", len(gotResp.Items))
				}
				return ""
			},
		},
	}

	return table
}

func create201(sd SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:       "admin-create-promotion",
			URL:        "/v1/promotions",
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPost,
			StatusCode: http.StatusCreated,
			Input: &promoapp.NewPromotion{
				RestaurantID:   ptr(sd.Restaurants[0].ID.String()),
				Code:           "SUMMER50",
				Name:           "Summer Sale 50",
				Description:    "50 percent discount",
				DiscountType:   "percentage",
				DiscountValue:  50.0,
				MinOrderAmount: 20.0,
				Enabled:        true,
			},
			GotResp: &promoapp.Promotion{},
			ExpResp: &promoapp.Promotion{
				RestaurantID:   ptr(sd.Restaurants[0].ID.String()),
				Code:           "SUMMER50",
				Name:           "Summer Sale 50",
				Description:    "50 percent discount",
				DiscountType:   "percentage",
				DiscountValue:  50.0,
				MinOrderAmount: 20.0,
				Enabled:        true,
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, ok := got.(*promoapp.Promotion)
				if !ok {
					return "error occurred"
				}
				expResp := exp.(*promoapp.Promotion)

				expResp.ID = gotResp.ID
				expResp.DateCreated = gotResp.DateCreated
				expResp.DateUpdated = gotResp.DateUpdated

				return cmp.Diff(gotResp, expResp)
			},
		},
	}

	return table
}

func create400(sd SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:       "missing-restaurant-id",
			URL:        "/v1/promotions",
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPost,
			StatusCode: http.StatusBadRequest,
			Input: &promoapp.NewPromotion{
				Code:           "NOREST10",
				Name:           "No Restaurant Promo",
				DiscountType:   "percentage",
				DiscountValue:  10.0,
				MinOrderAmount: 10.0,
				Enabled:        true,
			},
			GotResp: &errs.Error{},
			ExpResp: errs.Newf(errs.InvalidArgument, "validate: [{\"field\":\"restaurantId\",\"error\":\"restaurantId is a required field\"}]"),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
	}

	return table
}

func create403(sd SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:       "create-promotion-other-org",
			URL:        "/v1/promotions",
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPost,
			StatusCode: http.StatusForbidden,
			Input: &promoapp.NewPromotion{
				RestaurantID:   ptr(sd.Promotions[6].RestaurantID.String()),
				Code:           "CROSSORG10",
				Name:           "Cross Org Promo",
				DiscountType:   "percentage",
				DiscountValue:  10.0,
				MinOrderAmount: 10.0,
				Enabled:        true,
			},
			GotResp: &errs.Error{},
			ExpResp: errs.Newf(errs.PermissionDenied, "user not in organization %s", sd.Orgs[1].ID),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
	}

	return table
}

func update200(sd SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:       "admin-update-promotion",
			URL:        fmt.Sprintf("/v1/promotions/%s", sd.Promotions[4].ID),
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPut,
			StatusCode: http.StatusOK,
			Input: &promoapp.UpdatePromotion{
				Name: ptr("Updated Promo Name"),
			},
			GotResp: &promoapp.Promotion{},
			ExpResp: &promoapp.Promotion{
				ID:   sd.Promotions[4].ID.String(),
				Code: sd.Promotions[4].Code,
				Name: "Updated Promo Name",
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, ok := got.(*promoapp.Promotion)
				if !ok {
					return "error occurred"
				}
				expResp := exp.(*promoapp.Promotion)

				expResp.RestaurantID = gotResp.RestaurantID
				expResp.Description = gotResp.Description
				expResp.DiscountType = gotResp.DiscountType
				expResp.DiscountValue = gotResp.DiscountValue
				expResp.MinOrderAmount = gotResp.MinOrderAmount
				expResp.MaxDiscountAmount = gotResp.MaxDiscountAmount
				expResp.UsageLimit = gotResp.UsageLimit
				expResp.UsageCount = gotResp.UsageCount
				expResp.Enabled = gotResp.Enabled
				expResp.DateCreated = gotResp.DateCreated
				expResp.DateUpdated = gotResp.DateUpdated

				return cmp.Diff(gotResp, expResp)
			},
		},
	}

	return table
}

func update400(sd SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:       "update-clear-restaurant-id",
			URL:        fmt.Sprintf("/v1/promotions/%s", sd.Promotions[4].ID),
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPut,
			StatusCode: http.StatusBadRequest,
			Input: &promoapp.UpdatePromotion{
				RestaurantID: ptr(ptr("")),
			},
			GotResp: &errs.Error{},
			ExpResp: errs.Newf(errs.InvalidArgument, "promotion restaurantId cannot be cleared"),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
		{
			Name:       "update-invalid-percentage",
			URL:        fmt.Sprintf("/v1/promotions/%s", sd.Promotions[4].ID),
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPut,
			StatusCode: http.StatusBadRequest,
			Input: &promoapp.UpdatePromotion{
				DiscountType:  ptr("percentage"),
				DiscountValue: ptr(150.0),
			},
			GotResp: &errs.Error{},
			ExpResp: errs.Newf(errs.InvalidArgument, "validate: percentage discount value cannot exceed 100"),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
	}

	return table
}

func update403(sd SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:       "update-global-promotion",
			URL:        fmt.Sprintf("/v1/promotions/%s", sd.Promotions[0].ID),
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPut,
			StatusCode: http.StatusForbidden,
			Input: &promoapp.UpdatePromotion{
				Name: ptr("Updated Promo Name"),
			},
			GotResp: &errs.Error{},
			ExpResp: errs.Newf(errs.PermissionDenied, "promotion %s is global; global promotions are managed by admin tooling", sd.Promotions[0].ID),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
		{
			Name:       "update-other-org-promotion",
			URL:        fmt.Sprintf("/v1/promotions/%s", sd.Promotions[6].ID),
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPut,
			StatusCode: http.StatusForbidden,
			Input: &promoapp.UpdatePromotion{
				Name: ptr("Updated Promo Name"),
			},
			GotResp: &errs.Error{},
			ExpResp: errs.Newf(errs.PermissionDenied, "user not in organization %s", sd.Orgs[1].ID),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
	}

	return table
}

func delete200(sd SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:       "admin-delete-promotion",
			URL:        fmt.Sprintf("/v1/promotions/%s", sd.Promotions[5].ID),
			Token:      sd.Admins[0].Token,
			Method:     http.MethodDelete,
			StatusCode: http.StatusNoContent,
			CmpFunc: func(got any, exp any) string {
				return ""
			},
		},
	}

	return table
}

func delete403(sd SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:       "delete-global-promotion",
			URL:        fmt.Sprintf("/v1/promotions/%s", sd.Promotions[1].ID),
			Token:      sd.Admins[0].Token,
			Method:     http.MethodDelete,
			StatusCode: http.StatusForbidden,
			GotResp:    &errs.Error{},
			ExpResp:    errs.Newf(errs.PermissionDenied, "promotion %s is global; global promotions are managed by admin tooling", sd.Promotions[1].ID),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
		{
			Name:       "delete-other-org-promotion",
			URL:        fmt.Sprintf("/v1/promotions/%s", sd.Promotions[6].ID),
			Token:      sd.Admins[0].Token,
			Method:     http.MethodDelete,
			StatusCode: http.StatusForbidden,
			GotResp:    &errs.Error{},
			ExpResp:    errs.Newf(errs.PermissionDenied, "user not in organization %s", sd.Orgs[1].ID),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
	}

	return table
}

func queryByID404(sd SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:       "querybyid-not-found",
			URL:        fmt.Sprintf("/v1/promotions/%s", uuid.New()),
			Token:      sd.Admins[0].Token,
			Method:     http.MethodGet,
			StatusCode: http.StatusNotFound,
			GotResp:    &errs.Error{},
			ExpResp:    errs.Newf(errs.NotFound, "query by id: promotion not found"),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
	}

	return table
}

func update404(sd SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:       "update-not-found",
			URL:        fmt.Sprintf("/v1/promotions/%s", uuid.New()),
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPut,
			StatusCode: http.StatusNotFound,
			Input: &promoapp.UpdatePromotion{
				Name: ptr("Updated Promo Name"),
			},
			GotResp: &errs.Error{},
			ExpResp: errs.Newf(errs.NotFound, "query by id: promotion not found"),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
	}

	return table
}

func delete404(sd SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:       "delete-not-found",
			URL:        fmt.Sprintf("/v1/promotions/%s", uuid.New()),
			Token:      sd.Admins[0].Token,
			Method:     http.MethodDelete,
			StatusCode: http.StatusNotFound,
			GotResp:    &errs.Error{},
			ExpResp:    errs.Newf(errs.NotFound, "query by id: promotion not found"),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
	}

	return table
}

func ptr[T any](v T) *T {
	return &v
}
