package password_test

import (
	"strings"
	"testing"

	"github.com/warlck/food-flow/business/types/password"
)

func Test_Parse(t *testing.T) {
	t.Parallel()

	table := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"empty", "", true},
		{"11 chars rejected", strings.Repeat("a", 11), true},
		{"12 chars accepted", strings.Repeat("a", 12), false},
		{"64 chars accepted", strings.Repeat("a", 64), false},
		{"65 chars rejected", strings.Repeat("a", 65), true},
		{"space rejected", "abcd efgh ijkl", true},
		{"printable symbols accepted", `!@#$%^&*()_+-=[]{};':",./<>?` + "\\|`~", false},
		{"tilde at upper bound accepted", strings.Repeat("a", 11) + "~", false},
		{"DEL rejected", strings.Repeat("a", 11) + "\x7f", true},
		{"unicode rejected", "password-äbc", true},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			_, err := password.Parse(tt.value)
			if tt.wantErr && err == nil {
				t.Errorf("Parse(%q) should fail", tt.value)
			}
			if !tt.wantErr && err != nil {
				t.Errorf("Parse(%q) should succeed: %s", tt.value, err)
			}
		})
	}
}

func Test_ParseConfirm(t *testing.T) {
	t.Parallel()

	if _, err := password.ParseConfirm("valid-password-1", "valid-password-1"); err != nil {
		t.Errorf("matching valid passwords should succeed: %s", err)
	}

	if _, err := password.ParseConfirm("valid-password-1", "other-password-2"); err == nil {
		t.Error("mismatched passwords should fail")
	}

	if _, err := password.ParseConfirm("short", "short"); err == nil {
		t.Error("policy-violating password should fail even when matching")
	}
}
