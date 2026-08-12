package orderapi_test

import (
	"context"
	"fmt"

	"github.com/warlck/food-flow/app/sdk/apitest"
	"github.com/warlck/food-flow/app/sdk/auth"
	"github.com/warlck/food-flow/business/domain/addonbus"
	"github.com/warlck/food-flow/business/domain/categorybus"
	"github.com/warlck/food-flow/business/domain/menuitembus"
	"github.com/warlck/food-flow/business/domain/orderbus"
	"github.com/warlck/food-flow/business/domain/promobus"
	"github.com/warlck/food-flow/business/domain/restaurantbus"
	"github.com/warlck/food-flow/business/domain/userbus"
	"github.com/warlck/food-flow/business/sdk/dbtest"
	"github.com/warlck/food-flow/business/types/name"
	"github.com/warlck/food-flow/business/types/role"
)

func insertSeedData(db *dbtest.Database, ath *auth.Auth) (apitest.SeedData, error) {
	ctx := context.Background()
	busDomain := db.BusDomain

	// -------------------------------------------------------------------------
	// Create admin users

	admins, err := userbus.TestSeedUsers(ctx, 1, role.Admin, busDomain.User)
	if err != nil {
		return apitest.SeedData{}, fmt.Errorf("seeding admin users: %w", err)
	}

	tu1 := apitest.User{
		User:  admins[0],
		Token: apitest.Token(db.BusDomain.User, ath, admins[0].Email.Address),
	}

	// -------------------------------------------------------------------------
	// Create regular users

	users, err := userbus.TestSeedUsers(ctx, 2, role.User, busDomain.User)
	if err != nil {
		return apitest.SeedData{}, fmt.Errorf("seeding users: %w", err)
	}

	tu2 := apitest.User{
		User:  users[0],
		Token: apitest.Token(db.BusDomain.User, ath, users[0].Email.Address),
	}

	tu3 := apitest.User{
		User:  users[1],
		Token: apitest.Token(db.BusDomain.User, ath, users[1].Email.Address),
	}

	// -------------------------------------------------------------------------
	// Create restaurant

	restaurants, err := restaurantbus.TestSeedRestaurants(ctx, 1, busDomain.Restaurant)
	if err != nil {
		return apitest.SeedData{}, fmt.Errorf("seeding restaurants: %w", err)
	}

	tr1 := apitest.Restaurant{
		Restaurant: restaurants[0],
	}

	// -------------------------------------------------------------------------
	// Create categories

	categories, err := categorybus.TestSeedCategories(ctx, 1, restaurants[0].ID, busDomain.Category)
	if err != nil {
		return apitest.SeedData{}, fmt.Errorf("seeding categories: %w", err)
	}

	tc1 := apitest.Category{
		Category: categories[0],
	}

	// -------------------------------------------------------------------------
	// Create menu items

	menuItems, err := menuitembus.TestSeedMenuItems(ctx, 2, categories[0].ID, restaurants[0].ID, busDomain.MenuItem)
	if err != nil {
		return apitest.SeedData{}, fmt.Errorf("seeding menu items: %w", err)
	}

	tmi1 := apitest.MenuItem{
		MenuItem: menuItems[0],
	}

	tmi2 := apitest.MenuItem{
		MenuItem: menuItems[1],
	}

	// -------------------------------------------------------------------------
	// Create addons

	addons, err := addonbus.TestSeedAddons(ctx, 2, categories[0].ID, restaurants[0].ID, busDomain.Addon)
	if err != nil {
		return apitest.SeedData{}, fmt.Errorf("seeding addons: %w", err)
	}

	ta1 := apitest.Addon{
		Addon: addons[0],
	}

	ta2 := apitest.Addon{
		Addon: addons[1],
	}

	// -------------------------------------------------------------------------
	// Create orders

	orders, err := orderbus.TestSeedOrders(ctx, 2, restaurants[0].ID, menuItems, busDomain.Order)
	if err != nil {
		return apitest.SeedData{}, fmt.Errorf("seeding orders: %w", err)
	}

	to1 := apitest.Order{
		Order: orders[0],
	}

	to2 := apitest.Order{
		Order: orders[1],
	}

	// Create promotions (with 0.0 MinOrderAmount so integration tests aren't flaky regardless of seeded menu item prices)
	_, err = busDomain.Promo.Create(ctx, promobus.NewPromotion{
		Code:           "TESTPROMO1",
		Name:           name.MustParse("Test Promo 1"),
		Description:    "Test promo with no min order amount",
		DiscountType:   promobus.DiscountTypePercentage,
		DiscountValue:  10.0,
		MinOrderAmount: 0.0,
		Enabled:        true,
	})
	if err != nil {
		return apitest.SeedData{}, fmt.Errorf("seeding promotions: %w", err)
	}

	// -------------------------------------------------------------------------

	sd := apitest.SeedData{
		Admins:      []apitest.User{tu1},
		Users:       []apitest.User{tu2, tu3},
		Restaurants: []apitest.Restaurant{tr1},
		Categories:  []apitest.Category{tc1},
		MenuItems:   []apitest.MenuItem{tmi1, tmi2},
		Addons:      []apitest.Addon{ta1, ta2},
		Orders:      []apitest.Order{to1, to2},
	}

	return sd, nil
}
