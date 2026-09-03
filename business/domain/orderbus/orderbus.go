package orderbus

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/warlck/food-flow/business/domain/addonbus"
	"github.com/warlck/food-flow/business/domain/categorybus"
	"github.com/warlck/food-flow/business/domain/menuitembus"
	"github.com/warlck/food-flow/business/domain/modifiergroupbus"
	"github.com/warlck/food-flow/business/domain/modifieroptionbus"
	"github.com/warlck/food-flow/business/domain/promobus"
	"github.com/warlck/food-flow/business/domain/restaurantbus"
	"github.com/warlck/food-flow/business/sdk/order"
	"github.com/warlck/food-flow/business/sdk/page"
	"github.com/warlck/food-flow/business/types/money"
	"github.com/warlck/food-flow/foundation/logger"
)

// Set of error variables for addon/order validation.
var (
	// ErrMinSpendNotMet is returned when the order subtotal is below the restaurant's minimum spend requirement.
	ErrMinSpendNotMet = errorsNew("subtotal does not meet restaurant minimum spend")
)

func errorsNew(text string) error {
	return fmt.Errorf("%s", text)
}

// Business manages the set of APIs for order access.
type Business struct {
	log               *logger.Logger
	storer            Storer
	menuItemBus       *menuitembus.Business
	restaurantBus     *restaurantbus.Business
	addonBus          *addonbus.Business
	promoBus          *promobus.Business
	categoryBus       *categorybus.Business
	modifierGroupBus  *modifiergroupbus.Business
	modifierOptionBus *modifieroptionbus.Business
}

// NewBusiness constructs an orderbus business API for use.
func NewBusiness(
	log *logger.Logger,
	storer Storer,
	menuItemBus *menuitembus.Business,
	restaurantBus *restaurantbus.Business,
	addonBus *addonbus.Business,
	promoBus *promobus.Business,
	categoryBus *categorybus.Business,
	modifierGroupBus *modifiergroupbus.Business,
	modifierOptionBus *modifieroptionbus.Business,
) *Business {
	return &Business{
		log:               log,
		storer:            storer,
		menuItemBus:       menuItemBus,
		restaurantBus:     restaurantBus,
		addonBus:          addonBus,
		promoBus:          promoBus,
		categoryBus:       categoryBus,
		modifierGroupBus:  modifierGroupBus,
		modifierOptionBus: modifierOptionBus,
	}
}

// =============================================================================

// Create creates a new order in the system.
func (b *Business) Create(ctx context.Context, no NewOrder) (Order, error) {
	if len(no.Items) == 0 {
		return Order{}, ErrEmptyItems
	}

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
		return Order{}, ErrInvalidDeliveryAddress
	}

	now := time.Now()
	orderID := uuid.New()
	items := make([]OrderItem, len(no.Items))
	var subtotal float64

	for i, newItem := range no.Items {
		if newItem.Quantity < 1 {
			return Order{}, fmt.Errorf("item quantity must be >= 1")
		}

		menuItemID, err := uuid.Parse(newItem.MenuItemID)
		if err != nil {
			return Order{}, fmt.Errorf("invalid menu item ID: %w", err)
		}

		menuItem, err := b.menuItemBus.QueryByID(ctx, menuItemID)
		if err != nil {
			return Order{}, fmt.Errorf("menu item %s: %w", newItem.MenuItemID, err)
		}

		if menuItem.RestaurantID != restaurantID {
			return Order{}, fmt.Errorf("menu item %s does not belong to restaurant %s", newItem.MenuItemID, no.RestaurantID)
		}

		if !menuItem.Available {
			return Order{}, fmt.Errorf("%w: %s", ErrMenuItemUnavailable, menuItem.Name.String())
		}

		// Snapshot Category
		category, err := b.categoryBus.QueryByID(ctx, menuItem.CategoryID)
		if err != nil {
			return Order{}, fmt.Errorf("category lookup %s: %w", menuItem.CategoryID, err)
		}

		itemID := uuid.New()
		unitPrice := menuItem.Price.Value()

		// ---------------------------------------------------------------------
		// Process Modifiers
		// ---------------------------------------------------------------------
		itemGroups, err := b.modifierGroupBus.QueryAll(ctx, modifiergroupbus.QueryFilter{MenuItemID: &menuItem.ID}, modifiergroupbus.DefaultOrderBy)
		if err != nil {
			return Order{}, fmt.Errorf("query modifier groups for item %s: %w", menuItem.ID, err)
		}

		groupsMap := make(map[uuid.UUID]modifiergroupbus.ModifierGroup, len(itemGroups))
		for _, g := range itemGroups {
			groupsMap[g.ID] = g
		}

		// Check for duplicate submitted modifier groups
		submittedGroupIDs := make(map[uuid.UUID]bool, len(newItem.Modifiers))
		submittedByGroup := make(map[uuid.UUID][]NewOrderItemModifier, len(newItem.Modifiers))
		for _, m := range newItem.Modifiers {
			gID, err := uuid.Parse(m.ModifierGroupID)
			if err != nil {
				return Order{}, fmt.Errorf("invalid modifier group ID %s: %w", m.ModifierGroupID, err)
			}
			if submittedGroupIDs[gID] {
				return Order{}, fmt.Errorf("%w: duplicate modifier group %s", ErrModifierSelectionLimit, m.ModifierGroupID)
			}
			submittedGroupIDs[gID] = true
			submittedByGroup[gID] = append(submittedByGroup[gID], m)
		}

		// Validate modifier selections per group. An unavailable group is
		// suspended: it imposes no requirement and accepts no selection.
		var itemModifiers []OrderItemModifier
		for _, g := range itemGroups {
			submitted := submittedByGroup[g.ID]

			if !g.Available {
				if len(submitted) > 0 {
					return Order{}, fmt.Errorf("%w: group %s (%s) is unavailable", ErrModifierGroupUnavailable, g.ID, g.Name.String())
				}
				continue
			}

			count := len(submitted)
			if count < g.MinSelections {
				return Order{}, fmt.Errorf("%w: group %s (%s)", ErrModifierRequired, g.ID, g.Name.String())
			}
			if count > g.MaxSelections {
				return Order{}, fmt.Errorf("%w: group %s", ErrModifierSelectionLimit, g.ID)
			}

			for _, m := range submitted {
				optID, err := uuid.Parse(m.ModifierOptionID)
				if err != nil {
					return Order{}, fmt.Errorf("invalid modifier option ID %s: %w", m.ModifierOptionID, err)
				}

				option, err := b.modifierOptionBus.QueryByID(ctx, optID)
				if err != nil {
					return Order{}, fmt.Errorf("%w: option %s: %v", ErrModifierOptionNotFound, m.ModifierOptionID, err)
				}

				if option.ModifierGroupID != g.ID || option.RestaurantID != restaurantID {
					return Order{}, fmt.Errorf("%w: option %s does not belong to group %s", ErrModifierOptionNotFound, option.ID, g.ID)
				}

				if !option.Available {
					return Order{}, fmt.Errorf("%w: option %s (%s)", ErrModifierOptionUnavailable, option.ID, option.Name.String())
				}

				itemModifiers = append(itemModifiers, OrderItemModifier{
					ID:                 uuid.New(),
					OrderItemID:        itemID,
					ModifierGroupID:    g.ID,
					ModifierGroupName:  g.Name.String(),
					ModifierOptionID:   option.ID,
					ModifierOptionName: option.Name.String(),
					PriceDelta:         option.PriceDelta,
					DateCreated:        now,
				})

				unitPrice += option.PriceDelta.Value()
			}
		}

		// Check for any submitted modifier group that does not belong to this item
		for gID := range submittedGroupIDs {
			if _, exists := groupsMap[gID]; !exists {
				return Order{}, fmt.Errorf("%w: group %s does not belong to menu item %s", ErrModifierGroupNotFound, gID, menuItem.ID)
			}
		}

		// ---------------------------------------------------------------------
		// Process Addons
		// ---------------------------------------------------------------------
		menuItemAddons, err := b.addonBus.QueryAll(ctx, addonbus.QueryFilter{
			MenuItemID: &menuItem.ID,
		}, addonbus.DefaultOrderBy)
		if err != nil {
			return Order{}, fmt.Errorf("lookup menu item addons: %w", err)
		}

		assignedMap := make(map[uuid.UUID]addonbus.Addon, len(menuItemAddons))
		for _, a := range menuItemAddons {
			assignedMap[a.ID] = a
		}

		var itemAddons []OrderItemAddon
		var addonTotalPerItemUnit float64
		seenAddonIDs := make(map[uuid.UUID]bool, len(newItem.Addons))

		for _, newAddon := range newItem.Addons {
			addonID, err := uuid.Parse(newAddon.AddonID)
			if err != nil {
				return Order{}, fmt.Errorf("invalid addon ID: %w", err)
			}

			if seenAddonIDs[addonID] {
				return Order{}, fmt.Errorf("duplicate addon ID %s", newAddon.AddonID)
			}
			seenAddonIDs[addonID] = true

			addon, exists := assignedMap[addonID]
			if !exists {
				return Order{}, fmt.Errorf("%w: addon %s cannot be applied to menu item %s", ErrAddonNotAssigned, newAddon.AddonID, newItem.MenuItemID)
			}

			if addon.RestaurantID != restaurantID {
				return Order{}, fmt.Errorf("addon %s does not belong to restaurant %s", newAddon.AddonID, no.RestaurantID)
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

			addonTotalPerItemUnit += addon.Price.Value() * float64(newAddon.Quantity)
		}

		lineTotal := float64(newItem.Quantity) * (unitPrice + addonTotalPerItemUnit)
		subtotal += lineTotal

		items[i] = OrderItem{
			ID:                  itemID,
			OrderID:             orderID,
			CategoryID:          category.ID,
			CategoryName:        category.Name.String(),
			MenuItemID:          menuItemID,
			MenuItemName:        menuItem.Name.String(),
			MenuItemPrice:       menuItem.Price,
			Quantity:            newItem.Quantity,
			SpecialInstructions: newItem.SpecialInstructions,
			Modifiers:           itemModifiers,
			Addons:              itemAddons,
			DateCreated:         now,
		}
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
	deliveryFeeMoney, err := money.Parse(deliveryFee)
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

	// Create delivery address if delivery order
	var deliveryAddr *DeliveryAddress
	if no.OrderType == OrderTypeDelivery && no.DeliveryAddress != nil {
		deliveryAddr = &DeliveryAddress{
			ID:                   uuid.New(),
			OrderID:              orderID,
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

	// Build order entity
	order := Order{
		ID:                  orderID,
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
		DeliveryFee:         deliveryFeeMoney,
		Tax:                 taxFee,
		Total:               totalFee,
		SpecialInstructions: no.SpecialInstructions,
		Items:               items,
		DeliveryAddress:     deliveryAddr,
		DateCreated:         now,
		DateUpdated:         now,
	}

	// Persist order in database
	if err := b.storer.Create(ctx, order); err != nil {
		return Order{}, fmt.Errorf("create: %w", err)
	}

	// Increment promo code usage if applied
	if promoID != nil && b.promoBus != nil {
		if err := b.promoBus.IncrementUsage(ctx, *promoID); err != nil {
			b.log.Error(ctx, "failed to increment promo code usage", "promo_id", *promoID, "error", err)
		}
	}

	return order, nil
}

// Update modifies information about an order.
func (b *Business) Update(ctx context.Context, order Order, uos UpdateOrderStatus) (Order, error) {
	if uos.OrderStatus == OrderStatusOutForDelivery && order.OrderType != OrderTypeDelivery {
		return Order{}, ErrOutForDeliveryRequiresDelivery
	}

	if uos.OrderStatus != "" {
		order.OrderStatus = uos.OrderStatus
	}

	if uos.PaymentStatus != "" {
		order.PaymentStatus = uos.PaymentStatus
		if uos.PaymentStatus == PaymentStatusPaid && order.OrderStatus == OrderStatusPending {
			order.OrderStatus = OrderStatusConfirmed
		}
	}

	order.DateUpdated = time.Now()

	if err := b.storer.Update(ctx, order); err != nil {
		return Order{}, fmt.Errorf("update: %w", err)
	}

	return order, nil
}

// UpdateStatus updates an order's status and/or payment status.
func (b *Business) UpdateStatus(ctx context.Context, orderID uuid.UUID, uos UpdateOrderStatus) error {
	order, err := b.storer.QueryByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("query order: %w", err)
	}

	_, err = b.Update(ctx, order, uos)
	return err
}

// UpdateStripePaymentIntent updates the Stripe PaymentIntent ID for an order.
func (b *Business) UpdateStripePaymentIntent(ctx context.Context, orderID uuid.UUID, paymentIntentID string) error {
	order, err := b.storer.QueryByID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("query order: %w", err)
	}

	order.StripePaymentIntentID = paymentIntentID
	order.DateUpdated = time.Now()

	if err := b.storer.Update(ctx, order); err != nil {
		return fmt.Errorf("update payment intent: %w", err)
	}

	return nil
}

// UpdatePaymentIntent updates the Stripe PaymentIntent ID for an order.
func (b *Business) UpdatePaymentIntent(ctx context.Context, orderID uuid.UUID, paymentIntentID string) error {
	return b.UpdateStripePaymentIntent(ctx, orderID, paymentIntentID)
}

// UpdatePaymentStatus updates the payment status for an order.
func (b *Business) UpdatePaymentStatus(ctx context.Context, orderID uuid.UUID, paymentStatus string) error {
	return b.UpdateStatus(ctx, orderID, UpdateOrderStatus{PaymentStatus: paymentStatus})
}

// Cancel marks the specified order as cancelled.
func (b *Business) Cancel(ctx context.Context, orderID uuid.UUID) error {
	return b.Delete(ctx, orderID)
}

// QueryOrderMetrics computes and aggregates all pure order-domain analytics datasets.
func (b *Business) QueryOrderMetrics(ctx context.Context, filter InsightsFilter) (OrderMetrics, error) {
	summary, err := b.storer.QuerySalesSummary(ctx, filter)
	if err != nil {
		return OrderMetrics{}, fmt.Errorf("query sales summary: %w", err)
	}

	salesOverTime, err := b.storer.QuerySalesOverTime(ctx, filter)
	if err != nil {
		return OrderMetrics{}, fmt.Errorf("query sales over time: %w", err)
	}

	topItems, err := b.storer.QueryTopItemSales(ctx, filter, 5)
	if err != nil {
		return OrderMetrics{}, fmt.Errorf("query top items: %w", err)
	}

	allItemSales, err := b.storer.QueryAllItemSales(ctx, filter)
	if err != nil {
		return OrderMetrics{}, fmt.Errorf("query all item sales: %w", err)
	}

	topAddons, err := b.storer.QueryTopAddonSales(ctx, filter, 5)
	if err != nil {
		return OrderMetrics{}, fmt.Errorf("query top addons: %w", err)
	}

	orderTypes, err := b.storer.QueryOrderTypes(ctx, filter)
	if err != nil {
		return OrderMetrics{}, fmt.Errorf("query order types: %w", err)
	}

	peakHours, err := b.storer.QueryPeakHours(ctx, filter)
	if err != nil {
		return OrderMetrics{}, fmt.Errorf("query peak hours: %w", err)
	}

	return OrderMetrics{
		Summary:       summary,
		SalesOverTime: salesOverTime,
		TopItems:      topItems,
		AllItemSales:  allItemSales,
		TopAddons:     topAddons,
		OrderTypes:    orderTypes,
		PeakHours:     peakHours,
	}, nil
}

// Delete marks the specified order as cancelled.
func (b *Business) Delete(ctx context.Context, orderID uuid.UUID) error {
	if err := b.storer.Delete(ctx, orderID); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}

// Query retrieves a list of existing orders.
func (b *Business) Query(ctx context.Context, filter QueryFilter, orderBy order.By, pageNumber page.Page) ([]Order, error) {
	orders, err := b.storer.Query(ctx, filter, orderBy, pageNumber)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return orders, nil
}

// Count returns the total number of orders.
func (b *Business) Count(ctx context.Context, filter QueryFilter) (int, error) {
	return b.storer.Count(ctx, filter)
}

// QueryByID finds the order by the specified ID.
func (b *Business) QueryByID(ctx context.Context, orderID uuid.UUID) (Order, error) {
	order, err := b.storer.QueryByID(ctx, orderID)
	if err != nil {
		return Order{}, fmt.Errorf("query: orderID[%s]: %w", orderID, err)
	}

	return order, nil
}

// Helper function to round float to 2 decimal places.
func roundToTwoDecimals(val float64) float64 {
	return math.Round(val*100) / 100
}

// Helper function to validate status transitions.
func isValidStatusTransition(from, to string) bool {
	validTransitions := map[string][]string{
		OrderStatusPending:        {OrderStatusConfirmed, OrderStatusCancelled},
		OrderStatusConfirmed:      {OrderStatusPreparing, OrderStatusCancelled},
		OrderStatusPreparing:      {OrderStatusReady, OrderStatusCancelled},
		OrderStatusReady:          {OrderStatusOutForDelivery, OrderStatusCompleted, OrderStatusCancelled},
		OrderStatusOutForDelivery: {OrderStatusCompleted, OrderStatusCancelled},
		OrderStatusCompleted:      {}, // Terminal state
		OrderStatusCancelled:      {}, // Terminal state
	}

	allowed, exists := validTransitions[from]
	if !exists {
		return false
	}

	for _, s := range allowed {
		if s == to {
			return true
		}
	}

	return false
}
