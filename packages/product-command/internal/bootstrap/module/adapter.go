package module

import (
	"log/slog"
	"net/http"
	"product-command-module/internal/adapter/product"

	"connectrpc.com/grpcreflect"
	"github.com/iamKienb/api-contract/gen/product/productconnect"
	observabilityx "github.com/iamKienb/go-core/middleware/observability"
)

type AdapterModule struct {
	Mux *http.ServeMux
}

func NewAdapterModule(app *ApplicationModule, logger *slog.Logger) *AdapterModule {
	allInterceptors := observabilityx.InternalServerOption(logger)

	mux := http.NewServeMux()
	reflector := grpcreflect.NewStaticReflector(
		productconnect.ProductCommandServiceName,
	)

	productServer := product.NewProductServer(
		app.CreateProductExecutor,
	)

	mux.Handle(productconnect.NewProductCommandServiceHandler(productServer, allInterceptors))

	mux.Handle(grpcreflect.NewHandlerV1(reflector))
	mux.Handle(grpcreflect.NewHandlerV1Alpha(reflector))

	return &AdapterModule{Mux: mux}
}
