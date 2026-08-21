package organizationapi_test

import (
	"testing"

	"github.com/warlck/food-flow/app/sdk/apitest"
)

func Test_Organization(t *testing.T) {
	t.Parallel()

	test := apitest.New(t, "Test_Organization")

	sd, err := insertSeedData(test.DB, test.Auth)
	if err != nil {
		t.Fatalf("Seeding error: %s", err)
	}

	test.Run(t, queryMyOrgs200(sd), "querymyorgs-200")
	test.Run(t, queryMyOrgs401(sd), "querymyorgs-401")
}
