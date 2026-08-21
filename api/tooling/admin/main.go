package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/ardanlabs/conf/v3"
	"github.com/google/uuid"
	"github.com/warlck/food-flow/api/tooling/admin/commands"
	"github.com/warlck/food-flow/business/sdk/sqldb"
	"github.com/warlck/food-flow/foundation/logger"
)

var build = "develop"

type config struct {
	conf.Version
	Args conf.Args
	DB   struct {
		User         string `conf:"default:postgres"`
		Password     string `conf:"default:postgres,mask"`
		Host         string `conf:"default:database-service"`
		Name         string `conf:"default:postgres"`
		MaxIdleConns int    `conf:"default:0"`
		MaxOpenConns int    `conf:"default:0"`
		DisableTLS   bool   `conf:"default:true"`
	}
	Auth struct {
		KeysFolder string `conf:"default:infra/keys/"`
		DefaultKID string `conf:"default:local-dev"`
	}
}

func main() {
	log := logger.New(io.Discard, logger.LevelInfo, "ADMIN", func(context.Context) string { return "00000000-0000-0000-0000-000000000000" })

	if err := run(log); err != nil {
		if !errors.Is(err, commands.ErrHelp) {
			fmt.Println("msg", err)
		}
		os.Exit(1)
	}
}

func run(log *logger.Logger) error {
	cfg := config{
		Version: conf.Version{
			Build: build,
			Desc:  "copyright information here",
		},
	}

	const prefix = "SALES"
	help, err := conf.Parse(prefix, &cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
			return nil
		}

		out, err := conf.String(&cfg)
		if err != nil {
			return fmt.Errorf("generating config for output: %w", err)
		}
		log.Info(context.Background(), "startup", "config", out)

		return fmt.Errorf("parsing config: %w", err)
	}

	return processCommands(cfg.Args, log, cfg)
}

// processCommands handles the execution of the commands specified on
// the command line.
func processCommands(args conf.Args, log *logger.Logger, cfg config) error {
	dbConfig := sqldb.Config{
		User:         cfg.DB.User,
		Password:     cfg.DB.Password,
		Host:         cfg.DB.Host,
		Name:         cfg.DB.Name,
		MaxIdleConns: cfg.DB.MaxIdleConns,
		MaxOpenConns: cfg.DB.MaxOpenConns,
		DisableTLS:   cfg.DB.DisableTLS,
	}

	switch args.Num(0) {
	case "migrate":
		if err := commands.Migrate(dbConfig); err != nil {
			return fmt.Errorf("migrating database: %w", err)
		}

	case "seed":
		if err := commands.Seed(dbConfig); err != nil {
			return fmt.Errorf("seeding database: %w", err)
		}

	case "migrate-seed":
		if err := commands.Migrate(dbConfig); err != nil {
			return fmt.Errorf("migrating database: %w", err)
		}
		if err := commands.Seed(dbConfig); err != nil {
			return fmt.Errorf("seeding database: %w", err)
		}

	case "useradd":
		name := args.Num(1)
		email := args.Num(2)
		password := args.Num(3)
		if err := commands.UserAdd(log, dbConfig, name, email, password); err != nil {
			return fmt.Errorf("adding user: %w", err)
		}

	case "users":
		pageNumber := args.Num(1)
		rowsPerPage := args.Num(2)
		if err := commands.Users(log, dbConfig, pageNumber, rowsPerPage); err != nil {
			return fmt.Errorf("getting users: %w", err)
		}

	case "genkey":
		if err := commands.GenKey(); err != nil {
			return fmt.Errorf("key generation: %w", err)
		}

	case "gentoken":
		userID, err := uuid.Parse(args.Num(1))
		if err != nil {
			return fmt.Errorf("generating token: %w", err)
		}
		kid := args.Num(2)
		if kid == "" {
			kid = cfg.Auth.DefaultKID
		}
		if err := commands.GenToken(log, dbConfig, cfg.Auth.KeysFolder, userID, kid); err != nil {
			return fmt.Errorf("generating token: %w", err)
		}

	case "orgadd":
		name := args.Num(1)
		if err := commands.OrgAdd(log, dbConfig, name); err != nil {
			return fmt.Errorf("adding organization: %w", err)
		}

	case "orgadduser":
		orgID := args.Num(1)
		userID := args.Num(2)
		roleStr := args.Num(3)
		if err := commands.OrgAddUser(log, dbConfig, orgID, userID, roleStr); err != nil {
			return fmt.Errorf("adding user to organization: %w", err)
		}

	case "orgremoveuser":
		orgID := args.Num(1)
		userID := args.Num(2)
		if err := commands.OrgRemoveUser(log, dbConfig, orgID, userID); err != nil {
			return fmt.Errorf("removing user from organization: %w", err)
		}

	default:
		fmt.Println("migrate:    create the schema in the database")
		fmt.Println("seed:       add data to the database")
		fmt.Println("useradd:    add a new user to the database")
		fmt.Println("users:      get a list of users from the database")
		fmt.Println("genkey:     generate a set of private/public key files")
		fmt.Println("gentoken:   generate a JWT for a user with claims")
		fmt.Println("orgadd:     add a new organization to the database")
		fmt.Println("orgadduser: add a user to an organization")
		fmt.Println("orgremoveuser: remove a user from an organization")
		fmt.Println("provide a command to get more help.")
		return commands.ErrHelp
	}

	return nil
}
