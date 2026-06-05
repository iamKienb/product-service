package product

import (
	"context"
	"errors"
	"product-command-module/db/repository"
	"product-command-module/internal/domain/product"

	pgx "github.com/iamKienb/go-core/postgres"
	"github.com/jackc/pgx/v5/pgconn"
)

type productRepository struct {
	queries *repository.Queries
}

var productSlugConstraints = map[string]struct{}{
	"products_shop_id_slug_key": {},
	"uq_product_slug":           {},
}

func NewRepository(service pgx.PGXService) product.Repository {
	return &productRepository{
		queries: repository.New(service.GetPool()),
	}
}

func (r *productRepository) getQuerier(ctx context.Context) *repository.Queries {
	if tx := pgx.ExtractTx(ctx); tx != nil {
		return r.queries.WithTx(tx)
	}
	return r.queries
}

func (r *productRepository) IsDuplicateSlug(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		_, ok := productSlugConstraints[pgErr.ConstraintName]
		return pgErr.Code == "23505" && ok
	}

	return false
}
