package addonapp

import (
	"github.com/warlck/food-flow/business/domain/addonbus"
)

var defaultOrderBy = addonbus.DefaultOrderBy

var orderByFields = map[string]string{
	"addon_id":     addonbus.OrderByID,
	"name":         addonbus.OrderByName,
	"price":        addonbus.OrderByPrice,
	"date_created": addonbus.OrderByCreated,
	"rank":         addonbus.OrderByRank,
}
