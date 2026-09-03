package modifiergroupapp

import (
	"github.com/warlck/food-flow/business/domain/modifiergroupbus"
	"github.com/warlck/food-flow/business/sdk/order"
)

var defaultOrderBy = order.NewBy(modifiergroupbus.OrderByRank, order.ASC)

var orderByFields = map[string]string{
	"id":        modifiergroupbus.OrderByID,
	"rank":      modifiergroupbus.OrderByRank,
	"name":      modifiergroupbus.OrderByName,
	"available": modifiergroupbus.OrderByAvailable,
}
