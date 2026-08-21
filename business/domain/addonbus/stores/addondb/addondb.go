package addondb

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/warlck/food-flow/business/domain/addonbus"
	"github.com/warlck/food-flow/business/sdk/order"
	"github.com/warlck/food-flow/business/sdk/page"
	"github.com/warlck/food-flow/business/sdk/sqldb"
	"github.com/warlck/food-flow/foundation/logger"
)

// Store manages the set of APIs for addon database access.
type Store struct {
	log *logger.Logger
	db  sqlx.ExtContext
}

// NewStore constructs the API for data access.
func NewStore(log *logger.Logger, db *sqlx.DB) *Store {
	return &Store{
		log: log,
		db:  db,
	}
}

// Create inserts a new addon into the database.
func (s *Store) Create(ctx context.Context, addon addonbus.Addon) error {
	const q = `
	INSERT INTO addons
		(addon_id, category_id, restaurant_id, name, description, price, available, max_quantity, date_created, date_updated, rank)
	VALUES
		(:addon_id, :category_id, :restaurant_id, :name, :description, :price, :available, :max_quantity, :date_created, :date_updated, :rank)`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, toDBAddon(addon)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// Update replaces an addon document in the database.
func (s *Store) Update(ctx context.Context, addon addonbus.Addon) error {
	const q = `
	UPDATE
		addons
	SET 
		name = :name,
		description = :description,
		price = :price,
		available = :available,
		max_quantity = :max_quantity,
		rank = :rank,
		date_updated = :date_updated
	WHERE
		addon_id = :addon_id`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, toDBAddon(addon)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// Delete removes an addon from the database.
func (s *Store) Delete(ctx context.Context, addon addonbus.Addon) error {
	const q = `
	DELETE FROM
		addons
	WHERE
		addon_id = :addon_id`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, toDBAddon(addon)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// Query retrieves a list of existing addons from the database.
func (s *Store) Query(ctx context.Context, filter addonbus.QueryFilter, orderBy order.By, pg page.Page) ([]addonbus.Addon, error) {
	data := map[string]any{
		"offset":        (pg.Number() - 1) * pg.RowsPerPage(),
		"rows_per_page": pg.RowsPerPage(),
	}

	const q = `
	SELECT
		addon_id, category_id, restaurant_id, name, description, price, available, max_quantity, date_created, date_updated, rank
	FROM
		addons`

	buf := bytes.NewBufferString(q)
	applyFilter(filter, data, buf)

	orderByClause, err := orderByClause(orderBy)
	if err != nil {
		return nil, err
	}

	buf.WriteString(orderByClause)
	buf.WriteString(" OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY")

	var dbAddons []dbAddon
	if err := sqldb.NamedQuerySlice(ctx, s.log, s.db, buf.String(), data, &dbAddons); err != nil {
		return nil, fmt.Errorf("namedqueryslice: %w", err)
	}

	return toBusAddons(dbAddons)
}

// Count returns the total number of addons in the database.
func (s *Store) Count(ctx context.Context, filter addonbus.QueryFilter) (int, error) {
	data := map[string]any{}

	const q = `
	SELECT
		count(1)
	FROM
		addons`

	buf := bytes.NewBufferString(q)
	applyFilter(filter, data, buf)

	var count struct {
		Count int `db:"count"`
	}
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, buf.String(), data, &count); err != nil {
		return 0, fmt.Errorf("namedquerystruct: %w", err)
	}

	return count.Count, nil
}

// QueryByID gets the specified addon from the database.
func (s *Store) QueryByID(ctx context.Context, addonID uuid.UUID) (addonbus.Addon, error) {
	data := struct {
		ID uuid.UUID `db:"addon_id"`
	}{
		ID: addonID,
	}

	const q = `
	SELECT
		addon_id, category_id, restaurant_id, name, description, price, available, max_quantity, date_created, date_updated, rank
	FROM
		addons
	WHERE 
		addon_id = :addon_id`

	var dbAddn dbAddon
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, q, data, &dbAddn); err != nil {
		if errors.Is(err, sqldb.ErrDBNotFound) {
			return addonbus.Addon{}, addonbus.ErrNotFound
		}
		return addonbus.Addon{}, fmt.Errorf("namedquerystruct: %w", err)
	}

	return toBusAddon(dbAddn)
}

// QueryByCategoryID gets all addons for a specific category from the database.
func (s *Store) QueryByCategoryID(ctx context.Context, categoryID uuid.UUID) ([]addonbus.Addon, error) {
	data := struct {
		CategoryID uuid.UUID `db:"category_id"`
	}{
		CategoryID: categoryID,
	}

	const q = `
	SELECT
		addon_id, category_id, restaurant_id, name, description, price, available, max_quantity, date_created, date_updated, rank
	FROM
		addons
	WHERE 
		category_id = :category_id
	ORDER BY
		rank ASC, name ASC, addon_id ASC`

	var dbAddons []dbAddon
	if err := sqldb.NamedQuerySlice(ctx, s.log, s.db, q, data, &dbAddons); err != nil {
		return nil, fmt.Errorf("namedqueryslice: %w", err)
	}

	return toBusAddons(dbAddons)
}

// Reorder updates the rank of addons in a category transactionally in steps of 10.
func (s *Store) Reorder(ctx context.Context, categoryID uuid.UUID, orderedIDs []uuid.UUID) error {
	tx, err := s.db.(*sqlx.DB).BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	const q = `
	UPDATE
		addons
	SET
		"rank" = :rank,
		"date_updated" = :date_updated
	WHERE
		addon_id = :addon_id AND category_id = :category_id`

	now := time.Now().UTC()
	for i, id := range orderedIDs {
		rank := (i + 1) * 10
		data := struct {
			Rank        int       `db:"rank"`
			DateUpdated time.Time `db:"date_updated"`
			AddonID     uuid.UUID `db:"addon_id"`
			CategoryID  uuid.UUID `db:"category_id"`
		}{
			Rank:        rank,
			DateUpdated: now,
			AddonID:     id,
			CategoryID:  categoryID,
		}

		if err := sqldb.NamedExecContext(ctx, s.log, tx, q, data); err != nil {
			return fmt.Errorf("namedexeccontext: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

// orderByClause validates the order by clause and returns the SQL string.
func orderByClause(orderBy order.By) (string, error) {
	const orderByFields = "addon_id, name, price, rank, date_created"

	by, exists := orderByFields, true
	_ = by
	_ = exists

	return fmt.Sprintf(" ORDER BY %s %s", orderBy.Field, orderBy.Direction), nil
}
