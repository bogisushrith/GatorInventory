package response

import "ims-intro/pkg/service/dto"

type SummaryResponse struct {
	TotalRevenue  float64 `json:"totalRevenue"`
	TotalOrders   int     `json:"totalOrders"`
	TotalProducts int     `json:"totalProducts"`
	LowStockCount int     `json:"lowStockCount"`
}

type RevenuePoint struct {
	Date  string  `json:"date"`
	Value float64 `json:"value"`
}

type OrderPoint struct {
	Date  string `json:"date"`
	Value int    `json:"value"`
}

type ProductAnalytics struct {
	ProductID int64   `json:"productId"`
	Name      string  `json:"name"`
	TotalSold int     `json:"totalSold"`
	Revenue   float64 `json:"revenue"`
}

type LowStockProductResponse struct {
	ProductID int64  `json:"productId"`
	Name      string `json:"name"`
	Quantity  int    `json:"quantity"`
	Category  string `json:"category"`
}

func ToSummaryResponse(summary dto.AnalyticsSummary) SummaryResponse {
	return SummaryResponse{
		TotalRevenue:  summary.TotalRevenue,
		TotalOrders:   summary.TotalOrders,
		TotalProducts: summary.TotalProducts,
		LowStockCount: summary.LowStockCount,
	}
}

func ToRevenuePointResponse(points []dto.RevenueTrendPoint) []RevenuePoint {
	result := make([]RevenuePoint, 0, len(points))
	for _, point := range points {
		result = append(result, RevenuePoint{Date: point.Date, Value: point.Value})
	}
	return result
}

func ToOrderPointResponse(points []dto.OrderTrendPoint) []OrderPoint {
	result := make([]OrderPoint, 0, len(points))
	for _, point := range points {
		result = append(result, OrderPoint{Date: point.Date, Value: point.Value})
	}
	return result
}

func ToTopProductAnalyticsResponse(products []dto.TopProductMetric) []ProductAnalytics {
	result := make([]ProductAnalytics, 0, len(products))
	for _, item := range products {
		result = append(result, ProductAnalytics{
			ProductID: item.ProductID,
			Name:      item.Name,
			TotalSold: item.TotalSold,
			Revenue:   item.Revenue,
		})
	}
	return result
}

func ToLowStockProductResponse(products []dto.LowStockProduct) []LowStockProductResponse {
	result := make([]LowStockProductResponse, 0, len(products))
	for _, item := range products {
		result = append(result, LowStockProductResponse{
			ProductID: item.ID,
			Name:      item.Name,
			Quantity:  item.Quantity,
			Category:  item.Category,
		})
	}
	return result
}
