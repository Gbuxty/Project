package handlers

import (
	"Project/FeedService/internal/service"
	"Project/proto/gen"
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
)

type FeedHandlers struct {
	service   *service.FeedService
	logger    *zap.Logger
	secretKey string
	gen.UnimplementedFeedServiceServer
}

func NewFeedHandlers(service *service.FeedService, logger *zap.Logger, secretKey string) *FeedHandlers {
	return &FeedHandlers{
		service:   service,
		logger:    logger,
		secretKey: secretKey,
	}
}

func (h *FeedHandlers) CreatePost(ctx context.Context, req *gen.CreatePostRequest) (*gen.CreatePostResponse, error) {
	h.logger.Info("Create Post", zap.String("Content", req.Content))

	post, err := h.service.CreatePost(ctx, req.Content, req.ImageUrl)
	if err != nil {
		h.logger.Error("Failed to create post", zap.Error(err))
		return nil, fmt.Errorf("failed to create post: %w", err)
	}

	h.logger.Info("Create Post Successful", zap.String("Content", req.Content))
	return &gen.CreatePostResponse{
		Post: &gen.Post{
			Content:   post.Content,
			ImageUrl:  post.ImageURL,
			CreatedAt: post.CreatedAt.Format(time.RFC3339),
		},
	}, nil

}

func (h *FeedHandlers) GetAllPosts(ctx context.Context, req *gen.GetAllPostsRequest) (*gen.GetAllPostsResponse, error) {
	posts, totalPosts, err := h.service.GetAllPosts(ctx, int(req.Page), int(req.PageSize))
	if err != nil {
		h.logger.Error("Failed to get posts", zap.Error(err))
		return nil, fmt.Errorf("failed to get posts: %w", err)
	}

	var responsePosts []*gen.Post
	for _, post := range posts {
		responsePosts = append(responsePosts, &gen.Post{

			Content:   post.Content,
			ImageUrl:  post.ImageURL,
			CreatedAt: post.CreatedAt.Format(time.RFC3339),
		})
	}

	return &gen.GetAllPostsResponse{
		Posts:      responsePosts,
		TotalPosts: int32(totalPosts),
	}, nil
}
