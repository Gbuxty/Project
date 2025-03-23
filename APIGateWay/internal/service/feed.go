package service

import (
	"Project/APIGateWay/internal/domain"
	"Project/proto/gen"
	"context"
	"fmt"

	"google.golang.org/grpc"
)

type FeedService struct {
	client gen.FeedServiceClient
}

func NewFeedService(conn *grpc.ClientConn) *FeedService {
	return &FeedService{client: gen.NewFeedServiceClient(conn)}
}

func (s *FeedService) CreatePost(ctx context.Context, req *domain.CreatePostRequest) (*domain.PostResponse, error) {
	
	grpcReq := &gen.CreatePostRequest{
		Content:  req.Content,
		ImageUrl: req.ImageURL,
	}

	res, err := s.client.CreatePost(ctx, grpcReq)
	if err != nil {
		return nil, fmt.Errorf("Failed Create Post%w",err)
	}

	return &domain.PostResponse{
		Content:   res.Post.Content,
		ImageURL:  res.Post.ImageUrl,
		CreatedAt: res.Post.CreatedAt,
	}, nil

}

func (s *FeedService) GetAllPosts(ctx context.Context, page, pageSize int) ([]*domain.PostResponse, int, error) {
	grpcReq := &gen.GetAllPostsRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
	}

	res, err := s.client.GetAllPosts(ctx, grpcReq)
	if err != nil {
		return nil, 0, fmt.Errorf("Failed get all posts")
	}

	var posts []*domain.PostResponse
	for _, p := range res.Posts {
		posts = append(posts, &domain.PostResponse{

			Content:   p.Content,
			ImageURL:  p.ImageUrl,
			CreatedAt: p.CreatedAt,
		})
	}

	return posts, int(res.TotalPosts), nil
}
