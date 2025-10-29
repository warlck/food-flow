package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/warlck/food-flow/business/sdk/migrate"
	"github.com/warlck/food-flow/business/sdk/sqldb"
)

// Migrate creates the schema in the database.
func Migrate(cfg sqldb.Config) error {
	db, err := sqldb.Open(cfg)
	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := migrate.Migrate(ctx, db); err != nil {
		return fmt.Errorf("migrate database: %w", err)
	}

	fmt.Println("migrations complete")
	return nil
}
