package processor

import (
	"context"
	"fmt"
	"product-shared-module/alias"
	"product-shared-module/events"
	"product-worker-module/internal/application/port"
	"product-worker-module/internal/application/processor/handler"
	"time"
)

const (
	idemKeyTTL = 24 * time.Hour
	key        = "user-worker:key:%s"
)

type ProductEventProcessor struct {
	handlers    map[string]port.EventHandler
	workerCache port.WorkerCache
}

func NewProductEventProcessor(repo port.ESRepository, workerCache port.WorkerCache) port.EventProcessor {
	p := &ProductEventProcessor{
		handlers:    make(map[string]port.EventHandler),
		workerCache: workerCache,
	}

	p.handlers[events.TopicProductCreated] = handler.NewProductCreatedHandler(repo, alias.ProductAlias)
	p.handlers[events.TopicProductDeleted] = handler.NewProductDeletedHandler(repo, alias.ProductAlias)

	return p
}

func (p *ProductEventProcessor) Handle(ctx context.Context, msg port.Message) error {
	h, ok := p.handlers[msg.Topic]
	if !ok {
		return nil
	}

	idemKey := msg.IdempotencyKey()

	if idemKey != "" {
		key := fmt.Sprintf(key, idemKey)
		isNew, err := p.workerCache.SetNx(ctx, key, 1, idemKeyTTL)
		if err != nil {
			return err
		}

		if !isNew {
			return nil
		}
	}

	return h.Handle(ctx, msg.Value)
}
