package store

import (
	"context"
	"time"
)

// StoreInterface 是所有可用存储的接口
type StoreInterface interface {
	Get(ctx context.Context, key interface{}) (interface{}, error)
	GetWithTTL(ctx context.Context, key interface{}) (interface{}, time.Duration, error)
	Set(ctx context.Context, key interface{}, value interface{}, options *Options) error
	Delete(ctx context.Context, key interface{}) error
	Invalidate(ctx context.Context, options InvalidateOptions) error
	Clear(ctx context.Context) error
	GetType() string // 获取 Store 实现的名字
}
