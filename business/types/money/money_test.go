package money_test

import (
	"errors"
	"math"
	"testing"

	"github.com/warlck/food-flow/business/types/money"
)

func Test_Money(t *testing.T) {
	t.Parallel()

	t.Run("valid-values", func(t *testing.T) {
		tests := []float64{0.0, 0.01, 10.50, 999.99, 1_000_000.0, 99_999_999.99}
		for _, val := range tests {
			m, err := money.Parse(val)
			if err != nil {
				t.Fatalf("expected valid money for %.2f, got err: %v", val, err)
			}
			if m.Value() != val {
				t.Fatalf("expected value %.2f, got %.2f", val, m.Value())
			}
		}
	})

	t.Run("negative-value-rejected", func(t *testing.T) {
		_, err := money.Parse(-0.01)
		if err == nil {
			t.Fatal("expected error for negative money, got nil")
		}
	})

	t.Run("overflow-value-rejected", func(t *testing.T) {
		_, err := money.Parse(100_000_000.00)
		if err == nil {
			t.Fatal("expected error for money > 99_999_999.99, got nil")
		}
		if !errors.Is(err, money.ErrOverflow) {
			t.Fatalf("expected money.ErrOverflow, got: %v", err)
		}
	})

	t.Run("non-finite-value-rejected", func(t *testing.T) {
		for _, val := range []float64{math.NaN(), math.Inf(1), math.Inf(-1)} {
			if _, err := money.Parse(val); err == nil {
				t.Fatalf("expected error for non-finite money %v, got nil", val)
			}
		}
	})

	t.Run("must-parse", func(t *testing.T) {
		m := money.MustParse(123.45)
		if m.Value() != 123.45 {
			t.Fatalf("expected 123.45, got %.2f", m.Value())
		}
		if m.String() != "123.45" {
			t.Fatalf("expected '123.45', got '%s'", m.String())
		}
	})

	t.Run("marshal-text", func(t *testing.T) {
		m := money.MustParse(50.0)
		bytes, err := m.MarshalText()
		if err != nil {
			t.Fatalf("unexpected error marshaling text: %v", err)
		}
		if string(bytes) != "50.00" {
			t.Fatalf("expected '50.00', got '%s'", string(bytes))
		}
	})
}
