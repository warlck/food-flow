package modifieroptiondb

import (
	"bytes"
	"strings"

	"github.com/warlck/food-flow/business/domain/modifieroptionbus"
)

func applyFilter(filter modifieroptionbus.QueryFilter, data map[string]any, buf *bytes.Buffer) {
	var wc []string

	if filter.ID != nil {
		data["modifier_option_id"] = *filter.ID
		wc = append(wc, "modifier_option_id = :modifier_option_id")
	}

	if filter.ModifierGroupID != nil {
		data["modifier_group_id"] = *filter.ModifierGroupID
		wc = append(wc, "modifier_group_id = :modifier_group_id")
	}

	if filter.RestaurantID != nil {
		data["restaurant_id"] = *filter.RestaurantID
		wc = append(wc, "restaurant_id = :restaurant_id")
	}

	if filter.Name != nil {
		data["name"] = "%" + filter.Name.String() + "%"
		wc = append(wc, "name ILIKE :name")
	}

	if filter.Available != nil {
		data["available"] = *filter.Available
		wc = append(wc, "available = :available")
	}

	if len(wc) > 0 {
		buf.WriteString(" WHERE ")
		buf.WriteString(strings.Join(wc, " AND "))
	}
}
