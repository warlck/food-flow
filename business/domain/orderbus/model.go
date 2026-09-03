// Package orderbus provides business access to order domain.
package orderbus

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/warlck/food-flow/business/types/money"
)

// Set of error variables for order business logic.
var (
	ErrNotFound                  = errors.New("order not found")
	ErrEmptyItems                = errors.New("order must have at least one item")
	ErrInvalidDeliveryAddress    = errors.New("delivery address is required for delivery orders")
	ErrRestaurantNotFound        = errors.New("restaurant not found")
	ErrPromoNotFound             = errors.New("promo code not found")
	ErrPromoExpired              = errors.New("promo code is expired")
	ErrPromoMinSubtotalNotMet    = errors.New("promo code minimum subtotal not met")
	ErrPromoUsageLimitExceeded   = errors.New("promo code usage limit exceeded")
	ErrMenuItemUnavailable       = errors.New("menu item is unavailable")
	ErrModifierGroupNotFound     = errors.New("modifier group not found")
	ErrModifierOptionNotFound    = errors.New("modifier option not found")
	ErrModifierOptionForeign     = errors.New("modifier option is not valid for this item")
	ErrModifierRequired          = errors.New("modifier selection required")
	ErrModifierOptionUnavailable = errors.New("modifier option is unavailable")
	ErrModifierGroupUnavailable  = errors.New("modifier group is unavailable")
	ErrModifierSelectionLimit    = errors.New("modifier selection limit exceeded")
	ErrAddonNotAssigned           = errors.New("addon is not assigned to menu item")
	ErrAddonUnavailable           = errors.New("addon is unavailable")
	ErrAddonQuantityOutOfRange    = errors.New("addon quantity out of range")
	ErrDuplicateAddon             = errors.New("duplicate addon ID")
	ErrMenuItemRestaurantMismatch = errors.New("menu item does not belong to restaurant")
	ErrAddonRestaurantMismatch    = errors.New("addon does not belong to restaurant")
)

// Order statuses - lifecycle of an order
const (
	OrderStatusPending        = "pending"          // Order created, payment not confirmed
	OrderStatusConfirmed      = "confirmed"        // Payment confirmed, awaiting preparation
	OrderStatusPreparing      = "preparing"        // Restaurant is preparing the order
	OrderStatusReady          = "ready"            // Order ready for pickup/delivery
	OrderStatusOutForDelivery = "out_for_delivery" // Order is on its way to the customer (delivery orders only)
	OrderStatusCompleted      = "completed"        // Order delivered/picked up
	OrderStatusCancelled      = "cancelled"        // Order cancelled
)

// Payment statuses - payment lifecycle
const (
	PaymentStatusPending    = "pending"    // Payment not yet initiated
	PaymentStatusProcessing = "processing" // Payment being processed by Stripe
	PaymentStatusPaid       = "paid"       // Payment successful
	PaymentStatusFailed     = "failed"     // Payment failed
	PaymentStatusRefunded   = "refunded"   // Payment refunded
)

// Order types
const (
	OrderTypeDelivery = "delivery"
	OrderTypePickup   = "pickup"
)

// Payment methods
const (
	PaymentMethodCreditCard = "creditCard"
)

// =============================================================================

// Order represents a customer order in the system
type Order struct {
	ID                    uuid.UUID        // Unique order identifier
	RestaurantID          uuid.UUID        // Restaurant fulfilling the order
	CustomerName          string           // Customer's name
	CustomerEmail         string           // Customer's email
	CustomerPhone         string           // Customer's phone number
	OrderType             string           // "delivery" or "pickup"
	OrderStatus           string           // Current order status
	PaymentStatus         string           // Current payment status
	PaymentMethod         string           // How customer will pay
	PromoCode             string           // Promo code applied (if any)
	Subtotal              money.Money      // Sum of all items before fees/tax
	Discount              money.Money      // Discount amount applied
	DeliveryFee           money.Money      // Delivery fee (0 for pickup)
	Tax                   money.Money      // Tax amount
	Total                 money.Money      // Final total amount
	SpecialInstructions   string           // Order-level instructions
	StripePaymentIntentID string           // Stripe PaymentIntent ID
	Items                 []OrderItem      // Items in the order
	DeliveryAddress       *DeliveryAddress // Delivery address (nil for pickup)
	DateCreated           time.Time        // When order was created
	DateUpdated           time.Time        // Last update time
}

// OrderItem represents a menu item within an order
type OrderItem struct {
	ID                  uuid.UUID           // Unique order item identifier
	OrderID             uuid.UUID           // Parent order ID
	CategoryID          uuid.UUID           // Snapshot of category ID
	CategoryName        string              // Snapshot of category name
	MenuItemID          uuid.UUID           // Reference to menu item
	MenuItemName        string              // Snapshot of item name
	MenuItemPrice       money.Money         // Snapshot of item price
	Quantity            int                 // Quantity ordered
	SpecialInstructions string              // Item-specific instructions
	Modifiers           []OrderItemModifier // Modifiers applied to this item
	Addons              []OrderItemAddon    // Addons applied to this item
	DateCreated         time.Time           // When item was added
}

// OrderItemModifier represents a modifier option selected on an order item
type OrderItemModifier struct {
	ID                 uuid.UUID   // Unique order item modifier identifier
	OrderItemID        uuid.UUID   // Parent order item ID
	ModifierGroupID    uuid.UUID   // Reference to modifier group
	ModifierGroupName  string      // Snapshot of modifier group name
	ModifierOptionID   uuid.UUID   // Reference to modifier option
	ModifierOptionName string      // Snapshot of modifier option name
	PriceDelta         money.Money // Snapshot of modifier option price delta
	DateCreated        time.Time   // When modifier was added
}

// OrderItemAddon represents an addon applied to an order item
type OrderItemAddon struct {
	ID          uuid.UUID   // Unique order item addon identifier
	OrderItemID uuid.UUID   // Parent order item ID
	AddonID     uuid.UUID   // Reference to addon
	AddonName   string      // Snapshot of addon name
	AddonPrice  money.Money // Snapshot of addon price
	Quantity    int         // Quantity of this addon per menu item
	DateCreated time.Time   // When addon was added
}

// DeliveryAddress represents a delivery address for an order
type DeliveryAddress struct {
	ID                   uuid.UUID // Unique address identifier
	OrderID              uuid.UUID // Parent order ID
	Street               string    // Street address
	City                 string    // City
	State                string    // State/Province
	PostalCode           string    // Postal/ZIP code
	DeliveryInstructions string    // Delivery-specific instructions
	Latitude             *float64  // Destination latitude
	Longitude            *float64  // Destination longitude
	DateCreated          time.Time // When address was added
}

// =============================================================================

// NewOrder contains data for creating a new order
type NewOrder struct {
	RestaurantID        string
	CustomerName        string
	CustomerEmail       string
	CustomerPhone       string
	OrderType           string // "delivery" or "pickup"
	PaymentMethod       string // "creditCard"
	PromoCode           string // Promo code applied
	Items               []NewOrderItem
	DeliveryAddress     *NewDeliveryAddress
	SpecialInstructions string
}

// NewOrderItem contains data for adding an item to an order
type NewOrderItem struct {
	MenuItemID          string
	Quantity            int
	SpecialInstructions string
	Modifiers           []NewOrderItemModifier
	Addons              []NewOrderItemAddon
}

// NewOrderItemModifier contains data for adding a modifier option to an order item
type NewOrderItemModifier struct {
	ModifierGroupID  string
	ModifierOptionID string
}

// NewOrderItemAddon contains data for adding an addon to an order item
type NewOrderItemAddon struct {
	AddonID  string
	Quantity int
}

// NewDeliveryAddress contains delivery address data
type NewDeliveryAddress struct {
	Street               string
	City                 string
	State                string
	PostalCode           string
	DeliveryInstructions string
	Latitude             *float64
	Longitude            *float64
}

// UpdateOrderStatus contains data for updating order status
type UpdateOrderStatus struct {
	OrderStatus   string
	PaymentStatus string
}
