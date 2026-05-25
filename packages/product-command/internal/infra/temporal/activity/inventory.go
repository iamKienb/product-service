package activity

import (
	"context"
	"product-command-module/internal/application/port"
)

type InventoryActivity struct {
	client port.InventoryClient
}

func NewInventoryActivity(client port.InventoryClient) *InventoryActivity {
	return &InventoryActivity{client: client}
}

func (a *InventoryActivity) CreateInventory(ctx context.Context, req port.CreateInventoryRequest) error {
	return a.client.CreateInventory(ctx, req)
}

func (a *InventoryActivity) DeleteInventory(ctx context.Context, skuIDs []string) error {
	return a.client.DeleteInventory(ctx, skuIDs)
}
