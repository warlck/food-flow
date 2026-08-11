package promoapp

import (
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/warlck/food-flow/app/sdk/errs"
	"github.com/warlck/food-flow/business/domain/promobus"
	"github.com/warlck/food-flow/business/sdk/order"
	"github.com/warlck/food-flow/business/types/name"
)

type queryParams struct {
	Page         string
	Rows         string
	OrderBy      string
	ID           string
	Code         string
	Name         string
	RestaurantID string
	Enabled      string
}

func parseQueryParams(r *http.Request) (queryParams, error) {
	values := r.URL.Query()

	qp := queryParams{
		Page:         values.Get("page"),
		Rows:         values.Get("rows"),
		OrderBy:      values.Get("orderBy"),
		ID:           values.Get("id"),
		Code:         values.Get("code"),
		Name:         values.Get("name"),
		RestaurantID: values.Get("restaurant_id"),
		Enabled:      values.Get("enabled"),
	}

	return qp, nil
}

func parseFilter(qp queryParams) (promobus.QueryFilter, error) {
	var fieldErrors errs.FieldErrors
	var filter promobus.QueryFilter

	if qp.ID != "" {
		id, err := uuid.Parse(qp.ID)
		if err != nil {
			fieldErrors.Add("id", err)
		} else {
			filter.ID = &id
		}
	}

	if qp.Code != "" {
		filter.Code = &qp.Code
	}

	if qp.Name != "" {
		nme, err := name.Parse(qp.Name)
		if err != nil {
			fieldErrors.Add("name", err)
		} else {
			filter.Name = &nme
		}
	}

	if qp.RestaurantID != "" {
		restaurantID, err := uuid.Parse(qp.RestaurantID)
		if err != nil {
			fieldErrors.Add("restaurant_id", err)
		} else {
			filter.RestaurantID = &restaurantID
		}
	}

	if qp.Enabled != "" {
		enabled, err := strconv.ParseBool(qp.Enabled)
		if err != nil {
			fieldErrors.Add("enabled", err)
		} else {
			filter.Enabled = &enabled
		}
	}

	if fieldErrors != nil {
		return promobus.QueryFilter{}, fieldErrors.ToError()
	}

	return filter, nil
}

var defaultOrderBy = order.NewBy(promobus.OrderByCode, order.ASC)

var orderByFields = map[string]string{
	promobus.OrderByID:          promobus.OrderByID,
	promobus.OrderByCode:        promobus.OrderByCode,
	promobus.OrderByName:        promobus.OrderByName,
	promobus.OrderByDateCreated: promobus.OrderByDateCreated,
}
