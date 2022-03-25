package store

// InvalidateOptions 表示缓存失效可用选项
type InvalidateOptions struct {
	// Tags 允许指定与当前值相关联的标签
	Tags []string
}

// TagsValue returns the tags option value
func (o InvalidateOptions) TagsValue() []string {
	return o.Tags
}
