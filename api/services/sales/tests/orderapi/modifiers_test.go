package orderapi_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/warlck/food-flow/app/domain/orderapp"
	"github.com/warlck/food-flow/app/sdk/apitest"
	"github.com/warlck/food-flow/app/sdk/errs"
	"github.com/warlck/food-flow/business/domain/categorybus"
	"github.com/warlck/food-flow/business/domain/menuitembus"
	"github.com/warlck/food-flow/business/domain/modifiergroupbus"
	"github.com/warlck/food-flow/business/domain/modifieroptionbus"
	"github.com/warlck/food-flow/business/domain/organizationbus"
	"github.com/warlck/food-flow/business/domain/restaurantbus"
	"github.com/warlck/food-flow/business/domain/userbus"
	"github.com/warlck/food-flow/business/types/money"
	"github.com/warlck/food-flow/business/types/name"
	"github.com/warlck/food-flow/business/types/role"
)

// Test_OrderModifierValidation verifies at the HTTP layer that order
// modifier validation rejections surface as 400s carrying their domain
// message (review finding 24), and that a suspended required group does not
// block ordering end-to-end (review finding 5).
func Test_OrderModifierValidation(t *testing.T) {
	t.Parallel()

	test := apitest.New(t, "Test_OrderModifierValidation")

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

	newOrder := func(itemID string, mods []orderapp.NewOrderItemModifier) *orderapp.NewOrder {
		return &orderapp.NewOrder{
			RestaurantID:  restaurants[0].ID.String(),
			CustomerName:  "Modifier Customer",
			CustomerEmail: "modifier@example.com",
			CustomerPhone: "555-0200",
			OrderType:     "pickup",
			PaymentMethod: "creditCard",
			Items: []orderapp.NewOrderItem{
				{
					MenuItemID: itemID,
					Quantity:   1,
					Modifiers:  mods,
				},
			},
		}
	}

	// -------------------------------------------------------------------------
	// Fixture: an item whose required group is available; selecting nothing
	// must return the domain 400 message.

	reqItem, err := busDomain.MenuItem.Create(ctx, menuitembus.NewMenuItem{
		Name:         name.MustParse("Validated Burger"),
		Description:  "modifier validation fixture",
		Price:        money.MustParse(12.00),
		CategoryID:   categories[0].ID,
		RestaurantID: restaurants[0].ID,
	})
	if err != nil {
		t.Fatalf("seeding required-group item: %s", err)
	}

	reqGroup, err := busDomain.ModifierGroup.Create(ctx, modifiergroupbus.NewModifierGroup{
		MenuItemID:    reqItem.ID,
		RestaurantID:  restaurants[0].ID,
		Name:          name.MustParse("Required Spice Level"),
		MinSelections: 1,
		MaxSelections: 1,
		Available:     true,
	})
	if err != nil {
		t.Fatalf("seeding required group: %s", err)
	}

	_, err = busDomain.ModifierOption.Create(ctx, modifieroptionbus.NewModifierOption{
		ModifierGroupID: reqGroup.ID,
		RestaurantID:    restaurants[0].ID,
		Name:            name.MustParse("Mild Option"),
		PriceDelta:      money.MustParse(0),
	})
	if err != nil {
		t.Fatalf("seeding required group option: %s", err)
	}

	requiredTable := []apitest.Table{
		{
			Name:       "required-group-selection-missing",
			URL:        "/v1/orders",
			Token:      apitest.Token(busDomain, test.Auth, admins[0].Email.Address),
			Method:     http.MethodPost,
			StatusCode: http.StatusBadRequest,
			Input:      newOrder(reqItem.ID.String(), nil),
			GotResp:    &errs.Error{},
			ExpResp: errs.Newf(errs.InvalidArgument, "modifier selection required: group %s (%s)",
				reqGroup.ID.String(), reqGroup.Name.String()),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
	}
	test.Run(t, requiredTable, "modifiervalidation-required")

	// -------------------------------------------------------------------------
	// Fixture: an item whose required group is suspended. No selection must
	// succeed (201); a submitted selection must 400 with the group message.

	suspItem, err := busDomain.MenuItem.Create(ctx, menuitembus.NewMenuItem{
		Name:         name.MustParse("Suspended Burger"),
		Description:  "suspended group fixture",
		Price:        money.MustParse(12.00),
		CategoryID:   categories[0].ID,
		RestaurantID: restaurants[0].ID,
	})
	if err != nil {
		t.Fatalf("seeding suspended-group item: %s", err)
	}

	suspGroup, err := busDomain.ModifierGroup.Create(ctx, modifiergroupbus.NewModifierGroup{
		MenuItemID:    suspItem.ID,
		RestaurantID:  restaurants[0].ID,
		Name:          name.MustParse("Suspended Size Group"),
		MinSelections: 1,
		MaxSelections: 1,
		Available:     false,
	})
	if err != nil {
		t.Fatalf("seeding suspended group: %s", err)
	}

	suspOption, err := busDomain.ModifierOption.Create(ctx, modifieroptionbus.NewModifierOption{
		ModifierGroupID: suspGroup.ID,
		RestaurantID:    restaurants[0].ID,
		Name:            name.MustParse("Suspended Large Option"),
		PriceDelta:      money.MustParse(2.00),
	})
	if err != nil {
		t.Fatalf("seeding suspended group option: %s", err)
	}

	suspendedTable := []apitest.Table{
		{
			Name:       "suspended-required-group-no-selection-succeeds",
			URL:        "/v1/orders",
			Token:      apitest.Token(busDomain, test.Auth, admins[0].Email.Address),
			Method:     http.MethodPost,
			StatusCode: http.StatusCreated,
			Input:      newOrder(suspItem.ID.String(), nil),
			GotResp:    &orderapp.Order{},
			ExpResp:    &orderapp.Order{},
			CmpFunc: func(got any, exp any) string {
				gotResp, ok := got.(*orderapp.Order)
				if !ok {
					return "error occurred: could not cast response"
				}
				// Base 12.00 with the 10% seeded tax rate: tax 1.20, total 13.20.
				if gotResp.Subtotal != 12.00 {
					return fmt.Sprintf("subtotal: got %.2f, want 12.00", gotResp.Subtotal)
				}
				if gotResp.Tax != 1.20 {
					return fmt.Sprintf("tax: got %.2f, want 1.20", gotResp.Tax)
				}
				if gotResp.Total != 13.20 {
					return fmt.Sprintf("total: got %.2f, want 13.20", gotResp.Total)
				}
				return ""
			},
		},
		{
			Name:       "suspended-group-submitted-selection-rejected",
			URL:        "/v1/orders",
			Token:      apitest.Token(busDomain, test.Auth, admins[0].Email.Address),
			Method:     http.MethodPost,
			StatusCode: http.StatusBadRequest,
			Input: newOrder(suspItem.ID.String(), []orderapp.NewOrderItemModifier{
				{
					ModifierGroupID:  suspGroup.ID.String(),
					ModifierOptionID: suspOption.ID.String(),
				},
			}),
			GotResp: &errs.Error{},
			ExpResp: errs.Newf(errs.InvalidArgument, "modifier group is unavailable: group %s (%s) is unavailable",
				suspGroup.ID.String(), suspGroup.Name.String()),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
	}
	test.Run(t, suspendedTable, "modifiervalidation-suspended")

	// -------------------------------------------------------------------------
	// Fixture: an available group whose option is unavailable; selecting it
	// must 400 with the option message.

	unavailItem, err := busDomain.MenuItem.Create(ctx, menuitembus.NewMenuItem{
		Name:         name.MustParse("Unavailable Option Burger"),
		Description:  "unavailable option fixture",
		Price:        money.MustParse(12.00),
		CategoryID:   categories[0].ID,
		RestaurantID: restaurants[0].ID,
	})
	if err != nil {
		t.Fatalf("seeding unavailable-option item: %s", err)
	}

	unavailGroup, err := busDomain.ModifierGroup.Create(ctx, modifiergroupbus.NewModifierGroup{
		MenuItemID:    unavailItem.ID,
		RestaurantID:  restaurants[0].ID,
		Name:          name.MustParse("Sauce Selection"),
		MinSelections: 0,
		MaxSelections: 1,
		Available:     true,
	})
	if err != nil {
		t.Fatalf("seeding sauce group: %s", err)
	}

	availFalse := false
	unavailOption, err := busDomain.ModifierOption.Create(ctx, modifieroptionbus.NewModifierOption{
		ModifierGroupID: unavailGroup.ID,
		RestaurantID:    restaurants[0].ID,
		Name:            name.MustParse("Out of Stock Sauce"),
		PriceDelta:      money.MustParse(1.00),
		Available:       &availFalse,
	})
	if err != nil {
		t.Fatalf("seeding unavailable option: %s", err)
	}

	unavailableTable := []apitest.Table{
		{
			Name:       "unavailable-option-selection-rejected",
			URL:        "/v1/orders",
			Token:      apitest.Token(busDomain, test.Auth, admins[0].Email.Address),
			Method:     http.MethodPost,
			StatusCode: http.StatusBadRequest,
			Input: newOrder(unavailItem.ID.String(), []orderapp.NewOrderItemModifier{
				{
					ModifierGroupID:  unavailGroup.ID.String(),
					ModifierOptionID: unavailOption.ID.String(),
				},
			}),
			GotResp: &errs.Error{},
			ExpResp: errs.Newf(errs.InvalidArgument, "modifier option is unavailable: option %s (%s)",
				unavailOption.ID.String(), unavailOption.Name.String()),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
	}
	test.Run(t, unavailableTable, "modifiervalidation-unavailable-option")
}
