package addonapi_test

import (
	"testing"

	"github.com/warlck/food-flow/app/sdk/apitest"
)

func Test_Addon(t *testing.T) {
	t.Parallel()

	test := apitest.New(t, "Test_Addon")

	// -------------------------------------------------------------------------

	sd, err := insertSeedData(test.DB, test.Auth)
	if err != nil {
		t.Fatalf("Seeding error: %s", err)
	}

	// -------------------------------------------------------------------------

	test.Run(t, query200(sd), "query-200")
	test.Run(t, queryByID200(sd), "querybyid-200")

	test.Run(t, create201(sd), "create-201")
	test.Run(t, create400(sd), "create-400")
	test.Run(t, create401(sd), "create-401")

	test.Run(t, update200(sd), "update-200")
	test.Run(t, update400(sd), "update-400")
	test.Run(t, reorder200(sd), "reorder-200")
	test.Run(t, reorder400(sd), "reorder-400")
	test.Run(t, reorder401(sd), "reorder-401")

	test.Run(t, delete200(sd), "delete-200")
	test.Run(t, delete401(sd), "delete-401")
}
