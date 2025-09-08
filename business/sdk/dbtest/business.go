package dbtest

import (
	"github.com/jmoiron/sqlx"
	"github.com/warlck/food-flow/business/domain/userbus"
	"github.com/warlck/food-flow/business/domain/userbus/stores/userdb"
	"github.com/warlck/food-flow/foundation/logger"
)

// BusDomain represents all the business domain apis needed for testing.
type BusDomain struct {
	User *userbus.Business
}

func newBusDomains(log *logger.Logger, db *sqlx.DB) BusDomain {

	userStorage := userdb.NewStore(log, db)
	userBus := userbus.NewBusiness(log, userStorage)

	return BusDomain{
		User: userBus,
	}
}
