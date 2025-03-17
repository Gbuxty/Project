package handlers

import (
    "Project/APIGateWay/internal/domain"
    "Project/APIGateWay/internal/service"
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
    var req domain.AuthRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(401).JSON(fiber.Map{"error": "Invalid request body"})
    }

    success, err := h.authService.Register(c.Context(), &req)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }

    return c.JSON(fiber.Map{"success": success})
}

func (h *AuthHandlers) Login(c *fiber.Ctx) error {
    var req domain.AuthRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
    }

    tokens, user, err := h.authService.Login(c.Context(), &req)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }

    return c.JSON(fiber.Map{
        "tokens": tokens,
        "user":   user,
    })
}

func (h *AuthHandlers) Logout(c *fiber.Ctx) error {
    var req domain.LogoutRequest
  
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
    }

    success, err := h.authService.Logout(c.Context(), req.UserID)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }

    return c.JSON(fiber.Map{"success": success})
}

func (h *AuthHandlers) Refresh(c *fiber.Ctx) error {
    var req struct {
        RefreshToken string `json:"refresh_token"`
    }
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
    }

    tokens, err := h.authService.Refresh(c.Context(), req.RefreshToken)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }

    return c.JSON(tokens)
}

func (h *AuthHandlers) Me(c *fiber.Ctx) error {
    authHeader := c.Get("Authorization")
    if authHeader == "" {
        return c.Status(401).JSON(fiber.Map{"error": "Authorization header is required"})
    }

    accessToken := strings.TrimPrefix(authHeader, "Bearer ")

    user, err := h.authService.Me(c.Context(), accessToken)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }

    return c.JSON(user)
}

func (h *AuthHandlers) ConfirmEmail(c *fiber.Ctx) error {
    var req domain.ConfirmationRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
    }

    success, err := h.authService.ConfirmEmail(c.Context(), &req)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }

    return c.JSON(fiber.Map{"success": success})
}