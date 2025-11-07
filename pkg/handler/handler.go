package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"

	mp "github.com/Demon/backend-GO/pkg/handler/main_paths"
)


type Handler struct {
	MainHandler *mp.MainHandler
	db *sqlx.DB
}

func NewHandler(MainHandler *mp.MainHandler, db *sqlx.DB) *Handler {
	return &Handler{
		MainHandler: MainHandler,
		db : db,
	}
}


func (h *Handler)InitRoutes() *gin.Engine {
	router := gin.New()

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.MainHandler.Authorization.SignUp)
		auth.POST("/sign-in", h.MainHandler.Authorization.SignIn)
	}

	user := router.Group("/user")
	{
		user.PATCH("/chang", h.MainHandler.User.UpdateUser)
		user.GET("/get-profile", h.MainHandler.User.GetProfile)
	}

	orders := router.Group("/orders")
	{
		orders.POST("/create-order", h.MainHandler.Orders.CreateOrder)
	}

	test := router.Group("/test")
	{
		test.POST("/status", responce)
		test.POST("/count", h.getUserCount)
	}

	return router
}

func responce(c *gin.Context) {
	c.JSON(200, gin.H{"ok": true})
}

func (h *Handler) getUserCount(c *gin.Context) {
	var count int


	err := h.db.Get(&count, "SELECT COUNT(*) FROM users")
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Не удалось получить количество пользователей",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ok":    true,
		"count": count,
	})
}