package orderbus

import "time"

// QueryFilter holds filter criteria for querying orders
type QueryFilter struct {
	ID            *string
	RestaurantID  *string
	CustomerEmail *string
	OrderStatus   *string
	PaymentStatus *string
	OrderType     *string
	StartDate     *time.Time
	EndDate       *time.Time
}
