package orderapp

import (
	"context"
	"errors"
	"fmt"
	"math"
	"net/http"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/warlck/food-flow/app/sdk/errs"
	"github.com/warlck/food-flow/app/sdk/mid"
	"github.com/warlck/food-flow/business/domain/orderbus"
	"github.com/warlck/food-flow/business/domain/restaurantbus"
	"github.com/warlck/food-flow/business/sdk/page"
	"github.com/warlck/food-flow/foundation/web"
)

// AppSalesSummary represents the HTTP DTO for sales KPIs.
type AppSalesSummary struct {
	GrossSales        float64 `json:"grossSales"`
	NetSales          float64 `json:"netSales"`
	TotalOrders       int     `json:"totalOrders"`
	CompletedOrders   int     `json:"completedOrders"`
	CancelledOrders   int     `json:"cancelledOrders"`
	AverageOrderValue float64 `json:"averageOrderValue"`
	TotalDiscounts    float64 `json:"totalDiscounts"`
	TotalDeliveryFees float64 `json:"totalDeliveryFees"`
	TotalTax          float64 `json:"totalTax"`
	TotalCollected    float64 `json:"totalCollected"`
}

// AppTimeSeriesPoint represents daily revenue and order metrics.
type AppTimeSeriesPoint struct {
	Date           string  `json:"date"`
	GrossSales     float64 `json:"grossSales"`
	NetSales       float64 `json:"netSales"`
	TotalCollected float64 `json:"totalCollected"`
	OrderCount     int     `json:"orderCount"`
	AverageOrder   float64 `json:"averageOrder"`
}

// AppTopItemMetric represents best-selling menu items.
type AppTopItemMetric struct {
	MenuItemID   string  `json:"menuItemId"`
	MenuItemName string  `json:"menuItemName"`
	CategoryName string  `json:"categoryName"`
	QuantitySold int     `json:"quantitySold"`
	TotalRevenue float64 `json:"totalRevenue"`
}

// AppTopCategoryMetric represents category performance breakdown.
type AppTopCategoryMetric struct {
	CategoryID   string  `json:"categoryId"`
	CategoryName string  `json:"categoryName"`
	QuantitySold int     `json:"quantitySold"`
	TotalRevenue float64 `json:"totalRevenue"`
	Percentage   float64 `json:"percentage"`
}

// AppTopAddonMetric represents add-on modifier metrics.
type AppTopAddonMetric struct {
	AddonID      string  `json:"addonId"`
	AddonName    string  `json:"addonName"`
	QuantitySold int     `json:"quantitySold"`
	TotalRevenue float64 `json:"totalRevenue"`
}

// AppOrderTypeMetric represents fulfillment type share.
type AppOrderTypeMetric struct {
	OrderType    string  `json:"orderType"`
	Count        int     `json:"count"`
	TotalRevenue float64 `json:"totalRevenue"`
	Percentage   float64 `json:"percentage"`
}

// AppHourlyMetric represents 24-hour order volume and sales distribution.
type AppHourlyMetric struct {
	Hour         int     `json:"hour"`
	Count        int     `json:"count"`
	TotalRevenue float64 `json:"totalRevenue"`
}

// AppInsights represents the comprehensive analytics response payload.
type AppInsights struct {
	Summary       AppSalesSummary        `json:"summary"`
	SalesOverTime []AppTimeSeriesPoint   `json:"salesOverTime"`
	TopItems      []AppTopItemMetric     `json:"topItems"`
	TopCategories []AppTopCategoryMetric `json:"topCategories"`
	TopAddons     []AppTopAddonMetric    `json:"topAddons"`
	OrderTypes    []AppOrderTypeMetric   `json:"orderTypes"`
	PeakHours     []AppHourlyMetric      `json:"peakHours"`
}

// queryInsights handles GET /v1/insights.
func (a *app) queryInsights(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	qp := r.URL.Query()

	claims := mid.GetClaims(ctx)
	var filter orderbus.InsightsFilter
	var authorizedRestaurantIDs []uuid.UUID

	if rid := qp.Get("restaurant_id"); rid != "" {
		parsedID, err := uuid.Parse(rid)
		if err != nil {
			return errs.New(errs.InvalidArgument, fmt.Errorf("invalid restaurant_id: %w", err))
		}

		// Enforce tenant authorization: caller must belong to the restaurant's organization
		rest, err := a.restaurantBus.QueryByID(ctx, parsedID)
		if err != nil {
			if errors.Is(err, restaurantbus.ErrNotFound) {
				return errs.New(errs.NotFound, fmt.Errorf("restaurant not found"))
			}
			return fmt.Errorf("query restaurant: %w", err)
		}

		if !claims.IsOrgAuthorized(rest.OrganizationID) {
			return errs.New(errs.PermissionDenied, fmt.Errorf("user not authorized for this organization"))
		}

		filter.RestaurantIDs = []uuid.UUID{parsedID}
		authorizedRestaurantIDs = []uuid.UUID{parsedID}
	} else {
		// When no restaurant is specified, scope strictly to all restaurants in caller's authorized organization(s)
		for _, orgIDStr := range claims.OrganizationIDs {
			orgID, err := uuid.Parse(orgIDStr)
			if err != nil {
				continue
			}
			rests, err := a.restaurantBus.Query(ctx, restaurantbus.QueryFilter{
				OrganizationID: &orgID,
			}, restaurantbus.DefaultOrderBy, page.MustParse("1", "100"))
			if err != nil {
				return fmt.Errorf("query restaurants for org: %w", err)
			}
			for _, r := range rests {
				authorizedRestaurantIDs = append(authorizedRestaurantIDs, r.ID)
			}
		}

		if len(authorizedRestaurantIDs) == 0 {
			return web.Respond(ctx, w, AppInsights{
				Summary:       AppSalesSummary{},
				SalesOverTime: []AppTimeSeriesPoint{},
				TopItems:      []AppTopItemMetric{},
				TopCategories: []AppTopCategoryMetric{},
				TopAddons:     []AppTopAddonMetric{},
				OrderTypes:    []AppOrderTypeMetric{},
				PeakHours:     []AppHourlyMetric{},
			}, http.StatusOK)
		}

		filter.RestaurantIDs = authorizedRestaurantIDs
	}

	if sd := qp.Get("start_date"); sd != "" {
		parsedDate, err := time.Parse(time.RFC3339, sd)
		if err != nil {
			// Fallback to YYYY-MM-DD (start of day UTC)
			parsedDate, err = time.Parse("2006-01-02", sd)
			if err != nil {
				return errs.New(errs.InvalidArgument, fmt.Errorf("invalid start_date format: %w", err))
			}
			parsedDate = time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(), 0, 0, 0, 0, time.UTC)
		}
		filter.StartDate = &parsedDate
	}

	if ed := qp.Get("end_date"); ed != "" {
		parsedDate, err := time.Parse(time.RFC3339, ed)
		if err != nil {
			// Fallback to YYYY-MM-DD (end of day 23:59:59.999999999 UTC)
			parsedDate, err = time.Parse("2006-01-02", ed)
			if err != nil {
				return errs.New(errs.InvalidArgument, fmt.Errorf("invalid end_date format: %w", err))
			}
			parsedDate = time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(), 23, 59, 59, 999999999, time.UTC)
		}
		filter.EndDate = &parsedDate
	}

	// 1. Query pure order domain metrics
	metrics, err := a.orderBus.QueryOrderMetrics(ctx, filter)
	if err != nil {
		return fmt.Errorf("query order metrics: %w", err)
	}

	// 2. Build Top Items with Category Name from snapshots
	appTopItems := make([]AppTopItemMetric, len(metrics.TopItems))
	for i, item := range metrics.TopItems {
		catName := item.CategoryName
		if catName == "" {
			catName = "General"
		}

		appTopItems[i] = AppTopItemMetric{
			MenuItemID:   item.MenuItemID.String(),
			MenuItemName: item.MenuItemName,
			CategoryName: catName,
			QuantitySold: item.QuantitySold,
			TotalRevenue: item.TotalRevenue.Value(),
		}
	}

	// 3. Category Aggregation from All Item Sales snapshots
	type catAccumulator struct {
		id       uuid.UUID
		name     string
		quantity int
		revenue  float64
	}

	catTotals := make(map[uuid.UUID]*catAccumulator)
	var totalCategoryRevenue float64

	for _, sale := range metrics.AllItemSales {
		catID := sale.CategoryID
		catName := sale.CategoryName
		if catName == "" {
			catName = "General"
		}

		rev := sale.TotalRevenue.Value()
		totalCategoryRevenue += rev

		if agg, exists := catTotals[catID]; exists {
			agg.quantity += sale.QuantitySold
			agg.revenue += rev
		} else {
			catTotals[catID] = &catAccumulator{
				id:       catID,
				name:     catName,
				quantity: sale.QuantitySold,
				revenue:  rev,
			}
		}
	}

	var appTopCategories []AppTopCategoryMetric
	for _, agg := range catTotals {
		var pct float64
		if totalCategoryRevenue > 0 {
			pct = math.Round((agg.revenue/totalCategoryRevenue)*10000) / 100
		}
		appTopCategories = append(appTopCategories, AppTopCategoryMetric{
			CategoryID:   agg.id.String(),
			CategoryName: agg.name,
			QuantitySold: agg.quantity,
			TotalRevenue: agg.revenue,
			Percentage:   pct,
		})
	}

	// Sort categories by revenue descending
	sort.Slice(appTopCategories, func(i, j int) bool {
		return appTopCategories[i].TotalRevenue > appTopCategories[j].TotalRevenue
	})

	if len(appTopCategories) > 10 {
		appTopCategories = appTopCategories[:10]
	}

	// 6. Map other metrics to App DTOs
	appTimeSeries := make([]AppTimeSeriesPoint, len(metrics.SalesOverTime))
	for i, p := range metrics.SalesOverTime {
		appTimeSeries[i] = AppTimeSeriesPoint{
			Date:           p.Date,
			GrossSales:     p.GrossSales.Value(),
			NetSales:       p.NetSales.Value(),
			TotalCollected: p.TotalCollected.Value(),
			OrderCount:     p.OrderCount,
			AverageOrder:   p.AverageOrder.Value(),
		}
	}

	appTopAddons := make([]AppTopAddonMetric, len(metrics.TopAddons))
	for i, a := range metrics.TopAddons {
		appTopAddons[i] = AppTopAddonMetric{
			AddonID:      a.AddonID.String(),
			AddonName:    a.AddonName,
			QuantitySold: a.QuantitySold,
			TotalRevenue: a.TotalRevenue.Value(),
		}
	}

	appOrderTypes := make([]AppOrderTypeMetric, len(metrics.OrderTypes))
	for i, ot := range metrics.OrderTypes {
		appOrderTypes[i] = AppOrderTypeMetric{
			OrderType:    ot.OrderType,
			Count:        ot.Count,
			TotalRevenue: ot.TotalRevenue.Value(),
			Percentage:   ot.Percentage,
		}
	}

	appPeakHours := make([]AppHourlyMetric, len(metrics.PeakHours))
	for i, h := range metrics.PeakHours {
		appPeakHours[i] = AppHourlyMetric{
			Hour:         h.Hour,
			Count:        h.Count,
			TotalRevenue: h.TotalRevenue.Value(),
		}
	}

	response := AppInsights{
		Summary: AppSalesSummary{
			GrossSales:        metrics.Summary.GrossSales.Value(),
			NetSales:          metrics.Summary.NetSales.Value(),
			TotalOrders:       metrics.Summary.TotalOrders,
			CompletedOrders:   metrics.Summary.CompletedOrders,
			CancelledOrders:   metrics.Summary.CancelledOrders,
			AverageOrderValue: metrics.Summary.AverageOrderValue.Value(),
			TotalDiscounts:    metrics.Summary.TotalDiscounts.Value(),
			TotalDeliveryFees: metrics.Summary.TotalDeliveryFees.Value(),
			TotalTax:          metrics.Summary.TotalTax.Value(),
			TotalCollected:    metrics.Summary.TotalCollected.Value(),
		},
		SalesOverTime: appTimeSeries,
		TopItems:      appTopItems,
		TopCategories: appTopCategories,
		TopAddons:     appTopAddons,
		OrderTypes:    appOrderTypes,
		PeakHours:     appPeakHours,
	}

	return web.Respond(ctx, w, response, http.StatusOK)
}
