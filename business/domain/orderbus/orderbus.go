package orderbus

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/warlck/food-flow/business/domain/addonbus"
	"github.com/warlck/food-flow/business/domain/menuitembus"
	"github.com/warlck/food-flow/business/domain/promobus"
	"github.com/warlck/food-flow/business/domain/restaurantbus"
	"github.com/warlck/food-flow/business/sdk/order"
	"github.com/warlck/food-flow/business/sdk/page"
	"github.com/warlck/food-flow/business/types/money"
	"github.com/warlck/food-flow/foundation/logger"
)

// Set of error variables for addon validation during order creation.
var (
	// ErrAddonUnavailable is returned when an addon is not currently available.
	ErrAddonUnavailable = errors.New("addon is not available")

	// ErrAddonQuantityOutOfRange is returned when an addon quantity is below 1
	// or exceeds the addon's maximum quantity.
	ErrAddonQuantityOutOfRange = errors.New("addon quantity out of allowed range")

	// ErrAddonCategoryMismatch is returned when an addon does not belong to
	// the same category as the menu item it is applied to.
	ErrAddonCategoryMismatch = errors.New("addon does not belong to the menu item's category")

	// ErrMinSpendNotMet is returned when the order subtotal is below the restaurant's minimum spend requirement.
	ErrMinSpendNotMet = errors.New("subtotal does not meet restaurant minimum spend")
)

// Business manages the set of APIs for order access.
type Business struct {
	log           *logger.Logger
	storer        Storer
	menuItemBus   *menuitembus.Business
	restaurantBus *restaurantbus.Business
	addonBus      *addonbus.Business
	promoBus      *promobus.Business
}

// NewBusiness constructs a orderbus business API for use.
func NewBusiness(log *logger.Logger, storer Storer, menuItemBus *menuitembus.Business, restaurantBus *restaurantbus.Business, addonBus *addonbus.Business, promoBus *promobus.Business) *Business {
	return &Business{
		log:           log,
		storer:        storer,
		menuItemBus:   menuItemBus,
		restaurantBus: restaurantBus,
		addonBus:      addonBus,
		promoBus:      promoBus,
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
	restaurant, err := b.restaurantBus.QueryByID(ctx, restaurantID)
	if err != nil {
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

		itemID := uuid.New()

		// Process addons for this item, snapshotting name and price.
		var itemAddons []OrderItemAddon
		for _, newAddon := range newItem.Addons {
			addonID, err := uuid.Parse(newAddon.AddonID)
			if err != nil {
				return Order{}, fmt.Errorf("invalid addon ID: %w", err)
			}

			addon, err := b.addonBus.QueryByID(ctx, addonID)
			if err != nil {
				return Order{}, fmt.Errorf("addon %s: %w", newAddon.AddonID, err)
			}

			// Validate addon belongs to the restaurant
			if addon.RestaurantID.String() != no.RestaurantID {
				return Order{}, fmt.Errorf("addon %s does not belong to restaurant %s", newAddon.AddonID, no.RestaurantID)
			}

			// Validate addon belongs to the menu item's category
			if addon.CategoryID != menuItem.CategoryID {
				return Order{}, fmt.Errorf("%w: addon %s cannot be applied to menu item %s", ErrAddonCategoryMismatch, newAddon.AddonID, newItem.MenuItemID)
			}

			if !addon.Available {
				return Order{}, fmt.Errorf("%w: %s", ErrAddonUnavailable, addon.Name.String())
			}

			if newAddon.Quantity < 1 || newAddon.Quantity > addon.MaxQuantity {
				return Order{}, fmt.Errorf("%w: quantity %d for addon %s (max %d)", ErrAddonQuantityOutOfRange, newAddon.Quantity, newAddon.AddonID, addon.MaxQuantity)
			}

			itemAddons = append(itemAddons, OrderItemAddon{
				ID:          uuid.New(),
				OrderItemID: itemID,
				AddonID:     addonID,
				AddonName:   addon.Name.String(),
				AddonPrice:  addon.Price,
				Quantity:    newAddon.Quantity,
				DateCreated: now,
			})

			subtotal += addon.Price.Value() * float64(newAddon.Quantity) * float64(newItem.Quantity)
		}

		items[i] = OrderItem{
			ID:                  itemID,
			MenuItemID:          menuItemID,
			MenuItemName:        menuItem.Name.String(),
			MenuItemPrice:       menuItem.Price,
			Quantity:            newItem.Quantity,
			SpecialInstructions: newItem.SpecialInstructions,
			Addons:              itemAddons,
			DateCreated:         now,
		}

		subtotal += menuItem.Price.Value() * float64(newItem.Quantity)
	}

	// Round subtotal to 2 decimal places to avoid precision errors
	subtotal = roundToTwoDecimals(subtotal)

	// Validate restaurant minimum spend
	if restaurant.MinSpend > 0 && subtotal < restaurant.MinSpend {
		return Order{}, fmt.Errorf("%w: subtotal %.2f is less than minimum spend %.2f", ErrMinSpendNotMet, subtotal, restaurant.MinSpend)
	}

	// Validate promo code if provided
	var discountVal float64
	var promoCodeClean string
	var promoID *uuid.UUID

	if strings.TrimSpace(no.PromoCode) != "" && b.promoBus != nil {
		valRes, err := b.promoBus.ValidatePromoCode(ctx, no.PromoCode, &restaurantID, subtotal)
		if err != nil {
			return Order{}, fmt.Errorf("promo validation: %w", err)
		}
		if !valRes.Valid {
			return Order{}, fmt.Errorf("invalid promo code: %s", valRes.Reason)
		}
		discountVal = valRes.DiscountAmount
		promoCodeClean = valRes.Code
		if valRes.Promotion != nil {
			promoID = &valRes.Promotion.ID
		}
	}

	// Calculate taxable subtotal (subtotal - discount)
	taxableSubtotal := math.Max(0, subtotal-discountVal)

	// Calculate delivery fee based on the distance to the destination.
	var deliveryFee float64
	if no.OrderType == OrderTypeDelivery {
		if no.DeliveryAddress.Latitude == nil || no.DeliveryAddress.Longitude == nil {
			return Order{}, ErrDeliveryCoordinatesRequired
		}

		if restaurant.Latitude == nil || restaurant.Longitude == nil {
			return Order{}, ErrRestaurantLocationMissing
		}

		quote, err := deliveryQuote(*restaurant.Latitude, *restaurant.Longitude, *no.DeliveryAddress.Latitude, *no.DeliveryAddress.Longitude, restaurant.MaxDeliveryDistanceKm)
		if err != nil {
			return Order{}, err
		}

		if !quote.WithinLimit {
			return Order{}, fmt.Errorf("%w: %.2f km exceeds %.2f km limit", ErrDeliveryOutOfRange, quote.DistanceKm, quote.MaxDeliveryDistanceKm)
		}

		deliveryFee = quote.DeliveryFee.Value()
	}

	// Calculate tax using restaurant's tax rate on taxable subtotal and round to 2 decimal places
	tax := roundToTwoDecimals(taxableSubtotal * restaurant.TaxRate)

	// Calculate total and round to 2 decimal places
	total := roundToTwoDecimals(taxableSubtotal + deliveryFee + tax)

	subtotalFee, err := money.Parse(subtotal)
	if err != nil {
		return Order{}, err
	}
	discountFee, err := money.Parse(discountVal)
	if err != nil {
		return Order{}, err
	}
	deliveryFeeM, err := money.Parse(deliveryFee)
	if err != nil {
		return Order{}, err
	}
	taxFee, err := money.Parse(tax)
	if err != nil {
		return Order{}, err
	}
	totalFee, err := money.Parse(total)
	if err != nil {
		return Order{}, err
	}

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
			Latitude:             no.DeliveryAddress.Latitude,
			Longitude:            no.DeliveryAddress.Longitude,
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
		PromoCode:           promoCodeClean,
		Subtotal:            subtotalFee,
		Discount:            discountFee,
		DeliveryFee:         deliveryFeeM,
		Tax:                 taxFee,
		Total:               totalFee,
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

	// Increment promo code usage count if applicable
	if promoID != nil && b.promoBus != nil {
		if err := b.promoBus.IncrementUsage(ctx, *promoID); err != nil {
			if b.log != nil {
				b.log.Error(ctx, "failed to increment promo code usage", "promoID", *promoID, "err", err)
			}
		}
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

	if uo.OrderStatus == OrderStatusOutForDelivery && order.OrderType != OrderTypeDelivery {
		return ErrOutForDeliveryRequiresDelivery
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

// roundToTwoDecimals rounds a float64 to two decimal places.
func roundToTwoDecimals(value float64) float64 {
	return math.Round(value*100) / 100
}
