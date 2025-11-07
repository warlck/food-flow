package orderdb

import (
	"time"

	"github.com/google/uuid"
	"github.com/warlck/food-flow/business/domain/orderbus"
)

// dbOrder represents the database row for an order.
type dbOrder struct {
	ID                    string    `db:"order_id"`
	RestaurantID          string    `db:"restaurant_id"`
	CustomerName          string    `db:"customer_name"`
	CustomerEmail         string    `db:"customer_email"`
	CustomerPhone         string    `db:"customer_phone"`
	OrderType             string    `db:"order_type"`
	OrderStatus           string    `db:"order_status"`
	PaymentStatus         string    `db:"payment_status"`
	PaymentMethod         string    `db:"payment_method"`
	Subtotal              float64   `db:"subtotal"`
	DeliveryFee           float64   `db:"delivery_fee"`
	Tax                   float64   `db:"tax"`
	Total                 float64   `db:"total"`
	SpecialInstructions   string    `db:"special_instructions"`
	StripePaymentIntentID string    `db:"stripe_payment_intent_id"`
	DateCreated           time.Time `db:"date_created"`
	DateUpdated           time.Time `db:"date_updated"`
}

// dbOrderItem represents the database row for an order item.
type dbOrderItem struct {
	ID                  string    `db:"order_item_id"`
	OrderID             string    `db:"order_id"`
	MenuItemID          string    `db:"menu_item_id"`
	MenuItemName        string    `db:"menu_item_name"`
	MenuItemPrice       float64   `db:"menu_item_price"`
	Quantity            int       `db:"quantity"`
	SpecialInstructions string    `db:"special_instructions"`
	DateCreated         time.Time `db:"date_created"`
}

// dbDeliveryAddress represents the database row for a delivery address.
type dbDeliveryAddress struct {
	ID                   string    `db:"address_id"`
	OrderID              string    `db:"order_id"`
	Street               string    `db:"street"`
	City                 string    `db:"city"`
	State                string    `db:"state"`
	PostalCode           string    `db:"postal_code"`
	DeliveryInstructions string    `db:"delivery_instructions"`
	DateCreated          time.Time `db:"date_created"`
}

// =============================================================================

func toDBOrder(bus orderbus.Order) dbOrder {
	return dbOrder{
		ID:                    bus.ID.String(),
		RestaurantID:          bus.RestaurantID.String(),
		CustomerName:          bus.CustomerName,
		CustomerEmail:         bus.CustomerEmail,
		CustomerPhone:         bus.CustomerPhone,
		OrderType:             bus.OrderType,
		OrderStatus:           bus.OrderStatus,
		PaymentStatus:         bus.PaymentStatus,
		PaymentMethod:         bus.PaymentMethod,
		Subtotal:              bus.Subtotal,
		DeliveryFee:           bus.DeliveryFee,
		Tax:                   bus.Tax,
		Total:                 bus.Total,
		SpecialInstructions:   bus.SpecialInstructions,
		StripePaymentIntentID: bus.StripePaymentIntentID,
		DateCreated:           bus.DateCreated.UTC(),
		DateUpdated:           bus.DateUpdated.UTC(),
	}
}

func toBusOrder(dbo dbOrder, items []dbOrderItem, addr *dbDeliveryAddress) orderbus.Order {
	orderID := uuid.MustParse(dbo.ID)
	restaurantID := uuid.MustParse(dbo.RestaurantID)

	order := orderbus.Order{
		ID:                    orderID,
		RestaurantID:          restaurantID,
		CustomerName:          dbo.CustomerName,
		CustomerEmail:         dbo.CustomerEmail,
		CustomerPhone:         dbo.CustomerPhone,
		OrderType:             dbo.OrderType,
		OrderStatus:           dbo.OrderStatus,
		PaymentStatus:         dbo.PaymentStatus,
		PaymentMethod:         dbo.PaymentMethod,
		Subtotal:              dbo.Subtotal,
		DeliveryFee:           dbo.DeliveryFee,
		Tax:                   dbo.Tax,
		Total:                 dbo.Total,
		SpecialInstructions:   dbo.SpecialInstructions,
		StripePaymentIntentID: dbo.StripePaymentIntentID,
		DateCreated:           dbo.DateCreated.In(time.Local),
		DateUpdated:           dbo.DateUpdated.In(time.Local),
		Items:                 make([]orderbus.OrderItem, len(items)),
	}

	for i, item := range items {
		order.Items[i] = toBusOrderItem(item)
	}

	if addr != nil {
		da := toBusDeliveryAddress(*addr, orderID)
		order.DeliveryAddress = &da
	}

	return order
}

func toDBOrderItem(bus orderbus.OrderItem, orderID uuid.UUID) dbOrderItem {
	return dbOrderItem{
		ID:                  bus.ID.String(),
		OrderID:             orderID.String(),
		MenuItemID:          bus.MenuItemID.String(),
		MenuItemName:        bus.MenuItemName,
		MenuItemPrice:       bus.MenuItemPrice,
		Quantity:            bus.Quantity,
		SpecialInstructions: bus.SpecialInstructions,
		DateCreated:         bus.DateCreated.UTC(),
	}
}

func toBusOrderItem(dbo dbOrderItem) orderbus.OrderItem {
	return orderbus.OrderItem{
		ID:                  uuid.MustParse(dbo.ID),
		OrderID:             uuid.MustParse(dbo.OrderID),
		MenuItemID:          uuid.MustParse(dbo.MenuItemID),
		MenuItemName:        dbo.MenuItemName,
		MenuItemPrice:       dbo.MenuItemPrice,
		Quantity:            dbo.Quantity,
		SpecialInstructions: dbo.SpecialInstructions,
		DateCreated:         dbo.DateCreated.In(time.Local),
	}
}

func toDBDeliveryAddress(bus orderbus.DeliveryAddress, orderID uuid.UUID) dbDeliveryAddress {
	return dbDeliveryAddress{
		ID:                   bus.ID.String(),
		OrderID:              orderID.String(),
		Street:               bus.Street,
		City:                 bus.City,
		State:                bus.State,
		PostalCode:           bus.PostalCode,
		DeliveryInstructions: bus.DeliveryInstructions,
		DateCreated:          bus.DateCreated.UTC(),
	}
}

func toBusDeliveryAddress(dbo dbDeliveryAddress, orderID uuid.UUID) orderbus.DeliveryAddress {
	return orderbus.DeliveryAddress{
		ID:                   uuid.MustParse(dbo.ID),
		OrderID:              orderID,
		Street:               dbo.Street,
		City:                 dbo.City,
		State:                dbo.State,
		PostalCode:           dbo.PostalCode,
		DeliveryInstructions: dbo.DeliveryInstructions,
		DateCreated:          dbo.DateCreated.In(time.Local),
	}
}
