package modifieroptiondb

import (
	"fmt"

	"github.com/warlck/food-flow/business/domain/modifieroptionbus"
	"github.com/warlck/food-flow/business/sdk/order"
)

var orderByFields = map[string]string{
	modifieroptionbus.OrderByID:         "modifier_option_id",
	modifieroptionbus.OrderByRank:       "rank",
	modifieroptionbus.OrderByName:       "name",
	modifieroptionbus.OrderByPriceDelta: "price_delta",
	modifieroptionbus.OrderByAvailable:  "available",
}

func orderByClause(orderBy order.By) (string, error) {
	by, exists := orderByFields[orderBy.Field]
	if !exists {
		return "", fmt.Errorf("field %q does not exist", orderBy.Field)
	}

	if orderBy.Field == modifieroptionbus.OrderByRank {
		return fmt.Sprintf(" ORDER BY rank %s NULLS LAST, name ASC, modifier_option_id ASC", orderBy.Direction), nil
	}

	return " ORDER BY " + by + " " + orderBy.Direction, nil
}
