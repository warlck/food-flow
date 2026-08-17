package imagebus

import "github.com/warlck/food-flow/business/sdk/order"

// DefaultOrderBy represents the default sort order.
var DefaultOrderBy = order.NewBy(OrderByDateCreated, order.DESC)

// Sort field constants.
const (
	OrderByID          = "image_id"
	OrderByDateCreated = "date_created"
)
