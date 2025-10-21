package restaurantbus

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/warlck/food-flow/business/types/name"
)

// TestNewRestaurants is a helper method for testing.
func TestNewRestaurants(n int) []NewRestaurant {
	newRests := make([]NewRestaurant, n)

	idx := rand.Intn(10000)
	for i := range n {
		idx++

		nr := NewRestaurant{
			Name:        name.MustParse(fmt.Sprintf("Rest%d", idx)),
			Description: fmt.Sprintf("Description for Restaurant%d", idx),
			Address:     fmt.Sprintf("%d Main St", idx),
			Phone:       fmt.Sprintf("+1-555-%04d", idx),
			Email:       fmt.Sprintf("rest%d@test.com", idx),
			ImageURL:    fmt.Sprintf("image%d.jpg", idx),
		}

		newRests[i] = nr
	}

	return newRests
}

// TestSeedRestaurants is a helper method for testing.
func TestSeedRestaurants(ctx context.Context, n int, api *Business) ([]Restaurant, error) {
	newRests := TestNewRestaurants(n)

	rests := make([]Restaurant, len(newRests))
	for i, nr := range newRests {
		rest, err := api.Create(ctx, nr)
		if err != nil {
			return nil, fmt.Errorf("seeding restaurant: idx: %d : %w", i, err)
		}

		rests[i] = rest
	}

	return rests, nil
}
