package menuitemapi_test

import (
	"github.com/warlck/food-flow/app/domain/menuitemapp"
	"github.com/warlck/food-flow/business/domain/menuitembus"
)

func toAppMenuItem(item menuitembus.MenuItem) menuitemapp.MenuItem {
	return menuitemapp.MenuItem{
		ID:           item.ID.String(),
		CategoryID:   item.CategoryID.String(),
		RestaurantID: item.RestaurantID.String(),
		Name:         item.Name.String(),
		Description:  item.Description,
		Price:        item.Price.Value(),
		Available:    item.Available,
		ImageURL:     item.ImageURL,
		DateCreated:  item.DateCreated.Format("2006-01-02T15:04:05Z07:00"),
		DateUpdated:  item.DateUpdated.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func toAppMenuItemPtr(item menuitembus.MenuItem) *menuitemapp.MenuItem {
	appItem := toAppMenuItem(item)
	return &appItem
}

func toAppMenuItems(items []menuitembus.MenuItem) []menuitemapp.MenuItem {
	appItems := make([]menuitemapp.MenuItem, len(items))
	for i, item := range items {
		appItems[i] = toAppMenuItem(item)
	}
	return appItems
}
