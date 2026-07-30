package addonapp

import (
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/warlck/food-flow/app/sdk/errs"
	"github.com/warlck/food-flow/business/domain/addonbus"
)

type queryParams struct {
	Page         string
	Rows         string
	OrderBy      string
	ID           string
	CategoryID   string
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
		CategoryID:   values.Get("category_id"),
		RestaurantID: values.Get("restaurant_id"),
		Name:         values.Get("name"),
		Available:    values.Get("available"),
	}

	return qp, nil
}

func parseFilter(qp queryParams) (addonbus.QueryFilter, error) {
	var fieldErrors errs.FieldErrors
	var filter addonbus.QueryFilter

	if qp.ID != "" {
		id, err := uuid.Parse(qp.ID)
		switch err {
		case nil:
			filter.ID = &id
		default:
			fieldErrors.Add("addon_id", err)
		}
	}

	if qp.CategoryID != "" {
		categoryID, err := uuid.Parse(qp.CategoryID)
		switch err {
		case nil:
			filter.CategoryID = &categoryID
		default:
			fieldErrors.Add("category_id", err)
		}
	}

	if qp.RestaurantID != "" {
		restaurantID, err := uuid.Parse(qp.RestaurantID)
		switch err {
		case nil:
			filter.RestaurantID = &restaurantID
		default:
			fieldErrors.Add("restaurant_id", err)
		}
	}

	if qp.Name != "" {
		filter.Name = &qp.Name
	}

	if qp.Available != "" {
		available, err := strconv.ParseBool(qp.Available)
		switch err {
		case nil:
			filter.Available = &available
		default:
			fieldErrors.Add("available", err)
		}
	}

	if fieldErrors != nil {
		return addonbus.QueryFilter{}, fieldErrors
	}

	return filter, nil
}
