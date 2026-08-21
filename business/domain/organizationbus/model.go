package organizationbus

import (
	"time"

	"github.com/google/uuid"
	"github.com/warlck/food-flow/business/types/name"
	"github.com/warlck/food-flow/business/types/role"
)

// Organization represents information about an organization.
type Organization struct {
	ID          uuid.UUID
	Name        name.Name
	DateCreated time.Time
	DateUpdated time.Time
}

// NewOrganization contains information needed to create a new organization.
type NewOrganization struct {
	Name name.Name
}

// UpdateOrganization contains information needed to update an organization.
type UpdateOrganization struct {
	Name *name.Name
}

// OrganizationUser represents the mapping between an organization and a user.
type OrganizationUser struct {
	OrganizationID uuid.UUID
	UserID         uuid.UUID
	Role           role.Role
	DateCreated    time.Time
}

// NewOrganizationUser contains information needed to add a user to an organization.
type NewOrganizationUser struct {
	OrganizationID uuid.UUID
	UserID         uuid.UUID
	Role           role.Role
}
