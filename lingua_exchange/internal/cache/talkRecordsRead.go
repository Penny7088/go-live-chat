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
	talkRecordsReadCachePrefixKey = "talkRecordsRead:"
	// TalkRecordsReadExpireTime expire time
	TalkRecordsReadExpireTime = 5 * time.Minute
)

var _ TalkRecordsReadCache = (*talkRecordsReadCache)(nil)

// TalkRecordsReadCache cache interface
type TalkRecordsReadCache interface {
	Set(ctx context.Context, id uint64, data *model.TalkRecordsRead, duration time.Duration) error
	Get(ctx context.Context, id uint64) (*model.TalkRecordsRead, error)
	MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.TalkRecordsRead, error)
	MultiSet(ctx context.Context, data []*model.TalkRecordsRead, duration time.Duration) error
	Del(ctx context.Context, id uint64) error
	SetCacheWithNotFound(ctx context.Context, id uint64) error
}

// talkRecordsReadCache define a cache struct
type talkRecordsReadCache struct {
	cache cache.Cache
}

// NewTalkRecordsReadCache new a cache
func NewTalkRecordsReadCache(cacheType *model.CacheType) TalkRecordsReadCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.TalkRecordsRead{}
		})
		return &talkRecordsReadCache{cache: c}
	case "memory":
		c := cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.TalkRecordsRead{}
		})
		return &talkRecordsReadCache{cache: c}
	}

	return nil // no cache
}

// GetTalkRecordsReadCacheKey cache key
func (c *talkRecordsReadCache) GetTalkRecordsReadCacheKey(id uint64) string {
	return talkRecordsReadCachePrefixKey + utils.Uint64ToStr(id)
}

// Set write to cache
func (c *talkRecordsReadCache) Set(ctx context.Context, id uint64, data *model.TalkRecordsRead, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetTalkRecordsReadCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get cache value
func (c *talkRecordsReadCache) Get(ctx context.Context, id uint64) (*model.TalkRecordsRead, error) {
	var data *model.TalkRecordsRead
	cacheKey := c.GetTalkRecordsReadCacheKey(id)
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MultiSet multiple set cache
func (c *talkRecordsReadCache) MultiSet(ctx context.Context, data []*model.TalkRecordsRead, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetTalkRecordsReadCacheKey(v.ID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}

	return nil
}

// MultiGet multiple get cache, return key in map is id value
func (c *talkRecordsReadCache) MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.TalkRecordsRead, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetTalkRecordsReadCacheKey(v)
		keys = append(keys, cacheKey)
	}

	itemMap := make(map[string]*model.TalkRecordsRead)
	err := c.cache.MultiGet(ctx, keys, itemMap)
	if err != nil {
		return nil, err
	}

	retMap := make(map[uint64]*model.TalkRecordsRead)
	for _, id := range ids {
		val, ok := itemMap[c.GetTalkRecordsReadCacheKey(id)]
		if ok {
			retMap[id] = val
		}
	}

	return retMap, nil
}

// Del delete cache
func (c *talkRecordsReadCache) Del(ctx context.Context, id uint64) error {
	cacheKey := c.GetTalkRecordsReadCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetCacheWithNotFound set empty cache
func (c *talkRecordsReadCache) SetCacheWithNotFound(ctx context.Context, id uint64) error {
	cacheKey := c.GetTalkRecordsReadCacheKey(id)
	err := c.cache.SetCacheWithNotFound(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}
