package dto

type UserCreate struct {
	Username string
	Password string
	Role     string
}

type ProductCreate struct {
	Name     string
	Price    float32
	Quantity int64
	Category string
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
