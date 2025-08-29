package tests

import (
	"runtime/debug"
	"testing"
)

func Test_User(t *testing.T) {
	t.Parallel()

	// -------------------------------------------------------------------------

	apiTest := startTest(t, "Test_User")
	defer func() {
		if r := recover(); r != nil {
			t.Log(r)
			t.Error(string(debug.Stack()))
		}
		apiTest.DB.Teardown()
	}()

	// -------------------------------------------------------------------------

	sd, err := userSeedData(apiTest.DB, apiTest.Auth)
	if err != nil {
		t.Fatalf("Seeding error: %s", err)
	}

	// -------------------------------------------------------------------------

	apiTest.Run(t, query200(sd), "user-query-200")
	apiTest.Run(t, query400(sd), "user-querybyid-400")
	apiTest.Run(t, queryByID200(sd), "user-querybyid-200")

	apiTest.Run(t, create201(sd), "user-create-200")
	apiTest.Run(t, create401(sd), "user-create-401")
	apiTest.Run(t, create400(sd), "user-create-400")

}
