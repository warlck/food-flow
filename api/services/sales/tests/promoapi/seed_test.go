package promoapi_test

import (
	"context"
	"fmt"
	"github.com/warlck/food-flow/business/domain/organizationbus"

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
	Orgs       []organizationbus.Organization
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

	// Create regular users
	usrs, err := userbus.TestSeedUsers(ctx, 2, role.User, busDomain.User)
	if err != nil {
		return SeedData{}, fmt.Errorf("seeding users : %w", err)
	}

	// Seed restaurants
	orgs, err := organizationbus.TestSeedOrganizations(ctx, 2, busDomain.Organization)
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

	tu1 := apitest.User{
		User:  adminUsrs[0],
		Token: apitest.Token(db.BusDomain, ath, adminUsrs[0].Email.Address),
	}
	tu2 := apitest.User{
		User:  adminUsrs[1],
		Token: apitest.Token(db.BusDomain, ath, adminUsrs[1].Email.Address),
	}
	tu3 := apitest.User{
		User:  usrs[0],
		Token: apitest.Token(db.BusDomain, ath, usrs[0].Email.Address),
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

	rest0ID := rests[0].ID
	scopedPromo, err := busDomain.Promo.Create(ctx, promobus.NewPromotion{
		Code:           "SCOPED10",
		Name:           name.MustParse("Scoped Promotion"),
		DiscountType:   promobus.DiscountTypePercentage,
		DiscountValue:  10.0,
		MinOrderAmount: 10.0,
		RestaurantID:   &rest0ID,
		Enabled:        true,
	})
	if err != nil {
		return SeedData{}, fmt.Errorf("seeding scoped promo : %w", err)
	}
	promos = append(promos, scopedPromo)

	// Seed a restaurant and promotion belonging to a second organization the
	// admin users are NOT members of, for cross-organization authorization tests.
	otherRests, err := restaurantbus.TestSeedRestaurants(ctx, 1, busDomain.Restaurant, orgs[1].ID)
	if err != nil {
		return SeedData{}, fmt.Errorf("seeding other org restaurants : %w", err)
	}

	otherRestID := otherRests[0].ID
	otherOrgPromo, err := busDomain.Promo.Create(ctx, promobus.NewPromotion{
		Code:           "OTHERORG20",
		Name:           name.MustParse("Other Org Promotion"),
		DiscountType:   promobus.DiscountTypePercentage,
		DiscountValue:  20.0,
		MinOrderAmount: 10.0,
		RestaurantID:   &otherRestID,
		Enabled:        true,
	})
	if err != nil {
		return SeedData{}, fmt.Errorf("seeding other org promo : %w", err)
	}
	promos = append(promos, otherOrgPromo)

	sd := SeedData{
		SeedData: apitest.SeedData{
			Users:  []apitest.User{tu3},
			Admins: []apitest.User{tu1, tu2},
			Restaurants: []apitest.Restaurant{
				{Restaurant: rests[0]},
				{Restaurant: rests[1]},
			},
		},
		Orgs:       orgs,
		Promotions: promos,
	}

	return sd, nil
}
