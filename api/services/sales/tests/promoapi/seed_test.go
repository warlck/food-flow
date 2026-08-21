package promoapi_test

import (
	"github.com/warlck/food-flow/business/domain/organizationbus"
	"context"
	"fmt"

	"github.com/warlck/food-flow/app/sdk/apitest"
	"github.com/warlck/food-flow/app/sdk/auth"
	"github.com/warlck/food-flow/business/domain/promobus"
	"github.com/warlck/food-flow/business/domain/restaurantbus"
	"github.com/warlck/food-flow/business/domain/userbus"
	"github.com/warlck/food-flow/business/sdk/dbtest"
	"github.com/warlck/food-flow/business/types/name"
	"github.com/warlck/food-flow/business/types/role"
)

type SeedData struct {
	apitest.SeedData
	Promotions []promobus.Promotion
}

func insertSeedData(db *dbtest.Database, ath *auth.Auth) (SeedData, error) {
	ctx := context.Background()
	busDomain := db.BusDomain

	// Create admin users
	adminUsrs, err := userbus.TestSeedUsers(ctx, 2, role.Admin, busDomain.User)
	if err != nil {
		return SeedData{}, fmt.Errorf("seeding admin users : %w", err)
	}

	tu1 := apitest.User{
		User:  adminUsrs[0],
		Token: apitest.Token(db.BusDomain, ath, adminUsrs[0].Email.Address),
	}
	tu2 := apitest.User{
		User:  adminUsrs[1],
		Token: apitest.Token(db.BusDomain, ath, adminUsrs[1].Email.Address),
	}

	// Create regular users
	usrs, err := userbus.TestSeedUsers(ctx, 2, role.User, busDomain.User)
	if err != nil {
		return SeedData{}, fmt.Errorf("seeding users : %w", err)
	}

	tu3 := apitest.User{
		User:  usrs[0],
		Token: apitest.Token(db.BusDomain, ath, usrs[0].Email.Address),
	}

	// Seed restaurants
	orgs, err := organizationbus.TestSeedOrganizations(ctx, 1, busDomain.Organization)
	if err != nil {
		return SeedData{}, fmt.Errorf("seeding organizations: %w", err)
	}

	if _, err := busDomain.Organization.AddUser(ctx, organizationbus.NewOrganizationUser{
		OrganizationID: orgs[0].ID,
		UserID:         adminUsrs[0].ID,
		Role:           role.Admin,
	}); err != nil {
		return SeedData{}, fmt.Errorf("adding user to organization: %w", err)
	}
	if _, err := busDomain.Organization.AddUser(ctx, organizationbus.NewOrganizationUser{
		OrganizationID: orgs[0].ID,
		UserID:         adminUsrs[1].ID,
		Role:           role.Admin,
	}); err != nil {
		return SeedData{}, fmt.Errorf("adding user to organization: %w", err)
	}

	rests, err := restaurantbus.TestSeedRestaurants(ctx, 2, busDomain.Restaurant, orgs[0].ID)
	if err != nil {
		return SeedData{}, fmt.Errorf("seeding restaurants : %w", err)
	}

	// Seed promotions
	promos, err := promobus.TestSeedPromotions(ctx, 3, busDomain.Promo)
	if err != nil {
		return SeedData{}, fmt.Errorf("seeding promotions : %w", err)
	}

	// Additional unhappy-path test promotions
	disabledPromo, err := busDomain.Promo.Create(ctx, promobus.NewPromotion{
		Code:           "DISABLED10",
		Name:           name.MustParse("Disabled Promotion"),
		DiscountType:   promobus.DiscountTypePercentage,
		DiscountValue:  10.0,
		MinOrderAmount: 10.0,
		Enabled:        false,
	})
	if err != nil {
		return SeedData{}, fmt.Errorf("seeding disabled promo : %w", err)
	}
	promos = append(promos, disabledPromo)

	rest1ID := rests[1].ID
	restPromo, err := busDomain.Promo.Create(ctx, promobus.NewPromotion{
		Code:           "REST2ONLY",
		Name:           name.MustParse("Restaurant 2 Only Promotion"),
		DiscountType:   promobus.DiscountTypePercentage,
		DiscountValue:  15.0,
		MinOrderAmount: 10.0,
		RestaurantID:   &rest1ID,
		Enabled:        true,
	})
	if err != nil {
		return SeedData{}, fmt.Errorf("seeding rest promo : %w", err)
	}
	promos = append(promos, restPromo)

	sd := SeedData{
		SeedData: apitest.SeedData{
			Users:  []apitest.User{tu3},
			Admins: []apitest.User{tu1, tu2},
			Restaurants: []apitest.Restaurant{
				{Restaurant: rests[0]},
				{Restaurant: rests[1]},
			},
		},
		Promotions: promos,
	}

	return sd, nil
}
