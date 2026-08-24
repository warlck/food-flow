package categorybus

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
	ErrNotFound = errors.New("category not found")
)

// Storer interface declares the behavior this package needs to persist and
// retrieve data.
type Storer interface {
	Create(ctx context.Context, cat Category) error
	Update(ctx context.Context, cat Category) error
	Delete(ctx context.Context, cat Category) error
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
	now := time.Now()

	cat := Category{
		ID:           uuid.New(),
		Name:         nc.Name,
		Description:  nc.Description,
		RestaurantID: nc.RestaurantID,
		Enabled:      true,
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

	cat.DateUpdated = time.Now()

	if err := b.storer.Update(ctx, cat); err != nil {
		return Category{}, fmt.Errorf("update: %w", err)
	}

	return cat, nil
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
