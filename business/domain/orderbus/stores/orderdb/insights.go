package orderdb

import (
	"bytes"
	"context"
	"fmt"
	"math"

	"github.com/google/uuid"
	"github.com/warlck/food-flow/business/domain/orderbus"
	"github.com/warlck/food-flow/business/sdk/sqldb"
	"github.com/warlck/food-flow/business/types/money"
)

func toSafeMoney(val float64) money.Money {
	if val < 0 || math.IsNaN(val) || math.IsInf(val, 0) {
		val = 0
	}
	if val > 1_000_000_000 {
		val = 1_000_000_000
	}
	return money.MustParse(math.Round(val*100) / 100)
}

// applyInsightsFilter appends WHERE constraints for insights queries.
func (s *Store) applyInsightsFilter(filter orderbus.InsightsFilter, data map[string]any, buf *bytes.Buffer, prefix string) {
	var wc []string

	if filter.RestaurantID != nil && *filter.RestaurantID != "" {
		data["insights_restaurant_id"] = *filter.RestaurantID
		wc = append(wc, fmt.Sprintf("%srestaurant_id = :insights_restaurant_id", prefix))
	}

	if filter.StartDate != nil {
		data["insights_start_date"] = (*filter.StartDate).UTC()
		wc = append(wc, fmt.Sprintf("%sdate_created >= :insights_start_date", prefix))
	}

	if filter.EndDate != nil {
		data["insights_end_date"] = (*filter.EndDate).UTC()
		wc = append(wc, fmt.Sprintf("%sdate_created <= :insights_end_date", prefix))
	}

	if len(wc) > 0 {
		buf.WriteString(" WHERE ")
		for i, w := range wc {
			if i > 0 {
				buf.WriteString(" AND ")
			}
			buf.WriteString(w)
		}
	}
}

// QuerySalesSummary aggregates high-level sales KPI metrics.
func (s *Store) QuerySalesSummary(ctx context.Context, filter orderbus.InsightsFilter) (orderbus.SalesSummary, error) {
	data := map[string]any{}

	const q = `
	SELECT
		COALESCE(COUNT(1), 0) AS total_orders,
		COALESCE(COUNT(CASE WHEN o.order_status = 'completed' THEN 1 END), 0) AS completed_orders,
		COALESCE(COUNT(CASE WHEN o.order_status = 'cancelled' THEN 1 END), 0) AS cancelled_orders,
		COALESCE(SUM(CASE WHEN o.order_status != 'cancelled' THEN o.subtotal ELSE 0 END), 0) AS gross_sales,
		COALESCE(SUM(CASE WHEN o.order_status != 'cancelled' THEN GREATEST(0, (o.subtotal - o.discount)) ELSE 0 END), 0) AS net_sales,
		COALESCE(SUM(CASE WHEN o.order_status != 'cancelled' THEN o.discount ELSE 0 END), 0) AS total_discounts,
		COALESCE(SUM(CASE WHEN o.order_status != 'cancelled' THEN o.delivery_fee ELSE 0 END), 0) AS total_delivery_fees,
		COALESCE(SUM(CASE WHEN o.order_status != 'cancelled' THEN o.tax ELSE 0 END), 0) AS total_tax,
		COALESCE(SUM(CASE WHEN o.order_status != 'cancelled' THEN o.total ELSE 0 END), 0) AS total_collected
	FROM
		orders AS o`

	buf := bytes.NewBufferString(q)
	s.applyInsightsFilter(filter, data, buf, "o.")

	var row struct {
		TotalOrders       int     `db:"total_orders"`
		CompletedOrders   int     `db:"completed_orders"`
		CancelledOrders   int     `db:"cancelled_orders"`
		GrossSales        float64 `db:"gross_sales"`
		NetSales          float64 `db:"net_sales"`
		TotalDiscounts    float64 `db:"total_discounts"`
		TotalDeliveryFees float64 `db:"total_delivery_fees"`
		TotalTax          float64 `db:"total_tax"`
		TotalCollected    float64 `db:"total_collected"`
	}

	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, buf.String(), data, &row); err != nil {
		return orderbus.SalesSummary{}, fmt.Errorf("query sales summary: %w", err)
	}

	var aov float64
	validOrders := row.TotalOrders - row.CancelledOrders
	if validOrders > 0 {
		aov = row.TotalCollected / float64(validOrders)
	}

	return orderbus.SalesSummary{
		GrossSales:        toSafeMoney(row.GrossSales),
		NetSales:          toSafeMoney(row.NetSales),
		TotalOrders:       row.TotalOrders,
		CompletedOrders:   row.CompletedOrders,
		CancelledOrders:   row.CancelledOrders,
		AverageOrderValue: toSafeMoney(aov),
		TotalDiscounts:    toSafeMoney(row.TotalDiscounts),
		TotalDeliveryFees: toSafeMoney(row.TotalDeliveryFees),
		TotalTax:          toSafeMoney(row.TotalTax),
		TotalCollected:    toSafeMoney(row.TotalCollected),
	}, nil
}

// QuerySalesOverTime returns time-series daily sales metrics.
func (s *Store) QuerySalesOverTime(ctx context.Context, filter orderbus.InsightsFilter) ([]orderbus.TimeSeriesPoint, error) {
	data := map[string]any{}

	const q = `
	SELECT
		TO_CHAR(o.date_created, 'YYYY-MM-DD') AS date,
		COALESCE(SUM(CASE WHEN o.order_status != 'cancelled' THEN o.subtotal ELSE 0 END), 0) AS gross_sales,
		COALESCE(SUM(CASE WHEN o.order_status != 'cancelled' THEN GREATEST(0, (o.subtotal - o.discount)) ELSE 0 END), 0) AS net_sales,
		COALESCE(SUM(CASE WHEN o.order_status != 'cancelled' THEN o.total ELSE 0 END), 0) AS total_collected,
		COALESCE(COUNT(1), 0) AS order_count
	FROM
		orders AS o`

	buf := bytes.NewBufferString(q)
	s.applyInsightsFilter(filter, data, buf, "o.")
	buf.WriteString(" GROUP BY TO_CHAR(o.date_created, 'YYYY-MM-DD') ORDER BY date ASC")

	var rows []struct {
		Date           string  `db:"date"`
		GrossSales     float64 `db:"gross_sales"`
		NetSales       float64 `db:"net_sales"`
		TotalCollected float64 `db:"total_collected"`
		OrderCount     int     `db:"order_count"`
	}

	if err := sqldb.NamedQuerySlice(ctx, s.log, s.db, buf.String(), data, &rows); err != nil {
		return nil, fmt.Errorf("query sales over time: %w", err)
	}

	points := make([]orderbus.TimeSeriesPoint, len(rows))
	for i, r := range rows {
		var aov float64
		if r.OrderCount > 0 {
			aov = r.TotalCollected / float64(r.OrderCount)
		}
		points[i] = orderbus.TimeSeriesPoint{
			Date:           r.Date,
			GrossSales:     toSafeMoney(r.GrossSales),
			NetSales:       toSafeMoney(r.NetSales),
			TotalCollected: toSafeMoney(r.TotalCollected),
			OrderCount:     r.OrderCount,
			AverageOrder:   toSafeMoney(aov),
		}
	}

	return points, nil
}

// QueryTopItemSales returns best-selling menu items by quantity and revenue with limit.
func (s *Store) QueryTopItemSales(ctx context.Context, filter orderbus.InsightsFilter, limit int) ([]orderbus.ItemSalesMetric, error) {
	data := map[string]any{
		"insights_limit": limit,
	}

	const q = `
	SELECT
		oi.menu_item_id,
		oi.menu_item_name,
		COALESCE(SUM(oi.quantity), 0) AS quantity_sold,
		COALESCE(SUM(oi.menu_item_price * oi.quantity), 0) AS total_revenue
	FROM
		order_items AS oi
	JOIN
		orders AS o ON o.order_id = oi.order_id`

	buf := bytes.NewBufferString(q)
	s.applyInsightsFilter(filter, data, buf, "o.")
	if len(data) == 1 {
		buf.WriteString(" WHERE o.order_status != 'cancelled'")
	} else {
		buf.WriteString(" AND o.order_status != 'cancelled'")
	}
	buf.WriteString(" GROUP BY oi.menu_item_id, oi.menu_item_name ORDER BY total_revenue DESC LIMIT :insights_limit")

	var rows []struct {
		MenuItemID   uuid.UUID `db:"menu_item_id"`
		MenuItemName string    `db:"menu_item_name"`
		QuantitySold int       `db:"quantity_sold"`
		TotalRevenue float64   `db:"total_revenue"`
	}

	if err := sqldb.NamedQuerySlice(ctx, s.log, s.db, buf.String(), data, &rows); err != nil {
		return nil, fmt.Errorf("query top item sales: %w", err)
	}

	items := make([]orderbus.ItemSalesMetric, len(rows))
	for i, r := range rows {
		items[i] = orderbus.ItemSalesMetric{
			MenuItemID:   r.MenuItemID,
			MenuItemName: r.MenuItemName,
			QuantitySold: r.QuantitySold,
			TotalRevenue: toSafeMoney(r.TotalRevenue),
		}
	}

	return items, nil
}

// QueryAllItemSales returns all sold menu items for the period (for category revenue aggregation).
func (s *Store) QueryAllItemSales(ctx context.Context, filter orderbus.InsightsFilter) ([]orderbus.ItemSalesMetric, error) {
	data := map[string]any{}

	const q = `
	SELECT
		oi.menu_item_id,
		oi.menu_item_name,
		COALESCE(SUM(oi.quantity), 0) AS quantity_sold,
		COALESCE(SUM(oi.menu_item_price * oi.quantity), 0) AS total_revenue
	FROM
		order_items AS oi
	JOIN
		orders AS o ON o.order_id = oi.order_id`

	buf := bytes.NewBufferString(q)
	s.applyInsightsFilter(filter, data, buf, "o.")
	if len(data) == 0 {
		buf.WriteString(" WHERE o.order_status != 'cancelled'")
	} else {
		buf.WriteString(" AND o.order_status != 'cancelled'")
	}
	buf.WriteString(" GROUP BY oi.menu_item_id, oi.menu_item_name ORDER BY total_revenue DESC")

	var rows []struct {
		MenuItemID   uuid.UUID `db:"menu_item_id"`
		MenuItemName string    `db:"menu_item_name"`
		QuantitySold int       `db:"quantity_sold"`
		TotalRevenue float64   `db:"total_revenue"`
	}

	if err := sqldb.NamedQuerySlice(ctx, s.log, s.db, buf.String(), data, &rows); err != nil {
		return nil, fmt.Errorf("query all item sales: %w", err)
	}

	items := make([]orderbus.ItemSalesMetric, len(rows))
	for i, r := range rows {
		items[i] = orderbus.ItemSalesMetric{
			MenuItemID:   r.MenuItemID,
			MenuItemName: r.MenuItemName,
			QuantitySold: r.QuantitySold,
			TotalRevenue: toSafeMoney(r.TotalRevenue),
		}
	}

	return items, nil
}

// QueryTopAddonSales returns popular upsell add-on performance.
func (s *Store) QueryTopAddonSales(ctx context.Context, filter orderbus.InsightsFilter, limit int) ([]orderbus.TopAddonMetric, error) {
	data := map[string]any{
		"insights_limit": limit,
	}

	const q = `
	SELECT
		oia.addon_id,
		oia.addon_name,
		COALESCE(SUM(oia.quantity * oi.quantity), 0) AS quantity_sold,
		COALESCE(SUM(oia.addon_price * oia.quantity * oi.quantity), 0) AS total_revenue
	FROM
		order_item_addons AS oia
	JOIN
		order_items AS oi ON oi.order_item_id = oia.order_item_id
	JOIN
		orders AS o ON o.order_id = oi.order_id`

	buf := bytes.NewBufferString(q)
	s.applyInsightsFilter(filter, data, buf, "o.")
	if len(data) == 1 {
		buf.WriteString(" WHERE o.order_status != 'cancelled'")
	} else {
		buf.WriteString(" AND o.order_status != 'cancelled'")
	}
	buf.WriteString(" GROUP BY oia.addon_id, oia.addon_name ORDER BY total_revenue DESC LIMIT :insights_limit")

	var rows []struct {
		AddonID      uuid.UUID `db:"addon_id"`
		AddonName    string    `db:"addon_name"`
		QuantitySold int       `db:"quantity_sold"`
		TotalRevenue float64   `db:"total_revenue"`
	}

	if err := sqldb.NamedQuerySlice(ctx, s.log, s.db, buf.String(), data, &rows); err != nil {
		return nil, fmt.Errorf("query top addons: %w", err)
	}

	addons := make([]orderbus.TopAddonMetric, len(rows))
	for i, r := range rows {
		addons[i] = orderbus.TopAddonMetric{
			AddonID:      r.AddonID,
			AddonName:    r.AddonName,
			QuantitySold: r.QuantitySold,
			TotalRevenue: toSafeMoney(r.TotalRevenue),
		}
	}

	return addons, nil
}

// QueryOrderTypes returns distribution of fulfillment channels.
func (s *Store) QueryOrderTypes(ctx context.Context, filter orderbus.InsightsFilter) ([]orderbus.OrderTypeMetric, error) {
	data := map[string]any{}

	const q = `
	SELECT
		o.order_type,
		COALESCE(COUNT(1), 0) AS count,
		COALESCE(SUM(CASE WHEN o.order_status != 'cancelled' THEN o.total ELSE 0 END), 0) AS total_revenue
	FROM
		orders AS o`

	buf := bytes.NewBufferString(q)
	s.applyInsightsFilter(filter, data, buf, "o.")
	buf.WriteString(" GROUP BY o.order_type ORDER BY count DESC")

	var rows []struct {
		OrderType    string  `db:"order_type"`
		Count        int     `db:"count"`
		TotalRevenue float64 `db:"total_revenue"`
	}

	if err := sqldb.NamedQuerySlice(ctx, s.log, s.db, buf.String(), data, &rows); err != nil {
		return nil, fmt.Errorf("query order types: %w", err)
	}

	var sumCount int
	for _, r := range rows {
		sumCount += r.Count
	}

	types := make([]orderbus.OrderTypeMetric, len(rows))
	for i, r := range rows {
		var pct float64
		if sumCount > 0 {
			pct = (float64(r.Count) / float64(sumCount)) * 100
		}
		types[i] = orderbus.OrderTypeMetric{
			OrderType:    r.OrderType,
			Count:        r.Count,
			TotalRevenue: toSafeMoney(r.TotalRevenue),
			Percentage:   pct,
		}
	}

	return types, nil
}

// QueryPeakHours returns order volume and sales distribution across 24 hours.
func (s *Store) QueryPeakHours(ctx context.Context, filter orderbus.InsightsFilter) ([]orderbus.HourlyMetric, error) {
	data := map[string]any{}

	const q = `
	SELECT
		CAST(EXTRACT(HOUR FROM o.date_created) AS INT) AS hour,
		COALESCE(COUNT(1), 0) AS count,
		COALESCE(SUM(CASE WHEN o.order_status != 'cancelled' THEN o.total ELSE 0 END), 0) AS total_revenue
	FROM
		orders AS o`

	buf := bytes.NewBufferString(q)
	s.applyInsightsFilter(filter, data, buf, "o.")
	buf.WriteString(" GROUP BY EXTRACT(HOUR FROM o.date_created) ORDER BY hour ASC")

	var rows []struct {
		Hour         int     `db:"hour"`
		Count        int     `db:"count"`
		TotalRevenue float64 `db:"total_revenue"`
	}

	if err := sqldb.NamedQuerySlice(ctx, s.log, s.db, buf.String(), data, &rows); err != nil {
		return nil, fmt.Errorf("query peak hours: %w", err)
	}

	hours := make([]orderbus.HourlyMetric, len(rows))
	for i, r := range rows {
		hours[i] = orderbus.HourlyMetric{
			Hour:         r.Hour,
			Count:        r.Count,
			TotalRevenue: toSafeMoney(r.TotalRevenue),
		}
	}

	return hours, nil
}
