package modifiergroupdb

import (
	"bytes"
	"strings"

	"github.com/warlck/food-flow/business/domain/modifiergroupbus"
)

func applyFilter(filter modifiergroupbus.QueryFilter, data map[string]any, buf *bytes.Buffer) {
	var wc []string

	if filter.ID != nil {
		data["modifier_group_id"] = *filter.ID
		wc = append(wc, "modifier_group_id = :modifier_group_id")
	}

	if filter.MenuItemID != nil {
		data["menu_item_id"] = *filter.MenuItemID
		wc = append(wc, "menu_item_id = :menu_item_id")
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
