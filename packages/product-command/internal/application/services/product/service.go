package product

import (
	"context"
	"product-command-module/internal/application/commands/create_product"
	"product-command-module/internal/application/port"
	"product-command-module/internal/application/services/outbox"
	"product-command-module/internal/domain/product"
)

type Service interface {
	CreateProduct(ctx context.Context, cmd create_product.Command) (*create_product.Result, error)
	RollbackProduct(ctx context.Context, productID string) error
}

type productService struct {
	productRepo  product.Repository
	productCache port.ProductCache

	outboxService outbox.Service
	txManager     port.TxManager
}

func NewProductService(
	productRepo product.Repository,
	productCache port.ProductCache,

	outboxService outbox.Service,
	txManager port.TxManager,
) Service {
	return &productService{
		productRepo:   productRepo,
		productCache:  productCache,
		outboxService: outboxService,
		txManager:     txManager,
	}
}
