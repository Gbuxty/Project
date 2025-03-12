package server

import (
	"github.com/gofiber/fiber/v2"
	"Project/APIGateWay/internal/handlers"
	"Project/APIGateWay/internal/service"
	"github.com/arsmn/fiber-swagger/v2"
)

type Server struct {
	app *fiber.App
}

func NewServer(authService *service.AuthService) *Server {
	app := fiber.New()

	app.Static("/", "/app")

	app.Get("/swagger/*", swagger.New(swagger.Config{
		URL:"/swagger.yaml",
	}))
	
	authHandlers := handlers.NewAuthHandlers(authService)

	app.Post("/register", authHandlers.Register)
	app.Post("/login", authHandlers.Login)
	app.Post("/logout", authHandlers.Logout)
	app.Post("/refresh", authHandlers.Refresh)
	app.Get("/me", authHandlers.Me)
	app.Post("/confirm-email", authHandlers.ConfirmEmail)

	return &Server{
		app: app,
	}
}

func (s *Server) Start(address string) error {
	return s.app.Listen(address)
}