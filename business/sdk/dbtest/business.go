package dbtest

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/warlck/food-flow/business/domain/addonbus"
	"github.com/warlck/food-flow/business/domain/addonbus/stores/addondb"
	"github.com/warlck/food-flow/business/domain/categorybus"
	"github.com/warlck/food-flow/business/domain/categorybus/stores/categorydb"
	"github.com/warlck/food-flow/business/domain/imagebus"
	"github.com/warlck/food-flow/business/domain/imagebus/stores/imagedb"
	"github.com/warlck/food-flow/business/domain/menuitembus"
	"github.com/warlck/food-flow/business/domain/menuitembus/stores/menuitemdb"
	"github.com/warlck/food-flow/business/domain/modifiergroupbus"
	"github.com/warlck/food-flow/business/domain/modifiergroupbus/stores/modifiergroupdb"
	"github.com/warlck/food-flow/business/domain/modifieroptionbus"
	"github.com/warlck/food-flow/business/domain/modifieroptionbus/stores/modifieroptiondb"
	"github.com/warlck/food-flow/business/domain/orderbus"
	"github.com/warlck/food-flow/business/domain/orderbus/stores/orderdb"
	"github.com/warlck/food-flow/business/domain/organizationbus"
	"github.com/warlck/food-flow/business/domain/organizationbus/stores/organizationdb"
	"github.com/warlck/food-flow/business/domain/promobus"
	"github.com/warlck/food-flow/business/domain/promobus/stores/promodb"
	"github.com/warlck/food-flow/business/domain/restaurantbus"
	"github.com/warlck/food-flow/business/domain/restaurantbus/stores/restaurantdb"
	"github.com/warlck/food-flow/business/domain/userbus"
	"github.com/warlck/food-flow/business/domain/userbus/stores/userdb"
	"github.com/warlck/food-flow/foundation/logger"
	"github.com/warlck/food-flow/foundation/storage"
)

// BusDomain represents all the business domain apis needed for testing.
type BusDomain struct {
	User            *userbus.Business
	Organization    *organizationbus.Business
	Restaurant      *restaurantbus.Business
	Category        *categorybus.Business
	MenuItem        *menuitembus.Business
	ModifierGroup   *modifiergroupbus.Business
	ModifierOption  *modifieroptionbus.Business
	Order           *orderbus.Business
	Addon           *addonbus.Business
	Promo           *promobus.Business
	Image           *imagebus.Business
	ImageLocalStore storage.LocalStore
}

func newBusDomains(log *logger.Logger, db *sqlx.DB) (BusDomain, error) {

	userStorage := userdb.NewStore(log, db)
	userBus := userbus.NewBusiness(log, userStorage)

	orgStorage := organizationdb.NewStore(log, db)
	orgBus := organizationbus.NewBusiness(log, orgStorage, userBus)

	restaurantStorage := restaurantdb.NewStore(log, db)
	restaurantBus := restaurantbus.NewBusiness(log, restaurantStorage)

	categoryStorage := categorydb.NewStore(log, db)
	categoryBus := categorybus.NewBusiness(log, categoryStorage)

	menuItemStorage := menuitemdb.NewStore(log, db)
	menuItemBus := menuitembus.NewBusiness(log, menuItemStorage)

	modifierGroupStorage := modifiergroupdb.NewStore(log, db)
	modifierOptionStorage := modifieroptiondb.NewStore(log, db)

	modifierGroupBus := modifiergroupbus.NewBusiness(log, modifierGroupStorage, modifierOptionStorage)
	modifierOptionBus := modifieroptionbus.NewBusiness(log, modifierOptionStorage, modifierGroupStorage)

	addonStorage := addondb.NewStore(log, db)
	addonBus := addonbus.NewBusiness(log, addonStorage)

	promoStorage := promodb.NewStore(log, db)
	promoBus := promobus.NewBusiness(log, promoStorage)

	orderStorage := orderdb.NewStore(log, db)
	orderBus := orderbus.NewBusiness(log, orderStorage, menuItemBus, restaurantBus, addonBus, promoBus, categoryBus, modifierGroupBus, modifierOptionBus)

	imageStorage := imagedb.NewStore(log, db)
	imageDir, err := os.MkdirTemp("", "ff-image-test-*")
	if err != nil {
		return BusDomain{}, fmt.Errorf("creating image test dir: %w", err)
	}
	imageSigner, err := storage.NewSigner(context.Background(), storage.Config{
		Backend:      storage.BackendLocal,
		LocalDir:     imageDir,
		LocalBaseURL: "/v1/images/local",
		URLTTL:       15 * time.Minute,
	})
	if err != nil {
		return BusDomain{}, fmt.Errorf("creating image signer: %w", err)
	}
	imageBus := imagebus.NewBusiness(log, imageStorage, imageSigner, 0)

	imageLocalStore, _ := imageSigner.(storage.LocalStore)

	return BusDomain{
		User:            userBus,
		Organization:    orgBus,
		Restaurant:      restaurantBus,
		Category:        categoryBus,
		MenuItem:        menuItemBus,
		ModifierGroup:   modifierGroupBus,
		ModifierOption:  modifierOptionBus,
		Order:           orderBus,
		Addon:           addonBus,
		Promo:           promoBus,
		Image:           imageBus,
		ImageLocalStore: imageLocalStore,
	}, nil
}
