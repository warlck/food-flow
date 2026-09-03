// Package modifiergroupdb contains modifier group related CRUD functionality.
package modifiergroupdb

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/warlck/food-flow/business/domain/modifiergroupbus"
	"github.com/warlck/food-flow/business/sdk/order"
	"github.com/warlck/food-flow/business/sdk/page"
	"github.com/warlck/food-flow/business/sdk/sqldb"
	"github.com/warlck/food-flow/foundation/logger"
)

// Store manages the set of APIs for modifier group database access.
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

// Create inserts a new modifier group into the database.
func (s *Store) Create(ctx context.Context, group modifiergroupbus.ModifierGroup) error {
	const q = `
	INSERT INTO modifier_groups
		(modifier_group_id, menu_item_id, restaurant_id, name, description, min_selections, max_selections, available, rank, date_created, date_updated)
	VALUES
		(:modifier_group_id, :menu_item_id, :restaurant_id, :name, :description, :min_selections, :max_selections, :available, :rank, :date_created, :date_updated)`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, toDBModifierGroup(group)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// Update replaces a modifier group document in the database.
func (s *Store) Update(ctx context.Context, group modifiergroupbus.ModifierGroup) error {
	const q = `
	UPDATE
		modifier_groups
	SET 
		"name" = :name,
		"description" = :description,
		"min_selections" = :min_selections,
		"max_selections" = :max_selections,
		"available" = :available,
		"rank" = :rank,
		"date_updated" = :date_updated
	WHERE
		modifier_group_id = :modifier_group_id`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, toDBModifierGroup(group)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// Reorder verifies the exact current set and updates modifier group ranks in one transaction.
func (s *Store) Reorder(ctx context.Context, menuItemID uuid.UUID, orderedIDs []uuid.UUID) ([]modifiergroupbus.ModifierGroup, error) {
	tx, err := s.db.(*sqlx.DB).BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	data := struct {
		MenuItemID uuid.UUID `db:"menu_item_id"`
	}{
		MenuItemID: menuItemID,
	}

	// Lock the parent so concurrent child inserts cannot pass their foreign-key
	// check until the exact-set validation and rank updates commit.
	const lockParent = `
	SELECT
		menu_item_id
	FROM
		menu_items
	WHERE
		menu_item_id = :menu_item_id
	FOR UPDATE`

	var parent struct {
		MenuItemID uuid.UUID `db:"menu_item_id"`
	}
	if err := sqldb.NamedQueryStruct(ctx, s.log, tx, lockParent, data, &parent); err != nil {
		return nil, fmt.Errorf("lock menu item for modifier group reorder: %w", err)
	}

	const query = `
	SELECT
		modifier_group_id, menu_item_id, restaurant_id, name, description, min_selections, max_selections, available, rank, date_created, date_updated
	FROM
		modifier_groups
	WHERE
		menu_item_id = :menu_item_id
	ORDER BY
		modifier_group_id
	FOR UPDATE`

	var dbGroups []modifierGroup
	if err := sqldb.NamedQuerySlice(ctx, s.log, tx, query, data, &dbGroups); err != nil {
		return nil, fmt.Errorf("query modifier groups for reorder: %w", err)
	}

	groups, err := toBusModifierGroups(dbGroups)
	if err != nil {
		return nil, fmt.Errorf("convert modifier groups for reorder: %w", err)
	}
	if len(groups) != len(orderedIDs) {
		return nil, fmt.Errorf("%w: exact set mismatch: expected %d modifier groups, got %d",
			modifiergroupbus.ErrInvalidReorder, len(groups), len(orderedIDs))
	}

	existing := make(map[uuid.UUID]modifiergroupbus.ModifierGroup, len(groups))
	for _, group := range groups {
		existing[group.ID] = group
	}

	now := time.Now()
	reordered := make([]modifiergroupbus.ModifierGroup, len(orderedIDs))
	for i, id := range orderedIDs {
		group, exists := existing[id]
		if !exists {
			return nil, fmt.Errorf("%w: modifier group id %s does not belong to menu item %s",
				modifiergroupbus.ErrInvalidReorder, id, menuItemID)
		}
		rank := (i + 1) * 10
		group.Rank = &rank
		group.DateUpdated = now
		reordered[i] = group
	}

	const update = `
	UPDATE
		modifier_groups
	SET
		"rank" = :rank,
		"date_updated" = :date_updated
	WHERE
		modifier_group_id = :modifier_group_id
		AND menu_item_id = :menu_item_id`

	for _, group := range reordered {
		if err := sqldb.NamedExecContext(ctx, s.log, tx, update, toDBModifierGroup(group)); err != nil {
			return nil, fmt.Errorf("update modifier group rank: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}

	return reordered, nil
}

// Delete removes a modifier group from the database.
func (s *Store) Delete(ctx context.Context, group modifiergroupbus.ModifierGroup) error {
	data := struct {
		ID string `db:"modifier_group_id"`
	}{
		ID: group.ID.String(),
	}

	const q = `
	DELETE FROM
		modifier_groups
	WHERE
		modifier_group_id = :modifier_group_id`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, data); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// Query retrieves a list of existing modifier groups from the database.
func (s *Store) Query(ctx context.Context, filter modifiergroupbus.QueryFilter, orderBy order.By, page page.Page) ([]modifiergroupbus.ModifierGroup, error) {
	data := map[string]any{
		"offset":        (page.Number() - 1) * page.RowsPerPage(),
		"rows_per_page": page.RowsPerPage(),
	}

	const q = `
	SELECT
		modifier_group_id, menu_item_id, restaurant_id, name, description, min_selections, max_selections, available, rank, date_created, date_updated
	FROM
		modifier_groups`

	buf := bytes.NewBufferString(q)
	applyFilter(filter, data, buf)

	orderByClause, err := orderByClause(orderBy)
	if err != nil {
		return nil, err
	}

	buf.WriteString(orderByClause)
	buf.WriteString(" OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY")

	var dbGroups []modifierGroup
	if err := sqldb.NamedQuerySlice(ctx, s.log, s.db, buf.String(), data, &dbGroups); err != nil {
		return nil, fmt.Errorf("namedqueryslice: %w", err)
	}

	return toBusModifierGroups(dbGroups)
}

// QueryAll retrieves all modifier groups matching the filter from the database without pagination.
func (s *Store) QueryAll(ctx context.Context, filter modifiergroupbus.QueryFilter, orderBy order.By) ([]modifiergroupbus.ModifierGroup, error) {
	data := map[string]any{}

	const q = `
	SELECT
		modifier_group_id, menu_item_id, restaurant_id, name, description, min_selections, max_selections, available, rank, date_created, date_updated
	FROM
		modifier_groups`

	buf := bytes.NewBufferString(q)
	applyFilter(filter, data, buf)

	orderByClause, err := orderByClause(orderBy)
	if err != nil {
		return nil, err
	}

	buf.WriteString(orderByClause)

	var dbGroups []modifierGroup
	if err := sqldb.NamedQuerySlice(ctx, s.log, s.db, buf.String(), data, &dbGroups); err != nil {
		return nil, fmt.Errorf("namedqueryslice: %w", err)
	}

	return toBusModifierGroups(dbGroups)
}

// Count returns the total number of modifier groups in the DB.
func (s *Store) Count(ctx context.Context, filter modifiergroupbus.QueryFilter) (int, error) {
	data := map[string]any{}

	const q = `
	SELECT
		count(1)
	FROM
		modifier_groups`

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

// QueryByID gets the specified modifier group from the database.
func (s *Store) QueryByID(ctx context.Context, groupID uuid.UUID) (modifiergroupbus.ModifierGroup, error) {
	data := struct {
		ID string `db:"modifier_group_id"`
	}{
		ID: groupID.String(),
	}

	const q = `
	SELECT
		modifier_group_id, menu_item_id, restaurant_id, name, description, min_selections, max_selections, available, rank, date_created, date_updated
	FROM
		modifier_groups
	WHERE 
		modifier_group_id = :modifier_group_id`

	var dbGroup modifierGroup
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, q, data, &dbGroup); err != nil {
		if errors.Is(err, sqldb.ErrDBNotFound) {
			return modifiergroupbus.ModifierGroup{}, fmt.Errorf("db: %w", modifiergroupbus.ErrNotFound)
		}
		return modifiergroupbus.ModifierGroup{}, fmt.Errorf("db: %w", err)
	}

	return toBusModifierGroup(dbGroup)
}
