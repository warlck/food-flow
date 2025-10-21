package categorybus

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/google/uuid"
	"github.com/warlck/food-flow/business/types/name"
)

// TestNewCategories is a helper method for testing.
func TestNewCategories(n int, restaurantID uuid.UUID) []NewCategory {
	newCats := make([]NewCategory, n)

	idx := rand.Intn(10000)
	for i := range n {
		idx++

		nc := NewCategory{
			Name:         name.MustParse(fmt.Sprintf("Category%d", idx)),
			Description:  fmt.Sprintf("Description for Category%d", idx),
			RestaurantID: restaurantID,
		}

		newCats[i] = nc
	}

	return newCats
}

// TestSeedCategories is a helper method for testing.
func TestSeedCategories(ctx context.Context, n int, restaurantID uuid.UUID, api *Business) ([]Category, error) {
	newCats := TestNewCategories(n, restaurantID)

	cats := make([]Category, len(newCats))
	for i, nc := range newCats {
		cat, err := api.Create(ctx, nc)
		if err != nil {
			return nil, fmt.Errorf("seeding category: idx: %d : %w", i, err)
		}

		cats[i] = cat
	}

	return cats, nil
}
