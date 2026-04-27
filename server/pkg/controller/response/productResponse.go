package response

import "ims-intro/pkg/domain"

type ProductResponse struct {
	Id       int64   `json:"id"`
	Name     string  `json:"name"`
	Price    float32 `json:"price"`
	Quantity int     `json:"quantity"`
	Category string  `json:"category"`
	ImageURL string  `json:"image_url"`
}

type ProductPaginationResponse struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

type ProductListResponse struct {
	Data       []*ProductResponse        `json:"data"`
	Pagination ProductPaginationResponse `json:"pagination"`
}

func toProductResponse(product *domain.Product) *ProductResponse {
	return &ProductResponse{
		Id:       product.Id,
		Name:     product.Name,
		Price:    product.Price,
		Quantity: product.Quantity,
		Category: product.Category,
		ImageURL: product.ImageURL,
	}
}

func ToProductResponse(product *domain.Product) *ProductResponse {
	return toProductResponse(product)
}

func ToProductResponseList(products []*domain.Product) []*ProductResponse {
	var responses []*ProductResponse
	for _, product := range products {
		responses = append(responses, toProductResponse(product))
	}
	return responses
}

func ToPaginatedProductResponse(products []*domain.Product, page, limit int, total int64, totalPages int) *ProductListResponse {
	return &ProductListResponse{
		Data: ToProductResponseList(products),
		Pagination: ProductPaginationResponse{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}
}
