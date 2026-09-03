package orderapi_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/warlck/food-flow/app/domain/orderapp"
	"github.com/warlck/food-flow/app/sdk/apitest"
	"github.com/warlck/food-flow/app/sdk/errs"
	"github.com/warlck/food-flow/business/domain/categorybus"
	"github.com/warlck/food-flow/business/domain/menuitembus"
	"github.com/warlck/food-flow/business/domain/organizationbus"
	"github.com/warlck/food-flow/business/domain/restaurantbus"
	"github.com/warlck/food-flow/business/domain/userbus"
	"github.com/warlck/food-flow/business/types/money"
	"github.com/warlck/food-flow/business/types/name"
	"github.com/warlck/food-flow/business/types/role"
)

// Test_OrderMoneyBoundary verifies the money boundary at the order level and
// the insights level: amounts beyond the maximum supported value must surface
// as typed client errors (400), never clamped and never panicking.
//
// Fixture: a menu item priced at 90,000,000.00 in a restaurant with the seeded
// 10% tax rate, so a quantity-1 order totals exactly 99,000,000.00 (under the
// 99,999,999.99 cap) and two such orders aggregate to 198,000,000.00 (over
// the cap).
func Test_OrderMoneyBoundary(t *testing.T) {
	t.Parallel()

	test := apitest.New(t, "Test_OrderMoneyBoundary")

	// -------------------------------------------------------------------------

	ctx := context.Background()
	busDomain := test.DB.BusDomain

	admins, err := userbus.TestSeedUsers(ctx, 1, role.Admin, busDomain.User)
	if err != nil {
		t.Fatalf("seeding admin users: %s", err)
	}

	orgs, err := organizationbus.TestSeedOrganizations(ctx, 1, busDomain.Organization)
	if err != nil {
		t.Fatalf("seeding organizations: %s", err)
	}

	if _, err := busDomain.Organization.AddUser(ctx, organizationbus.NewOrganizationUser{
		OrganizationID: orgs[0].ID,
		UserID:         admins[0].ID,
		Role:           role.Admin,
	}); err != nil {
		t.Fatalf("adding admin to organization: %s", err)
	}

	restaurants, err := restaurantbus.TestSeedRestaurants(ctx, 1, busDomain.Restaurant, orgs[0].ID)
	if err != nil {
		t.Fatalf("seeding restaurants: %s", err)
	}

	categories, err := categorybus.TestSeedCategories(ctx, 1, restaurants[0].ID, busDomain.Category)
	if err != nil {
		t.Fatalf("seeding categories: %s", err)
	}

	mi, err := busDomain.MenuItem.Create(ctx, menuitembus.NewMenuItem{
		Name:         name.MustParse("Big Ticket Burger"),
		Description:  "Money boundary test fixture",
		Price:        money.MustParse(90_000_000.00),
		CategoryID:   categories[0].ID,
		RestaurantID: restaurants[0].ID,
	})
	if err != nil {
		t.Fatalf("seeding expensive menu item: %s", err)
	}

	adminToken := apitest.Token(busDomain, test.Auth, admins[0].Email.Address)

	// -------------------------------------------------------------------------
	// Order creation: values under the cap succeed; values over the cap
	// return a typed 400 instead of a clamp or an internal error.

	newOrder := func(qty int) *orderapp.NewOrder {
		return &orderapp.NewOrder{
			RestaurantID:  restaurants[0].ID.String(),
			CustomerName:  "Boundary Customer",
			CustomerEmail: "boundary@example.com",
			CustomerPhone: "555-0100",
			OrderType:     "pickup",
			PaymentMethod: "creditCard",
			Items: []orderapp.NewOrderItem{
				{
					MenuItemID: mi.ID.String(),
					Quantity:   qty,
				},
			},
		}
	}

	createTable := []apitest.Table{
		{
			Name:       "order-total-at-boundary-succeeds",
			URL:        "/v1/orders",
			Token:      adminToken,
			Method:     http.MethodPost,
			StatusCode: http.StatusCreated,
			Input:      newOrder(1),
			GotResp:    &orderapp.Order{},
			ExpResp:    &orderapp.Order{},
			CmpFunc: func(got any, exp any) string {
				gotResp, ok := got.(*orderapp.Order)
				if !ok {
					return "error occurred: could not cast response"
				}
				// 90,000,000.00 subtotal + 10% tax = 99,000,000.00 total.
				if gotResp.Total != 99_000_000.00 {
					return "order total: got " + money.MustParse(gotResp.Total).String() + ", want 99000000.00"
				}
				return ""
			},
		},
		{
			Name:       "order-total-above-cap-rejected",
			URL:        "/v1/orders",
			Token:      adminToken,
			Method:     http.MethodPost,
			StatusCode: http.StatusBadRequest,
			Input:      newOrder(2),
			GotResp:    &errs.Error{},
			ExpResp:    errs.Newf(errs.InvalidArgument, "order total exceeds the maximum supported amount"),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
	}
	test.Run(t, createTable, "moneyboundary-create")

	// -------------------------------------------------------------------------
	// Second under-cap order so the aggregate gross sales reach
	// 198,000,000.00, beyond the maximum supported amount.

	secondOrderTable := []apitest.Table{
		{
			Name:       "second-order-total-at-boundary-succeeds",
			URL:        "/v1/orders",
			Token:      adminToken,
			Method:     http.MethodPost,
			StatusCode: http.StatusCreated,
			Input:      newOrder(1),
			GotResp:    &orderapp.Order{},
			ExpResp:    &orderapp.Order{},
			CmpFunc: func(got any, exp any) string {
				gotResp, ok := got.(*orderapp.Order)
				if !ok {
					return "error occurred: could not cast response"
				}
				if gotResp.Total != 99_000_000.00 {
					return "order total: got " + money.MustParse(gotResp.Total).String() + ", want 99000000.00"
				}
				return ""
			},
		},
	}
	test.Run(t, secondOrderTable, "moneyboundary-second-order")

	// -------------------------------------------------------------------------
	// Insights: the aggregated gross sales now exceed the cap and must
	// return a typed 400 instead of clamping or panicking.

	insightsTable := []apitest.Table{
		{
			Name:       "insights-aggregate-above-cap-rejected",
			URL:        "/v1/insights",
			Token:      adminToken,
			Method:     http.MethodGet,
			StatusCode: http.StatusBadRequest,
			GotResp:    &errs.Error{},
			ExpResp:    errs.Newf(errs.InvalidArgument, "aggregated sales total exceeds the maximum supported amount"),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
	}
	test.Run(t, insightsTable, "moneyboundary-insights")
}
