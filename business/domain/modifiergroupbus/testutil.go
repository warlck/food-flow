package modifiergroupbus

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/google/uuid"
	"github.com/warlck/food-flow/business/types/name"
)

// TestNewModifierGroups generates a list of new modifier groups for testing.
func TestNewModifierGroups(n int, menuItemID uuid.UUID, restaurantID uuid.UUID) []NewModifierGroup {
	newGroups := make([]NewModifierGroup, n)

	idx := rand.Intn(10000)
	for i := range n {
		idx++

		rankVal := (i + 1) * 10
		ng := NewModifierGroup{
			MenuItemID:    menuItemID,
			RestaurantID:  restaurantID,
			Name:          name.MustParse(fmt.Sprintf("ModifierGroup%d", idx)),
			Description:   fmt.Sprintf("Description for ModifierGroup%d", idx),
			MinSelections: 1,
			MaxSelections: 1,
			Available:     true,
			Rank:          &rankVal,
		}

		newGroups[i] = ng
	}

	return newGroups
}

// TestSeedModifierGroups seeds n modifier groups for testing.
func TestSeedModifierGroups(ctx context.Context, n int, menuItemID uuid.UUID, restaurantID uuid.UUID, api *Business) ([]ModifierGroup, error) {
	newGroups := TestNewModifierGroups(n, menuItemID, restaurantID)

	groups := make([]ModifierGroup, len(newGroups))
	for i, ng := range newGroups {
		group, err := api.Create(ctx, ng)
		if err != nil {
			return nil, fmt.Errorf("seeding modifier group: idx: %d : %w", i, err)
		}

		groups[i] = group
	}

	return groups, nil
}
