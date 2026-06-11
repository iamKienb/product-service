package service

import (
	"context"
	"fmt"
	"product-query-module/internal/application/queries/get_product_by_sku_ids"
	"product-query-module/internal/application/service/models"

	"github.com/elastic/go-elasticsearch/v8"
)

type ProductService interface {
	SearchBySkuIDs(ctx context.Context, qry get_product_by_sku_ids.Query) ([]*get_product_by_sku_ids.Result, error)
}

type productService struct {
	index    string
	esClient *elasticsearch.TypedClient
}

const (
	productStatusActive  = "ACTIVE"
	variantInnerHitName  = "sku_variants"
	errCheckoutSkuAbsent = "checkout product validation failed: sku %s not found or invalid for shop %s"
	errCheckoutSkuStatus = "checkout product validation failed: sku %s is currently status %s"
	errShopNotFound      = "shop is not found"
)

func NewProductService(esClient *elasticsearch.TypedClient, index string) *productService {
	return &productService{
		esClient: esClient,
		index:    index,
	}
}

func (s *productService) SearchBySkuIDs(ctx context.Context, qry get_product_by_sku_ids.Query) ([]*get_product_by_sku_ids.Result, error) {
	if len(qry.SkuIDs) == 0 {
		return nil, nil
	}

	variantBuilder := NewQueryBuilder().
		FilterTerms("variants.id", qry.SkuIDs).
		FilterTerm("variants.status", productStatusActive)

	searchQueryBody := NewQueryBuilder().
		MustTerm("shop_id", qry.ShopID).
		FilterTerm("status", productStatusActive).
		Nested("variants", variantBuilder, variantInnerHitName).
		Build()

	esResult, err := SearchDocuments[models.Product](ctx, s.esClient, s.index, searchQueryBody)
	if err != nil {
		return nil, err
	}

	if len(esResult.Hits) == 0 {
		return nil, fmt.Errorf(errShopNotFound)
	}

	foundVariantsMap := make(map[string]*get_product_by_sku_ids.Result)

	for _, rawResult := range esResult.Hits {
		productDoc := rawResult.Source
		if productDoc.Status != productStatusActive {
			continue
		}

		matchedVariants, err := DecodeInnerHits[models.ProductVariant](rawResult.InnerHits, variantInnerHitName)
		if err != nil {
			return nil, fmt.Errorf("failed to decode nested inner hits variants: %w", err)
		}

		for _, variant := range matchedVariants {
			foundVariantsMap[variant.SkuID] = &get_product_by_sku_ids.Result{
				ShopID:      productDoc.ShopID,
				ProductID:   productDoc.ID,
				ProductName: productDoc.Name,
				SkuID:       variant.SkuID,
				SkuCode:     variant.SkuCode,
				Price:       variant.Price,
				ImageURL:    variant.ImageURL,
				Status:      variant.Status,
			}
		}

	}

	skuResults := make([]*get_product_by_sku_ids.Result, 0, len(qry.SkuIDs))
	for _, reqSkuID := range qry.SkuIDs {
		cachedResult, exists := foundVariantsMap[reqSkuID]
		if !exists {
			return nil, fmt.Errorf(errCheckoutSkuAbsent, reqSkuID, qry.ShopID)
		}

		if cachedResult.Status != productStatusActive {
			return nil, fmt.Errorf(errCheckoutSkuStatus, reqSkuID, cachedResult.Status)
		}

		skuResults = append(skuResults, cachedResult)
	}

	return skuResults, nil
}
