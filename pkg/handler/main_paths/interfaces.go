package main_paths

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	SignUp(c *gin.Context)
}

type MainHandler struct {
	Authorization
}

func NewMainPaths(db *sqlx.DB) *MainHandler {
	return &MainHandler{
		Authorization: NewAuthPostgres(db),
	}
}
