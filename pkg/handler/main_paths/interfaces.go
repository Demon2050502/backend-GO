package main_paths

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	SignUp(c *gin.Context)
	SignIn(c *gin.Context)
}

type User interface {
	UpdateUser(*gin.Context)
	GetProfile(*gin.Context)
}

type MainHandler struct {
	Authorization
	User
}

func NewMainPaths(db *sqlx.DB) *MainHandler {
	return &MainHandler{
		Authorization: NewAuthPostgres(db),
		User: NewUser(db),
	}
}
