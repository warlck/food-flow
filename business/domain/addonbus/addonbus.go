package addonbus

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/warlck/food-flow/business/sdk/order"
	"github.com/warlck/food-flow/business/sdk/page"
	"github.com/warlck/food-flow/foundation/logger"
)

// Set of error variables for CRUD operations.
var (
	ErrNotFound       = errors.New("addon not found")
	ErrInvalidOrder   = errors.New("invalid addon order")
	ErrInvalidReorder = errors.New("invalid reorder set")
	ErrDuplicateName  = errors.New("addon name already exists for this menu item")
)

// Storer interface declares the behavior this package needs to persist and retrieve data.
type Storer interface {
	Create(ctx context.Context, addon Addon) error
	Update(ctx context.Context, addon Addon) error
	Delete(ctx context.Context, addon Addon) error
	Query(ctx context.Context, filter QueryFilter, orderBy order.By, page page.Page) ([]Addon, error)
	QueryAll(ctx context.Context, filter QueryFilter, orderBy order.By) ([]Addon, error)
	Count(ctx context.Context, filter QueryFilter) (int, error)
	QueryByID(ctx context.Context, addonID uuid.UUID) (Addon, error)
	Reorder(ctx context.Context, menuItemID uuid.UUID, addonIDs []uuid.UUID) error
}

// Business manages the set of APIs for addon access.
type Business struct {
	log    *logger.Logger
	storer Storer
}

// NewBusiness constructs an addon business API for use.
func NewBusiness(log *logger.Logger, storer Storer) *Business {
	return &Business{
		log:    log,
		storer: storer,
	}
}

// Create adds a new addon to the menu item.
func (b *Business) Create(ctx context.Context, na NewAddon) (Addon, error) {
	if na.Price.Value() < 0 {
		return Addon{}, fmt.Errorf("price must be >= 0")
	}

	maxQty := na.MaxQuantity
	if maxQty <= 0 {
		maxQty = 10
	}

	available := true
	if na.Available != nil {
		available = *na.Available
	}

	if na.Rank != nil && *na.Rank < 1 {
		return Addon{}, fmt.Errorf("rank must be >= 1")
	}

	// Check duplicate name within the menu item
	existing, err := b.storer.QueryAll(ctx, QueryFilter{
		MenuItemID: &na.MenuItemID,
		Name:       &na.Name,
	}, DefaultOrderBy)
	if err != nil {
		return Addon{}, fmt.Errorf("checking duplicate name: %w", err)
	}
	for _, a := range existing {
		if strings.EqualFold(a.Name.String(), na.Name.String()) {
			return Addon{}, ErrDuplicateName
		}
	}

	now := time.Now()

	addon := Addon{
		ID:           uuid.New(),
		MenuItemID:   na.MenuItemID,
		RestaurantID: na.RestaurantID,
		Name:         na.Name,
		Description:  na.Description,
		Price:        na.Price,
		Available:    available,
		MaxQuantity:  maxQty,
		Rank:         na.Rank,
		DateCreated:  now,
		DateUpdated:  now,
	}

	if err := b.storer.Create(ctx, addon); err != nil {
		return Addon{}, fmt.Errorf("create: %w", err)
	}

	return addon, nil
}

// Update modifies information about an addon.
func (b *Business) Update(ctx context.Context, addon Addon, ua UpdateAddon) (Addon, error) {
	if ua.Name != nil {
		if !strings.EqualFold(ua.Name.String(), addon.Name.String()) {
			existing, err := b.storer.QueryAll(ctx, QueryFilter{
				MenuItemID: &addon.MenuItemID,
				Name:       ua.Name,
			}, DefaultOrderBy)
			if err != nil {
				return Addon{}, fmt.Errorf("checking duplicate name: %w", err)
			}
			for _, a := range existing {
				if a.ID != addon.ID && strings.EqualFold(a.Name.String(), ua.Name.String()) {
					return Addon{}, ErrDuplicateName
				}
			}
		}
		addon.Name = *ua.Name
	}

	if ua.Description != nil {
		addon.Description = *ua.Description
	}

	if ua.Price != nil {
		if ua.Price.Value() < 0 {
			return Addon{}, fmt.Errorf("price must be >= 0")
		}
		addon.Price = *ua.Price
	}

	if ua.Available != nil {
		addon.Available = *ua.Available
	}

	if ua.MaxQuantity != nil {
		if *ua.MaxQuantity <= 0 {
			return Addon{}, fmt.Errorf("max_quantity must be >= 1")
		}
		addon.MaxQuantity = *ua.MaxQuantity
	}

	if ua.Rank.Present {
		if ua.Rank.Value != nil && *ua.Rank.Value < 1 {
			return Addon{}, fmt.Errorf("rank must be >= 1")
		}
		addon.Rank = ua.Rank.Value
	}

	addon.DateUpdated = time.Now()

	if err := b.storer.Update(ctx, addon); err != nil {
		return Addon{}, fmt.Errorf("update: %w", err)
	}

	return addon, nil
}

// Delete removes the specified addon.
func (b *Business) Delete(ctx context.Context, addon Addon) error {
	if err := b.storer.Delete(ctx, addon); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}

// Query retrieves a list of existing addons.
func (b *Business) Query(ctx context.Context, filter QueryFilter, orderBy order.By, page page.Page) ([]Addon, error) {
	addons, err := b.storer.Query(ctx, filter, orderBy, page)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return addons, nil
}

// QueryAll retrieves all addons matching the filter without pagination.
func (b *Business) QueryAll(ctx context.Context, filter QueryFilter, orderBy order.By) ([]Addon, error) {
	addons, err := b.storer.QueryAll(ctx, filter, orderBy)
	if err != nil {
		return nil, fmt.Errorf("queryall: %w", err)
	}

	return addons, nil
}

// Count returns the total number of addons.
func (b *Business) Count(ctx context.Context, filter QueryFilter) (int, error) {
	return b.storer.Count(ctx, filter)
}

// QueryByID finds the addon by the specified ID.
func (b *Business) QueryByID(ctx context.Context, addonID uuid.UUID) (Addon, error) {
	addon, err := b.storer.QueryByID(ctx, addonID)
	if err != nil {
		return Addon{}, fmt.Errorf("query: addonID[%s]: %w", addonID, err)
	}

	return addon, nil
}

// Reorder updates the ranks of addons on a menu item in a transaction.
func (b *Business) Reorder(ctx context.Context, menuItemID uuid.UUID, orderedAddonIDs []uuid.UUID) ([]Addon, error) {
	if len(orderedAddonIDs) == 0 {
		return nil, fmt.Errorf("%w: ordered IDs cannot be empty", ErrInvalidReorder)
	}

	seen := make(map[uuid.UUID]bool, len(orderedAddonIDs))
	for _, id := range orderedAddonIDs {
		if seen[id] {
			return nil, fmt.Errorf("%w: duplicate addon id %s", ErrInvalidReorder, id)
		}
		seen[id] = true
	}

	existing, err := b.storer.QueryAll(ctx, QueryFilter{
		MenuItemID: &menuItemID,
	}, DefaultOrderBy)
	if err != nil {
		return nil, fmt.Errorf("query addons for reorder: %w", err)
	}

	if len(existing) != len(orderedAddonIDs) {
		return nil, fmt.Errorf("%w: exact set mismatch: expected %d addons, got %d", ErrInvalidReorder, len(existing), len(orderedAddonIDs))
	}

	existingMap := make(map[uuid.UUID]Addon, len(existing))
	for _, a := range existing {
		existingMap[a.ID] = a
	}

	for _, id := range orderedAddonIDs {
		if _, exists := existingMap[id]; !exists {
			return nil, fmt.Errorf("%w: addon id %s does not belong to menu item %s", ErrInvalidReorder, id, menuItemID)
		}
	}

	if err := b.storer.Reorder(ctx, menuItemID, orderedAddonIDs); err != nil {
		return nil, fmt.Errorf("reorder addons: %w", err)
	}

	return b.storer.QueryAll(ctx, QueryFilter{
		MenuItemID: &menuItemID,
	}, DefaultOrderBy)
}
