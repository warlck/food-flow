package apitest

import (
	"net/http/httptest"
	"testing"

	"github.com/warlck/food-flow/app/sdk/mux"

	authbuild "github.com/warlck/food-flow/api/services/auth/build/all"
	salesbuild "github.com/warlck/food-flow/api/services/sales/build/all"
	"github.com/warlck/food-flow/app/sdk/auth"
	"github.com/warlck/food-flow/app/sdk/authclient"
	"github.com/warlck/food-flow/business/sdk/dbtest"
)

// New initialized the system to run a test.
func New(t *testing.T, testName string) *Test {
	db := dbtest.New(t, testName)

	// -------------------------------------------------------------------------

	auth := auth.New(auth.Config{
		Log:       db.Log,
		KeyLookup: &KeyStore{},
	})

	// -------------------------------------------------------------------------

	server := httptest.NewServer(mux.WebAPI(mux.Config{
		Log: db.Log,
		DB:  db.DB,
		BusConfig: mux.BusConfig{
			UserBus:       db.BusDomain.User,
			RestaurantBus: db.BusDomain.Restaurant,
			CategoryBus:   db.BusDomain.Category,
			MenuItemBus:   db.BusDomain.MenuItem,
			OrderBus:      db.BusDomain.Order,
			AddonBus:      db.BusDomain.Addon,
		},
		Auth: auth,
	}, authbuild.Routes()))

	authClient := authclient.New(db.Log, server.URL)

	// -------------------------------------------------------------------------

	mux := mux.WebAPI(mux.Config{
		Log: db.Log,
		DB:  db.DB,
		BusConfig: mux.BusConfig{
			UserBus:       db.BusDomain.User,
			RestaurantBus: db.BusDomain.Restaurant,
			CategoryBus:   db.BusDomain.Category,
			MenuItemBus:   db.BusDomain.MenuItem,
			OrderBus:      db.BusDomain.Order,
			AddonBus:      db.BusDomain.Addon,
		},
		AuthClient: authClient,
	}, salesbuild.Routes())

	return &Test{
		DB:   db,
		Auth: auth,
		mux:  mux,
	}
}
