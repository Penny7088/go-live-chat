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
	groupNoticeCachePrefixKey = "groupNotice:"
	// GroupNoticeExpireTime expire time
	GroupNoticeExpireTime = 5 * time.Minute
)

var _ GroupNoticeCache = (*groupNoticeCache)(nil)

// GroupNoticeCache cache interface
type GroupNoticeCache interface {
	Set(ctx context.Context, id uint64, data *model.GroupNotice, duration time.Duration) error
	Get(ctx context.Context, id uint64) (*model.GroupNotice, error)
	MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.GroupNotice, error)
	MultiSet(ctx context.Context, data []*model.GroupNotice, duration time.Duration) error
	Del(ctx context.Context, id uint64) error
	SetCacheWithNotFound(ctx context.Context, id uint64) error
}

// groupNoticeCache define a cache struct
type groupNoticeCache struct {
	cache cache.Cache
}

// NewGroupNoticeCache new a cache
func NewGroupNoticeCache(cacheType *model.CacheType) GroupNoticeCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.GroupNotice{}
		})
		return &groupNoticeCache{cache: c}
	case "memory":
		c := cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.GroupNotice{}
		})
		return &groupNoticeCache{cache: c}
	}

	return nil // no cache
}

// GetGroupNoticeCacheKey cache key
func (c *groupNoticeCache) GetGroupNoticeCacheKey(id uint64) string {
	return groupNoticeCachePrefixKey + utils.Uint64ToStr(id)
}

// Set write to cache
func (c *groupNoticeCache) Set(ctx context.Context, id uint64, data *model.GroupNotice, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetGroupNoticeCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get cache value
func (c *groupNoticeCache) Get(ctx context.Context, id uint64) (*model.GroupNotice, error) {
	var data *model.GroupNotice
	cacheKey := c.GetGroupNoticeCacheKey(id)
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MultiSet multiple set cache
func (c *groupNoticeCache) MultiSet(ctx context.Context, data []*model.GroupNotice, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetGroupNoticeCacheKey(v.ID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}

	return nil
}

// MultiGet multiple get cache, return key in map is id value
func (c *groupNoticeCache) MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.GroupNotice, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetGroupNoticeCacheKey(v)
		keys = append(keys, cacheKey)
	}

	itemMap := make(map[string]*model.GroupNotice)
	err := c.cache.MultiGet(ctx, keys, itemMap)
	if err != nil {
		return nil, err
	}

	retMap := make(map[uint64]*model.GroupNotice)
	for _, id := range ids {
		val, ok := itemMap[c.GetGroupNoticeCacheKey(id)]
		if ok {
			retMap[id] = val
		}
	}

	return retMap, nil
}

// Del delete cache
func (c *groupNoticeCache) Del(ctx context.Context, id uint64) error {
	cacheKey := c.GetGroupNoticeCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetCacheWithNotFound set empty cache
func (c *groupNoticeCache) SetCacheWithNotFound(ctx context.Context, id uint64) error {
	cacheKey := c.GetGroupNoticeCacheKey(id)
	err := c.cache.SetCacheWithNotFound(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}
