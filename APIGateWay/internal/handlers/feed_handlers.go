package handlers

import (
	"Project/APIGateWay/internal/domain"
	"Project/APIGateWay/internal/service"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc/metadata"
)

type FeedHandlers struct {
	feedService *service.FeedService
	auth        service.AuthService //ЭТО ПОЛНАЯ ХУЙНЯ БЕЗ ЭТОГО НЕ МОГУ ИСПОЛЬЗОВАТЬ Me функцию из authservice которая принимает токен что бы валидировать юзера
}

func NewFeedHandlers(feedService *service.FeedService) *FeedHandlers {
	return &FeedHandlers{
		feedService: feedService,
	}
}

func (h *FeedHandlers) CreatePost(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(http.StatusUnauthorized).JSON(domain.ErrorResponse{Error: "Authorization header is missing"})
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == authHeader {
		return c.Status(http.StatusUnauthorized).JSON(domain.ErrorResponse{Error: "Invalid authorization header format"})
	}

	user, err := h.auth.Me(c.Context(), token)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(domain.ErrorResponse{Error: "Invalid token"})
	}

	var req domain.CreatePostRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(domain.ErrorResponse{Error: "Invalid request body"})
	}

	

	md := metadata.Pairs("user_id", user.ID)
	ctx := metadata.NewOutgoingContext(c.Context(), md)

	post, err := h.feedService.CreatePost(ctx, &req)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(domain.ErrorResponse{Error: "Failed Create Post"})
	}

	return c.JSON(post)
}

func (h *FeedHandlers) GetAllPosts(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("pageSize", 10)

	posts, totalPosts, err := h.feedService.GetAllPosts(c.Context(), page, pageSize)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(domain.ErrorResponse{Error: "Failed Get all posts"})
	}

	return c.JSON(domain.AllPostResponse{Posts: posts, TotalPosts: totalPosts})
}
