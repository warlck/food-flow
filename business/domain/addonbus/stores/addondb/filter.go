package addondb

import (
	"bytes"
	"strings"

	"github.com/warlck/food-flow/business/domain/addonbus"
)

func applyFilter(filter addonbus.QueryFilter, data map[string]any, buf *bytes.Buffer) {
	var wc []string

	if filter.ID != nil {
		data["addon_id"] = *filter.ID
		wc = append(wc, "addon_id = :addon_id")
	}

	if filter.CategoryID != nil {
		data["category_id"] = *filter.CategoryID
		wc = append(wc, "category_id = :category_id")
	}

	if filter.RestaurantID != nil {
		data["restaurant_id"] = *filter.RestaurantID
		wc = append(wc, "restaurant_id = :restaurant_id")
	}

	if filter.Name != nil {
		data["name"] = "%" + *filter.Name + "%"
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
