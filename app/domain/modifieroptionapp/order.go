package modifieroptionapp

import (
	"github.com/warlck/food-flow/business/domain/modifieroptionbus"
	"github.com/warlck/food-flow/business/sdk/order"
)

var defaultOrderBy = order.NewBy(modifieroptionbus.OrderByRank, order.ASC)

var orderByFields = map[string]string{
	"id":          modifieroptionbus.OrderByID,
	"rank":        modifieroptionbus.OrderByRank,
	"name":        modifieroptionbus.OrderByName,
	"price_delta": modifieroptionbus.OrderByPriceDelta,
	"available":   modifieroptionbus.OrderByAvailable,
}
