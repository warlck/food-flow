package userbus_test

import (
	"context"
	"fmt"
	"net/mail"
	"os"
	"runtime/debug"
	"testing"

	"github.com/google/go-cmp/cmp"
	"golang.org/x/crypto/bcrypt"

	"github.com/warlck/food-flow/business/domain/userbus"
	"github.com/warlck/food-flow/business/sdk/dbtest"
	"github.com/warlck/food-flow/business/sdk/unittest"
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
func Test_User(t *testing.T) {
	t.Parallel()

	db := dbtest.NewDatabase(t, c, "Test_User")
	defer func() {
		if r := recover(); r != nil {
			t.Log(r)
			t.Error(string(debug.Stack()))
		}
		db.Teardown()
	}()

	_, err := userSeedData(db)
	if err != nil {
		t.Fatalf("Seeding error: %s", err)
	}

	// -------------------------------------------------------------------------

	unittest.Run(t, create(db), "userbus-create")

}

// =============================================================================

// =============================================================================

func userSeedData(db *dbtest.Database) (dbtest.SeedData, error) {
	ctx := context.Background()
	busDomain := db.BusDomain

	usrs, err := userbus.TestSeedUsers(ctx, 2, userbus.RoleAdmin, busDomain.User)
	if err != nil {
		return dbtest.SeedData{}, fmt.Errorf("seeding users : %w", err)
	}

	tu1 := dbtest.User{
		User: usrs[0],
	}

	tu2 := dbtest.User{
		User: usrs[1],
	}

	// -------------------------------------------------------------------------

	usrs, err = userbus.TestSeedUsers(ctx, 2, userbus.RoleUser, busDomain.User)
	if err != nil {
		return dbtest.SeedData{}, fmt.Errorf("seeding users : %w", err)
	}

	tu3 := dbtest.User{
		User: usrs[0],
	}

	tu4 := dbtest.User{
		User: usrs[1],
	}

	// -------------------------------------------------------------------------

	sd := dbtest.SeedData{
		Users:  []dbtest.User{tu3, tu4},
		Admins: []dbtest.User{tu1, tu2},
	}

	return sd, nil
}

// =============================================================================
func create(db *dbtest.Database) []unittest.Table {
	email, _ := mail.ParseAddress("adil@gmail.com")

	table := []unittest.Table{
		{
			Name: "basic",
			ExpResp: userbus.User{
				Name:       "Adil Zitdinov",
				Email:      *email,
				Roles:      []userbus.Role{userbus.RoleAdmin},
				Department: "IT",
				Enabled:    true,
			},
			ExcFunc: func(ctx context.Context) any {
				nu := userbus.NewUser{
					Name:       "Adil Zitdinov",
					Email:      *email,
					Roles:      []userbus.Role{userbus.RoleAdmin},
					Department: "IT",
					Password:   "123",
				}

				resp, err := db.BusDomain.User.Create(ctx, nu)
				if err != nil {
					return err
				}

				return resp
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.(userbus.User)
				if !exists {
					return "error occurred"
				}

				if err := bcrypt.CompareHashAndPassword(gotResp.PasswordHash, []byte("123")); err != nil {
					return err.Error()
				}

				expResp := exp.(userbus.User)

				expResp.ID = gotResp.ID
				expResp.PasswordHash = gotResp.PasswordHash
				expResp.DateCreated = gotResp.DateCreated
				expResp.DateUpdated = gotResp.DateUpdated

				return cmp.Diff(gotResp, expResp)
			},
		},
	}

	return table
}
