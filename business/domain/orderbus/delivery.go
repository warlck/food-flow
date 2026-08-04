package orderbus

import (
	"errors"
	"fmt"
	"math"

	"github.com/warlck/food-flow/business/types/money"
)

// Set of error variables for delivery operations.
var (
	// ErrDeliveryCoordinatesRequired is returned when a delivery order is
	// missing destination coordinates.
	ErrDeliveryCoordinatesRequired = errors.New("delivery address coordinates required for delivery orders")

	// ErrRestaurantLocationMissing is returned when the restaurant has no
	// coordinates configured and delivery distance cannot be calculated.
	ErrRestaurantLocationMissing = errors.New("restaurant location not configured")

	// ErrDeliveryOutOfRange is returned when the delivery destination exceeds
	// the restaurant's maximum delivery distance.
	ErrDeliveryOutOfRange = errors.New("delivery destination outside delivery range")
)

const (
	// DeliveryFeePerKm is charged per kilometre beyond the free distance.
	DeliveryFeePerKm = 1.00

	// DeliveryFreeDistanceKm is the distance covered free of charge.
	DeliveryFreeDistanceKm = 1.0
)

// earthRadiusKm is the mean radius of the Earth in kilometres.
const earthRadiusKm = 6371.0

// DistanceKm calculates the great-circle distance between two coordinates
// using the Haversine formula.
func DistanceKm(lat1, lng1, lat2, lng2 float64) float64 {
	rad := math.Pi / 180
	dLat := (lat2 - lat1) * rad
	dLng := (lng2 - lng1) * rad

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*rad)*math.Cos(lat2*rad)*math.Sin(dLng/2)*math.Sin(dLng/2)

	return 2 * earthRadiusKm * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
}

// CalculateDeliveryFee returns the delivery fee for a distance in kilometres.
// The first DeliveryFreeDistanceKm kilometres are free, then DeliveryFeePerKm
// is charged per kilometre.
func CalculateDeliveryFee(distanceKm float64) float64 {
	if distanceKm <= DeliveryFreeDistanceKm {
		return 0
	}

	return roundToTwoDecimals((distanceKm - DeliveryFreeDistanceKm) * DeliveryFeePerKm)
}

// DeliveryQuote describes the delivery fee and range check for a destination.
type DeliveryQuote struct {
	DistanceKm            float64     // Distance between restaurant and destination
	DeliveryFee           money.Money // Fee charged for the distance
	MaxDeliveryDistanceKm float64     // Restaurant limit (0 means unlimited)
	WithinLimit           bool        // Whether the destination can be served
}

// DeliveryQuote calculates the delivery quote given restaurant and destination coordinates.
func (b *Business) DeliveryQuote(restLat, restLng, destLat, destLng float64, maxDeliveryDistanceKm float64) (DeliveryQuote, error) {
	return deliveryQuote(restLat, restLng, destLat, destLng, maxDeliveryDistanceKm)
}

// deliveryQuote calculates the delivery quote given restaurant and destination coordinates.
func deliveryQuote(restLat, restLng, destLat, destLng float64, maxDeliveryDistanceKm float64) (DeliveryQuote, error) {
	distance := DistanceKm(restLat, restLng, destLat, destLng)

	withinLimit := maxDeliveryDistanceKm <= 0 || distance <= maxDeliveryDistanceKm

	fee, err := money.Parse(CalculateDeliveryFee(distance))
	if err != nil {
		return DeliveryQuote{}, fmt.Errorf("parse delivery fee: %w", err)
	}

	return DeliveryQuote{
		DistanceKm:            roundToTwoDecimals(distance),
		DeliveryFee:           fee,
		MaxDeliveryDistanceKm: maxDeliveryDistanceKm,
		WithinLimit:           withinLimit,
	}, nil
}
