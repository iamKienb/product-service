package module

import (
	"context"
	"fmt"
	"product-worker-module/internal/application/port"
	"product-worker-module/internal/bootstrap/config"
	"product-worker-module/internal/infra/cache"
	"product-worker-module/internal/infra/elasticsearch"

	esx "github.com/iamKienb/go-core/elasticsearch"
	kafkax "github.com/iamKienb/go-core/kafka"
	redisx "github.com/iamKienb/go-core/redis"
)

type InfraModule struct {
	ESService    esx.ESXService
	RedisService redisx.RedisXService
	Kafka        kafkax.KafkaXService
	ESRepo       port.ESRepository
	workerCache  port.WorkerCache
}

func NewInfraModule(ctx context.Context, cfg *config.ProductWorkerConfig) (*InfraModule, error) {
	esService, err := esx.New(cfg.ES)
	if err != nil {
		return nil, fmt.Errorf("elasticsearch: %w", err)
	}

	redisService, err := redisx.New(cfg.Redis)
	if err != nil {
		return nil, fmt.Errorf("redis: %w", err)
	}

	kafka, err := kafkax.New(cfg.Kafka)
	if err != nil {
		return nil, fmt.Errorf("kafka: %w", err)
	}

	return &InfraModule{
		ESService:    esService,
		RedisService: redisService,
		Kafka:        kafka,
		ESRepo:       elasticsearch.NewESRepository(esService, esService.GetClient()),
		workerCache:  cache.NewWorkerCache(redisService.GetClient()),
	}, nil
}
