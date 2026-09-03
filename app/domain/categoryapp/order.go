package categoryapp

import (
	"github.com/warlck/food-flow/business/domain/categorybus"
	"github.com/warlck/food-flow/business/sdk/order"
)

var defaultOrderBy = order.NewBy(categorybus.OrderByRank, order.ASC)

var orderByFields = map[string]string{
	"id":          categorybus.OrderByID,
	"category_id": categorybus.OrderByID,
	"name":        categorybus.OrderByName,
	"enabled":     categorybus.OrderByEnabled,
	"rank":        categorybus.OrderByRank,
}
