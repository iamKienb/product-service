package port

import (
	"context"
)

type ProductCache interface {
	GetSlugFromBloomFilter(ctx context.Context, slug string) (int, error)
	AddSlugToBloomFilter(ctx context.Context, slug string) error
}
