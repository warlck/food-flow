package modifiergroupapp

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/warlck/food-flow/app/sdk/errs"
	"github.com/warlck/food-flow/business/domain/modifiergroupbus"
	"github.com/warlck/food-flow/business/types/name"
)

type queryParams struct {
	Page         string
	Rows         string
	OrderBy      string
	ID           string
	MenuItemID   string
	RestaurantID string
	Name         string
	Available    string
}

func parseQueryParams(r *http.Request) (queryParams, error) {
	values := r.URL.Query()

	return queryParams{
		Page:         values.Get("page"),
		Rows:         values.Get("rows"),
		OrderBy:      values.Get("orderBy"),
		ID:           values.Get("id"),
		MenuItemID:   values.Get("menu_item_id"),
		RestaurantID: values.Get("restaurant_id"),
		Name:         values.Get("name"),
		Available:    values.Get("available"),
	}, nil
}

func parseFilter(qp queryParams) (modifiergroupbus.QueryFilter, error) {
	var filter modifiergroupbus.QueryFilter

	if qp.ID != "" {
		id, err := uuid.Parse(qp.ID)
		if err != nil {
			return modifiergroupbus.QueryFilter{}, errs.NewFieldErrors("id", err)
		}
		filter.ID = &id
	}

	if qp.MenuItemID != "" {
		id, err := uuid.Parse(qp.MenuItemID)
		if err != nil {
			return modifiergroupbus.QueryFilter{}, errs.NewFieldErrors("menu_item_id", err)
		}
		filter.MenuItemID = &id
	}

	if qp.RestaurantID != "" {
		id, err := uuid.Parse(qp.RestaurantID)
		if err != nil {
			return modifiergroupbus.QueryFilter{}, errs.NewFieldErrors("restaurant_id", err)
		}
		filter.RestaurantID = &id
	}

	if qp.Name != "" {
		nme, err := name.Parse(qp.Name)
		if err != nil {
			return modifiergroupbus.QueryFilter{}, errs.NewFieldErrors("name", err)
		}
		filter.Name = &nme
	}

	if qp.Available != "" {
		avail, err := strconv.ParseBool(qp.Available)
		if err != nil {
			return modifiergroupbus.QueryFilter{}, errs.NewFieldErrors("available", fmt.Errorf("invalid bool: %w", err))
		}
		filter.Available = &avail
	}

	return filter, nil
}
