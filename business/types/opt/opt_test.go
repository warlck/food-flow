package opt_test

import (
	"encoding/json"
	"testing"

	"github.com/warlck/food-flow/business/types/opt"
)

type testStruct struct {
	Rank opt.NullInt `json:"rank"`
}

func Test_NullInt(t *testing.T) {
	t.Parallel()

	t.Run("omitted", func(t *testing.T) {
		var s testStruct
		if err := json.Unmarshal([]byte(`{}`), &s); err != nil {
			t.Fatalf("unexpected unmarshal error: %v", err)
		}
		if s.Rank.Present {
			t.Fatal("expected Rank.Present to be false for omitted field")
		}
		if s.Rank.Value != nil {
			t.Fatal("expected Rank.Value to be nil")
		}
	})

	t.Run("explicit-value", func(t *testing.T) {
		var s testStruct
		if err := json.Unmarshal([]byte(`{"rank": 10}`), &s); err != nil {
			t.Fatalf("unexpected unmarshal error: %v", err)
		}
		if !s.Rank.Present {
			t.Fatal("expected Rank.Present to be true")
		}
		if s.Rank.Value == nil || *s.Rank.Value != 10 {
			t.Fatalf("expected Rank.Value to be 10, got %v", s.Rank.Value)
		}
	})

	t.Run("explicit-null", func(t *testing.T) {
		var s testStruct
		if err := json.Unmarshal([]byte(`{"rank": null}`), &s); err != nil {
			t.Fatalf("unexpected unmarshal error: %v", err)
		}
		if !s.Rank.Present {
			t.Fatal("expected Rank.Present to be true for explicit null")
		}
		if s.Rank.Value != nil {
			t.Fatalf("expected Rank.Value to be nil for explicit null, got %v", s.Rank.Value)
		}
	})

	t.Run("invalid-type", func(t *testing.T) {
		var s testStruct
		if err := json.Unmarshal([]byte(`{"rank": "not_an_int"}`), &s); err == nil {
			t.Fatal("expected error for non-int json value")
		}
	})

	t.Run("equal", func(t *testing.T) {
		n1 := opt.NewNullInt(10)
		n2 := opt.NewNullInt(10)
		n3 := opt.NewNullInt(20)
		nNull := opt.NewNullIntNull()
		nOmit := opt.NullInt{}

		if !n1.Equal(n2) {
			t.Fatal("expected n1 == n2")
		}
		if n1.Equal(n3) {
			t.Fatal("expected n1 != n3")
		}
		if n1.Equal(nNull) {
			t.Fatal("expected n1 != nNull")
		}
		if nNull.Equal(nOmit) {
			t.Fatal("expected nNull != nOmit")
		}
	})
}
