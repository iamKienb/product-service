package module

import (
	"context"
	"fmt"
	"product-worker-module/internal/application/port"
	"product-worker-module/internal/bootstrap/config"
	"product-worker-module/internal/infra/elasticsearch"

	esx "github.com/iamKienb/go-core/elasticsearch"
	kafkax "github.com/iamKienb/go-core/kafka"
)

type InfraModule struct {
	ESService esx.ESXService
	Kafka     kafkax.KafkaXService
	ESRepo    port.ESRepository
}

func NewInfraModule(ctx context.Context, cfg *config.ProductWorkerConfig) (*InfraModule, error) {
	esService, err := esx.New(cfg.ES)
	if err != nil {
		return nil, fmt.Errorf("elasticsearch: %w", err)
	}
	kafka, err := kafkax.New(cfg.Kafka)
	if err != nil {
		return nil, fmt.Errorf("kafka: %w", err)
	}
	return &InfraModule{ESService: esService, Kafka: kafka, ESRepo: elasticsearch.NewESRepository(esService)}, nil
}
