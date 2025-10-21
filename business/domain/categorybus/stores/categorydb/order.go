package categorydb

import (
	"fmt"

	"github.com/warlck/food-flow/business/domain/categorybus"
	"github.com/warlck/food-flow/business/sdk/order"
)

var orderByFields = map[string]string{
	categorybus.OrderByID:      "category_id",
	categorybus.OrderByName:    "name",
	categorybus.OrderByEnabled: "enabled",
}

func orderByClause(orderBy order.By) (string, error) {
	by, exists := orderByFields[orderBy.Field]
	if !exists {
		return "", fmt.Errorf("field %q does not exist", orderBy.Field)
	}

	return " ORDER BY " + by + " " + orderBy.Direction, nil
}
