package port

import "context"

type ESRepository interface {
	SyncData(ctx context.Context, index string, id string, data any) error
	Delete(ctx context.Context, index string, id string) error
}
