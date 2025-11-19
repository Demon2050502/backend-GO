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

type CreateBidRequest struct {
	Token       string  `json:"token" binding:"required"`
	OrderID     int64   `json:"order_id" binding:"required" db:"order_id"`
	PortfolioID int64   `json:"portfolio_id" binding:"required" db:"portfolio_id"`
	Proposal    string  `json:"proposal"`
	Price       float64 `json:"price" binding:"required"`
}

type BidResponse struct {
	ID          int64   `json:"id"`
	OrderID     int64   `json:"order_id" db:"order_id"`
	ExecutorID  int64   `json:"executor_id" db:"executor_id"`
	PortfolioID int64   `json:"portfolio_id" db:"portfolio_id"`
	Proposal    string  `json:"proposal"`
	Price       float64 `json:"price"`
	Status      string  `json:"status"`
	CreatedAt   string  `json:"created_at" db:"created_at"`
}
