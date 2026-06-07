package cache

import (
	"context"
	"fmt"
	"product-command-module/internal/application/port"

	redisx "github.com/iamKienb/go-core/redis"
)

const slugKey = "product-command:slug:%s"

type productCache struct {
	cache redisx.RedisXService
}

func NewProductCache(service redisx.RedisXService) port.ProductCache {
	return &productCache{cache: service}
}

func (c *productCache) IsSlugKnown(ctx context.Context, slug string) (bool, error) {
	return c.cache.Exists(ctx, fmt.Sprintf(slugKey, slug))
}

func (c *productCache) RememberSlug(ctx context.Context, slug string) error {
	return c.cache.Set(ctx, fmt.Sprintf(slugKey, slug), "exists", 0)
}
