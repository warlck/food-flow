package organizationbus

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/warlck/food-flow/business/types/name"
)

// TestNewOrganizations is a helper method for testing.
func TestNewOrganizations(n int) []NewOrganization {
	newOrgs := make([]NewOrganization, n)

	idx := rand.Intn(10000)
	for i := range n {
		idx++

		no := NewOrganization{
			Name: name.MustParse(fmt.Sprintf("Org%d", idx)),
		}

		newOrgs[i] = no
	}

	return newOrgs
}

// TestSeedOrganizations is a helper method for testing.
func TestSeedOrganizations(ctx context.Context, n int, bus *Business) ([]Organization, error) {
	newOrgs := TestNewOrganizations(n)

	orgs := make([]Organization, len(newOrgs))
	for i, no := range newOrgs {
		org, err := bus.Create(ctx, no)
		if err != nil {
			return nil, fmt.Errorf("seeding organization: idx: %d : %w", i, err)
		}

		orgs[i] = org
	}

	return orgs, nil
}
