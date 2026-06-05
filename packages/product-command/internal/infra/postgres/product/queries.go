package product

import (
	"context"
	"errors"
	"fmt"
	"time"

	"product-command-module/db/repository"
	domain_product "product-command-module/internal/domain/product"
	"product-command-module/internal/domain/shared"

	"github.com/iamKienb/go-core/postgres/conv"
	pgxv5 "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func (r *productRepository) CheckSlugExists(ctx context.Context, shopID shared.ShopID, slug string) (bool, error) {
	exists, err := r.getQuerier(ctx).CheckProductSlugExists(ctx, repository.CheckProductSlugExistsParams{
		ShopID: conv.UUID(shopID),
		Slug:   slug,
	})
	if err != nil {
		return false, fmt.Errorf("infra: check product slug exists: %w", err)
	}

	return exists, nil
}

func (r *productRepository) GetProductByID(ctx context.Context, productID shared.ProductID) (*domain_product.Product, error) {
	row, err := r.getQuerier(ctx).GetProductByID(ctx, conv.UUID(productID))
	if err != nil {
		if errors.Is(err, pgxv5.ErrNoRows) {
			return nil, domain_product.ErrProductNotFound
		}
		return nil, fmt.Errorf("infra: get product by id: %w", err)
	}

	variants, err := r.getQuerier(ctx).ListVariantsByProductID(ctx, conv.UUID(productID))
	if err != nil {
		return nil, fmt.Errorf("infra: list product variants: %w", err)
	}

	return toDomainProduct(row, variants), nil
}

func (r *productRepository) GetProductByShopAndSlug(ctx context.Context, shopID shared.ShopID, slug string) (*domain_product.Product, error) {
	row, err := r.getQuerier(ctx).GetProductByShopAndSlug(ctx, repository.GetProductByShopAndSlugParams{
		ShopID: conv.UUID(shopID),
		Slug:   slug,
	})
	if err != nil {
		if errors.Is(err, pgxv5.ErrNoRows) {
			return nil, domain_product.ErrProductNotFound
		}
		return nil, fmt.Errorf("infra: get product by shop and slug: %w", err)
	}

	variants, err := r.getQuerier(ctx).ListVariantsByProductID(ctx, row.ID)
	if err != nil {
		return nil, fmt.Errorf("infra: list product variants by slug: %w", err)
	}

	return toDomainProduct(row, variants), nil
}

func toDomainProduct(row repository.Product, variants []repository.ProductVariant) *domain_product.Product {
	product := &domain_product.Product{
		ID:          shared.ProductID(row.ID.Bytes),
		ShopID:      shared.ShopID(row.ShopID.Bytes),
		Name:        row.Name,
		Slug:        row.Slug,
		Description: row.Description,
		Brand:       row.Brand,
		ThumbUrl:    row.ThumbUrl,
		VideoUrl:    row.VideoUrl,
		PriceMin:    row.PriceMin,
		PriceMax:    row.PriceMax,
		Status:      domain_product.ProductStatus(row.Status),
		HasVariant:  row.HasVariant.Bool,
		CreatedBy:   shared.UserID(row.CreatedBy.Bytes),
		UpdatedBy:   userIDPointer(row.UpdatedBy),
		CreatedAt:   row.CreatedAt.Time,
		UpdatedAt:   timePointer(row.UpdatedAt),
		Variants:    toDomainVariants(variants),
	}
	return product
}

func toDomainVariants(rows []repository.ProductVariant) []domain_product.ProductVariant {
	variants := make([]domain_product.ProductVariant, 0, len(rows))
	for _, row := range rows {
		variants = append(variants, domain_product.ProductVariant{
			SkuID:     shared.SkuID(row.SkuID.Bytes),
			ProductID: shared.ProductID(row.ProductID.Bytes),
			ShopID:    shared.ShopID(row.ShopID.Bytes),
			SkuCode:   row.SkuCode,
			Price:     row.Price,
			Currency:  row.Currency,
			ImageUrl:  row.ImageUrl,
			Status:    domain_product.VariantStatus(row.Status),
			IsDefault: row.IsDefault.Bool,
			CreatedBy: shared.UserID(row.CreatedBy.Bytes),
			UpdatedBy: userIDPointer(row.UpdatedBy),
			CreatedAt: row.CreatedAt.Time,
			UpdatedAt: timePointer(row.UpdatedAt),
		})
	}
	return variants
}

func userIDPointer(value pgtype.UUID) *shared.UserID {
	if !value.Valid {
		return nil
	}
	userID := shared.UserID(value.Bytes)
	return &userID
}

func timePointer(value pgtype.Timestamptz) *time.Time {
	if !value.Valid {
		return nil
	}
	result := value.Time
	return &result
}
