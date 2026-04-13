package domain

import "time"

type Order struct {
	ID        int
	UserID    int64
	UserName  string
	Status    string
	CreatedAt time.Time
	Items     []OrderItem
}

type OrderItem struct {
	ID        int
	OrderID   int
	ProductID int
	ProductName string
	ProductPrice float32
	Quantity  int
}
