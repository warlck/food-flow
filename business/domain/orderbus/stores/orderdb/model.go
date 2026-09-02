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
	PromoCode             *string   `db:"promo_code"`
	Subtotal              float64   `db:"subtotal"`
	Discount              float64   `db:"discount"`
	DeliveryFee           float64   `db:"delivery_fee"`
	Tax                   float64   `db:"tax"`
	Total                 float64   `db:"total"`
	SpecialInstructions   *string   `db:"special_instructions"`
	StripePaymentIntentID *string   `db:"stripe_payment_intent_id"`
	DateCreated           time.Time `db:"date_created"`
	DateUpdated           time.Time `db:"date_updated"`
}

// dbOrderItem represents the database row for an order item.
type dbOrderItem struct {
	ID                  uuid.UUID `db:"order_item_id"`
	OrderID             uuid.UUID `db:"order_id"`
	CategoryID          uuid.UUID `db:"category_id"`
	CategoryName        string    `db:"category_name"`
	MenuItemID          uuid.UUID `db:"menu_item_id"`
	MenuItemName        string    `db:"menu_item_name"`
	MenuItemPrice       float64   `db:"menu_item_price"`
	Quantity            int       `db:"quantity"`
	SpecialInstructions *string   `db:"special_instructions"`
	DateCreated         time.Time `db:"date_created"`
}

// dbOrderItemModifier represents the database row for an order item modifier.
type dbOrderItemModifier struct {
	ID                 uuid.UUID `db:"order_item_modifier_id"`
	OrderItemID        uuid.UUID `db:"order_item_id"`
	ModifierGroupID    uuid.UUID `db:"modifier_group_id"`
	ModifierGroupName  string    `db:"modifier_group_name"`
	ModifierOptionID   uuid.UUID `db:"modifier_option_id"`
	ModifierOptionName string    `db:"modifier_option_name"`
	PriceDelta         float64   `db:"price_delta"`
	DateCreated        time.Time `db:"date_created"`
}

// dbOrderItemAddon represents the database row for an order item addon.
type dbOrderItemAddon struct {
	ID          uuid.UUID `db:"order_item_addon_id"`
	OrderItemID uuid.UUID `db:"order_item_id"`
	AddonID     uuid.UUID `db:"addon_id"`
	AddonName   string    `db:"addon_name"`
	AddonPrice  float64   `db:"addon_price"`
	Quantity    int       `db:"quantity"`
	DateCreated time.Time `db:"date_created"`
}

// dbDeliveryAddress represents the database row for a delivery address.
type dbDeliveryAddress struct {
	ID                   uuid.UUID `db:"address_id"`
	OrderID              uuid.UUID `db:"order_id"`
	Street               string    `db:"street"`
	City                 string    `db:"city"`
	State                string    `db:"state"`
	PostalCode           string    `db:"postal_code"`
	DeliveryInstructions *string   `db:"delivery_instructions"`
	Latitude             *float64  `db:"latitude"`
	Longitude            *float64  `db:"longitude"`
	DateCreated          time.Time `db:"date_created"`
}

// =============================================================================

func toDBOrder(bus orderbus.Order) dbOrder {
	var promoCode *string
	if bus.PromoCode != "" {
		promoCode = &bus.PromoCode
	}

	var specialInstructions *string
	if bus.SpecialInstructions != "" {
		specialInstructions = &bus.SpecialInstructions
	}

	var stripePaymentIntentID *string
	if bus.StripePaymentIntentID != "" {
		stripePaymentIntentID = &bus.StripePaymentIntentID
	}

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
		PromoCode:             promoCode,
		Subtotal:              bus.Subtotal.Value(),
		Discount:              bus.Discount.Value(),
		DeliveryFee:           bus.DeliveryFee.Value(),
		Tax:                   bus.Tax.Value(),
		Total:                 bus.Total.Value(),
		SpecialInstructions:   specialInstructions,
		StripePaymentIntentID: stripePaymentIntentID,
		DateCreated:           bus.DateCreated.UTC(),
		DateUpdated:           bus.DateUpdated.UTC(),
	}
}

func toBusOrder(dbo dbOrder, items []dbOrderItem, modifiers []dbOrderItemModifier, addons []dbOrderItemAddon, addr *dbDeliveryAddress) (orderbus.Order, error) {
	subtotal, err := money.Parse(dbo.Subtotal)
	if err != nil {
		return orderbus.Order{}, fmt.Errorf("parse subtotal: %w", err)
	}

	discount, err := money.Parse(dbo.Discount)
	if err != nil {
		return orderbus.Order{}, fmt.Errorf("parse discount: %w", err)
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

	var promoCodeStr string
	if dbo.PromoCode != nil {
		promoCodeStr = *dbo.PromoCode
	}

	var specialInstructionsStr string
	if dbo.SpecialInstructions != nil {
		specialInstructionsStr = *dbo.SpecialInstructions
	}

	var stripePaymentIntentIDStr string
	if dbo.StripePaymentIntentID != nil {
		stripePaymentIntentIDStr = *dbo.StripePaymentIntentID
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
		PromoCode:             promoCodeStr,
		Subtotal:              subtotal,
		Discount:              discount,
		DeliveryFee:           deliveryFee,
		Tax:                   tax,
		Total:                 total,
		SpecialInstructions:   specialInstructionsStr,
		StripePaymentIntentID: stripePaymentIntentIDStr,
		DateCreated:           dbo.DateCreated.In(time.Local),
		DateUpdated:           dbo.DateUpdated.In(time.Local),
		Items:                 make([]orderbus.OrderItem, len(items)),
	}

	modsByItem := make(map[uuid.UUID][]orderbus.OrderItemModifier)
	for _, mod := range modifiers {
		modsByItem[mod.OrderItemID] = append(modsByItem[mod.OrderItemID], toBusOrderItemModifier(mod))
	}

	addonsByItem := make(map[uuid.UUID][]orderbus.OrderItemAddon)
	for _, addon := range addons {
		addonsByItem[addon.OrderItemID] = append(addonsByItem[addon.OrderItemID], toBusOrderItemAddon(addon))
	}

	for i, item := range items {
		busItem := toBusOrderItem(item)
		busItem.Modifiers = modsByItem[item.ID]
		busItem.Addons = addonsByItem[item.ID]
		order.Items[i] = busItem
	}

	if addr != nil {
		da := toBusDeliveryAddress(*addr, dbo.ID)
		order.DeliveryAddress = &da
	}

	return order, nil
}

func toDBOrderItem(bus orderbus.OrderItem, orderID uuid.UUID) dbOrderItem {
	var specialInstructions *string
	if bus.SpecialInstructions != "" {
		specialInstructions = &bus.SpecialInstructions
	}

	return dbOrderItem{
		ID:                  bus.ID,
		OrderID:             orderID,
		CategoryID:          bus.CategoryID,
		CategoryName:        bus.CategoryName,
		MenuItemID:          bus.MenuItemID,
		MenuItemName:        bus.MenuItemName,
		MenuItemPrice:       bus.MenuItemPrice.Value(),
		Quantity:            bus.Quantity,
		SpecialInstructions: specialInstructions,
		DateCreated:         bus.DateCreated.UTC(),
	}
}

func toBusOrderItem(dbo dbOrderItem) orderbus.OrderItem {
	price, err := money.Parse(dbo.MenuItemPrice)
	if err != nil {
		price = money.Money{}
	}

	var specialInstructionsStr string
	if dbo.SpecialInstructions != nil {
		specialInstructionsStr = *dbo.SpecialInstructions
	}

	return orderbus.OrderItem{
		ID:                  dbo.ID,
		OrderID:             dbo.OrderID,
		CategoryID:          dbo.CategoryID,
		CategoryName:        dbo.CategoryName,
		MenuItemID:          dbo.MenuItemID,
		MenuItemName:        dbo.MenuItemName,
		MenuItemPrice:       price,
		Quantity:            dbo.Quantity,
		SpecialInstructions: specialInstructionsStr,
		DateCreated:         dbo.DateCreated.In(time.Local),
	}
}

func toDBOrderItemModifier(bus orderbus.OrderItemModifier, orderItemID uuid.UUID) dbOrderItemModifier {
	return dbOrderItemModifier{
		ID:                 bus.ID,
		OrderItemID:        orderItemID,
		ModifierGroupID:    bus.ModifierGroupID,
		ModifierGroupName:  bus.ModifierGroupName,
		ModifierOptionID:   bus.ModifierOptionID,
		ModifierOptionName: bus.ModifierOptionName,
		PriceDelta:         bus.PriceDelta.Value(),
		DateCreated:        bus.DateCreated.UTC(),
	}
}

func toBusOrderItemModifier(dbo dbOrderItemModifier) orderbus.OrderItemModifier {
	priceDelta, err := money.Parse(dbo.PriceDelta)
	if err != nil {
		priceDelta = money.Money{}
	}

	return orderbus.OrderItemModifier{
		ID:                 dbo.ID,
		OrderItemID:        dbo.OrderItemID,
		ModifierGroupID:    dbo.ModifierGroupID,
		ModifierGroupName:  dbo.ModifierGroupName,
		ModifierOptionID:   dbo.ModifierOptionID,
		ModifierOptionName: dbo.ModifierOptionName,
		PriceDelta:         priceDelta,
		DateCreated:        dbo.DateCreated.In(time.Local),
	}
}

func toDBOrderItemAddon(bus orderbus.OrderItemAddon, orderItemID uuid.UUID) dbOrderItemAddon {
	return dbOrderItemAddon{
		ID:          bus.ID,
		OrderItemID: orderItemID,
		AddonID:     bus.AddonID,
		AddonName:   bus.AddonName,
		AddonPrice:  bus.AddonPrice.Value(),
		Quantity:    bus.Quantity,
		DateCreated: bus.DateCreated.UTC(),
	}
}

func toBusOrderItemAddon(dbo dbOrderItemAddon) orderbus.OrderItemAddon {
	price, err := money.Parse(dbo.AddonPrice)
	if err != nil {
		price = money.Money{}
	}

	return orderbus.OrderItemAddon{
		ID:          dbo.ID,
		OrderItemID: dbo.OrderItemID,
		AddonID:     dbo.AddonID,
		AddonName:   dbo.AddonName,
		AddonPrice:  price,
		Quantity:    dbo.Quantity,
		DateCreated: dbo.DateCreated.In(time.Local),
	}
}

func toDBDeliveryAddress(bus orderbus.DeliveryAddress, orderID uuid.UUID) dbDeliveryAddress {
	var deliveryInstructions *string
	if bus.DeliveryInstructions != "" {
		deliveryInstructions = &bus.DeliveryInstructions
	}

	return dbDeliveryAddress{
		ID:                   bus.ID,
		OrderID:              orderID,
		Street:               bus.Street,
		City:                 bus.City,
		State:                bus.State,
		PostalCode:           bus.PostalCode,
		DeliveryInstructions: deliveryInstructions,
		Latitude:             bus.Latitude,
		Longitude:            bus.Longitude,
		DateCreated:          bus.DateCreated.UTC(),
	}
}

func toBusDeliveryAddress(dbo dbDeliveryAddress, orderID uuid.UUID) orderbus.DeliveryAddress {
	var deliveryInstructionsStr string
	if dbo.DeliveryInstructions != nil {
		deliveryInstructionsStr = *dbo.DeliveryInstructions
	}

	return orderbus.DeliveryAddress{
		ID:                   dbo.ID,
		OrderID:              orderID,
		Street:               dbo.Street,
		City:                 dbo.City,
		State:                dbo.State,
		PostalCode:           dbo.PostalCode,
		DeliveryInstructions: deliveryInstructionsStr,
		Latitude:             dbo.Latitude,
		Longitude:            dbo.Longitude,
		DateCreated:          dbo.DateCreated.In(time.Local),
	}
}
