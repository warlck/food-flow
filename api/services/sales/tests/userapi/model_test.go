package userapi_test

import (
	"time"

	"github.com/warlck/food-flow/app/domain/userapi"
	"github.com/warlck/food-flow/business/domain/userbus"
	"github.com/warlck/food-flow/business/types/role"
)

func toAppUser(bus userbus.User) userapi.User {
	return userapi.User{
		ID:          bus.ID.String(),
		Name:        bus.Name.String(),
		Email:       bus.Email.Address,
		Roles:       role.ParseToString(bus.Roles),
		Department:  bus.Department.String(),
		Enabled:     bus.Enabled,
		DateCreated: bus.DateCreated.Format(time.RFC3339),
		DateUpdated: bus.DateUpdated.Format(time.RFC3339),
	}
}

func toAppUsers(users []userbus.User) []userapi.User {
	items := make([]userapi.User, len(users))
	for i, usr := range users {
		items[i] = toAppUser(usr)
	}

	return items
}

func toAppUserPtr(bus userbus.User) *userapi.User {
	appUsr := toAppUser(bus)
	return &appUsr
}
