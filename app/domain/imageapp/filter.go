package imageapp

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/warlck/food-flow/app/sdk/errs"
	"github.com/warlck/food-flow/business/domain/imagebus"
	"github.com/warlck/food-flow/business/sdk/order"
)

type queryParams struct {
	Page         string
	Rows         string
	OrderBy      string
	RestaurantID string
	EntityType   string
	Status       string
}

func parseQueryParams(r *http.Request) (queryParams, error) {
	values := r.URL.Query()

	qp := queryParams{
		Page:         values.Get("page"),
		Rows:         values.Get("rows"),
		OrderBy:      values.Get("orderBy"),
		RestaurantID: values.Get("restaurant_id"),
		EntityType:   values.Get("entity_type"),
		Status:       values.Get("status"),
	}

	return qp, nil
}

func parseFilter(qp queryParams) (imagebus.QueryFilter, error) {
	var fieldErrors errs.FieldErrors
	var filter imagebus.QueryFilter

	if qp.RestaurantID != "" {
		restaurantID, err := uuid.Parse(qp.RestaurantID)
		if err != nil {
			fieldErrors.Add("restaurant_id", err)
		} else {
			filter.RestaurantID = &restaurantID
		}
	}

	if qp.EntityType != "" {
		switch qp.EntityType {
		case imagebus.EntityTypeRestaurant, imagebus.EntityTypeMenuItem:
			filter.EntityType = &qp.EntityType
		default:
			fieldErrors.Add("entity_type", fmt.Errorf("invalid entity type %q", qp.EntityType))
		}
	}

	if qp.Status != "" {
		switch qp.Status {
		case imagebus.StatusPending, imagebus.StatusConfirmed:
			filter.Status = &qp.Status
		default:
			fieldErrors.Add("status", fmt.Errorf("invalid status %q", qp.Status))
		}
	}

	if fieldErrors != nil {
		return imagebus.QueryFilter{}, fieldErrors.ToError()
	}

	return filter, nil
}

var defaultOrderBy = order.NewBy(imagebus.OrderByDateCreated, order.DESC)

var orderByFields = map[string]string{
	imagebus.OrderByID:          imagebus.OrderByID,
	imagebus.OrderByDateCreated: imagebus.OrderByDateCreated,
}
