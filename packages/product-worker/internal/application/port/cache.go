package port

import (
	"context"
	"time"
)

type WorkerCache interface {
	SetNx(ctx context.Context, key string, data any, ttl time.Duration) (bool, error)
}
