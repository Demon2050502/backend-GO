package handler

import (
	"github.com/Demon/backend-GO/pkg/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	auth *service.Services
}

func NewHandler(auth *service.Services) *Handler {
	return &Handler{auth: auth}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.SignUp)
		auth.POST("/sign-in", h.SignIn)
	}

	return router
}


// func responce(c *gin.Context) {
//     c.JSON(200, gin.H{"ok": true})
// }