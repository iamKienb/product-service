package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"product-worker-module/internal/application/port"
	"time"

	"github.com/redis/go-redis/v9"
)

type workerCache struct {
	client *redis.Client
}

func NewWorkerCache(client *redis.Client) port.WorkerCache {
	return &workerCache{
		client: client,
	}
}

func (c *workerCache) SetNx(ctx context.Context, key string, data any, ttl time.Duration) (bool, error) {
	data, err := json.Marshal(data)
	if err != nil {
		return false, fmt.Errorf("marshal data %w", err)
	}

	isNew, err := c.client.SetNX(ctx, key, data, ttl).Result()
	if err != nil {
		return false, fmt.Errorf("set nx %w", err)
	}

	return isNew, nil
}
