package product

import (
	"context"
	"product-command-module/internal/application/port"
	"product-command-module/internal/domain/product"
	domain_shared "product-command-module/internal/domain/shared"
)

func (s *productService) RollbackProduct(ctx context.Context, productID string) error {
	parseProductID, err := domain_shared.ParseToRawID[domain_shared.ProductID](productID)
	if err != nil {
		return s.wrapError(product.ErrInvalidProductID)
	}

	prod, err := s.productRepo.GetProductByID(ctx, parseProductID)
	if err != nil {
		return s.wrapError(err)
	}

	prod.MarkAsDeleted()

	if err := s.txManager.WithTx(ctx, func(txCtx context.Context) error {
		if err := s.productRepo.DeleteProduct(txCtx, parseProductID); err != nil {
			return err
		}

		productEvents := prod.FlushEvents()
		if len(productEvents) > 0 {
			if err := s.outboxService.Publish(txCtx, port.OutboxParam{
				AggregateID:   prod.ID.RawID(),
				AggregateType: prod.Type(),
				Events:        productEvents,
			}); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return s.wrapError(err)
	}

	return nil
}
