package product

import (
	"context"
	"product-command-module/internal/application/commands/create_product"

	"connectrpc.com/connect"
	"github.com/iamKienb/api-contract/gen/product"
	"github.com/iamKienb/api-contract/gen/product/productconnect"
	authx "github.com/iamKienb/go-core/middleware/auth"
)

type productServer struct {
	createProductExecutor create_product.Executor
}

func NewProductServer(
	createProductExecutor create_product.Executor,
) *productServer {
	return &productServer{
		createProductExecutor: createProductExecutor,
	}
}

func (s *productServer) CreateProduct(ctx context.Context, req *connect.Request[product.CreateProductsRequest]) (*connect.Response[product.CreateProductsResponse], error) {
	currentUser := authx.GetUserInfoFromCtx(ctx)
	cmd, err := ToCreateProductCommand(currentUser.UserID, req.Msg)
	if err != nil {
		return nil, err
	}

	result, err := s.createProductExecutor.Execute(ctx, cmd)
	if err != nil {
		return nil, mapError(err)
	}

	return connect.NewResponse(ToCreateProductResponse(result)), nil
}

var _ productconnect.ProductCommandServiceHandler = (*productServer)(nil)
