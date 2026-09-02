package modifieroptionbus

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/google/uuid"
	"github.com/warlck/food-flow/business/types/money"
	"github.com/warlck/food-flow/business/types/name"
)

// TestNewModifierOptions generates a list of new modifier options for testing.
func TestNewModifierOptions(n int, modifierGroupID uuid.UUID, restaurantID uuid.UUID) []NewModifierOption {
	newOptions := make([]NewModifierOption, n)

	idx := rand.Intn(10000)
	for i := range n {
		idx++

		avail := true
		rankVal := (i + 1) * 10
		no := NewModifierOption{
			ModifierGroupID: modifierGroupID,
			RestaurantID:    restaurantID,
			Name:            name.MustParse(fmt.Sprintf("Option%d", idx)),
			Description:     fmt.Sprintf("Description for Option%d", idx),
			PriceDelta:      money.MustParse(float64(i) * 1.50),
			Available:       &avail,
			Rank:            &rankVal,
		}

		newOptions[i] = no
	}

	return newOptions
}

// TestSeedModifierOptions seeds n modifier options for testing.
func TestSeedModifierOptions(ctx context.Context, n int, modifierGroupID uuid.UUID, restaurantID uuid.UUID, api *Business) ([]ModifierOption, error) {
	newOptions := TestNewModifierOptions(n, modifierGroupID, restaurantID)

	options := make([]ModifierOption, len(newOptions))
	for i, no := range newOptions {
		option, err := api.Create(ctx, no)
		if err != nil {
			return nil, fmt.Errorf("seeding modifier option: idx: %d : %w", i, err)
		}

		options[i] = option
	}

	return options, nil
}
