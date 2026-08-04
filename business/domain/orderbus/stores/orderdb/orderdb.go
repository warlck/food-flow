// Package orderdb contains order related CRUD functionality.
package orderdb

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/warlck/food-flow/business/domain/orderbus"
	"github.com/warlck/food-flow/business/sdk/order"
	"github.com/warlck/food-flow/business/sdk/page"
	"github.com/warlck/food-flow/business/sdk/sqldb"
	"github.com/warlck/food-flow/foundation/logger"
)

// Store manages the set of APIs for order database access.
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

// Create inserts a new order into the database.
func (s *Store) Create(ctx context.Context, order orderbus.Order) error {
	// Use transaction to insert order, items, and address
	tx, err := s.db.(*sqlx.DB).BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Insert order
	const orderSQL = `
	INSERT INTO orders (
		order_id, restaurant_id, customer_name, customer_email, customer_phone,
		order_type, order_status, payment_status, payment_method,
		subtotal, delivery_fee, tax, total, special_instructions,
		stripe_payment_intent_id, date_created, date_updated
	) VALUES (
		:order_id, :restaurant_id, :customer_name, :customer_email, :customer_phone,
		:order_type, :order_status, :payment_status, :payment_method,
		:subtotal, :delivery_fee, :tax, :total, :special_instructions,
		:stripe_payment_intent_id, :date_created, :date_updated
	)`

	if err := sqldb.NamedExecContext(ctx, s.log, tx, orderSQL, toDBOrder(order)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	// Insert order items
	for _, item := range order.Items {
		const itemSQL = `
		INSERT INTO order_items (
			order_item_id, order_id, menu_item_id, menu_item_name,
			menu_item_price, quantity, special_instructions, date_created
		) VALUES (
			:order_item_id, :order_id, :menu_item_id, :menu_item_name,
			:menu_item_price, :quantity, :special_instructions, :date_created
		)`

		if err := sqldb.NamedExecContext(ctx, s.log, tx, itemSQL, toDBOrderItem(item, order.ID)); err != nil {
			return fmt.Errorf("namedexeccontext: %w", err)
		}
	}

	// Insert delivery address if present
	if order.DeliveryAddress != nil {
		const addressSQL = `
		INSERT INTO delivery_addresses (
			address_id, order_id, street, city, state, postal_code,
			delivery_instructions, latitude, longitude, date_created
		) VALUES (
			:address_id, :order_id, :street, :city, :state, :postal_code,
			:delivery_instructions, :latitude, :longitude, :date_created
		)`

		if err := sqldb.NamedExecContext(ctx, s.log, tx, addressSQL, toDBDeliveryAddress(*order.DeliveryAddress, order.ID)); err != nil {
			return fmt.Errorf("namedexeccontext: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit: %w", err)
	}

	return nil
}

// Update replaces an order document in the database.
func (s *Store) Update(ctx context.Context, order orderbus.Order) error {
	const q = `
	UPDATE
		orders
	SET
		order_status = :order_status,
		payment_status = :payment_status,
		stripe_payment_intent_id = :stripe_payment_intent_id,
		date_updated = :date_updated
	WHERE
		order_id = :order_id`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, toDBOrder(order)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// Delete removes an order from the database (soft delete by setting status to cancelled).
func (s *Store) Delete(ctx context.Context, orderID uuid.UUID) error {
	data := struct {
		ID uuid.UUID `db:"order_id"`
	}{
		ID: orderID,
	}

	const q = `
	UPDATE
		orders
	SET
		order_status = 'cancelled',
		date_updated = NOW()
	WHERE
		order_id = :order_id`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, data); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// Query retrieves a list of existing orders from the database.
func (s *Store) Query(ctx context.Context, filter orderbus.QueryFilter, orderBy order.By, pageNumber page.Page) ([]orderbus.Order, error) {
	data := map[string]any{
		"offset":        (pageNumber.Number() - 1) * pageNumber.RowsPerPage(),
		"rows_per_page": pageNumber.RowsPerPage(),
	}

	const q = `
	SELECT
		o.order_id, o.restaurant_id, o.customer_name, o.customer_email, o.customer_phone,
		o.order_type, o.order_status, o.payment_status, o.payment_method,
		o.subtotal, o.delivery_fee, o.tax, o.total, o.special_instructions,
		o.stripe_payment_intent_id, o.date_created, o.date_updated
	FROM
		orders AS o`

	buf := bytes.NewBufferString(q)
	s.applyFilter(filter, data, buf)

	orderByClause, err := orderByClause(orderBy)
	if err != nil {
		return nil, err
	}

	buf.WriteString(orderByClause)
	buf.WriteString(" OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY")

	var dbOrders []dbOrder
	if err := sqldb.NamedQuerySlice(ctx, s.log, s.db, buf.String(), data, &dbOrders); err != nil {
		return nil, fmt.Errorf("namedqueryslice: %w", err)
	}

	// Load order items and delivery addresses for each order
	orders := make([]orderbus.Order, len(dbOrders))
	for i, dbo := range dbOrders {
		// Load items
		items, err := s.queryOrderItems(ctx, dbo.ID)
		if err != nil {
			return nil, err
		}

		// Load delivery address if delivery order
		var addr *dbDeliveryAddress
		if dbo.OrderType == orderbus.OrderTypeDelivery {
			a, err := s.queryDeliveryAddress(ctx, dbo.ID)
			if err != nil && !errors.Is(err, sqldb.ErrDBNotFound) {
				return nil, err
			}
			if err == nil {
				addr = &a
			}
		}

		order, err := toBusOrder(dbo, items, addr)
		if err != nil {
			return nil, fmt.Errorf("tobusorder: %w", err)
		}
		orders[i] = order
	}

	return orders, nil
}

// Count returns the total number of orders in the DB.
func (s *Store) Count(ctx context.Context, filter orderbus.QueryFilter) (int, error) {
	data := map[string]any{}

	const q = `
	SELECT
		COUNT(1) AS count
	FROM
		orders AS o`

	buf := bytes.NewBufferString(q)
	s.applyFilter(filter, data, buf)

	var count struct {
		Count int `db:"count"`
	}
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, buf.String(), data, &count); err != nil {
		return 0, fmt.Errorf("namedquerystruct: %w", err)
	}

	return count.Count, nil
}

// QueryByID gets the specified order from the database.
func (s *Store) QueryByID(ctx context.Context, orderID uuid.UUID) (orderbus.Order, error) {
	data := struct {
		ID uuid.UUID `db:"order_id"`
	}{
		ID: orderID,
	}

	const q = `
	SELECT
		o.order_id, o.restaurant_id, o.customer_name, o.customer_email, o.customer_phone,
		o.order_type, o.order_status, o.payment_status, o.payment_method,
		o.subtotal, o.delivery_fee, o.tax, o.total, o.special_instructions,
		o.stripe_payment_intent_id, o.date_created, o.date_updated
	FROM
		orders AS o
	WHERE
		o.order_id = :order_id`

	var dbo dbOrder
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, q, data, &dbo); err != nil {
		return orderbus.Order{}, fmt.Errorf("namedquerystruct: %w", err)
	}

	// Load items
	items, err := s.queryOrderItems(ctx, orderID)
	if err != nil {
		return orderbus.Order{}, err
	}

	// Load delivery address if delivery order
	var addr *dbDeliveryAddress
	if dbo.OrderType == orderbus.OrderTypeDelivery {
		a, err := s.queryDeliveryAddress(ctx, orderID)
		if err != nil && !errors.Is(err, sqldb.ErrDBNotFound) {
			return orderbus.Order{}, err
		}
		if err == nil {
			addr = &a
		}
	}

	order, err := toBusOrder(dbo, items, addr)
	if err != nil {
		return orderbus.Order{}, fmt.Errorf("tobusorder: %w", err)
	}

	return order, nil
}

// queryOrderItems retrieves all items for an order.
func (s *Store) queryOrderItems(ctx context.Context, orderID uuid.UUID) ([]dbOrderItem, error) {
	data := struct {
		OrderID uuid.UUID `db:"order_id"`
	}{
		OrderID: orderID,
	}

	const q = `
	SELECT
		order_item_id, order_id, menu_item_id, menu_item_name,
		menu_item_price, quantity, special_instructions, date_created
	FROM
		order_items
	WHERE
		order_id = :order_id
	ORDER BY
		date_created`

	var items []dbOrderItem
	if err := sqldb.NamedQuerySlice(ctx, s.log, s.db, q, data, &items); err != nil {
		return nil, fmt.Errorf("namedqueryslice: %w", err)
	}

	return items, nil
}

// queryDeliveryAddress retrieves the delivery address for an order.
func (s *Store) queryDeliveryAddress(ctx context.Context, orderID uuid.UUID) (dbDeliveryAddress, error) {
	data := struct {
		OrderID uuid.UUID `db:"order_id"`
	}{
		OrderID: orderID,
	}

	const q = `
	SELECT
		address_id, order_id, street, city, state, postal_code,
		delivery_instructions, latitude, longitude, date_created
	FROM
		delivery_addresses
	WHERE
		order_id = :order_id`

	var addr dbDeliveryAddress
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, q, data, &addr); err != nil {
		return dbDeliveryAddress{}, fmt.Errorf("namedquerystruct: %w", err)
	}

	return addr, nil
}

// applyFilter applies the filter to the query.
func (s *Store) applyFilter(filter orderbus.QueryFilter, data map[string]any, buf *bytes.Buffer) {
	var wc []string

	if filter.ID != nil {
		data["order_id"] = *filter.ID
		wc = append(wc, "o.order_id = :order_id")
	}

	if filter.RestaurantID != nil {
		data["restaurant_id"] = *filter.RestaurantID
		wc = append(wc, "o.restaurant_id = :restaurant_id")
	}

	if filter.CustomerEmail != nil {
		data["customer_email"] = *filter.CustomerEmail
		wc = append(wc, "o.customer_email = :customer_email")
	}

	if filter.OrderStatus != nil {
		data["order_status"] = *filter.OrderStatus
		wc = append(wc, "o.order_status = :order_status")
	}

	if filter.PaymentStatus != nil {
		data["payment_status"] = *filter.PaymentStatus
		wc = append(wc, "o.payment_status = :payment_status")
	}

	if filter.OrderType != nil {
		data["order_type"] = *filter.OrderType
		wc = append(wc, "o.order_type = :order_type")
	}

	if filter.StartDate != nil {
		data["start_date"] = *filter.StartDate
		wc = append(wc, "o.date_created >= :start_date")
	}

	if filter.EndDate != nil {
		data["end_date"] = *filter.EndDate
		wc = append(wc, "o.date_created <= :end_date")
	}

	if len(wc) > 0 {
		buf.WriteString(" WHERE ")
		buf.WriteString(wc[0])
		for i := 1; i < len(wc); i++ {
			buf.WriteString(" AND ")
			buf.WriteString(wc[i])
		}
	}
}
