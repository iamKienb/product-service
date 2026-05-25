package temporal

import (
	"product-command-module/internal/infra/temporal/activity"
	"product-command-module/internal/infra/temporal/workflow"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

type Worker struct {
	worker worker.Worker
}

type Registry struct {
	ProductAct   *activity.ProductActivity
	InventoryAct *activity.InventoryActivity
}

func NewWorker(
	temporalClient client.Client,
	taskQueue string,
	registry Registry,
) *Worker {
	newWorker := worker.New(temporalClient, taskQueue, worker.Options{
		MaxConcurrentActivityExecutionSize: 10,
	})

	newWorker.RegisterWorkflow(workflow.CreateProductWorkflow)

	newWorker.RegisterActivity(registry.ProductAct)
	newWorker.RegisterActivity(registry.InventoryAct)

	return &Worker{worker: newWorker}
}

func (w *Worker) Start() error {
	return w.worker.Start()
}
func (w *Worker) Stop() {
	w.worker.Stop()
}
