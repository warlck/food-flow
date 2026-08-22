package restaurantapi_test

import (
	"context"
	"fmt"

	"github.com/warlck/food-flow/app/sdk/apitest"
	"github.com/warlck/food-flow/app/sdk/auth"
	"github.com/warlck/food-flow/business/domain/addonbus"
	"github.com/warlck/food-flow/business/domain/categorybus"
	"github.com/warlck/food-flow/business/domain/menuitembus"
	"github.com/warlck/food-flow/business/domain/organizationbus"
	"github.com/warlck/food-flow/business/domain/restaurantbus"
	"github.com/warlck/food-flow/business/domain/userbus"
	"github.com/warlck/food-flow/business/sdk/dbtest"
	"github.com/warlck/food-flow/business/types/role"
)

func insertSeedData(db *dbtest.Database, ath *auth.Auth) (apitest.SeedData, error) {
	ctx := context.Background()
	busDomain := db.BusDomain

	// Create admin users for auth
	adminUsrs, err := userbus.TestSeedUsers(ctx, 2, role.Admin, busDomain.User)
	if err != nil {
		return apitest.SeedData{}, fmt.Errorf("seeding admin users : %w", err)
	}

	// -------------------------------------------------------------------------

	// Create regular users
	usrs, err := userbus.TestSeedUsers(ctx, 3, role.User, busDomain.User)
	if err != nil {
		return apitest.SeedData{}, fmt.Errorf("seeding users : %w", err)
	}

	// -------------------------------------------------------------------------

	// Seed restaurants
	orgs, err := organizationbus.TestSeedOrganizations(ctx, 1, busDomain.Organization)
	if err != nil {
		return apitest.SeedData{}, fmt.Errorf("seeding organizations: %w", err)
	}

	// Associate admins with organization
	if _, err := busDomain.Organization.AddUser(ctx, organizationbus.NewOrganizationUser{
		OrganizationID: orgs[0].ID,
		UserID:         adminUsrs[0].ID,
		Role:           role.Admin,
	}); err != nil {
		return apitest.SeedData{}, fmt.Errorf("adding user to organization: %w", err)
	}
	if _, err := busDomain.Organization.AddUser(ctx, organizationbus.NewOrganizationUser{
		OrganizationID: orgs[0].ID,
		UserID:         adminUsrs[1].ID,
		Role:           role.Admin,
	}); err != nil {
		return apitest.SeedData{}, fmt.Errorf("adding user to organization: %w", err)
	}

	rests, err := restaurantbus.TestSeedRestaurants(ctx, 4, busDomain.Restaurant, orgs[0].ID)
	if err != nil {
		return apitest.SeedData{}, fmt.Errorf("seeding restaurants : %w", err)
	}

	// -------------------------------------------------------------------------

	// Seed categories for the first restaurant (for detailed query test)
	cats, err := categorybus.TestSeedCategories(ctx, 2, rests[0].ID, busDomain.Category)
	if err != nil {
		return apitest.SeedData{}, fmt.Errorf("seeding categories : %w", err)
	}

	// -------------------------------------------------------------------------

	// Seed menu items for the categories
	items, err := menuitembus.TestSeedMenuItems(ctx, 3, cats[0].ID, rests[0].ID, busDomain.MenuItem)
	if err != nil {
		return apitest.SeedData{}, fmt.Errorf("seeding menu items : %w", err)
	}

	items2, err := menuitembus.TestSeedMenuItems(ctx, 2, cats[1].ID, rests[0].ID, busDomain.MenuItem)
	if err != nil {
		return apitest.SeedData{}, fmt.Errorf("seeding menu items : %w", err)
	}

	// -------------------------------------------------------------------------

	// Seed a dedicated ranked category on the second restaurant so the
	// details-endpoint rank ordering can be tested without disturbing the
	// unranked seeds that assert backward compatibility on rests[0].
	rankedCats, err := categorybus.TestSeedCategories(ctx, 1, rests[1].ID, busDomain.Category)
	if err != nil {
		return apitest.SeedData{}, fmt.Errorf("seeding ranked category : %w", err)
	}

	rankedItems, err := menuitembus.TestSeedMenuItems(ctx, 3, rankedCats[0].ID, rests[1].ID, busDomain.MenuItem)
	if err != nil {
		return apitest.SeedData{}, fmt.Errorf("seeding ranked menu items : %w", err)
	}

	// rankedItems[0] -> rank 20, rankedItems[1] -> rank 10, rankedItems[2] -> unranked.
	if rankedItems[0], err = busDomain.MenuItem.Update(ctx, rankedItems[0], menuitembus.UpdateMenuItem{Rank: dbtest.IntPointer(20)}); err != nil {
		return apitest.SeedData{}, fmt.Errorf("ranking menu item 0 : %w", err)
	}
	if rankedItems[1], err = busDomain.MenuItem.Update(ctx, rankedItems[1], menuitembus.UpdateMenuItem{Rank: dbtest.IntPointer(10)}); err != nil {
		return apitest.SeedData{}, fmt.Errorf("ranking menu item 1 : %w", err)
	}

	rankedAddons, err := addonbus.TestSeedAddons(ctx, 3, rankedCats[0].ID, rests[1].ID, busDomain.Addon)
	if err != nil {
		return apitest.SeedData{}, fmt.Errorf("seeding ranked addons : %w", err)
	}

	// rankedAddons[0] -> rank 20, rankedAddons[1] -> rank 10, rankedAddons[2] -> unranked.
	if rankedAddons[0], err = busDomain.Addon.Update(ctx, rankedAddons[0], addonbus.UpdateAddon{Rank: dbtest.IntPointer(20)}); err != nil {
		return apitest.SeedData{}, fmt.Errorf("ranking addon 0 : %w", err)
	}
	if rankedAddons[1], err = busDomain.Addon.Update(ctx, rankedAddons[1], addonbus.UpdateAddon{Rank: dbtest.IntPointer(10)}); err != nil {
		return apitest.SeedData{}, fmt.Errorf("ranking addon 1 : %w", err)
	}

	// -------------------------------------------------------------------------

	tu1 := apitest.User{
		User:  adminUsrs[0],
		Token: apitest.Token(db.BusDomain, ath, adminUsrs[0].Email.Address),
	}

	tu2 := apitest.User{
		User:  adminUsrs[1],
		Token: apitest.Token(db.BusDomain, ath, adminUsrs[1].Email.Address),
	}

	tu3 := apitest.User{
		User:  usrs[0],
		Token: apitest.Token(db.BusDomain, ath, usrs[0].Email.Address),
	}

	tu4 := apitest.User{
		User:  usrs[1],
		Token: apitest.Token(db.BusDomain, ath, usrs[1].Email.Address),
	}

	tu5 := apitest.User{
		User:  usrs[2],
		Token: apitest.Token(db.BusDomain, ath, usrs[2].Email.Address),
	}
	sd := apitest.SeedData{
		Users:  []apitest.User{tu3, tu4, tu5},
		Admins: []apitest.User{tu1, tu2},
		Restaurants: []apitest.Restaurant{
			{Restaurant: rests[0]},
			{Restaurant: rests[1]},
			{Restaurant: rests[2]},
			{Restaurant: rests[3]},
		},
		Categories: []apitest.Category{
			{Category: cats[0]},
			{Category: cats[1]},
			{Category: rankedCats[0]},
		},
		MenuItems: []apitest.MenuItem{
			{MenuItem: items[0]},
			{MenuItem: items[1]},
			{MenuItem: items[2]},
			{MenuItem: items2[0]},
			{MenuItem: items2[1]},
			{MenuItem: rankedItems[0]},
			{MenuItem: rankedItems[1]},
			{MenuItem: rankedItems[2]},
		},
		Addons: []apitest.Addon{
			{Addon: rankedAddons[0]},
			{Addon: rankedAddons[1]},
			{Addon: rankedAddons[2]},
		},
	}

	return sd, nil
}
