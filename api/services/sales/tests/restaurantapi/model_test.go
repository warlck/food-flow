package restaurantapi_test

import (
	"time"

	restaurantapi "github.com/warlck/food-flow/app/domain/restaurantapp"
	"github.com/warlck/food-flow/business/domain/restaurantbus"
)

func toAppRestaurantPtr(rest restaurantbus.Restaurant) *restaurantapi.Restaurant {
	return &restaurantapi.Restaurant{
		ID:          rest.ID.String(),
		Name:        rest.Name.String(),
		Description: rest.Description,
		Address:     rest.Address,
		Phone:       rest.Phone,
		Email:       rest.Email,
		ImageURL:    rest.ImageURL,
		Enabled:     rest.Enabled,
		DateCreated: rest.DateCreated.Format(time.RFC3339),
		DateUpdated: rest.DateUpdated.Format(time.RFC3339),
	}
}
