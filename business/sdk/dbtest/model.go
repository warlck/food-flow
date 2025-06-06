package dbtest

import "github.com/warlck/food-flow/business/domain/userbus"

// User represents an app user specified for the test.
type User struct {
	userbus.User
}

// SeedData represents data that was seeded for the test.
type SeedData struct {
	Users  []User
	Admins []User
}
