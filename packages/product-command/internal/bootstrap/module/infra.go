package module

import (
	"context"
	"fmt"
	"net/http"
	"product-command-module/internal/application/port"
	"product-command-module/internal/bootstrap/config"
	"product-command-module/internal/domain/product"
	"product-command-module/internal/infra/grpc_client"
	outboxPg "product-command-module/internal/infra/postgres/outbox"
	"product-command-module/internal/infra/temporal/runner"
	"time"

	pgx "github.com/iamKienb/go-core/postgres"
	redisx "github.com/iamKienb/go-core/redis"
	"go.temporal.io/sdk/client"
)

type InfraModule struct {
	TemporalClient client.Client
	WorkflowRunner runner.Runner
	PGService      pgx.PGXService
	RedisService   redisx.RedisXService

	ProductRepo  product.Repository
	ProductCache port.ProductCache
	OutboxRepo   port.OutboxRepository

	InventoryClient port.InventoryClient
	ShopClient      port.ShopClient

	TxManager port.TxManager
}

func NewInfraModule(ctx context.Context, cfg *config.ProductCommandConfig) (*InfraModule, error) {
	tClient, err := client.Dial(client.Options{
		HostPort:  cfg.Temporal.Address,
		Namespace: cfg.Temporal.Namespace,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to dial temporal: %w", err)
	}

	httpClient := &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			IdleConnTimeout:     90 * time.Second,
			MaxIdleConnsPerHost: 20,
		},
	}

	apiGatewayAddr := cfg.Gateway.APIGatewayAddr

	pgService, err := pgx.New(cfg.Postgres)
	if err != nil {
		return nil, fmt.Errorf("postgres: %w", err)
	}

	redisService, err := redisx.New(cfg.Redis)
	if err != nil {
		return nil, fmt.Errorf("redis: %w", err)
	}

	workflowRunner := runner.NewWorkflowRunner(tClient, cfg.Temporal)

	return &InfraModule{
		PGService:    pgService,
		RedisService: redisService,

		TemporalClient: tClient,
		WorkflowRunner: workflowRunner,

		OutboxRepo: outboxPg.NewRepository(pgService),

		ShopClient: grpc_client.NewShopClient(httpClient, apiGatewayAddr),

		InventoryClient: grpc_client.NewInventoryClient(httpClient, apiGatewayAddr),

		TxManager: pgx.NewTxManager(pgService.GetPool()),
	}, nil
}
