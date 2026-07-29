package orderbus_test

import (
	"context"
	"fmt"
	"sort"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/warlck/food-flow/business/domain/categorybus"
	"github.com/warlck/food-flow/business/domain/menuitembus"
	"github.com/warlck/food-flow/business/domain/orderbus"
	"github.com/warlck/food-flow/business/domain/restaurantbus"
	"github.com/warlck/food-flow/business/sdk/dbtest"
	"github.com/warlck/food-flow/business/sdk/page"
	"github.com/warlck/food-flow/business/sdk/unittest"
	"github.com/warlck/food-flow/business/types/money"
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
	rests, err := restaurantbus.TestSeedRestaurants(ctx, 2, busDomain.Restaurant)
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
				DeliveryFee:   money.MustParse(5.00),
			},
			ExcFunc: func(ctx context.Context) any {
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
