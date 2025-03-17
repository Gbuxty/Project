package handlers

import (
	"Project/APIGateWay/internal/domain"
	"Project/APIGateWay/internal/service"
	"strings"

	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc/metadata"
)

type FeedHandlers struct {
	feedService *service.FeedService
}

func NewFeedHandlers(feedService *service.FeedService) *FeedHandlers {
	return &FeedHandlers{
		feedService: feedService,
	}
}

func (h *FeedHandlers) CreatePost(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(401).JSON(fiber.Map{"error": "Authorization header is missing"})
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == authHeader {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid authorization header format"})
	}

	var req domain.CreatePostRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	md := metadata.Pairs("authorization", "Bearer "+token)
	ctx := metadata.NewOutgoingContext(c.Context(), md)

	post, err := h.feedService.CreatePost(ctx, &req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(post)
}

func (h *FeedHandlers) GetAllPosts(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("pageSize", 10)

	posts, totalPosts, err := h.feedService.GetAllPosts(c.Context(), page, pageSize)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"posts":       posts,
		"total_posts": totalPosts,
	})
}
