package orderbus

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/warlck/food-flow/business/domain/menuitembus"
	"github.com/warlck/food-flow/business/domain/restaurantbus"
	"github.com/warlck/food-flow/business/sdk/order"
	"github.com/warlck/food-flow/business/sdk/page"
	"github.com/warlck/food-flow/foundation/logger"
)

// Business manages the set of APIs for order access.
type Business struct {
	log           *logger.Logger
	storer        Storer
	menuItemBus   *menuitembus.Business
	restaurantBus *restaurantbus.Business
}

// NewBusiness constructs a orderbus business API for use.
func NewBusiness(log *logger.Logger, storer Storer, menuItemBus *menuitembus.Business, restaurantBus *restaurantbus.Business) *Business {
	return &Business{
		log:           log,
		storer:        storer,
		menuItemBus:   menuItemBus,
		restaurantBus: restaurantBus,
	}
}

// =============================================================================

// Create creates a new order in the system.
func (b *Business) Create(ctx context.Context, no NewOrder) (Order, error) {
	// Parse restaurant ID
	restaurantID, err := uuid.Parse(no.RestaurantID)
	if err != nil {
		return Order{}, fmt.Errorf("invalid restaurant ID: %w", err)
	}

	// Validate restaurant exists
	if _, err := b.restaurantBus.QueryByID(ctx, restaurantID); err != nil {
		return Order{}, fmt.Errorf("restaurant validation: %w", err)
	}

	// Validate delivery address for delivery orders
	if no.OrderType == OrderTypeDelivery && no.DeliveryAddress == nil {
		return Order{}, fmt.Errorf("delivery address required for delivery orders")
	}

	// Create order items and calculate totals
	now := time.Now()
	items := make([]OrderItem, len(no.Items))
	var subtotal float64

	for i, newItem := range no.Items {
		// Parse menu item ID
		menuItemID, err := uuid.Parse(newItem.MenuItemID)
		if err != nil {
			return Order{}, fmt.Errorf("invalid menu item ID: %w", err)
		}

		// Get menu item details
		menuItem, err := b.menuItemBus.QueryByID(ctx, menuItemID)
		if err != nil {
			return Order{}, fmt.Errorf("menu item %s: %w", newItem.MenuItemID, err)
		}

		// Validate menu item belongs to the restaurant
		if menuItem.RestaurantID.String() != no.RestaurantID {
			return Order{}, fmt.Errorf("menu item %s does not belong to restaurant %s", newItem.MenuItemID, no.RestaurantID)
		}

		items[i] = OrderItem{
			ID:                  uuid.New(),
			MenuItemID:          menuItemID,
			MenuItemName:        menuItem.Name.String(),
			MenuItemPrice:       menuItem.Price.Value(),
			Quantity:            newItem.Quantity,
			SpecialInstructions: newItem.SpecialInstructions,
			DateCreated:         now,
		}

		subtotal += menuItem.Price.Value() * float64(newItem.Quantity)
	}

	// Calculate delivery fee
	var deliveryFee float64
	if no.OrderType == OrderTypeDelivery {
		deliveryFee = 5.00 // Fixed delivery fee for now
	}

	// Calculate tax (example: 8%)
	tax := subtotal * 0.08

	// Calculate total
	total := subtotal + deliveryFee + tax

	// Create delivery address if provided
	var deliveryAddress *DeliveryAddress
	if no.DeliveryAddress != nil {
		deliveryAddress = &DeliveryAddress{
			ID:                   uuid.New(),
			Street:               no.DeliveryAddress.Street,
			City:                 no.DeliveryAddress.City,
			State:                no.DeliveryAddress.State,
			PostalCode:           no.DeliveryAddress.PostalCode,
			DeliveryInstructions: no.DeliveryAddress.DeliveryInstructions,
			DateCreated:          now,
		}
	}

	// Create the order
	order := Order{
		ID:                  uuid.New(),
		RestaurantID:        restaurantID,
		CustomerName:        no.CustomerName,
		CustomerEmail:       no.CustomerEmail,
		CustomerPhone:       no.CustomerPhone,
		OrderType:           no.OrderType,
		OrderStatus:         OrderStatusPending,
		PaymentStatus:       PaymentStatusPending,
		PaymentMethod:       no.PaymentMethod,
		Subtotal:            subtotal,
		DeliveryFee:         deliveryFee,
		Tax:                 tax,
		Total:               total,
		SpecialInstructions: no.SpecialInstructions,
		Items:               items,
		DeliveryAddress:     deliveryAddress,
		DateCreated:         now,
		DateUpdated:         now,
	}

	// Store the order
	if err := b.storer.Create(ctx, order); err != nil {
		return Order{}, fmt.Errorf("create: %w", err)
	}

	return order, nil
}

// Query retrieves a list of orders with filtering and pagination.
func (b *Business) Query(ctx context.Context, filter QueryFilter, orderBy order.By, pageNumber page.Page) ([]Order, error) {
	orders, err := b.storer.Query(ctx, filter, orderBy, pageNumber)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return orders, nil
}

// QueryByID retrieves a single order by ID.
func (b *Business) QueryByID(ctx context.Context, orderID uuid.UUID) (Order, error) {
	order, err := b.storer.QueryByID(ctx, orderID)
	if err != nil {
		return Order{}, fmt.Errorf("query: orderID[%s]: %w", orderID, err)
	}

	return order, nil
}

// Count returns the total number of orders matching the filter.
func (b *Business) Count(ctx context.Context, filter QueryFilter) (int, error) {
	count, err := b.storer.Count(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("count: %w", err)
	}

	return count, nil
}

// =============================================================================

// UpdateStatus updates the order and/or payment status.
func (b *Business) UpdateStatus(ctx context.Context, orderID uuid.UUID, uo UpdateOrderStatus) error {
	order, err := b.QueryByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("query: %w", err)
	}

	if uo.OrderStatus != "" {
		order.OrderStatus = uo.OrderStatus
	}

	if uo.PaymentStatus != "" {
		order.PaymentStatus = uo.PaymentStatus
	}

	order.DateUpdated = time.Now()

	if err := b.storer.Update(ctx, order); err != nil {
		return fmt.Errorf("update: %w", err)
	}

	return nil
}

// UpdateStripePaymentIntent updates the Stripe PaymentIntent ID for an order.
func (b *Business) UpdateStripePaymentIntent(ctx context.Context, orderID uuid.UUID, paymentIntentID string) error {
	order, err := b.QueryByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("query: %w", err)
	}

	order.StripePaymentIntentID = paymentIntentID
	order.PaymentStatus = PaymentStatusProcessing
	order.DateUpdated = time.Now()

	if err := b.storer.Update(ctx, order); err != nil {
		return fmt.Errorf("update: %w", err)
	}

	return nil
}

// Cancel cancels an order.
func (b *Business) Cancel(ctx context.Context, orderID uuid.UUID) error {
	order, err := b.QueryByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("query: %w", err)
	}

	// Only allow cancellation of pending or confirmed orders
	if order.OrderStatus != OrderStatusPending && order.OrderStatus != OrderStatusConfirmed {
		return fmt.Errorf("cannot cancel order in %s status", order.OrderStatus)
	}

	order.OrderStatus = OrderStatusCancelled
	order.DateUpdated = time.Now()

	if err := b.storer.Update(ctx, order); err != nil {
		return fmt.Errorf("update: %w", err)
	}

	return nil
}
