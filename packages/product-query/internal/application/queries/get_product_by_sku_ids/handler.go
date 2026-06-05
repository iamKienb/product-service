package get_product_by_sku_ids

import "context"

type productService interface {
	SearchBySkuIDs(ctx context.Context, qry Query) ([]*Result, error)
}

type handler struct {
	service productService
}

func NewHandler(service productService) Executor {
	return &handler{service: service}
}

func (h *handler) Execute(ctx context.Context, qry Query) ([]*Result, error) {
	return h.service.SearchBySkuIDs(ctx, qry)
}
