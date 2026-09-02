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
func TestNewAddons(n int, restaurantID uuid.UUID) []NewAddon {
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

		avail := true
		na := NewAddon{
			RestaurantID: restaurantID,
			Name:         name.MustParse(fmt.Sprintf("%s%d", addonNames[nameIdx], idx+i)),
			Description:  addonDescriptions[nameIdx],
			Price:        money.MustParse(price),
			Available:    &avail,
			MaxQuantity:  3,
		}

		newAddons[i] = na
	}

	return newAddons
}

// TestSeedAddons is a helper method for testing.
func TestSeedAddons(ctx context.Context, n int, restaurantID uuid.UUID, bus *Business) ([]Addon, error) {
	newAddons := TestNewAddons(n, restaurantID)

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

// TestAssignAddonsToMenuItem assigns addons to a menu item with ranks.
func TestAssignAddonsToMenuItem(ctx context.Context, menuItemID uuid.UUID, restaurantID uuid.UUID, addonIDs []uuid.UUID, bus *Business) ([]MenuItemAddonInfo, error) {
	assignments := make([]ItemAddonAssignment, len(addonIDs))
	for i, id := range addonIDs {
		rankVal := (i + 1) * 10
		assignments[i] = ItemAddonAssignment{
			AddonID: id,
			Rank:    &rankVal,
		}
	}

	return bus.ReplaceMenuItemAddons(ctx, menuItemID, restaurantID, assignments)
}
