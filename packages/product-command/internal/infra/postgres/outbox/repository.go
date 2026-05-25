package outbox

import (
	"context"
	"product-command-module/db/repository"
	"product-command-module/internal/application/port"

	pgx "github.com/iamKienb/go-core/postgres"
)

type outboxRepository struct {
	queries *repository.Queries
}

func NewRepository(service pgx.PGXService) port.OutboxRepository {
	return &outboxRepository{
		queries: repository.New(service.GetPool()),
	}
}

func (r *outboxRepository) getQuerier(ctx context.Context) *repository.Queries {
	if tx := pgx.ExtractTx(ctx); tx != nil {
		return r.queries.WithTx(tx)
	}

	return r.queries
}
