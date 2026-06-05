package product

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	domain_product "product-command-module/internal/domain/product"
	"product-command-module/internal/domain/shared"
)

func (r *productRepository) CheckSlugExists(ctx context.Context, slug string) (bool, error) {
	_, err := r.getQuerier(ctx).CountProductBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("infra:postgres: count by slug: %w", err)
	}

	return true, nil
}

func (r *productRepository) GetProductByID(ctx context.Context, productID shared.ProductID) (*domain_product.Product, error) {
	return nil, domain_product.ErrProductNotFound
}
