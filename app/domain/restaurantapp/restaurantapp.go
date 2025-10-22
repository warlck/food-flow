package restaurantapp

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/warlck/food-flow/app/sdk/errs"
	"github.com/warlck/food-flow/app/sdk/query"
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
}

// newApp constructs a handlers for route access.
func newApp(restaurantBus *restaurantbus.Business, categoryBus *categorybus.Business, menuItemBus *menuitembus.Business) *app {
	return &app{
		restaurantBus: restaurantBus,
		categoryBus:   categoryBus,
		menuItemBus:   menuItemBus,
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
		ID:          res.ID.String(),
		Name:        res.Name.String(),
		Description: res.Description,
		Address:     res.Address,
		Phone:       res.Phone,
		Email:       res.Email,
		ImageURL:    res.ImageURL,
		Enabled:     res.Enabled,
		DateCreated: res.DateCreated.Format("2006-01-02T15:04:05Z07:00"),
		DateUpdated: res.DateUpdated.Format("2006-01-02T15:04:05Z07:00"),
		Categories:  make([]Category, 0, len(categories)),
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
		menuItemOrderBy := order.NewBy(menuitembus.OrderByID, order.ASC)
		menuItems, err := a.menuItemBus.Query(ctx, menuItemFilter, menuItemOrderBy, pg)
		if err != nil {
			return fmt.Errorf("query menu items: categoryID[%s]: %w", cat.ID, err)
		}

		// Convert menu items
		for _, item := range menuItems {
			appMenuItem := MenuItem{
				ID:          item.ID.String(),
				Name:        item.Name.String(),
				Description: item.Description,
				Price:       item.Price.Value(),
				ImageURL:    item.ImageURL,
				Available:   item.Available,
			}
			appCategory.MenuItems = append(appCategory.MenuItems, appMenuItem)
		}

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
