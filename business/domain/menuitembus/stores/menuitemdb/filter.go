package menuitemdb

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/warlck/food-flow/business/domain/menuitembus"
)

func applyFilter(filter menuitembus.QueryFilter, data map[string]any, buf *bytes.Buffer) {
	var wc []string

	if filter.ID != nil {
		data["menu_item_id"] = *filter.ID
		wc = append(wc, "menu_item_id = :menu_item_id")
	}

	if filter.Name != nil {
		data["name"] = fmt.Sprintf("%%%s%%", *filter.Name)
		wc = append(wc, "name LIKE :name")
	}

	if filter.CategoryID != nil {
		data["category_id"] = *filter.CategoryID
		wc = append(wc, "category_id = :category_id")
	}

	if filter.RestaurantID != nil {
		data["restaurant_id"] = *filter.RestaurantID
		wc = append(wc, "restaurant_id = :restaurant_id")
	}

	if filter.MinPrice != nil {
		data["min_price"] = filter.MinPrice.Value()
		wc = append(wc, "price >= :min_price")
	}

	if filter.MaxPrice != nil {
		data["max_price"] = filter.MaxPrice.Value()
		wc = append(wc, "price <= :max_price")
	}

	if filter.Available != nil {
		data["available"] = *filter.Available
		wc = append(wc, "available = :available")
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
