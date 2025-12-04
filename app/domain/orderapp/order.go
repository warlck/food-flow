package orderapp

import (
	"github.com/warlck/food-flow/business/domain/orderbus"
	"github.com/warlck/food-flow/business/sdk/order"
)

var defaultOrderBy = order.NewBy("date_created", order.DESC)

var orderByFields = map[string]string{
	"order_id": orderbus.OrderByID,
	"status":   orderbus.OrderByStatus,
	"total":    orderbus.OrderByTotal,
	"date":     orderbus.OrderByDateCreated,
}
