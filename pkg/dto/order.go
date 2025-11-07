package dto

type CreateOrderRequest struct {
	Token       string  `json:"token" db:"token"`
	Title       string  `json:"title" binding:"required" db:"title"`
	Description string  `json:"description" binding:"required" db:"description"`
	CategoryID  int     `json:"category_id" db:"category_id"`
	Price       float64 `json:"price" binding:"required" db:"price"`
	PriceMax    float64 `json:"price_max" binding:"required" db:"price_max"`
	TimeStart   string  `json:"time_start" binding:"required" db:"time_start"`
	TimeEnd     string  `json:"time_end" binding:"required" db:"time_end"`
	IsVisible   bool    `json:"is_visible" db:"is_visible"`
}

type CreateOrderResponse struct {
	ID int64 `json:"id"`
}
