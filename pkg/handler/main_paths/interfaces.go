package main_paths

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	SignUp(c *gin.Context)
	SignIn(c *gin.Context)
}

type ChangUser interface {

}

type MainHandler struct {
	Authorization
	ChangUser
}

func NewMainPaths(db *sqlx.DB) *MainHandler {
	return &MainHandler{
		Authorization: NewAuthPostgres(db),
		ChangUser: NewChangUser(db),
	}
}
