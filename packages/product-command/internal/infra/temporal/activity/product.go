package activity

import (
	"context"
	"product-command-module/internal/application/commands/create_product"
)

type ProductService interface {
	CreateProduct(ctx context.Context, cmd create_product.Command) (*create_product.Result, error)
	RollbackProduct(ctx context.Context, productID string) error
}

type ProductActivity struct {
	service ProductService
}

func NewProductActivity(service ProductService) *ProductActivity {
	return &ProductActivity{service: service}
}

func (a *ProductActivity) CreateProduct(ctx context.Context, cmd create_product.Command) (*create_product.Result, error) {
	return a.service.CreateProduct(ctx, cmd)
}

func (a *ProductActivity) RollbackProduct(ctx context.Context, productID string) error {
	return a.service.RollbackProduct(ctx, productID)
}
