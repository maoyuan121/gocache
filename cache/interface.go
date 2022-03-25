package cache

import (
	"context"
	"time"

	"github.com/eko/gocache/v2/codec"
	"github.com/eko/gocache/v2/store"
)

// CacheInterface 表示所有缓存的接口 (aggregate, metric, memory, redis，…)
type CacheInterface interface {
	Get(ctx context.Context, key interface{}) (interface{}, error)
	Set(ctx context.Context, key, object interface{}, options *store.Options) error
	Delete(ctx context.Context, key interface{}) error
	Invalidate(ctx context.Context, options store.InvalidateOptions) error
	Clear(ctx context.Context) error
	GetType() string // 获取 cache 实现名
}

// SetterCacheInterface 表示允许存储的缓存接口 (例如:memory, redis，…)
type SetterCacheInterface interface {
	CacheInterface
	GetWithTTL(ctx context.Context, key interface{}) (interface{}, time.Duration, error)

	GetCodec() codec.CodecInterface
}
