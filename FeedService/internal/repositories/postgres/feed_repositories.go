package postgres

import (
	"Project/FeedService/internal/domain/models"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type FeedRepositories struct {
	db *pgxpool.Pool
}


func NewPostRepositories(db *pgxpool.Pool) (*FeedRepositories, error) {
	return &FeedRepositories{
		db: db,
	}, nil
}

func (r *FeedRepositories) CloseDb() {
	r.db.Close()
}

func (r *FeedRepositories) CreatePost(ctx context.Context, post *models.Post) error {
    query := `
        INSERT INTO posts (user_id, content, image_url, created_at)
        VALUES ($1, $2, $3, $4)
    `
    _, err := r.db.Exec(ctx, query, post.UserID, post.Content, post.ImageURL, post.CreatedAt)
    if err != nil {
        return fmt.Errorf("failed to create post: %w", err)
    }
    return nil
}

func (r *FeedRepositories) GetALLPosts(ctx context.Context, page, pageSize int) ([]models.Post, int, error) {
    offset := (page - 1) * pageSize

   
    query := `
        SELECT id, user_id, content, image_url, created_at
        FROM posts
        ORDER BY created_at DESC
        LIMIT $1 OFFSET $2
    `
    rows, err := r.db.Query(ctx, query, pageSize, offset)
    if err != nil {
        return nil, 0, fmt.Errorf("failed to get posts: %w", err)
    }
    defer rows.Close()

    var posts []models.Post
    for rows.Next() {
        var post models.Post
        if err := rows.Scan(&post.ID, &post.UserID, &post.Content, &post.ImageURL, &post.CreatedAt); err != nil {
            return nil, 0, fmt.Errorf("failed to scan post: %w", err)
        }
        posts = append(posts, post)
    }
	
    var totalPosts int
    countQuery := `SELECT COUNT(*) FROM posts`
    if err := r.db.QueryRow(ctx, countQuery).Scan(&totalPosts); err != nil {
        return nil, 0, fmt.Errorf("failed to count posts: %w", err)
    }

    return posts, totalPosts, nil
}