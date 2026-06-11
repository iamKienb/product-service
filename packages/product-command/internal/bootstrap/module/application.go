package module

import (
	"product-command-module/internal/application/commands/create_product"
	"product-command-module/internal/application/services/outbox"
	"product-command-module/internal/application/services/product"
)

type ApplicationModule struct {
	ProductService        product.Service
	CreateProductExecutor create_product.Executor
}

func NewApplicationModule(infra *InfraModule) *ApplicationModule {
	outboxService := outbox.NewOutboxService(infra.OutboxRepo)

	productService := product.NewProductService(
		infra.ProductRepo,
		infra.ProductCache,
		outboxService,
		infra.TxManager,
	)

	createProductCommandHandler := create_product.NewHandler(infra.ShopClient, infra.WorkflowRunner)

	return &ApplicationModule{
		ProductService:        productService,
		CreateProductExecutor: createProductCommandHandler,
	}
}
