package tests

import (
	"fmt"
	"os"
	"testing"

	"github.com/warlck/food-flow/apis/services/sales-api/mux"
	"github.com/warlck/food-flow/business/sdk/dbtest"
	"github.com/warlck/food-flow/business/web/apitest"
	"github.com/warlck/food-flow/business/web/auth"
	"github.com/warlck/food-flow/foundation/docker"
)

var c *docker.Container

func TestMain(m *testing.M) {
	code, err := run(m)
	if err != nil {
		fmt.Println(err)
	}

	os.Exit(code)
}

func run(m *testing.M) (int, error) {
	var err error

	c, err = dbtest.StartDB()
	if err != nil {
		return 1, err
	}
	defer dbtest.StopDB(c)

	return m.Run(), nil
}

func startTest(t *testing.T, testName string) *apitest.Test {
	db := dbtest.NewDatabase(t, c, testName)

	// -------------------------------------------------------------------------

	auth, err := auth.NewAuth(auth.Config{
		Log:       db.Log,
		KeyLookup: &apitest.KeyStore{},
		Issuer:    "service api",
	})
	if err != nil {
		t.Fatal(err)
	}

	// -------------------------------------------------------------------------

	mux := mux.WebAPI(mux.Config{
		Log:  db.Log,
		Auth: auth,
		DB:   db.DB,
	})

	return apitest.New(db, auth, mux)
}
