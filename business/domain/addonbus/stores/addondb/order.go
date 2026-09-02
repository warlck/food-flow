package addondb

import (
	"fmt"

	"github.com/warlck/food-flow/business/domain/addonbus"
	"github.com/warlck/food-flow/business/sdk/order"
)

var orderByFields = map[string]string{
	addonbus.OrderByID:          "addon_id",
	addonbus.OrderByName:        "name",
	addonbus.OrderByPrice:       "price",
	addonbus.OrderByDateCreated: "date_created",
	addonbus.OrderByDateUpdated: "date_updated",
}

func orderByClause(orderBy order.By) (string, error) {
	by, exists := orderByFields[orderBy.Field]
	if !exists {
		return "", fmt.Errorf("field %q does not exist", orderBy.Field)
	}

	return " ORDER BY " + by + " " + orderBy.Direction, nil
}
