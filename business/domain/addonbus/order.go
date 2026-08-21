package addonbus

import "github.com/warlck/food-flow/business/sdk/order"

// OrderBy field constants for addon queries.
const (
	OrderByID      = "addon_id"
	OrderByName    = "name"
	OrderByPrice   = "price"
	OrderByRank    = "rank"
	OrderByCreated = "date_created"
)

// DefaultOrderBy represents the default order by value.
var DefaultOrderBy = order.NewBy(OrderByName, order.ASC)
