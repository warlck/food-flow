package orderapi_test

import (
	"testing"

	"github.com/warlck/food-flow/app/sdk/apitest"
)

func Test_Order(t *testing.T) {
	t.Parallel()

	test := apitest.New(t, "Test_Order")

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

	test.Run(t, updateStatus200(sd), "updatestatus-200")
	test.Run(t, updateStatus401(sd), "updatestatus-401")

	test.Run(t, cancel200(sd), "cancel-200")
	test.Run(t, cancel401(sd), "cancel-401")
}
