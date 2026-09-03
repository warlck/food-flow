// Package modifieroptionbus provides business access to modifier option domain.
package modifieroptionbus

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/warlck/food-flow/business/domain/modifiergroupbus"
	"github.com/warlck/food-flow/business/sdk/order"
	"github.com/warlck/food-flow/business/sdk/page"
	"github.com/warlck/food-flow/foundation/logger"
)

// Set of error variables for CRUD operations.
var (
	ErrNotFound       = errors.New("modifier option not found")
	ErrInvalidReorder = errors.New("invalid reorder set")

	// ErrLastAvailableOption is returned when disabling or deleting the last
	// available option of an active required group; disable the group or the
	// menu item first.
	ErrLastAvailableOption = errors.New("cannot disable or delete the last available option of an active required group")
)

// Storer interface declares the behavior this package needs to persist and retrieve data.
type Storer interface {
	Create(ctx context.Context, opt ModifierOption) error
	Update(ctx context.Context, opt ModifierOption) error
	Delete(ctx context.Context, opt ModifierOption) error
	Reorder(ctx context.Context, options []ModifierOption) error
	Query(ctx context.Context, filter QueryFilter, orderBy order.By, page page.Page) ([]ModifierOption, error)
	QueryAll(ctx context.Context, filter QueryFilter, orderBy order.By) ([]ModifierOption, error)
	Count(ctx context.Context, filter QueryFilter) (int, error)
	CountAvailable(ctx context.Context, groupID uuid.UUID) (int, error)
	QueryByID(ctx context.Context, optionID uuid.UUID) (ModifierOption, error)
}

// GroupStorer declares the behavior this package needs from the modifier
// group store to enforce the required-group availability invariant.
type GroupStorer interface {
	QueryByID(ctx context.Context, groupID uuid.UUID) (modifiergroupbus.ModifierGroup, error)
}

// Business manages the set of APIs for modifier option access.
type Business struct {
	log         *logger.Logger
	storer      Storer
	groupStorer GroupStorer
}

// NewBusiness constructs a modifier option business API for use.
func NewBusiness(log *logger.Logger, storer Storer, groupStorer GroupStorer) *Business {
	return &Business{
		log:         log,
		storer:      storer,
		groupStorer: groupStorer,
	}
}

// Create adds a new modifier option to the system.
func (b *Business) Create(ctx context.Context, no NewModifierOption) (ModifierOption, error) {
	if no.PriceDelta.Value() < 0 {
		return ModifierOption{}, fmt.Errorf("price_delta must be >= 0")
	}
	if no.Rank != nil && *no.Rank < 1 {
		return ModifierOption{}, fmt.Errorf("rank must be >= 1")
	}

	available := true
	if no.Available != nil {
		available = *no.Available
	}

	now := time.Now()

	opt := ModifierOption{
		ID:              uuid.New(),
		ModifierGroupID: no.ModifierGroupID,
		RestaurantID:    no.RestaurantID,
		Name:            no.Name,
		Description:     no.Description,
		PriceDelta:      no.PriceDelta,
		Available:       available,
		Rank:            no.Rank,
		DateCreated:     now,
		DateUpdated:     now,
	}

	if err := b.storer.Create(ctx, opt); err != nil {
		return ModifierOption{}, fmt.Errorf("create: %w", err)
	}

	return opt, nil
}

// Update modifies information about a modifier option.
func (b *Business) Update(ctx context.Context, opt ModifierOption, uo UpdateModifierOption) (ModifierOption, error) {
	wasAvailable := opt.Available

	if uo.Name != nil {
		opt.Name = *uo.Name
	}
	if uo.Description != nil {
		opt.Description = *uo.Description
	}
	if uo.PriceDelta != nil {
		if uo.PriceDelta.Value() < 0 {
			return ModifierOption{}, fmt.Errorf("price_delta must be >= 0")
		}
		opt.PriceDelta = *uo.PriceDelta
	}
	if uo.Available != nil {
		opt.Available = *uo.Available
	}
	if uo.Rank.Present {
		if uo.Rank.Value != nil && *uo.Rank.Value < 1 {
			return ModifierOption{}, fmt.Errorf("rank must be >= 1")
		}
		opt.Rank = uo.Rank.Value
	}

	// Disabling the last available option of an active required group would
	// make the item unorderable; the group or item must be disabled first.
	if wasAvailable && !opt.Available {
		if err := b.guardLastAvailableOption(ctx, opt); err != nil {
			return ModifierOption{}, err
		}
	}

	opt.DateUpdated = time.Now()

	if err := b.storer.Update(ctx, opt); err != nil {
		return ModifierOption{}, fmt.Errorf("update: %w", err)
	}

	return opt, nil
}

// Reorder updates the ranks of modifier options within a group.
func (b *Business) Reorder(ctx context.Context, modifierGroupID uuid.UUID, orderedIDs []uuid.UUID) ([]ModifierOption, error) {
	if len(orderedIDs) == 0 {
		return nil, fmt.Errorf("%w: ordered IDs cannot be empty", ErrInvalidReorder)
	}

	seen := make(map[uuid.UUID]bool, len(orderedIDs))
	for _, id := range orderedIDs {
		if seen[id] {
			return nil, fmt.Errorf("%w: duplicate modifier option id %s", ErrInvalidReorder, id)
		}
		seen[id] = true
	}

	existing, err := b.storer.QueryAll(ctx, QueryFilter{ModifierGroupID: &modifierGroupID}, order.NewBy(OrderByID, order.ASC))
	if err != nil {
		return nil, fmt.Errorf("query modifier options for reorder: %w", err)
	}

	if len(existing) != len(orderedIDs) {
		return nil, fmt.Errorf("%w: exact set mismatch: expected %d modifier options, got %d", ErrInvalidReorder, len(existing), len(orderedIDs))
	}

	existingMap := make(map[uuid.UUID]ModifierOption, len(existing))
	for _, option := range existing {
		existingMap[option.ID] = option
	}

	now := time.Now()
	reordered := make([]ModifierOption, len(orderedIDs))
	for i, id := range orderedIDs {
		option, exists := existingMap[id]
		if !exists {
			return nil, fmt.Errorf("%w: modifier option id %s does not belong to group %s", ErrInvalidReorder, id, modifierGroupID)
		}
		rankVal := (i + 1) * 10
		option.Rank = &rankVal
		option.DateUpdated = now
		reordered[i] = option
	}

	if err := b.storer.Reorder(ctx, reordered); err != nil {
		return nil, fmt.Errorf("reorder modifier options: %w", err)
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

// Delete removes the specified modifier option.
func (b *Business) Delete(ctx context.Context, opt ModifierOption) error {
	if opt.Available {
		if err := b.guardLastAvailableOption(ctx, opt); err != nil {
			return err
		}
	}

	if err := b.storer.Delete(ctx, opt); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}

// guardLastAvailableOption rejects the mutation when the option is the last
// available option of an active required group. The option must still be
// persisted as available when this runs, so a count of one means last.
func (b *Business) guardLastAvailableOption(ctx context.Context, opt ModifierOption) error {
	group, err := b.groupStorer.QueryByID(ctx, opt.ModifierGroupID)
	if err != nil {
		return fmt.Errorf("query modifier group: %w", err)
	}

	if !group.Available || group.MinSelections == 0 {
		return nil
	}

	count, err := b.storer.CountAvailable(ctx, opt.ModifierGroupID)
	if err != nil {
		return fmt.Errorf("count available options: %w", err)
	}

	if count <= 1 {
		return ErrLastAvailableOption
	}

	return nil
}

// Query retrieves a list of existing modifier options.
func (b *Business) Query(ctx context.Context, filter QueryFilter, orderBy order.By, page page.Page) ([]ModifierOption, error) {
	options, err := b.storer.Query(ctx, filter, orderBy, page)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return options, nil
}

// QueryAll retrieves all modifier options matching the filter without pagination.
func (b *Business) QueryAll(ctx context.Context, filter QueryFilter, orderBy order.By) ([]ModifierOption, error) {
	options, err := b.storer.QueryAll(ctx, filter, orderBy)
	if err != nil {
		return nil, fmt.Errorf("queryall: %w", err)
	}

	return options, nil
}

// Count returns the total number of modifier options.
func (b *Business) Count(ctx context.Context, filter QueryFilter) (int, error) {
	return b.storer.Count(ctx, filter)
}

// QueryByID finds the modifier option by the specified ID.
func (b *Business) QueryByID(ctx context.Context, optionID uuid.UUID) (ModifierOption, error) {
	opt, err := b.storer.QueryByID(ctx, optionID)
	if err != nil {
		return ModifierOption{}, fmt.Errorf("query: optionID[%s]: %w", optionID, err)
	}

	return opt, nil
}
