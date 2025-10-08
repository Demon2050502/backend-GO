package main_paths

import (
	"github.com/jmoiron/sqlx"
)


type ConDataBase struct {
	db *sqlx.DB
}

func NewChangUser(db *sqlx.DB) *ConDataBase {
	return &ConDataBase{db: db}
}

// заменя данных пользователя
