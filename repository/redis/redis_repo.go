package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/yonesbeheshti/ai-gateway/adapter/redis"
	"github.com/yonesbeheshti/ai-gateway/entity"
)

//const (
//	cachePrefix string = "chat:cache:"
//)

type Cache struct {
	redisAdapter *redisadapter.RedisAdapter
}

func New(redisAdapter *redisadapter.RedisAdapter) *Cache {
	return &Cache{redisAdapter: redisAdapter}
}

func (c *Cache) Set(ctx context.Context, key string, resp *entity.ChatResponse, ttl time.Duration) error {
	data, err := json.Marshal(resp)
	if err != nil {
		return fmt.Errorf("marshal response for cache error: %w", err)
	}
	return c.redisAdapter.Client().Set(ctx, key, data, ttl).Err()

}
func (c *Cache) Get(ctx context.Context, key string) (*entity.ChatResponse, bool, error) {
	val, err := c.redisAdapter.Client().Get(ctx, key).Bytes()
	if errors.Is(err, redis.Nil) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, fmt.Errorf("redis get error: %w", err)
	}
	var resp entity.ChatResponse
	if err := json.Unmarshal(val, &resp); err != nil {
		return nil, false, fmt.Errorf("unmarshal cached response error: %w", err)
	}
	return &resp, true, nil
}
