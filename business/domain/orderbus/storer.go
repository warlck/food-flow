package orderbus

import (
	"context"

	"github.com/google/uuid"
	"github.com/warlck/food-flow/business/sdk/order"
	"github.com/warlck/food-flow/business/sdk/page"
)

// Storer interface declares the behavior this package needs to persist and
// retrieve data.
type Storer interface {
	Create(ctx context.Context, order Order) error
	Query(ctx context.Context, filter QueryFilter, orderBy order.By, page page.Page) ([]Order, error)
	QueryByID(ctx context.Context, orderID uuid.UUID) (Order, error)
	Update(ctx context.Context, order Order) error
	Delete(ctx context.Context, orderID uuid.UUID) error
	Count(ctx context.Context, filter QueryFilter) (int, error)
}
