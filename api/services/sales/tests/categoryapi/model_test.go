package categoryapi_test

import (
	"time"

	categoryapi "github.com/warlck/food-flow/app/domain/categoryapp"
	"github.com/warlck/food-flow/business/domain/categorybus"
)

func toAppCategoryPtr(cat categorybus.Category) *categoryapi.Category {
	return &categoryapi.Category{
		ID:           cat.ID.String(),
		RestaurantID: cat.RestaurantID.String(),
		Name:         cat.Name.String(),
		Description:  cat.Description,
		Enabled:      cat.Enabled,
		DateCreated:  cat.DateCreated.Format(time.RFC3339),
		DateUpdated:  cat.DateUpdated.Format(time.RFC3339),
	}
}
