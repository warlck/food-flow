package restaurantbus

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
	ErrNotFound = errors.New("restaurant not found")
)

// Storer interface declares the behavior this package needs to persist and
// retrieve data.
type Storer interface {
	Create(ctx context.Context, res Restaurant) error
	Update(ctx context.Context, res Restaurant) error
	Delete(ctx context.Context, res Restaurant) error
	Query(ctx context.Context, filter QueryFilter, orderBy order.By, page page.Page) ([]Restaurant, error)
	Count(ctx context.Context, filter QueryFilter) (int, error)
	QueryByID(ctx context.Context, restaurantID uuid.UUID) (Restaurant, error)
}

// Business manages the set of APIs for restaurant access.
type Business struct {
	log    *logger.Logger
	storer Storer
}

// NewBusiness constructs a restaurant business API for use.
func NewBusiness(log *logger.Logger, storer Storer) *Business {
	return &Business{
		log:    log,
		storer: storer,
	}
}

// Create adds a new restaurant to the system.
func (b *Business) Create(ctx context.Context, nr NewRestaurant) (Restaurant, error) {
	now := time.Now()

	res := Restaurant{
		ID:          uuid.New(),
		Name:        nr.Name,
		Description: nr.Description,
		Address:     nr.Address,
		Phone:       nr.Phone,
		Email:       nr.Email,
		ImageURL:    nr.ImageURL,
		Enabled:     true,
		DateCreated: now,
		DateUpdated: now,
	}

	if err := b.storer.Create(ctx, res); err != nil {
		return Restaurant{}, fmt.Errorf("create: %w", err)
	}

	return res, nil
}

// Update modifies information about a restaurant.
func (b *Business) Update(ctx context.Context, res Restaurant, ur UpdateRestaurant) (Restaurant, error) {
	if ur.Name != nil {
		res.Name = *ur.Name
	}

	if ur.Description != nil {
		res.Description = *ur.Description
	}

	if ur.Address != nil {
		res.Address = *ur.Address
	}

	if ur.Phone != nil {
		res.Phone = *ur.Phone
	}

	if ur.Email != nil {
		res.Email = *ur.Email
	}

	if ur.ImageURL != nil {
		res.ImageURL = *ur.ImageURL
	}

	if ur.Enabled != nil {
		res.Enabled = *ur.Enabled
	}

	res.DateUpdated = time.Now()

	if err := b.storer.Update(ctx, res); err != nil {
		return Restaurant{}, fmt.Errorf("update: %w", err)
	}

	return res, nil
}

// Delete removes the specified restaurant.
func (b *Business) Delete(ctx context.Context, res Restaurant) error {
	if err := b.storer.Delete(ctx, res); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}

// Query retrieves a list of existing restaurants.
func (b *Business) Query(ctx context.Context, filter QueryFilter, orderBy order.By, page page.Page) ([]Restaurant, error) {
	restaurants, err := b.storer.Query(ctx, filter, orderBy, page)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return restaurants, nil
}

// Count returns the total number of restaurants.
func (b *Business) Count(ctx context.Context, filter QueryFilter) (int, error) {
	return b.storer.Count(ctx, filter)
}

// QueryByID finds the restaurant by the specified ID.
func (b *Business) QueryByID(ctx context.Context, restaurantID uuid.UUID) (Restaurant, error) {
	restaurant, err := b.storer.QueryByID(ctx, restaurantID)
	if err != nil {
		return Restaurant{}, fmt.Errorf("query: restaurantID[%s]: %w", restaurantID, err)
	}

	return restaurant, nil
}
