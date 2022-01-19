package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gerifield/service-scaling/scale2/model"
	"github.com/go-redis/redis/v8"
)

type RedisCache struct {
	backend *redis.Client
}

func NewRedis(backend *redis.Client) *RedisCache {
	return &RedisCache{
		backend: backend,
	}
}

func (rc *RedisCache) GetMessages(ctx context.Context, key string) ([]model.Message, error) {
	resp, err := rc.backend.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var ret []model.Message
	err = json.Unmarshal([]byte(resp), &ret)

	return ret, err
}

func (rc *RedisCache) SaveMessages(ctx context.Context, key string, value []model.Message, ttl time.Duration) error {
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return rc.backend.Set(ctx, key, string(b), ttl).Err()
}

func (rc *RedisCache) Invalidate(ctx context.Context, key string) error {
	return rc.backend.Del(ctx, key).Err()
}
