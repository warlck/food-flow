package promodb

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/warlck/food-flow/business/domain/promobus"
	"github.com/warlck/food-flow/business/sdk/order"
	"github.com/warlck/food-flow/business/sdk/page"
	"github.com/warlck/food-flow/business/sdk/sqldb"
	"github.com/warlck/food-flow/foundation/logger"
)

type Store struct {
	log *logger.Logger
	db  sqlx.ExtContext
}

func NewStore(log *logger.Logger, db *sqlx.DB) *Store {
	return &Store{
		log: log,
		db:  db,
	}
}

func (s *Store) Create(ctx context.Context, promo promobus.Promotion) error {
	const q = `
	INSERT INTO promotions
		(promotion_id, restaurant_id, code, name, description, discount_type, discount_value, min_order_amount, max_discount_amount, usage_limit, usage_count, start_date, end_date, enabled, date_created, date_updated)
	VALUES
		(:promotion_id, :restaurant_id, :code, :name, :description, :discount_type, :discount_value, :min_order_amount, :max_discount_amount, :usage_limit, :usage_count, :start_date, :end_date, :enabled, :date_created, :date_updated)`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, toDBPromotion(promo)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

func (s *Store) Update(ctx context.Context, promo promobus.Promotion) error {
	const q = `
	UPDATE
		promotions
	SET
		"restaurant_id"       = :restaurant_id,
		"code"                = :code,
		"name"                = :name,
		"description"         = :description,
		"discount_type"       = :discount_type,
		"discount_value"      = :discount_value,
		"min_order_amount"    = :min_order_amount,
		"max_discount_amount" = :max_discount_amount,
		"usage_limit"         = :usage_limit,
		"usage_count"         = :usage_count,
		"start_date"          = :start_date,
		"end_date"            = :end_date,
		"enabled"             = :enabled,
		"date_updated"        = :date_updated
	WHERE
		promotion_id = :promotion_id`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, toDBPromotion(promo)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

func (s *Store) Delete(ctx context.Context, promo promobus.Promotion) error {
	data := struct {
		ID string `db:"promotion_id"`
	}{
		ID: promo.ID.String(),
	}

	const q = `
	DELETE FROM
		promotions
	WHERE
		promotion_id = :promotion_id`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, data); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

func (s *Store) Query(ctx context.Context, filter promobus.QueryFilter, orderBy order.By, page page.Page) ([]promobus.Promotion, error) {
	data := map[string]any{
		"offset":        (page.Number() - 1) * page.RowsPerPage(),
		"rows_per_page": page.RowsPerPage(),
	}

	const q = `
	SELECT
		promotion_id, restaurant_id, code, name, description, discount_type, discount_value, min_order_amount, max_discount_amount, usage_limit, usage_count, start_date, end_date, enabled, date_created, date_updated
	FROM
		promotions`

	buf := bytes.NewBufferString(q)
	s.applyFilter(filter, data, buf)

	orderByClause, err := orderByClause(orderBy)
	if err != nil {
		return nil, err
	}
	buf.WriteString(" " + orderByClause)
	buf.WriteString(" OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY")

	var dbPromos []promotion
	if err := sqldb.NamedQuerySlice(ctx, s.log, s.db, buf.String(), data, &dbPromos); err != nil {
		return nil, fmt.Errorf("namedqueryslice: %w", err)
	}

	return toBusPromotions(dbPromos)
}

func (s *Store) Count(ctx context.Context, filter promobus.QueryFilter) (int, error) {
	data := map[string]any{}

	const q = `
	SELECT
		count(1)
	FROM
		promotions`

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

func (s *Store) QueryByID(ctx context.Context, promotionID uuid.UUID) (promobus.Promotion, error) {
	data := struct {
		ID string `db:"promotion_id"`
	}{
		ID: promotionID.String(),
	}

	const q = `
	SELECT
		promotion_id, restaurant_id, code, name, description, discount_type, discount_value, min_order_amount, max_discount_amount, usage_limit, usage_count, start_date, end_date, enabled, date_created, date_updated
	FROM
		promotions
	WHERE
		promotion_id = :promotion_id`

	var dbPromo promotion
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, q, data, &dbPromo); err != nil {
		if errors.Is(err, sqldb.ErrDBNotFound) {
			return promobus.Promotion{}, promobus.ErrNotFound
		}
		return promobus.Promotion{}, fmt.Errorf("namedquerystruct: %w", err)
	}

	return toBusPromotion(dbPromo)
}

func (s *Store) QueryByCode(ctx context.Context, code string) (promobus.Promotion, error) {
	data := struct {
		Code string `db:"code"`
	}{
		Code: code,
	}

	const q = `
	SELECT
		promotion_id, restaurant_id, code, name, description, discount_type, discount_value, min_order_amount, max_discount_amount, usage_limit, usage_count, start_date, end_date, enabled, date_created, date_updated
	FROM
		promotions
	WHERE
		UPPER(code) = UPPER(:code)`

	var dbPromo promotion
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, q, data, &dbPromo); err != nil {
		if errors.Is(err, sqldb.ErrDBNotFound) {
			return promobus.Promotion{}, promobus.ErrNotFound
		}
		return promobus.Promotion{}, fmt.Errorf("namedquerystruct: %w", err)
	}

	return toBusPromotion(dbPromo)
}

func (s *Store) IncrementUsage(ctx context.Context, promotionID uuid.UUID) error {
	data := struct {
		ID string `db:"promotion_id"`
	}{
		ID: promotionID.String(),
	}

	const q = `
	UPDATE
		promotions
	SET
		usage_count = usage_count + 1
	WHERE
		promotion_id = :promotion_id
		AND (usage_limit IS NULL OR usage_count < usage_limit)`

	res, err := sqldb.NamedExecResultContext(ctx, s.log, s.db, q, data)
	if err != nil {
		return fmt.Errorf("namedexecresultcontext: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rowsaffected: %w", err)
	}

	if rows == 0 {
		return promobus.ErrLimitReached
	}

	return nil
}

func (s *Store) applyFilter(filter promobus.QueryFilter, data map[string]any, buf *bytes.Buffer) {
	var wc []string

	if filter.ID != nil {
		data["promotion_id"] = *filter.ID
		wc = append(wc, "promotion_id = :promotion_id")
	}

	if filter.Code != nil {
		data["code"] = *filter.Code
		wc = append(wc, "UPPER(code) = UPPER(:code)")
	}

	if filter.Name != nil {
		data["name"] = fmt.Sprintf("%%%s%%", filter.Name.String())
		wc = append(wc, "name ILIKE :name")
	}

	if filter.RestaurantID != nil {
		data["restaurant_id"] = *filter.RestaurantID
		wc = append(wc, "(restaurant_id = :restaurant_id OR restaurant_id IS NULL)")
	}

	if filter.Enabled != nil {
		data["enabled"] = *filter.Enabled
		wc = append(wc, "enabled = :enabled")
	}

	if len(wc) > 0 {
		buf.WriteString(" WHERE ")
		for i, w := range wc {
			if i > 0 {
				buf.WriteString(" AND ")
			}
			buf.WriteString(w)
		}
	}
}

var orderByFields = map[string]string{
	promobus.OrderByID:          "promotion_id",
	promobus.OrderByCode:        "code",
	promobus.OrderByName:        "name",
	promobus.OrderByDateCreated: "date_created",
}

func orderByClause(orderBy order.By) (string, error) {
	by, ok := orderByFields[orderBy.Field]
	if !ok {
		return "", fmt.Errorf("field %q does not exist", orderBy.Field)
	}
	return "ORDER BY " + by + " " + orderBy.Direction, nil
}
