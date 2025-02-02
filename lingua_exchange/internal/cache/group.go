package cache

import (
	"context"
	"strings"
	"time"

	"github.com/zhufuyi/sponge/pkg/cache"
	"github.com/zhufuyi/sponge/pkg/encoding"
	"github.com/zhufuyi/sponge/pkg/utils"

	"lingua_exchange/internal/model"
)

const (
	// cache prefix key, must end with a colon
	groupCachePrefixKey = "group:"
	// GroupExpireTime expire time
	GroupExpireTime = 5 * time.Minute
)

var _ GroupCache = (*groupCache)(nil)

// GroupCache cache interface
type GroupCache interface {
	Set(ctx context.Context, id uint64, data *model.Group, duration time.Duration) error
	Get(ctx context.Context, id uint64) (*model.Group, error)
	MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.Group, error)
	MultiSet(ctx context.Context, data []*model.Group, duration time.Duration) error
	Del(ctx context.Context, id uint64) error
	SetCacheWithNotFound(ctx context.Context, id uint64) error
}

// groupCache define a cache struct
type groupCache struct {
	cache cache.Cache
}

// NewGroupCache new a cache
func NewGroupCache(cacheType *model.CacheType) GroupCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.Group{}
		})
		return &groupCache{cache: c}
	case "memory":
		c := cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.Group{}
		})
		return &groupCache{cache: c}
	}

	return nil // no cache
}

// GetGroupCacheKey cache key
func (c *groupCache) GetGroupCacheKey(id uint64) string {
	return groupCachePrefixKey + utils.Uint64ToStr(id)
}

// Set write to cache
func (c *groupCache) Set(ctx context.Context, id uint64, data *model.Group, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetGroupCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get cache value
func (c *groupCache) Get(ctx context.Context, id uint64) (*model.Group, error) {
	var data *model.Group
	cacheKey := c.GetGroupCacheKey(id)
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MultiSet multiple set cache
func (c *groupCache) MultiSet(ctx context.Context, data []*model.Group, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetGroupCacheKey(uint64(v.ID))
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}

	return nil
}

// MultiGet multiple get cache, return key in map is id value
func (c *groupCache) MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.Group, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetGroupCacheKey(v)
		keys = append(keys, cacheKey)
	}

	itemMap := make(map[string]*model.Group)
	err := c.cache.MultiGet(ctx, keys, itemMap)
	if err != nil {
		return nil, err
	}

	retMap := make(map[uint64]*model.Group)
	for _, id := range ids {
		val, ok := itemMap[c.GetGroupCacheKey(id)]
		if ok {
			retMap[id] = val
		}
	}

	return retMap, nil
}

// Del delete cache
func (c *groupCache) Del(ctx context.Context, id uint64) error {
	cacheKey := c.GetGroupCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetCacheWithNotFound set empty cache
func (c *groupCache) SetCacheWithNotFound(ctx context.Context, id uint64) error {
	cacheKey := c.GetGroupCacheKey(id)
	err := c.cache.SetCacheWithNotFound(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}
