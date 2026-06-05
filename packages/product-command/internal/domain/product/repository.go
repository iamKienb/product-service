package product

import (
	"context"
	"product-command-module/internal/domain/shared"
)

type QueryRepository interface {
	GetProductByID(ctx context.Context, productID shared.ProductID) (*Product, error)
	GetProductByShopAndSlug(ctx context.Context, shopID shared.ShopID, slug string) (*Product, error)
	CheckSlugExists(ctx context.Context, shopID shared.ShopID, slug string) (bool, error)
}

type CommandRepository interface {
	CreateProduct(ctx context.Context, product *Product) error
	// UpdateVariant(ctx context.Context, skuID shared.SkuID, variant Variant) error

	DeleteProduct(ctx context.Context, productID shared.ProductID) error
}

type Repository interface {
	QueryRepository
	CommandRepository
}
