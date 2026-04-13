package domain

type CartItem struct {
	ID           int64
	UserID       int64
	ProductID    int64
	Quantity     int
	ProductName  string
	ProductPrice float32
	ProductStock int
}
