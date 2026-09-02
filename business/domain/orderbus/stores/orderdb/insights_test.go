package orderdb

import (
	"errors"
	"math"
	"strings"
	"testing"

	"github.com/warlck/food-flow/business/types/money"
)

// TestToMoney verifies the aggregate-to-money conversion never clamps:
// valid aggregates round to two decimals while negative, non-finite, and
// above-cap aggregates return typed errors the caller must propagate.
func TestToMoney(t *testing.T) {
	tests := []struct {
		name     string
		val      float64
		expMoney money.Money
		expErr   error
	}{
		{
			name:     "normal-value",
			val:      12.34,
			expMoney: money.MustParse(12.34),
		},
		{
			name:     "rounds-to-two-decimals-half-up",
			val:      12.345,
			expMoney: money.MustParse(12.35),
		},
		{
			name:     "zero-value",
			val:      0,
			expMoney: money.MustParse(0),
		},
		{
			name:     "zero-from-empty-aggregate",
			val:      0.0001,
			expMoney: money.MustParse(0),
		},
		{
			name:   "negative-aggregate-rejected",
			val:    -0.01,
			expErr: errors.New("invalid money -0.01"),
		},
		{
			name:   "above-cap-aggregate-rejected",
			val:    100_000_000.00,
			expErr: money.ErrOverflow,
		},
		{
			name:   "nan-aggregate-rejected",
			val:    math.NaN(),
			expErr: errors.New("value must be finite"),
		},
		{
			name:   "positive-infinity-aggregate-rejected",
			val:    math.Inf(1),
			expErr: errors.New("value must be finite"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := toMoney(tc.val)
			if tc.expErr != nil {
				if err == nil {
					t.Fatalf("expected error %v, got money %v", tc.expErr, got)
				}
				if !strings.Contains(err.Error(), tc.expErr.Error()) {
					t.Fatalf("expected error containing %q, got %q", tc.expErr.Error(), err.Error())
				}
				if tc.expErr == money.ErrOverflow && !errors.Is(err, money.ErrOverflow) {
					t.Fatalf("expected money.ErrOverflow in error chain, got %v", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !got.Equal(tc.expMoney) {
				t.Fatalf("expected money %v, got %v", tc.expMoney, got)
			}
		})
	}
}
