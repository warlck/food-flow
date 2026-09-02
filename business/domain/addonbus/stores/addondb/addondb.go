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
		(addon_id, restaurant_id, name, description, price, available, max_quantity, date_created, date_updated)
	VALUES
		(:addon_id, :restaurant_id, :name, :description, :price, :available, :max_quantity, :date_created, :date_updated)`

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
		addon_id, restaurant_id, name, description, price, available, max_quantity, date_created, date_updated
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
		addon_id, restaurant_id, name, description, price, available, max_quantity, date_created, date_updated
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
		addon_id, restaurant_id, name, description, price, available, max_quantity, date_created, date_updated
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

// QueryMenuItemAddons gets all assigned addons for a menu item.
func (s *Store) QueryMenuItemAddons(ctx context.Context, menuItemID uuid.UUID) ([]addonbus.MenuItemAddonInfo, error) {
	data := struct {
		MenuItemID uuid.UUID `db:"menu_item_id"`
	}{
		MenuItemID: menuItemID,
	}

	const q = `
	SELECT
		a.addon_id,
		a.restaurant_id,
		a.name,
		a.description,
		a.price,
		a.available,
		a.max_quantity,
		a.date_created,
		a.date_updated,
		mia.rank AS association_rank
	FROM
		menu_item_addons AS mia
	JOIN
		addons AS a ON a.addon_id = mia.addon_id
	WHERE
		mia.menu_item_id = :menu_item_id
	ORDER BY
		mia.rank ASC NULLS LAST, a.name ASC, a.addon_id ASC`

	var rows []dbMenuItemAddonRow
	if err := sqldb.NamedQuerySlice(ctx, s.log, s.db, q, data, &rows); err != nil {
		return nil, fmt.Errorf("namedqueryslice: %w", err)
	}

	return toBusMenuItemAddons(rows)
}

// ReplaceMenuItemAddons replaces all addon associations for a menu item in a transaction.
func (s *Store) ReplaceMenuItemAddons(ctx context.Context, menuItemID uuid.UUID, restaurantID uuid.UUID, assignments []addonbus.ItemAddonAssignment) error {
	tx, err := s.db.(*sqlx.DB).BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	const delQ = `DELETE FROM menu_item_addons WHERE menu_item_id = :menu_item_id`
	if err := sqldb.NamedExecContext(ctx, s.log, tx, delQ, map[string]any{"menu_item_id": menuItemID}); err != nil {
		return fmt.Errorf("delete menu_item_addons: %w", err)
	}

	const insQ = `
	INSERT INTO menu_item_addons
		(menu_item_id, addon_id, restaurant_id, rank, date_created)
	VALUES
		(:menu_item_id, :addon_id, :restaurant_id, :rank, :date_created)`

	now := time.Now().UTC()
	for _, a := range assignments {
		row := struct {
			MenuItemID   uuid.UUID `db:"menu_item_id"`
			AddonID      uuid.UUID `db:"addon_id"`
			RestaurantID uuid.UUID `db:"restaurant_id"`
			Rank         *int      `db:"rank"`
			DateCreated  time.Time `db:"date_created"`
		}{
			MenuItemID:   menuItemID,
			AddonID:      a.AddonID,
			RestaurantID: restaurantID,
			Rank:         a.Rank,
			DateCreated:  now,
		}

		if err := sqldb.NamedExecContext(ctx, s.log, tx, insQ, row); err != nil {
			return fmt.Errorf("insert menu_item_addons: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

// ReorderMenuItemAddons updates the ranks of assigned addons in a transaction.
func (s *Store) ReorderMenuItemAddons(ctx context.Context, menuItemID uuid.UUID, assignments []addonbus.ItemAddonAssignment) error {
	tx, err := s.db.(*sqlx.DB).BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	const q = `
	UPDATE
		menu_item_addons
	SET
		rank = :rank
	WHERE
		menu_item_id = :menu_item_id AND addon_id = :addon_id`

	for _, a := range assignments {
		data := struct {
			MenuItemID uuid.UUID `db:"menu_item_id"`
			AddonID    uuid.UUID `db:"addon_id"`
			Rank       *int      `db:"rank"`
		}{
			MenuItemID: menuItemID,
			AddonID:    a.AddonID,
			Rank:       a.Rank,
		}

		if err := sqldb.NamedExecContext(ctx, s.log, tx, q, data); err != nil {
			return fmt.Errorf("update menu_item_addons: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}
