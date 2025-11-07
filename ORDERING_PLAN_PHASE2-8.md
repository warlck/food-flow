# Implementation Plan - Phases 2-8

## Phase 2: Application Domain Layer

**Status:** ⬜ Not Started  
**Directory:** `app/domain/orderapp/`  
**Estimated Time:** 2 days

This document continues from Phase 1 in the main ORDERING_IMPLEMENTATION_PLAN.md

---

## Phase 3: Database Schema & Migrations

**Status:** ⬜ Not Started  
**File:** `business/sdk/migrate/sql/migrate.sql`  
**Estimated Time:** 1-2 days

### 3.1 Database Tables

#### Tasks
- [ ] Create `orders` table
- [ ] Create `order_items` table
- [ ] Create `delivery_addresses` table  
- [ ] Create `payment_transactions` table
- [ ] Add indexes for performance
- [ ] Write rollback migration

#### Orders Table
```sql
-- Version: 1.05
-- Description: Create table orders
CREATE TABLE orders (
    order_id              UUID           NOT NULL,
    restaurant_id         UUID           NOT NULL,
    customer_name         TEXT           NOT NULL,
    customer_email        TEXT           NOT NULL,
    customer_phone        TEXT           NOT NULL,
    order_type            TEXT           NOT NULL,
    order_status          TEXT           NOT NULL,
    payment_status        TEXT           NOT NULL,
    payment_method        TEXT           NOT NULL,
    subtotal              NUMERIC(10, 2) NOT NULL,
    delivery_fee          NUMERIC(10, 2) NOT NULL DEFAULT 0,
    tax                   NUMERIC(10, 2) NOT NULL,
    total                 NUMERIC(10, 2) NOT NULL,
    special_instructions  TEXT           NULL,
    stripe_payment_intent_id TEXT        NULL,
    date_created          TIMESTAMP      NOT NULL,
    date_updated          TIMESTAMP      NOT NULL,

    PRIMARY KEY (order_id),
    FOREIGN KEY (restaurant_id) REFERENCES restaurants(restaurant_id) ON DELETE CASCADE
);
```

#### Order Items Table
```sql
-- Version: 1.06
-- Description: Create table order_items
CREATE TABLE order_items (
    order_item_id        UUID           NOT NULL,
    order_id             UUID           NOT NULL,
    menu_item_id         UUID           NOT NULL,
    menu_item_name       TEXT           NOT NULL,
    menu_item_price      NUMERIC(10, 2) NOT NULL,
    quantity             INT            NOT NULL,
    special_instructions TEXT           NULL,
    date_created         TIMESTAMP      NOT NULL,

    PRIMARY KEY (order_item_id),
    FOREIGN KEY (order_id) REFERENCES orders(order_id) ON DELETE CASCADE,
    FOREIGN KEY (menu_item_id) REFERENCES menu_items(menu_item_id)
);
```

#### Delivery Addresses Table
```sql
-- Version: 1.07
-- Description: Create table delivery_addresses
CREATE TABLE delivery_addresses (
    address_id            UUID      NOT NULL,
    order_id              UUID      NOT NULL UNIQUE,
    street                TEXT      NOT NULL,
    city                  TEXT      NOT NULL,
    state                 TEXT      NOT NULL,
    postal_code           TEXT      NOT NULL,
    delivery_instructions TEXT      NULL,
    date_created          TIMESTAMP NOT NULL,

    PRIMARY KEY (address_id),
    FOREIGN KEY (order_id) REFERENCES orders(order_id) ON DELETE CASCADE
);
```

#### Payment Transactions Table
```sql
-- Version: 1.08
-- Description: Create table payment_transactions
CREATE TABLE payment_transactions (
    transaction_id          UUID           NOT NULL,
    order_id                UUID           NOT NULL,
    stripe_payment_intent_id TEXT          NULL,
    stripe_charge_id        TEXT           NULL,
    payment_method          TEXT           NOT NULL,
    amount                  NUMERIC(10, 2) NOT NULL,
    currency                TEXT           NOT NULL DEFAULT 'usd',
    status                  TEXT           NOT NULL,
    error_message           TEXT           NULL,
    metadata                JSONB          NULL,
    date_created            TIMESTAMP      NOT NULL,
    date_updated            TIMESTAMP      NOT NULL,

    PRIMARY KEY (transaction_id),
    FOREIGN KEY (order_id) REFERENCES orders(order_id) ON DELETE CASCADE
);
```

#### Performance Indexes
```sql
-- Version: 1.09
-- Description: Create indexes for orders
CREATE INDEX idx_orders_restaurant_id ON orders(restaurant_id);
CREATE INDEX idx_orders_customer_email ON orders(customer_email);
CREATE INDEX idx_orders_order_status ON orders(order_status);
CREATE INDEX idx_orders_payment_status ON orders(payment_status);
CREATE INDEX idx_orders_date_created ON orders(date_created DESC);
CREATE INDEX idx_order_items_order_id ON order_items(order_id);
CREATE INDEX idx_delivery_addresses_order_id ON delivery_addresses(order_id);
CREATE INDEX idx_payment_transactions_order_id ON payment_transactions(order_id);
CREATE INDEX idx_payment_transactions_stripe_id ON payment_transactions(stripe_payment_intent_id);
```

---

### 3.2 Database Store Implementation

#### Tasks
- [ ] Create `business/domain/orderbus/stores/orderdb/` directory
- [ ] Implement `orderdb.go` - main store
- [ ] Implement `model.go` - DB models
- [ ] Implement `filter.go` - filter SQL generation
- [ ] Write tests for store operations

#### Store Structure
```go
// File: business/domain/orderbus/stores/orderdb/orderdb.go

package orderdb

import (
    "context"
    "database/sql"
    "fmt"
    
    "github.com/jmoiron/sqlx"
    "github.com/warlck/food-flow/business/domain/orderbus"
    "github.com/warlck/food-flow/business/sdk/order"
    "github.com/warlck/food-flow/business/sdk/page"
)

type Store struct {
    db *sqlx.DB
}

func NewStore(db *sqlx.DB) *Store {
    return &Store{
        db: db,
    }
}

// Create inserts a new order and related items
func (s *Store) Create(ctx context.Context, order orderbus.Order) error {
    // Use transaction to insert order, items, and address
    tx, err := s.db.BeginTxx(ctx, nil)
    if err != nil {
        return fmt.Errorf("begin transaction: %w", err)
    }
    defer tx.Rollback()
    
    // Insert order
    const orderSQL = `
    INSERT INTO orders (
        order_id, restaurant_id, customer_name, customer_email, customer_phone,
        order_type, order_status, payment_status, payment_method,
        subtotal, delivery_fee, tax, total, special_instructions,
        stripe_payment_intent_id, date_created, date_updated
    ) VALUES (
        :order_id, :restaurant_id, :customer_name, :customer_email, :customer_phone,
        :order_type, :order_status, :payment_status, :payment_method,
        :subtotal, :delivery_fee, :tax, :total, :special_instructions,
        :stripe_payment_intent_id, :date_created, :date_updated
    )`
    
    if _, err := tx.NamedExecContext(ctx, orderSQL, toDBOrder(order)); err != nil {
        return fmt.Errorf("insert order: %w", err)
    }
    
    // Insert order items
    for _, item := range order.Items {
        const itemSQL = `
        INSERT INTO order_items (
            order_item_id, order_id, menu_item_id, menu_item_name,
            menu_item_price, quantity, special_instructions, date_created
        ) VALUES (
            :order_item_id, :order_id, :menu_item_id, :menu_item_name,
            :menu_item_price, :quantity, :special_instructions, :date_created
        )`
        
        if _, err := tx.NamedExecContext(ctx, itemSQL, toDBOrderItem(item, order.ID)); err != nil {
            return fmt.Errorf("insert order item: %w", err)
        }
    }
    
    // Insert delivery address if present
    if order.DeliveryAddress != nil {
        const addressSQL = `
        INSERT INTO delivery_addresses (
            address_id, order_id, street, city, state, postal_code,
            delivery_instructions, date_created
        ) VALUES (
            :address_id, :order_id, :street, :city, :state, :postal_code,
            :delivery_instructions, :date_created
        )`
        
        if _, err := tx.NamedExecContext(ctx, addressSQL, toDBDeliveryAddress(*order.DeliveryAddress, order.ID)); err != nil {
            return fmt.Errorf("insert delivery address: %w", err)
        }
    }
    
    return tx.Commit()
}

// Query retrieves orders with filtering
func (s *Store) Query(ctx context.Context, filter orderbus.QueryFilter, orderBy order.By, pg page.Page) ([]orderbus.Order, error) {
    // Implementation with JOIN queries
    return nil, nil
}

// QueryByID retrieves a single order with all related data
func (s *Store) QueryByID(ctx context.Context, orderID string) (orderbus.Order, error) {
    // Implementation with JOIN queries
    return orderbus.Order{}, nil
}

// Update updates an existing order
func (s *Store) Update(ctx context.Context, order orderbus.Order) error {
    const sql = `
    UPDATE orders SET
        order_status = :order_status,
        payment_status = :payment_status,
        stripe_payment_intent_id = :stripe_payment_intent_id,
        date_updated = :date_updated
    WHERE order_id = :order_id`
    
    if _, err := s.db.NamedExecContext(ctx, sql, toDBOrder(order)); err != nil {
        return fmt.Errorf("update order: %w", err)
    }
    
    return nil
}

// Delete soft-deletes or cancels an order
func (s *Store) Delete(ctx context.Context, orderID string) error {
    const sql = `
    UPDATE orders SET
        order_status = 'cancelled',
        date_updated = NOW()
    WHERE order_id = $1`
    
    if _, err := s.db.ExecContext(ctx, sql, orderID); err != nil {
        return fmt.Errorf("delete order: %w", err)
    }
    
    return nil
}

// Count returns total number of orders matching filter
func (s *Store) Count(ctx context.Context, filter orderbus.QueryFilter) (int, error) {
    // Implementation
    return 0, nil
}
```

---

### 3.3 Database Models

#### Tasks
- [ ] Create DB-specific models with tags
- [ ] Implement conversion between domain and DB models

```go
// File: business/domain/orderbus/stores/orderdb/model.go

package orderdb

import (
    "time"
    
    "github.com/warlck/food-flow/business/domain/orderbus"
)

// dbOrder is the database representation
type dbOrder struct {
    ID                     string    `db:"order_id"`
    RestaurantID           string    `db:"restaurant_id"`
    CustomerName           string    `db:"customer_name"`
    CustomerEmail          string    `db:"customer_email"`
    CustomerPhone          string    `db:"customer_phone"`
    OrderType              string    `db:"order_type"`
    OrderStatus            string    `db:"order_status"`
    PaymentStatus          string    `db:"payment_status"`
    PaymentMethod          string    `db:"payment_method"`
    Subtotal               float64   `db:"subtotal"`
    DeliveryFee            float64   `db:"delivery_fee"`
    Tax                    float64   `db:"tax"`
    Total                  float64   `db:"total"`
    SpecialInstructions    string    `db:"special_instructions"`
    StripePaymentIntentID  string    `db:"stripe_payment_intent_id"`
    DateCreated            time.Time `db:"date_created"`
    DateUpdated            time.Time `db:"date_updated"`
}

type dbOrderItem struct {
    ID                   string    `db:"order_item_id"`
    OrderID              string    `db:"order_id"`
    MenuItemID           string    `db:"menu_item_id"`
    MenuItemName         string    `db:"menu_item_name"`
    MenuItemPrice        float64   `db:"menu_item_price"`
    Quantity             int       `db:"quantity"`
    SpecialInstructions  string    `db:"special_instructions"`
    DateCreated          time.Time `db:"date_created"`
}

type dbDeliveryAddress struct {
    ID                    string    `db:"address_id"`
    OrderID               string    `db:"order_id"`
    Street                string    `db:"street"`
    City                  string    `db:"city"`
    State                 string    `db:"state"`
    PostalCode            string    `db:"postal_code"`
    DeliveryInstructions  string    `db:"delivery_instructions"`
    DateCreated           time.Time `db:"date_created"`
}

// Conversion functions
func toDBOrder(order orderbus.Order) dbOrder {
    return dbOrder{
        ID:                    order.ID,
        RestaurantID:          order.RestaurantID,
        CustomerName:          order.CustomerName,
        CustomerEmail:         order.CustomerEmail,
        CustomerPhone:         order.CustomerPhone,
        OrderType:             order.OrderType,
        OrderStatus:           order.OrderStatus,
        PaymentStatus:         order.PaymentStatus,
        PaymentMethod:         order.PaymentMethod,
        Subtotal:              order.Subtotal,
        DeliveryFee:           order.DeliveryFee,
        Tax:                   order.Tax,
        Total:                 order.Total,
        SpecialInstructions:   order.SpecialInstructions,
        StripePaymentIntentID: order.StripePaymentIntentID,
        DateCreated:           order.DateCreated,
        DateUpdated:           order.DateUpdated,
    }
}

func toBusOrder(dbo dbOrder, items []dbOrderItem, addr *dbDeliveryAddress) orderbus.Order {
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
        Subtotal:              dbo.Subtotal,
        DeliveryFee:           dbo.DeliveryFee,
        Tax:                   dbo.Tax,
        Total:                 dbo.Total,
        SpecialInstructions:   dbo.SpecialInstructions,
        StripePaymentIntentID: dbo.StripePaymentIntentID,
        DateCreated:           dbo.DateCreated,
        DateUpdated:           dbo.DateUpdated,
        Items:                 make([]orderbus.OrderItem, len(items)),
    }
    
    for i, item := range items {
        order.Items[i] = toBusOrderItem(item)
    }
    
    if addr != nil {
        da := toBusDeliveryAddress(*addr)
        order.DeliveryAddress = &da
    }
    
    return order
}

func toDBOrderItem(item orderbus.OrderItem, orderID string) dbOrderItem {
    return dbOrderItem{
        ID:                  item.ID,
        OrderID:             orderID,
        MenuItemID:          item.MenuItemID,
        MenuItemName:        item.MenuItemName,
        MenuItemPrice:       item.MenuItemPrice,
        Quantity:            item.Quantity,
        SpecialInstructions: item.SpecialInstructions,
        DateCreated:         item.DateCreated,
    }
}

func toBusOrderItem(dbo dbOrderItem) orderbus.OrderItem {
    return orderbus.OrderItem{
        ID:                  dbo.ID,
        OrderID:             dbo.OrderID,
        MenuItemID:          dbo.MenuItemID,
        MenuItemName:        dbo.MenuItemName,
        MenuItemPrice:       dbo.MenuItemPrice,
        Quantity:            dbo.Quantity,
        SpecialInstructions: dbo.SpecialInstructions,
        DateCreated:         dbo.DateCreated,
    }
}

func toDBDeliveryAddress(addr orderbus.DeliveryAddress, orderID string) dbDeliveryAddress {
    return dbDeliveryAddress{
        ID:                   addr.ID,
        OrderID:              orderID,
        Street:               addr.Street,
        City:                 addr.City,
        State:                addr.State,
        PostalCode:           addr.PostalCode,
        DeliveryInstructions: addr.DeliveryInstructions,
        DateCreated:          addr.DateCreated,
    }
}

func toBusDeliveryAddress(dbo dbDeliveryAddress) orderbus.DeliveryAddress {
    return orderbus.DeliveryAddress{
        ID:                   dbo.ID,
        OrderID:              dbo.OrderID,
        Street:               dbo.Street,
        City:                 dbo.City,
        State:                dbo.State,
        PostalCode:           dbo.PostalCode,
        DeliveryInstructions: dbo.DeliveryInstructions,
        DateCreated:          dbo.DateCreated,
    }
}
```

---

**Phase 3 Summary:**
- ✅ Database schema defined with 4 tables
- ✅ Indexes added for performance
- ✅ Store implementation with transactions
- ✅ DB model conversions implemented
- ✅ Ready for REST API implementation

---

*To be continued with Phases 4-8 in next steps...*

## Phase 4: REST API Implementation

**Status:** ⬜ Not Started  
**Directory:** `app/domain/orderapi/`  
**Estimated Time:** 3-4 days

### 4.1 API Route Definitions

#### Tasks
- [ ] Create `app/domain/orderapi/` directory
- [ ] Implement `route.go` - define all routes
- [ ] Implement `orderapi.go` - route handlers
- [ ] Add middleware for authentication
- [ ] Add request validation
- [ ] Add error handling

#### Routes Structure
```go
// File: app/domain/orderapi/route.go

package orderapi

import (
    "github.com/warlck/food-flow/app/sdk/mid"
    "github.com/warlck/food-flow/business/domain/orderbus"
    "github.com/warlck/food-flow/foundation/web"
)

type Config struct {
    OrderBus orderbus.Core
}

// Routes adds order routes to the web application
func Routes(app *web.App, cfg Config) {
    const version = "v1"
    
    // Public routes (no auth required)
    app.HandleFunc("POST", version, "/orders", cfg.create)
    app.HandleFunc("GET", version, "/orders/:id", cfg.queryByID)
    
    // Protected routes (require authentication)
    authen := mid.Authenticate(cfg.Auth)
    app.HandleFunc("GET", version, "/orders", cfg.query, authen)
    app.HandleFunc("PATCH", version, "/orders/:id/status", cfg.updateStatus, authen)
    app.HandleFunc("DELETE", version, "/orders/:id", cfg.cancel, authen)
    
    // Payment routes
    app.HandleFunc("POST", version, "/orders/:id/payment/intent", cfg.createPaymentIntent)
    app.HandleFunc("POST", version, "/orders/:id/payment/confirm", cfg.confirmPayment)
    
    // Webhook route (for Stripe)
    app.HandleFunc("POST", version, "/webhooks/stripe", cfg.handleStripeWebhook)
}
```

---

### 4.2 API Handler Implementation

#### Tasks
- [ ] Implement `create` handler
- [ ] Implement `query` and `queryByID` handlers
- [ ] Implement `updateStatus` handler
- [ ] Implement `cancel` handler
- [ ] Implement payment-related handlers
- [ ] Implement webhook handler

#### Create Order Handler
```go
// File: app/domain/orderapi/orderapi.go

package orderapi

import (
    "context"
    "net/http"
    
    "github.com/warlck/food-flow/app/sdk/errs"
    "github.com/warlck/food-flow/business/domain/orderbus"
    "github.com/warlck/food-flow/foundation/web"
)

type api struct {
    orderBus orderbus.Core
}

// create handles creating a new order
func (a *api) create(ctx context.Context, r *http.Request) web.Encoder {
    var req CreateOrderRequest
    if err := web.Decode(r, &req); err != nil {
        return errs.NewError(err)
    }
    
    // Validate delivery address for delivery orders
    if req.OrderType == "delivery" && req.DeliveryAddress == nil {
        return errs.NewFieldErrors("deliveryAddress", "required for delivery orders")
    }
    
    // Convert to business model
    newOrder := toBusOrder(req)
    
    // Create order
    order, err := a.orderBus.Create(ctx, newOrder)
    if err != nil {
        return errs.NewError(err)
    }
    
    return toAppOrder(order)
}

// queryByID retrieves a single order by ID
func (a *api) queryByID(ctx context.Context, r *http.Request) web.Encoder {
    orderID := web.Param(r, "id")
    
    order, err := a.orderBus.QueryByID(ctx, orderID)
    if err != nil {
        return errs.NewError(err)
    }
    
    return toAppOrder(order)
}

// query retrieves orders with filtering and pagination
func (a *api) query(ctx context.Context, r *http.Request) web.Encoder {
    qp := parseQueryParams(r.URL.Query())
    
    // Parse filter
    filter := toBusinessFilter(qp)
    
    // Parse order by
    orderBy, err := order.Parse(orderByFields, qp.OrderBy, defaultOrderBy)
    if err != nil {
        return errs.NewFieldErrors("orderBy", err)
    }
    
    // Parse pagination
    page, err := page.Parse(r)
    if err != nil {
        return errs.NewFieldErrors("page", err)
    }
    
    // Query orders
    orders, err := a.orderBus.Query(ctx, filter, orderBy, page)
    if err != nil {
        return errs.NewError(err)
    }
    
    // Get total count
    total, err := a.orderBus.Count(ctx, filter)
    if err != nil {
        return errs.NewError(err)
    }
    
    // Convert to response
    items := make([]OrderResponse, len(orders))
    for i, order := range orders {
        items[i] = toAppOrder(order)
    }
    
    return page.NewResponse(items, total, page.Number, page.RowsPerPage)
}

// updateStatus updates order or payment status
func (a *api) updateStatus(ctx context.Context, r *http.Request) web.Encoder {
    orderID := web.Param(r, "id")
    
    var req UpdateOrderStatusRequest
    if err := web.Decode(r, &req); err != nil {
        return errs.NewError(err)
    }
    
    status := orderbus.UpdateOrderStatus{
        OrderStatus:   req.OrderStatus,
        PaymentStatus: req.PaymentStatus,
    }
    
    if err := a.orderBus.UpdateStatus(ctx, orderID, status); err != nil {
        return errs.NewError(err)
    }
    
    return web.StatusOK
}

// cancel cancels an order
func (a *api) cancel(ctx context.Context, r *http.Request) web.Encoder {
    orderID := web.Param(r, "id")
    
    if err := a.orderBus.Cancel(ctx, orderID); err != nil {
        return errs.NewError(err)
    }
    
    return web.StatusOK
}
```

---

### 4.3 Register Routes in Sales Service

#### Tasks
- [ ] Update `api/services/sales/build/all/all.go`
- [ ] Add orderbus to dependency injection
- [ ] Import and register orderapi routes

```go
// File: api/services/sales/build/all/all.go

package all

import (
    "github.com/warlck/food-flow/app/domain/orderapi"
    "github.com/warlck/food-flow/business/domain/orderbus"
    "github.com/warlck/food-flow/business/domain/orderbus/stores/orderdb"
    // ... other imports
)

func Routes() add.Routes {
    return func(app *web.App, cfg add.Config) {
        // ... existing routes (checkapp, userapp, etc.)
        
        // Order routes
        orderBus := orderbus.NewBusiness(
            orderdb.NewStore(cfg.DB),
            cfg.MenuItemBus,
            cfg.RestaurantBus,
        )
        orderapi.Routes(app, orderapi.Config{
            OrderBus: orderBus,
        })
    }
}
```

---

### 4.4 API Testing

#### Tasks
- [ ] Create integration tests for all endpoints
- [ ] Test order creation flow
- [ ] Test status updates
- [ ] Test error handling
- [ ] Test validation

```go
// File: app/domain/orderapi/orderapi_test.go

package orderapi_test

import (
    "testing"
    
    "github.com/warlck/food-flow/app/sdk/apitest"
)

func TestCreate(t *testing.T) {
    // Test order creation endpoint
}

func TestQuery(t *testing.T) {
    // Test querying orders
}

func TestUpdateStatus(t *testing.T) {
    // Test status updates
}
```

---

**Phase 4 Summary:**
- ✅ RESTful API endpoints defined
- ✅ Handlers implemented with validation
- ✅ Routes registered in sales service
- ✅ Ready for Stripe integration


## Phase 5: Stripe Payment Integration

**Status:** ⬜ Not Started  
**Directory:** `business/domain/paymentbus/`  
**Estimated Time:** 3-4 days

### 5.1 Stripe SDK Setup

#### Tasks
- [ ] Create Stripe account (test mode)
- [ ] Obtain API keys (publishable and secret)
- [ ] Add Stripe Go SDK: `go get github.com/stripe/stripe-go/v80`
- [ ] Configure environment variables
- [ ] Test API connection

```bash
# Add to go.mod
go get github.com/stripe/stripe-go/v80
```

```yaml
# Environment variables (infra/k8s/base/sales-pod.yaml)
env:
  - name: STRIPE_SECRET_KEY
    valueFrom:
      secretKeyRef:
        name: stripe-secrets
        key: secret-key
  - name: STRIPE_PUBLISHABLE_KEY
    valueFrom:
      secretKeyRef:
        name: stripe-secrets
        key: publishable-key
  - name: STRIPE_WEBHOOK_SECRET
    valueFrom:
      secretKeyRef:
        name: stripe-secrets
        key: webhook-secret
```

---

### 5.2 Payment Business Layer

#### Tasks
- [ ] Create `business/domain/paymentbus/` directory
- [ ] Implement `paymentbus.go` - core logic
- [ ] Implement `stripe.go` - Stripe client wrapper
- [ ] Implement `model.go` - payment models
- [ ] Implement webhook handling
- [ ] Add tests

#### Payment Service Interface
```go
// File: business/domain/paymentbus/paymentbus.go

package paymentbus

import (
    "context"
    
    "github.com/warlck/food-flow/business/domain/orderbus"
)

// Core defines the payment service interface
type Core interface {
    // CreatePaymentIntent creates a Stripe PaymentIntent
    CreatePaymentIntent(ctx context.Context, order orderbus.Order) (PaymentIntent, error)
    
    // ConfirmPayment confirms a payment was successful
    ConfirmPayment(ctx context.Context, paymentIntentID string) error
    
    // RefundPayment refunds a payment
    RefundPayment(ctx context.Context, paymentIntentID string, amount float64) error
    
    // GetPaymentStatus retrieves payment status from Stripe
    GetPaymentStatus(ctx context.Context, paymentIntentID string) (PaymentStatus, error)
    
    // HandleWebhook processes Stripe webhook events
    HandleWebhook(ctx context.Context, payload []byte, signature string) error
}

// PaymentIntent represents a Stripe PaymentIntent
type PaymentIntent struct {
    ID           string
    ClientSecret string
    Amount       int64  // Amount in cents
    Currency     string
    Status       string
    OrderID      string
}

// PaymentStatus represents the current status of a payment
type PaymentStatus struct {
    ID             string
    Status         string
    Amount         int64
    Currency       string
    ErrorMessage   string
    PaymentMethod  string
}

// Config holds Stripe configuration
type Config struct {
    SecretKey      string
    PublishableKey string
    WebhookSecret  string
}
```

---

### 5.3 Stripe Client Implementation

#### Tasks
- [ ] Implement Stripe PaymentIntent creation
- [ ] Implement payment confirmation
- [ ] Implement refund logic
- [ ] Implement webhook signature verification

```go
// File: business/domain/paymentbus/stripe.go

package paymentbus

import (
    "context"
    "fmt"
    
    "github.com/stripe/stripe-go/v80"
    "github.com/stripe/stripe-go/v80/paymentintent"
    "github.com/stripe/stripe-go/v80/refund"
    "github.com/stripe/stripe-go/v80/webhook"
    "github.com/warlck/food-flow/business/domain/orderbus"
)

type business struct {
    config   Config
    orderBus orderbus.Core
}

// NewBusiness creates a new payment business
func NewBusiness(cfg Config, orderBus orderbus.Core) Core {
    stripe.Key = cfg.SecretKey
    return &business{
        config:   cfg,
        orderBus: orderBus,
    }
}

// CreatePaymentIntent creates a new Stripe PaymentIntent
func (b *business) CreatePaymentIntent(ctx context.Context, order orderbus.Order) (PaymentIntent, error) {
    // Convert total to cents (Stripe uses smallest currency unit)
    amountCents := int64(order.Total * 100)
    
    params := &stripe.PaymentIntentParams{
        Amount:   stripe.Int64(amountCents),
        Currency: stripe.String("usd"),
        Params: stripe.Params{
            Metadata: map[string]string{
                "order_id":      order.ID,
                "restaurant_id": order.RestaurantID,
                "customer_email": order.CustomerEmail,
            },
        },
        AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
            Enabled: stripe.Bool(true),
        },
    }
    
    pi, err := paymentintent.New(params)
    if err != nil {
        return PaymentIntent{}, fmt.Errorf("create payment intent: %w", err)
    }
    
    // Update order with PaymentIntent ID
    if err := b.orderBus.UpdateStripePaymentIntent(ctx, order.ID, pi.ID); err != nil {
        return PaymentIntent{}, fmt.Errorf("update order payment intent: %w", err)
    }
    
    return PaymentIntent{
        ID:           pi.ID,
        ClientSecret: pi.ClientSecret,
        Amount:       pi.Amount,
        Currency:     string(pi.Currency),
        Status:       string(pi.Status),
        OrderID:      order.ID,
    }, nil
}

// ConfirmPayment confirms that a payment was successful
func (b *business) ConfirmPayment(ctx context.Context, paymentIntentID string) error {
    pi, err := paymentintent.Get(paymentIntentID, nil)
    if err != nil {
        return fmt.Errorf("get payment intent: %w", err)
    }
    
    // Check if payment succeeded
    if pi.Status != stripe.PaymentIntentStatusSucceeded {
        return fmt.Errorf("payment not successful: %s", pi.Status)
    }
    
    // Get order ID from metadata
    orderID, ok := pi.Metadata["order_id"]
    if !ok {
        return fmt.Errorf("order_id not found in payment intent metadata")
    }
    
    // Update order payment status
    return b.orderBus.UpdateStatus(ctx, orderID, orderbus.UpdateOrderStatus{
        OrderStatus:   orderbus.OrderStatusConfirmed,
        PaymentStatus: orderbus.PaymentStatusPaid,
    })
}

// RefundPayment refunds a payment
func (b *business) RefundPayment(ctx context.Context, paymentIntentID string, amount float64) error {
    amountCents := int64(amount * 100)
    
    params := &stripe.RefundParams{
        PaymentIntent: stripe.String(paymentIntentID),
    }
    
    if amount > 0 {
        params.Amount = stripe.Int64(amountCents)
    }
    
    _, err := refund.New(params)
    if err != nil {
        return fmt.Errorf("create refund: %w", err)
    }
    
    return nil
}

// GetPaymentStatus retrieves the current payment status
func (b *business) GetPaymentStatus(ctx context.Context, paymentIntentID string) (PaymentStatus, error) {
    pi, err := paymentintent.Get(paymentIntentID, nil)
    if err != nil {
        return PaymentStatus{}, fmt.Errorf("get payment intent: %w", err)
    }
    
    return PaymentStatus{
        ID:       pi.ID,
        Status:   string(pi.Status),
        Amount:   pi.Amount,
        Currency: string(pi.Currency),
    }, nil
}

// HandleWebhook processes Stripe webhook events
func (b *business) HandleWebhook(ctx context.Context, payload []byte, signature string) error {
    event, err := webhook.ConstructEvent(payload, signature, b.config.WebhookSecret)
    if err != nil {
        return fmt.Errorf("verify webhook signature: %w", err)
    }
    
    switch event.Type {
    case "payment_intent.succeeded":
        var pi stripe.PaymentIntent
        if err := event.Data.UnmarshalJSON(&pi); err != nil {
            return fmt.Errorf("unmarshal payment intent: %w", err)
        }
        
        // Update order status
        orderID, ok := pi.Metadata["order_id"]
        if !ok {
            return fmt.Errorf("order_id not found in metadata")
        }
        
        return b.orderBus.UpdateStatus(ctx, orderID, orderbus.UpdateOrderStatus{
            OrderStatus:   orderbus.OrderStatusConfirmed,
            PaymentStatus: orderbus.PaymentStatusPaid,
        })
        
    case "payment_intent.payment_failed":
        var pi stripe.PaymentIntent
        if err := event.Data.UnmarshalJSON(&pi); err != nil {
            return fmt.Errorf("unmarshal payment intent: %w", err)
        }
        
        orderID, ok := pi.Metadata["order_id"]
        if !ok {
            return fmt.Errorf("order_id not found in metadata")
        }
        
        return b.orderBus.UpdateStatus(ctx, orderID, orderbus.UpdateOrderStatus{
            PaymentStatus: orderbus.PaymentStatusFailed,
        })
        
    default:
        // Ignore other event types
        return nil
    }
}
```

---

### 5.4 Payment API Handlers

#### Tasks
- [ ] Implement payment intent creation endpoint
- [ ] Implement payment confirmation endpoint
- [ ] Implement Stripe webhook endpoint
- [ ] Add webhook signature verification

```go
// File: app/domain/orderapi/payment.go

package orderapi

import (
    "context"
    "io"
    "net/http"
    
    "github.com/warlck/food-flow/app/sdk/errs"
    "github.com/warlck/food-flow/business/domain/paymentbus"
    "github.com/warlck/food-flow/foundation/web"
)

// createPaymentIntent creates a Stripe PaymentIntent for an order
func (a *api) createPaymentIntent(ctx context.Context, r *http.Request) web.Encoder {
    orderID := web.Param(r, "id")
    
    // Get order
    order, err := a.orderBus.QueryByID(ctx, orderID)
    if err != nil {
        return errs.NewError(err)
    }
    
    // Only create payment intent for credit card payments
    if order.PaymentMethod != "creditCard" {
        return errs.NewError("payment method must be creditCard")
    }
    
    // Create PaymentIntent
    pi, err := a.paymentBus.CreatePaymentIntent(ctx, order)
    if err != nil {
        return errs.NewError(err)
    }
    
    return PaymentIntentResponse{
        ClientSecret: pi.ClientSecret,
        OrderID:      pi.OrderID,
        Amount:       float64(pi.Amount) / 100,
        Currency:     pi.Currency,
    }
}

// confirmPayment confirms a payment was successful
func (a *api) confirmPayment(ctx context.Context, r *http.Request) web.Encoder {
    orderID := web.Param(r, "id")
    
    // Get order
    order, err := a.orderBus.QueryByID(ctx, orderID)
    if err != nil {
        return errs.NewError(err)
    }
    
    if order.StripePaymentIntentID == "" {
        return errs.NewError("no payment intent found for order")
    }
    
    // Confirm payment
    if err := a.paymentBus.ConfirmPayment(ctx, order.StripePaymentIntentID); err != nil {
        return errs.NewError(err)
    }
    
    return web.StatusOK
}

// handleStripeWebhook handles Stripe webhook events
func (a *api) handleStripeWebhook(ctx context.Context, r *http.Request) web.Encoder {
    payload, err := io.ReadAll(r.Body)
    if err != nil {
        return errs.NewError(err)
    }
    
    signature := r.Header.Get("Stripe-Signature")
    
    if err := a.paymentBus.HandleWebhook(ctx, payload, signature); err != nil {
        return errs.NewError(err)
    }
    
    return web.StatusOK
}
```

---

### 5.5 Configure Stripe Webhook

#### Tasks
- [ ] Set up webhook endpoint in Stripe Dashboard
- [ ] Configure webhook URL: `https://your-domain.com/v1/webhooks/stripe`
- [ ] Select events to listen for:
  - `payment_intent.succeeded`
  - `payment_intent.payment_failed`
- [ ] Save webhook secret for signature verification

---

**Phase 5 Summary:**
- ✅ Stripe SDK integrated
- ✅ Payment service implemented
- ✅ Payment API endpoints created
- ✅ Webhook handling configured
- ✅ Ready for frontend integration

