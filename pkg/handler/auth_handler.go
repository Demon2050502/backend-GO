package handler

import (
	"errors"
	"net/http"

	"github.com/Demon/backend-GO/pkg/dto"
	"github.com/Demon/backend-GO/pkg/repository"
	"github.com/Demon/backend-GO/pkg/service"
	"github.com/gin-gonic/gin"
)

// SignUp обработчик регистрации
func (h *Handler) SignUp(c *gin.Context) {
	var req dto.SignUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.auth.SignUp(c.Request.Context(), &req)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrDuplicateUsername):
			c.JSON(http.StatusConflict, gin.H{"error": "username already exists"})
		case errors.Is(err, repository.ErrDuplicateEmail):
			c.JSON(http.StatusConflict, gin.H{"error": "email already exists"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// SignIn обработчик авторизации
func (h *Handler) SignIn(c *gin.Context) {
	var req dto.SignInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.auth.SignIn(c.Request.Context(), &req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}
