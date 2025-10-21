package categorydb

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/warlck/food-flow/business/domain/categorybus"
)

func applyFilter(filter categorybus.QueryFilter, data map[string]any, buf *bytes.Buffer) {
	var wc []string

	if filter.ID != nil {
		data["category_id"] = *filter.ID
		wc = append(wc, "category_id = :category_id")
	}

	if filter.Name != nil {
		data["name"] = fmt.Sprintf("%%%s%%", *filter.Name)
		wc = append(wc, "name LIKE :name")
	}

	if filter.RestaurantID != nil {
		data["restaurant_id"] = *filter.RestaurantID
		wc = append(wc, "restaurant_id = :restaurant_id")
	}

	if filter.Enabled != nil {
		data["enabled"] = *filter.Enabled
		wc = append(wc, "enabled = :enabled")
	}

	if filter.StartCreatedDate != nil {
		data["start_date_created"] = *filter.StartCreatedDate
		wc = append(wc, "date_created >= :start_date_created")
	}

	if filter.EndCreatedDate != nil {
		data["end_date_created"] = *filter.EndCreatedDate
		wc = append(wc, "date_created <= :end_date_created")
	}

	if len(wc) > 0 {
		buf.WriteString(" WHERE ")
		buf.WriteString(strings.Join(wc, " AND "))
	}
}
