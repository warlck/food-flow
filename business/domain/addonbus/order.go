package addonbus

import "github.com/warlck/food-flow/business/sdk/order"

// OrderBy field constants for addon queries.
const (
	OrderByID          = "addon_id"
	OrderByMenuItemID  = "menu_item_id"
	OrderByName        = "name"
	OrderByPrice       = "price"
	OrderByRank        = "rank"
	OrderByDateCreated = "date_created"
	OrderByDateUpdated = "date_updated"
)

// DefaultOrderBy represents the default order by value.
var DefaultOrderBy = order.NewBy(OrderByRank, order.ASC)
