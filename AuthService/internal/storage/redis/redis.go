package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type Client struct {
	client *redis.Client
}

func NewClient(addr string) *Client {
	return &Client{
		client: redis.NewClient(&redis.Options{
			Addr: addr,
			DB:   0,
		})}
}

type TokenCacheRepository struct{//зачем мне это??
	
}

func getAccessTokenKey(userID string) string {
	return fmt.Sprintf("access_token:%s", userID)
}

func (c *Client) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	key=getAccessTokenKey(key)
	return c.client.Set(ctx, key, value, ttl).Err()
}

func (c *Client) Get(ctx context.Context, key string) (string, error) {
	return c.client.Get(ctx, key).Result()
}

func (c *Client) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}
