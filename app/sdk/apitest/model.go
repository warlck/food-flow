package apitest

import (
	"github.com/warlck/food-flow/business/domain/categorybus"
	"github.com/warlck/food-flow/business/domain/menuitembus"
	"github.com/warlck/food-flow/business/domain/orderbus"
	"github.com/warlck/food-flow/business/domain/restaurantbus"
	"github.com/warlck/food-flow/business/domain/userbus"
)

// User extends the userbus.User with token for api test support.
type User struct {
	userbus.User
	Token string
}

// Restaurant extends the restaurantbus.Restaurant for api test support.
type Restaurant struct {
	restaurantbus.Restaurant
}

// Category extends the categorybus.Category for api test support.
type Category struct {
	categorybus.Category
}

// MenuItem extends the menuitembus.MenuItem for api test support.
type MenuItem struct {
	menuitembus.MenuItem
}

// Order extends the orderbus.Order for api test support.
type Order struct {
	orderbus.Order
}

// SeedData represents data for api tests.
type SeedData struct {
	Users       []User
	Admins      []User
	Restaurants []Restaurant
	Categories  []Category
	MenuItems   []MenuItem
	Orders      []Order
}

// Table represent fields needed for running an api test.
type Table struct {
	Name       string
	URL        string
	Token      string
	Method     string
	StatusCode int
	Input      any
	GotResp    any
	ExpResp    any
	CmpFunc    func(got any, exp any) string
}
