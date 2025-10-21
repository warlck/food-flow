package menuitemapp

import (
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/warlck/food-flow/app/sdk/errs"
	"github.com/warlck/food-flow/business/domain/menuitembus"
	"github.com/warlck/food-flow/business/types/money"
	"github.com/warlck/food-flow/business/types/name"
)

type queryParams struct {
	Page             string
	Rows             string
	OrderBy          string
	ID               string
	Name             string
	CategoryID       string
	RestaurantID     string
	MinPrice         string
	MaxPrice         string
	Available        string
	StartCreatedDate string
	EndCreatedDate   string
}

func parseQueryParams(r *http.Request) (queryParams, error) {
	values := r.URL.Query()

	qp := queryParams{
		Page:             values.Get("page"),
		Rows:             values.Get("rows"),
		OrderBy:          values.Get("orderBy"),
		ID:               values.Get("menu_item_id"),
		Name:             values.Get("name"),
		CategoryID:       values.Get("category_id"),
		RestaurantID:     values.Get("restaurant_id"),
		MinPrice:         values.Get("min_price"),
		MaxPrice:         values.Get("max_price"),
		Available:        values.Get("available"),
		StartCreatedDate: values.Get("start_created_date"),
		EndCreatedDate:   values.Get("end_created_date"),
	}

	return qp, nil
}

func parseFilter(qp queryParams) (menuitembus.QueryFilter, error) {
	var fieldErrors errs.FieldErrors
	var filter menuitembus.QueryFilter

	if qp.ID != "" {
		id, err := uuid.Parse(qp.ID)
		switch err {
		case nil:
			filter.ID = &id
		default:
			fieldErrors.Add("menu_item_id", err)
		}
	}

	if qp.Name != "" {
		name, err := name.Parse(qp.Name)
		switch err {
		case nil:
			filter.Name = &name
		default:
			fieldErrors.Add("name", err)
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

	if qp.MinPrice != "" {
		minPrice, err := strconv.ParseFloat(qp.MinPrice, 64)
		if err != nil {
			fieldErrors.Add("min_price", err)
		} else {
			price, err := money.Parse(minPrice)
			if err != nil {
				fieldErrors.Add("min_price", err)
			} else {
				filter.MinPrice = &price
			}
		}
	}

	if qp.MaxPrice != "" {
		maxPrice, err := strconv.ParseFloat(qp.MaxPrice, 64)
		if err != nil {
			fieldErrors.Add("max_price", err)
		} else {
			price, err := money.Parse(maxPrice)
			if err != nil {
				fieldErrors.Add("max_price", err)
			} else {
				filter.MaxPrice = &price
			}
		}
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

	if qp.StartCreatedDate != "" {
		t, err := time.Parse(time.RFC3339, qp.StartCreatedDate)
		switch err {
		case nil:
			filter.StartCreatedDate = &t
		default:
			fieldErrors.Add("start_created_date", err)
		}
	}

	if qp.EndCreatedDate != "" {
		t, err := time.Parse(time.RFC3339, qp.EndCreatedDate)
		switch err {
		case nil:
			filter.EndCreatedDate = &t
		default:
			fieldErrors.Add("end_created_date", err)
		}
	}

	if fieldErrors != nil {
		return menuitembus.QueryFilter{}, fieldErrors.ToError()
	}

	return filter, nil
}
