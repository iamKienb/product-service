package grpc_client

import (
	"context"
	"net/http"
	"product-command-module/internal/application/port"

	"connectrpc.com/connect"
	"github.com/iamKienb/api-contract/gen/shop"
	"github.com/iamKienb/api-contract/gen/shop/shopconnect"
)

type shopClient struct {
	client shopconnect.ShopCommandServiceClient
}

func NewShopClient(httpClient *http.Client, baseURL string) port.ShopClient {
	return &shopClient{
		client: shopconnect.NewShopCommandServiceClient(httpClient, baseURL),
	}
}

func (c *shopClient) CheckPermission(ctx context.Context, req port.CheckPermissionRequest) (port.CheckPermissionResponse, error) {
	result, err := c.client.CheckPermission(ctx, connect.NewRequest(&shop.CheckPermissionRequest{
		ShopId: req.ShopID,
		UserId: req.UserID,
		Action: req.Action,
	}))

	if err != nil {
		return port.CheckPermissionResponse{}, err
	}

	return port.CheckPermissionResponse{
		IsAllowed: result.Msg.GetIsAllowed(),
		Message:   result.Msg.GetMessage(),
	}, nil
}
