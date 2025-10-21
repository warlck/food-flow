package restaurantdb

import (
	"fmt"

	"github.com/warlck/food-flow/business/domain/restaurantbus"
	"github.com/warlck/food-flow/business/sdk/order"
)

var orderByFields = map[string]string{
	restaurantbus.OrderByID:      "restaurant_id",
	restaurantbus.OrderByName:    "name",
	restaurantbus.OrderByEnabled: "enabled",
}

func orderByClause(orderBy order.By) (string, error) {
	by, exists := orderByFields[orderBy.Field]
	if !exists {
		return "", fmt.Errorf("field %q does not exist", orderBy.Field)
	}

	return " ORDER BY " + by + " " + orderBy.Direction, nil
}
