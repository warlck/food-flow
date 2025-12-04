package orderapp

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/warlck/food-flow/app/sdk/errs"
	"github.com/warlck/food-flow/business/domain/orderbus"
)

type queryParams struct {
	Page             string
	Rows             string
	OrderBy          string
	ID               string
	RestaurantID     string
	CustomerEmail    string
	OrderStatus      string
	PaymentStatus    string
	OrderType        string
	StartCreatedDate string
	EndCreatedDate   string
}

func parseQueryParams(r *http.Request) (queryParams, error) {
	values := r.URL.Query()

	qp := queryParams{
		Page:             values.Get("page"),
		Rows:             values.Get("rows"),
		OrderBy:          values.Get("orderBy"),
		ID:               values.Get("order_id"),
		RestaurantID:     values.Get("restaurant_id"),
		CustomerEmail:    values.Get("customer_email"),
		OrderStatus:      values.Get("order_status"),
		PaymentStatus:    values.Get("payment_status"),
		OrderType:        values.Get("order_type"),
		StartCreatedDate: values.Get("start_date"),
		EndCreatedDate:   values.Get("end_date"),
	}

	return qp, nil
}

func parseFilter(qp queryParams) (orderbus.QueryFilter, error) {
	var fieldErrors errs.FieldErrors
	var filter orderbus.QueryFilter

	if qp.ID != "" {
		id, err := uuid.Parse(qp.ID)
		switch err {
		case nil:
			idStr := id.String()
			filter.ID = &idStr
		default:
			fieldErrors.Add("order_id", err)
		}
	}

	if qp.RestaurantID != "" {
		restaurantID, err := uuid.Parse(qp.RestaurantID)
		switch err {
		case nil:
			ridStr := restaurantID.String()
			filter.RestaurantID = &ridStr
		default:
			fieldErrors.Add("restaurant_id", err)
		}
	}

	if qp.CustomerEmail != "" {
		filter.CustomerEmail = &qp.CustomerEmail
	}

	if qp.OrderStatus != "" {
		filter.OrderStatus = &qp.OrderStatus
	}

	if qp.PaymentStatus != "" {
		filter.PaymentStatus = &qp.PaymentStatus
	}

	if qp.OrderType != "" {
		filter.OrderType = &qp.OrderType
	}

	if qp.StartCreatedDate != "" {
		t, err := time.Parse(time.RFC3339, qp.StartCreatedDate)
		switch err {
		case nil:
			filter.StartDate = &t
		default:
			fieldErrors.Add("start_date", err)
		}
	}

	if qp.EndCreatedDate != "" {
		t, err := time.Parse(time.RFC3339, qp.EndCreatedDate)
		switch err {
		case nil:
			filter.EndDate = &t
		default:
			fieldErrors.Add("end_date", err)
		}
	}

	if fieldErrors != nil {
		return orderbus.QueryFilter{}, fieldErrors.ToError()
	}

	return filter, nil
}
