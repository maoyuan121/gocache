[![Test](https://github.com/eko/gocache/actions/workflows/all.yml/badge.svg?branch=master)](https://github.com/eko/gocache/actions/workflows/all.yml)
[![GoDoc](https://godoc.org/github.com/eko/gocache?status.png)](https://godoc.org/github.com/eko/gocache)
[![GoReportCard](https://goreportcard.com/badge/github.com/eko/gocache)](https://goreportcard.com/report/github.com/eko/gocache)
[![codecov](https://codecov.io/gh/eko/gocache/branch/master/graph/badge.svg)](https://codecov.io/gh/eko/gocache)

Gocache
=======

Gocache 是一个 Go 缓存库
这是一个可扩展的缓存库，为缓存数据带来了许多特性。

## 概要

以下是它带来的细节:

* ✅ 多缓存存储:实际上在内存，redis，或你自己的自定义存储
* ✅ 链式缓存:使用多个缓存，有优先级顺序(例如内存，然后回退到 redis 共享缓存)
* ✅ 一个可加载的缓存:允许你调用一个回调函数把你的数据放回缓存
* ✅ 一个度量缓存，让你存储关于你的缓存使用的度量(命中，错过，设置成功，设置错误，…)
* ✅ 将缓存值作为结构自动 marshal/unmarshal 的marshaler
* ✅ 在存储中定义默认值，并在设置数据时覆盖它们
* ✅ 缓存过期时间和/或使用标记无效

## Built-in stores

* [Memory (bigcache)](https://github.com/allegro/bigcache) (allegro/bigcache)
* [Memory (ristretto)](https://github.com/dgraph-io/ristretto) (dgraph-io/ristretto)
* [Memory (go-cache)](https://github.com/patrickmn/go-cache) (patrickmn/go-cache)
* [Memcache](https://github.com/bradfitz/gomemcache) (bradfitz/memcache)
* [Redis](https://github.com/go-redis/redis/v8) (go-redis/redis)
* [Freecache](https://github.com/coocood/freecache) (coocood/freecache)
* [Pegasus](https://pegasus.apache.org/) ([apache/incubator-pegasus](https://github.com/apache/incubator-pegasus)) [benchmark](https://pegasus.apache.org/overview/benchmark/)
* More to come soon

## Built-in metrics providers

* [Prometheus](https://github.com/prometheus/client_golang)

## Available cache features in detail

### A simple cache

下面是一个简单的 Redis 缓存实例化，但你也可以看看其他可用的存储:

#### Memcache

```go
memcacheStore := store.NewMemcache(
	memcache.New("10.0.0.1:11211", "10.0.0.2:11211", "10.0.0.3:11212"),
	&store.Options{
		Expiration: 10*time.Second,
	},
)

cacheManager := cache.New(memcacheStore)
err := cacheManager.Set(ctx, "my-key", []byte("my-value"), &store.Options{
	Expiration: 15*time.Second, //  覆盖在  store 中定义的默认值 10 秒
})
if err != nil {
    panic(err)
}

value := cacheManager.Get(ctx, "my-key")

cacheManager.Delete(ctx, "my-key")

cacheManager.Clear(ctx) // 清除整个缓存，以防您想要清除所有缓存
```

#### Memory (using Bigcache)

```go
bigcacheClient, _ := bigcache.NewBigCache(bigcache.DefaultConfig(5 * time.Minute))
bigcacheStore := store.NewBigcache(bigcacheClient, nil) // No options provided (as second argument)

cacheManager := cache.New(bigcacheStore)
err := cacheManager.Set(ctx, "my-key", []byte("my-value"), nil)
if err != nil {
    panic(err)
}

value := cacheManager.Get(ctx, "my-key")
```

#### Memory (using Ristretto)

```go
ristrettoCache, err := ristretto.NewCache(&ristretto.Config{
	NumCounters: 1000,
	MaxCost: 100,
	BufferItems: 64,
})
if err != nil {
    panic(err)
}
ristrettoStore := store.NewRistretto(ristrettoCache, nil)

cacheManager := cache.New(ristrettoStore)
err := cacheManager.Set(ctx, "my-key", "my-value", &store.Options{Cost: 2})
if err != nil {
    panic(err)
}

value := cacheManager.Get(ctx, "my-key")

cacheManager.Delete(ctx, "my-key")
```

#### Memory (using Go-cache)

```go
gocacheClient := gocache.New(5*time.Minute, 10*time.Minute)
gocacheStore := store.NewGoCache(gocacheClient, nil)

cacheManager := cache.New(gocacheStore)
err := cacheManager.Set(ctx, "my-key", []byte("my-value"), nil)
if err != nil {
	panic(err)
}

value, err := cacheManager.Get(ctx, "my-key")
if err != nil {
	panic(err)
}
fmt.Printf("%s", value)
```

#### Redis

```go
redisStore := store.NewRedis(redis.NewClient(&redis.Options{
	Addr: "127.0.0.1:6379",
}), nil)

cacheManager := cache.New(redisStore)
err := cacheManager.Set("my-key", "my-value", &store.Options{Expiration: 15*time.Second})
if err != nil {
    panic(err)
}

value, err := cacheManager.Get(ctx, "my-key")
switch err {
	case nil:
		fmt.Printf("Get the key '%s' from the redis cache. Result: %s", "my-key", value)
	case redis.Nil:
		fmt.Printf("Failed to find the key '%s' from the redis cache.", "my-key")
	default:
	    fmt.Printf("Failed to get the value from the redis cache with key '%s': %v", "my-key", err)
}
```

#### Freecache

```go
freecacheStore := store.NewFreecache(freecache.NewCache(1000), &Options{
	Expiration: 10 * time.Second,
})

cacheManager := cache.New(freecacheStore)
err := cacheManager.Set(ctx, "by-key", []byte("my-value"), opts)
if err != nil {
    panic(err)
}

value := cacheManager.Get(ctx, "my-key")
```

#### Pegasus

```go
pegasusStore, err := store.NewPegasus(&store.OptionsPegasus{
    MetaServers: []string{"127.0.0.1:34601", "127.0.0.1:34602", "127.0.0.1:34603"},
})

if err != nil {
    fmt.Println(err)
    return
}

cacheManager := cache.New(pegasusStore)
err = cacheManager.Set(ctx, "my-key", "my-value", &store.Options{
    Expiration: 10 * time.Second,
})
if err != nil {
    panic(err)
}

value, _ := cacheManager.Get(ctx, "my-key")
```

### 链式缓存

在这里，我们将按照以下顺序链接缓存:首先在内存中使用 Ristretto 存储，然后在 Redis (作为一个后备):

```go
// Initialize Ristretto cache and Redis client
ristrettoCache, err := ristretto.NewCache(&ristretto.Config{NumCounters: 1000, MaxCost: 100, BufferItems: 64})
if err != nil {
    panic(err)
}

redisClient := redis.NewClient(&redis.Options{Addr: "127.0.0.1:6379"})

// Initialize stores
ristrettoStore := store.NewRistretto(ristrettoCache, nil)
redisStore := store.NewRedis(redisClient, &store.Options{Expiration: 5*time.Second})

// Initialize chained cache
cacheManager := cache.NewChain(
    cache.New(ristrettoStore),
    cache.New(redisStore),
)

// ... Then, do what you want with your cache
```

`Chain` 缓存也会把数据放回之前的缓存中，所以在这种情况下，如果 ristretto 的缓存中没有数据，但是 redis 有，数据也会被放回 ristretto(内存) 缓存中。

### A loadable cache

这个缓存将提供一个加载函数，作为一个可调用函数，并将你的数据设置回缓存，以防它们不可用:

```go
// Initialize Redis client and store
redisClient := redis.NewClient(&redis.Options{Addr: "127.0.0.1:6379"})
redisStore := store.NewRedis(redisClient, nil)

// Initialize a load function that loads your data from a custom source
loadFunction := func(ctx context.Context, key interface{}) (interface{}, error) {
    // ... retrieve value from available source
    return &Book{ID: 1, Name: "My test amazing book", Slug: "my-test-amazing-book"}, nil
}

// Initialize loadable cache
cacheManager := cache.NewLoadable(
	loadFunction,
	cache.New(redisStore),
)

// ... Then, you can get your data and your function will automatically put them in cache(s)
```

当然，你也可以传递一个 `Chain` 缓存到 `Loadable` 缓存中，所以如果你的数据在所有缓存中不可用，它会把它带回所有缓存中。

### A metric cache to retrieve cache statistics

这个缓存将根据您传递给它的度量提供程序记录度量。在这里，我们使用普罗米修斯的供应商:

```go
// Initialize Redis client and store
redisClient := redis.NewClient(&redis.Options{Addr: "127.0.0.1:6379"})
redisStore := store.NewRedis(redisClient, nil)

// Initializes Prometheus metrics service
promMetrics := metrics.NewPrometheus("my-test-app")

// Initialize metric cache
cacheManager := cache.NewMetric(
	promMetrics,
	cache.New(redisStore),
)

// ... Then, you can get your data and metrics will be observed by Prometheus
```

### A marshaler wrapper

一些缓存，如 Redis 存储和返回值作为一个字符串，所以你必须 marshal/unmarshal 你的结构，如果你想缓存一个对象。
这就是为什么我们带来了一个封送服务，包装你的缓存，让你的工作:

```go
// Initialize Redis client and store
redisClient := redis.NewClient(&redis.Options{Addr: "127.0.0.1:6379"})
redisStore := store.NewRedis(redisClient, nil)

// Initialize chained cache
cacheManager := cache.NewMetric(
	promMetrics,
	cache.New(redisStore),
)

// Initializes marshaler
marshal := marshaler.New(cacheManager)

key := BookQuery{Slug: "my-test-amazing-book"}
value := Book{ID: 1, Name: "My test amazing book", Slug: "my-test-amazing-book"}

err = marshal.Set(ctx, key, value)
if err != nil {
    panic(err)
}

returnedValue, err := marshal.Get(ctx, key, new(Book))
if err != nil {
    panic(err)
}

// Then, do what you want with the  value

marshal.Delete(ctx, "my-key")
```

你需要做的唯一一件事就是在调用' `.Get()` 方法时指定你想要你的值被反编组的结构作为第二个参数。

### Cache invalidation using tags

您可以将一些标记附加到您创建的 item 上，以便稍后可以轻松地使其中一些标记失效。

标签使用您为缓存选择的相同存储存储。

下面是使用的例子：

```go
// Initialize Redis client and store
redisClient := redis.NewClient(&redis.Options{Addr: "127.0.0.1:6379"})
redisStore := store.NewRedis(redisClient, nil)

// Initialize chained cache
cacheManager := cache.NewMetric(
	promMetrics,
	cache.New(redisStore),
)

// Initializes marshaler
marshal := marshaler.New(cacheManager)

key := BookQuery{Slug: "my-test-amazing-book"}
value := Book{ID: 1, Name: "My test amazing book", Slug: "my-test-amazing-book"}

// 在缓存中设置了一个 item，并附上了一个 "book" 标签
err = marshal.Set(ctx, key, value, store.Options{Tags: []string{"book"}})
if err != nil {
    panic(err)
}

// 清除所有附上了 "book" 标签的 item
err := marshal.Invalidate(ctx, store.InvalidateOptions{Tags: []string{"book"}})
if err != nil {
    panic(err)
}

returnedValue, err := marshal.Get(ctx, key, new(Book))
if err != nil {
	// Should be triggered because item has been deleted so it cannot be found.
    panic(err)
}
```

将其与缓存上的过期时间混合在一起，可以更好地控制数据的缓存方式。

### Write your own custom cache

Cache respect the following interface so you can write your own (proprietary?) cache logic if needed by implementing the following interface:

```go
type CacheInterface interface {
	Get(ctx context.Context, key interface{}) (interface{}, error)
	Set(ctx context.Context, key, object interface{}, options *store.Options) error
	Delete(ctx context.Context, key interface{}) error
	Invalidate(ctx context.Context, options store.InvalidateOptions) error
	Clear(ctx context.Context) error
	GetType() string
}
```

Or, in case you use a setter cache, also implement the `GetCodec()` method:

```go
type SetterCacheInterface interface {
	CacheInterface
	GetWithTTL(ctx context.Context, key interface{}) (interface{}, time.Duration, error)

	GetCodec() codec.CodecInterface
}
```

As all caches available in this library implement `CacheInterface`, you will be able to mix your own caches with your own.

### Write your own custom store

You also have the ability to write your own custom store by implementing the following interface:

```go
type StoreInterface interface {
	Get(ctx context.Context, key interface{}) (interface{}, error)
	GetWithTTL(ctx context.Context, key interface{}) (interface{}, time.Duration, error)
	Set(ctx context.Context, key interface{}, value interface{}, options *Options) error
	Delete(ctx context.Context, key interface{}) error
	Invalidate(ctx context.Context, options InvalidateOptions) error
	Clear(ctx context.Context) error
	GetType() string
}
```

Of course, I suggest you to have a look at current caches or stores to implement your own.

### Custom cache key generator

You can implement the following interface in order to generate a custom cache key:

```go
type CacheKeyGenerator interface {
	GetCacheKey() string
}
```

### Benchmarks

![Benchmarks](https://raw.githubusercontent.com/eko/gocache/master/misc/benchmarks.jpeg)

## Community

Please feel free to contribute on this library and do not hesitate to open an issue if you want to discuss about a feature.

## Run tests

Generate mocks:
```bash
$ go get github.com/golang/mock/mockgen
$ make mocks
```

Test suite can be run with:

```bash
$ make test # run unit test
$ make benchmark-store # run benchmark test
```
