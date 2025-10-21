package menuitembus

import "github.com/warlck/food-flow/business/sdk/order"

// DefaultOrderBy represents the default way we sort.
var DefaultOrderBy = order.NewBy(OrderByID, order.ASC)

// Set of fields that the results can be ordered by.
const (
	OrderByID        = "menu_item_id"
	OrderByName      = "name"
	OrderByPrice     = "price"
	OrderByAvailable = "available"
)
