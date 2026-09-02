package addonbus

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/google/uuid"
	"github.com/warlck/food-flow/business/types/money"
	"github.com/warlck/food-flow/business/types/name"
)

// TestNewAddons is a helper method for testing.
func TestNewAddons(n int, menuItemID uuid.UUID, restaurantID uuid.UUID) []NewAddon {
	newAddons := make([]NewAddon, n)

	addonNames := []string{
		"Extra Cheese", "Extra Meat", "Jalapenos", "Garlic Sauce",
		"Hummus", "Feta Cheese", "Grilled Vegetables", "Spicy Sauce",
		"Extra Sauce", "Avocado", "Bacon", "Mushrooms",
	}

	addonDescriptions := []string{
		"Additional melted cheese", "Additional portion of meat", "Spicy jalapeno peppers", "Extra garlic sauce on the side",
		"Side of creamy hummus", "Crumbled feta cheese", "Assorted grilled vegetables", "Hot chili sauce",
		"Extra sauce portion", "Fresh avocado slices", "Crispy bacon strips", "Sauteed mushrooms",
	}

	addonPrices := []float64{2.00, 4.00, 1.50, 1.00, 3.00, 2.50, 3.50, 1.00, 1.50, 3.00, 2.50, 2.00}

	idx := rand.Intn(10000)
	for i := range n {
		nameIdx := i % len(addonNames)
		price := addonPrices[nameIdx]

		avail := true
		rankVal := (i + 1) * 10
		na := NewAddon{
			MenuItemID:   menuItemID,
			RestaurantID: restaurantID,
			Name:         name.MustParse(fmt.Sprintf("%s%d", addonNames[nameIdx], idx+i)),
			Description:  addonDescriptions[nameIdx],
			Price:        money.MustParse(price),
			Available:    &avail,
			MaxQuantity:  3,
			Rank:         &rankVal,
		}

		newAddons[i] = na
	}

	return newAddons
}

// TestSeedAddons is a helper method for testing.
func TestSeedAddons(ctx context.Context, n int, menuItemID uuid.UUID, restaurantID uuid.UUID, bus *Business) ([]Addon, error) {
	newAddons := TestNewAddons(n, menuItemID, restaurantID)

	addons := make([]Addon, len(newAddons))
	for i, na := range newAddons {
		addon, err := bus.Create(ctx, na)
		if err != nil {
			return nil, fmt.Errorf("seeding addon: idx: %d : %w", i, err)
		}

		addons[i] = addon
	}

	return addons, nil
}
