package workflow

import (
	"fmt"
	"product-command-module/internal/application/commands/create_product"
	"product-command-module/internal/application/port"
	"product-command-module/internal/bootstrap/config"
	"product-command-module/internal/infra/temporal/activity"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func CreateProductWorkflow(ctx workflow.Context, cmd create_product.Command, cfg config.TemporalConfig) (*create_product.Result, error) {
	retryPolicy := &temporal.RetryPolicy{
		InitialInterval:    time.Second,
		BackoffCoefficient: 2.0,
		MaximumInterval:    10 * time.Second,
		MaximumAttempts:    3,
	}

	activityCtx := workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		TaskQueue:              cfg.ProductTaskQueue,
		ScheduleToCloseTimeout: cfg.ActivityTimeout,
		RetryPolicy:            retryPolicy,
	})

	rollbackCtx := workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		TaskQueue:              cfg.ProductTaskQueue,
		ScheduleToCloseTimeout: cfg.RollbackTimeout,
		RetryPolicy: &temporal.RetryPolicy{
			MaximumAttempts: 5,
		},
	})

	var productAct *activity.ProductActivity
	var inventoryAct *activity.InventoryActivity
	var compensations []func(workflow.Context)
	executeRollback := func() {
		dCtx, _ := workflow.NewDisconnectedContext(rollbackCtx)
		for i := len(compensations) - 1; i >= 0; i-- {
			compensations[i](dCtx)
		}
	}

	var result create_product.Result
	if err := workflow.ExecuteActivity(activityCtx, productAct.CreateProduct, cmd).Get(ctx, &result); err != nil {
		return nil, fmt.Errorf("create product failed: %w", err)
	}

	compensations = append(compensations, func(rCtx workflow.Context) {
		_ = workflow.ExecuteActivity(rCtx, productAct.RollbackProduct, result.ProductID).Get(rCtx, nil)
	})

	items := make([]port.SkuItem, 0, len(result.SkuItems))
	for _, item := range result.SkuItems {
		items = append(items, port.SkuItem{
			SkuID:    item.SkuID.String(),
			Quantity: item.Quantity,
		})
	}

	params := port.CreateInventoryRequest{
		ShopID: cmd.ShopID.String(),
		Items:  items,
	}

	if err := workflow.ExecuteActivity(activityCtx, inventoryAct.CreateInventories, params).Get(ctx, nil); err != nil {
		executeRollback()
		return nil, fmt.Errorf("create sku failed: %w", err)
	}

	return &result, nil
}
