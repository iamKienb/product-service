package module

import (
	"product-worker-module/internal/application/port"
	"product-worker-module/internal/application/processor"
)

type ApplicationModule struct {
	EventProcessor port.EventProcessor
}

func NewApplicationModule(infra *InfraModule) *ApplicationModule {
	return &ApplicationModule{
		EventProcessor: processor.NewProductEventProcessor(infra.ESRepo, infra.workerCache),
	}
}
