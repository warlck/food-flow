package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/warlck/food-flow/business/domain/organizationbus"
	"github.com/warlck/food-flow/business/domain/organizationbus/stores/organizationdb"
	"github.com/warlck/food-flow/business/domain/userbus"
	"github.com/warlck/food-flow/business/domain/userbus/stores/userdb"
	"github.com/warlck/food-flow/business/sdk/sqldb"
	"github.com/warlck/food-flow/business/types/name"
	"github.com/warlck/food-flow/foundation/logger"
)

// OrgAdd adds a new organization into the database.
func OrgAdd(log *logger.Logger, cfg sqldb.Config, nme string) error {
	if nme == "" {
		fmt.Println("help: orgadd <name>")
		return ErrHelp
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

	parsedName, err := name.Parse(nme)
	if err != nil {
		return fmt.Errorf("parsing name: %w", err)
	}

	no := organizationbus.NewOrganization{
		Name: parsedName,
	}

	org, err := orgBus.Create(ctx, no)
	if err != nil {
		return fmt.Errorf("create organization: %w", err)
	}

	fmt.Println("organization id:", org.ID)
	return nil
}
