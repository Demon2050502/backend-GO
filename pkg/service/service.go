package service

import (
	"context"

	"github.com/Demon/backend-GO/pkg/dto"
	"github.com/Demon/backend-GO/pkg/repository"
)

type AuthService interface {
	SignUp(ctx context.Context, req *dto.SignUpRequest) (*dto.SignUpResponse, error)
	SignIn(ctx context.Context, req *dto.SignInRequest) (*dto.SignInResponse, error)
}

	
type Services struct {
	AuthService
}

func NewService(repos *repository.Repositories) *Services {
	return &Services{
		AuthService: NewAuthService(repos.AuthRepository),
	}
}
