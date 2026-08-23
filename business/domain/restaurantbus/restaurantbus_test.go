package restaurantbus_test

import (
	"context"
	"fmt"
	"github.com/warlck/food-flow/business/domain/organizationbus"
	"sort"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/warlck/food-flow/business/domain/restaurantbus"
	"github.com/warlck/food-flow/business/sdk/dbtest"
	"github.com/warlck/food-flow/business/sdk/page"
	"github.com/warlck/food-flow/business/sdk/unittest"
	"github.com/warlck/food-flow/business/types/name"
)

func Test_Restaurant(t *testing.T) {
	t.Parallel()

	db := dbtest.New(t, "Test_Restaurant")

	sd, err := insertSeedData(db.BusDomain)
	if err != nil {
		t.Fatalf("Seeding error: %s", err)
	}

	// -------------------------------------------------------------------------

	unittest.Run(t, query(db.BusDomain, sd), "query")
	unittest.Run(t, create(db.BusDomain, sd), "create")
	unittest.Run(t, update(db.BusDomain, sd), "update")
	unittest.Run(t, delete(db.BusDomain, sd), "delete")
}

// =============================================================================

func insertSeedData(busDomain dbtest.BusDomain) (unittest.SeedData, error) {
	ctx := context.Background()

	orgs, err := organizationbus.TestSeedOrganizations(ctx, 1, busDomain.Organization)
	if err != nil {
		return unittest.SeedData{}, fmt.Errorf("seeding organizations: %w", err)
	}
	rests, err := restaurantbus.TestSeedRestaurants(ctx, 4, busDomain.Restaurant, orgs[0].ID)
	if err != nil {
		return unittest.SeedData{}, fmt.Errorf("seeding restaurants : %w", err)
	}

	tr1 := unittest.Restaurant{
		Restaurant: rests[0],
	}

	tr2 := unittest.Restaurant{
		Restaurant: rests[1],
	}

	tr3 := unittest.Restaurant{
		Restaurant: rests[2],
	}

	tr4 := unittest.Restaurant{
		Restaurant: rests[3],
	}

	// -------------------------------------------------------------------------

	sd := unittest.SeedData{
		Restaurants: []unittest.Restaurant{tr1, tr2, tr3, tr4},
	}

	return sd, nil
}

// =============================================================================

func query(busDomain dbtest.BusDomain, sd unittest.SeedData) []unittest.Table {
	rests := make([]restaurantbus.Restaurant, 0, len(sd.Restaurants))

	for _, rest := range sd.Restaurants {
		rests = append(rests, rest.Restaurant)
	}

	sort.Slice(rests, func(i, j int) bool {
		return rests[i].ID.String() <= rests[j].ID.String()
	})

	table := []unittest.Table{
		{
			Name:    "all",
			ExpResp: rests,
			ExcFunc: func(ctx context.Context) any {
				filter := restaurantbus.QueryFilter{
					Name: dbtest.NamePointer("Rest"),
				}

				resp, err := busDomain.Restaurant.Query(ctx, filter, restaurantbus.DefaultOrderBy, page.MustParse("1", "10"))
				if err != nil {
					return err
				}

				return resp
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.([]restaurantbus.Restaurant)
				if !exists {
					return "error occurred"
				}

				expResp := exp.([]restaurantbus.Restaurant)

				for i := range gotResp {
					if gotResp[i].DateCreated.Format(time.RFC3339) == expResp[i].DateCreated.Format(time.RFC3339) {
						expResp[i].DateCreated = gotResp[i].DateCreated
					}

					if gotResp[i].DateUpdated.Format(time.RFC3339) == expResp[i].DateUpdated.Format(time.RFC3339) {
						expResp[i].DateUpdated = gotResp[i].DateUpdated
					}
				}

				return cmp.Diff(gotResp, expResp)
			},
		},
		{
			Name:    "byid",
			ExpResp: sd.Restaurants[0].Restaurant,
			ExcFunc: func(ctx context.Context) any {
				resp, err := busDomain.Restaurant.QueryByID(ctx, sd.Restaurants[0].ID)
				if err != nil {
					return err
				}

				return resp
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.(restaurantbus.Restaurant)
				if !exists {
					return "error occurred"
				}

				expResp := exp.(restaurantbus.Restaurant)

				if gotResp.DateCreated.Format(time.RFC3339) == expResp.DateCreated.Format(time.RFC3339) {
					expResp.DateCreated = gotResp.DateCreated
				}

				if gotResp.DateUpdated.Format(time.RFC3339) == expResp.DateUpdated.Format(time.RFC3339) {
					expResp.DateUpdated = gotResp.DateUpdated
				}

				return cmp.Diff(gotResp, expResp)
			},
		},
	}

	return table
}

func create(busDomain dbtest.BusDomain, sd unittest.SeedData) []unittest.Table {
	table := []unittest.Table{
		{
			Name: "basic",
			ExpResp: restaurantbus.Restaurant{
				OrganizationID: sd.Restaurants[0].OrganizationID,
				Name:           name.MustParse("The Italian Place"),
				Description:    "Authentic Italian cuisine with traditional recipes",
				Address:        "123 Pasta Lane",
				Phone:          "+1-555-1234",
				Email:          "info@italianplace.com",
				ImageURL:       "italian.jpg",
				LogoURL:        "italian_logo.jpg",
				OperatingHours: restaurantbus.DefaultOperatingHours(),
				Enabled:        true,
			},
			ExcFunc: func(ctx context.Context) any {
				nr := restaurantbus.NewRestaurant{
					OrganizationID: sd.Restaurants[0].OrganizationID,
					Name:           name.MustParse("The Italian Place"),
					Description:    "Authentic Italian cuisine with traditional recipes",
					Address:        "123 Pasta Lane",
					Phone:          "+1-555-1234",
					Email:          "info@italianplace.com",
					ImageURL:       "italian.jpg",
					LogoURL:        "italian_logo.jpg",
					OperatingHours: restaurantbus.DefaultOperatingHours(),
				}

				resp, err := busDomain.Restaurant.Create(ctx, nr)
				if err != nil {
					return err
				}

				return resp
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.(restaurantbus.Restaurant)
				if !exists {
					return "error occurred"
				}

				expResp := exp.(restaurantbus.Restaurant)

				expResp.ID = gotResp.ID
				expResp.DateCreated = gotResp.DateCreated
				expResp.DateUpdated = gotResp.DateUpdated

				return cmp.Diff(gotResp, expResp)
			},
		},
	}

	return table
}

func update(busDomain dbtest.BusDomain, sd unittest.SeedData) []unittest.Table {
	table := []unittest.Table{
		{
			Name: "basic",
			ExpResp: restaurantbus.Restaurant{
				ID:                    sd.Restaurants[0].ID,
				OrganizationID:        sd.Restaurants[0].OrganizationID,
				Name:                  name.MustParse("Updated Rest Name"),
				Description:           "Updated description for this restaurant",
				Address:               "456 New Street",
				Phone:                 "+1-555-9999",
				Email:                 "updated@example.com",
				ImageURL:              "updated.jpg",
				LogoURL:               "updated_logo.jpg",
				OperatingHours:        sd.Restaurants[0].OperatingHours,
				Enabled:               false,
				Latitude:              sd.Restaurants[0].Latitude,
				Longitude:             sd.Restaurants[0].Longitude,
				MaxDeliveryDistanceKm: sd.Restaurants[0].MaxDeliveryDistanceKm,
				TaxRate:               sd.Restaurants[0].TaxRate,
				DateCreated:           sd.Restaurants[0].DateCreated,
			},
			ExcFunc: func(ctx context.Context) any {
				ur := restaurantbus.UpdateRestaurant{
					Name:        ptr(name.MustParse("Updated Rest Name")),
					Description: ptr("Updated description for this restaurant"),
					Address:     ptr("456 New Street"),
					Phone:       ptr("+1-555-9999"),
					Email:       ptr("updated@example.com"),
					ImageURL:    ptr("updated.jpg"),
					LogoURL:     ptr("updated_logo.jpg"),
					Enabled:     ptr(false),
					MinSpend:    ptr(25.00),
				}

				resp, err := busDomain.Restaurant.Update(ctx, sd.Restaurants[0].Restaurant, ur)
				if err != nil {
					return err
				}

				return resp
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.(restaurantbus.Restaurant)
				if !exists {
					return "error occurred"
				}

				expResp := exp.(restaurantbus.Restaurant)
				expResp.MinSpend = 25.00
				expResp.DateUpdated = gotResp.DateUpdated

				return cmp.Diff(gotResp, expResp)
			},
		},
		{
			Name: "min-spend-update",
			ExpResp: restaurantbus.Restaurant{
				ID:                    sd.Restaurants[1].ID,
				OrganizationID:        sd.Restaurants[1].OrganizationID,
				Name:                  sd.Restaurants[1].Name,
				Description:           sd.Restaurants[1].Description,
				Address:               sd.Restaurants[1].Address,
				Phone:                 sd.Restaurants[1].Phone,
				Email:                 sd.Restaurants[1].Email,
				ImageURL:              sd.Restaurants[1].ImageURL,
				LogoURL:               sd.Restaurants[1].LogoURL,
				OperatingHours:        sd.Restaurants[1].OperatingHours,
				Enabled:               sd.Restaurants[1].Enabled,
				Latitude:              sd.Restaurants[1].Latitude,
				Longitude:             sd.Restaurants[1].Longitude,
				MaxDeliveryDistanceKm: sd.Restaurants[1].MaxDeliveryDistanceKm,
				MinSpend:              50.00,
				TaxRate:               sd.Restaurants[1].TaxRate,
				DateCreated:           sd.Restaurants[1].DateCreated,
			},
			ExcFunc: func(ctx context.Context) any {
				ur := restaurantbus.UpdateRestaurant{
					MinSpend: ptr(50.00),
				}

				resp, err := busDomain.Restaurant.Update(ctx, sd.Restaurants[1].Restaurant, ur)
				if err != nil {
					return err
				}

				return resp
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.(restaurantbus.Restaurant)
				if !exists {
					return "error occurred"
				}

				expResp := exp.(restaurantbus.Restaurant)
				expResp.DateUpdated = gotResp.DateUpdated

				return cmp.Diff(gotResp, expResp)
			},
		},
		{
			Name: "logo-url-update",
			ExpResp: restaurantbus.Restaurant{
				ID:                    sd.Restaurants[2].ID,
				OrganizationID:        sd.Restaurants[2].OrganizationID,
				Name:                  sd.Restaurants[2].Name,
				Description:           sd.Restaurants[2].Description,
				Address:               sd.Restaurants[2].Address,
				Phone:                 sd.Restaurants[2].Phone,
				Email:                 sd.Restaurants[2].Email,
				ImageURL:              sd.Restaurants[2].ImageURL,
				LogoURL:               "brand_new_logo.png",
				OperatingHours:        sd.Restaurants[2].OperatingHours,
				Enabled:               sd.Restaurants[2].Enabled,
				Latitude:              sd.Restaurants[2].Latitude,
				Longitude:             sd.Restaurants[2].Longitude,
				MaxDeliveryDistanceKm: sd.Restaurants[2].MaxDeliveryDistanceKm,
				MinSpend:              sd.Restaurants[2].MinSpend,
				TaxRate:               sd.Restaurants[2].TaxRate,
				DateCreated:           sd.Restaurants[2].DateCreated,
			},
			ExcFunc: func(ctx context.Context) any {
				ur := restaurantbus.UpdateRestaurant{
					LogoURL: ptr("brand_new_logo.png"),
				}

				resp, err := busDomain.Restaurant.Update(ctx, sd.Restaurants[2].Restaurant, ur)
				if err != nil {
					return err
				}

				return resp
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.(restaurantbus.Restaurant)
				if !exists {
					return "error occurred"
				}

				expResp := exp.(restaurantbus.Restaurant)
				expResp.DateUpdated = gotResp.DateUpdated

				return cmp.Diff(gotResp, expResp)
			},
		},
		{
			Name: "operating-hours-update",
			ExpResp: restaurantbus.Restaurant{
				ID:             sd.Restaurants[3].ID,
				OrganizationID: sd.Restaurants[3].OrganizationID,
				Name:           sd.Restaurants[3].Name,
				Description:    sd.Restaurants[3].Description,
				Address:        sd.Restaurants[3].Address,
				Phone:          sd.Restaurants[3].Phone,
				Email:          sd.Restaurants[3].Email,
				ImageURL:       sd.Restaurants[3].ImageURL,
				LogoURL:        sd.Restaurants[3].LogoURL,
				OperatingHours: restaurantbus.OperatingHours{
					"monday":    {Open: "08:00", Close: "20:00", IsClosed: false},
					"tuesday":   {Open: "08:00", Close: "20:00", IsClosed: false},
					"wednesday": {Open: "08:00", Close: "20:00", IsClosed: false},
					"thursday":  {Open: "08:00", Close: "20:00", IsClosed: false},
					"friday":    {Open: "08:00", Close: "21:00", IsClosed: false},
					"saturday":  {Open: "09:00", Close: "21:00", IsClosed: false},
					"sunday":    {Open: "09:00", Close: "18:00", IsClosed: true},
				},
				Enabled:               sd.Restaurants[3].Enabled,
				Latitude:              sd.Restaurants[3].Latitude,
				Longitude:             sd.Restaurants[3].Longitude,
				MaxDeliveryDistanceKm: sd.Restaurants[3].MaxDeliveryDistanceKm,
				MinSpend:              sd.Restaurants[3].MinSpend,
				TaxRate:               sd.Restaurants[3].TaxRate,
				DateCreated:           sd.Restaurants[3].DateCreated,
			},
			ExcFunc: func(ctx context.Context) any {
				customHours := restaurantbus.OperatingHours{
					"monday":    {Open: "08:00", Close: "20:00", IsClosed: false},
					"tuesday":   {Open: "08:00", Close: "20:00", IsClosed: false},
					"wednesday": {Open: "08:00", Close: "20:00", IsClosed: false},
					"thursday":  {Open: "08:00", Close: "20:00", IsClosed: false},
					"friday":    {Open: "08:00", Close: "21:00", IsClosed: false},
					"saturday":  {Open: "09:00", Close: "21:00", IsClosed: false},
					"sunday":    {Open: "09:00", Close: "18:00", IsClosed: true},
				}
				ur := restaurantbus.UpdateRestaurant{
					OperatingHours: &customHours,
				}

				resp, err := busDomain.Restaurant.Update(ctx, sd.Restaurants[3].Restaurant, ur)
				if err != nil {
					return err
				}

				return resp
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.(restaurantbus.Restaurant)
				if !exists {
					return "error occurred"
				}

				expResp := exp.(restaurantbus.Restaurant)
				expResp.DateUpdated = gotResp.DateUpdated

				return cmp.Diff(gotResp, expResp)
			},
		},
	}

	return table
}

func delete(busDomain dbtest.BusDomain, sd unittest.SeedData) []unittest.Table {
	table := []unittest.Table{
		{
			Name:    "basic",
			ExpResp: nil,
			ExcFunc: func(ctx context.Context) any {
				if err := busDomain.Restaurant.Delete(ctx, sd.Restaurants[1].Restaurant); err != nil {
					return err
				}

				return nil
			},
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
	}

	return table
}

// ptr is a helper function for getting the address of a value.
func ptr[T any](v T) *T {
	return &v
}
