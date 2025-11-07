package main_paths

import (
	"net/http"

	dto "github.com/Demon/backend-GO/pkg/dto"
	auxpath "github.com/Demon/backend-GO/pkg/handler/auxiliary_paths"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)


type OrderPostgres struct {
	db *sqlx.DB
}

func NewOrderPostgres(db *sqlx.DB) *OrderPostgres {
	return &OrderPostgres{db: db}
}

func (h *OrderPostgres) CreateOrder(c *gin.Context) {
	ctx := c.Request.Context()

	var req dto.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
			"details": err.Error(),
		})
		return
	}

	if req.Token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "missing token",
		})
		return
	}

	// Получаем userID
	userID, err := auxpath.GetUserIDFromToken(req.Token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "invalid token",
			"details": err.Error(),
			"message": "Неверный или просроченный токен",
		})
		return
	}

	// Вставляем заказ
	query := `
        INSERT INTO orders (title, description, category_id, price, price_max, time_start, time_end, is_visible, client_id)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
        RETURNING id
    `

	var orderID int64
	err = h.db.GetContext(ctx, &orderID, query,
		req.Title, req.Description, req.CategoryID,
		req.Price, req.PriceMax, req.TimeStart, req.TimeEnd,
		req.IsVisible, userID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "db error",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.CreateOrderResponse{ID: orderID})
}
