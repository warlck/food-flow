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
	"github.com/warlck/food-flow/foundation/logger"
)

// OrgRemoveUser removes a user from an organization.
func OrgRemoveUser(log *logger.Logger, cfg sqldb.Config, orgID string, userID string) error {
	if orgID == "" || userID == "" {
		fmt.Println("help: orgremoveuser <organization_id> <user_id>")
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

	db, err := sqldb.Open(cfg)
	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	userBus := userbus.NewBusiness(log, userdb.NewStore(log, db))
	orgBus := organizationbus.NewBusiness(log, organizationdb.NewStore(log, db), userBus)

	ou := organizationbus.OrganizationUser{
		OrganizationID: oID,
		UserID:         uID,
	}

	if err := orgBus.RemoveUser(ctx, ou); err != nil {
		return fmt.Errorf("remove user from org: %w", err)
	}

	fmt.Println("User removed from organization successfully.")
	return nil
}
