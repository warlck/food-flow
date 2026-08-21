package organizationapp

import (
	"time"

	"github.com/warlck/food-flow/business/domain/organizationbus"
)

type Organization struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	DateCreated time.Time `json:"date_created"`
	DateUpdated time.Time `json:"date_updated"`
}

func toAppOrganization(bus organizationbus.Organization) Organization {
	return Organization{
		ID:          bus.ID.String(),
		Name:        bus.Name.String(),
		DateCreated: bus.DateCreated,
		DateUpdated: bus.DateUpdated,
	}
}

func toAppOrganizations(orgs []organizationbus.Organization) []Organization {
	appOrgs := make([]Organization, len(orgs))
	for i, org := range orgs {
		appOrgs[i] = toAppOrganization(org)
	}
	return appOrgs
}
