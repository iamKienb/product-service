package create_product

import (
	"context"
	"product-command-module/internal/application/port"

	"github.com/iamKienb/go-core/app_error"
)

type workflowRunner interface {
	CreateProduct(ctx context.Context, cmd Command) (*Result, error)
}

type handler struct {
	client   port.ShopClient
	workflow workflowRunner
}

func NewHandler(shopClient port.ShopClient, workflow workflowRunner) Executor {
	return &handler{
		client:   shopClient,
		workflow: workflow,
	}
}

func (h *handler) Execute(ctx context.Context, cmd Command) (*Result, error) {
	result, err := h.client.CheckPermission(ctx, port.CheckPermissionRequest{
		ShopID: cmd.ShopID.String(),
		UserID: cmd.UserID.String(),
		Action: cmd.Action,
	})
	if err != nil {
		return nil, err
	}

	if !result.IsAllowed {
		return nil, app_error.Forbidden(result.Message)
	}

	return h.workflow.CreateProduct(ctx, cmd)
}
