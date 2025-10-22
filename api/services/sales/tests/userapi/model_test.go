package userapi_test

import (
	"github.com/warlck/food-flow/app/domain/userapp"
	"github.com/warlck/food-flow/business/domain/userbus"
)

func toAppUser(bus userbus.User) userapp.User {
	return userapp.ToAppUser(bus)
}

func toAppUsers(users []userbus.User) []userapp.User {
	return userapp.ToAppUsers(users)
}

func toAppUserPtr(bus userbus.User) *userapp.User {
	appUsr := toAppUser(bus)
	return &appUsr
}
