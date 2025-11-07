package orderbus

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/warlck/food-flow/business/domain/menuitembus"
)

// TestSeedOrders is a helper method for testing.
func TestSeedOrders(ctx context.Context, n int, restaurantID uuid.UUID, menuItems []menuitembus.MenuItem, bus *Business) ([]Order, error) {
	orders := make([]Order, n)

	for i := 0; i < n; i++ {
		orderType := OrderTypePickup
		var deliveryAddr *NewDeliveryAddress

		if i%2 == 0 {
			orderType = OrderTypeDelivery
			deliveryAddr = &NewDeliveryAddress{
				Street:               fmt.Sprintf("%d Main St", i+100),
				City:                 "Test City",
				State:                "CA",
				PostalCode:           "12345",
				DeliveryInstructions: "Leave at door",
			}
		}

		items := make([]NewOrderItem, len(menuItems))
		for j, mi := range menuItems {
			items[j] = NewOrderItem{
				MenuItemID:          mi.ID.String(),
				Quantity:            j + 1,
				SpecialInstructions: fmt.Sprintf("Special request %d", j),
			}
		}

		no := NewOrder{
			RestaurantID:        restaurantID.String(),
			CustomerName:        fmt.Sprintf("Customer %d", i),
			CustomerEmail:       fmt.Sprintf("customer%d@test.com", i),
			CustomerPhone:       fmt.Sprintf("555-000%d", i),
			OrderType:           orderType,
			PaymentMethod:       PaymentMethodCreditCard,
			Items:               items,
			DeliveryAddress:     deliveryAddr,
			SpecialInstructions: fmt.Sprintf("Order note %d", i),
		}

		order, err := bus.Create(ctx, no)
		if err != nil {
			return nil, fmt.Errorf("seeding order: %w", err)
		}

		orders[i] = order
	}

	return orders, nil
}
