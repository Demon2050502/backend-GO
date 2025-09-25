package models

import "time"

type User struct {
	ID           int64      `db:"id" json:"id"`
	FirstName    string     `db:"firstname" json:"firstname"`
	LastName     string     `db:"lastname" json:"lastname"`
	Username     string     `db:"username" json:"username"`
	Email        string     `db:"email" json:"email"`
	Phone        *string    `db:"phone" json:"phone,omitempty"`
	PasswordHash string     `db:"password_hash" json:"-"`
	AvatarURL    *string    `db:"avatar_url" json:"avatar_url,omitempty"`
	Bio          *string    `db:"bio" json:"bio,omitempty"`
	Rating       float64    `db:"rating" json:"rating"`
	Balance      float64    `db:"balance" json:"balance"`
	TypeID       int        `db:"type_id" json:"type_id"`
	IsActive     bool       `db:"is_active" json:"is_active"`
	CreatedAt    time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at" json:"updated_at"`
	LastLogin    time.Time  `db:"last_login" json:"last_login"`
}
