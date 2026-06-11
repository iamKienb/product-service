package bootstrap

import (
	"context"
	"fmt"
	"log/slog"
	"product-shared-module/events"
	"product-worker-module/internal/bootstrap/config"
	"product-worker-module/internal/bootstrap/module"
	"sync"

	configx "github.com/iamKienb/go-core/config"
	kafkax "github.com/iamKienb/go-core/kafka"
)

type App struct {
	logger    *slog.Logger
	infra     *module.InfraModule
	consumers []*kafkax.Consumer
}

func NewApp(logger *slog.Logger) *App { return &App{logger: logger} }

func (a *App) Start(ctx context.Context) error {
	cfg, err := configx.Loader[config.ProductWorkerConfig]()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}
	infra, err := module.NewInfraModule(ctx, cfg)
	if err != nil {
		return err
	}
	a.infra = infra
	cfg.Consumer.Topics = events.Topics
	application := module.NewApplicationModule(infra)
	consumers := make([]*kafkax.Consumer, 0, len(cfg.Consumer.Topics))
	for _, topic := range cfg.Consumer.Topics {
		consumer, err := kafkax.NewConsumer(infra.Kafka, cfg.Consumer, a.logger, application.EventProcessor)
		if err != nil {
			return fmt.Errorf("create consumer for topic %s: %w", topic, err)
		}
		consumers = append(consumers, consumer)
	}

	a.consumers = consumers
	a.logger.Info("starting product worker", slog.Any("topics", cfg.Consumer.Topics))

	errCh := make(chan error, len(consumers))
	var wg sync.WaitGroup
	for _, consumer := range consumers {
		wg.Add(1)
		go func(consumer *kafkax.Consumer) {
			defer wg.Done()
			errCh <- consumer.Start(ctx)
		}(consumer)
	}
	go func() { wg.Wait(); close(errCh) }()
	for {
		select {
		case <-ctx.Done():
			return nil
		case err, ok := <-errCh:
			if !ok {
				return nil
			}
			if err != nil {
				return fmt.Errorf("consumer loop: %w", err)
			}
		}
	}
}

func (a *App) Stop(ctx context.Context) error {
	if a.infra != nil {
		for _, consumer := range a.consumers {
			if consumer != nil {
				_ = consumer.Close()
			}
		}
		if a.infra.Kafka != nil {
			_ = a.infra.Kafka.Close()
		}
		if a.infra.ESService != nil {
			_ = a.infra.ESService.Close(ctx)
		}
	}
	return nil
}
