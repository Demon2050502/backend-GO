package main_paths

import (
	"net/http"

	dto "github.com/Demon/backend-GO/pkg/dto"
	auxpath "github.com/Demon/backend-GO/pkg/handler/auxiliary_paths"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)


type ConOrderBD struct {
	db *sqlx.DB
}

func NewOrder(db *sqlx.DB) *ConOrderBD {
	return &ConOrderBD{db: db}
}

func (h *ConOrderBD) CreateOrder(c *gin.Context) {
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

func (h *ConOrderBD) CreateBid(c *gin.Context) {
    ctx := c.Request.Context()

    var req dto.CreateBidRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error":   "invalid request",
            "details": err.Error(),
        })
        return
    }

    if req.Token == "" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
        return
    }

    executorID, err := auxpath.GetUserIDFromToken(req.Token)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{
            "error":   "invalid token",
            "details": err.Error(),
        })
        return
    }

    // Проверяем, что заказ существует
    var exists bool
    err = h.db.GetContext(ctx, &exists,
        `SELECT EXISTS(SELECT 1 FROM orders WHERE id = $1)`, req.OrderID)
    if err != nil || !exists {
        c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
        return
    }

    // Проверяем, что портфолио принадлежит пользователю
    err = h.db.GetContext(ctx, &exists,
        `SELECT EXISTS(SELECT 1 FROM portfolios WHERE id = $1 AND user_id = $2)`,
        req.PortfolioID, executorID)
    if err != nil || !exists {
        c.JSON(http.StatusForbidden, gin.H{"error": "portfolio does not belong to user"})
        return
    }

    // Создаём отклик
    query := `
        INSERT INTO order_bids (order_id, executor_id, portfolio_id, proposal, price)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id, order_id, executor_id, portfolio_id, proposal, price, status, created_at
    `

    var bid dto.BidResponse
    err = h.db.GetContext(ctx, &bid, query,
        req.OrderID, executorID, req.PortfolioID, req.Proposal, req.Price)

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error":   "db error",
            "details": err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, bid)
}
