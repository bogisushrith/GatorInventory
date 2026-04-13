package request

type AddCartItemRequest struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

type UpdateCartItemRequest struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

type RemoveCartItemRequest struct {
	ProductID int `json:"product_id"`
}
