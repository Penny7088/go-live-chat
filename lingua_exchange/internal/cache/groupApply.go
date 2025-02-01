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
	groupApplyCachePrefixKey = "groupApply:"
	// GroupApplyExpireTime expire time
	GroupApplyExpireTime = 5 * time.Minute
)

var _ GroupApplyCache = (*groupApplyCache)(nil)

// GroupApplyCache cache interface
type GroupApplyCache interface {
	Set(ctx context.Context, id uint64, data *model.GroupApply, duration time.Duration) error
	Get(ctx context.Context, id uint64) (*model.GroupApply, error)
	MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.GroupApply, error)
	MultiSet(ctx context.Context, data []*model.GroupApply, duration time.Duration) error
	Del(ctx context.Context, id uint64) error
	SetCacheWithNotFound(ctx context.Context, id uint64) error
}

// groupApplyCache define a cache struct
type groupApplyCache struct {
	cache cache.Cache
}

// NewGroupApplyCache new a cache
func NewGroupApplyCache(cacheType *model.CacheType) GroupApplyCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.GroupApply{}
		})
		return &groupApplyCache{cache: c}
	case "memory":
		c := cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.GroupApply{}
		})
		return &groupApplyCache{cache: c}
	}

	return nil // no cache
}

// GetGroupApplyCacheKey cache key
func (c *groupApplyCache) GetGroupApplyCacheKey(id uint64) string {
	return groupApplyCachePrefixKey + utils.Uint64ToStr(id)
}

// Set write to cache
func (c *groupApplyCache) Set(ctx context.Context, id uint64, data *model.GroupApply, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetGroupApplyCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get cache value
func (c *groupApplyCache) Get(ctx context.Context, id uint64) (*model.GroupApply, error) {
	var data *model.GroupApply
	cacheKey := c.GetGroupApplyCacheKey(id)
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MultiSet multiple set cache
func (c *groupApplyCache) MultiSet(ctx context.Context, data []*model.GroupApply, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetGroupApplyCacheKey(v.ID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}

	return nil
}

// MultiGet multiple get cache, return key in map is id value
func (c *groupApplyCache) MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.GroupApply, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetGroupApplyCacheKey(v)
		keys = append(keys, cacheKey)
	}

	itemMap := make(map[string]*model.GroupApply)
	err := c.cache.MultiGet(ctx, keys, itemMap)
	if err != nil {
		return nil, err
	}

	retMap := make(map[uint64]*model.GroupApply)
	for _, id := range ids {
		val, ok := itemMap[c.GetGroupApplyCacheKey(id)]
		if ok {
			retMap[id] = val
		}
	}

	return retMap, nil
}

// Del delete cache
func (c *groupApplyCache) Del(ctx context.Context, id uint64) error {
	cacheKey := c.GetGroupApplyCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetCacheWithNotFound set empty cache
func (c *groupApplyCache) SetCacheWithNotFound(ctx context.Context, id uint64) error {
	cacheKey := c.GetGroupApplyCacheKey(id)
	err := c.cache.SetCacheWithNotFound(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}
