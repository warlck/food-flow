package modifiergroupdb

import (
	"fmt"

	"github.com/warlck/food-flow/business/domain/modifiergroupbus"
	"github.com/warlck/food-flow/business/sdk/order"
)

var orderByFields = map[string]string{
	modifiergroupbus.OrderByID:        "modifier_group_id",
	modifiergroupbus.OrderByRank:      "rank",
	modifiergroupbus.OrderByName:      "name",
	modifiergroupbus.OrderByAvailable: "available",
}

func orderByClause(orderBy order.By) (string, error) {
	by, exists := orderByFields[orderBy.Field]
	if !exists {
		return "", fmt.Errorf("field %q does not exist", orderBy.Field)
	}

	if orderBy.Field == modifiergroupbus.OrderByRank {
		return fmt.Sprintf(" ORDER BY rank %s NULLS LAST, name ASC, modifier_group_id ASC", orderBy.Direction), nil
	}

	return " ORDER BY " + by + " " + orderBy.Direction, nil
}
