package organizationapp

import (
	"time"

	"github.com/warlck/food-flow/business/domain/organizationbus"
)

type Organization struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	DateCreated string `json:"dateCreated"`
	DateUpdated string `json:"dateUpdated"`
}

// ToAppOrganization converts a business organization to an app organization.
func ToAppOrganization(bus organizationbus.Organization) Organization {
	return Organization{
		ID:          bus.ID.String(),
		Name:        bus.Name.String(),
		DateCreated: bus.DateCreated.Format(time.RFC3339),
		DateUpdated: bus.DateUpdated.Format(time.RFC3339),
	}
}

// ToAppOrganizations converts a slice of business organizations to app organizations.
func ToAppOrganizations(orgs []organizationbus.Organization) []Organization {
	appOrgs := make([]Organization, len(orgs))
	for i, org := range orgs {
		appOrgs[i] = ToAppOrganization(org)
	}
	return appOrgs
}
