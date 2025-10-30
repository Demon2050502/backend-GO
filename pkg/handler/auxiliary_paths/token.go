package auxiliary_paths

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func getSecretKey() ([]byte, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, errors.New("JWT_SECRET environment variable is not set")
	}
	return []byte(jwtSecret), nil
}

func GenerateToken(userID int64, username string) (string, error) {
	secret, err  := getSecretKey()
	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	println("Generating token for userID:", userID)

	return token.SignedString(secret)
}

// GetUserIDFromToken извлекает user_id из токена.
func GetUserIDFromToken(tokenString string) (int64, error) {
	if tokenString == "" {
		return 0, errors.New("empty token")
	}

	secret, err  := getSecretKey()
	if err != nil {
		return 0, err
	}

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		// Проверяем, что используется нужный метод подписи
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return secret, nil
	})
	if err != nil || !token.Valid {
		return 0, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid claims")
	}

	// Извлекаем user_id из claims
	if v, ok := claims["user_id"]; ok {
		switch id := v.(type) {
		case float64:
			return int64(id), nil
		case int64:
			return id, nil
		case string:
			var userID int64
			_, err := fmt.Sscan(id, &userID)
			if err == nil {
				return userID, nil
			}
		}
	}

	return 0, errors.New("user_id not found in token")
}