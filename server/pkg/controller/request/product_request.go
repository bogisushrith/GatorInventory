package request

import "ims-intro/pkg/service/dto"

type AddProductRequest struct {
	Name     string  `json:"name"`
	Price    float32 `json:"price"`
	Quantity int     `json:"quantity"`
	Category string  `json:"category"`
	ImageURL string  `json:"image_url"`
}

type UpdateProductStockRequest struct {
	Quantity int `json:"quantity"`
}

func (request *AddProductRequest) ToModel() *dto.ProductCreate {
	return &dto.ProductCreate{
		Name:     request.Name,
		Price:    request.Price,
		Quantity: request.Quantity,
		Category: request.Category,
		ImageURL: request.ImageURL,
	}
}
