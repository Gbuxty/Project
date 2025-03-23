package server

import (
	"Project/APIGateWay/internal/handlers"
	"Project/APIGateWay/internal/service"

	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
)

type Server struct {
	app *fiber.App
}

func NewServer(authService *service.AuthService, feedService *service.FeedService) *Server {
	app := fiber.New()

	app.Static("/", "/app")

	app.Get("/swagger/*", swagger.New(swagger.Config{
		URL: "/swagger.yaml",
	}))

	authHandlers := handlers.NewAuthHandlers(authService)
	feedHandlers := handlers.NewFeedHandlers(feedService)

	app.Post("/register", authHandlers.Register)
	app.Post("/login", authHandlers.Login)
	app.Post("/logout", authHandlers.Logout)
	app.Post("/refresh", authHandlers.Refresh)
	app.Get("/me", authHandlers.Me)
	app.Post("/confirm-email", authHandlers.ConfirmEmail)

	app.Post("/posts", feedHandlers.CreatePost)
	app.Get("/posts/all", feedHandlers.GetAllPosts)

	return &Server{
		app: app,
	}
}

func (s *Server) Start(address string) error {
	return s.app.Listen(address)
}
