package response

import (
	"time"

	"ims-intro/pkg/service/dto"
)

type UserAnalyticsSummaryResponse struct {
	TotalOrders   int     `json:"totalOrders"`
	TotalSpent    float64 `json:"totalSpent"`
	PendingOrders int     `json:"pendingOrders"`
}

type UserRecentOrderResponse struct {
	OrderID      int       `json:"orderId"`
	ProductNames string    `json:"productNames"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"createdAt"`
	TotalPrice   float64   `json:"totalPrice"`
}

type UserTopProductResponse struct {
	ProductID     int64  `json:"productId"`
	Name          string `json:"name"`
	PurchaseCount int    `json:"purchaseCount"`
}

type UserRecommendationResponse struct {
	ProductID int64   `json:"productId"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	Category  string  `json:"category"`
	Reason    string  `json:"reason"`
}

func ToUserAnalyticsSummaryResponse(summary dto.UserAnalyticsSummary) UserAnalyticsSummaryResponse {
	return UserAnalyticsSummaryResponse{
		TotalOrders:   summary.TotalOrders,
		TotalSpent:    summary.TotalSpent,
		PendingOrders: summary.PendingOrders,
	}
}

func ToUserRecentOrderResponse(items []dto.UserRecentOrder) []UserRecentOrderResponse {
	result := make([]UserRecentOrderResponse, 0, len(items))
	for _, item := range items {
		result = append(result, UserRecentOrderResponse{
			OrderID:      item.OrderID,
			ProductNames: item.ProductNames,
			Status:       item.Status,
			CreatedAt:    item.CreatedAt,
			TotalPrice:   item.TotalPrice,
		})
	}
	return result
}

func ToUserTopProductResponse(items []dto.UserTopProduct) []UserTopProductResponse {
	result := make([]UserTopProductResponse, 0, len(items))
	for _, item := range items {
		result = append(result, UserTopProductResponse{
			ProductID:     item.ProductID,
			Name:          item.Name,
			PurchaseCount: item.PurchaseCount,
		})
	}
	return result
}

func ToUserRecommendationResponse(items []dto.UserRecommendation) []UserRecommendationResponse {
	result := make([]UserRecommendationResponse, 0, len(items))
	for _, item := range items {
		result = append(result, UserRecommendationResponse{
			ProductID: item.ProductID,
			Name:      item.Name,
			Price:     item.Price,
			Category:  item.Category,
			Reason:    item.Reason,
		})
	}
	return result
}
