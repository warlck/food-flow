package name_test

import (
	"strings"
	"testing"

	"github.com/warlck/food-flow/business/types/name"
)

func Test_Name(t *testing.T) {
	t.Parallel()

	t.Run("valid-names", func(t *testing.T) {
		validCases := []string{
			"Margherita Pizza",
			"100 Percent",
			"Coca-Cola",
			"Chef's Special",
			"Niño's Tacos",
			"Burger (Large)",
			"Coca-Cola (330ml)",
			"Meal (For 2)",
			"Ramen [Spicy]",
			"[Vegan] Dumplings",
			"Combo {Special}",
			"{Chef's Pick} Burger",
			"Set (For 2) [Mild] {Value}",
			"Tomato, Mozzarella",
			"Dr. Pepper (Large)",
			"St. Louis Ribs",
			"Special: Wagyu Beef",
			"Combo 1.5L: Family Pack, (Cold)",
			"Coffee / Tea",
			"Fish \\ Chips",
			"Combo A | Option B",
			"item_with_underscore",
			"Super-Hot!",
			"Feeling Hungry?",
			"Café au Lait (Hot)",
			"Crème Brûlée",
			"Hot",                    // Min length 3
			strings.Repeat("a", 100), // Max length 100
		}

		for _, tc := range validCases {
			n, err := name.Parse(tc)
			if err != nil {
				t.Fatalf("expected valid name for %q, got err: %v", tc, err)
			}
			if n.String() != tc {
				t.Fatalf("expected %q, got %q", tc, n.String())
			}
		}
	})

	t.Run("invalid-names-rejected", func(t *testing.T) {
		invalidCases := []struct {
			input  string
			reason string
		}{
			{"", "empty string"},
			{"A", "length 1"},
			{"Bi", "length 2"},
			{strings.Repeat("a", 101), "length 101"},
			{"12\" Pizza", "double quote"},
			{"Ben & Jerry's", "ampersand"},
			{"Fish & Chips", "ampersand"},
			{"Soup + Salad", "plus sign"},
			{"100% Beef", "percent sign"},
			{"Combo #1", "hash sign"},
			{"Lunch @ Diner", "at sign"},
			{"~200g Steak", "tilde"},
			{"Chef's Pick*", "asterisk"},
			{"Item; Extra", "semicolon"},
			{"<script>alert(1)</script>", "html script tag"},
			{"Burger <Special>", "angle brackets"},
			{"Item\nName", "newline"},
			{"Item\tName", "tab"},
			{"Item\x00Name", "null byte"},
		}

		for _, tc := range invalidCases {
			_, err := name.Parse(tc.input)
			if err == nil {
				t.Fatalf("expected error for %s (%q), got nil", tc.reason, tc.input)
			}
		}
	})

	t.Run("must-parse", func(t *testing.T) {
		n := name.MustParse("Burger (Large)")
		if n.String() != "Burger (Large)" {
			t.Fatalf("expected 'Burger (Large)', got %q", n.String())
		}

		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected panic for invalid MustParse, got nil")
			}
		}()
		name.MustParse("Burger & Fries")
	})

	t.Run("parse-null", func(t *testing.T) {
		// Empty string is valid null
		nullEmpty, err := name.ParseNull("")
		if err != nil {
			t.Fatalf("unexpected error for empty ParseNull: %v", err)
		}
		if nullEmpty.Valid() {
			t.Fatal("expected Valid() == false for empty ParseNull")
		}
		if nullEmpty.String() != "NULL" {
			t.Fatalf("expected 'NULL', got %q", nullEmpty.String())
		}

		// Valid non-empty string
		nullValid, err := name.ParseNull("Kitchen (Main)")
		if err != nil {
			t.Fatalf("unexpected error for valid ParseNull: %v", err)
		}
		if !nullValid.Valid() {
			t.Fatal("expected Valid() == true for valid ParseNull")
		}
		if nullValid.String() != "Kitchen (Main)" {
			t.Fatalf("expected 'Kitchen (Main)', got %q", nullValid.String())
		}

		// Invalid non-empty string
		_, err = name.ParseNull("Kitchen & Bath")
		if err == nil {
			t.Fatal("expected error for invalid ParseNull, got nil")
		}
	})

	t.Run("must-parse-null", func(t *testing.T) {
		nullEmpty := name.MustParseNull("")
		if nullEmpty.Valid() {
			t.Fatal("expected Valid() == false")
		}

		nullValid := name.MustParseNull("Sales (Global)")
		if !nullValid.Valid() || nullValid.String() != "Sales (Global)" {
			t.Fatalf("expected valid null name, got %q", nullValid.String())
		}

		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected panic for invalid MustParseNull, got nil")
			}
		}()
		name.MustParseNull("ab")
	})

	t.Run("equal-and-marshal", func(t *testing.T) {
		n1 := name.MustParse("Item (Large)")
		n2 := name.MustParse("Item (Large)")
		n3 := name.MustParse("Item (Small)")

		if !n1.Equal(n2) {
			t.Fatal("expected n1.Equal(n2) == true")
		}
		if n1.Equal(n3) {
			t.Fatal("expected n1.Equal(n3) == false")
		}

		b, err := n1.MarshalText()
		if err != nil {
			t.Fatalf("unexpected error marshaling text: %v", err)
		}
		if string(b) != "Item (Large)" {
			t.Fatalf("expected 'Item (Large)', got %q", string(b))
		}

		null1 := name.MustParseNull("Department (IT)")
		null2 := name.MustParseNull("Department (IT)")
		null3 := name.MustParseNull("")

		if !null1.Equal(null2) {
			t.Fatal("expected null1.Equal(null2) == true")
		}
		if null1.Equal(null3) {
			t.Fatal("expected null1.Equal(null3) == false")
		}

		bNull, err := null1.MarshalText()
		if err != nil {
			t.Fatalf("unexpected error marshaling text: %v", err)
		}
		if string(bNull) != "Department (IT)" {
			t.Fatalf("expected 'Department (IT)', got %q", string(bNull))
		}
	})
}
