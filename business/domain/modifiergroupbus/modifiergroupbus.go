// Package modifiergroupbus provides business access to modifier group domain.
package modifiergroupbus

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
	ErrNotFound       = errors.New("modifier group not found")
	ErrInvalidReorder = errors.New("invalid reorder set")
)

// Storer interface declares the behavior this package needs to persist and retrieve data.
type Storer interface {
	Create(ctx context.Context, group ModifierGroup) error
	Update(ctx context.Context, group ModifierGroup) error
	Delete(ctx context.Context, group ModifierGroup) error
	Reorder(ctx context.Context, groups []ModifierGroup) error
	Query(ctx context.Context, filter QueryFilter, orderBy order.By, page page.Page) ([]ModifierGroup, error)
	QueryAll(ctx context.Context, filter QueryFilter, orderBy order.By) ([]ModifierGroup, error)
	Count(ctx context.Context, filter QueryFilter) (int, error)
	QueryByID(ctx context.Context, groupID uuid.UUID) (ModifierGroup, error)
}

// Business manages the set of APIs for modifier group access.
type Business struct {
	log    *logger.Logger
	storer Storer
}

// NewBusiness constructs a modifier group business API for use.
func NewBusiness(log *logger.Logger, storer Storer) *Business {
	return &Business{
		log:    log,
		storer: storer,
	}
}

// Create adds a new modifier group to the system.
func (b *Business) Create(ctx context.Context, ng NewModifierGroup) (ModifierGroup, error) {
	if ng.MinSelections != 0 && ng.MinSelections != 1 {
		return ModifierGroup{}, fmt.Errorf("min_selections must be 0 or 1")
	}
	if ng.MaxSelections != 1 {
		return ModifierGroup{}, fmt.Errorf("max_selections must be 1 in v1")
	}
	if ng.MinSelections > ng.MaxSelections {
		return ModifierGroup{}, fmt.Errorf("min_selections cannot exceed max_selections")
	}
	if ng.Rank != nil && *ng.Rank < 1 {
		return ModifierGroup{}, fmt.Errorf("rank must be >= 1")
	}

	now := time.Now()

	group := ModifierGroup{
		ID:            uuid.New(),
		MenuItemID:    ng.MenuItemID,
		RestaurantID:  ng.RestaurantID,
		Name:          ng.Name,
		Description:   ng.Description,
		MinSelections: ng.MinSelections,
		MaxSelections: ng.MaxSelections,
		Available:     ng.Available,
		Rank:          ng.Rank,
		DateCreated:   now,
		DateUpdated:   now,
	}

	if err := b.storer.Create(ctx, group); err != nil {
		return ModifierGroup{}, fmt.Errorf("create: %w", err)
	}

	return group, nil
}

// Update modifies information about a modifier group.
func (b *Business) Update(ctx context.Context, group ModifierGroup, ug UpdateModifierGroup) (ModifierGroup, error) {
	if ug.Name != nil {
		group.Name = *ug.Name
	}
	if ug.Description != nil {
		group.Description = *ug.Description
	}
	if ug.MinSelections != nil {
		if *ug.MinSelections != 0 && *ug.MinSelections != 1 {
			return ModifierGroup{}, fmt.Errorf("min_selections must be 0 or 1")
		}
		group.MinSelections = *ug.MinSelections
	}
	if ug.MaxSelections != nil {
		if *ug.MaxSelections != 1 {
			return ModifierGroup{}, fmt.Errorf("max_selections must be 1 in v1")
		}
		group.MaxSelections = *ug.MaxSelections
	}
	if group.MinSelections > group.MaxSelections {
		return ModifierGroup{}, fmt.Errorf("min_selections cannot exceed max_selections")
	}
	if ug.Available != nil {
		group.Available = *ug.Available
	}
	if ug.Rank.Present {
		if ug.Rank.Value != nil && *ug.Rank.Value < 1 {
			return ModifierGroup{}, fmt.Errorf("rank must be >= 1")
		}
		group.Rank = ug.Rank.Value
	}

	group.DateUpdated = time.Now()

	if err := b.storer.Update(ctx, group); err != nil {
		return ModifierGroup{}, fmt.Errorf("update: %w", err)
	}

	return group, nil
}

// Reorder updates the ranks of modifier groups within a menu item.
func (b *Business) Reorder(ctx context.Context, menuItemID uuid.UUID, orderedIDs []uuid.UUID) ([]ModifierGroup, error) {
	if len(orderedIDs) == 0 {
		return nil, fmt.Errorf("%w: ordered IDs cannot be empty", ErrInvalidReorder)
	}

	seen := make(map[uuid.UUID]bool, len(orderedIDs))
	for _, id := range orderedIDs {
		if seen[id] {
			return nil, fmt.Errorf("%w: duplicate modifier group id %s", ErrInvalidReorder, id)
		}
		seen[id] = true
	}

	existing, err := b.storer.QueryAll(ctx, QueryFilter{MenuItemID: &menuItemID}, order.NewBy(OrderByID, order.ASC))
	if err != nil {
		return nil, fmt.Errorf("query modifier groups for reorder: %w", err)
	}

	if len(existing) != len(orderedIDs) {
		return nil, fmt.Errorf("%w: exact set mismatch: expected %d modifier groups, got %d", ErrInvalidReorder, len(existing), len(orderedIDs))
	}

	existingMap := make(map[uuid.UUID]ModifierGroup, len(existing))
	for _, grp := range existing {
		existingMap[grp.ID] = grp
	}

	now := time.Now()
	reordered := make([]ModifierGroup, len(orderedIDs))
	for i, id := range orderedIDs {
		grp, exists := existingMap[id]
		if !exists {
			return nil, fmt.Errorf("%w: modifier group id %s does not belong to menu item %s", ErrInvalidReorder, id, menuItemID)
		}
		rankVal := (i + 1) * 10
		grp.Rank = &rankVal
		grp.DateUpdated = now
		reordered[i] = grp
	}

	if err := b.storer.Reorder(ctx, reordered); err != nil {
		return nil, fmt.Errorf("reorder modifier groups: %w", err)
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

// Delete removes the specified modifier group.
func (b *Business) Delete(ctx context.Context, group ModifierGroup) error {
	if err := b.storer.Delete(ctx, group); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}

// Query retrieves a list of existing modifier groups.
func (b *Business) Query(ctx context.Context, filter QueryFilter, orderBy order.By, page page.Page) ([]ModifierGroup, error) {
	groups, err := b.storer.Query(ctx, filter, orderBy, page)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return groups, nil
}

// QueryAll retrieves all modifier groups matching the filter without pagination.
func (b *Business) QueryAll(ctx context.Context, filter QueryFilter, orderBy order.By) ([]ModifierGroup, error) {
	groups, err := b.storer.QueryAll(ctx, filter, orderBy)
	if err != nil {
		return nil, fmt.Errorf("queryall: %w", err)
	}

	return groups, nil
}

// Count returns the total number of modifier groups.
func (b *Business) Count(ctx context.Context, filter QueryFilter) (int, error) {
	return b.storer.Count(ctx, filter)
}

// QueryByID finds the modifier group by the specified ID.
func (b *Business) QueryByID(ctx context.Context, groupID uuid.UUID) (ModifierGroup, error) {
	group, err := b.storer.QueryByID(ctx, groupID)
	if err != nil {
		return ModifierGroup{}, fmt.Errorf("query: groupID[%s]: %w", groupID, err)
	}

	return group, nil
}
