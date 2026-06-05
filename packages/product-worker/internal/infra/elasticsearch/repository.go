package elasticsearch

import (
	"context"
	"product-worker-module/internal/application/port"

	esx "github.com/iamKienb/go-core/elasticsearch"
)

type esRepository struct{ service esx.ESXService }

func NewESRepository(service esx.ESXService) port.ESRepository {
	return &esRepository{service: service}
}
func (r *esRepository) SyncData(ctx context.Context, index string, id string, data any) error {
	return r.service.Sync(ctx, index, id, data)
}
func (r *esRepository) Delete(ctx context.Context, index string, id string) error {
	return r.service.Delete(ctx, index, id)
}
