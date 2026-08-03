package orderdb

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/warlck/food-flow/business/domain/orderbus"
	"github.com/warlck/food-flow/business/types/money"
)

// dbOrder represents the database row for an order.
type dbOrder struct {
	ID                    uuid.UUID `db:"order_id"`
	RestaurantID          uuid.UUID `db:"restaurant_id"`
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
	ID                  uuid.UUID `db:"order_item_id"`
	OrderID             uuid.UUID `db:"order_id"`
	MenuItemID          uuid.UUID `db:"menu_item_id"`
	MenuItemName        string    `db:"menu_item_name"`
	MenuItemPrice       float64   `db:"menu_item_price"`
	Quantity            int       `db:"quantity"`
	SpecialInstructions string    `db:"special_instructions"`
	DateCreated         time.Time `db:"date_created"`
}

// dbDeliveryAddress represents the database row for a delivery address.
type dbDeliveryAddress struct {
	ID                   uuid.UUID `db:"address_id"`
	OrderID              uuid.UUID `db:"order_id"`
	Street               string    `db:"street"`
	City                 string    `db:"city"`
	State                string    `db:"state"`
	PostalCode           string    `db:"postal_code"`
	DeliveryInstructions string    `db:"delivery_instructions"`
	Latitude             *float64  `db:"latitude"`
	Longitude            *float64  `db:"longitude"`
	DateCreated          time.Time `db:"date_created"`
}

// =============================================================================

func toDBOrder(bus orderbus.Order) dbOrder {
	return dbOrder{
		ID:                    bus.ID,
		RestaurantID:          bus.RestaurantID,
		CustomerName:          bus.CustomerName,
		CustomerEmail:         bus.CustomerEmail,
		CustomerPhone:         bus.CustomerPhone,
		OrderType:             bus.OrderType,
		OrderStatus:           bus.OrderStatus,
		PaymentStatus:         bus.PaymentStatus,
		PaymentMethod:         bus.PaymentMethod,
		Subtotal:              bus.Subtotal.Value(),
		DeliveryFee:           bus.DeliveryFee.Value(),
		Tax:                   bus.Tax.Value(),
		Total:                 bus.Total.Value(),
		SpecialInstructions:   bus.SpecialInstructions,
		StripePaymentIntentID: bus.StripePaymentIntentID,
		DateCreated:           bus.DateCreated.UTC(),
		DateUpdated:           bus.DateUpdated.UTC(),
	}
}

func toBusOrder(dbo dbOrder, items []dbOrderItem, addr *dbDeliveryAddress) (orderbus.Order, error) {
	subtotal, err := money.Parse(dbo.Subtotal)
	if err != nil {
		return orderbus.Order{}, fmt.Errorf("parse subtotal: %w", err)
	}

	deliveryFee, err := money.Parse(dbo.DeliveryFee)
	if err != nil {
		return orderbus.Order{}, fmt.Errorf("parse delivery fee: %w", err)
	}

	tax, err := money.Parse(dbo.Tax)
	if err != nil {
		return orderbus.Order{}, fmt.Errorf("parse tax: %w", err)
	}

	total, err := money.Parse(dbo.Total)
	if err != nil {
		return orderbus.Order{}, fmt.Errorf("parse total: %w", err)
	}

	order := orderbus.Order{
		ID:                    dbo.ID,
		RestaurantID:          dbo.RestaurantID,
		CustomerName:          dbo.CustomerName,
		CustomerEmail:         dbo.CustomerEmail,
		CustomerPhone:         dbo.CustomerPhone,
		OrderType:             dbo.OrderType,
		OrderStatus:           dbo.OrderStatus,
		PaymentStatus:         dbo.PaymentStatus,
		PaymentMethod:         dbo.PaymentMethod,
		Subtotal:              subtotal,
		DeliveryFee:           deliveryFee,
		Tax:                   tax,
		Total:                 total,
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
		da := toBusDeliveryAddress(*addr, dbo.ID)
		order.DeliveryAddress = &da
	}

	return order, nil
}

func toDBOrderItem(bus orderbus.OrderItem, orderID uuid.UUID) dbOrderItem {
	return dbOrderItem{
		ID:                  bus.ID,
		OrderID:             orderID,
		MenuItemID:          bus.MenuItemID,
		MenuItemName:        bus.MenuItemName,
		MenuItemPrice:       bus.MenuItemPrice.Value(),
		Quantity:            bus.Quantity,
		SpecialInstructions: bus.SpecialInstructions,
		DateCreated:         bus.DateCreated.UTC(),
	}
}

func toBusOrderItem(dbo dbOrderItem) orderbus.OrderItem {
	price, err := money.Parse(dbo.MenuItemPrice)
	if err != nil {
		// This should never happen with valid database data
		// but we need to handle it gracefully
		price = money.Money{}
	}

	return orderbus.OrderItem{
		ID:                  dbo.ID,
		OrderID:             dbo.OrderID,
		MenuItemID:          dbo.MenuItemID,
		MenuItemName:        dbo.MenuItemName,
		MenuItemPrice:       price,
		Quantity:            dbo.Quantity,
		SpecialInstructions: dbo.SpecialInstructions,
		DateCreated:         dbo.DateCreated.In(time.Local),
	}
}

func toDBDeliveryAddress(bus orderbus.DeliveryAddress, orderID uuid.UUID) dbDeliveryAddress {
	return dbDeliveryAddress{
		ID:                   bus.ID,
		OrderID:              orderID,
		Street:               bus.Street,
		City:                 bus.City,
		State:                bus.State,
		PostalCode:           bus.PostalCode,
		DeliveryInstructions: bus.DeliveryInstructions,
		Latitude:             bus.Latitude,
		Longitude:            bus.Longitude,
		DateCreated:          bus.DateCreated.UTC(),
	}
}

func toBusDeliveryAddress(dbo dbDeliveryAddress, orderID uuid.UUID) orderbus.DeliveryAddress {
	return orderbus.DeliveryAddress{
		ID:                   dbo.ID,
		OrderID:              orderID,
		Street:               dbo.Street,
		City:                 dbo.City,
		State:                dbo.State,
		PostalCode:           dbo.PostalCode,
		DeliveryInstructions: dbo.DeliveryInstructions,
		Latitude:             dbo.Latitude,
		Longitude:            dbo.Longitude,
		DateCreated:          dbo.DateCreated.In(time.Local),
	}
}
