package promobus

import "github.com/warlck/food-flow/business/sdk/order"

// DefaultOrderBy represents the default sort order.
var DefaultOrderBy = order.NewBy(OrderByCode, order.ASC)

// Sort field constants.
const (
	OrderByID          = "promotion_id"
	OrderByCode        = "code"
	OrderByName        = "name"
	OrderByDateCreated = "date_created"
)
