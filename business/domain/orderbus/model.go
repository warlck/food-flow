// Package orderbus provides business access to order domain.
package orderbus

import (
	"time"

	"github.com/google/uuid"
	"github.com/warlck/food-flow/business/types/money"
)

// Order statuses - lifecycle of an order
const (
	OrderStatusPending   = "pending"   // Order created, payment not confirmed
	OrderStatusConfirmed = "confirmed" // Payment confirmed, awaiting preparation
	OrderStatusPreparing = "preparing" // Restaurant is preparing the order
	OrderStatusReady     = "ready"     // Order ready for pickup/delivery
	OrderStatusCompleted = "completed" // Order delivered/picked up
	OrderStatusCancelled = "cancelled" // Order cancelled
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
	Subtotal              money.Money      // Sum of all items before fees/tax
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
	ID                  uuid.UUID   // Unique order item identifier
	OrderID             uuid.UUID   // Parent order ID
	MenuItemID          uuid.UUID   // Reference to menu item
	MenuItemName        string      // Snapshot of item name
	MenuItemPrice       money.Money // Snapshot of item price
	Quantity            int         // Quantity ordered
	SpecialInstructions string      // Item-specific instructions
	DateCreated         time.Time   // When item was added
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
	Items               []NewOrderItem
	DeliveryAddress     *NewDeliveryAddress
	SpecialInstructions string
}

// NewOrderItem contains data for adding an item to an order
type NewOrderItem struct {
	MenuItemID          string
	Quantity            int
	SpecialInstructions string
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
