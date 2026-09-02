package modifiergroupbus

import "github.com/warlck/food-flow/business/sdk/order"

// DefaultOrderBy represents the default way we sort modifier groups.
var DefaultOrderBy = order.NewBy(OrderByRank, order.ASC)

// Set of fields that the results can be ordered by.
const (
	OrderByID        = "modifier_group_id"
	OrderByRank      = "rank"
	OrderByName      = "name"
	OrderByAvailable = "available"
)
