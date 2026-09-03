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
		(addon_id, menu_item_id, restaurant_id, name, description, price, available, max_quantity, rank, date_created, date_updated)
	VALUES
		(:addon_id, :menu_item_id, :restaurant_id, :name, :description, :price, :available, :max_quantity, :rank, :date_created, :date_updated)`

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
		addon_id, menu_item_id, restaurant_id, name, description, price, available, max_quantity, rank, date_created, date_updated
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

// QueryAll retrieves all addons matching the filter without pagination.
func (s *Store) QueryAll(ctx context.Context, filter addonbus.QueryFilter, orderBy order.By) ([]addonbus.Addon, error) {
	data := map[string]any{}

	const q = `
	SELECT
		addon_id, menu_item_id, restaurant_id, name, description, price, available, max_quantity, rank, date_created, date_updated
	FROM
		addons`

	buf := bytes.NewBufferString(q)
	applyFilter(filter, data, buf)

	orderByClause, err := orderByClause(orderBy)
	if err != nil {
		return nil, err
	}

	buf.WriteString(orderByClause)

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
		addon_id, menu_item_id, restaurant_id, name, description, price, available, max_quantity, rank, date_created, date_updated
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

// Reorder updates the ranks of addons for a menu item in a transaction.
func (s *Store) Reorder(ctx context.Context, menuItemID uuid.UUID, addonIDs []uuid.UUID) error {
	tx, err := s.db.(*sqlx.DB).BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	now := time.Now().UTC()

	const q = `
	UPDATE
		addons
	SET
		rank = :rank,
		date_updated = :date_updated
	WHERE
		addon_id = :addon_id AND menu_item_id = :menu_item_id`

	for i, id := range addonIDs {
		rankVal := (i + 1) * 10
		data := struct {
			AddonID     uuid.UUID `db:"addon_id"`
			MenuItemID  uuid.UUID `db:"menu_item_id"`
			Rank        int       `db:"rank"`
			DateUpdated time.Time `db:"date_updated"`
		}{
			AddonID:     id,
			MenuItemID:  menuItemID,
			Rank:        rankVal,
			DateUpdated: now,
		}

		if err := sqldb.NamedExecContext(ctx, s.log, tx, q, data); err != nil {
			return fmt.Errorf("update addon rank: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}
