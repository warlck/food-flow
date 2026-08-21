package imageapi_test

import (
	"context"
	"fmt"

	"github.com/warlck/food-flow/app/sdk/apitest"
	"github.com/warlck/food-flow/app/sdk/auth"
	"github.com/warlck/food-flow/business/domain/organizationbus"
	"github.com/warlck/food-flow/business/domain/restaurantbus"
	"github.com/warlck/food-flow/business/domain/userbus"
	"github.com/warlck/food-flow/business/sdk/dbtest"
	"github.com/warlck/food-flow/business/types/role"
)

func insertSeedData(db *dbtest.Database, ath *auth.Auth) (apitest.SeedData, error) {
	ctx := context.Background()
	busDomain := db.BusDomain

	// Create admin user
	adminUsrs, err := userbus.TestSeedUsers(ctx, 1, role.Admin, busDomain.User)
	if err != nil {
		return apitest.SeedData{}, fmt.Errorf("seeding admin users : %w", err)
	}

	// Create regular user
	usrs, err := userbus.TestSeedUsers(ctx, 1, role.User, busDomain.User)
	if err != nil {
		return apitest.SeedData{}, fmt.Errorf("seeding users : %w", err)
	}

	// Create restaurant
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

	rests, err := restaurantbus.TestSeedRestaurants(ctx, 1, busDomain.Restaurant, orgs[0].ID)
	if err != nil {
		return apitest.SeedData{}, fmt.Errorf("seeding restaurants : %w", err)
	}

	tr1 := apitest.Restaurant{
		Restaurant: rests[0],
	}

	

	tu1 := apitest.User{
		User:  adminUsrs[0],
		Token: apitest.Token(db.BusDomain, ath, adminUsrs[0].Email.Address),
	}

	tu2 := apitest.User{
		User:  usrs[0],
		Token: apitest.Token(db.BusDomain, ath, usrs[0].Email.Address),
	}
	sd := apitest.SeedData{
		Admins:      []apitest.User{tu1},
		Users:       []apitest.User{tu2},
		Restaurants: []apitest.Restaurant{tr1},
	}

	return sd, nil
}
