package main_paths

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	dto "github.com/Demon/backend-GO/pkg/dto"
	auxpath "github.com/Demon/backend-GO/pkg/handler/auxiliary_paths"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)


type ConUserBD struct {
	db *sqlx.DB
}

func NewUser(db *sqlx.DB) *ConUserBD {
	return &ConUserBD{db: db}
}

func (h *ConUserBD) GetProfile(c *gin.Context) {
	ctx := c.Request.Context()

	token := c.Query("token")
	if token == "" {
		var body struct {
			Token string `json:"token"`
		}
		if err := c.ShouldBindJSON(&body); err == nil {
			token = body.Token
		}
	}

	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "missing token",
			"message": "Токен не предоставлен",
		})
		return
	}

	userID, err := auxpath.GetUserIDFromToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "invalid token",
			"details": err.Error(),
			"message": "Неверный или просроченный токен",
		})
		return
	}

	var profile dto.GetProfileResponse

	query := `
		SELECT 
			firstname,
			lastname,
			username,
			email,
			phone,
			avatar_url,
			bio,
			type_id,
			is_active,
			to_char(created_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"') as created_at,
			to_char(updated_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"') as updated_at
		FROM users
		WHERE id = $1
	`

	if err := h.db.GetContext(ctx, &profile, query, userID); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "user not found",
				"message": "Пользователь не найден",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "db error",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, profile)
}

func (h *ConUserBD) UpdateUser(c *gin.Context) {
	ctx := c.Request.Context()

	// 1) Парсим тело запроса
	var req dto.ChangUpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request",
			"details": err.Error(),
			"message": "Неправильные входные данные",
		})
		return
	}

	if req.Token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "missing token",
			"message": "Токен не предоставлен",
		})
		return
	}

	// 2) Получаем userID из токена
	userID, err := auxpath.GetUserIDFromToken(req.Token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "invalid token",
			"details": err.Error(),
			"message": "Неверный или просроченный токен",
		})
		return
	}

	// Проверяем username в отдельной функции
	usernameWillChange, _, err := h.CheckUserUniqueFields(ctx, userID, req.Username, req.Email)

	if err != nil {
		if strings.Contains(err.Error(), "username already taken") {
			c.JSON(http.StatusConflict, gin.H{"error": "username already taken", "message": "Имя пользователя уже занято"})
			return
		}
		if strings.Contains(err.Error(), "email already taken") {
			c.JSON(http.StatusConflict, gin.H{"error": "email already taken", "message": "Email уже используется другим пользователем"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Собираем UPDATE динамически
	setParts := []string{}
	args := []interface{}{}
	argPos := 1
	add := func(field string, value interface{}) {
		setParts = append(setParts, fmt.Sprintf("%s = $%d", field, argPos))
		args = append(args, value)
		argPos++
	}

	if req.FirstName != nil {
		add("firstname", *req.FirstName)
	}
	if req.LastName != nil {
		add("lastname", *req.LastName)
	}
	if req.Username != nil {
		add("username", *req.Username)
	}
	if req.Email != nil {
		add("email", *req.Email)
	}
	if req.Phone != nil {
		add("phone", *req.Phone)
	}
	if req.AvatarURL != nil {
		add("avatar_url", *req.AvatarURL)
	}
	if req.Bio != nil {
		add("bio", *req.Bio)
	}
	if req.TypeID != nil {
		add("type_id", *req.TypeID)
	}

	if len(setParts) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no fields to update", "message": "Нет полей для обновления"})
		return
	}

	setParts = append(setParts, "updated_at = NOW()")
	query := fmt.Sprintf(`
		UPDATE users
		SET %s
		WHERE id = $%d
		RETURNING firstname, lastname, username, email, phone, avatar_url, bio, type_id
	`, strings.Join(setParts, ", "), argPos)
	args = append(args, userID)

	var user dto.ChangUserResponse
	if err := h.db.GetContext(ctx, &user, query, args...); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		if strings.Contains(err.Error(), "duplicate key") {
			c.JSON(http.StatusConflict, gin.H{
				"error": "unique constraint",
				"details": err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error", "details": err.Error()})
		return
	}

	// Генерация нового токена при изменении username
	if usernameWillChange {
		newToken, err := auxpath.GenerateToken(userID, user.Username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate new token", "details": err.Error()})
			return
		}
		user.Token = newToken
	} else {
		user.Token = req.Token
	}

	c.JSON(http.StatusOK, user)
}

func (h *ConUserBD) CreatePortfolio(c *gin.Context) {
	ctx := c.Request.Context()

	var req dto.CreatePortfolioRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
			"details": err.Error(),
		})
		return
	}

	userID, err := auxpath.GetUserIDFromToken(req.Token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid token",
			"details": err.Error(),
		})
		return
	}

	tx, err := h.db.BeginTxx(ctx, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "db error",
			"details": err.Error(),
		})
		return
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	// Создаём портфолио
	var portfolioID int64

	isPublic := true
	if req.IsPublic != nil {
		isPublic = *req.IsPublic
	}

	queryPortfolio := `
		INSERT INTO portfolios (user_id, title, description, category_id, is_public)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	err = tx.GetContext(ctx, &portfolioID, queryPortfolio,
		userID,
		req.Title,
		req.Description,
		req.CategoryID,
		isPublic,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "db error",
			"details": err.Error(),
		})
		return
	}

	// Добавляем файлы
	if len(req.Files) > 0 {
		queryFile := `
			INSERT INTO portfolio_files (portfolio_id, file_url, file_type, title, description)
			VALUES ($1, $2, $3, $4, $5)
		`

		for _, f := range req.Files {
			_, err = tx.ExecContext(ctx, queryFile,
				portfolioID,
				f.FileURL,
				f.FileType,
				f.Title,
				f.Description,
			)

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "db error",
					"details": err.Error(),
				})
				return
			}
		}
	}

	// Коммитим изменения
	if err = tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "db error",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.CreatePortfolioResponse{
		ID: portfolioID,
	})
}

func (h *ConUserBD) GetPortfolios(c *gin.Context) {
	ctx := c.Request.Context()

	var body struct{ Token string `json:"token"` }
	token := c.GetHeader("Authorization")
	if token == "" {
		_ = c.ShouldBindJSON(&body)
		token = body.Token
	}
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
		return
	}

	userID, err := auxpath.GetUserIDFromToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	query := `
        SELECT id, title, description, category_id, is_public
        FROM portfolios
        WHERE user_id = $1
        ORDER BY created_at DESC
    `
	var portfolios []dto.PortfolioItem
	err = h.db.SelectContext(ctx, &portfolios, query, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "db error",
			"details": err.Error(),
		})
		return
	}

	
		fmt.Println(portfolios)


	for i := range portfolios {
		qFiles := `
            SELECT id, file_url, file_type, title, description
            FROM portfolio_files
            WHERE portfolio_id = $1
            ORDER BY uploaded_at DESC
        `
		var files []dto.PortfolioFile
		err := h.db.SelectContext(ctx, &files, qFiles, portfolios[i].ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "file fetch error",
				"details": err.Error(),
			})
			return
		}
		portfolios[i].Files = files

		
	}

	c.JSON(http.StatusOK, dto.GetPortfoliosResponse{
		Portfolios: portfolios,
	})
}

func (h *ConUserBD) UpdatePortfolio(c *gin.Context) {
    ctx := c.Request.Context()

    // Parse JSON
    var req dto.UpdatePortfolioRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error":   "invalid request",
            "details": err.Error(),
        })
        return
    }

    userID, err := auxpath.GetUserIDFromToken(req.Token)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{
            "error":   "invalid token",
            "details": err.Error(),
        })
        return
    }

    // Check portfolio exists and belongs to the user
    err = h.db.GetContext(ctx, &userID,
        `SELECT user_id FROM portfolios WHERE id = $1`,
        req.PortfolioID,
    )
    if err == sql.ErrNoRows {
        c.JSON(http.StatusNotFound, gin.H{"error": "portfolio not found"})
        return
    }
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "db error", "details": err.Error()})
        return
    }

    // Start transaction
    tx, err := h.db.BeginTxx(ctx, nil)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "transaction error", "details": err.Error()})
        return
    }

    // Build UPDATE query
    setParts := []string{}
    args := []interface{}{}
    arg := 1

    add := func(field string, value interface{}) {
        setParts = append(setParts, fmt.Sprintf("%s = $%d", field, arg))
        args = append(args, value)
        arg++
    }

    if req.Title != nil {
        add("title", *req.Title)
    }
    if req.Description != nil {
        add("description", *req.Description)
    }
    if req.CategoryID != nil {
        add("category_id", *req.CategoryID)
    }
    if req.IsPublic != nil {
        add("is_public", *req.IsPublic)
    }

    if len(setParts) > 0 {
        setParts = append(setParts, "updated_at = NOW()")
        args = append(args, req.PortfolioID)

        query := fmt.Sprintf(`
            UPDATE portfolios
            SET %s
            WHERE id = $%d
        `, strings.Join(setParts, ", "), arg)

        if _, err := tx.ExecContext(ctx, query, args...); err != nil {
            tx.Rollback()
            c.JSON(http.StatusInternalServerError, gin.H{"error": "update failed", "details": err.Error()})
            return
        }
    }

    // ---- Update Files (if included) ----
    if req.Files != nil {
        _, err := tx.ExecContext(ctx,
            `DELETE FROM portfolio_files WHERE portfolio_id = $1`,
            req.PortfolioID,
        )
        if err != nil {
            tx.Rollback()
            c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to clear files", "details": err.Error()})
            return
        }

        for _, f := range *req.Files {
            _, err = tx.ExecContext(ctx,
                `INSERT INTO portfolio_files (portfolio_id, file_url, file_type, title, description)
                 VALUES ($1, $2, $3, $4, $5)`,
                req.PortfolioID, f.FileURL, f.FileType, f.Title, f.Description,
            )
            if err != nil {
                tx.Rollback()
                c.JSON(http.StatusInternalServerError, gin.H{"error": "file insert error", "details": err.Error()})
                return
            }
        }
    }

    // Commit transaction
    if err := tx.Commit(); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "commit error", "details": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"status": "portfolio updated"})
}


func (h *ConUserBD) CheckUserUniqueFields(
	ctx context.Context,
	userID int64,
	newUsername *string,
	newEmail *string,
) (usernameChanged bool, emailChanged bool, err error) {

	// --- Проверка username ---
	if newUsername != nil {
		var currentUsername string
		if err := h.db.GetContext(ctx, &currentUsername, `SELECT username FROM users WHERE id = $1`, userID); err != nil {
			if err == sql.ErrNoRows {
				return false, false, fmt.Errorf("user not found")
			}
			return false, false, fmt.Errorf("db error: %w", err)
		}

		// Если новый username отличается — проверяем, что он свободен
		if *newUsername != currentUsername {
			var count int
			if err := h.db.GetContext(ctx, &count, `
				SELECT COUNT(1) FROM users WHERE username = $1 AND id != $2
			`, *newUsername, userID); err != nil {
				return false, false, fmt.Errorf("db error: %w", err)
			}
			if count > 0 {
				return false, false, fmt.Errorf("username already taken")
			}
			usernameChanged = true
		}
	}

	// --- Проверка email ---
	if newEmail != nil {
		var currentEmail sql.NullString
		if err := h.db.GetContext(ctx, &currentEmail, `SELECT email FROM users WHERE id = $1`, userID); err != nil {
			if err == sql.ErrNoRows {
				return false, false, fmt.Errorf("user not found")
			}
			return false, false, fmt.Errorf("db error: %w", err)
		}

		// Если email изменился — проверяем, что он свободен
		if !currentEmail.Valid || *newEmail != currentEmail.String {
			var count int
			if err := h.db.GetContext(ctx, &count, `
				SELECT COUNT(1) FROM users WHERE email = $1 AND id != $2
			`, *newEmail, userID); err != nil {
				return false, false, fmt.Errorf("db error: %w", err)
			}
			if count > 0 {
				return false, false, fmt.Errorf("email already taken")
			}
			emailChanged = true
		}
	}

	return usernameChanged, emailChanged, nil
}
