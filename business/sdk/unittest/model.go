package unittest

import (
	"context"

	"github.com/warlck/food-flow/business/domain/categorybus"
	"github.com/warlck/food-flow/business/domain/menuitembus"
	"github.com/warlck/food-flow/business/domain/restaurantbus"
	"github.com/warlck/food-flow/business/domain/userbus"
)

// User represents an app user specified for the test.
type User struct {
	userbus.User
}

// Restaurant represents a restaurant specified for the test.
type Restaurant struct {
	restaurantbus.Restaurant
}

// Category represents a category specified for the test.
type Category struct {
	categorybus.Category
}

// MenuItem represents a menu item specified for the test.
type MenuItem struct {
	menuitembus.MenuItem
}

// SeedData represents data that was seeded for the test.
type SeedData struct {
	Users       []User
	Admins      []User
	Restaurants []Restaurant
	Categories  []Category
	MenuItems   []MenuItem
}

// Table represent fields needed for running an unit test.
type Table struct {
	Name    string
	ExpResp any
	ExcFunc func(ctx context.Context) any
	CmpFunc func(got any, exp any) string
}
