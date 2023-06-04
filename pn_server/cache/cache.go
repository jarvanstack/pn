package cache

import (
	"time"

	"github.com/karlseguin/ccache/v3"
)

type CacheI interface {
	Get(key string) (string, bool)
	Set(key string, value string, t time.Duration)
}

// =====================================================================================================================

var _ CacheI = (*Cache)(nil)

type Cache struct {
	cache *ccache.Cache[string]
}

func New(size int64) *Cache {
	return &Cache{
		cache: ccache.New(ccache.Configure[string]().GetsPerPromote(1).MaxSize(int64(size)).ItemsToPrune(500)),
	}
}

func (c *Cache) Get(key string) (string, bool) {
	item := c.cache.Get(key)
	if item == nil {
		return "", false
	}
	return item.Value(), true
}

func (c *Cache) Set(key string, value string, t time.Duration) {
	c.cache.Set(key, value, t)
}
