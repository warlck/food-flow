package modifieroptionbus

import "github.com/warlck/food-flow/business/sdk/order"

// DefaultOrderBy represents the default way we sort modifier options.
var DefaultOrderBy = order.NewBy(OrderByRank, order.ASC)

// Set of fields that the results can be ordered by.
const (
	OrderByID         = "modifier_option_id"
	OrderByRank       = "rank"
	OrderByName       = "name"
	OrderByPriceDelta = "price_delta"
	OrderByAvailable  = "available"
)
