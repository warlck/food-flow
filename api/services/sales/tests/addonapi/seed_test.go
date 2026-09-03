package addonapi_test

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

	// Seed organizations
	orgs, err := organizationbus.TestSeedOrganizations(ctx, 2, busDomain.Organization)
	if err != nil {
		return apitest.SeedData{}, fmt.Errorf("seeding organizations: %w", err)
	}

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

	rests, err := restaurantbus.TestSeedRestaurants(ctx, 2, busDomain.Restaurant, orgs[0].ID)
	if err != nil {
		return apitest.SeedData{}, fmt.Errorf("seeding restaurants : %w", err)
	}

	// -------------------------------------------------------------------------

	// Seed categories for restaurant
	cats, err := categorybus.TestSeedCategories(ctx, 2, rests[0].ID, busDomain.Category)
	if err != nil {
		return apitest.SeedData{}, fmt.Errorf("seeding categories : %w", err)
	}

	items1, err := menuitembus.TestSeedMenuItems(ctx, 2, cats[0].ID, rests[0].ID, busDomain.MenuItem)
	if err != nil {
		return apitest.SeedData{}, fmt.Errorf("seeding menu items : %w", err)
	}

	// -------------------------------------------------------------------------

	// Seed addons for menu item 0
	addons1, err := addonbus.TestSeedAddons(ctx, 4, items1[0].ID, rests[0].ID, busDomain.Addon)
	if err != nil {
		return apitest.SeedData{}, fmt.Errorf("seeding addons : %w", err)
	}

	// Seed a restaurant, category, menu item, and addons in the second organization,
	// which the admin users are not members of.
	otherRests, err := restaurantbus.TestSeedRestaurants(ctx, 1, busDomain.Restaurant, orgs[1].ID)
	if err != nil {
		return apitest.SeedData{}, fmt.Errorf("seeding other org restaurant : %w", err)
	}

	otherCats, err := categorybus.TestSeedCategories(ctx, 1, otherRests[0].ID, busDomain.Category)
	if err != nil {
		return apitest.SeedData{}, fmt.Errorf("seeding other org category : %w", err)
	}

	otherItems, err := menuitembus.TestSeedMenuItems(ctx, 1, otherCats[0].ID, otherRests[0].ID, busDomain.MenuItem)
	if err != nil {
		return apitest.SeedData{}, fmt.Errorf("seeding other org menu items : %w", err)
	}

	otherAddons, err := addonbus.TestSeedAddons(ctx, 2, otherItems[0].ID, otherRests[0].ID, busDomain.Addon)
	if err != nil {
		return apitest.SeedData{}, fmt.Errorf("seeding other org addons : %w", err)
	}

	allAddons := append(addons1, otherAddons...)
	appAddons := make([]apitest.Addon, len(allAddons))
	for i, a := range allAddons {
		appAddons[i] = apitest.Addon{Addon: a}
	}

	allItems := append(items1, otherItems...)
	appItems := make([]apitest.MenuItem, len(allItems))
	for i, m := range allItems {
		appItems[i] = apitest.MenuItem{MenuItem: m}
	}

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
		Users:         []apitest.User{tu3, tu4, tu5},
		Admins:        []apitest.User{tu1, tu2},
		Organizations: []apitest.Organization{{Organization: orgs[0]}, {Organization: orgs[1]}},
		Restaurants: []apitest.Restaurant{
			{Restaurant: rests[0]},
			{Restaurant: rests[1]},
			{Restaurant: otherRests[0]},
		},
		Categories: []apitest.Category{
			{Category: cats[0]},
			{Category: cats[1]},
			{Category: otherCats[0]},
		},
		MenuItems: appItems,
		Addons:    appAddons,
	}

	return sd, nil
}
