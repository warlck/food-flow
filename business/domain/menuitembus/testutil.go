package menuitembus

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/google/uuid"
	"github.com/warlck/food-flow/business/types/money"
	"github.com/warlck/food-flow/business/types/name"
)

// TestNewMenuItems is a helper method for testing.
func TestNewMenuItems(n int, categoryID, restaurantID uuid.UUID) []NewMenuItem {
	newItems := make([]NewMenuItem, n)

	idx := rand.Intn(10000)
	for i := range n {
		idx++

		price := float64(idx%100) + 0.99
		ni := NewMenuItem{
			Name:         name.MustParse(fmt.Sprintf("MenuItem%d", idx)),
			Description:  fmt.Sprintf("Description for MenuItem%d", idx),
			Price:        money.MustParse(price),
			CategoryID:   categoryID,
			RestaurantID: restaurantID,
			ImageURL:     fmt.Sprintf("item%d.jpg", idx),
		}

		newItems[i] = ni
	}

	return newItems
}

// TestSeedMenuItems is a helper method for testing.
func TestSeedMenuItems(ctx context.Context, n int, categoryID, restaurantID uuid.UUID, api *Business) ([]MenuItem, error) {
	newItems := TestNewMenuItems(n, categoryID, restaurantID)

	items := make([]MenuItem, len(newItems))
	for i, ni := range newItems {
		item, err := api.Create(ctx, ni)
		if err != nil {
			return nil, fmt.Errorf("seeding menu item: idx: %d : %w", i, err)
		}

		items[i] = item
	}

	return items, nil
}
