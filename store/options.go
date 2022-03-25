package store

import (
	"time"
)

// Options 是所有 store 的选项
type Options struct {
	// Cost 指的是设置值时该项所使用的内存容量
	// 实际上，它似乎只有 Ristretto 库使用
	Cost int64

	// Expiration 允许在设置值时指定过期时间
	Expiration time.Duration

	// Tags 允许指定与当前值相关联的标签
	Tags []string
}

// CostValue returns the allocated memory capacity
func (o Options) CostValue() int64 {
	return o.Cost
}

// ExpirationValue returns the expiration option value
func (o Options) ExpirationValue() time.Duration {
	return o.Expiration
}

// TagsValue returns the tags option value
func (o Options) TagsValue() []string {
	return o.Tags
}
