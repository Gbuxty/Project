package service

import (
	"Project/FeedService/internal/domain/models"
	"Project/FeedService/internal/repositories/postgres"

	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
)

type FeedService struct {
	repo      FeedRepositories
	secretKey string
}

type FeedRepositories interface {
	CreatePost(ctx context.Context, post *models.Post) error
	GetALLPosts(ctx context.Context, page, pageSize int) ([]models.Post, int, error)
}

func NewFeedService(feedRepo *postgres.FeedRepositories, secretKey string) *FeedService {
	return &FeedService{
		repo:      feedRepo,
		secretKey: secretKey,
	}
}

func (s *FeedService) CreatePost(ctx context.Context, content, imageURL string) (*models.Post, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("metadata is missing")
	}

	userIDSTR := md.Get("user_id")[0]

	userID, err := uuid.Parse(userIDSTR)
	if err != nil {
		return nil, fmt.Errorf("failed parse userid string")
	}

	post := &models.Post{
		ID:        uuid.New(),
		UserID:    userID,
		Content:   content,
		ImageURL:  imageURL,
		CreatedAt: time.Now(),
	}

	if err := s.repo.CreatePost(ctx, post); err != nil {
		return nil, fmt.Errorf("failed to create post: %w", err)
	}

	return post, nil
}

func (s *FeedService) GetAllPosts(ctx context.Context, page, pageSize int) ([]models.Post, int, error) {
	posts, totalPosts, err := s.repo.GetALLPosts(ctx, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get posts: %w", err)
	}

	return posts, totalPosts, nil
}
