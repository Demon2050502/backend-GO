package repository

import (
	"context"

	"github.com/Demon/backend-GO/pkg/models"
	"github.com/jmoiron/sqlx"
)

type AuthRepository interface {
	CreateUser(ctx context.Context, u *models.User) (int64, error)
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByIdentifier(ctx context.Context, identifier string) (*models.User, error)
	UpdateLastLogin(ctx context.Context, userID int64) error
}

type Repositories struct {
	AuthRepository

}

func NewRepositories(db *sqlx.DB) *Repositories {
	return &Repositories{
		AuthRepository: NewAuthRepository(db),
	}
}
