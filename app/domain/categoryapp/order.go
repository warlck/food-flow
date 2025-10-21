package categoryapp

import (
	"github.com/warlck/food-flow/business/domain/categorybus"
	"github.com/warlck/food-flow/business/sdk/order"
)

var defaultOrderBy = order.NewBy("category_id", order.ASC)

var orderByFields = map[string]string{
	"category_id": categorybus.OrderByID,
	"name":        categorybus.OrderByName,
	"enabled":     categorybus.OrderByEnabled,
}
