package menuitemdb

import (
	"fmt"

	"github.com/warlck/food-flow/business/domain/menuitembus"
	"github.com/warlck/food-flow/business/sdk/order"
)

var orderByFields = map[string]string{
	menuitembus.OrderByID:        "menu_item_id",
	menuitembus.OrderByName:      "name",
	menuitembus.OrderByPrice:     "price",
	menuitembus.OrderByAvailable: "available",
}

func orderByClause(orderBy order.By) (string, error) {
	by, exists := orderByFields[orderBy.Field]
	if !exists {
		return "", fmt.Errorf("field %q does not exist", orderBy.Field)
	}

	return " ORDER BY " + by + " " + orderBy.Direction, nil
}
