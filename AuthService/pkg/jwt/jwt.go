package jwt

import (
	"fmt"

	"time"

	"github.com/dchest/uniuri"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

func GenerateToken(userID uuid.UUID, email string, secretKey string, TokenTTL time.Duration) (string, time.Time, error) {

	expiresAt := time.Now().Add(TokenTTL)

	claims := jwt.MapClaims{
		"user_id": userID.String(),
		"email":   email,
		"exp":     expiresAt.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {

		return "", time.Time{}, fmt.Errorf("failed to generate token: %w", err)
	}

	return signedToken, expiresAt, nil
}

func GenerateConfirmationCode() string {
	return uniuri.NewLen(6)
}

func ExtractUserIDFromToken(tokenString string, secretKey string) (uuid.UUID, error) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {

			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if err != nil {

		return uuid.Nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {

		return uuid.Nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {

		return uuid.Nil, fmt.Errorf("failed to extract claims from token")
	}

	userIDStr, ok := claims["user_id"].(string)
	if !ok {

		return uuid.Nil, fmt.Errorf("user_id not found in token claims")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to parse user_id as UUID: %w", err)
	}

	return userID, nil
}
