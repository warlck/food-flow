package restaurantapp

import (
	"context"
	"fmt"
	"net/http"
	"sort"

	"github.com/google/uuid"
	"github.com/warlck/food-flow/app/sdk/errs"
	"github.com/warlck/food-flow/app/sdk/query"
	"github.com/warlck/food-flow/business/domain/addonbus"
	"github.com/warlck/food-flow/business/domain/categorybus"
	"github.com/warlck/food-flow/business/domain/menuitembus"
	"github.com/warlck/food-flow/business/domain/restaurantbus"
	"github.com/warlck/food-flow/business/sdk/order"
	"github.com/warlck/food-flow/business/sdk/page"
	"github.com/warlck/food-flow/foundation/web"
)

// Handlers manages the set of restaurant endpoints.
type app struct {
	restaurantBus *restaurantbus.Business
	categoryBus   *categorybus.Business
	menuItemBus   *menuitembus.Business
	addonBus      *addonbus.Business
}

// newApp constructs a handlers for route access.
func newApp(restaurantBus *restaurantbus.Business, categoryBus *categorybus.Business, menuItemBus *menuitembus.Business, addonBus *addonbus.Business) *app {
	return &app{
		restaurantBus: restaurantBus,
		categoryBus:   categoryBus,
		menuItemBus:   menuItemBus,
		addonBus:      addonBus,
	}
}

// Create adds a new restaurant to the system.
func (a *app) create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var app NewRestaurant
	if err := web.Decode(r, &app); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	nr, err := toBusNewRestaurant(app)
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	res, err := a.restaurantBus.Create(ctx, nr)
	if err != nil {
		return fmt.Errorf("create: restaurant[%+v]: %w", res, err)
	}

	return web.Respond(ctx, w, ToAppRestaurant(res), http.StatusCreated)
}

// Query retrieves a list of restaurants based on query parameters.
func (a *app) query(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	qp, err := parseQueryParams(r)
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	pg, err := page.Parse(qp.Page, qp.Rows)
	if err != nil {
		return errs.NewFieldErrors("page", err)
	}

	filter, err := parseFilter(qp)
	if err != nil {
		return err.(*errs.Error)
	}

	orderBy, err := order.Parse(orderByFields, qp.OrderBy, defaultOrderBy)
	if err != nil {
		return errs.NewFieldErrors("order", err)
	}

	restaurants, err := a.restaurantBus.Query(ctx, filter, orderBy, pg)
	if err != nil {
		return errs.Newf(errs.Internal, "query: %s", err)
	}

	total, err := a.restaurantBus.Count(ctx, filter)
	if err != nil {
		return errs.Newf(errs.Internal, "count: %s", err)
	}

	result := query.NewResult(ToAppRestaurants(restaurants), total, pg)
	return web.Respond(ctx, w, result, http.StatusOK)
}

// QueryByID retrieves a restaurant by its ID.
func (a *app) queryByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	restaurantIDStr := web.Param(r, "restaurant_id")

	restaurantID, err := uuid.Parse(restaurantIDStr)
	if err != nil {
		return errs.NewFieldErrors("restaurant_id", err)
	}

	res, err := a.restaurantBus.QueryByID(ctx, restaurantID)
	if err != nil {
		return fmt.Errorf("querybyid: restaurantID[%s]: %w", restaurantID, err)
	}

	return web.Respond(ctx, w, ToAppRestaurant(res), http.StatusOK)
}

// QueryByIDWithDetails retrieves a restaurant by its ID along with all categories and menu items.
func (a *app) queryByIDWithDetails(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	restaurantIDStr := web.Param(r, "restaurant_id")

	restaurantID, err := uuid.Parse(restaurantIDStr)
	if err != nil {
		return errs.NewFieldErrors("restaurant_id", err)
	}

	// Fetch the restaurant
	res, err := a.restaurantBus.QueryByID(ctx, restaurantID)
	if err != nil {
		return fmt.Errorf("querybyid: restaurantID[%s]: %w", restaurantID, err)
	}

	// Fetch categories for this restaurant
	categoryFilter := categorybus.QueryFilter{
		RestaurantID: &restaurantID,
	}
	// Use a large page to get all categories (adjust if needed)
	pg := page.MustParse("1", "100")
	categoryOrderBy := order.NewBy(categorybus.OrderByID, order.ASC)
	categories, err := a.categoryBus.Query(ctx, categoryFilter, categoryOrderBy, pg)
	if err != nil {
		return fmt.Errorf("query categories: restaurantID[%s]: %w", restaurantID, err)
	}

	// Build the restaurant with nested data
	appRestaurant := RestaurantWithMenuCategories{
		ID:                    res.ID.String(),
		Name:                  res.Name.String(),
		Description:           res.Description,
		Address:               res.Address,
		Phone:                 res.Phone,
		Email:                 res.Email,
		ImageURL:              res.ImageURL,
		Enabled:               res.Enabled,
		Latitude:              res.Latitude,
		Longitude:             res.Longitude,
		MaxDeliveryDistanceKm: res.MaxDeliveryDistanceKm,
		TaxRate:               res.TaxRate,
		DateCreated:           res.DateCreated.Format("2006-01-02T15:04:05Z07:00"),
		DateUpdated:           res.DateUpdated.Format("2006-01-02T15:04:05Z07:00"),
		Categories:            make([]Category, 0, len(categories)),
	}

	// For each category, fetch its menu items
	for _, cat := range categories {
		appCategory := Category{
			ID:          cat.ID.String(),
			Name:        cat.Name.String(),
			Description: cat.Description,
			Enabled:     cat.Enabled,
			MenuItems:   make([]MenuItem, 0),
		}

		// Fetch menu items for this category
		menuItemFilter := menuitembus.QueryFilter{
			CategoryID: &cat.ID,
		}
		// Sort menu items by price (cheapest first) so the frontend can render the default
		// item per category without client-side sorting.
		menuItemOrderBy := order.NewBy(menuitembus.OrderByPrice, order.ASC)
		menuItems, err := a.menuItemBus.Query(ctx, menuItemFilter, menuItemOrderBy, pg)
		if err != nil {
			return fmt.Errorf("query menu items: categoryID[%s]: %w", cat.ID, err)
		}

		// Fetch addons for this category (shared across menu items).
		addons, err := a.addonBus.QueryByCategoryID(ctx, cat.ID)
		if err != nil {
			return fmt.Errorf("query addons: categoryID[%s]: %w", cat.ID, err)
		}

		// Convert addons to app layer (always return a non-nil slice in JSON).
		appAddons := []Addon{}
		for _, addon := range addons {
			if addon.Available {
				appAddons = append(appAddons, Addon{
					ID:          addon.ID.String(),
					Name:        addon.Name.String(),
					Description: addon.Description,
					Price:       addon.Price.Value(),
					Available:   addon.Available,
					MaxQuantity: addon.MaxQuantity,
				})
			}
		}

		// Convert menu items and attach category addons.
		for _, item := range menuItems {
			appMenuItem := MenuItem{
				ID:          item.ID.String(),
				Name:        item.Name.String(),
				Description: item.Description,
				Price:       item.Price.Value(),
				ImageURL:    item.ImageURL,
				Available:   item.Available,
				Addons:      appAddons,
			}
			appCategory.MenuItems = append(appCategory.MenuItems, appMenuItem)
		}

		// Make ordering deterministic (price asc, then name, then id).
		sort.Slice(appCategory.MenuItems, func(i, j int) bool {
			if appCategory.MenuItems[i].Price != appCategory.MenuItems[j].Price {
				return appCategory.MenuItems[i].Price < appCategory.MenuItems[j].Price
			}
			if appCategory.MenuItems[i].Name != appCategory.MenuItems[j].Name {
				return appCategory.MenuItems[i].Name < appCategory.MenuItems[j].Name
			}
			return appCategory.MenuItems[i].ID < appCategory.MenuItems[j].ID
		})

		appRestaurant.Categories = append(appRestaurant.Categories, appCategory)
	}

	return web.Respond(ctx, w, appRestaurant, http.StatusOK)
}

// Update modifies a restaurant.
func (a *app) update(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	restaurantIDStr := web.Param(r, "restaurant_id")

	restaurantID, err := uuid.Parse(restaurantIDStr)
	if err != nil {
		return errs.NewFieldErrors("restaurant_id", err)
	}

	var app UpdateRestaurant
	if err := web.Decode(r, &app); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	ur, err := toBusUpdateRestaurant(app)
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	res, err := a.restaurantBus.QueryByID(ctx, restaurantID)
	if err != nil {
		return fmt.Errorf("querybyid: restaurantID[%s]: %w", restaurantID, err)
	}

	updRes, err := a.restaurantBus.Update(ctx, res, ur)
	if err != nil {
		return errs.Newf(errs.Internal, "update: restaurantID[%s] ur[%+v]: %s", restaurantID, ur, err)
	}

	return web.Respond(ctx, w, ToAppRestaurant(updRes), http.StatusOK)
}

// Delete removes a restaurant.
func (a *app) delete(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	restaurantIDStr := web.Param(r, "restaurant_id")

	restaurantID, err := uuid.Parse(restaurantIDStr)
	if err != nil {
		return errs.NewFieldErrors("restaurant_id", err)
	}

	res, err := a.restaurantBus.QueryByID(ctx, restaurantID)
	if err != nil {
		return fmt.Errorf("querybyid: restaurantID[%s]: %w", restaurantID, err)
	}

	if err := a.restaurantBus.Delete(ctx, res); err != nil {
		return errs.Newf(errs.Internal, "delete: restaurantID[%s]: %s", restaurantID, err)
	}

	return web.Respond(ctx, w, struct{}{}, http.StatusNoContent)
}
