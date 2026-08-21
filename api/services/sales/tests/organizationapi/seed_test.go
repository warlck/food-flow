package organizationapi_test

import (
	"context"
	"fmt"

	"github.com/warlck/food-flow/app/sdk/apitest"
	"github.com/warlck/food-flow/app/sdk/auth"
	"github.com/warlck/food-flow/business/domain/organizationbus"
	"github.com/warlck/food-flow/business/domain/userbus"
	"github.com/warlck/food-flow/business/sdk/dbtest"
	"github.com/warlck/food-flow/business/types/role"
)

func insertSeedData(db *dbtest.Database, ath *auth.Auth) (apitest.SeedData, error) {
	ctx := context.Background()
	busDomain := db.BusDomain

	adminUsrs, err := userbus.TestSeedUsers(ctx, 1, role.Admin, busDomain.User)
	if err != nil {
		return apitest.SeedData{}, fmt.Errorf("seeding admin users: %w", err)
	}

	usrs, err := userbus.TestSeedUsers(ctx, 1, role.User, busDomain.User)
	if err != nil {
		return apitest.SeedData{}, fmt.Errorf("seeding users: %w", err)
	}

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

	tu1 := apitest.User{
		User:  adminUsrs[0],
		Token: apitest.Token(db.BusDomain, ath, adminUsrs[0].Email.Address),
	}

	tu2 := apitest.User{
		User:  usrs[0],
		Token: apitest.Token(db.BusDomain, ath, usrs[0].Email.Address),
	}

	sd := apitest.SeedData{
		Admins:        []apitest.User{tu1},
		Users:         []apitest.User{tu2},
		Organizations: []apitest.Organization{{Organization: orgs[0]}},
	}

	return sd, nil
}
