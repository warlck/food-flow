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
func TestNewAddons(n int, categoryID, restaurantID uuid.UUID) []NewAddon {
	newAddons := make([]NewAddon, n)

	addonNames := []string{
		"Extra Cheese", "Extra Meat", "Jalapeños", "Garlic Sauce",
		"Hummus", "Feta Cheese", "Grilled Vegetables", "Spicy Sauce",
		"Extra Sauce", "Avocado", "Bacon", "Mushrooms",
	}

	addonDescriptions := []string{
		"Additional melted cheese", "Additional portion of meat", "Spicy jalapeño peppers", "Extra garlic sauce on the side",
		"Side of creamy hummus", "Crumbled feta cheese", "Assorted grilled vegetables", "Hot chili sauce",
		"Extra sauce portion", "Fresh avocado slices", "Crispy bacon strips", "Sautéed mushrooms",
	}

	addonPrices := []float64{2.00, 4.00, 1.50, 1.00, 3.00, 2.50, 3.50, 1.00, 1.50, 3.00, 2.50, 2.00}

	idx := rand.Intn(10000)
	for i := range n {
		nameIdx := i % len(addonNames)
		price := addonPrices[nameIdx]

		na := NewAddon{
			CategoryID:   categoryID,
			RestaurantID: restaurantID,
			Name:         name.MustParse(fmt.Sprintf("%s%d", addonNames[nameIdx], idx+i)),
			Description:  addonDescriptions[nameIdx],
			Price:        money.MustParse(price),
			MaxQuantity:  3,
		}

		newAddons[i] = na
	}

	return newAddons
}

// TestSeedAddons is a helper method for testing.
func TestSeedAddons(ctx context.Context, n int, categoryID, restaurantID uuid.UUID, bus *Business) ([]Addon, error) {
	newAddons := TestNewAddons(n, categoryID, restaurantID)

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

// TestSeedAddonsForMenuItem is a helper method for testing that creates
// addons with specific names for a category.
func TestSeedAddonsForMenuItem(ctx context.Context, categoryID, restaurantID uuid.UUID, bus *Business) ([]Addon, error) {
	addonData := []struct {
		name        string
		description string
		price       float64
		maxQty      int
	}{
		{"Extra Cheese", "Additional melted cheese", 2.00, 3},
		{"Extra Meat", "Additional portion of meat", 4.00, 2},
		{"Garlic Sauce", "Extra garlic sauce on the side", 1.00, 3},
		{"Spicy Sauce", "Hot chili sauce", 1.00, 3},
	}

	addons := make([]Addon, len(addonData))
	for i, data := range addonData {
		na := NewAddon{
			CategoryID:   categoryID,
			RestaurantID: restaurantID,
			Name:         name.MustParse(data.name),
			Description:  data.description,
			Price:        money.MustParse(data.price),
			MaxQuantity:  data.maxQty,
		}

		addon, err := bus.Create(ctx, na)
		if err != nil {
			return nil, fmt.Errorf("seeding addon %s: %w", data.name, err)
		}

		addons[i] = addon
	}

	return addons, nil
}
