package restaurantapp

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/warlck/food-flow/app/sdk/errs"
	"github.com/warlck/food-flow/app/sdk/mid"
	"github.com/warlck/food-flow/app/sdk/query"
	"github.com/warlck/food-flow/business/domain/addonbus"
	"github.com/warlck/food-flow/business/domain/categorybus"
	"github.com/warlck/food-flow/business/domain/menuitembus"
	"github.com/warlck/food-flow/business/domain/modifiergroupbus"
	"github.com/warlck/food-flow/business/domain/modifieroptionbus"
	"github.com/warlck/food-flow/business/domain/restaurantbus"
	"github.com/warlck/food-flow/business/sdk/order"
	"github.com/warlck/food-flow/business/sdk/page"
	"github.com/warlck/food-flow/foundation/web"
)

// Handlers manages the set of restaurant endpoints.
type app struct {
	restaurantBus     *restaurantbus.Business
	categoryBus       *categorybus.Business
	menuItemBus       *menuitembus.Business
	modifierGroupBus  *modifiergroupbus.Business
	modifierOptionBus *modifieroptionbus.Business
	addonBus          *addonbus.Business
}

// newApp constructs a handlers for route access.
func newApp(
	restaurantBus *restaurantbus.Business,
	categoryBus *categorybus.Business,
	menuItemBus *menuitembus.Business,
	modifierGroupBus *modifiergroupbus.Business,
	modifierOptionBus *modifieroptionbus.Business,
	addonBus *addonbus.Business,
) *app {
	return &app{
		restaurantBus:     restaurantBus,
		categoryBus:       categoryBus,
		menuItemBus:       menuItemBus,
		modifierGroupBus:  modifierGroupBus,
		modifierOptionBus: modifierOptionBus,
		addonBus:          addonBus,
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

	claims := mid.GetClaims(ctx)
	if !claims.IsOrgAuthorized(nr.OrganizationID) {
		return errs.Newf(errs.PermissionDenied, "user not in organization %s", nr.OrganizationID)
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

// QueryByIDWithDetails retrieves a restaurant by its ID along with all enabled categories and full menu hierarchy.
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

	// Fetch categories for this restaurant (sorted by rank ASC NULLS LAST, name ASC, category_id ASC)
	categoryFilter := categorybus.QueryFilter{
		RestaurantID: &restaurantID,
	}
	categories, err := a.categoryBus.QueryAll(ctx, categoryFilter, categorybus.DefaultOrderBy)
	if err != nil {
		return fmt.Errorf("query categories: restaurantID[%s]: %w", restaurantID, err)
	}

	// Fetch menu items for this restaurant
	menuItemFilter := menuitembus.QueryFilter{
		RestaurantID: &restaurantID,
	}
	menuItems, err := a.menuItemBus.QueryAll(ctx, menuItemFilter, menuitembus.DefaultOrderBy)
	if err != nil {
		return fmt.Errorf("query menu items: restaurantID[%s]: %w", restaurantID, err)
	}

	// Fetch modifier groups for this restaurant
	modGroups, err := a.modifierGroupBus.QueryAll(ctx, modifiergroupbus.QueryFilter{RestaurantID: &restaurantID}, modifiergroupbus.DefaultOrderBy)
	if err != nil {
		return fmt.Errorf("query modifier groups: restaurantID[%s]: %w", restaurantID, err)
	}

	// Fetch modifier options for this restaurant
	modOptions, err := a.modifierOptionBus.QueryAll(ctx, modifieroptionbus.QueryFilter{RestaurantID: &restaurantID}, modifieroptionbus.DefaultOrderBy)
	if err != nil {
		return fmt.Errorf("query modifier options: restaurantID[%s]: %w", restaurantID, err)
	}

	// Group modifier options by group ID
	optionsByGroup := make(map[uuid.UUID][]Option)
	for _, opt := range modOptions {
		optionsByGroup[opt.ModifierGroupID] = append(optionsByGroup[opt.ModifierGroupID], Option{
			ID:          opt.ID.String(),
			Name:        opt.Name.String(),
			Description: opt.Description,
			PriceDelta:  opt.PriceDelta.Value(),
			Available:   opt.Available,
			Rank:        opt.Rank,
		})
	}

	// Group modifier groups by menu item ID
	groupsByItem := make(map[uuid.UUID][]ModifierGroup)
	for _, grp := range modGroups {
		opts := optionsByGroup[grp.ID]
		if opts == nil {
			opts = []Option{}
		}
		groupsByItem[grp.MenuItemID] = append(groupsByItem[grp.MenuItemID], ModifierGroup{
			ID:            grp.ID.String(),
			Name:          grp.Name.String(),
			Description:   grp.Description,
			MinSelections: grp.MinSelections,
			MaxSelections: grp.MaxSelections,
			Available:     grp.Available,
			Rank:          grp.Rank,
			Options:       opts,
		})
	}

	// Group menu items by category ID
	itemsByCategory := make(map[uuid.UUID][]MenuItem)
	for _, item := range menuItems {
		itemGroups := groupsByItem[item.ID]
		if itemGroups == nil {
			itemGroups = []ModifierGroup{}
		}

		// Query assigned addons for this item
		assignedAddons, err := a.addonBus.QueryMenuItemAddons(ctx, item.ID)
		if err != nil {
			return fmt.Errorf("query menu item addons: %w", err)
		}

		appAddons := make([]Addon, 0, len(assignedAddons))
		for _, aInfo := range assignedAddons {
			appAddons = append(appAddons, Addon{
				ID:          aInfo.Addon.ID.String(),
				Name:        aInfo.Addon.Name.String(),
				Description: aInfo.Addon.Description,
				Price:       aInfo.Addon.Price.Value(),
				Available:   aInfo.Addon.Available,
				MaxQuantity: aInfo.Addon.MaxQuantity,
				Rank:        aInfo.Rank,
			})
		}

		// Calculate orderable
		isOrderable := item.Available
		if isOrderable {
			for _, grp := range itemGroups {
				if grp.Available && grp.MinSelections > 0 {
					hasAvailOpt := false
					for _, opt := range grp.Options {
						if opt.Available {
							hasAvailOpt = true
							break
						}
					}
					if !hasAvailOpt {
						isOrderable = false
						break
					}
				}
			}
		}

		appMenuItem := MenuItem{
			ID:             item.ID.String(),
			Name:           item.Name.String(),
			Description:    item.Description,
			Price:          item.Price.Value(),
			ImageURL:       item.ImageURL,
			Available:      item.Available,
			Orderable:      isOrderable,
			Rank:           item.Rank,
			ModifierGroups: itemGroups,
			Addons:         appAddons,
		}

		itemsByCategory[item.CategoryID] = append(itemsByCategory[item.CategoryID], appMenuItem)
	}

	// Build the response with enabled categories only
	appCategories := make([]Category, 0, len(categories))
	for _, cat := range categories {
		if !cat.Enabled {
			continue
		}

		catItems := itemsByCategory[cat.ID]
		if catItems == nil {
			catItems = []MenuItem{}
		}

		appCategories = append(appCategories, Category{
			ID:          cat.ID.String(),
			Name:        cat.Name.String(),
			Description: cat.Description,
			Enabled:     cat.Enabled,
			Rank:        cat.Rank,
			MenuItems:   catItems,
		})
	}

	appRestaurant := RestaurantWithMenuCategories{
		ID:                    res.ID.String(),
		Name:                  res.Name.String(),
		Description:           res.Description,
		Address:               res.Address,
		Phone:                 res.Phone,
		Email:                 res.Email,
		ImageURL:              res.ImageURL,
		LogoURL:               res.LogoURL,
		OperatingHours:        ToAppOperatingHours(res.OperatingHours),
		Enabled:               res.Enabled,
		Latitude:              res.Latitude,
		Longitude:             res.Longitude,
		MaxDeliveryDistanceKm: res.MaxDeliveryDistanceKm,
		MinSpend:              res.MinSpend,
		TaxRate:               res.TaxRate,
		Categories:            appCategories,
		DateCreated:           res.DateCreated.Format("2006-01-02T15:04:05Z07:00"),
		DateUpdated:           res.DateUpdated.Format("2006-01-02T15:04:05Z07:00"),
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

	claims := mid.GetClaims(ctx)
	if !claims.IsOrgAuthorized(res.OrganizationID) {
		return errs.Newf(errs.PermissionDenied, "user not in organization %s", res.OrganizationID)
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

	claims := mid.GetClaims(ctx)
	if !claims.IsOrgAuthorized(res.OrganizationID) {
		return errs.Newf(errs.PermissionDenied, "user not in organization %s", res.OrganizationID)
	}

	if err := a.restaurantBus.Delete(ctx, res); err != nil {
		return errs.Newf(errs.Internal, "delete: restaurantID[%s]: %s", restaurantID, err)
	}

	return web.Respond(ctx, w, struct{}{}, http.StatusNoContent)
}
