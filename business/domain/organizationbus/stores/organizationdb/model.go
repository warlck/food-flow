package organizationdb

import (
	"time"

	"github.com/google/uuid"
	"github.com/warlck/food-flow/business/domain/organizationbus"
	"github.com/warlck/food-flow/business/types/name"
	"github.com/warlck/food-flow/business/types/role"
)

// dbOrganization represents an individual organization.
type dbOrganization struct {
	ID          uuid.UUID `db:"organization_id"`
	Name        string    `db:"name"`
	DateCreated time.Time `db:"date_created"`
	DateUpdated time.Time `db:"date_updated"`
}

func toDBOrganization(bus organizationbus.Organization) dbOrganization {
	return dbOrganization{
		ID:          bus.ID,
		Name:        bus.Name.String(),
		DateCreated: bus.DateCreated,
		DateUpdated: bus.DateUpdated,
	}
}

func toBusOrganization(db dbOrganization) (organizationbus.Organization, error) {
	nme, err := name.Parse(db.Name)
	if err != nil {
		return organizationbus.Organization{}, err
	}

	bus := organizationbus.Organization{
		ID:          db.ID,
		Name:        nme,
		DateCreated: db.DateCreated.In(time.Local),
		DateUpdated: db.DateUpdated.In(time.Local),
	}

	return bus, nil
}

func toBusOrganizations(dbs []dbOrganization) ([]organizationbus.Organization, error) {
	bus := make([]organizationbus.Organization, len(dbs))
	for i, db := range dbs {
		var err error
		bus[i], err = toBusOrganization(db)
		if err != nil {
			return nil, err
		}
	}
	return bus, nil
}

// dbOrganizationUser represents the mapping between an organization and a user.
type dbOrganizationUser struct {
	OrganizationID uuid.UUID `db:"organization_id"`
	UserID         uuid.UUID `db:"user_id"`
	Role           string    `db:"role"`
	DateCreated    time.Time `db:"date_created"`
}

func toDBOrganizationUser(bus organizationbus.OrganizationUser) dbOrganizationUser {
	return dbOrganizationUser{
		OrganizationID: bus.OrganizationID,
		UserID:         bus.UserID,
		Role:           bus.Role.String(),
		DateCreated:    bus.DateCreated.UTC(),
	}
}

func toBusOrganizationUser(db dbOrganizationUser) (organizationbus.OrganizationUser, error) {
	r, err := role.Parse(db.Role)
	if err != nil {
		return organizationbus.OrganizationUser{}, err
	}

	bus := organizationbus.OrganizationUser{
		OrganizationID: db.OrganizationID,
		UserID:         db.UserID,
		Role:           r,
		DateCreated:    db.DateCreated.In(time.Local),
	}

	return bus, nil
}

func toBusOrganizationUsers(dbs []dbOrganizationUser) ([]organizationbus.OrganizationUser, error) {
	bus := make([]organizationbus.OrganizationUser, len(dbs))
	for i, db := range dbs {
		var err error
		bus[i], err = toBusOrganizationUser(db)
		if err != nil {
			return nil, err
		}
	}
	return bus, nil
}
