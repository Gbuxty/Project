package jwt

import (
	"context"
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
)

func ExtractTokenFromContext(ctx context.Context) (string, error) {
    md, ok := metadata.FromIncomingContext(ctx)
    if !ok {
        return "", fmt.Errorf("metadata is not provided")
    }

    authHeader := md.Get("authorization")
    if len(authHeader) == 0 {
        return "", fmt.Errorf("authorization header is not provided")
    }

   
    token := strings.TrimPrefix(authHeader[0], "Bearer ")
    if token == authHeader[0] {
        return "", fmt.Errorf("invalid authorization header format")
    }

    return token, nil
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