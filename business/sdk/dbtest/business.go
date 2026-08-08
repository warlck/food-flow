package dbtest

import (
	"github.com/jmoiron/sqlx"
	"github.com/warlck/food-flow/business/domain/addonbus"
	"github.com/warlck/food-flow/business/domain/addonbus/stores/addondb"
	"github.com/warlck/food-flow/business/domain/categorybus"
	"github.com/warlck/food-flow/business/domain/categorybus/stores/categorydb"
	"github.com/warlck/food-flow/business/domain/menuitembus"
	"github.com/warlck/food-flow/business/domain/menuitembus/stores/menuitemdb"
	"github.com/warlck/food-flow/business/domain/orderbus"
	"github.com/warlck/food-flow/business/domain/orderbus/stores/orderdb"
	"github.com/warlck/food-flow/business/domain/restaurantbus"
	"github.com/warlck/food-flow/business/domain/restaurantbus/stores/restaurantdb"
	"github.com/warlck/food-flow/business/domain/userbus"
	"github.com/warlck/food-flow/business/domain/userbus/stores/userdb"
	"github.com/warlck/food-flow/foundation/logger"
)

// BusDomain represents all the business domain apis needed for testing.
type BusDomain struct {
	User       *userbus.Business
	Restaurant *restaurantbus.Business
	Category   *categorybus.Business
	MenuItem   *menuitembus.Business
	Order      *orderbus.Business
	Addon      *addonbus.Business
}

func newBusDomains(log *logger.Logger, db *sqlx.DB) BusDomain {

	userStorage := userdb.NewStore(log, db)
	userBus := userbus.NewBusiness(log, userStorage)

	restaurantStorage := restaurantdb.NewStore(log, db)
	restaurantBus := restaurantbus.NewBusiness(log, restaurantStorage)

	categoryStorage := categorydb.NewStore(log, db)
	categoryBus := categorybus.NewBusiness(log, categoryStorage)

	menuItemStorage := menuitemdb.NewStore(log, db)
	menuItemBus := menuitembus.NewBusiness(log, menuItemStorage)

	addonStorage := addondb.NewStore(log, db)
	addonBus := addonbus.NewBusiness(log, addonStorage)

	orderStorage := orderdb.NewStore(log, db)
	orderBus := orderbus.NewBusiness(log, orderStorage, menuItemBus, restaurantBus, addonBus)

	return BusDomain{
		User:       userBus,
		Restaurant: restaurantBus,
		Category:   categoryBus,
		MenuItem:   menuItemBus,
		Order:      orderBus,
		Addon:      addonBus,
	}
}
