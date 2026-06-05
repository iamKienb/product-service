package bootstrap

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"product-query-module/internal/bootstrap/config"
	"product-query-module/internal/bootstrap/module"
	"strconv"
	"time"

	configx "github.com/iamKienb/go-core/config"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type App struct {
	logger *slog.Logger
	server *http.Server
	infra  *module.InfraModule
}

func NewApp(logger *slog.Logger) *App { return &App{logger: logger} }
func (a *App) Start(ctx context.Context) error {
	cfg, err := configx.Loader[config.ProductQueryConfig]()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}
	infra, err := module.NewInfraModule(cfg)
	if err != nil {
		return fmt.Errorf("infra: %w", err)
	}
	a.infra = infra
	application := module.NewApplicationModule(infra)
	adapter := module.NewAdapterModule(application, a.logger)
	a.server = &http.Server{Addr: ":" + strconv.Itoa(cfg.Server.GrpcPort), Handler: h2c.NewHandler(adapter.Mux, &http2.Server{}), ReadTimeout: 10 * time.Second, WriteTimeout: 30 * time.Second, IdleTimeout: 60 * time.Second}
	if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("server: %w", err)
	}
	return nil
}
func (a *App) Stop(ctx context.Context) error {
	if a.server != nil {
		_ = a.server.Shutdown(ctx)
	}
	if a.infra != nil && a.infra.ESService != nil {
		return a.infra.ESService.Close(ctx)
	}
	return nil
}
