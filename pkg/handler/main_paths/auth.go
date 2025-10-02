package main_paths

import (
	"errors"
	"net/http"
	"regexp"
	"time"
	"unicode"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"

	dto "github.com/Demon/backend-GO/pkg/dto"
	auxpath "github.com/Demon/backend-GO/pkg/handler/auxiliary_paths"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

// SignUp — обработчик регистрации пользователя.
func (h *AuthPostgres) SignUp(c *gin.Context) {
	var req dto.SignUpRequest

	// 1) Валидация входных данных на основе тегов binding в dto.SignUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request", "details": err.Error()})
		return
	}

	// 2) Проверяем пароль и хешируем пароль
	if err := validatePassword(req.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}
	passwordHash := string(hashed)

	// 3) Подготовим значения для вставки
	typeID := 0
	if req.TypeID != nil {
		typeID = *req.TypeID
	}

	// 4) Выполняем INSERT ... RETURNING id, created_at
	insertQuery := `
		INSERT INTO users
			(firstname, lastname, username, email, phone, password_hash, avatar_url, bio, type_id)
		VALUES
			($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at
	`

	var newID int64
	var createdAt time.Time

	// QueryRowx возвращает одну строку — используем Scan для получения возвращаемых полей.
	row := h.db.QueryRowx(
		insertQuery,
		req.Name,
		req.Name,
		req.Username,
		req.Email,
		req.Phone,       // *string или nil
		passwordHash,
		req.AvatarURL,   // *string или nil
		req.Bio,         // *string или nil
		typeID,
	)

	if err := row.Scan(&newID, &createdAt); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user", "details": err.Error()})
		return
	}

	// Генерация JWT-токена отдельной функцией
	tokenString, err := auxpath.GenerateToken(newID, req.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	// 5) Формируем ответ по dto.SignUpResponse
	resp := dto.SignUpResponse{
		Token: tokenString,
	}
	resp.User.ID = newID
	resp.User.Username = req.Username
	resp.User.Email = req.Email

	c.JSON(http.StatusCreated, resp)
}

func validatePassword(pwd string) error {
	// Проверка на пробелы
	for _, r := range pwd {
		if unicode.IsSpace(r) {
			return errors.New("password must not contain spaces")
		}
	}

	// Проверка на кириллицу
	// Диапазон \u0400–\u04FF охватывает кириллицу
	cyrillicRegex := regexp.MustCompile(`[\p{Cyrillic}]`)
	if cyrillicRegex.MatchString(pwd) {
		return errors.New("password must not contain Cyrillic characters")
	}

	return nil
}
