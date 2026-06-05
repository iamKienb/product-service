package product

import (
	"context"
	"errors"

	"product-command-module/internal/application/commands/create_product"
	"product-command-module/internal/application/port"
	"product-command-module/internal/domain/product"
)

func (s *productService) CreateProduct(ctx context.Context, cmd create_product.Command) (*create_product.Result, error) {
	if err := validateCreateProductCommand(cmd); err != nil {
		return nil, err
	}

	existing, err := s.findExistingProductBySlug(ctx, cmd)
	if err == nil {
		return createProductResultFromExisting(cmd, existing)
	}
	if !errors.Is(err, product.ErrProductNotFound) {
		return nil, err
	}

	newProduct := product.NewProduct(product.NewProductParams{
		ShopID:      cmd.ShopID,
		UserID:      cmd.UserID,
		Name:        cmd.Name,
		Slug:        cmd.Slug,
		Description: cmd.Description,
		Brand:       cmd.Brand,
		ThumbUrl:    cmd.ThumbUrl,
		VideoUrl:    cmd.VideoUrl,
		HasVariant:  cmd.HasVariant,
	})

	if cmd.HasVariant {
		for _, attr := range cmd.Attributes {
			if err := newProduct.AddAttribute(attr.Name, attr.Values); err != nil {
				return nil, err
			}
		}
	}

	skuItems := make([]create_product.SkuItem, 0, len(cmd.Variants))
	for index, variant := range cmd.Variants {
		isDefault := index == 0

		variantParam := product.ProductVariantParams{
			SkuCode:             variant.SkuCode,
			Price:               variant.Price,
			Currency:            variant.Currency,
			ImageUrl:            variant.ImageUrl,
			AttributeValueNames: variant.AttributeValueNames,
		}

		skuID, err := newProduct.AddVariant(variantParam, isDefault)
		if err != nil {
			return nil, err
		}
		skuItems = append(skuItems, create_product.SkuItem{
			SkuID:    skuID,
			Quantity: variant.Quantity,
		})
	}

	newProduct.MarkAsCreated()

	if err := s.txManager.WithTx(ctx, func(ctx context.Context) error {
		if err := s.productRepo.CreateProduct(ctx, newProduct); err != nil {
			return err
		}

		if productEvents := newProduct.FlushEvents(); len(productEvents) > 0 {
			if err := s.outboxService.Publish(ctx, port.OutboxParam{
				AggregateID:   newProduct.ID.RawID(),
				AggregateType: newProduct.Type(),
				Events:        productEvents,
			}); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		if errors.Is(err, product.ErrProductSlugTaken) {
			existing, readErr := s.productRepo.GetProductByShopAndSlug(ctx, cmd.ShopID, cmd.Slug)
			if readErr == nil {
				return createProductResultFromExisting(cmd, existing)
			}
		}
		return nil, err
	}

	bgCtx := context.WithoutCancel(ctx)
	go func() {
		_ = s.productCache.AddSlugToBloomFilter(bgCtx, cmd.Slug)
	}()

	return &create_product.Result{
		ShopID:    cmd.ShopID,
		ProductID: newProduct.ID,
		SkuItems:  skuItems,
	}, nil
}

func validateCreateProductCommand(cmd create_product.Command) error {
	if len(cmd.Variants) == 0 {
		return product.ErrNoSKU
	}
	for _, variant := range cmd.Variants {
		if variant.Quantity < 0 {
			return product.ErrInvalidSkuQuantity
		}
	}
	return nil
}

func (s *productService) findExistingProductBySlug(ctx context.Context, cmd create_product.Command) (*product.Product, error) {
	exists, err := s.productCache.GetSlugFromBloomFilter(ctx, cmd.Slug)
	if err != nil {
		return nil, err
	}

	if exists == 0 {
		return nil, product.ErrProductNotFound
	}

	isDuplicateSlug, err := s.productRepo.CheckSlugExists(ctx, cmd.ShopID, cmd.Slug)
	if err != nil {
		return nil, err
	}
	if !isDuplicateSlug {
		return nil, product.ErrProductNotFound
	}

	return s.productRepo.GetProductByShopAndSlug(ctx, cmd.ShopID, cmd.Slug)
}

func createProductResultFromExisting(cmd create_product.Command, existing *product.Product) (*create_product.Result, error) {
	if existing.ShopID != cmd.ShopID {
		return nil, product.ErrProductSlugTaken
	}

	variantsByCode := make(map[string]product.ProductVariant, len(existing.Variants))
	for _, variant := range existing.Variants {
		variantsByCode[variant.SkuCode] = variant
	}

	skuItems := make([]create_product.SkuItem, 0, len(cmd.Variants))
	for _, requested := range cmd.Variants {
		variant, exists := variantsByCode[requested.SkuCode]
		if !exists {
			return nil, product.ErrProductSlugTaken
		}
		skuItems = append(skuItems, create_product.SkuItem{
			SkuID:    variant.SkuID,
			Quantity: requested.Quantity,
		})
	}

	return &create_product.Result{
		ShopID:    existing.ShopID,
		ProductID: existing.ID,
		SkuItems:  skuItems,
	}, nil
}
