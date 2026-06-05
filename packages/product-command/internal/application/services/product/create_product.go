package product

import (
	"context"
	"product-command-module/internal/application/commands/create_product"
	"product-command-module/internal/application/port"
	"product-command-module/internal/domain/product"
)

func (s *productService) CreateProduct(ctx context.Context, cmd create_product.Command) (*create_product.Result, error) {
	if err := s.checkSlugAvailable(ctx, cmd.Slug); err != nil {
		return nil, s.wrapError(err)
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
			newProduct.AddAttribute(attr.Name, attr.Values)
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

		skuID := newProduct.AddVariant(variantParam, isDefault)
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
		return nil, s.wrapError(err)
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

func (s *productService) checkSlugAvailable(ctx context.Context, slug string) error {
	exists, err := s.productCache.GetSlugFromBloomFilter(ctx, slug)
	if err != nil {
		return err
	}

	if exists == 0 {
		return nil
	}

	isDuplicateSlug, err := s.productRepo.CheckSlugExists(ctx, slug)
	if err != nil {
		return err
	}
	if isDuplicateSlug {
		return product.ErrProductSlugTaken
	}

	return nil
}
