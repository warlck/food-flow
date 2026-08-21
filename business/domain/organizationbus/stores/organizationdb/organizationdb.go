// Package organizationdb contains organization related CRUD functionality.
package organizationdb

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/warlck/food-flow/business/domain/organizationbus"
	"github.com/warlck/food-flow/business/sdk/sqldb"
	"github.com/warlck/food-flow/foundation/logger"
)

// Store manages the set of APIs for organization database access.
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

// Create inserts a new organization into the database.
func (s *Store) Create(ctx context.Context, org organizationbus.Organization) error {
	const q = `
	INSERT INTO organizations
		(organization_id, name, date_created, date_updated)
	VALUES
		(:organization_id, :name, :date_created, :date_updated)`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, toDBOrganization(org)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// Update replaces an organization document in the database.
func (s *Store) Update(ctx context.Context, org organizationbus.Organization) error {
	const q = `
	UPDATE
		organizations
	SET 
		"name" = :name,
		"date_updated" = :date_updated
	WHERE
		organization_id = :organization_id`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, toDBOrganization(org)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// Delete removes an organization from the database.
func (s *Store) Delete(ctx context.Context, org organizationbus.Organization) error {
	data := struct {
		ID string `db:"organization_id"`
	}{
		ID: org.ID.String(),
	}

	const q = `
	DELETE FROM
		organizations
	WHERE
		organization_id = :organization_id`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, data); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// QueryByID gets the specified organization from the database.
func (s *Store) QueryByID(ctx context.Context, organizationID uuid.UUID) (organizationbus.Organization, error) {
	data := struct {
		ID string `db:"organization_id"`
	}{
		ID: organizationID.String(),
	}

	const q = `
	SELECT
		organization_id, name, date_created, date_updated
	FROM
		organizations
	WHERE 
		organization_id = :organization_id`

	var dbOrg dbOrganization
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, q, data, &dbOrg); err != nil {
		if errors.Is(err, sqldb.ErrDBNotFound) {
			return organizationbus.Organization{}, fmt.Errorf("db: %w", organizationbus.ErrNotFound)
		}
		return organizationbus.Organization{}, fmt.Errorf("namedquerystruct: %w", err)
	}

	return toBusOrganization(dbOrg)
}

// =========================================================================
// Organization Users

// AddUser adds a user to an organization.
func (s *Store) AddUser(ctx context.Context, orgUser organizationbus.OrganizationUser) error {
	const q = `
	INSERT INTO organization_users
		(organization_id, user_id, role, date_created)
	VALUES
		(:organization_id, :user_id, :role, :date_created)`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, toDBOrganizationUser(orgUser)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// RemoveUser removes a user from an organization.
func (s *Store) RemoveUser(ctx context.Context, orgUser organizationbus.OrganizationUser) error {
	data := struct {
		OrganizationID string `db:"organization_id"`
		UserID         string `db:"user_id"`
	}{
		OrganizationID: orgUser.OrganizationID.String(),
		UserID:         orgUser.UserID.String(),
	}

	const q = `
	DELETE FROM
		organization_users
	WHERE
		organization_id = :organization_id AND user_id = :user_id`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, data); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// QueryOrgsForUser retrieves a list of organizations a user belongs to.
func (s *Store) QueryOrgsForUser(ctx context.Context, userID uuid.UUID) ([]organizationbus.Organization, error) {
	data := struct {
		UserID string `db:"user_id"`
	}{
		UserID: userID.String(),
	}

	const q = `
	SELECT
		o.organization_id, o.name, o.date_created, o.date_updated
	FROM
		organizations AS o
	INNER JOIN
		organization_users AS ou ON o.organization_id = ou.organization_id
	WHERE
		ou.user_id = :user_id`

	var dbs []dbOrganization
	if err := sqldb.NamedQuerySlice(ctx, s.log, s.db, q, data, &dbs); err != nil {
		if errors.Is(err, sqldb.ErrDBNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("namedqueryslice: %w", err)
	}

	return toBusOrganizations(dbs)
}

// QueryUsersForOrg retrieves a list of users for an organization.
func (s *Store) QueryUsersForOrg(ctx context.Context, organizationID uuid.UUID) ([]organizationbus.OrganizationUser, error) {
	data := struct {
		OrganizationID string `db:"organization_id"`
	}{
		OrganizationID: organizationID.String(),
	}

	const q = `
	SELECT
		organization_id, user_id, role, date_created
	FROM
		organization_users
	WHERE
		organization_id = :organization_id`

	var dbs []dbOrganizationUser
	if err := sqldb.NamedQuerySlice(ctx, s.log, s.db, q, data, &dbs); err != nil {
		if errors.Is(err, sqldb.ErrDBNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("namedqueryslice: %w", err)
	}

	return toBusOrganizationUsers(dbs)
}
