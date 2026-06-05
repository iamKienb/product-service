package grpc_client

import (
	"context"
	"net/http"

	"product-command-module/internal/application/port"

	"connectrpc.com/connect"
	"github.com/iamKienb/api-contract/gen/inventory"
	"github.com/iamKienb/api-contract/gen/inventory/inventoryconnect"
)

type inventoryClient struct {
	client inventoryconnect.InventoryCommandServiceClient
}

func NewInventoryClient(httpClient *http.Client, baseURL string) port.InventoryClient {
	return &inventoryClient{
		client: inventoryconnect.NewInventoryCommandServiceClient(httpClient, baseURL),
	}
}

func (c *inventoryClient) CreateInventories(ctx context.Context, req port.CreateInventoryRequest) error {
	items := make([]*inventory.InventoryItem, 0, len(req.Items))
	for _, item := range req.Items {
		items = append(items, &inventory.InventoryItem{
			SkuId:    item.SkuID,
			Quantity: item.Quantity,
		})
	}

	payload := connect.NewRequest(&inventory.CreateInventoriesRequest{
		ShopId: req.ShopID,
		Items:  items,
	})

	_, err := c.client.CreateInventories(ctx, payload)

	return err
}

func (c *inventoryClient) DeleteInventory(ctx context.Context, skuIDs []string) error {
	payload := connect.NewRequest(&inventory.DeleteInventoryRequest{
		SkuIds: skuIDs,
	})

	_, err := c.client.DeleteInventory(ctx, payload)
	return err
}
