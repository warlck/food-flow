package organizationbus

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/warlck/food-flow/business/domain/userbus"
	"github.com/warlck/food-flow/business/types/role"
	"github.com/warlck/food-flow/foundation/logger"
)

// Set of error variables for CRUD operations.
var (
	ErrNotFound      = errors.New("organization not found")
	ErrNotAuthorized = errors.New("user is not authorized to be associated with an organization")
)

// Storer interface declares the behavior this package needs to persist and
// retrieve data.
type Storer interface {
	Create(ctx context.Context, org Organization) error
	Update(ctx context.Context, org Organization) error
	Delete(ctx context.Context, org Organization) error
	QueryByID(ctx context.Context, organizationID uuid.UUID) (Organization, error)
	AddUser(ctx context.Context, orgUser OrganizationUser) error
	RemoveUser(ctx context.Context, orgUser OrganizationUser) error
	QueryOrgsForUser(ctx context.Context, userID uuid.UUID) ([]Organization, error)
	QueryUsersForOrg(ctx context.Context, organizationID uuid.UUID) ([]OrganizationUser, error)
}

// UserDelegate defines the behavior needed for user queries.
type UserDelegate interface {
	QueryByID(ctx context.Context, userID uuid.UUID) (userbus.User, error)
}

// Business manages the set of APIs for organization access.
type Business struct {
	log     *logger.Logger
	storer  Storer
	userBus UserDelegate
}

// NewBusiness constructs an organization business API for use.
func NewBusiness(log *logger.Logger, storer Storer, userBus UserDelegate) *Business {
	return &Business{
		log:     log,
		storer:  storer,
		userBus: userBus,
	}
}

// Create adds a new organization to the system.
func (b *Business) Create(ctx context.Context, no NewOrganization) (Organization, error) {
	now := time.Now()

	org := Organization{
		ID:          uuid.New(),
		Name:        no.Name,
		DateCreated: now,
		DateUpdated: now,
	}

	if err := b.storer.Create(ctx, org); err != nil {
		return Organization{}, fmt.Errorf("create: %w", err)
	}

	return org, nil
}

// Update modifies information about an organization.
func (b *Business) Update(ctx context.Context, org Organization, uo UpdateOrganization) (Organization, error) {
	if uo.Name != nil {
		org.Name = *uo.Name
	}

	org.DateUpdated = time.Now()

	if err := b.storer.Update(ctx, org); err != nil {
		return Organization{}, fmt.Errorf("update: %w", err)
	}

	return org, nil
}

// Delete removes an organization from the system.
func (b *Business) Delete(ctx context.Context, org Organization) error {
	if err := b.storer.Delete(ctx, org); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}

// QueryByID finds the organization identified by a given ID.
func (b *Business) QueryByID(ctx context.Context, organizationID uuid.UUID) (Organization, error) {
	org, err := b.storer.QueryByID(ctx, organizationID)
	if err != nil {
		return Organization{}, fmt.Errorf("query: %w", err)
	}

	return org, nil
}

// AddUser adds a user to an organization, ensuring they are an ADMIN.
func (b *Business) AddUser(ctx context.Context, nou NewOrganizationUser) (OrganizationUser, error) {
	usr, err := b.userBus.QueryByID(ctx, nou.UserID)
	if err != nil {
		return OrganizationUser{}, fmt.Errorf("query user: %w", err)
	}

	isAdmin := false
	for _, r := range usr.Roles {
		if r == role.Admin {
			isAdmin = true
			break
		}
	}

	if !isAdmin {
		return OrganizationUser{}, ErrNotAuthorized
	}

	ou := OrganizationUser{
		OrganizationID: nou.OrganizationID,
		UserID:         nou.UserID,
		Role:           nou.Role,
		DateCreated:    time.Now(),
	}

	if err := b.storer.AddUser(ctx, ou); err != nil {
		return OrganizationUser{}, fmt.Errorf("add user: %w", err)
	}

	return ou, nil
}

// RemoveUser removes a user from an organization.
func (b *Business) RemoveUser(ctx context.Context, orgUser OrganizationUser) error {
	if err := b.storer.RemoveUser(ctx, orgUser); err != nil {
		return fmt.Errorf("remove user: %w", err)
	}

	return nil
}

// QueryOrgsForUser retrieves a list of organizations a user belongs to.
func (b *Business) QueryOrgsForUser(ctx context.Context, userID uuid.UUID) ([]Organization, error) {
	orgs, err := b.storer.QueryOrgsForUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("query orgs: %w", err)
	}

	return orgs, nil
}

// QueryUsersForOrg retrieves a list of users for an organization.
func (b *Business) QueryUsersForOrg(ctx context.Context, organizationID uuid.UUID) ([]OrganizationUser, error) {
	users, err := b.storer.QueryUsersForOrg(ctx, organizationID)
	if err != nil {
		return nil, fmt.Errorf("query users: %w", err)
	}

	return users, nil
}
