package handlers

import (
	"Project/APIGateWay/internal/domain"
	"Project/APIGateWay/internal/service"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type AuthHandlers struct {
	authService *service.AuthService
}

func NewAuthHandlers(authService *service.AuthService) *AuthHandlers {
	return &AuthHandlers{
		authService: authService,
	}
}

func (h *AuthHandlers) Register(c *fiber.Ctx) error {
    var req domain.RegisterRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusUnauthorized).JSON(domain.ErrorResponse{Error: "Invalid request body"})
	}

	success, err := h.authService.Register(c.Context(), &req)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(domain.ErrorResponse{Error: "Unsuccess register"})
	}

	return c.JSON(success)
}

func (h *AuthHandlers) Login(c *fiber.Ctx) error {
    var req domain.LoginRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusUnauthorized).JSON(domain.ErrorResponse{Error: "Invalid request body"})
	}

	tokens, user, err := h.authService.Login(c.Context(), &req)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(domain.ErrorResponse{Error: "Invalid Login or Password"})
	}

	return c.JSON(domain.LoginResponse{Tokens: *tokens,User: *user})
}

func (h *AuthHandlers) Logout(c *fiber.Ctx) error {
	var req domain.LogoutRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusUnauthorized).JSON(domain.ErrorResponse{Error: "Invalid request body"})
	}

	success, err := h.authService.Logout(c.Context(), req.UserID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(domain.ErrorResponse{Error: "Failed logout"})
	}

	return c.JSON(success)
}

func (h *AuthHandlers) Refresh(c *fiber.Ctx) error {
    var req domain.RefreshRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusUnauthorized).JSON(domain.ErrorResponse{Error: "Invalid request body"})
	}

	tokens, err := h.authService.Refresh(c.Context(), req.RefreshToken)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(domain.ErrorResponse{Error: "Failed Refresh tokens"})
	}

	return c.JSON(tokens)
}

func (h *AuthHandlers) Me(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(http.StatusUnauthorized).JSON(domain.ErrorResponse{Error: "Failed get header"})
	}

	accessToken := strings.TrimPrefix(authHeader, "Bearer ")

	user, err := h.authService.Me(c.Context(), accessToken)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(domain.ErrorResponse{Error: "Failed Get Info of ME"})
	}

	return c.JSON(user)
}

func (h *AuthHandlers) ConfirmEmail(c *fiber.Ctx) error {
	var req domain.ConfirmationRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusUnauthorized).JSON(domain.ErrorResponse{Error: "Invalid request body"})
	}

	success, err := h.authService.ConfirmEmail(c.Context(), &req)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(domain.ErrorResponse{Error: "Failed Send Confim Token"})
	}

	return c.JSON(success)
}

func CheckToken(c *fiber.Ctx, authService *service.AuthService) (*domain.User, error) {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return nil, fiber.NewError(http.StatusUnauthorized, "Authorization header is missing")
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == authHeader {
		return nil, fiber.NewError(http.StatusUnauthorized, "Invalid authorization header format")
	}

	user, err := authService.Me(c.Context(), token)
	if err != nil {
		return nil, fiber.NewError(http.StatusUnauthorized, "Invalid token")
	}

	return user, nil
}