package product

import (
	"context"
	"product-query-module/internal/application/queries/get_product_by_sku_ids"

	"connectrpc.com/connect"
	"github.com/iamKienb/api-contract/gen/product"
	"github.com/iamKienb/api-contract/gen/product/productconnect"
)

type productServer struct {
	productconnect.UnimplementedProductQueryServiceHandler
	getProductSkuExecutor get_product_by_sku_ids.Executor
}

func NewProductServer(
	getProductSkuExecutor get_product_by_sku_ids.Executor,
) *productServer {
	return &productServer{
		getProductSkuExecutor: getProductSkuExecutor,
	}
}

func (s *productServer) GetProductsBySkuIDs(ctx context.Context, req *connect.Request[product.GetProductsBySkuIDsRequest]) (*connect.Response[product.GetProductsBySkuIDsResponse], error) {
	qry, err := toGetProductBySkuIDsQuery(req.Msg)
	if err != nil {
		return nil, err
	}

	result, err := s.getProductSkuExecutor.Execute(ctx, qry)
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(toGetProductBySkuIDsResponse(result)), nil
}

var _ productconnect.ProductQueryServiceHandler = (*productServer)(nil)
