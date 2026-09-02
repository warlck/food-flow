package addonbus

import "github.com/warlck/food-flow/business/sdk/order"

// OrderBy field constants for addon queries.
const (
	OrderByID          = "addon_id"
	OrderByName        = "name"
	OrderByPrice       = "price"
	OrderByDateCreated = "date_created"
	OrderByDateUpdated = "date_updated"
)

// DefaultOrderBy represents the default order by value.
var DefaultOrderBy = order.NewBy(OrderByName, order.ASC)
