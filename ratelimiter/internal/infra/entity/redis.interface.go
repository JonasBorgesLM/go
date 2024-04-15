package entity

import (
	"context"
	"time"
)

type RedisInterface interface {
	AddBlock(ctx context.Context, key string, blockDurationMilliseconds int64) (*time.Time, error)
	GetBlock(ctx context.Context, key string) (*time.Time, error)
	IncrementAccessCount(ctx context.Context, key string, maxAccessCount int64) (bool, error)
}
