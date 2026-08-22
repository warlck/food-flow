package menuitemapp

import (
	"github.com/warlck/food-flow/business/domain/menuitembus"
	"github.com/warlck/food-flow/business/sdk/order"
)

var defaultOrderBy = order.NewBy("menu_item_id", order.ASC)

var orderByFields = map[string]string{
	"menu_item_id": menuitembus.OrderByID,
	"name":         menuitembus.OrderByName,
	"price":        menuitembus.OrderByPrice,
	"available":    menuitembus.OrderByAvailable,
	"rank":         menuitembus.OrderByRank,
}
