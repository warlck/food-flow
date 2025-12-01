package orderbus

import "github.com/warlck/food-flow/business/sdk/order"

// Set of fields that the results can be ordered by.
const (
	OrderByID          = "order_id"
	OrderByDateCreated = "date_created"
	OrderByTotal       = "total"
	OrderByStatus      = "order_status"
)

var DefaultOrderBy = order.NewBy(OrderByDateCreated, order.DESC)
