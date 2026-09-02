package menuitembus

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
	ErrNotFound       = errors.New("menu item not found")
	ErrInvalidOrder   = errors.New("invalid menu item order")
	ErrInvalidReorder = errors.New("invalid reorder set")
)

// Storer interface declares the behavior this package needs to persist and retrieve data.
type Storer interface {
	Create(ctx context.Context, item MenuItem) error
	Update(ctx context.Context, item MenuItem) error
	Delete(ctx context.Context, item MenuItem) error
	Reorder(ctx context.Context, items []MenuItem) error
	Query(ctx context.Context, filter QueryFilter, orderBy order.By, page page.Page) ([]MenuItem, error)
	QueryAll(ctx context.Context, filter QueryFilter, orderBy order.By) ([]MenuItem, error)
	Count(ctx context.Context, filter QueryFilter) (int, error)
	QueryByID(ctx context.Context, menuItemID uuid.UUID) (MenuItem, error)
	QueryByCategoryID(ctx context.Context, categoryID uuid.UUID) ([]MenuItem, error)
}

// Business manages the set of APIs for menu item access.
type Business struct {
	log    *logger.Logger
	storer Storer
}

// NewBusiness constructs a menu item business API for use.
func NewBusiness(log *logger.Logger, storer Storer) *Business {
	return &Business{
		log:    log,
		storer: storer,
	}
}

// Create adds a new menu item to the system.
func (b *Business) Create(ctx context.Context, ni NewMenuItem) (MenuItem, error) {
	if ni.Rank != nil && *ni.Rank < 1 {
		return MenuItem{}, fmt.Errorf("rank must be >= 1")
	}

	now := time.Now()

	item := MenuItem{
		ID:           uuid.New(),
		Name:         ni.Name,
		Description:  ni.Description,
		Price:        ni.Price,
		CategoryID:   ni.CategoryID,
		RestaurantID: ni.RestaurantID,
		ImageURL:     ni.ImageURL,
		Available:    true,
		Rank:         ni.Rank,
		DateCreated:  now,
		DateUpdated:  now,
	}

	if err := b.storer.Create(ctx, item); err != nil {
		return MenuItem{}, fmt.Errorf("create: %w", err)
	}

	return item, nil
}

// Update modifies information about a menu item.
func (b *Business) Update(ctx context.Context, item MenuItem, ui UpdateMenuItem) (MenuItem, error) {
	if ui.Name != nil {
		item.Name = *ui.Name
	}
	if ui.Description != nil {
		item.Description = *ui.Description
	}
	if ui.Price != nil {
		item.Price = *ui.Price
	}
	if ui.CategoryID != nil {
		item.CategoryID = *ui.CategoryID
	}
	if ui.ImageURL != nil {
		item.ImageURL = *ui.ImageURL
	}
	if ui.Available != nil {
		item.Available = *ui.Available
	}
	if ui.Rank.Present {
		if ui.Rank.Value != nil && *ui.Rank.Value < 1 {
			return MenuItem{}, fmt.Errorf("rank must be >= 1")
		}
		item.Rank = ui.Rank.Value
	}

	item.DateUpdated = time.Now()

	if err := b.storer.Update(ctx, item); err != nil {
		return MenuItem{}, fmt.Errorf("update: %w", err)
	}

	return item, nil
}

// Reorder updates the display rank of all menu items in a category.
func (b *Business) Reorder(ctx context.Context, categoryID uuid.UUID, orderedIDs []uuid.UUID) ([]MenuItem, error) {
	if len(orderedIDs) == 0 {
		return nil, fmt.Errorf("%w: ordered IDs cannot be empty", ErrInvalidReorder)
	}

	seen := make(map[uuid.UUID]bool, len(orderedIDs))
	for _, id := range orderedIDs {
		if seen[id] {
			return nil, fmt.Errorf("%w: duplicate menu item id %s", ErrInvalidReorder, id)
		}
		seen[id] = true
	}

	existing, err := b.storer.QueryByCategoryID(ctx, categoryID)
	if err != nil {
		return nil, fmt.Errorf("query category menu items: %w", err)
	}

	if len(existing) != len(orderedIDs) {
		return nil, fmt.Errorf("%w: exact set mismatch: expected %d menu items, got %d", ErrInvalidReorder, len(existing), len(orderedIDs))
	}

	existingMap := make(map[uuid.UUID]MenuItem, len(existing))
	for _, item := range existing {
		existingMap[item.ID] = item
	}

	now := time.Now()
	reordered := make([]MenuItem, len(orderedIDs))
	for i, id := range orderedIDs {
		item, exists := existingMap[id]
		if !exists {
			return nil, fmt.Errorf("%w: menu item id %s does not belong to category %s", ErrInvalidReorder, id, categoryID)
		}
		rankVal := (i + 1) * 10
		item.Rank = &rankVal
		item.DateUpdated = now
		reordered[i] = item
	}

	if err := b.storer.Reorder(ctx, reordered); err != nil {
		return nil, fmt.Errorf("reorder menu items: %w", err)
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

// Delete removes the specified menu item.
func (b *Business) Delete(ctx context.Context, item MenuItem) error {
	if err := b.storer.Delete(ctx, item); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}

// Query retrieves a list of existing menu items.
func (b *Business) Query(ctx context.Context, filter QueryFilter, orderBy order.By, page page.Page) ([]MenuItem, error) {
	items, err := b.storer.Query(ctx, filter, orderBy, page)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return items, nil
}

// QueryAll retrieves all menu items matching the filter without pagination.
func (b *Business) QueryAll(ctx context.Context, filter QueryFilter, orderBy order.By) ([]MenuItem, error) {
	items, err := b.storer.QueryAll(ctx, filter, orderBy)
	if err != nil {
		return nil, fmt.Errorf("queryall: %w", err)
	}

	return items, nil
}

// Count returns the total number of menu items.
func (b *Business) Count(ctx context.Context, filter QueryFilter) (int, error) {
	return b.storer.Count(ctx, filter)
}

// QueryByID finds the menu item by the specified ID.
func (b *Business) QueryByID(ctx context.Context, menuItemID uuid.UUID) (MenuItem, error) {
	item, err := b.storer.QueryByID(ctx, menuItemID)
	if err != nil {
		return MenuItem{}, fmt.Errorf("query: menuItemID[%s]: %w", menuItemID, err)
	}

	return item, nil
}

// QueryByCategoryID finds all menu items for a specific category.
func (b *Business) QueryByCategoryID(ctx context.Context, categoryID uuid.UUID) ([]MenuItem, error) {
	items, err := b.storer.QueryByCategoryID(ctx, categoryID)
	if err != nil {
		return nil, fmt.Errorf("query: categoryID[%s]: %w", categoryID, err)
	}

	return items, nil
}
