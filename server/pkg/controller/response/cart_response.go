package response

import "ims-intro/pkg/domain"

type CartItemResponse struct {
	ID           int64   `json:"id"`
	UserID       int64   `json:"user_id"`
	ProductID    int64   `json:"product_id"`
	Quantity     int     `json:"quantity"`
	ProductName  string  `json:"product_name"`
	ProductPrice float32 `json:"product_price"`
	ProductStock int     `json:"product_stock"`
}

func ToCartItemResponse(item domain.CartItem) CartItemResponse {
	return CartItemResponse{
		ID:           item.ID,
		UserID:       item.UserID,
		ProductID:    item.ProductID,
		Quantity:     item.Quantity,
		ProductName:  item.ProductName,
		ProductPrice: item.ProductPrice,
		ProductStock: item.ProductStock,
	}
}

func ToCartItemResponseList(items []domain.CartItem) []CartItemResponse {
	responses := make([]CartItemResponse, 0, len(items))
	for _, item := range items {
		responses = append(responses, ToCartItemResponse(item))
	}
	return responses
}
