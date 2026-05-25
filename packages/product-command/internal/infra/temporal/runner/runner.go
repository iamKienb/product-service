package runner

import (
	"context"
	"product-command-module/internal/application/commands/create_product"
	"product-command-module/internal/bootstrap/config"

	"go.temporal.io/sdk/client"
)

type Runner interface {
	CreateProduct(ctx context.Context, cmd create_product.Command) (*create_product.Result, error)
}

type workflowRunner struct {
	temporalClient client.Client
	temporalCfg    config.TemporalConfig
}

func NewWorkflowRunner(temporalClient client.Client, cfg config.TemporalConfig) Runner {
	return &workflowRunner{temporalClient: temporalClient, temporalCfg: cfg}
}
