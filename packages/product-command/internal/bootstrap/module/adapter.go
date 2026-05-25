package module

import (
	"log/slog"
	"net/http"
	"product-command-module/internal/adapter/product"

	"connectrpc.com/connect"
	"connectrpc.com/grpcreflect"
	"github.com/iamKienb/api-contract/gen/product/productconnect"
	authx "github.com/iamKienb/go-core/middleware/auth"
	observabilityx "github.com/iamKienb/go-core/middleware/observability"
)

type AdapterModule struct {
	Mux *http.ServeMux
}

func NewAdapterModule(app *ApplicationModule, logger *slog.Logger) *AdapterModule {
	var interceptors []connect.Interceptor

	tracingInterceptor, err := observabilityx.TracingInterceptor()
	if err != nil {
		logger.Error("failed to initialize tracing interceptor", slog.Any("error", err))
	} else {
		interceptors = append(interceptors, tracingInterceptor)
	}

	interceptors = append(interceptors,
		observabilityx.RecoveryInterceptor(logger),
		authx.RequestContextInterceptor(),
		authx.AuthInternalInterceptor(),
		observabilityx.LoggingInterceptor(logger),
		observabilityx.ValidationRequestInterceptor(),
		observabilityx.ErrorResponseInterceptor(logger),
	)

	allInterceptors := connect.WithInterceptors(interceptors...)

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
