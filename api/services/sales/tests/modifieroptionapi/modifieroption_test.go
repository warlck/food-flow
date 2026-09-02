package modifieroptionapi_test

import (
	"testing"

	"github.com/warlck/food-flow/app/sdk/apitest"
)

func Test_ModifierOption(t *testing.T) {
	t.Parallel()

	test := apitest.New(t, "Test_ModifierOption")

	// -------------------------------------------------------------------------

	sd, err := insertSeedData(test.DB, test.Auth)
	if err != nil {
		t.Fatalf("Seeding error: %s", err)
	}

	// -------------------------------------------------------------------------

	test.Run(t, queryByID200(sd), "querybyid-200")
	test.Run(t, queryByID404(sd), "querybyid-404")

	test.Run(t, update404(sd), "update-404")

	test.Run(t, delete404(sd), "delete-404")
}
