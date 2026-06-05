package product

import (
	"context"
	"fmt"

	"product-command-module/db/repository"
	domain_product "product-command-module/internal/domain/product"
	"product-command-module/internal/domain/shared"

	"github.com/iamKienb/go-core/postgres/conv"
	"github.com/jackc/pgx/v5/pgtype"
)

func (r *productRepository) CreateProduct(ctx context.Context, product *domain_product.Product) error {
	if err := r.getQuerier(ctx).CreateProduct(ctx, toCreateProductParams(product)); err != nil {
		return fmt.Errorf("infra: create product: %w", err)
	}

	if len(product.Attributes) > 0 {
		if err := r.getQuerier(ctx).BatchCreateAttributes(ctx, toBatchCreateAttributesParams(product)); err != nil {
			return fmt.Errorf("infra: create product attributes: %w", err)
		}
	}

	if len(product.AttributeValues) > 0 {
		if err := r.getQuerier(ctx).BatchCreateAttributeValues(ctx, toBatchCreateAttributeValuesParams(product.AttributeValues)); err != nil {
			return fmt.Errorf("infra: create product attribute values: %w", err)
		}
	}

	if len(product.Variants) > 0 {
		if err := r.getQuerier(ctx).BatchCreateVariants(ctx, toBatchCreateVariantsParams(product)); err != nil {
			return fmt.Errorf("infra: create product variants: %w", err)
		}

		linkParams := toBatchLinkVariantAttributesParams(product.Variants)
		if len(linkParams.SkuIds) > 0 {
			if err := r.getQuerier(ctx).BatchLinkVariantAttributes(ctx, linkParams); err != nil {
				return fmt.Errorf("infra: link product variant attributes: %w", err)
			}
		}
	}

	return nil
}

func (r *productRepository) DeleteProduct(ctx context.Context, productID shared.ProductID) error {
	return domain_product.ErrProductNotFound
}

func toCreateProductParams(product *domain_product.Product) repository.CreateProductParams {
	return repository.CreateProductParams{
		ID:          conv.UUID(product.ID),
		ShopID:      conv.UUID(product.ShopID),
		Name:        product.Name,
		Slug:        product.Slug,
		Description: product.Description,
		Brand:       product.Brand,
		ThumbUrl:    product.ThumbUrl,
		VideoUrl:    product.VideoUrl,
		PriceMin:    product.PriceMin,
		PriceMax:    product.PriceMax,
		Status:      string(product.Status),
		HasVariant:  product.HasVariant,
		CreatedBy:   conv.UUID(product.CreatedBy),
		CreatedAt:   conv.TimeStampZ(&product.CreatedAt),
	}
}

func toBatchCreateAttributesParams(product *domain_product.Product) repository.BatchCreateAttributesParams {
	params := repository.BatchCreateAttributesParams{
		Ids:       make([]pgtype.UUID, 0, len(product.Attributes)),
		ProductID: conv.UUID(product.ID),
		Names:     make([]string, 0, len(product.Attributes)),
	}

	for _, attribute := range product.Attributes {
		params.Ids = append(params.Ids, conv.UUID(attribute.ID))
		params.Names = append(params.Names, attribute.Name)
	}

	return params
}

func toBatchCreateAttributeValuesParams(values []domain_product.AttributeValue) repository.BatchCreateAttributeValuesParams {
	params := repository.BatchCreateAttributeValuesParams{
		Ids:                 make([]pgtype.UUID, 0, len(values)),
		ProductAttributeIds: make([]pgtype.UUID, 0, len(values)),
		Values:              make([]string, 0, len(values)),
	}

	for _, value := range values {
		params.Ids = append(params.Ids, conv.UUID(value.ID))
		params.ProductAttributeIds = append(params.ProductAttributeIds, conv.UUID(value.AttributeID))
		params.Values = append(params.Values, value.Name)
	}

	return params
}

func toBatchCreateVariantsParams(product *domain_product.Product) repository.BatchCreateVariantsParams {
	params := repository.BatchCreateVariantsParams{
		SkuIds:     make([]pgtype.UUID, 0, len(product.Variants)),
		ProductID:  conv.UUID(product.ID),
		ShopID:     conv.UUID(product.ShopID),
		SkuCodes:   make([]string, 0, len(product.Variants)),
		Prices:     make([]int32, 0, len(product.Variants)),
		Currencies: make([]string, 0, len(product.Variants)),
		ImageUrls:  make([]string, 0, len(product.Variants)),
		Status:     make([]string, 0, len(product.Variants)),
		IsDefaults: make([]bool, 0, len(product.Variants)),
		CreatedBys: make([]pgtype.UUID, 0, len(product.Variants)),
		CreatedAts: make([]pgtype.Timestamptz, 0, len(product.Variants)),
	}

	for _, variant := range product.Variants {
		params.SkuIds = append(params.SkuIds, conv.UUID(variant.SkuID))
		params.SkuCodes = append(params.SkuCodes, variant.SkuCode)
		params.Prices = append(params.Prices, int32(variant.Price))
		params.Currencies = append(params.Currencies, variant.Currency)
		params.ImageUrls = append(params.ImageUrls, variant.ImageUrl)
		params.Status = append(params.Status, string(domain_product.VariantStatusActive))
		params.IsDefaults = append(params.IsDefaults, variant.IsDefault)
		params.CreatedBys = append(params.CreatedBys, conv.UUID(variant.CreatedBy))
		params.CreatedAts = append(params.CreatedAts, conv.TimeStampZ(&variant.CreatedAt))
	}

	return params
}

func toBatchLinkVariantAttributesParams(variants []domain_product.ProductVariant) repository.BatchLinkVariantAttributesParams {
	params := repository.BatchLinkVariantAttributesParams{
		SkuIds:            []pgtype.UUID{},
		AttributeValueIds: []pgtype.UUID{},
	}

	for _, variant := range variants {
		for _, attributeValueID := range variant.AttributeValueIDs {
			params.SkuIds = append(params.SkuIds, conv.UUID(variant.SkuID))
			params.AttributeValueIds = append(params.AttributeValueIds, conv.UUID(attributeValueID))
		}
	}

	return params
}
