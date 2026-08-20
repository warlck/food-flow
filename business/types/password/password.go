// Package password represents a password in the system.
package password

import (
	"fmt"
	"regexp"
)

// Password represents a password in the system.
type Password struct {
	value string
}

// String returns the value of the password.
func (p Password) String() string {
	return p.value
}

// Equal provides support for the go-cmp package and testing.
func (p Password) Equal(n2 Password) bool {
	return p.value == n2.value
}

// MarshalText provides support for logging and any marshal needs.
func (p Password) MarshalText() ([]byte, error) {
	return []byte(p.value), nil
}

// =============================================================================

// Passwords are 12-64 printable ASCII characters (no spaces or control
// characters). The policy is enforced at the app edge (user creation,
// password change, useradd tooling); login bcrypt-compares the raw string
// and never parses, so existing passwords keep working.
var passwordRegEx = regexp.MustCompile(`^[\x21-\x7E]{12,64}$`)

// Parse parses the string value and returns a password if the value complies
// with the rules for a password.
func Parse(value string) (Password, error) {
	if !passwordRegEx.MatchString(value) {
		return Password{}, fmt.Errorf("invalid password %q", value)
	}

	return Password{value}, nil
}

// MustParse parses the string value and returns a password if the value
// complies with the rules for a password. If an error occurs the function panics.
// Use it for testing only
func MustParse(value string) Password {
	password, err := Parse(value)
	if err != nil {
		panic(err)
	}

	return password
}

func ParseConfirm(pass string, confirm string) (Password, error) {
	p, err := Parse(pass)
	if err != nil {
		return Password{}, err
	}

	if pass != confirm {
		return Password{}, fmt.Errorf("passwords do not match")
	}

	return p, nil
}

func ParseConfirmPointers(pass *string, confirm *string) (Password, error) {
	if pass == nil || confirm == nil {
		return Password{}, fmt.Errorf("passwords do not match")
	}

	return ParseConfirm(*pass, *confirm)
}
