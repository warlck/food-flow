package categorybus

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/warlck/food-flow/business/sdk/order"
	"github.com/warlck/food-flow/business/sdk/page"
	"github.com/warlck/food-flow/foundation/logger"
)

// Set of error variables for CRUD operations.
var (
	ErrNotFound       = errors.New("category not found")
	ErrInvalidReorder = errors.New("invalid reorder set")
)

// Storer interface declares the behavior this package needs to persist and
// retrieve data.
type Storer interface {
	Create(ctx context.Context, cat Category) error
	Update(ctx context.Context, cat Category) error
	Delete(ctx context.Context, cat Category) error
	Reorder(ctx context.Context, categories []Category) error
	Query(ctx context.Context, filter QueryFilter, orderBy order.By, page page.Page) ([]Category, error)
	QueryAll(ctx context.Context, filter QueryFilter, orderBy order.By) ([]Category, error)
	Count(ctx context.Context, filter QueryFilter) (int, error)
	QueryByID(ctx context.Context, categoryID uuid.UUID) (Category, error)
}

// Business manages the set of APIs for category access.
type Business struct {
	log    *logger.Logger
	storer Storer
}

// NewBusiness constructs a category business API for use.
func NewBusiness(log *logger.Logger, storer Storer) *Business {
	return &Business{
		log:    log,
		storer: storer,
	}
}

// Create adds a new category to the system.
func (b *Business) Create(ctx context.Context, nc NewCategory) (Category, error) {
	if nc.Rank != nil && *nc.Rank < 1 {
		return Category{}, fmt.Errorf("rank must be >= 1")
	}

	now := time.Now()

	cat := Category{
		ID:           uuid.New(),
		Name:         nc.Name,
		Description:  nc.Description,
		RestaurantID: nc.RestaurantID,
		Enabled:      true,
		Rank:         nc.Rank,
		DateCreated:  now,
		DateUpdated:  now,
	}

	if err := b.storer.Create(ctx, cat); err != nil {
		return Category{}, fmt.Errorf("create: %w", err)
	}

	return cat, nil
}

// Update modifies information about a category.
func (b *Business) Update(ctx context.Context, cat Category, uc UpdateCategory) (Category, error) {
	if uc.Name != nil {
		cat.Name = *uc.Name
	}

	if uc.Description != nil {
		cat.Description = *uc.Description
	}

	if uc.Enabled != nil {
		cat.Enabled = *uc.Enabled
	}

	if uc.Rank.Present {
		if uc.Rank.Value != nil && *uc.Rank.Value < 1 {
			return Category{}, fmt.Errorf("rank must be >= 1")
		}
		cat.Rank = uc.Rank.Value
	}

	cat.DateUpdated = time.Now()

	if err := b.storer.Update(ctx, cat); err != nil {
		return Category{}, fmt.Errorf("update: %w", err)
	}

	return cat, nil
}

// Reorder renumbers the rank of categories in the given restaurant.
func (b *Business) Reorder(ctx context.Context, restaurantID uuid.UUID, orderedIDs []uuid.UUID) ([]Category, error) {
	if len(orderedIDs) == 0 {
		return nil, fmt.Errorf("%w: ordered IDs cannot be empty", ErrInvalidReorder)
	}

	// Check duplicates
	seen := make(map[uuid.UUID]bool, len(orderedIDs))
	for _, id := range orderedIDs {
		if seen[id] {
			return nil, fmt.Errorf("%w: duplicate category id %s", ErrInvalidReorder, id)
		}
		seen[id] = true
	}

	existing, err := b.storer.QueryAll(ctx, QueryFilter{RestaurantID: &restaurantID}, order.NewBy(OrderByID, order.ASC))
	if err != nil {
		return nil, fmt.Errorf("query categories for reorder: %w", err)
	}

	if len(existing) != len(orderedIDs) {
		return nil, fmt.Errorf("%w: exact set mismatch: expected %d categories, got %d", ErrInvalidReorder, len(existing), len(orderedIDs))
	}

	existingMap := make(map[uuid.UUID]Category, len(existing))
	for _, cat := range existing {
		existingMap[cat.ID] = cat
	}

	now := time.Now()
	reordered := make([]Category, len(orderedIDs))
	for i, id := range orderedIDs {
		cat, exists := existingMap[id]
		if !exists {
			return nil, fmt.Errorf("%w: category id %s does not belong to restaurant %s", ErrInvalidReorder, id, restaurantID)
		}
		rankVal := (i + 1) * 10
		cat.Rank = &rankVal
		cat.DateUpdated = now
		reordered[i] = cat
	}

	if err := b.storer.Reorder(ctx, reordered); err != nil {
		return nil, fmt.Errorf("reorder categories: %w", err)
	}

	sort.SliceStable(reordered, func(i, j int) bool {
		r1 := reordered[i].Rank
		r2 := reordered[j].Rank
		if r1 != nil && r2 != nil && *r1 != *r2 {
			return *r1 < *r2
		}
		if r1 != nil && r2 == nil {
			return true
		}
		if r1 == nil && r2 != nil {
			return false
		}
		if reordered[i].Name.String() != reordered[j].Name.String() {
			return reordered[i].Name.String() < reordered[j].Name.String()
		}
		return reordered[i].ID.String() < reordered[j].ID.String()
	})

	return reordered, nil
}

// Delete removes the specified category.
func (b *Business) Delete(ctx context.Context, cat Category) error {
	if err := b.storer.Delete(ctx, cat); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}

// Query retrieves a list of existing categories.
func (b *Business) Query(ctx context.Context, filter QueryFilter, orderBy order.By, page page.Page) ([]Category, error) {
	categories, err := b.storer.Query(ctx, filter, orderBy, page)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return categories, nil
}

// QueryAll retrieves all categories matching the filter without pagination.
func (b *Business) QueryAll(ctx context.Context, filter QueryFilter, orderBy order.By) ([]Category, error) {
	categories, err := b.storer.QueryAll(ctx, filter, orderBy)
	if err != nil {
		return nil, fmt.Errorf("queryall: %w", err)
	}

	return categories, nil
}

// Count returns the total number of categories.
func (b *Business) Count(ctx context.Context, filter QueryFilter) (int, error) {
	return b.storer.Count(ctx, filter)
}

// QueryByID finds the category by the specified ID.
func (b *Business) QueryByID(ctx context.Context, categoryID uuid.UUID) (Category, error) {
	category, err := b.storer.QueryByID(ctx, categoryID)
	if err != nil {
		return Category{}, fmt.Errorf("query: categoryID[%s]: %w", categoryID, err)
	}

	return category, nil
}
