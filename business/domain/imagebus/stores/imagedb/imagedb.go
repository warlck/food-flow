// Package imagedb provides database access for the images table.
package imagedb

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/warlck/food-flow/business/domain/imagebus"
	"github.com/warlck/food-flow/business/sdk/order"
	"github.com/warlck/food-flow/business/sdk/page"
	"github.com/warlck/food-flow/business/sdk/sqldb"
	"github.com/warlck/food-flow/foundation/logger"
)

// Store manages image persistence.
type Store struct {
	log *logger.Logger
	db  sqlx.ExtContext
}

// NewStore constructs an image Store.
func NewStore(log *logger.Logger, db *sqlx.DB) *Store {
	return &Store{
		log: log,
		db:  db,
	}
}

// Create inserts a new image record.
func (s *Store) Create(ctx context.Context, img imagebus.Image) error {
	const q = `
	INSERT INTO images
		(image_id, restaurant_id, entity_type, object_path, public_url, content_type, size_bytes, status, uploaded_by, date_created, date_updated)
	VALUES
		(:image_id, :restaurant_id, :entity_type, :object_path, :public_url, :content_type, :size_bytes, :status, :uploaded_by, :date_created, :date_updated)`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, toDBImage(img)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// Update mutates the mutable fields of an image record.
func (s *Store) Update(ctx context.Context, img imagebus.Image) error {
	const q = `
	UPDATE
		images
	SET
		"entity_type"  = :entity_type,
		"object_path"  = :object_path,
		"public_url"   = :public_url,
		"content_type" = :content_type,
		"size_bytes"   = :size_bytes,
		"status"       = :status,
		"uploaded_by"  = :uploaded_by,
		"date_updated" = :date_updated
	WHERE
		image_id = :image_id`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, toDBImage(img)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// Delete removes an image record.
func (s *Store) Delete(ctx context.Context, img imagebus.Image) error {
	data := struct {
		ID string `db:"image_id"`
	}{
		ID: img.ID.String(),
	}

	const q = `
	DELETE FROM
		images
	WHERE
		image_id = :image_id`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, data); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// Query retrieves a list of images based on the filter.
func (s *Store) Query(ctx context.Context, filter imagebus.QueryFilter, orderBy order.By, page page.Page) ([]imagebus.Image, error) {
	data := map[string]any{
		"offset":        (page.Number() - 1) * page.RowsPerPage(),
		"rows_per_page": page.RowsPerPage(),
	}

	const q = `
	SELECT
		image_id, restaurant_id, entity_type, object_path, public_url, content_type, size_bytes, status, uploaded_by, date_created, date_updated
	FROM
		images`

	buf := bytes.NewBufferString(q)
	s.applyFilter(filter, data, buf)

	orderByClause, err := orderByClause(orderBy)
	if err != nil {
		return nil, err
	}
	buf.WriteString(" " + orderByClause)
	buf.WriteString(" OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY")

	var dbImages []image
	if err := sqldb.NamedQuerySlice(ctx, s.log, s.db, buf.String(), data, &dbImages); err != nil {
		return nil, fmt.Errorf("namedqueryslice: %w", err)
	}

	return toBusImages(dbImages), nil
}

// Count returns the number of images matching the filter.
func (s *Store) Count(ctx context.Context, filter imagebus.QueryFilter) (int, error) {
	data := map[string]any{}

	const q = `
	SELECT
		count(1)
	FROM
		images`

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

// QueryByID retrieves an image by its ID.
func (s *Store) QueryByID(ctx context.Context, imageID uuid.UUID) (imagebus.Image, error) {
	data := struct {
		ID string `db:"image_id"`
	}{
		ID: imageID.String(),
	}

	const q = `
	SELECT
		image_id, restaurant_id, entity_type, object_path, public_url, content_type, size_bytes, status, uploaded_by, date_created, date_updated
	FROM
		images
	WHERE
		image_id = :image_id`

	var dbImage image
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, q, data, &dbImage); err != nil {
		if errors.Is(err, sqldb.ErrDBNotFound) {
			return imagebus.Image{}, imagebus.ErrNotFound
		}
		return imagebus.Image{}, fmt.Errorf("namedquerystruct: %w", err)
	}

	return toBusImage(dbImage), nil
}

func (s *Store) applyFilter(filter imagebus.QueryFilter, data map[string]any, buf *bytes.Buffer) {
	var wc []string

	if filter.ID != nil {
		data["image_id"] = *filter.ID
		wc = append(wc, "image_id = :image_id")
	}

	if filter.RestaurantID != nil {
		data["restaurant_id"] = *filter.RestaurantID
		wc = append(wc, "restaurant_id = :restaurant_id")
	}

	if filter.EntityType != nil {
		data["entity_type"] = *filter.EntityType
		wc = append(wc, "entity_type = :entity_type")
	}

	if filter.Status != nil {
		data["status"] = *filter.Status
		wc = append(wc, "status = :status")
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
	imagebus.OrderByID:          "image_id",
	imagebus.OrderByDateCreated: "date_created",
}

func orderByClause(orderBy order.By) (string, error) {
	by, ok := orderByFields[orderBy.Field]
	if !ok {
		return "", fmt.Errorf("field %q does not exist", orderBy.Field)
	}
	return "ORDER BY " + by + " " + orderBy.Direction, nil
}
