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

	// Order Metrics
	QuerySalesSummary(ctx context.Context, filter InsightsFilter) (SalesSummary, error)
	QuerySalesOverTime(ctx context.Context, filter InsightsFilter) ([]TimeSeriesPoint, error)
	QueryTopItemSales(ctx context.Context, filter InsightsFilter, limit int) ([]ItemSalesMetric, error)
	QueryAllItemSales(ctx context.Context, filter InsightsFilter) ([]ItemSalesMetric, error)
	QueryTopAddonSales(ctx context.Context, filter InsightsFilter, limit int) ([]TopAddonMetric, error)
	QueryOrderTypes(ctx context.Context, filter InsightsFilter) ([]OrderTypeMetric, error)
	QueryPeakHours(ctx context.Context, filter InsightsFilter) ([]HourlyMetric, error)
}
