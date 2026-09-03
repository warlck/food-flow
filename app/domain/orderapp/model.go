package orderapp

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/warlck/food-flow/app/sdk/errs"
	"github.com/warlck/food-flow/business/domain/orderbus"
)

// Order represents information about an order for API responses.
type Order struct {
	ID                    string           `json:"id"`
	RestaurantID          string           `json:"restaurantId"`
	CustomerName          string           `json:"customerName"`
	CustomerEmail         string           `json:"customerEmail"`
	CustomerPhone         string           `json:"customerPhone"`
	OrderType             string           `json:"orderType"`
	OrderStatus           string           `json:"orderStatus"`
	PaymentStatus         string           `json:"paymentStatus"`
	PaymentMethod         string           `json:"paymentMethod"`
	PromoCode             string           `json:"promoCode,omitempty"`
	Subtotal              float64          `json:"subtotal"`
	Discount              float64          `json:"discount"`
	DeliveryFee           float64          `json:"deliveryFee"`
	Tax                   float64          `json:"tax"`
	Total                 float64          `json:"total"`
	SpecialInstructions   string           `json:"specialInstructions,omitempty"`
	StripePaymentIntentID string           `json:"stripePaymentIntentId,omitempty"`
	Items                 []OrderItem      `json:"items"`
	DeliveryAddress       *DeliveryAddress `json:"deliveryAddress,omitempty"`
	DateCreated           string           `json:"dateCreated"`
	DateUpdated           string           `json:"dateUpdated"`
}

// OrderItem represents an item in an order.
type OrderItem struct {
	ID                  string              `json:"id"`
	CategoryID          string              `json:"categoryId"`
	CategoryName        string              `json:"categoryName"`
	MenuItemID          string              `json:"menuItemId"`
	MenuItemName        string              `json:"menuItemName"`
	MenuItemPrice       float64             `json:"menuItemPrice"`
	Quantity            int                 `json:"quantity"`
	SpecialInstructions string              `json:"specialInstructions,omitempty"`
	Modifiers           []OrderItemModifier `json:"modifiers"`
	Addons              []OrderItemAddon    `json:"addons"`
	DateCreated         string              `json:"dateCreated"`
}

// OrderItemModifier represents a modifier option applied to an order item.
type OrderItemModifier struct {
	ID                 string  `json:"id"`
	ModifierGroupID    string  `json:"modifierGroupId"`
	ModifierGroupName  string  `json:"modifierGroupName"`
	ModifierOptionID   string  `json:"modifierOptionId"`
	ModifierOptionName string  `json:"modifierOptionName"`
	PriceDelta         float64 `json:"priceDelta"`
}

// OrderItemAddon represents an addon applied to an order item.
type OrderItemAddon struct {
	ID         string  `json:"id"`
	AddonID    string  `json:"addonId"`
	AddonName  string  `json:"addonName"`
	AddonPrice float64 `json:"addonPrice"`
	Quantity   int     `json:"quantity"`
}

// DeliveryAddress represents a delivery address for an order.
type DeliveryAddress struct {
	ID                   string   `json:"id"`
	Street               string   `json:"street"`
	City                 string   `json:"city"`
	State                string   `json:"state"`
	PostalCode           string   `json:"postalCode"`
	DeliveryInstructions string   `json:"deliveryInstructions,omitempty"`
	Latitude             *float64 `json:"latitude,omitempty"`
	Longitude            *float64 `json:"longitude,omitempty"`
	DateCreated          string   `json:"dateCreated"`
}

// Encode implements the encoder interface.
func (app Order) Encode() ([]byte, string, error) {
	data, err := json.Marshal(app)
	return data, "application/json", err
}

// ToAppOrder converts a business layer order to an app layer order.
func ToAppOrder(bus orderbus.Order) Order {
	items := make([]OrderItem, len(bus.Items))
	for i, item := range bus.Items {
		modifiers := make([]OrderItemModifier, len(item.Modifiers))
		for j, mod := range item.Modifiers {
			modifiers[j] = OrderItemModifier{
				ID:                 mod.ID.String(),
				ModifierGroupID:    mod.ModifierGroupID.String(),
				ModifierGroupName:  mod.ModifierGroupName,
				ModifierOptionID:   mod.ModifierOptionID.String(),
				ModifierOptionName: mod.ModifierOptionName,
				PriceDelta:         mod.PriceDelta.Value(),
			}
		}

		addons := make([]OrderItemAddon, len(item.Addons))
		for j, addon := range item.Addons {
			addons[j] = OrderItemAddon{
				ID:         addon.ID.String(),
				AddonID:    addon.AddonID.String(),
				AddonName:  addon.AddonName,
				AddonPrice: addon.AddonPrice.Value(),
				Quantity:   addon.Quantity,
			}
		}

		items[i] = OrderItem{
			ID:                  item.ID.String(),
			CategoryID:          item.CategoryID.String(),
			CategoryName:        item.CategoryName,
			MenuItemID:          item.MenuItemID.String(),
			MenuItemName:        item.MenuItemName,
			MenuItemPrice:       item.MenuItemPrice.Value(),
			Quantity:            item.Quantity,
			SpecialInstructions: item.SpecialInstructions,
			Modifiers:           modifiers,
			Addons:              addons,
			DateCreated:         item.DateCreated.Format(time.RFC3339),
		}
	}

	var deliveryAddr *DeliveryAddress
	if bus.DeliveryAddress != nil {
		deliveryAddr = &DeliveryAddress{
			ID:                   bus.DeliveryAddress.ID.String(),
			Street:               bus.DeliveryAddress.Street,
			City:                 bus.DeliveryAddress.City,
			State:                bus.DeliveryAddress.State,
			PostalCode:           bus.DeliveryAddress.PostalCode,
			DeliveryInstructions: bus.DeliveryAddress.DeliveryInstructions,
			Latitude:             bus.DeliveryAddress.Latitude,
			Longitude:            bus.DeliveryAddress.Longitude,
			DateCreated:          bus.DeliveryAddress.DateCreated.Format(time.RFC3339),
		}
	}

	return Order{
		ID:                    bus.ID.String(),
		RestaurantID:          bus.RestaurantID.String(),
		CustomerName:          bus.CustomerName,
		CustomerEmail:         bus.CustomerEmail,
		CustomerPhone:         bus.CustomerPhone,
		OrderType:             bus.OrderType,
		OrderStatus:           bus.OrderStatus,
		PaymentStatus:         bus.PaymentStatus,
		PaymentMethod:         bus.PaymentMethod,
		PromoCode:             bus.PromoCode,
		Subtotal:              bus.Subtotal.Value(),
		Discount:              bus.Discount.Value(),
		DeliveryFee:           bus.DeliveryFee.Value(),
		Tax:                   bus.Tax.Value(),
		Total:                 bus.Total.Value(),
		SpecialInstructions:   bus.SpecialInstructions,
		StripePaymentIntentID: bus.StripePaymentIntentID,
		Items:                 items,
		DeliveryAddress:       deliveryAddr,
		DateCreated:           bus.DateCreated.Format(time.RFC3339),
		DateUpdated:           bus.DateUpdated.Format(time.RFC3339),
	}
}

// ToAppOrders converts a slice of business layer orders to app layer orders.
func ToAppOrders(orders []orderbus.Order) []Order {
	app := make([]Order, len(orders))
	for i, ord := range orders {
		app[i] = ToAppOrder(ord)
	}

	return app
}

// =============================================================================

// NewOrder defines the data needed to add a new order.
type NewOrder struct {
	RestaurantID        string              `json:"restaurantId" validate:"required"`
	CustomerName        string              `json:"customerName" validate:"required"`
	CustomerEmail       string              `json:"customerEmail" validate:"required,email"`
	CustomerPhone       string              `json:"customerPhone" validate:"required"`
	OrderType           string              `json:"orderType" validate:"required,oneof=pickup delivery"`
	PaymentMethod       string              `json:"paymentMethod" validate:"required,oneof=creditCard cash"`
	PromoCode           string              `json:"promoCode"`
	SpecialInstructions string              `json:"specialInstructions"`
	Items               []NewOrderItem      `json:"items" validate:"required,min=1,dive"`
	DeliveryAddress     *NewDeliveryAddress `json:"deliveryAddress"`
}

// NewOrderItem defines the data needed to add an item to an order.
type NewOrderItem struct {
	MenuItemID          string                 `json:"menuItemId" validate:"required"`
	Quantity            int                    `json:"quantity" validate:"required,min=1"`
	SpecialInstructions string                 `json:"specialInstructions"`
	Modifiers           []NewOrderItemModifier `json:"modifiers" validate:"omitempty,dive"`
	Addons              []NewOrderItemAddon    `json:"addons" validate:"omitempty,dive"`
}

// NewOrderItemModifier defines the data needed to add a modifier option to an order item.
type NewOrderItemModifier struct {
	ModifierGroupID  string `json:"modifierGroupId" validate:"required,uuid"`
	ModifierOptionID string `json:"modifierOptionId" validate:"required,uuid"`
}

// NewOrderItemAddon defines the data needed to add an addon to an order item.
type NewOrderItemAddon struct {
	AddonID  string `json:"addonId" validate:"required"`
	Quantity int    `json:"quantity" validate:"required,min=1"`
}

// NewDeliveryAddress defines the data needed to add a delivery address.
type NewDeliveryAddress struct {
	Street               string   `json:"street" validate:"required"`
	City                 string   `json:"city" validate:"required"`
	State                string   `json:"state" validate:"required"`
	PostalCode           string   `json:"postalCode" validate:"required"`
	DeliveryInstructions string   `json:"deliveryInstructions"`
	Latitude             *float64 `json:"latitude" validate:"required,latitude"`
	Longitude            *float64 `json:"longitude" validate:"required,longitude"`
}

// Decode implements the decoder interface.
func (app *NewOrder) Decode(data []byte) error {
	return json.Unmarshal(data, app)
}

// Validate checks the data in the model is considered clean.
func (app NewOrder) Validate() error {
	if err := errs.Check(app); err != nil {
		return fmt.Errorf("validate: %w", err)
	}

	return nil
}

func toBusNewOrder(app NewOrder) (orderbus.NewOrder, error) {
	restaurantID, err := uuid.Parse(app.RestaurantID)
	if err != nil {
		return orderbus.NewOrder{}, fmt.Errorf("parse restaurantID: %w", err)
	}

	items := make([]orderbus.NewOrderItem, len(app.Items))
	for i, item := range app.Items {
		modifiers := make([]orderbus.NewOrderItemModifier, len(item.Modifiers))
		for j, mod := range item.Modifiers {
			modifiers[j] = orderbus.NewOrderItemModifier{
				ModifierGroupID:  mod.ModifierGroupID,
				ModifierOptionID: mod.ModifierOptionID,
			}
		}

		addons := make([]orderbus.NewOrderItemAddon, len(item.Addons))
		for j, addon := range item.Addons {
			addons[j] = orderbus.NewOrderItemAddon{
				AddonID:  addon.AddonID,
				Quantity: addon.Quantity,
			}
		}

		items[i] = orderbus.NewOrderItem{
			MenuItemID:          item.MenuItemID,
			Quantity:            item.Quantity,
			SpecialInstructions: item.SpecialInstructions,
			Modifiers:           modifiers,
			Addons:              addons,
		}
	}

	var deliveryAddr *orderbus.NewDeliveryAddress
	if app.DeliveryAddress != nil {
		deliveryAddr = &orderbus.NewDeliveryAddress{
			Street:               app.DeliveryAddress.Street,
			City:                 app.DeliveryAddress.City,
			State:                app.DeliveryAddress.State,
			PostalCode:           app.DeliveryAddress.PostalCode,
			DeliveryInstructions: app.DeliveryAddress.DeliveryInstructions,
			Latitude:             app.DeliveryAddress.Latitude,
			Longitude:            app.DeliveryAddress.Longitude,
		}
	}

	return orderbus.NewOrder{
		RestaurantID:        restaurantID.String(),
		CustomerName:        app.CustomerName,
		CustomerEmail:       app.CustomerEmail,
		CustomerPhone:       app.CustomerPhone,
		OrderType:           app.OrderType,
		PaymentMethod:       app.PaymentMethod,
		PromoCode:           app.PromoCode,
		Items:               items,
		DeliveryAddress:     deliveryAddr,
		SpecialInstructions: app.SpecialInstructions,
	}, nil
}

// UpdateOrderStatus defines the data needed to update an order's status.
type UpdateOrderStatus struct {
	OrderStatus   string `json:"orderStatus" validate:"omitempty,oneof=pending confirmed preparing ready out_for_delivery completed cancelled"`
	PaymentStatus string `json:"paymentStatus" validate:"omitempty,oneof=pending processing paid failed refunded"`
}

// Decode implements the decoder interface.
func (app *UpdateOrderStatus) Decode(data []byte) error {
	return json.Unmarshal(data, app)
}

// Validate checks the data in the model is considered clean.
func (app UpdateOrderStatus) Validate() error {
	if err := errs.Check(app); err != nil {
		return fmt.Errorf("validate: %w", err)
	}

	return nil
}

func toBusUpdateOrderStatus(app UpdateOrderStatus) orderbus.UpdateOrderStatus {
	return orderbus.UpdateOrderStatus{
		OrderStatus:   app.OrderStatus,
		PaymentStatus: app.PaymentStatus,
	}
}

// DeliveryQuote represents delivery fee calculation results.
type DeliveryQuote struct {
	DistanceKm            float64 `json:"distanceKm"`
	DeliveryFee           float64 `json:"deliveryFee"`
	MaxDeliveryDistanceKm float64 `json:"maxDeliveryDistanceKm"`
	WithinLimit           bool    `json:"withinLimit"`
}

// Encode implements the encoder interface.
func (app DeliveryQuote) Encode() ([]byte, string, error) {
	data, err := json.Marshal(app)
	return data, "application/json", err
}

// ToAppDeliveryQuote converts a business layer delivery quote to an app layer delivery quote.
func ToAppDeliveryQuote(bus orderbus.DeliveryQuote) DeliveryQuote {
	return DeliveryQuote{
		DistanceKm:            bus.DistanceKm,
		DeliveryFee:           bus.DeliveryFee.Value(),
		MaxDeliveryDistanceKm: bus.MaxDeliveryDistanceKm,
		WithinLimit:           bus.WithinLimit,
	}
}
