package categoryapp

import (
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/warlck/food-flow/app/sdk/errs"
	"github.com/warlck/food-flow/business/domain/categorybus"
	"github.com/warlck/food-flow/business/types/name"
)

type queryParams struct {
	Page             string
	Rows             string
	OrderBy          string
	ID               string
	Name             string
	RestaurantID     string
	Enabled          string
	StartCreatedDate string
	EndCreatedDate   string
}

func parseQueryParams(r *http.Request) (queryParams, error) {
	values := r.URL.Query()

	qp := queryParams{
		Page:             values.Get("page"),
		Rows:             values.Get("rows"),
		OrderBy:          values.Get("orderBy"),
		ID:               values.Get("category_id"),
		Name:             values.Get("name"),
		RestaurantID:     values.Get("restaurant_id"),
		Enabled:          values.Get("enabled"),
		StartCreatedDate: values.Get("start_created_date"),
		EndCreatedDate:   values.Get("end_created_date"),
	}

	return qp, nil
}

func parseFilter(qp queryParams) (categorybus.QueryFilter, error) {
	var fieldErrors errs.FieldErrors
	var filter categorybus.QueryFilter

	if qp.ID != "" {
		id, err := uuid.Parse(qp.ID)
		switch err {
		case nil:
			filter.ID = &id
		default:
			fieldErrors.Add("category_id", err)
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

	if qp.RestaurantID != "" {
		restaurantID, err := uuid.Parse(qp.RestaurantID)
		switch err {
		case nil:
			filter.RestaurantID = &restaurantID
		default:
			fieldErrors.Add("restaurant_id", err)
		}
	}

	if qp.Enabled != "" {
		enabled, err := strconv.ParseBool(qp.Enabled)
		switch err {
		case nil:
			filter.Enabled = &enabled
		default:
			fieldErrors.Add("enabled", err)
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
		return categorybus.QueryFilter{}, fieldErrors.ToError()
	}

	return filter, nil
}
