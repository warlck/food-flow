package restaurantapp

import (
	"github.com/warlck/food-flow/business/domain/restaurantbus"
	"github.com/warlck/food-flow/business/sdk/order"
)

var defaultOrderBy = order.NewBy("restaurant_id", order.ASC)

var orderByFields = map[string]string{
	"restaurant_id": restaurantbus.OrderByID,
	"name":          restaurantbus.OrderByName,
	"enabled":       restaurantbus.OrderByEnabled,
}
