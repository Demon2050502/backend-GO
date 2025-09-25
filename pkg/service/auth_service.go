package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Demon/backend-GO/pkg/dto"
	"github.com/Demon/backend-GO/pkg/models"
	"github.com/Demon/backend-GO/pkg/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmptyRequest       = errors.New("empty request")
)

type authService struct {
	repos repository.AuthRepository
}

func NewAuthService(repos repository.AuthRepository) *authService {
	return &authService{repos: repos}
}


func (s *authService) SignUp(ctx context.Context, req *dto.SignUpRequest) (*dto.SignUpResponse, error) {
	if req == nil {
		return nil, ErrEmptyRequest
	}

	// uniqueness checks
	if u, err := s.repos.GetByUsername(ctx, req.Username); err != nil {
		return nil, fmt.Errorf("check username: %w", err)
	} else if u != nil {
		return nil, repository.ErrDuplicateUsername
	}

	if u, err := s.repos.GetByEmail(ctx, req.Email); err != nil {
		return nil, fmt.Errorf("check email: %w", err)
	} else if u != nil {
		return nil, repository.ErrDuplicateEmail
	}

	// hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	now := time.Now().UTC()
	typeID := 1
	if req.TypeID != nil && *req.TypeID > 0 {
		typeID = *req.TypeID
	}

	user := &models.User{
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Username:     req.Username,
		Email:        req.Email,
		Phone:        req.Phone,
		PasswordHash: string(hashed),
		AvatarURL:    req.AvatarURL,
		Bio:          req.Bio,
		Rating:       0.0,
		Balance:      0.0,
		TypeID:       typeID,
		IsActive:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
		LastLogin:    now,
	}

	id, err := s.repos.CreateUser(ctx, user)
	if err != nil {
		// propagate friendly repo errors
		if errors.Is(err, repository.ErrDuplicateUsername) {
			return nil, repository.ErrDuplicateUsername
		}
		if errors.Is(err, repository.ErrDuplicateEmail) {
			return nil, repository.ErrDuplicateEmail
		}
		return nil, fmt.Errorf("create user: %w", err)
	}

	resp := &dto.SignUpResponse{
		ID:        id,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: now.Format(time.RFC3339),
	}
	return resp, nil
}

func (s *authService) SignIn(ctx context.Context, req *dto.SignInRequest) (*dto.SignInResponse, error) {
	if req == nil {
		return nil, ErrEmptyRequest
	}

	user, err := s.repos.GetByIdentifier(ctx, req.Identifier)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	if user == nil {
		return nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	// best-effort last_login update
	_ = s.repos.UpdateLastLogin(ctx, user.ID)

	token, err := GenerateToken(user.ID, user.Username, user.TypeID)
	if err != nil {
		return nil, fmt.Errorf("generate token: %w", err)
	}

	var out dto.SignInResponse
	out.Token = token
	out.User.ID = user.ID
	out.User.Username = user.Username
	out.User.Email = user.Email
	out.User.Balance = user.Balance
	out.User.TypeID = user.TypeID

	return &out, nil
}

// временно
func GenerateToken(userID int64, username string, typeID int) (string, error) {
	return "dacqwcqczdfqvqwvq", nil
}