package processor

import (
	"context"
	"encoding/json"
	"product-shared-module/alias"
	"product-shared-module/events"
	"product-worker-module/internal/application/port"
	"product-worker-module/internal/application/processor/handler"
)

type ProductEventProcessor struct{ handlers map[string]port.EventHandler }

func NewProductEventProcessor(repo port.ESRepository) port.EventProcessor {
	p := &ProductEventProcessor{handlers: make(map[string]port.EventHandler)}

	p.handlers[events.TopicProductCreated] = handler.NewProductCreatedHandler(repo, alias.ProductAlias)
	p.handlers[events.TopicProductDeleted] = handler.NewProductDeletedHandler(repo, alias.ProductAlias)

	return p
}

func (p *ProductEventProcessor) Handle(ctx context.Context, msg port.Message) error {
	h, ok := p.handlers[msg.Topic]
	if !ok {
		return nil
	}
	var raw json.RawMessage = msg.Value
	return h.Handle(ctx, raw)
}
