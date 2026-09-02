package modifiergroupapi_test

import (
	"testing"

	"github.com/warlck/food-flow/app/sdk/apitest"
)

func Test_ModifierGroup(t *testing.T) {
	t.Parallel()

	test := apitest.New(t, "Test_ModifierGroup")

	// -------------------------------------------------------------------------

	sd, err := insertSeedData(test.DB, test.Auth)
	if err != nil {
		t.Fatalf("Seeding error: %s", err)
	}

	// -------------------------------------------------------------------------

	test.Run(t, query200(sd), "query-200")
	test.Run(t, query400(sd), "query-400")
	test.Run(t, query403(sd), "query-403")
	test.Run(t, queryByID200(sd), "querybyid-200")
	test.Run(t, queryByID403(sd), "querybyid-403")
	test.Run(t, queryByID404(sd), "querybyid-404")

	test.Run(t, update200(sd), "update-200")
	test.Run(t, update400(sd), "update-400")
	test.Run(t, update404(sd), "update-404")

	test.Run(t, delete404(sd), "delete-404")
}
