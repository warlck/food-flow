package addonapi_test

import (
	"context"
	"fmt"

	"github.com/warlck/food-flow/app/sdk/apitest"
	"github.com/warlck/food-flow/app/sdk/auth"
	"github.com/warlck/food-flow/business/domain/addonbus"
	"github.com/warlck/food-flow/business/domain/categorybus"
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

	// -------------------------------------------------------------------------

	// Seed addons for category
	addons, err := addonbus.TestSeedAddons(ctx, 2, cats[0].ID, rests[0].ID, busDomain.Addon)
	if err != nil {
		return apitest.SeedData{}, fmt.Errorf("seeding addons : %w", err)
	}

	appAddons := make([]apitest.Addon, len(addons))
	for i, a := range addons {
		appAddons[i] = apitest.Addon{Addon: a}
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
		Users:  []apitest.User{tu3, tu4, tu5},
		Admins: []apitest.User{tu1, tu2},
		Restaurants: []apitest.Restaurant{
			{Restaurant: rests[0]},
			{Restaurant: rests[1]},
		},
		Categories: []apitest.Category{
			{Category: cats[0]},
			{Category: cats[1]},
		},
		Addons: appAddons,
	}

	return sd, nil
}
