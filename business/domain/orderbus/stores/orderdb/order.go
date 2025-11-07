package orderdb

import (
	"fmt"

	"github.com/warlck/food-flow/business/domain/orderbus"
	"github.com/warlck/food-flow/business/sdk/order"
)

var orderByFields = map[string]string{
	orderbus.OrderByID:          "order_id",
	orderbus.OrderByDateCreated: "date_created",
	orderbus.OrderByTotal:       "total",
	orderbus.OrderByStatus:      "order_status",
}

func orderByClause(orderBy order.By) (string, error) {
	by, exists := orderByFields[orderBy.Field]
	if !exists {
		return "", fmt.Errorf("field %q does not exist", orderBy.Field)
	}

	return " ORDER BY " + by + " " + orderBy.Direction, nil
}
