package elasticsearch

import (
	"context"
	"product-worker-module/internal/application/port"

	"github.com/elastic/go-elasticsearch/v8"
	esx "github.com/iamKienb/go-core/elasticsearch"
)

type esRepository struct {
	service esx.ESXService
	client  *elasticsearch.TypedClient
}

func NewESRepository(service esx.ESXService, client *elasticsearch.TypedClient) port.ESRepository {
	return &esRepository{
		service: service,
		client:  client,
	}
}
func (r *esRepository) SyncData(ctx context.Context, index string, id string, data any) error {
	return r.service.Sync(ctx, index, id, data)
}

func (r *esRepository) Delete(ctx context.Context, index string, id string) error {
	return r.service.Delete(ctx, index, id)
}
