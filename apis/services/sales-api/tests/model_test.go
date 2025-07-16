package tests

import (
	"time"

	"github.com/warlck/food-flow/apis/services/sales-api/handlers/userapi"
	"github.com/warlck/food-flow/business/domain/userbus"
	"github.com/warlck/food-flow/business/sdk/errs"
)

func toErrorPtr(err errs.Error) *errs.Error {
	return &err
}

func toAppUser(usr userbus.User) userapi.User {
	roles := make([]string, len(usr.Roles))
	for i, role := range usr.Roles {
		roles[i] = role.Name()
	}

	return userapi.User{
		ID:           usr.ID.String(),
		Name:         usr.Name,
		Email:        usr.Email.Address,
		Roles:        roles,
		PasswordHash: nil, // This field is not marshalled.
		Department:   usr.Department,
		Enabled:      usr.Enabled,
		DateCreated:  usr.DateCreated.Format(time.RFC3339),
		DateUpdated:  usr.DateUpdated.Format(time.RFC3339),
	}
}

func toAppUsers(users []userbus.User) []userapi.User {
	items := make([]userapi.User, len(users))
	for i, usr := range users {
		items[i] = toAppUser(usr)
	}

	return items
}

func toAppUserPtr(usr userbus.User) *userapi.User {
	appUsr := toAppUser(usr)
	return &appUsr
}
