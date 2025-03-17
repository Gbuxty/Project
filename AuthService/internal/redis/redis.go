package redis

import (
	"context"
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

func (c *Client) Set(ctx context.Context,key string,value interface{},ttl time.Duration)error{
	return c.client.Set(ctx, key, value, ttl).Err()
}

func (c *Client) Get(ctx context.Context, key string) (string, error) {
    return c.client.Get(ctx, key).Result()
}

func (c *Client) Delete(ctx context.Context, key string) error {
    return c.client.Del(ctx, key).Err()
}