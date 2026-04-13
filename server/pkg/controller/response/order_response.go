package response

import (
	"time"

	"ims-intro/pkg/domain"
)

type CreateOrderResponse struct {
	OrderID int `json:"order_id"`
}

type OrderItemResponse struct {
	ID        int `json:"id"`
	OrderID   int `json:"order_id"`
	ProductID int `json:"product_id"`
	ProductName string `json:"product_name"`
	ProductPrice float32 `json:"product_price"`
	Quantity  int `json:"quantity"`
}

type OrderResponse struct {
	ID        int                 `json:"id"`
	UserID    int64               `json:"user_id"`
	UserName  string              `json:"user_name"`
	Status    string              `json:"status"`
	CreatedAt time.Time           `json:"created_at"`
	Items     []OrderItemResponse `json:"items"`
}

func ToOrderItemResponse(item domain.OrderItem) OrderItemResponse {
	return OrderItemResponse{ID: item.ID, OrderID: item.OrderID, ProductID: item.ProductID, ProductName: item.ProductName, ProductPrice: item.ProductPrice, Quantity: item.Quantity}
}

func ToOrderResponse(order domain.Order) OrderResponse {
	items := make([]OrderItemResponse, 0, len(order.Items))
	for _, item := range order.Items {
		items = append(items, ToOrderItemResponse(item))
	}
	return OrderResponse{ID: order.ID, UserID: order.UserID, UserName: order.UserName, Status: order.Status, CreatedAt: order.CreatedAt, Items: items}
}

func ToOrderResponseList(orders []domain.Order) []OrderResponse {
	orderResponses := make([]OrderResponse, 0, len(orders))
	for _, order := range orders {
		orderResponses = append(orderResponses, ToOrderResponse(order))
	}
	return orderResponses
}
