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

	product, err := r.toDomainProductWithRelations(ctx, row)
	if err != nil {
		return nil, err
	}
	return product, nil
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

	product, err := r.toDomainProductWithRelations(ctx, row)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (r *productRepository) toDomainProductWithRelations(ctx context.Context, row repository.Product) (*domain_product.Product, error) {
	variants, err := r.getQuerier(ctx).ListVariantsByProductID(ctx, row.ID)
	if err != nil {
		return nil, fmt.Errorf("infra: list product variants: %w", err)
	}

	attributes, err := r.getQuerier(ctx).ListAttributesByProductID(ctx, row.ID)
	if err != nil {
		return nil, fmt.Errorf("infra: list product attributes: %w", err)
	}

	attributeValues, err := r.getQuerier(ctx).ListAttributeValuesByProductID(ctx, row.ID)
	if err != nil {
		return nil, fmt.Errorf("infra: list product attribute values: %w", err)
	}

	variantAttributeValues, err := r.getQuerier(ctx).ListVariantAttributeValuesByProductID(ctx, row.ID)
	if err != nil {
		return nil, fmt.Errorf("infra: list product variant attribute values: %w", err)
	}

	return toDomainProduct(row, variants, attributes, attributeValues, variantAttributeValues), nil
}

func toDomainProduct(
	row repository.Product,
	variants []repository.ProductVariant,
	attributes []repository.ProductAttribute,
	attributeValues []repository.AttributeValue,
	variantAttributeValues []repository.ProductAttributeValue,
) *domain_product.Product {
	product := &domain_product.Product{
		ID:              shared.ProductID(row.ID.Bytes),
		ShopID:          shared.ShopID(row.ShopID.Bytes),
		Name:            row.Name,
		Slug:            row.Slug,
		Description:     row.Description,
		Brand:           row.Brand,
		ThumbUrl:        row.ThumbUrl,
		VideoUrl:        row.VideoUrl,
		PriceMin:        row.PriceMin,
		PriceMax:        row.PriceMax,
		Status:          domain_product.ProductStatus(row.Status),
		HasVariant:      row.HasVariant.Bool,
		CreatedBy:       shared.UserID(row.CreatedBy.Bytes),
		UpdatedBy:       userIDPointer(row.UpdatedBy),
		CreatedAt:       row.CreatedAt.Time,
		UpdatedAt:       timePointer(row.UpdatedAt),
		Attributes:      toDomainAttributes(attributes),
		AttributeValues: toDomainAttributeValues(attributeValues),
		Variants:        toDomainVariants(variants, variantAttributeValues),
	}
	return product
}

func toDomainAttributes(rows []repository.ProductAttribute) []domain_product.Attribute {
	attributes := make([]domain_product.Attribute, 0, len(rows))
	for _, row := range rows {
		attributes = append(attributes, domain_product.Attribute{
			ID:        shared.AttributeID(row.ID.Bytes),
			ProductID: shared.ProductID(row.ProductID.Bytes),
			Name:      row.Name,
		})
	}
	return attributes
}

func toDomainAttributeValues(rows []repository.AttributeValue) []domain_product.AttributeValue {
	values := make([]domain_product.AttributeValue, 0, len(rows))
	for _, row := range rows {
		values = append(values, domain_product.AttributeValue{
			ID:          shared.AttributeValueID(row.ID.Bytes),
			AttributeID: shared.AttributeID(row.ProductAttributeID.Bytes),
			Name:        row.Name,
		})
	}
	return values
}

func toDomainVariants(rows []repository.ProductVariant, attributeRows []repository.ProductAttributeValue) []domain_product.ProductVariant {
	attributeValueIDsBySku := make(map[shared.SkuID][]shared.AttributeValueID, len(rows))
	for _, row := range attributeRows {
		skuID := shared.SkuID(row.SkuID.Bytes)
		attributeValueIDsBySku[skuID] = append(attributeValueIDsBySku[skuID], shared.AttributeValueID(row.AttributeValueID.Bytes))
	}

	variants := make([]domain_product.ProductVariant, 0, len(rows))
	for _, row := range rows {
		skuID := shared.SkuID(row.SkuID.Bytes)
		variants = append(variants, domain_product.ProductVariant{
			SkuID:             skuID,
			ProductID:         shared.ProductID(row.ProductID.Bytes),
			ShopID:            shared.ShopID(row.ShopID.Bytes),
			SkuCode:           row.SkuCode,
			Price:             row.Price,
			Currency:          row.Currency,
			ImageUrl:          row.ImageUrl,
			Status:            domain_product.VariantStatus(row.Status),
			IsDefault:         row.IsDefault.Bool,
			AttributeValueIDs: attributeValueIDsBySku[skuID],
			CreatedBy:         shared.UserID(row.CreatedBy.Bytes),
			UpdatedBy:         userIDPointer(row.UpdatedBy),
			CreatedAt:         row.CreatedAt.Time,
			UpdatedAt:         timePointer(row.UpdatedAt),
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
