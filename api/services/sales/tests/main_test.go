package tests

import (
	"fmt"
	"os"
	"testing"

	"github.com/warlck/food-flow/app/sdk/apitest"
	"github.com/warlck/food-flow/business/sdk/dbtest"
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
	// db := dbtest.NewDatabase(t, c, testName)

	// // -------------------------------------------------------------------------

	// ac := authclient.New(log, "http://auth:3000")

	// // -------------------------------------------------------------------------

	// mux := mux.WebAPI(mux.Config{
	// 	Log:  db.Log,
	// 	Auth: ac,
	// 	DB:   db.DB,
	// })

	// return apitest.New(db, ac, mux)
	return nil
}
