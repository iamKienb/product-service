package port

import (
	"context"
)

type ProductCache interface {
	IsSlugKnown(ctx context.Context, slug string) (bool, error)
	RememberSlug(ctx context.Context, slug string) error
}
