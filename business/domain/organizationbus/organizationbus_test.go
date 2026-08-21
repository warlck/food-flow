package organizationbus_test

import (
	"context"
	"errors"
	"testing"

	"github.com/warlck/food-flow/business/domain/organizationbus"
	"github.com/warlck/food-flow/business/domain/userbus"
	"github.com/warlck/food-flow/business/sdk/dbtest"
	"github.com/warlck/food-flow/business/types/name"
	"github.com/warlck/food-flow/business/types/role"
)

func Test_Organization(t *testing.T) {
	t.Parallel()

	db := dbtest.New(t, "Test_Organization")
	ctx := context.Background()

	// 1. Seed Users (1 Admin, 1 Regular User)
	adminUsers, err := userbus.TestSeedUsers(ctx, 1, role.Admin, db.BusDomain.User)
	if err != nil {
		t.Fatalf("seeding admin users: %s", err)
	}
	adminUser := adminUsers[0]

	regUsers, err := userbus.TestSeedUsers(ctx, 1, role.User, db.BusDomain.User)
	if err != nil {
		t.Fatalf("seeding regular users: %s", err)
	}
	regUser := regUsers[0]

	// 2. Create Organization
	no := organizationbus.NewOrganization{
		Name: name.MustParse("Acme Food Group"),
	}

	org, err := db.BusDomain.Organization.Create(ctx, no)
	if err != nil {
		t.Fatalf("Create organization: %s", err)
	}

	if org.Name.String() != "Acme Food Group" {
		t.Fatalf("expected name Acme Food Group, got %s", org.Name.String())
	}

	// 3. QueryByID
	fetched, err := db.BusDomain.Organization.QueryByID(ctx, org.ID)
	if err != nil {
		t.Fatalf("QueryByID: %s", err)
	}
	if fetched.ID != org.ID {
		t.Fatalf("expected id %s, got %s", org.ID, fetched.ID)
	}

	// 4. Update Organization (Name)
	newName := name.MustParse("Acme Global Foods")
	uo := organizationbus.UpdateOrganization{
		Name: &newName,
	}

	updated, err := db.BusDomain.Organization.Update(ctx, org, uo)
	if err != nil {
		t.Fatalf("Update organization: %s", err)
	}
	if updated.Name.String() != "Acme Global Foods" {
		t.Fatalf("expected name Acme Global Foods, got %s", updated.Name.String())
	}

	// Verify persistence in DB
	refetched, err := db.BusDomain.Organization.QueryByID(ctx, org.ID)
	if err != nil {
		t.Fatalf("QueryByID after update: %s", err)
	}
	if refetched.Name.String() != "Acme Global Foods" {
		t.Fatalf("expected name persisted as Acme Global Foods, got %s", refetched.Name.String())
	}

	// 5. Add User Rejection (Non-Admin User)
	_, err = db.BusDomain.Organization.AddUser(ctx, organizationbus.NewOrganizationUser{
		OrganizationID: org.ID,
		UserID:         regUser.ID,
		Role:           role.User,
	})
	if !errors.Is(err, organizationbus.ErrNotAuthorized) {
		t.Fatalf("expected ErrNotAuthorized when adding non-admin user, got: %v", err)
	}

	// 6. Add User Success (Admin User)
	ou, err := db.BusDomain.Organization.AddUser(ctx, organizationbus.NewOrganizationUser{
		OrganizationID: org.ID,
		UserID:         adminUser.ID,
		Role:           role.Admin,
	})
	if err != nil {
		t.Fatalf("AddUser admin: %s", err)
	}
	if ou.UserID != adminUser.ID || ou.OrganizationID != org.ID {
		t.Fatalf("unexpected OrganizationUser: %+v", ou)
	}

	// 7. Query Orgs for User
	orgs, err := db.BusDomain.Organization.QueryOrgsForUser(ctx, adminUser.ID)
	if err != nil {
		t.Fatalf("QueryOrgsForUser: %s", err)
	}
	if len(orgs) != 1 || orgs[0].ID != org.ID {
		t.Fatalf("expected 1 organization %s, got %+v", org.ID, orgs)
	}

	// 8. Query Users for Org
	orgUsers, err := db.BusDomain.Organization.QueryUsersForOrg(ctx, org.ID)
	if err != nil {
		t.Fatalf("QueryUsersForOrg: %s", err)
	}
	if len(orgUsers) != 1 || orgUsers[0].UserID != adminUser.ID {
		t.Fatalf("expected 1 user %s in org, got %+v", adminUser.ID, orgUsers)
	}

	// 9. Remove User
	if err := db.BusDomain.Organization.RemoveUser(ctx, ou); err != nil {
		t.Fatalf("RemoveUser: %s", err)
	}

	orgsAfterRemove, err := db.BusDomain.Organization.QueryOrgsForUser(ctx, adminUser.ID)
	if err != nil {
		t.Fatalf("QueryOrgsForUser after remove: %s", err)
	}
	if len(orgsAfterRemove) != 0 {
		t.Fatalf("expected 0 organizations after remove, got %+v", orgsAfterRemove)
	}

	// 10. Delete Organization
	if err := db.BusDomain.Organization.Delete(ctx, org); err != nil {
		t.Fatalf("Delete organization: %s", err)
	}

	_, err = db.BusDomain.Organization.QueryByID(ctx, org.ID)
	if !errors.Is(err, organizationbus.ErrNotFound) {
		t.Fatalf("expected ErrNotFound after delete, got %v", err)
	}
}
