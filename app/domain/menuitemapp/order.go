package menuitemapp

import (
	"github.com/warlck/food-flow/business/domain/menuitembus"
	"github.com/warlck/food-flow/business/sdk/order"
)

var defaultOrderBy = order.NewBy(menuitembus.OrderByRank, order.ASC)

var orderByFields = map[string]string{
	"id":           menuitembus.OrderByID,
	"menu_item_id": menuitembus.OrderByID,
	"name":         menuitembus.OrderByName,
	"price":        menuitembus.OrderByPrice,
	"available":    menuitembus.OrderByAvailable,
	"rank":         menuitembus.OrderByRank,
}
