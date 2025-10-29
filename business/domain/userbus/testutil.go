package userbus

import (
	"context"
	"fmt"
	"math/rand"
	"net/mail"

	"github.com/warlck/food-flow/business/types/name"
	"github.com/warlck/food-flow/business/types/password"
	"github.com/warlck/food-flow/business/types/role"
)

// TestNewUsers is a helper method for testing.
func TestNewUsers(n int, rle role.Role) []NewUser {
	newUsrs := make([]NewUser, n)

	idx := rand.Intn(10000)
	for i := range n {
		idx++

		nu := NewUser{
			Name:       name.MustParse(fmt.Sprintf("Name%d", idx)),
			Email:      mail.Address{Address: fmt.Sprintf("Email%d@gmail.com", idx)},
			Roles:      []role.Role{rle},
			Department: name.MustParseNull(fmt.Sprintf("Department%d", idx)),
			Password:   password.MustParse("123"),
		}

		newUsrs[i] = nu
	}

	return newUsrs
}

// TestSeedUsers is a helper method for testing.
func TestSeedUsers(ctx context.Context, n int, rle role.Role, api *Business) ([]User, error) {
	newUsrs := TestNewUsers(n, rle)

	usrs := make([]User, len(newUsrs))
	for i, nu := range newUsrs {
		usr, err := api.Create(ctx, nu)
		if err != nil {
			return nil, fmt.Errorf("seeding user: idx: %d : %w", i, err)
		}

		usrs[i] = usr
	}

	return usrs, nil
}
