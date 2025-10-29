package commands

import (
	"context"
	"fmt"
	"net/mail"
	"time"

	"github.com/warlck/food-flow/business/domain/userbus"
	"github.com/warlck/food-flow/business/domain/userbus/stores/userdb"
	"github.com/warlck/food-flow/business/sdk/sqldb"
	"github.com/warlck/food-flow/business/types/name"
	"github.com/warlck/food-flow/business/types/password"
	"github.com/warlck/food-flow/business/types/role"
	"github.com/warlck/food-flow/foundation/logger"
)

// UserAdd adds new users into the database.
func UserAdd(log *logger.Logger, cfg sqldb.Config, nme string, email string, pass string) error {
	if nme == "" || email == "" || pass == "" {
		fmt.Println("help: useradd <name> <email> <password>")
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

	addr, err := mail.ParseAddress(email)
	if err != nil {
		return fmt.Errorf("parsing email: %w", err)
	}

	nu := userbus.NewUser{
		Name:     name.MustParse(nme),
		Email:    *addr,
		Password: password.MustParse(pass),
		Roles:    []role.Role{role.Admin, role.User},
	}

	usr, err := userBus.Create(ctx, nu)
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}

	fmt.Println("user id:", usr.ID)
	return nil
}
