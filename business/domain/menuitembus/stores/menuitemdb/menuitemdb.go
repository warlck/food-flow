// Package menuitemdb contains menu item related CRUD functionality.
package menuitemdb

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/warlck/food-flow/business/domain/menuitembus"
	"github.com/warlck/food-flow/business/sdk/order"
	"github.com/warlck/food-flow/business/sdk/page"
	"github.com/warlck/food-flow/business/sdk/sqldb"
	"github.com/warlck/food-flow/foundation/logger"
)

// Store manages the set of APIs for menu item database access.
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

// Create inserts a new menu item into the database.
func (s *Store) Create(ctx context.Context, item menuitembus.MenuItem) error {
	const q = `
	INSERT INTO menu_items
		(menu_item_id, name, description, price, category_id, restaurant_id, image_url, available, date_created, date_updated, rank)
	VALUES
		(:menu_item_id, :name, :description, :price, :category_id, :restaurant_id, :image_url, :available, :date_created, :date_updated, :rank)`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, toDBMenuItem(item)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// Update replaces a menu item document in the database.
func (s *Store) Update(ctx context.Context, item menuitembus.MenuItem) error {
	const q = `
	UPDATE
		menu_items
	SET 
		"name" = :name,
		"description" = :description,
		"price" = :price,
		"category_id" = :category_id,
		"image_url" = :image_url,
		"available" = :available,
		"rank" = :rank,
		"date_updated" = :date_updated
	WHERE
		menu_item_id = :menu_item_id`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, toDBMenuItem(item)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// Delete removes a menu item from the database.
func (s *Store) Delete(ctx context.Context, item menuitembus.MenuItem) error {
	data := struct {
		ID string `db:"menu_item_id"`
	}{
		ID: item.ID.String(),
	}

	const q = `
	DELETE FROM
		menu_items
	WHERE
		menu_item_id = :menu_item_id`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, data); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// Query retrieves a list of existing menu items from the database.
func (s *Store) Query(ctx context.Context, filter menuitembus.QueryFilter, orderBy order.By, page page.Page) ([]menuitembus.MenuItem, error) {
	data := map[string]any{
		"offset":        (page.Number() - 1) * page.RowsPerPage(),
		"rows_per_page": page.RowsPerPage(),
	}

	const q = `
	SELECT
		menu_item_id, name, description, price, category_id, restaurant_id, image_url, available, date_created, date_updated, rank
	FROM
		menu_items`

	buf := bytes.NewBufferString(q)
	applyFilter(filter, data, buf)

	orderByClause, err := orderByClause(orderBy)
	if err != nil {
		return nil, err
	}

	buf.WriteString(orderByClause)
	buf.WriteString(" OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY")

	var dbItems []menuItem
	if err := sqldb.NamedQuerySlice(ctx, s.log, s.db, buf.String(), data, &dbItems); err != nil {
		return nil, fmt.Errorf("namedqueryslice: %w", err)
	}

	return toBusMenuItems(dbItems)
}

// QueryAll retrieves all menu items matching the filter from the database without pagination.
func (s *Store) QueryAll(ctx context.Context, filter menuitembus.QueryFilter, orderBy order.By) ([]menuitembus.MenuItem, error) {
	data := map[string]any{}

	const q = `
	SELECT
		menu_item_id, name, description, price, category_id, restaurant_id, image_url, available, date_created, date_updated, rank
	FROM
		menu_items`

	buf := bytes.NewBufferString(q)
	applyFilter(filter, data, buf)

	orderByClause, err := orderByClause(orderBy)
	if err != nil {
		return nil, err
	}

	buf.WriteString(orderByClause)

	var dbItems []menuItem
	if err := sqldb.NamedQuerySlice(ctx, s.log, s.db, buf.String(), data, &dbItems); err != nil {
		return nil, fmt.Errorf("namedqueryslice: %w", err)
	}

	return toBusMenuItems(dbItems)
}

// Count returns the total number of menu items in the DB.
func (s *Store) Count(ctx context.Context, filter menuitembus.QueryFilter) (int, error) {
	data := map[string]any{}

	const q = `
	SELECT
		count(1)
	FROM
		menu_items`

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

// QueryByID gets the specified menu item from the database.
func (s *Store) QueryByID(ctx context.Context, menuItemID uuid.UUID) (menuitembus.MenuItem, error) {
	data := struct {
		ID string `db:"menu_item_id"`
	}{
		ID: menuItemID.String(),
	}

	const q = `
	SELECT
		menu_item_id, name, description, price, category_id, restaurant_id, image_url, available, date_created, date_updated, rank
	FROM
		menu_items
	WHERE 
		menu_item_id = :menu_item_id`

	var dbItem menuItem
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, q, data, &dbItem); err != nil {
		if errors.Is(err, sqldb.ErrDBNotFound) {
			return menuitembus.MenuItem{}, fmt.Errorf("db: %w", menuitembus.ErrNotFound)
		}
		return menuitembus.MenuItem{}, fmt.Errorf("db: %w", err)
	}

	return toBusMenuItem(dbItem)
}

// QueryByCategoryID gets all menu items for a specific category from the database.
func (s *Store) QueryByCategoryID(ctx context.Context, categoryID uuid.UUID) ([]menuitembus.MenuItem, error) {
	data := struct {
		CategoryID uuid.UUID `db:"category_id"`
	}{
		CategoryID: categoryID,
	}

	const q = `
	SELECT
		menu_item_id, name, description, price, category_id, restaurant_id, image_url, available, date_created, date_updated, rank
	FROM
		menu_items
	WHERE
		category_id = :category_id
	ORDER BY
		rank ASC NULLS LAST, price ASC, name ASC, menu_item_id ASC`

	var dbItems []menuItem
	if err := sqldb.NamedQuerySlice(ctx, s.log, s.db, q, data, &dbItems); err != nil {
		return nil, fmt.Errorf("namedqueryslice: %w", err)
	}

	return toBusMenuItems(dbItems)
}

// Reorder updates the rank of menu items in a category transactionally in steps of 10.
func (s *Store) Reorder(ctx context.Context, items []menuitembus.MenuItem) error {
	tx, err := s.db.(*sqlx.DB).BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	const q = `
	UPDATE
		menu_items
	SET
		"rank" = :rank,
		"date_updated" = :date_updated
	WHERE
		menu_item_id = :menu_item_id`

	for _, itm := range items {
		if err := sqldb.NamedExecContext(ctx, s.log, tx, q, toDBMenuItem(itm)); err != nil {
			return fmt.Errorf("namedexeccontext: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}
