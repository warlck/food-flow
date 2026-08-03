package orderbus_test

import (
	"testing"

	"github.com/warlck/food-flow/business/domain/orderbus"
)

func Test_DistanceKm(t *testing.T) {
	t.Parallel()

	table := []struct {
		name           string
		lat1, lng1     float64
		lat2, lng2     float64
		expMin, expMax float64
	}{
		{
			name: "same-point",
			lat1: 1.29305, lng1: 103.86020,
			lat2: 1.29305, lng2: 103.86020,
			expMin: 0, expMax: 0.001,
		},
		{
			name: "one-degree-longitude-at-equator",
			lat1: 0, lng1: 0,
			lat2: 0, lng2: 1,
			expMin: 111.0, expMax: 111.4,
		},
		{
			name: "singapore-cross-town",
			lat1: 1.29305, lng1: 103.86020,
			lat2: 1.35208, lng2: 103.94400,
			expMin: 11.0, expMax: 12.0,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			got := orderbus.DistanceKm(tt.lat1, tt.lng1, tt.lat2, tt.lng2)
			if got < tt.expMin || got > tt.expMax {
				t.Fatalf("distance %.4f km outside expected range [%.4f, %.4f]", got, tt.expMin, tt.expMax)
			}
		})
	}
}

func Test_CalculateDeliveryFee(t *testing.T) {
	t.Parallel()

	table := []struct {
		name       string
		distanceKm float64
		expFee     float64
	}{
		{"zero-distance", 0, 0},
		{"within-free-distance", 0.5, 0},
		{"exactly-free-distance", 1.0, 0},
		{"just-beyond-free-distance", 1.01, 0.01},
		{"three-point-five-km", 3.5, 2.50},
		{"ten-point-seven-seven-km", 10.77, 9.77},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			got := orderbus.CalculateDeliveryFee(tt.distanceKm)
			if got != tt.expFee {
				t.Fatalf("fee for %.2f km: got %.2f, expected %.2f", tt.distanceKm, got, tt.expFee)
			}
		})
	}
}

func Test_DeliveryQuote(t *testing.T) {
	t.Parallel()

	bus := &orderbus.Business{}

	// Test case within delivery limit
	quote, err := bus.DeliveryQuote(1.29305, 103.86020, 1.29305, 103.86020, 10.0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !quote.WithinLimit {
		t.Fatalf("expected within limit to be true")
	}
	if quote.DeliveryFee.Value() != 0 {
		t.Fatalf("expected 0 delivery fee, got %.2f", quote.DeliveryFee.Value())
	}

	// Test case exceeding delivery limit
	quoteExceed, err := bus.DeliveryQuote(1.29305, 103.86020, 1.35208, 103.94400, 5.0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if quoteExceed.WithinLimit {
		t.Fatalf("expected within limit to be false for distance exceeding max limit")
	}
}
