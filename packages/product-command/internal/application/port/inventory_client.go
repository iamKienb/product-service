package port

import (
	"context"
)

type SkuItem struct {
	SkuID    string
	Quantity int32
}

type CreateInventoryRequest struct {
	ShopID    string
	ProductID string
	Items     []SkuItem
}

type InventoryClient interface {
	CreateInventory(ctx context.Context, req CreateInventoryRequest) error
	DeleteInventory(ctx context.Context, skuIDs []string) error
}
