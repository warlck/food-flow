package orderbus_test

import (
	"github.com/warlck/food-flow/business/domain/organizationbus"
	"context"
	"errors"
	"fmt"
	"math"
	"sort"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/warlck/food-flow/business/domain/addonbus"
	"github.com/warlck/food-flow/business/domain/categorybus"
	"github.com/warlck/food-flow/business/domain/menuitembus"
	"github.com/warlck/food-flow/business/domain/orderbus"
	"github.com/warlck/food-flow/business/domain/promobus"
	"github.com/warlck/food-flow/business/domain/restaurantbus"
	"github.com/warlck/food-flow/business/sdk/dbtest"
	"github.com/warlck/food-flow/business/sdk/page"
	"github.com/warlck/food-flow/business/sdk/unittest"
	"github.com/warlck/food-flow/business/types/money"
	"github.com/warlck/food-flow/business/types/name"
)

var equateTime = cmp.Comparer(func(x, y time.Time) bool {
	if x.IsZero() || y.IsZero() {
		return x.IsZero() == y.IsZero()
	}
	return x.Truncate(time.Microsecond).Equal(y.Truncate(time.Microsecond))
})

func Test_Order(t *testing.T) {
	t.Parallel()

	db := dbtest.New(t, "Test_Order")

	sd, err := insertSeedData(db.BusDomain)
	if err != nil {
		t.Fatalf("Seeding error: %s", err)
	}

	// -------------------------------------------------------------------------

	unittest.Run(t, query(db.BusDomain, sd), "query")
	unittest.Run(t, create(db.BusDomain, sd), "create")
	unittest.Run(t, updateStatus(db.BusDomain, sd), "update-status")
	unittest.Run(t, cancel(db.BusDomain, sd), "cancel")
}

// =============================================================================

func insertSeedData(busDomain dbtest.BusDomain) (unittest.SeedData, error) {
	ctx := context.Background()

	// Seed restaurants
	orgs, err := organizationbus.TestSeedOrganizations(ctx, 1, busDomain.Organization)
	if err != nil {
		return unittest.SeedData{}, fmt.Errorf("seeding organizations: %w", err)
	}
	rests, err := restaurantbus.TestSeedRestaurants(ctx, 2, busDomain.Restaurant, orgs[0].ID)
	if err != nil {
		return unittest.SeedData{}, fmt.Errorf("seeding restaurants : %w", err)
	}

	// Seed categories
	cats1, err := categorybus.TestSeedCategories(ctx, 2, rests[0].ID, busDomain.Category)
	if err != nil {
		return unittest.SeedData{}, fmt.Errorf("seeding categories : %w", err)
	}

	cats2, err := categorybus.TestSeedCategories(ctx, 2, rests[1].ID, busDomain.Category)
	if err != nil {
		return unittest.SeedData{}, fmt.Errorf("seeding categories : %w", err)
	}

	// Seed menu items
	items1, err := menuitembus.TestSeedMenuItems(ctx, 3, cats1[0].ID, rests[0].ID, busDomain.MenuItem)
	if err != nil {
		return unittest.SeedData{}, fmt.Errorf("seeding menu items : %w", err)
	}

	items2, err := menuitembus.TestSeedMenuItems(ctx, 3, cats2[0].ID, rests[1].ID, busDomain.MenuItem)
	if err != nil {
		return unittest.SeedData{}, fmt.Errorf("seeding menu items : %w", err)
	}

	// Seed addons
	addons1, err := addonbus.TestSeedAddons(ctx, 2, cats1[0].ID, rests[0].ID, busDomain.Addon)
	if err != nil {
		return unittest.SeedData{}, fmt.Errorf("seeding addons : %w", err)
	}

	addons1OtherCat, err := addonbus.TestSeedAddons(ctx, 1, cats1[1].ID, rests[0].ID, busDomain.Addon)
	if err != nil {
		return unittest.SeedData{}, fmt.Errorf("seeding addons : %w", err)
	}

	addons2, err := addonbus.TestSeedAddons(ctx, 1, cats2[0].ID, rests[1].ID, busDomain.Addon)
	if err != nil {
		return unittest.SeedData{}, fmt.Errorf("seeding addons : %w", err)
	}

	// Seed orders
	orders1, err := orderbus.TestSeedOrders(ctx, 2, rests[0].ID, []menuitembus.MenuItem{items1[0], items1[1]}, busDomain.Order)
	if err != nil {
		return unittest.SeedData{}, fmt.Errorf("seeding orders : %w", err)
	}

	orders2, err := orderbus.TestSeedOrders(ctx, 2, rests[1].ID, []menuitembus.MenuItem{items2[0], items2[1]}, busDomain.Order)
	if err != nil {
		return unittest.SeedData{}, fmt.Errorf("seeding orders : %w", err)
	}

	// -------------------------------------------------------------------------

	sd := unittest.SeedData{
		Restaurants: []unittest.Restaurant{
			{Restaurant: rests[0]},
			{Restaurant: rests[1]},
		},
		Categories: []unittest.Category{
			{Category: cats1[0]},
			{Category: cats1[1]},
			{Category: cats2[0]},
			{Category: cats2[1]},
		},
		MenuItems: []unittest.MenuItem{
			{MenuItem: items1[0]},
			{MenuItem: items1[1]},
			{MenuItem: items1[2]},
			{MenuItem: items2[0]},
			{MenuItem: items2[1]},
			{MenuItem: items2[2]},
		},
		Orders: []unittest.Order{
			{Order: orders1[0]},
			{Order: orders1[1]},
			{Order: orders2[0]},
			{Order: orders2[1]},
		},
		Addons: []unittest.Addon{
			{Addon: addons1[0]},
			{Addon: addons1[1]},
			{Addon: addons1OtherCat[0]},
			{Addon: addons2[0]},
		},
	}

	return sd, nil
}

// =============================================================================

func query(busDomain dbtest.BusDomain, sd unittest.SeedData) []unittest.Table {
	orders := make([]orderbus.Order, 0, len(sd.Orders))

	for _, order := range sd.Orders {
		orders = append(orders, order.Order)
	}

	// Sort by DateCreated DESC to match DefaultOrderBy
	sort.Slice(orders, func(i, j int) bool {
		return orders[i].DateCreated.After(orders[j].DateCreated)
	})

	table := []unittest.Table{
		{
			Name:    "all",
			ExpResp: orders,
			ExcFunc: func(ctx context.Context) any {
				resp, err := busDomain.Order.Query(ctx, orderbus.QueryFilter{}, orderbus.DefaultOrderBy, page.MustParse("1", "10"))
				if err != nil {
					return err
				}

				return resp
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, ok := got.([]orderbus.Order)
				if !ok {
					return "error occurred"
				}

				expResp := exp.([]orderbus.Order)

				for i := range gotResp {
					// Normalize item OrderIDs from database
					for j := range gotResp[i].Items {
						if j < len(expResp[i].Items) {
							expResp[i].Items[j].OrderID = gotResp[i].Items[j].OrderID
						}
					}

					// Normalize address OrderID if present
					if gotResp[i].DeliveryAddress != nil && expResp[i].DeliveryAddress != nil {
						expResp[i].DeliveryAddress.OrderID = gotResp[i].DeliveryAddress.OrderID
					}
				}

				return cmp.Diff(gotResp, expResp, equateTime)
			},
		},
		{
			Name:    "byid",
			ExpResp: sd.Orders[0].Order,
			ExcFunc: func(ctx context.Context) any {
				resp, err := busDomain.Order.QueryByID(ctx, sd.Orders[0].ID)
				if err != nil {
					return err
				}

				return resp
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.(orderbus.Order)
				if !exists {
					return "error occurred"
				}

				expResp := exp.(orderbus.Order)

				// Normalize monetary values from database (rounded to 2 decimals)
				expResp.Subtotal = gotResp.Subtotal
				expResp.Tax = gotResp.Tax
				expResp.Total = gotResp.Total

				// Normalize item OrderIDs from database
				for i := range gotResp.Items {
					if i < len(expResp.Items) {
						expResp.Items[i].OrderID = gotResp.Items[i].OrderID
					}
				}

				// Normalize address OrderID if present
				if gotResp.DeliveryAddress != nil && expResp.DeliveryAddress != nil {
					expResp.DeliveryAddress.OrderID = gotResp.DeliveryAddress.OrderID
				}

				return cmp.Diff(gotResp, expResp, equateTime)
			},
		},
	}

	return table
}

func create(busDomain dbtest.BusDomain, sd unittest.SeedData) []unittest.Table {

	// Expected totals for the "order-with-addons" case: menu item 0 (qty 2)
	// with addon 0 (qty 1) and addon 1 (qty 2), plus menu item 1 (qty 1)
	// without addons. Addon cost is addon price * addon qty * item qty.
	expAddonSubtotal := sd.MenuItems[0].Price.Value()*2 +
		sd.Addons[0].Price.Value()*1*2 +
		sd.Addons[1].Price.Value()*2*2 +
		sd.MenuItems[1].Price.Value()*1
	expAddonSubtotal = math.Round(expAddonSubtotal*100) / 100
	expAddonTax := math.Round(expAddonSubtotal*sd.Restaurants[0].TaxRate*100) / 100
	expAddonTotal := math.Round((expAddonSubtotal+expAddonTax)*100) / 100

	table := []unittest.Table{
		{
			Name: "pickup-order",
			ExpResp: orderbus.Order{
				RestaurantID:    sd.Restaurants[0].ID,
				CustomerName:    "John Doe",
				CustomerEmail:   "john@example.com",
				CustomerPhone:   "555-1234",
				OrderType:       orderbus.OrderTypePickup,
				OrderStatus:     orderbus.OrderStatusPending,
				PaymentStatus:   orderbus.PaymentStatusPending,
				PaymentMethod:   orderbus.PaymentMethodCreditCard,
				DeliveryFee:     money.MustParse(0),
				DeliveryAddress: nil,
			},
			ExcFunc: func(ctx context.Context) any {
				no := orderbus.NewOrder{
					RestaurantID:  sd.Restaurants[0].ID.String(),
					CustomerName:  "John Doe",
					CustomerEmail: "john@example.com",
					CustomerPhone: "555-1234",
					OrderType:     orderbus.OrderTypePickup,
					PaymentMethod: orderbus.PaymentMethodCreditCard,
					Items: []orderbus.NewOrderItem{
						{
							MenuItemID: sd.MenuItems[0].ID.String(),
							Quantity:   2,
						},
						{
							MenuItemID: sd.MenuItems[1].ID.String(),
							Quantity:   1,
						},
					},
				}

				resp, err := busDomain.Order.Create(ctx, no)
				if err != nil {
					return err
				}

				return resp
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.(orderbus.Order)
				if !exists {
					return "error occurred"
				}

				expResp := exp.(orderbus.Order)

				expResp.ID = gotResp.ID
				expResp.Subtotal = gotResp.Subtotal
				expResp.Tax = gotResp.Tax
				expResp.Total = gotResp.Total
				expResp.Items = gotResp.Items
				expResp.DateCreated = gotResp.DateCreated
				expResp.DateUpdated = gotResp.DateUpdated

				return cmp.Diff(gotResp, expResp)
			},
		},
		{
			Name: "pickup-order-with-promo",
			ExpResp: orderbus.Order{
				RestaurantID:  sd.Restaurants[0].ID,
				CustomerName:  "Promo User",
				CustomerEmail: "promo@example.com",
				CustomerPhone: "555-9999",
				OrderType:     orderbus.OrderTypePickup,
				OrderStatus:   orderbus.OrderStatusPending,
				PaymentStatus: orderbus.PaymentStatusPending,
				PaymentMethod: orderbus.PaymentMethodCreditCard,
				PromoCode:     "SAVE10PERCENT",
				DeliveryFee:   money.MustParse(0),
			},
			ExcFunc: func(ctx context.Context) any {
				// Seed a promo code first
				_, err := busDomain.Promo.Create(ctx, promobus.NewPromotion{
					Code:          "SAVE10PERCENT",
					Name:          name.MustParse("Save 10 Percent"),
					DiscountType:  "percentage",
					DiscountValue: 10,
					Enabled:       true,
				})
				if err != nil {
					return err
				}

				no := orderbus.NewOrder{
					RestaurantID:  sd.Restaurants[0].ID.String(),
					CustomerName:  "Promo User",
					CustomerEmail: "promo@example.com",
					CustomerPhone: "555-9999",
					OrderType:     orderbus.OrderTypePickup,
					PaymentMethod: orderbus.PaymentMethodCreditCard,
					PromoCode:     "SAVE10PERCENT",
					Items: []orderbus.NewOrderItem{
						{
							MenuItemID: sd.MenuItems[0].ID.String(),
							Quantity:   2,
						},
					},
				}

				resp, err := busDomain.Order.Create(ctx, no)
				if err != nil {
					return err
				}

				return resp
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.(orderbus.Order)
				if !exists {
					return "error occurred"
				}

				expResp := exp.(orderbus.Order)

				if gotResp.PromoCode != "SAVE10PERCENT" {
					return fmt.Sprintf("expected PromoCode SAVE10PERCENT, got %s", gotResp.PromoCode)
				}
				if gotResp.Discount.Value() <= 0 {
					return fmt.Sprintf("expected positive Discount, got %.2f", gotResp.Discount.Value())
				}

				expSubtotal := sd.MenuItems[0].Price.Value() * 2
				expDiscount := math.Round(expSubtotal*0.10*100) / 100
				expTaxable := expSubtotal - expDiscount
				expTax := math.Round(expTaxable*sd.Restaurants[0].TaxRate*100) / 100
				expTotal := expTaxable + expTax

				if math.Abs(gotResp.Discount.Value()-expDiscount) > 0.01 {
					return fmt.Sprintf("discount mismatch: got %.2f, want %.2f", gotResp.Discount.Value(), expDiscount)
				}
				if math.Abs(gotResp.Tax.Value()-expTax) > 0.01 {
					return fmt.Sprintf("tax mismatch: got %.2f, want %.2f", gotResp.Tax.Value(), expTax)
				}
				if math.Abs(gotResp.Total.Value()-expTotal) > 0.01 {
					return fmt.Sprintf("total mismatch: got %.2f, want %.2f", gotResp.Total.Value(), expTotal)
				}

				expResp.ID = gotResp.ID
				expResp.Subtotal = gotResp.Subtotal
				expResp.Discount = gotResp.Discount
				expResp.Tax = gotResp.Tax
				expResp.Total = gotResp.Total
				expResp.Items = gotResp.Items
				expResp.DateCreated = gotResp.DateCreated
				expResp.DateUpdated = gotResp.DateUpdated

				return cmp.Diff(gotResp, expResp)
			},
		},
		{
			Name:    "pickup-order-invalid-promo",
			ExpResp: errors.New("invalid promo code: Invalid promo code"),
			ExcFunc: func(ctx context.Context) any {
				no := orderbus.NewOrder{
					RestaurantID:  sd.Restaurants[0].ID.String(),
					CustomerName:  "Bad Promo",
					CustomerEmail: "badpromo@example.com",
					CustomerPhone: "555-0000",
					OrderType:     orderbus.OrderTypePickup,
					PaymentMethod: orderbus.PaymentMethodCreditCard,
					PromoCode:     "NOTREAL99",
					Items: []orderbus.NewOrderItem{
						{
							MenuItemID: sd.MenuItems[0].ID.String(),
							Quantity:   1,
						},
					},
				}

				_, err := busDomain.Order.Create(ctx, no)
				return err
			},
			CmpFunc: func(got any, exp any) string {
				gotErr, exists := got.(error)
				if !exists {
					return "error expected"
				}

				expErr := exp.(error)
				if gotErr.Error() != expErr.Error() {
					return fmt.Sprintf("got error %q, want %q", gotErr.Error(), expErr.Error())
				}
				return ""
			},
		},
		{
			Name:    "pickup-order-disabled-promo",
			ExpResp: errors.New("invalid promo code: Promo code is inactive"),
			ExcFunc: func(ctx context.Context) any {
				_, err := busDomain.Promo.Create(ctx, promobus.NewPromotion{
					Code:           "DISABLEDPROMO",
					Name:           name.MustParse("Disabled Promo"),
					DiscountType:   "percentage",
					DiscountValue:  10,
					MinOrderAmount: 0,
					Enabled:        false,
				})
				if err != nil {
					return err
				}

				no := orderbus.NewOrder{
					RestaurantID:  sd.Restaurants[0].ID.String(),
					CustomerName:  "Disabled Promo User",
					CustomerEmail: "disabledpromo@example.com",
					CustomerPhone: "555-0001",
					OrderType:     orderbus.OrderTypePickup,
					PaymentMethod: orderbus.PaymentMethodCreditCard,
					PromoCode:     "DISABLEDPROMO",
					Items: []orderbus.NewOrderItem{
						{
							MenuItemID: sd.MenuItems[0].ID.String(),
							Quantity:   1,
						},
					},
				}

				_, err = busDomain.Order.Create(ctx, no)
				return err
			},
			CmpFunc: func(got any, exp any) string {
				gotErr, exists := got.(error)
				if !exists {
					return "error expected"
				}

				expErr := exp.(error)
				if gotErr.Error() != expErr.Error() {
					return fmt.Sprintf("got error %q, want %q", gotErr.Error(), expErr.Error())
				}
				return ""
			},
		},
		{
			Name:    "pickup-order-subtotal-below-min-promo",
			ExpResp: errors.New("invalid promo code: Minimum order subtotal of $500.00 required for this promo code"),
			ExcFunc: func(ctx context.Context) any {
				_, err := busDomain.Promo.Create(ctx, promobus.NewPromotion{
					Code:           "HIGHMINPROMO",
					Name:           name.MustParse("High Min Promo"),
					DiscountType:   "percentage",
					DiscountValue:  10,
					MinOrderAmount: 500.0,
					Enabled:        true,
				})
				if err != nil {
					return err
				}

				no := orderbus.NewOrder{
					RestaurantID:  sd.Restaurants[0].ID.String(),
					CustomerName:  "Below Min User",
					CustomerEmail: "belowmin@example.com",
					CustomerPhone: "555-0002",
					OrderType:     orderbus.OrderTypePickup,
					PaymentMethod: orderbus.PaymentMethodCreditCard,
					PromoCode:     "HIGHMINPROMO",
					Items: []orderbus.NewOrderItem{
						{
							MenuItemID: sd.MenuItems[0].ID.String(),
							Quantity:   1,
						},
					},
				}

				_, err = busDomain.Order.Create(ctx, no)
				return err
			},
			CmpFunc: func(got any, exp any) string {
				gotErr, exists := got.(error)
				if !exists {
					return "error expected"
				}

				expErr := exp.(error)
				if gotErr.Error() != expErr.Error() {
					return fmt.Sprintf("got error %q, want %q", gotErr.Error(), expErr.Error())
				}
				return ""
			},
		},
		{
			Name:    "pickup-order-expired-promo",
			ExpResp: errors.New("invalid promo code: Promo code has expired"),
			ExcFunc: func(ctx context.Context) any {
				past := time.Now().Add(-24 * time.Hour)
				_, err := busDomain.Promo.Create(ctx, promobus.NewPromotion{
					Code:           "EXPIREDORDERPROMO",
					Name:           name.MustParse("Expired Order Promo"),
					DiscountType:   "percentage",
					DiscountValue:  10,
					MinOrderAmount: 0,
					EndDate:        &past,
					Enabled:        true,
				})
				if err != nil {
					return err
				}

				no := orderbus.NewOrder{
					RestaurantID:  sd.Restaurants[0].ID.String(),
					CustomerName:  "Expired Promo User",
					CustomerEmail: "expired@example.com",
					CustomerPhone: "555-0003",
					OrderType:     orderbus.OrderTypePickup,
					PaymentMethod: orderbus.PaymentMethodCreditCard,
					PromoCode:     "EXPIREDORDERPROMO",
					Items: []orderbus.NewOrderItem{
						{
							MenuItemID: sd.MenuItems[0].ID.String(),
							Quantity:   1,
						},
					},
				}

				_, err = busDomain.Order.Create(ctx, no)
				return err
			},
			CmpFunc: func(got any, exp any) string {
				gotErr, exists := got.(error)
				if !exists {
					return "error expected"
				}

				expErr := exp.(error)
				if gotErr.Error() != expErr.Error() {
					return fmt.Sprintf("got error %q, want %q", gotErr.Error(), expErr.Error())
				}
				return ""
			},
		},
		{
			Name:    "pickup-order-restaurant-mismatch-promo",
			ExpResp: errors.New("invalid promo code: Promo code is not applicable to this restaurant"),
			ExcFunc: func(ctx context.Context) any {
				rest1ID := sd.Restaurants[1].ID
				_, err := busDomain.Promo.Create(ctx, promobus.NewPromotion{
					Code:           "REST1ONLYPROMO",
					Name:           name.MustParse("Rest 1 Promo"),
					DiscountType:   "percentage",
					DiscountValue:  10,
					MinOrderAmount: 0,
					RestaurantID:   &rest1ID,
					Enabled:        true,
				})
				if err != nil {
					return err
				}

				// Attempt order at Restaurant 0 with Restaurant 1's promo
				no := orderbus.NewOrder{
					RestaurantID:  sd.Restaurants[0].ID.String(),
					CustomerName:  "Rest Mismatch User",
					CustomerEmail: "mismatch@example.com",
					CustomerPhone: "555-0004",
					OrderType:     orderbus.OrderTypePickup,
					PaymentMethod: orderbus.PaymentMethodCreditCard,
					PromoCode:     "REST1ONLYPROMO",
					Items: []orderbus.NewOrderItem{
						{
							MenuItemID: sd.MenuItems[0].ID.String(),
							Quantity:   1,
						},
					},
				}

				_, err = busDomain.Order.Create(ctx, no)
				return err
			},
			CmpFunc: func(got any, exp any) string {
				gotErr, exists := got.(error)
				if !exists {
					return "error expected"
				}

				expErr := exp.(error)
				if gotErr.Error() != expErr.Error() {
					return fmt.Sprintf("got error %q, want %q", gotErr.Error(), expErr.Error())
				}
				return ""
			},
		},
		{
			Name:    "pickup-order-usage-limit-exceeded-on-second-order",
			ExpResp: errors.New("invalid promo code: Promo code usage limit reached"),
			ExcFunc: func(ctx context.Context) any {
				limit1 := 1
				_, err := busDomain.Promo.Create(ctx, promobus.NewPromotion{
					Code:           "LIMIT1PROMO",
					Name:           name.MustParse("Limit 1 Promo"),
					DiscountType:   "percentage",
					DiscountValue:  10,
					MinOrderAmount: 0,
					UsageLimit:     &limit1,
					Enabled:        true,
				})
				if err != nil {
					return err
				}

				// First order succeeds and increments usage count to 1
				no1 := orderbus.NewOrder{
					RestaurantID:  sd.Restaurants[0].ID.String(),
					CustomerName:  "First Order User",
					CustomerEmail: "user1@example.com",
					CustomerPhone: "555-0010",
					OrderType:     orderbus.OrderTypePickup,
					PaymentMethod: orderbus.PaymentMethodCreditCard,
					PromoCode:     "LIMIT1PROMO",
					Items: []orderbus.NewOrderItem{
						{
							MenuItemID: sd.MenuItems[0].ID.String(),
							Quantity:   1,
						},
					},
				}
				if _, err := busDomain.Order.Create(ctx, no1); err != nil {
					return fmt.Errorf("first order failed: %w", err)
				}

				// Second order fails atomically on increment because limit has been reached
				no2 := orderbus.NewOrder{
					RestaurantID:  sd.Restaurants[0].ID.String(),
					CustomerName:  "Second Order User",
					CustomerEmail: "user2@example.com",
					CustomerPhone: "555-0011",
					OrderType:     orderbus.OrderTypePickup,
					PaymentMethod: orderbus.PaymentMethodCreditCard,
					PromoCode:     "LIMIT1PROMO",
					Items: []orderbus.NewOrderItem{
						{
							MenuItemID: sd.MenuItems[0].ID.String(),
							Quantity:   1,
						},
					},
				}

				_, err = busDomain.Order.Create(ctx, no2)
				return err
			},
			CmpFunc: func(got any, exp any) string {
				gotErr, exists := got.(error)
				if !exists {
					return "error expected"
				}

				expErr := exp.(error)
				if gotErr.Error() != expErr.Error() {
					return fmt.Sprintf("got error %q, want %q", gotErr.Error(), expErr.Error())
				}
				return ""
			},
		},
		{
			Name:    "min-spend-not-met",
			ExpResp: errors.New("subtotal does not meet restaurant minimum spend: subtotal 10.00 is less than minimum spend 100.00"),
			ExcFunc: func(ctx context.Context) any {
				// Update restaurant minimum spend to 100.00
				ms := 100.00
				_, err := busDomain.Restaurant.Update(ctx, sd.Restaurants[1].Restaurant, restaurantbus.UpdateRestaurant{
					MinSpend: &ms,
				})
				if err != nil {
					return fmt.Errorf("updating restaurant min spend: %w", err)
				}

				no := orderbus.NewOrder{
					RestaurantID:  sd.Restaurants[1].ID.String(),
					CustomerName:  "Min Spend Tester",
					CustomerEmail: "minspend@example.com",
					CustomerPhone: "555-0099",
					OrderType:     orderbus.OrderTypePickup,
					PaymentMethod: orderbus.PaymentMethodCreditCard,
					Items: []orderbus.NewOrderItem{
						{
							MenuItemID: sd.MenuItems[3].ID.String(),
							Quantity:   1, // item price is < 100.00
						},
					},
				}

				_, err = busDomain.Order.Create(ctx, no)
				return err
			},
			CmpFunc: func(got any, exp any) string {
				gotErr, exists := got.(error)
				if !exists {
					return "error expected"
				}

				if !errors.Is(gotErr, orderbus.ErrMinSpendNotMet) {
					return fmt.Sprintf("got error %v, want ErrMinSpendNotMet", gotErr)
				}

				expMsg := fmt.Sprintf("subtotal does not meet restaurant minimum spend: subtotal %.2f is less than minimum spend 100.00", sd.MenuItems[3].Price.Value())
				if gotErr.Error() != expMsg {
					return fmt.Sprintf("got error %q, want %q", gotErr.Error(), expMsg)
				}
				return ""
			},
		},
		{
			Name: "delivery-order",
			ExpResp: orderbus.Order{
				RestaurantID:  sd.Restaurants[0].ID,
				CustomerName:  "Jane Smith",
				CustomerEmail: "jane@example.com",
				CustomerPhone: "555-5678",
				OrderType:     orderbus.OrderTypeDelivery,
				OrderStatus:   orderbus.OrderStatusPending,
				PaymentStatus: orderbus.PaymentStatusPending,
				PaymentMethod: orderbus.PaymentMethodCreditCard,
				DeliveryFee: money.MustParse(orderbus.CalculateDeliveryFee(orderbus.DistanceKm(
					*sd.Restaurants[0].Latitude, *sd.Restaurants[0].Longitude, 1.30719, 103.87434,
				))),
			},
			ExcFunc: func(ctx context.Context) any {
				lat := 1.30719
				lng := 103.87434

				no := orderbus.NewOrder{
					RestaurantID:  sd.Restaurants[0].ID.String(),
					CustomerName:  "Jane Smith",
					CustomerEmail: "jane@example.com",
					CustomerPhone: "555-5678",
					OrderType:     orderbus.OrderTypeDelivery,
					PaymentMethod: orderbus.PaymentMethodCreditCard,
					Items: []orderbus.NewOrderItem{
						{
							MenuItemID: sd.MenuItems[0].ID.String(),
							Quantity:   1,
						},
					},
					DeliveryAddress: &orderbus.NewDeliveryAddress{
						Street:     "123 Main St",
						City:       "Anytown",
						State:      "CA",
						PostalCode: "12345",
						Latitude:   &lat,
						Longitude:  &lng,
					},
				}

				resp, err := busDomain.Order.Create(ctx, no)
				if err != nil {
					return err
				}

				return resp
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.(orderbus.Order)
				if !exists {
					return "error occurred"
				}

				expResp := exp.(orderbus.Order)

				expResp.ID = gotResp.ID
				expResp.Subtotal = gotResp.Subtotal
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
			Name:    "delivery-missing-coordinates",
			ExpResp: orderbus.ErrDeliveryCoordinatesRequired,
			ExcFunc: func(ctx context.Context) any {
				no := orderbus.NewOrder{
					RestaurantID:  sd.Restaurants[0].ID.String(),
					CustomerName:  "No Coords",
					CustomerEmail: "nocoords@example.com",
					CustomerPhone: "555-9999",
					OrderType:     orderbus.OrderTypeDelivery,
					PaymentMethod: orderbus.PaymentMethodCreditCard,
					Items: []orderbus.NewOrderItem{
						{
							MenuItemID: sd.MenuItems[0].ID.String(),
							Quantity:   1,
						},
					},
					DeliveryAddress: &orderbus.NewDeliveryAddress{
						Street:     "123 Main St",
						City:       "Anytown",
						State:      "CA",
						PostalCode: "12345",
					},
				}

				_, err := busDomain.Order.Create(ctx, no)
				return err
			},
			CmpFunc: func(got any, exp any) string {
				gotErr, exists := got.(error)
				if !exists {
					return "expected an error"
				}

				if !errors.Is(gotErr, exp.(error)) {
					return "different error"
				}

				return ""
			},
		},
		{
			Name:    "delivery-out-of-range",
			ExpResp: orderbus.ErrDeliveryOutOfRange,
			ExcFunc: func(ctx context.Context) any {
				// Roughly 150 km away from the seeded restaurant,
				// beyond its 10 km delivery limit.
				lat := 2.64305
				lng := 103.86020

				no := orderbus.NewOrder{
					RestaurantID:  sd.Restaurants[0].ID.String(),
					CustomerName:  "Far Away",
					CustomerEmail: "far@example.com",
					CustomerPhone: "555-8888",
					OrderType:     orderbus.OrderTypeDelivery,
					PaymentMethod: orderbus.PaymentMethodCreditCard,
					Items: []orderbus.NewOrderItem{
						{
							MenuItemID: sd.MenuItems[0].ID.String(),
							Quantity:   1,
						},
					},
					DeliveryAddress: &orderbus.NewDeliveryAddress{
						Street:     "1 Far Road",
						City:       "Faraway",
						State:      "CA",
						PostalCode: "99999",
						Latitude:   &lat,
						Longitude:  &lng,
					},
				}

				_, err := busDomain.Order.Create(ctx, no)
				return err
			},
			CmpFunc: func(got any, exp any) string {
				gotErr, exists := got.(error)
				if !exists {
					return "expected an error"
				}

				if !errors.Is(gotErr, exp.(error)) {
					return "different error"
				}

				return ""
			},
		},
		{
			Name:    "zero-tax-order",
			ExpResp: float64(0.0),
			ExcFunc: func(ctx context.Context) any {
				// Create a restaurant with TaxRate = 0.0
				nr := restaurantbus.NewRestaurant{
					OrganizationID:        sd.Restaurants[0].OrganizationID,
					Name:                  name.MustParse("Zero Tax Bistro"),
					Description:           "Tax free dining",
					Address:               "123 Free St",
					Phone:                 "+1-555-0000",
					Email:                 "zerotax@test.com",
					TaxRate:               0.0,
					MaxDeliveryDistanceKm: 10,
				}
				rest, err := busDomain.Restaurant.Create(ctx, nr)
				if err != nil {
					return err
				}

				cats, err := categorybus.TestSeedCategories(ctx, 1, rest.ID, busDomain.Category)
				if err != nil {
					return err
				}

				items, err := menuitembus.TestSeedMenuItems(ctx, 1, cats[0].ID, rest.ID, busDomain.MenuItem)
				if err != nil {
					return err
				}

				no := orderbus.NewOrder{
					RestaurantID:  rest.ID.String(),
					CustomerName:  "Zero Tax Customer",
					CustomerEmail: "zerotaxcust@example.com",
					CustomerPhone: "555-0000",
					OrderType:     orderbus.OrderTypePickup,
					PaymentMethod: orderbus.PaymentMethodCreditCard,
					Items: []orderbus.NewOrderItem{
						{
							MenuItemID: items[0].ID.String(),
							Quantity:   1,
						},
					},
				}

				order, err := busDomain.Order.Create(ctx, no)
				if err != nil {
					return err
				}

				return order.Tax.Value()
			},
			CmpFunc: func(got any, exp any) string {
				gotTax, ok := got.(float64)
				if !ok {
					return "expected float64 tax value"
				}
				if gotTax != 0.0 {
					return fmt.Sprintf("expected tax 0.0, got %.2f", gotTax)
				}
				return ""
			},
		},
		{
			Name:    "order-with-addons",
			ExpResp: nil,
			ExcFunc: func(ctx context.Context) any {
				no := orderbus.NewOrder{
					RestaurantID:  sd.Restaurants[0].ID.String(),
					CustomerName:  "Addon Customer",
					CustomerEmail: "addon@example.com",
					CustomerPhone: "555-4321",
					OrderType:     orderbus.OrderTypePickup,
					PaymentMethod: orderbus.PaymentMethodCreditCard,
					Items: []orderbus.NewOrderItem{
						{
							MenuItemID: sd.MenuItems[0].ID.String(),
							Quantity:   2,
							Addons: []orderbus.NewOrderItemAddon{
								{AddonID: sd.Addons[0].ID.String(), Quantity: 1},
								{AddonID: sd.Addons[1].ID.String(), Quantity: 2},
							},
						},
						{
							MenuItemID: sd.MenuItems[1].ID.String(),
							Quantity:   1,
						},
					},
				}

				order, err := busDomain.Order.Create(ctx, no)
				if err != nil {
					return err
				}

				// Query the order back to verify addons were persisted.
				resp, err := busDomain.Order.QueryByID(ctx, order.ID)
				if err != nil {
					return err
				}

				return resp
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, ok := got.(orderbus.Order)
				if !ok {
					return "error occurred"
				}

				if gotResp.Subtotal.Value() != expAddonSubtotal {
					return fmt.Sprintf("subtotal: got %.2f, want %.2f", gotResp.Subtotal.Value(), expAddonSubtotal)
				}

				if gotResp.Tax.Value() != expAddonTax {
					return fmt.Sprintf("tax: got %.2f, want %.2f", gotResp.Tax.Value(), expAddonTax)
				}

				if gotResp.Total.Value() != expAddonTotal {
					return fmt.Sprintf("total: got %.2f, want %.2f", gotResp.Total.Value(), expAddonTotal)
				}

				// Find the item with addons by menu item ID since item
				// order is not guaranteed (items share a timestamp).
				var addonItem *orderbus.OrderItem
				for i := range gotResp.Items {
					if gotResp.Items[i].MenuItemID == sd.MenuItems[0].ID {
						addonItem = &gotResp.Items[i]
					}
				}

				if addonItem == nil {
					return "order item with addons not found"
				}

				if len(addonItem.Addons) != 2 {
					return fmt.Sprintf("expected 2 addons, got %d", len(addonItem.Addons))
				}

				expAddons := map[uuid.UUID]struct {
					name     string
					price    float64
					quantity int
				}{
					sd.Addons[0].ID: {sd.Addons[0].Name.String(), sd.Addons[0].Price.Value(), 1},
					sd.Addons[1].ID: {sd.Addons[1].Name.String(), sd.Addons[1].Price.Value(), 2},
				}

				for _, a := range addonItem.Addons {
					expA, ok := expAddons[a.AddonID]
					if !ok {
						return fmt.Sprintf("unexpected addon %s", a.AddonID)
					}

					if a.AddonName != expA.name {
						return fmt.Sprintf("addon name: got %q, want %q", a.AddonName, expA.name)
					}

					if a.AddonPrice.Value() != expA.price {
						return fmt.Sprintf("addon price: got %.2f, want %.2f", a.AddonPrice.Value(), expA.price)
					}

					if a.Quantity != expA.quantity {
						return fmt.Sprintf("addon quantity: got %d, want %d", a.Quantity, expA.quantity)
					}

					if a.OrderItemID != addonItem.ID {
						return "addon not linked to its order item"
					}
				}

				return ""
			},
		},
		{
			Name:    "addon-not-found",
			ExpResp: addonbus.ErrNotFound,
			ExcFunc: func(ctx context.Context) any {
				no := orderbus.NewOrder{
					RestaurantID:  sd.Restaurants[0].ID.String(),
					CustomerName:  "Ghost Addon",
					CustomerEmail: "ghost@example.com",
					CustomerPhone: "555-1111",
					OrderType:     orderbus.OrderTypePickup,
					PaymentMethod: orderbus.PaymentMethodCreditCard,
					Items: []orderbus.NewOrderItem{
						{
							MenuItemID: sd.MenuItems[0].ID.String(),
							Quantity:   1,
							Addons: []orderbus.NewOrderItemAddon{
								{AddonID: uuid.New().String(), Quantity: 1},
							},
						},
					},
				}

				_, err := busDomain.Order.Create(ctx, no)
				return err
			},
			CmpFunc: func(got any, exp any) string {
				gotErr, exists := got.(error)
				if !exists {
					return "expected an error"
				}

				if !errors.Is(gotErr, exp.(error)) {
					return "different error"
				}

				return ""
			},
		},
		{
			Name:    "addon-quantity-zero",
			ExpResp: orderbus.ErrAddonQuantityOutOfRange,
			ExcFunc: func(ctx context.Context) any {
				no := orderbus.NewOrder{
					RestaurantID:  sd.Restaurants[0].ID.String(),
					CustomerName:  "Zero Qty",
					CustomerEmail: "zeroqty@example.com",
					CustomerPhone: "555-2222",
					OrderType:     orderbus.OrderTypePickup,
					PaymentMethod: orderbus.PaymentMethodCreditCard,
					Items: []orderbus.NewOrderItem{
						{
							MenuItemID: sd.MenuItems[0].ID.String(),
							Quantity:   1,
							Addons: []orderbus.NewOrderItemAddon{
								{AddonID: sd.Addons[0].ID.String(), Quantity: 0},
							},
						},
					},
				}

				_, err := busDomain.Order.Create(ctx, no)
				return err
			},
			CmpFunc: func(got any, exp any) string {
				gotErr, exists := got.(error)
				if !exists {
					return "expected an error"
				}

				if !errors.Is(gotErr, exp.(error)) {
					return "different error"
				}

				return ""
			},
		},
		{
			Name:    "addon-quantity-exceeds-max",
			ExpResp: orderbus.ErrAddonQuantityOutOfRange,
			ExcFunc: func(ctx context.Context) any {
				no := orderbus.NewOrder{
					RestaurantID:  sd.Restaurants[0].ID.String(),
					CustomerName:  "Max Qty",
					CustomerEmail: "maxqty@example.com",
					CustomerPhone: "555-3333",
					OrderType:     orderbus.OrderTypePickup,
					PaymentMethod: orderbus.PaymentMethodCreditCard,
					Items: []orderbus.NewOrderItem{
						{
							MenuItemID: sd.MenuItems[0].ID.String(),
							Quantity:   1,
							Addons: []orderbus.NewOrderItemAddon{
								{AddonID: sd.Addons[0].ID.String(), Quantity: sd.Addons[0].MaxQuantity + 1},
							},
						},
					},
				}

				_, err := busDomain.Order.Create(ctx, no)
				return err
			},
			CmpFunc: func(got any, exp any) string {
				gotErr, exists := got.(error)
				if !exists {
					return "expected an error"
				}

				if !errors.Is(gotErr, exp.(error)) {
					return "different error"
				}

				return ""
			},
		},
		{
			Name:    "addon-unavailable",
			ExpResp: orderbus.ErrAddonUnavailable,
			ExcFunc: func(ctx context.Context) any {
				addon, err := busDomain.Addon.Create(ctx, addonbus.NewAddon{
					CategoryID:   sd.Categories[0].ID,
					RestaurantID: sd.Restaurants[0].ID,
					Name:         name.MustParse("Unavailable Addon"),
					Description:  "temporarily unavailable",
					Price:        money.MustParse(1.50),
					MaxQuantity:  2,
				})
				if err != nil {
					return err
				}

				available := false
				if _, err := busDomain.Addon.Update(ctx, addon, addonbus.UpdateAddon{Available: &available}); err != nil {
					return err
				}

				no := orderbus.NewOrder{
					RestaurantID:  sd.Restaurants[0].ID.String(),
					CustomerName:  "Unavailable Addon Customer",
					CustomerEmail: "unavailable@example.com",
					CustomerPhone: "555-4444",
					OrderType:     orderbus.OrderTypePickup,
					PaymentMethod: orderbus.PaymentMethodCreditCard,
					Items: []orderbus.NewOrderItem{
						{
							MenuItemID: sd.MenuItems[0].ID.String(),
							Quantity:   1,
							Addons: []orderbus.NewOrderItemAddon{
								{AddonID: addon.ID.String(), Quantity: 1},
							},
						},
					},
				}

				_, err = busDomain.Order.Create(ctx, no)
				return err
			},
			CmpFunc: func(got any, exp any) string {
				gotErr, exists := got.(error)
				if !exists {
					return "expected an error"
				}

				if !errors.Is(gotErr, exp.(error)) {
					return "different error"
				}

				return ""
			},
		},
		{
			Name:    "addon-wrong-category",
			ExpResp: orderbus.ErrAddonCategoryMismatch,
			ExcFunc: func(ctx context.Context) any {
				no := orderbus.NewOrder{
					RestaurantID:  sd.Restaurants[0].ID.String(),
					CustomerName:  "Wrong Category",
					CustomerEmail: "wrongcat@example.com",
					CustomerPhone: "555-5555",
					OrderType:     orderbus.OrderTypePickup,
					PaymentMethod: orderbus.PaymentMethodCreditCard,
					Items: []orderbus.NewOrderItem{
						{
							MenuItemID: sd.MenuItems[0].ID.String(),
							Quantity:   1,
							Addons: []orderbus.NewOrderItemAddon{
								{AddonID: sd.Addons[2].ID.String(), Quantity: 1},
							},
						},
					},
				}

				_, err := busDomain.Order.Create(ctx, no)
				return err
			},
			CmpFunc: func(got any, exp any) string {
				gotErr, exists := got.(error)
				if !exists {
					return "expected an error"
				}

				if !errors.Is(gotErr, exp.(error)) {
					return "different error"
				}

				return ""
			},
		},
		{
			Name:    "addon-wrong-restaurant",
			ExpResp: nil,
			ExcFunc: func(ctx context.Context) any {
				no := orderbus.NewOrder{
					RestaurantID:  sd.Restaurants[0].ID.String(),
					CustomerName:  "Wrong Restaurant",
					CustomerEmail: "wrongrest@example.com",
					CustomerPhone: "555-6666",
					OrderType:     orderbus.OrderTypePickup,
					PaymentMethod: orderbus.PaymentMethodCreditCard,
					Items: []orderbus.NewOrderItem{
						{
							MenuItemID: sd.MenuItems[0].ID.String(),
							Quantity:   1,
							Addons: []orderbus.NewOrderItemAddon{
								{AddonID: sd.Addons[3].ID.String(), Quantity: 1},
							},
						},
					},
				}

				_, err := busDomain.Order.Create(ctx, no)
				return err
			},
			CmpFunc: func(got any, exp any) string {
				gotErr, exists := got.(error)
				if !exists {
					return "expected an error"
				}

				if gotErr == nil {
					return "expected an error"
				}

				return ""
			},
		},
	}

	return table
}

func updateStatus(busDomain dbtest.BusDomain, sd unittest.SeedData) []unittest.Table {
	table := []unittest.Table{
		{
			Name: "update-order-status",
			ExpResp: orderbus.Order{
				ID:                    sd.Orders[0].ID,
				RestaurantID:          sd.Orders[0].RestaurantID,
				CustomerName:          sd.Orders[0].CustomerName,
				CustomerEmail:         sd.Orders[0].CustomerEmail,
				CustomerPhone:         sd.Orders[0].CustomerPhone,
				OrderType:             sd.Orders[0].OrderType,
				OrderStatus:           orderbus.OrderStatusConfirmed,
				PaymentStatus:         orderbus.PaymentStatusPaid,
				PaymentMethod:         sd.Orders[0].PaymentMethod,
				Subtotal:              sd.Orders[0].Subtotal,
				DeliveryFee:           sd.Orders[0].DeliveryFee,
				Tax:                   sd.Orders[0].Tax,
				Total:                 sd.Orders[0].Total,
				SpecialInstructions:   sd.Orders[0].SpecialInstructions,
				StripePaymentIntentID: sd.Orders[0].StripePaymentIntentID,
				Items:                 sd.Orders[0].Items,
				DeliveryAddress:       sd.Orders[0].DeliveryAddress,
				DateCreated:           sd.Orders[0].DateCreated,
			},
			ExcFunc: func(ctx context.Context) any {
				us := orderbus.UpdateOrderStatus{
					OrderStatus:   orderbus.OrderStatusConfirmed,
					PaymentStatus: orderbus.PaymentStatusPaid,
				}

				if err := busDomain.Order.UpdateStatus(ctx, sd.Orders[0].ID, us); err != nil {
					return err
				}

				resp, err := busDomain.Order.QueryByID(ctx, sd.Orders[0].ID)
				if err != nil {
					return err
				}

				return resp
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.(orderbus.Order)
				if !exists {
					return "error occurred"
				}

				expResp := exp.(orderbus.Order)

				// Normalize DateUpdated since it changes
				expResp.DateUpdated = gotResp.DateUpdated

				// Normalize monetary values from database
				expResp.Subtotal = gotResp.Subtotal
				expResp.Tax = gotResp.Tax
				expResp.Total = gotResp.Total

				// Normalize item OrderIDs from database
				for i := range gotResp.Items {
					if i < len(expResp.Items) {
						expResp.Items[i].OrderID = gotResp.Items[i].OrderID
					}
				}

				// Normalize address OrderID if present
				if gotResp.DeliveryAddress != nil && expResp.DeliveryAddress != nil {
					expResp.DeliveryAddress.OrderID = gotResp.DeliveryAddress.OrderID
				}

				return cmp.Diff(gotResp, expResp, equateTime)
			},
		},
		{
			Name: "update-status-out-for-delivery",
			ExpResp: orderbus.Order{
				ID:                    sd.Orders[2].ID,
				RestaurantID:          sd.Orders[2].RestaurantID,
				CustomerName:          sd.Orders[2].CustomerName,
				CustomerEmail:         sd.Orders[2].CustomerEmail,
				CustomerPhone:         sd.Orders[2].CustomerPhone,
				OrderType:             sd.Orders[2].OrderType,
				OrderStatus:           orderbus.OrderStatusOutForDelivery,
				PaymentStatus:         sd.Orders[2].PaymentStatus,
				PaymentMethod:         sd.Orders[2].PaymentMethod,
				Subtotal:              sd.Orders[2].Subtotal,
				DeliveryFee:           sd.Orders[2].DeliveryFee,
				Tax:                   sd.Orders[2].Tax,
				Total:                 sd.Orders[2].Total,
				SpecialInstructions:   sd.Orders[2].SpecialInstructions,
				StripePaymentIntentID: sd.Orders[2].StripePaymentIntentID,
				Items:                 sd.Orders[2].Items,
				DeliveryAddress:       sd.Orders[2].DeliveryAddress,
				DateCreated:           sd.Orders[2].DateCreated,
			},
			ExcFunc: func(ctx context.Context) any {
				us := orderbus.UpdateOrderStatus{
					OrderStatus: orderbus.OrderStatusOutForDelivery,
				}

				if err := busDomain.Order.UpdateStatus(ctx, sd.Orders[2].ID, us); err != nil {
					return err
				}

				resp, err := busDomain.Order.QueryByID(ctx, sd.Orders[2].ID)
				if err != nil {
					return err
				}

				return resp
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.(orderbus.Order)
				if !exists {
					return "error occurred"
				}

				expResp := exp.(orderbus.Order)

				// Normalize DateUpdated since it changes
				expResp.DateUpdated = gotResp.DateUpdated

				// Normalize monetary values from database
				expResp.Subtotal = gotResp.Subtotal
				expResp.Tax = gotResp.Tax
				expResp.Total = gotResp.Total

				// Normalize item OrderIDs from database
				for i := range gotResp.Items {
					if i < len(expResp.Items) {
						expResp.Items[i].OrderID = gotResp.Items[i].OrderID
					}
				}

				// Normalize address OrderID if present
				if gotResp.DeliveryAddress != nil && expResp.DeliveryAddress != nil {
					expResp.DeliveryAddress.OrderID = gotResp.DeliveryAddress.OrderID
				}

				return cmp.Diff(gotResp, expResp, equateTime)
			},
		},
		{
			Name:    "update-status-out-for-delivery-pickup",
			ExpResp: orderbus.ErrOutForDeliveryRequiresDelivery,
			ExcFunc: func(ctx context.Context) any {
				us := orderbus.UpdateOrderStatus{
					OrderStatus: orderbus.OrderStatusOutForDelivery,
				}

				return busDomain.Order.UpdateStatus(ctx, sd.Orders[3].ID, us)
			},
			CmpFunc: func(got any, exp any) string {
				gotErr, exists := got.(error)
				if !exists {
					return "expected an error"
				}

				if !errors.Is(gotErr, exp.(error)) {
					return "different error"
				}

				return ""
			},
		},
	}

	return table
}

func cancel(busDomain dbtest.BusDomain, sd unittest.SeedData) []unittest.Table {
	table := []unittest.Table{
		{
			Name: "cancel-order",
			ExpResp: orderbus.Order{
				ID:                    sd.Orders[1].ID,
				RestaurantID:          sd.Orders[1].RestaurantID,
				CustomerName:          sd.Orders[1].CustomerName,
				CustomerEmail:         sd.Orders[1].CustomerEmail,
				CustomerPhone:         sd.Orders[1].CustomerPhone,
				OrderType:             sd.Orders[1].OrderType,
				OrderStatus:           orderbus.OrderStatusCancelled,
				PaymentStatus:         sd.Orders[1].PaymentStatus,
				PaymentMethod:         sd.Orders[1].PaymentMethod,
				Subtotal:              sd.Orders[1].Subtotal,
				DeliveryFee:           sd.Orders[1].DeliveryFee,
				Tax:                   sd.Orders[1].Tax,
				Total:                 sd.Orders[1].Total,
				SpecialInstructions:   sd.Orders[1].SpecialInstructions,
				StripePaymentIntentID: sd.Orders[1].StripePaymentIntentID,
				Items:                 sd.Orders[1].Items,
				DeliveryAddress:       sd.Orders[1].DeliveryAddress,
				DateCreated:           sd.Orders[1].DateCreated,
			},
			ExcFunc: func(ctx context.Context) any {
				if err := busDomain.Order.Cancel(ctx, sd.Orders[1].ID); err != nil {
					return err
				}

				resp, err := busDomain.Order.QueryByID(ctx, sd.Orders[1].ID)
				if err != nil {
					return err
				}

				return resp
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.(orderbus.Order)
				if !exists {
					return "error occurred"
				}

				expResp := exp.(orderbus.Order)

				// Normalize DateUpdated since it changes
				expResp.DateUpdated = gotResp.DateUpdated

				// Normalize monetary values from database
				expResp.Subtotal = gotResp.Subtotal
				expResp.Tax = gotResp.Tax
				expResp.Total = gotResp.Total

				// Normalize item OrderIDs from database
				for i := range gotResp.Items {
					if i < len(expResp.Items) {
						expResp.Items[i].OrderID = gotResp.Items[i].OrderID
					}
				}

				// Normalize address OrderID if present
				if gotResp.DeliveryAddress != nil && expResp.DeliveryAddress != nil {
					expResp.DeliveryAddress.OrderID = gotResp.DeliveryAddress.OrderID
				}

				return cmp.Diff(gotResp, expResp, equateTime)
			},
		},
	}

	return table
}
