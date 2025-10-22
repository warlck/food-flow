package restaurantapi_test

import (
	"github.com/warlck/food-flow/app/domain/restaurantapp"
	"github.com/warlck/food-flow/business/domain/restaurantbus"
)

func toAppRestaurant(rest restaurantbus.Restaurant) restaurantapp.Restaurant {
	return restaurantapp.ToAppRestaurant(rest)
}

func toAppRestaurantPtr(rest restaurantbus.Restaurant) *restaurantapp.Restaurant {
	restapp := restaurantapp.ToAppRestaurant(rest)
	return &restapp
}
