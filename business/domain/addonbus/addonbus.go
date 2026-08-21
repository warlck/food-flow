package addonbus

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/warlck/food-flow/business/sdk/order"
	"github.com/warlck/food-flow/business/sdk/page"
	"github.com/warlck/food-flow/foundation/logger"
)

// Set of error variables for CRUD operations.
var (
	ErrNotFound = errors.New("addon not found")
)

// Storer interface declares the behavior this package needs to persist and
// retrieve data.
type Storer interface {
	Create(ctx context.Context, addon Addon) error
	Update(ctx context.Context, addon Addon) error
	Delete(ctx context.Context, addon Addon) error
	Query(ctx context.Context, filter QueryFilter, orderBy order.By, page page.Page) ([]Addon, error)
	Count(ctx context.Context, filter QueryFilter) (int, error)
	QueryByID(ctx context.Context, addonID uuid.UUID) (Addon, error)
	QueryByCategoryID(ctx context.Context, categoryID uuid.UUID) ([]Addon, error)
	Reorder(ctx context.Context, categoryID uuid.UUID, orderedIDs []uuid.UUID) error
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

// Create adds a new addon to the system.
func (b *Business) Create(ctx context.Context, na NewAddon) (Addon, error) {
	now := time.Now()

	maxQty := na.MaxQuantity
	if maxQty <= 0 {
		maxQty = 10 // default max quantity
	}

	addon := Addon{
		ID:           uuid.New(),
		CategoryID:   na.CategoryID,
		RestaurantID: na.RestaurantID,
		Name:         na.Name,
		Description:  na.Description,
		Price:        na.Price,
		Available:    true,
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
		addon.Name = *ua.Name
	}

	if ua.Description != nil {
		addon.Description = *ua.Description
	}

	if ua.Price != nil {
		addon.Price = *ua.Price
	}

	if ua.Available != nil {
		addon.Available = *ua.Available
	}

	if ua.MaxQuantity != nil {
		addon.MaxQuantity = *ua.MaxQuantity
	}

	if ua.Rank != nil {
		addon.Rank = ua.Rank
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

// QueryByCategoryID finds all addons for a specific category.
func (b *Business) QueryByCategoryID(ctx context.Context, categoryID uuid.UUID) ([]Addon, error) {
	addons, err := b.storer.QueryByCategoryID(ctx, categoryID)
	if err != nil {
		return nil, fmt.Errorf("query: categoryID[%s]: %w", categoryID, err)
	}

	return addons, nil
}

// Reorder updates the display rank of all addons in a category.
func (b *Business) Reorder(ctx context.Context, categoryID uuid.UUID, orderedIDs []uuid.UUID) error {
	addons, err := b.storer.QueryByCategoryID(ctx, categoryID)
	if err != nil {
		return fmt.Errorf("query category addons: %w", err)
	}

	if len(addons) != len(orderedIDs) {
		return errors.New("orderedIds must contain all addons in the category exactly once")
	}

	addonMap := make(map[uuid.UUID]bool, len(addons))
	for _, a := range addons {
		addonMap[a.ID] = true
	}

	for _, id := range orderedIDs {
		if !addonMap[id] {
			return errors.New("orderedIds contains invalid or duplicate addon id")
		}
		delete(addonMap, id)
	}

	if err := b.storer.Reorder(ctx, categoryID, orderedIDs); err != nil {
		return fmt.Errorf("reorder: %w", err)
	}

	return nil
}
