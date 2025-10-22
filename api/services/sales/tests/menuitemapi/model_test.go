package menuitemapi_test

import (
	"github.com/warlck/food-flow/app/domain/menuitemapp"
	"github.com/warlck/food-flow/business/domain/menuitembus"
)

func toAppMenuItem(item menuitembus.MenuItem) menuitemapp.MenuItem {
	return menuitemapp.ToAppMenuItem(item)
}

func toAppMenuItemPtr(item menuitembus.MenuItem) *menuitemapp.MenuItem {
	appItem := toAppMenuItem(item)
	return &appItem
}

func toAppMenuItems(items []menuitembus.MenuItem) []menuitemapp.MenuItem {
	appItems := menuitemapp.ToAppMenuItems(items)
	return appItems
}
