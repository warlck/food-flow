package modifieroptionapp

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/warlck/food-flow/app/sdk/errs"
	"github.com/warlck/food-flow/business/domain/modifieroptionbus"
	"github.com/warlck/food-flow/business/types/name"
)

type queryParams struct {
	Page            string
	Rows            string
	OrderBy         string
	ID              string
	ModifierGroupID string
	RestaurantID    string
	Name            string
	Available       string
}

func parseQueryParams(r *http.Request) (queryParams, error) {
	values := r.URL.Query()

	return queryParams{
		Page:            values.Get("page"),
		Rows:            values.Get("rows"),
		OrderBy:         values.Get("orderBy"),
		ID:              values.Get("id"),
		ModifierGroupID: values.Get("modifier_group_id"),
		RestaurantID:    values.Get("restaurant_id"),
		Name:            values.Get("name"),
		Available:       values.Get("available"),
	}, nil
}

func parseFilter(qp queryParams) (modifieroptionbus.QueryFilter, error) {
	var filter modifieroptionbus.QueryFilter

	if qp.ID != "" {
		id, err := uuid.Parse(qp.ID)
		if err != nil {
			return modifieroptionbus.QueryFilter{}, errs.NewFieldErrors("id", err)
		}
		filter.ID = &id
	}

	if qp.ModifierGroupID != "" {
		id, err := uuid.Parse(qp.ModifierGroupID)
		if err != nil {
			return modifieroptionbus.QueryFilter{}, errs.NewFieldErrors("modifier_group_id", err)
		}
		filter.ModifierGroupID = &id
	}

	if qp.RestaurantID != "" {
		id, err := uuid.Parse(qp.RestaurantID)
		if err != nil {
			return modifieroptionbus.QueryFilter{}, errs.NewFieldErrors("restaurant_id", err)
		}
		filter.RestaurantID = &id
	}

	if qp.Name != "" {
		nme, err := name.Parse(qp.Name)
		if err != nil {
			return modifieroptionbus.QueryFilter{}, errs.NewFieldErrors("name", err)
		}
		filter.Name = &nme
	}

	if qp.Available != "" {
		avail, err := strconv.ParseBool(qp.Available)
		if err != nil {
			return modifieroptionbus.QueryFilter{}, errs.NewFieldErrors("available", fmt.Errorf("invalid bool: %w", err))
		}
		filter.Available = &avail
	}

	return filter, nil
}
