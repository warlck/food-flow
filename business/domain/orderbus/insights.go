package orderbus

import (
	"time"

	"github.com/google/uuid"
	"github.com/warlck/food-flow/business/types/money"
)

// InsightsFilter holds filtering criteria for analytics queries.
type InsightsFilter struct {
	RestaurantID *string
	StartDate    *time.Time
	EndDate      *time.Time
}

// SalesSummary represents aggregate KPI figures for a given date range.
type SalesSummary struct {
	GrossSales        money.Money
	NetSales          money.Money
	TotalOrders       int
	CompletedOrders   int
	CancelledOrders   int
	AverageOrderValue money.Money
	TotalDiscounts    money.Money
	TotalDeliveryFees money.Money
	TotalTax          money.Money
}

// TimeSeriesPoint represents daily sales performance in a timeline.
type TimeSeriesPoint struct {
	Date         string
	GrossSales   money.Money
	NetSales     money.Money
	OrderCount   int
	AverageOrder money.Money
}

// ItemSalesMetric represents order-level sales numbers for a menu item.
type ItemSalesMetric struct {
	MenuItemID   uuid.UUID
	MenuItemName string
	QuantitySold int
	TotalRevenue money.Money
}

// TopAddonMetric represents popular upsell add-on performance.
type TopAddonMetric struct {
	AddonID      uuid.UUID
	AddonName    string
	QuantitySold int
	TotalRevenue money.Money
}

// OrderTypeMetric represents fulfillment channel breakdown (e.g. delivery, pickup).
type OrderTypeMetric struct {
	OrderType    string
	Count        int
	TotalRevenue money.Money
	Percentage   float64
}

// HourlyMetric represents 24-hour order distribution for identifying peak hours.
type HourlyMetric struct {
	Hour         int
	Count        int
	TotalRevenue money.Money
}

// OrderMetrics combines all pure order-domain analytics datasets.
type OrderMetrics struct {
	Summary       SalesSummary
	SalesOverTime []TimeSeriesPoint
	TopItems      []ItemSalesMetric
	AllItemSales  []ItemSalesMetric
	TopAddons     []TopAddonMetric
	OrderTypes    []OrderTypeMetric
	PeakHours     []HourlyMetric
}
