// Package modifieroptiondb contains modifier option related CRUD functionality.
package modifieroptiondb

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/warlck/food-flow/business/domain/modifieroptionbus"
	"github.com/warlck/food-flow/business/sdk/order"
	"github.com/warlck/food-flow/business/sdk/page"
	"github.com/warlck/food-flow/business/sdk/sqldb"
	"github.com/warlck/food-flow/foundation/logger"
)

// Store manages the set of APIs for modifier option database access.
type Store struct {
	log *logger.Logger
	db  sqlx.ExtContext
}

// NewStore constructs the api for data access.
func NewStore(log *logger.Logger, db *sqlx.DB) *Store {
	return &Store{
		log: log,
		db:  db,
	}
}

// Create inserts a new modifier option into the database.
func (s *Store) Create(ctx context.Context, opt modifieroptionbus.ModifierOption) error {
	const q = `
	INSERT INTO modifier_options
		(modifier_option_id, modifier_group_id, restaurant_id, name, description, price_delta, available, rank, date_created, date_updated)
	VALUES
		(:modifier_option_id, :modifier_group_id, :restaurant_id, :name, :description, :price_delta, :available, :rank, :date_created, :date_updated)`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, toDBModifierOption(opt)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// Update replaces a modifier option document in the database.
func (s *Store) Update(ctx context.Context, opt modifieroptionbus.ModifierOption) error {
	const q = `
	UPDATE
		modifier_options
	SET 
		"name" = :name,
		"description" = :description,
		"price_delta" = :price_delta,
		"available" = :available,
		"rank" = :rank,
		"date_updated" = :date_updated
	WHERE
		modifier_option_id = :modifier_option_id`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, toDBModifierOption(opt)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// Reorder updates the rank of a list of modifier options atomically in a transaction.
func (s *Store) Reorder(ctx context.Context, options []modifieroptionbus.ModifierOption) error {
	tx, err := s.db.(*sqlx.DB).BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	const q = `
	UPDATE
		modifier_options
	SET
		"rank" = :rank,
		"date_updated" = :date_updated
	WHERE
		modifier_option_id = :modifier_option_id`

	for _, opt := range options {
		if err := sqldb.NamedExecContext(ctx, s.log, tx, q, toDBModifierOption(opt)); err != nil {
			return fmt.Errorf("namedexeccontext: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit: %w", err)
	}

	return nil
}

// Delete removes a modifier option from the database.
func (s *Store) Delete(ctx context.Context, opt modifieroptionbus.ModifierOption) error {
	data := struct {
		ID string `db:"modifier_option_id"`
	}{
		ID: opt.ID.String(),
	}

	const q = `
	DELETE FROM
		modifier_options
	WHERE
		modifier_option_id = :modifier_option_id`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, data); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// Query retrieves a list of existing modifier options from the database.
func (s *Store) Query(ctx context.Context, filter modifieroptionbus.QueryFilter, orderBy order.By, page page.Page) ([]modifieroptionbus.ModifierOption, error) {
	data := map[string]any{
		"offset":        (page.Number() - 1) * page.RowsPerPage(),
		"rows_per_page": page.RowsPerPage(),
	}

	const q = `
	SELECT
		modifier_option_id, modifier_group_id, restaurant_id, name, description, price_delta, available, rank, date_created, date_updated
	FROM
		modifier_options`

	buf := bytes.NewBufferString(q)
	applyFilter(filter, data, buf)

	orderByClause, err := orderByClause(orderBy)
	if err != nil {
		return nil, err
	}

	buf.WriteString(orderByClause)
	buf.WriteString(" OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY")

	var dbOptions []modifierOption
	if err := sqldb.NamedQuerySlice(ctx, s.log, s.db, buf.String(), data, &dbOptions); err != nil {
		return nil, fmt.Errorf("namedqueryslice: %w", err)
	}

	return toBusModifierOptions(dbOptions)
}

// QueryAll retrieves all modifier options matching the filter from the database without pagination.
func (s *Store) QueryAll(ctx context.Context, filter modifieroptionbus.QueryFilter, orderBy order.By) ([]modifieroptionbus.ModifierOption, error) {
	data := map[string]any{}

	const q = `
	SELECT
		modifier_option_id, modifier_group_id, restaurant_id, name, description, price_delta, available, rank, date_created, date_updated
	FROM
		modifier_options`

	buf := bytes.NewBufferString(q)
	applyFilter(filter, data, buf)

	orderByClause, err := orderByClause(orderBy)
	if err != nil {
		return nil, err
	}

	buf.WriteString(orderByClause)

	var dbOptions []modifierOption
	if err := sqldb.NamedQuerySlice(ctx, s.log, s.db, buf.String(), data, &dbOptions); err != nil {
		return nil, fmt.Errorf("namedqueryslice: %w", err)
	}

	return toBusModifierOptions(dbOptions)
}

// Count returns the total number of modifier options in the DB.
func (s *Store) Count(ctx context.Context, filter modifieroptionbus.QueryFilter) (int, error) {
	data := map[string]any{}

	const q = `
	SELECT
		count(1)
	FROM
		modifier_options`

	buf := bytes.NewBufferString(q)
	applyFilter(filter, data, buf)

	var count struct {
		Count int `db:"count"`
	}
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, buf.String(), data, &count); err != nil {
		return 0, fmt.Errorf("db: %w", err)
	}

	return count.Count, nil
}

// QueryByID gets the specified modifier option from the database.
func (s *Store) QueryByID(ctx context.Context, optionID uuid.UUID) (modifieroptionbus.ModifierOption, error) {
	data := struct {
		ID string `db:"modifier_option_id"`
	}{
		ID: optionID.String(),
	}

	const q = `
	SELECT
		modifier_option_id, modifier_group_id, restaurant_id, name, description, price_delta, available, rank, date_created, date_updated
	FROM
		modifier_options
	WHERE 
		modifier_option_id = :modifier_option_id`

	var dbOption modifierOption
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, q, data, &dbOption); err != nil {
		if errors.Is(err, sqldb.ErrDBNotFound) {
			return modifieroptionbus.ModifierOption{}, fmt.Errorf("db: %w", modifieroptionbus.ErrNotFound)
		}
		return modifieroptionbus.ModifierOption{}, fmt.Errorf("db: %w", err)
	}

	return toBusModifierOption(dbOption)
}
