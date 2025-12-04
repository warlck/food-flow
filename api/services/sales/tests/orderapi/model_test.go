package orderapi_test

import (
	"github.com/warlck/food-flow/app/domain/orderapp"
	"github.com/warlck/food-flow/business/domain/orderbus"
)

func toAppOrder(bus orderbus.Order) orderapp.Order {
	return orderapp.ToAppOrder(bus)
}

func toAppOrders(orders []orderbus.Order) []orderapp.Order {
	return orderapp.ToAppOrders(orders)
}

func toAppOrderPtr(bus orderbus.Order) *orderapp.Order {
	appOrd := toAppOrder(bus)
	return &appOrd
}
