package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisClient(redisOptions *redis.Options) *RedisClient {
	client := redis.NewClient(redisOptions)
	return &RedisClient{client: client, ctx: context.Background()}
}

func (r *RedisClient) Set(key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal cache value: %w", err)
	}

	return r.client.Set(r.ctx, key, data, time.Minute*5).Err()
}

func (r *RedisClient) Get(key string) (string, error) {
	return r.client.Get(r.ctx, key).Result()
}

func (r *RedisClient) Close() {
	r.client.Close()
}
