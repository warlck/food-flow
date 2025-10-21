// Package restaurantdb contains restaurant related CRUD functionality.
package restaurantdb

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/warlck/food-flow/business/domain/restaurantbus"
	"github.com/warlck/food-flow/business/sdk/order"
	"github.com/warlck/food-flow/business/sdk/page"
	"github.com/warlck/food-flow/business/sdk/sqldb"
	"github.com/warlck/food-flow/foundation/logger"
)

// Store manages the set of APIs for restaurant database access.
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

// Create inserts a new restaurant into the database.
func (s *Store) Create(ctx context.Context, rest restaurantbus.Restaurant) error {
	const q = `
	INSERT INTO restaurants
		(restaurant_id, name, description, address, phone, email, image_url, enabled, date_created, date_updated)
	VALUES
		(:restaurant_id, :name, :description, :address, :phone, :email, :image_url, :enabled, :date_created, :date_updated)`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, toDBRestaurant(rest)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// Update replaces a restaurant document in the database.
func (s *Store) Update(ctx context.Context, rest restaurantbus.Restaurant) error {
	const q = `
	UPDATE
		restaurants
	SET 
		"name" = :name,
		"description" = :description,
		"address" = :address,
		"phone" = :phone,
		"email" = :email,
		"image_url" = :image_url,
		"enabled" = :enabled,
		"date_updated" = :date_updated
	WHERE
		restaurant_id = :restaurant_id`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, toDBRestaurant(rest)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// Delete removes a restaurant from the database.
func (s *Store) Delete(ctx context.Context, rest restaurantbus.Restaurant) error {
	data := struct {
		ID string `db:"restaurant_id"`
	}{
		ID: rest.ID.String(),
	}

	const q = `
	DELETE FROM
		restaurants
	WHERE
		restaurant_id = :restaurant_id`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, data); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// Query retrieves a list of existing restaurants from the database.
func (s *Store) Query(ctx context.Context, filter restaurantbus.QueryFilter, orderBy order.By, page page.Page) ([]restaurantbus.Restaurant, error) {
	data := map[string]any{
		"offset":        (page.Number() - 1) * page.RowsPerPage(),
		"rows_per_page": page.RowsPerPage(),
	}

	const q = `
	SELECT
		restaurant_id, name, description, address, phone, email, image_url, enabled, date_created, date_updated
	FROM
		restaurants`

	buf := bytes.NewBufferString(q)
	applyFilter(filter, data, buf)

	orderByClause, err := orderByClause(orderBy)
	if err != nil {
		return nil, err
	}

	buf.WriteString(orderByClause)
	buf.WriteString(" OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY")

	var dbRests []restaurant
	if err := sqldb.NamedQuerySlice(ctx, s.log, s.db, buf.String(), data, &dbRests); err != nil {
		return nil, fmt.Errorf("namedqueryslice: %w", err)
	}

	return toBusRestaurants(dbRests)
}

// Count returns the total number of restaurants in the DB.
func (s *Store) Count(ctx context.Context, filter restaurantbus.QueryFilter) (int, error) {
	data := map[string]any{}

	const q = `
	SELECT
		count(1)
	FROM
		restaurants`

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

// QueryByID gets the specified restaurant from the database.
func (s *Store) QueryByID(ctx context.Context, restaurantID uuid.UUID) (restaurantbus.Restaurant, error) {
	data := struct {
		ID string `db:"restaurant_id"`
	}{
		ID: restaurantID.String(),
	}

	const q = `
	SELECT
		restaurant_id, name, description, address, phone, email, image_url, enabled, date_created, date_updated
	FROM
		restaurants
	WHERE 
		restaurant_id = :restaurant_id`

	var dbRest restaurant
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, q, data, &dbRest); err != nil {
		if errors.Is(err, sqldb.ErrDBNotFound) {
			return restaurantbus.Restaurant{}, fmt.Errorf("db: %w", restaurantbus.ErrNotFound)
		}
		return restaurantbus.Restaurant{}, fmt.Errorf("db: %w", err)
	}

	return toBusRestaurant(dbRest)
}
