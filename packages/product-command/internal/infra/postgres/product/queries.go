package product

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

func (r *productRepository) CheckSlugExists(ctx context.Context, slug string) (bool, error) {
	_, err := r.getQuerier(ctx).CountBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("infra:postgres: count by slug: %w", err)
	}

	return true, nil
}
