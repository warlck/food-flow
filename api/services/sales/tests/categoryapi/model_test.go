package categoryapi_test

import (
	"github.com/warlck/food-flow/app/domain/categoryapp"
	"github.com/warlck/food-flow/business/domain/categorybus"
)

func toAppCategoryPtr(cat categorybus.Category) *categoryapp.Category {
	catApp := categoryapp.ToAppCategory(cat)
	return &catApp
}

func toAppCategory(cat categorybus.Category) categoryapp.Category {
	return categoryapp.ToAppCategory(cat)
}
