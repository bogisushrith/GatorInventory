package request

import "ims-intro/pkg/service/dto"

type CreateOrderRequest struct {
	Items []CreateOrderItemRequest `json:"items"`
}

type CreateOrderItemRequest struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

func (request *CreateOrderRequest) ToModel() *dto.OrderCreate {
	items := make([]dto.OrderItemCreate, 0, len(request.Items))
	for _, item := range request.Items {
		items = append(items, dto.OrderItemCreate{ProductID: item.ProductID, Quantity: item.Quantity})
	}
	return &dto.OrderCreate{Items: items}
}
