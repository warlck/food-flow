# Ordering & Checkout Flow Implementation Plan

**Project:** Food Flow  
**Feature:** Complete Order Management with Stripe Payment Processing  
**Created:** November 6, 2025  
**Status:** 🟡 Planning

---

## Table of Contents
1. [Overview](#overview)
2. [Phase 1: Business Domain Models](#phase-1-business-domain-models)
3. [Phase 2: Application Domain Layer](#phase-2-application-domain-layer)
4. [Phase 3: Database Schema & Migrations](#phase-3-database-schema--migrations)
5. [Phase 4: REST API Implementation](#phase-4-rest-api-implementation)
6. [Phase 5: Stripe Payment Integration](#phase-5-stripe-payment-integration)
7. [Phase 6: Frontend Integration](#phase-6-frontend-integration)
8. [Phase 7: Testing & Deployment](#phase-7-testing--deployment)
9. [Phase 8: Monitoring & Error Handling](#phase-8-monitoring--error-handling)

---

## Overview

### Goals
- ✅ Implement complete order management system
- ✅ Integrate Stripe for online payment processing
- ✅ Support both delivery and pickup orders
- ✅ Handle "pay at location" option
- ✅ Secure payment processing with webhook handling
- ✅ Real-time order status tracking
- ✅ Follow clean architecture principles (business → app → api)

### Tech Stack
- **Backend:** Go 1.24.3, PostgreSQL
- **Frontend:** React 18, TypeScript, Vite
- **Payment:** Stripe API (Go SDK + Stripe.js)
- **Infrastructure:** Kubernetes, Docker

### Architecture Layers
```
┌─────────────────────────────────────┐
│     REST API (app/domain/orderapi)  │  ← Phase 4
├─────────────────────────────────────┤
│  Application Layer (app/domain/*)   │  ← Phase 2
├─────────────────────────────────────┤
│   Business Layer (business/domain)  │  ← Phase 1
├─────────────────────────────────────┤
│        Database (PostgreSQL)        │  ← Phase 3
└─────────────────────────────────────┘
```

### Payment Flow with Stripe
```
1. Customer adds items to cart
   ↓
2. Customer proceeds to checkout
   ↓
3. Backend creates Order (status: "pending")
   ↓
4. Backend creates Stripe PaymentIntent
   ↓
5. Backend returns clientSecret to frontend
   ↓
6. Frontend loads Stripe.js and collects payment
   ↓
7. Customer enters card details
   ↓
8. Frontend confirms payment with Stripe
   ↓
9. Stripe processes payment
   ↓
10. Stripe sends webhook to backend
    ↓
11. Backend updates order status (payment_status: "paid")
    ↓
12. Customer sees order confirmation
```

---

## Phase 1: Business Domain Models

**Status:** ⬜ Not Started  
**Directory:** `business/domain/orderbus/`  
**Estimated Time:** 2-3 days

### 1.1 Core Domain Constants

#### Tasks
- [ ] Define order status constants
- [ ] Define payment status constants
- [ ] Define order type constants
- [ ] Define payment method constants

#### Constants Definition
```go
// File: business/domain/orderbus/model.go

package orderbus

// Order statuses - lifecycle of an order
const (
    OrderStatusPending    = "pending"      // Order created, payment not confirmed
    OrderStatusConfirmed  = "confirmed"    // Payment confirmed, awaiting preparation
    OrderStatusPreparing  = "preparing"    // Restaurant is preparing the order
    OrderStatusReady      = "ready"        // Order ready for pickup/delivery
    OrderStatusCompleted  = "completed"    // Order delivered/picked up
    OrderStatusCancelled  = "cancelled"    // Order cancelled
)

// Payment statuses - payment lifecycle
const (
    PaymentStatusPending    = "pending"      // Payment not yet initiated
    PaymentStatusProcessing = "processing"   // Payment being processed by Stripe
    PaymentStatusPaid       = "paid"         // Payment successful
    PaymentStatusFailed     = "failed"       // Payment failed
    PaymentStatusRefunded   = "refunded"     // Payment refunded
)

// Order types
const (
    OrderTypeDelivery = "delivery"
    OrderTypePickup   = "pickup"
)

// Payment methods
const (
    PaymentMethodCreditCard    = "creditCard"
    PaymentMethodPayAtLocation = "payAtLocation"
)
```

---

### 1.2 Domain Models

#### Tasks
- [ ] Create `Order` domain model
- [ ] Create `OrderItem` domain model
- [ ] Create `DeliveryAddress` domain model
- [ ] Create `NewOrder` input model
- [ ] Create `UpdateOrder` input model
- [ ] Add validation helpers

#### Order Model
```go
// File: business/domain/orderbus/model.go

package orderbus

import "time"

// Order represents a customer order in the system
type Order struct {
    ID                   string           // Unique order identifier
    RestaurantID         string           // Restaurant fulfilling the order
    CustomerName         string           // Customer's name
    CustomerEmail        string           // Customer's email
    CustomerPhone        string           // Customer's phone number
    OrderType            string           // "delivery" or "pickup"
    OrderStatus          string           // Current order status
    PaymentStatus        string           // Current payment status
    PaymentMethod        string           // How customer will pay
    Subtotal             float64          // Sum of all items before fees/tax
    DeliveryFee          float64          // Delivery fee (0 for pickup)
    Tax                  float64          // Tax amount
    Total                float64          // Final total amount
    SpecialInstructions  string           // Order-level instructions
    StripePaymentIntentID string          // Stripe PaymentIntent ID
    Items                []OrderItem      // Items in the order
    DeliveryAddress      *DeliveryAddress // Delivery address (nil for pickup)
    DateCreated          time.Time        // When order was created
    DateUpdated          time.Time        // Last update time
}

// OrderItem represents a menu item within an order
type OrderItem struct {
    ID                   string    // Unique order item identifier
    OrderID              string    // Parent order ID
    MenuItemID           string    // Reference to menu item
    MenuItemName         string    // Snapshot of item name
    MenuItemPrice        float64   // Snapshot of item price
    Quantity             int       // Quantity ordered
    SpecialInstructions  string    // Item-specific instructions
    DateCreated          time.Time // When item was added
}

// DeliveryAddress represents a delivery address for an order
type DeliveryAddress struct {
    ID                    string    // Unique address identifier
    OrderID               string    // Parent order ID
    Street                string    // Street address
    City                  string    // City
    State                 string    // State/Province
    PostalCode            string    // Postal/ZIP code
    DeliveryInstructions  string    // Delivery-specific instructions
    DateCreated           time.Time // When address was added
}
```

---

### 1.3 Input Models (for creating/updating orders)

#### Tasks
- [ ] Create `NewOrder` struct
- [ ] Create `NewOrderItem` struct
- [ ] Create `NewDeliveryAddress` struct
- [ ] Create `UpdateOrderStatus` struct

#### Input Models
```go
// File: business/domain/orderbus/model.go

// NewOrder contains data for creating a new order
type NewOrder struct {
    RestaurantID         string
    CustomerName         string
    CustomerEmail        string
    CustomerPhone        string
    OrderType            string // "delivery" or "pickup"
    PaymentMethod        string // "creditCard" or "payAtLocation"
    Items                []NewOrderItem
    DeliveryAddress      *NewDeliveryAddress
    SpecialInstructions  string
}

// NewOrderItem contains data for adding an item to an order
type NewOrderItem struct {
    MenuItemID           string
    Quantity             int
    SpecialInstructions  string
}

// NewDeliveryAddress contains delivery address data
type NewDeliveryAddress struct {
    Street               string
    City                 string
    State                string
    PostalCode           string
    DeliveryInstructions string
}

// UpdateOrderStatus contains data for updating order status
type UpdateOrderStatus struct {
    OrderStatus   string
    PaymentStatus string
}
```

---

### 1.4 Business Interface

#### Tasks
- [ ] Define `Core` interface with all business operations
- [ ] Define `Storer` interface for data persistence
- [ ] Add filter and query interfaces

#### Business Interface
```go
// File: business/domain/orderbus/orderbus.go

package orderbus

import (
    "context"
    
    "github.com/warlck/food-flow/business/sdk/order"
    "github.com/warlck/food-flow/business/sdk/page"
)

// Core manages order business logic
type Core interface {
    // Create creates a new order
    Create(ctx context.Context, no NewOrder) (Order, error)
    
    // Query retrieves orders with filtering and pagination
    Query(ctx context.Context, filter QueryFilter, orderBy order.By, page page.Page) ([]Order, error)
    
    // QueryByID retrieves a single order by ID
    QueryByID(ctx context.Context, orderID string) (Order, error)
    
    // UpdateStatus updates order and payment status
    UpdateStatus(ctx context.Context, orderID string, status UpdateOrderStatus) error
    
    // UpdateStripePaymentIntent updates the Stripe PaymentIntent ID
    UpdateStripePaymentIntent(ctx context.Context, orderID string, paymentIntentID string) error
    
    // Cancel cancels an order
    Cancel(ctx context.Context, orderID string) error
    
    // Count returns total number of orders matching filter
    Count(ctx context.Context, filter QueryFilter) (int, error)
}

// Storer interface for data persistence
type Storer interface {
    Create(ctx context.Context, order Order) error
    Query(ctx context.Context, filter QueryFilter, orderBy order.By, page page.Page) ([]Order, error)
    QueryByID(ctx context.Context, orderID string) (Order, error)
    Update(ctx context.Context, order Order) error
    Delete(ctx context.Context, orderID string) error
    Count(ctx context.Context, filter QueryFilter) (int, error)
}
```

---


### 1.5 Query Filter & Order By

#### Tasks
- [ ] Create `QueryFilter` struct
- [ ] Create order by definitions
- [ ] Implement filter parsing from URL parameters

```go
// File: business/domain/orderbus/filter.go

package orderbus

import "time"

// QueryFilter holds filter criteria for querying orders
type QueryFilter struct {
    ID              *string    
    RestaurantID    *string    
    CustomerEmail   *string    
    OrderStatus     *string    
    PaymentStatus   *string    
    OrderType       *string    
    StartDate       *time.Time 
    EndDate         *time.Time 
}
```

```go
// File: business/domain/orderbus/order.go

package orderbus

import "github.com/warlck/food-flow/business/sdk/order"

// Order by field definitions
const (
    OrderByID          = "order_id"
    OrderByDateCreated = "date_created"
    OrderByTotal       = "total"
    OrderByStatus      = "order_status"
)

var defaultOrderBy = order.NewBy(OrderByDateCreated, order.DESC)

// Map of valid order by fields
var orderByFields = map[string]string{
    "order_id":     OrderByID,
    "date_created": OrderByDateCreated,
    "total":        OrderByTotal,
    "status":       OrderByStatus,
}
```

---

### 1.6 Testing

#### Tasks
- [ ] Write unit tests for business logic
- [ ] Test order creation with various scenarios
- [ ] Test status updates
- [ ] Test validation functions

```go
// File: business/domain/orderbus/orderbus_test.go

package orderbus_test

import (
    "context"
    "testing"
    
    "github.com/warlck/food-flow/business/domain/orderbus"
)

func TestCreate(t *testing.T) {
    // Test order creation
}

func TestUpdateStatus(t *testing.T) {
    // Test status updates
}

func TestValidation(t *testing.T) {
    // Test validation functions
}
```

---

**Phase 1 Summary:**
- ✅ Domain models defined
- ✅ Business logic interface created
- ✅ Validation helpers implemented
- ✅ Ready for application layer

