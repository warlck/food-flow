// Package categorydb contains category related CRUD functionality.
package categorydb

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/warlck/food-flow/business/domain/categorybus"
	"github.com/warlck/food-flow/business/sdk/order"
	"github.com/warlck/food-flow/business/sdk/page"
	"github.com/warlck/food-flow/business/sdk/sqldb"
	"github.com/warlck/food-flow/foundation/logger"
)

// Store manages the set of APIs for category database access.
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

// Create inserts a new category into the database.
func (s *Store) Create(ctx context.Context, cat categorybus.Category) error {
	const q = `
	INSERT INTO categories
		(category_id, name, description, restaurant_id, enabled, date_created, date_updated)
	VALUES
		(:category_id, :name, :description, :restaurant_id, :enabled, :date_created, :date_updated)`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, toDBCategory(cat)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// Update replaces a category document in the database.
func (s *Store) Update(ctx context.Context, cat categorybus.Category) error {
	const q = `
	UPDATE
		categories
	SET 
		"name" = :name,
		"description" = :description,
		"enabled" = :enabled,
		"date_updated" = :date_updated
	WHERE
		category_id = :category_id`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, toDBCategory(cat)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// Delete removes a category from the database.
func (s *Store) Delete(ctx context.Context, cat categorybus.Category) error {
	data := struct {
		ID string `db:"category_id"`
	}{
		ID: cat.ID.String(),
	}

	const q = `
	DELETE FROM
		categories
	WHERE
		category_id = :category_id`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, data); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// Query retrieves a list of existing categories from the database.
func (s *Store) Query(ctx context.Context, filter categorybus.QueryFilter, orderBy order.By, page page.Page) ([]categorybus.Category, error) {
	data := map[string]any{
		"offset":        (page.Number() - 1) * page.RowsPerPage(),
		"rows_per_page": page.RowsPerPage(),
	}

	const q = `
	SELECT
		category_id, name, description, restaurant_id, enabled, date_created, date_updated
	FROM
		categories`

	buf := bytes.NewBufferString(q)
	applyFilter(filter, data, buf)

	orderByClause, err := orderByClause(orderBy)
	if err != nil {
		return nil, err
	}

	buf.WriteString(orderByClause)
	buf.WriteString(" OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY")

	var dbCats []category
	if err := sqldb.NamedQuerySlice(ctx, s.log, s.db, buf.String(), data, &dbCats); err != nil {
		return nil, fmt.Errorf("namedqueryslice: %w", err)
	}

	return toBusCategories(dbCats)
}

// QueryAll retrieves all categories matching the filter from the database without pagination.
func (s *Store) QueryAll(ctx context.Context, filter categorybus.QueryFilter, orderBy order.By) ([]categorybus.Category, error) {
	data := map[string]any{}

	const q = `
	SELECT
		category_id, name, description, restaurant_id, enabled, date_created, date_updated
	FROM
		categories`

	buf := bytes.NewBufferString(q)
	applyFilter(filter, data, buf)

	orderByClause, err := orderByClause(orderBy)
	if err != nil {
		return nil, err
	}

	buf.WriteString(orderByClause)

	var dbCats []category
	if err := sqldb.NamedQuerySlice(ctx, s.log, s.db, buf.String(), data, &dbCats); err != nil {
		return nil, fmt.Errorf("namedqueryslice: %w", err)
	}

	return toBusCategories(dbCats)
}

// Count returns the total number of categories in the DB.
func (s *Store) Count(ctx context.Context, filter categorybus.QueryFilter) (int, error) {
	data := map[string]any{}

	const q = `
	SELECT
		count(1)
	FROM
		categories`

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

// QueryByID gets the specified category from the database.
func (s *Store) QueryByID(ctx context.Context, categoryID uuid.UUID) (categorybus.Category, error) {
	data := struct {
		ID string `db:"category_id"`
	}{
		ID: categoryID.String(),
	}

	const q = `
	SELECT
		category_id, name, description, restaurant_id, enabled, date_created, date_updated
	FROM
		categories
	WHERE 
		category_id = :category_id`

	var dbCat category
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, q, data, &dbCat); err != nil {
		if errors.Is(err, sqldb.ErrDBNotFound) {
			return categorybus.Category{}, fmt.Errorf("db: %w", categorybus.ErrNotFound)
		}
		return categorybus.Category{}, fmt.Errorf("db: %w", err)
	}

	return toBusCategory(dbCat)
}
