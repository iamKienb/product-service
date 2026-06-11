package runner

import (
	"context"
	"fmt"
	"product-command-module/internal/application/commands/create_product"
	"product-command-module/internal/infra/temporal/workflow"

	authx "github.com/iamKienb/go-core/middleware/auth"
	enumspb "go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/client"
)

func (r *workflowRunner) CreateProduct(ctx context.Context, cmd create_product.Command) (*create_product.Result, error) {
	requestID := authx.GetRequestID(ctx)
	run, err := r.temporalClient.ExecuteWorkflow(ctx, client.StartWorkflowOptions{
		ID:                       fmt.Sprintf("create-product-%s-%s-%s", cmd.ShopID.String(), cmd.Slug, requestID),
		TaskQueue:                r.temporalCfg.ProductTaskQueue,
		WorkflowIDConflictPolicy: enumspb.WORKFLOW_ID_CONFLICT_POLICY_USE_EXISTING,
		WorkflowIDReusePolicy:    enumspb.WORKFLOW_ID_REUSE_POLICY_REJECT_DUPLICATE,
	}, workflow.CreateProductWorkflow, cmd, r.temporalCfg)

	if err != nil {
		return nil, err
	}

	var output create_product.Result
	if err := run.Get(ctx, &output); err != nil {
		return nil, fmt.Errorf("saga error: %w", err)
	}

	return &output, nil
}
