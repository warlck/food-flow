package addonapp

import (
	"github.com/warlck/food-flow/business/domain/addonbus"
)

var defaultOrderBy = addonbus.DefaultOrderBy

var orderByFields = map[string]string{
	"addon_id":     addonbus.OrderByID,
	"menu_item_id": addonbus.OrderByMenuItemID,
	"name":         addonbus.OrderByName,
	"price":        addonbus.OrderByPrice,
	"rank":         addonbus.OrderByRank,
	"dateCreated":  addonbus.OrderByDateCreated,
	"dateUpdated":  addonbus.OrderByDateUpdated,
}
