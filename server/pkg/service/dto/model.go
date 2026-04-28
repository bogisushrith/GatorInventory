package dto

import "time"

type UserCreate struct {
	Username string
	Email    string
	Password string
	Role     string
}

type ProductCreate struct {
	Name     string
	Price    float32
	Quantity int
	Category string
	ImageURL string
}

type ProductListQuery struct {
	Page     int
	Limit    int
	Search   string
	Category string
	MinPrice *float64
	MaxPrice *float64
}

type UserSummary struct {
	ID       int64
	Username string
	Email    string
	Role     string
}

type OrderItemCreate struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

type OrderCreate struct {
	Items []OrderItemCreate `json:"items"`
}

type OrderListQuery struct {
	UserID   int64
	Role     string
	Search   string
	Status   string
	DateFrom *time.Time
	DateTo   *time.Time
}

type AnalyticsQuery struct {
	Days               int
	LowStockThreshold  int
	TopProductsLimit   int
	LowStockItemsLimit int
}

type AnalyticsSummary struct {
	TotalRevenue  float64
	TotalOrders   int
	TotalProducts int
	LowStockCount int
}

type RevenueTrendPoint struct {
	Date  string
	Value float64
}

type OrderTrendPoint struct {
	Date  string
	Value int
}

type TopProductMetric struct {
	ProductID int64
	Name      string
	TotalSold int
	Revenue   float64
}

type LowStockProduct struct {
	ID       int64
	Name     string
	Quantity int
	Category string
}

type UserAnalyticsSummary struct {
	TotalOrders   int
	TotalSpent    float64
	PendingOrders int
}

type UserRecentOrder struct {
	OrderID      int
	ProductNames string
	Status       string
	CreatedAt    time.Time
	TotalPrice   float64
}

type UserTopProduct struct {
	ProductID     int64
	Name          string
	PurchaseCount int
}

type UserRecommendation struct {
	ProductID int64
	Name      string
	Price     float64
	Category  string
	Reason    string
}
