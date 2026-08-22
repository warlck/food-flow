package menuitemapi_test

import (
	"context"
	"fmt"

	"github.com/warlck/food-flow/app/sdk/apitest"
	"github.com/warlck/food-flow/app/sdk/auth"
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

	// Seed organizations; the second org's resources are used for
	// cross-organization authorization tests (admins are NOT members of it).
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

	// Seed categories for the restaurants
	cats, err := categorybus.TestSeedCategories(ctx, 2, rests[0].ID, busDomain.Category)
	if err != nil {
		return apitest.SeedData{}, fmt.Errorf("seeding categories : %w", err)
	}

	cats2, err := categorybus.TestSeedCategories(ctx, 1, rests[1].ID, busDomain.Category)
	if err != nil {
		return apitest.SeedData{}, fmt.Errorf("seeding categories : %w", err)
	}

	// -------------------------------------------------------------------------

	// Seed menu items for the categories
	items, err := menuitembus.TestSeedMenuItems(ctx, 2, cats[0].ID, rests[0].ID, busDomain.MenuItem)
	if err != nil {
		return apitest.SeedData{}, fmt.Errorf("seeding menu items : %w", err)
	}

	items2, err := menuitembus.TestSeedMenuItems(ctx, 2, cats[1].ID, rests[0].ID, busDomain.MenuItem)
	if err != nil {
		return apitest.SeedData{}, fmt.Errorf("seeding menu items : %w", err)
	}

	items3, err := menuitembus.TestSeedMenuItems(ctx, 2, cats2[0].ID, rests[1].ID, busDomain.MenuItem)
	if err != nil {
		return apitest.SeedData{}, fmt.Errorf("seeding menu items : %w", err)
	}

	// Seed a restaurant, category, and items in the second organization,
	// which the admin users are not members of.
	otherRests, err := restaurantbus.TestSeedRestaurants(ctx, 1, busDomain.Restaurant, orgs[1].ID)
	if err != nil {
		return apitest.SeedData{}, fmt.Errorf("seeding other org restaurant : %w", err)
	}

	otherCats, err := categorybus.TestSeedCategories(ctx, 1, otherRests[0].ID, busDomain.Category)
	if err != nil {
		return apitest.SeedData{}, fmt.Errorf("seeding other org category : %w", err)
	}

	otherItems, err := menuitembus.TestSeedMenuItems(ctx, 2, otherCats[0].ID, otherRests[0].ID, busDomain.MenuItem)
	if err != nil {
		return apitest.SeedData{}, fmt.Errorf("seeding other org menu items : %w", err)
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
			{Category: cats2[0]},
			{Category: otherCats[0]},
		},
		MenuItems: []apitest.MenuItem{
			{MenuItem: items[0]},
			{MenuItem: items[1]},
			{MenuItem: items2[0]},
			{MenuItem: items2[1]},
			{MenuItem: items3[0]},
			{MenuItem: items3[1]},
			{MenuItem: otherItems[0]},
			{MenuItem: otherItems[1]},
		},
	}

	return sd, nil
}
