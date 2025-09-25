package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Demon/backend-GO/pkg/models"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrDuplicateUsername = errors.New("username already exists")
	ErrDuplicateEmail    = errors.New("email already exists")
)



type authRepo struct {
	db *sqlx.DB
}

func NewAuthRepository(db *sqlx.DB) AuthRepository {
	return &authRepo{db: db}
}

func (r *authRepo) CreateUser(ctx context.Context, u *models.User) (int64, error) {
	q := `
INSERT INTO users
(firstname, lastname, username, email, phone, password_hash, avatar_url, bio, rating, balance, type_id, is_active, created_at, updated_at, last_login)
VALUES
($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15)
RETURNING id
`
	var id int64
	err := r.db.QueryRowxContext(ctx, q,
		u.FirstName, u.LastName, u.Username, u.Email, u.Phone, u.PasswordHash,
		u.AvatarURL, u.Bio, u.Rating, u.Balance, u.TypeID, u.IsActive,
		u.CreatedAt, u.UpdatedAt, u.LastLogin,
	).Scan(&id)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			switch pqErr.Constraint {
			case "users_username_key", "username_key":
				return 0, ErrDuplicateUsername
			case "users_email_key", "email_key":
				return 0, ErrDuplicateEmail
			default:
				// fallback detection
				if containsIgnoreCase(pqErr.Message, "username") {
					return 0, ErrDuplicateUsername
				}
				if containsIgnoreCase(pqErr.Message, "email") {
					return 0, ErrDuplicateEmail
				}
			}
		}
		return 0, fmt.Errorf("create user: %w", err)
	}
	return id, nil
}

func (r *authRepo) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	var u models.User
	q := `SELECT * FROM users WHERE username = $1 LIMIT 1`
	if err := r.db.GetContext(ctx, &u, q, username); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get by username: %w", err)
	}
	return &u, nil
}

func (r *authRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var u models.User
	q := `SELECT * FROM users WHERE email = $1 LIMIT 1`
	if err := r.db.GetContext(ctx, &u, q, email); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get by email: %w", err)
	}
	return &u, nil
}

func (r *authRepo) GetByIdentifier(ctx context.Context, identifier string) (*models.User, error) {
	u, err := r.GetByUsername(ctx, identifier)
	if err != nil {
		return nil, err
	}
	if u != nil {
		return u, nil
	}
	return r.GetByEmail(ctx, identifier)
}

func (r *authRepo) UpdateLastLogin(ctx context.Context, userID int64) error {
	q := `UPDATE users SET last_login = $1, updated_at = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, q, time.Now().UTC(), userID)
	if err != nil {
		return fmt.Errorf("update last_login: %w", err)
	}
	return nil
}

// helper
func containsIgnoreCase(s, sub string) bool {
	if s == "" || sub == "" {
		return false
	}
	return stringContainsFold(s, sub)
}

func stringContainsFold(s, substr string) bool {
	// use strings
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}


