package module

import (
	"product-query-module/internal/application/queries/get_product_by_sku_ids"
	productservice "product-query-module/internal/application/service"
	"product-shared-module/alias"
)

type ApplicationModule struct {
	GetProductSkuExecutor get_product_by_sku_ids.Executor
}

func NewApplicationModule(infra *InfraModule) *ApplicationModule {
	productService := productservice.NewProductService(infra.ESService.GetClient(), alias.ProductAlias)

	return &ApplicationModule{
		GetProductSkuExecutor: get_product_by_sku_ids.NewHandler(productService),
	}
}
