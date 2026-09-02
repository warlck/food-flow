package addonapp

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/warlck/food-flow/app/sdk/errs"
	"github.com/warlck/food-flow/business/domain/addonbus"
	"github.com/warlck/food-flow/business/types/name"
)

type queryParams struct {
	Page         string
	Rows         string
	OrderBy      string
	ID           string
	RestaurantID string
	Name         string
	Available    string
}

func parseQueryParams(r *http.Request) (queryParams, error) {
	values := r.URL.Query()

	qp := queryParams{
		Page:         values.Get("page"),
		Rows:         values.Get("rows"),
		OrderBy:      values.Get("orderBy"),
		ID:           values.Get("addon_id"),
		RestaurantID: values.Get("restaurant_id"),
		Name:         values.Get("name"),
		Available:    values.Get("available"),
	}

	return qp, nil
}

func parseFilter(qp queryParams) (addonbus.QueryFilter, error) {
	var filter addonbus.QueryFilter

	if qp.ID != "" {
		id, err := uuid.Parse(qp.ID)
		if err != nil {
			return addonbus.QueryFilter{}, errs.NewFieldErrors("addon_id", err)
		}
		filter.ID = &id
	}

	if qp.RestaurantID != "" {
		restaurantID, err := uuid.Parse(qp.RestaurantID)
		if err != nil {
			return addonbus.QueryFilter{}, errs.NewFieldErrors("restaurant_id", err)
		}
		filter.RestaurantID = &restaurantID
	}

	if qp.Name != "" {
		nme, err := name.Parse(qp.Name)
		if err != nil {
			return addonbus.QueryFilter{}, errs.NewFieldErrors("name", err)
		}
		filter.Name = &nme
	}

	if qp.Available != "" {
		available, err := strconv.ParseBool(qp.Available)
		if err != nil {
			return addonbus.QueryFilter{}, errs.NewFieldErrors("available", fmt.Errorf("invalid bool: %w", err))
		}
		filter.Available = &available
	}

	return filter, nil
}
