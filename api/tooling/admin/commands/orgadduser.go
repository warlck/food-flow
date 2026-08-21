package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/warlck/food-flow/business/domain/organizationbus"
	"github.com/warlck/food-flow/business/domain/organizationbus/stores/organizationdb"
	"github.com/warlck/food-flow/business/domain/userbus"
	"github.com/warlck/food-flow/business/domain/userbus/stores/userdb"
	"github.com/warlck/food-flow/business/sdk/sqldb"
	"github.com/warlck/food-flow/business/types/role"
	"github.com/warlck/food-flow/foundation/logger"
)

// OrgAddUser adds a user to an organization.
func OrgAddUser(log *logger.Logger, cfg sqldb.Config, orgID string, userID string, roleStr string) error {
	if orgID == "" || userID == "" || roleStr == "" {
		fmt.Println("help: orgadduser <organization_id> <user_id> <role>")
		return ErrHelp
	}

	oID, err := uuid.Parse(orgID)
	if err != nil {
		return fmt.Errorf("parsing org id: %w", err)
	}

	uID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("parsing user id: %w", err)
	}

	r, err := role.Parse(roleStr)
	if err != nil {
		return fmt.Errorf("parsing role: %w", err)
	}

	db, err := sqldb.Open(cfg)
	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	userBus := userbus.NewBusiness(log, userdb.NewStore(log, db))
	orgBus := organizationbus.NewBusiness(log, organizationdb.NewStore(log, db), userBus)

	nou := organizationbus.NewOrganizationUser{
		OrganizationID: oID,
		UserID:         uID,
		Role:           r,
	}

	if _, err := orgBus.AddUser(ctx, nou); err != nil {
		return fmt.Errorf("add user to org: %w", err)
	}

	fmt.Println("User added to organization successfully.")
	return nil
}
