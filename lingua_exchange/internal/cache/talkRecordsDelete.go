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
	talkRecordsDeleteCachePrefixKey = "talkRecordsDelete:"
	// TalkRecordsDeleteExpireTime expire time
	TalkRecordsDeleteExpireTime = 5 * time.Minute
)

var _ TalkRecordsDeleteCache = (*talkRecordsDeleteCache)(nil)

// TalkRecordsDeleteCache cache interface
type TalkRecordsDeleteCache interface {
	Set(ctx context.Context, id uint64, data *model.TalkRecordsDelete, duration time.Duration) error
	Get(ctx context.Context, id uint64) (*model.TalkRecordsDelete, error)
	MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.TalkRecordsDelete, error)
	MultiSet(ctx context.Context, data []*model.TalkRecordsDelete, duration time.Duration) error
	Del(ctx context.Context, id uint64) error
	SetCacheWithNotFound(ctx context.Context, id uint64) error
}

// talkRecordsDeleteCache define a cache struct
type talkRecordsDeleteCache struct {
	cache cache.Cache
}

// NewTalkRecordsDeleteCache new a cache
func NewTalkRecordsDeleteCache(cacheType *model.CacheType) TalkRecordsDeleteCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.TalkRecordsDelete{}
		})
		return &talkRecordsDeleteCache{cache: c}
	case "memory":
		c := cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.TalkRecordsDelete{}
		})
		return &talkRecordsDeleteCache{cache: c}
	}

	return nil // no cache
}

// GetTalkRecordsDeleteCacheKey cache key
func (c *talkRecordsDeleteCache) GetTalkRecordsDeleteCacheKey(id uint64) string {
	return talkRecordsDeleteCachePrefixKey + utils.Uint64ToStr(id)
}

// Set write to cache
func (c *talkRecordsDeleteCache) Set(ctx context.Context, id uint64, data *model.TalkRecordsDelete, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetTalkRecordsDeleteCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get cache value
func (c *talkRecordsDeleteCache) Get(ctx context.Context, id uint64) (*model.TalkRecordsDelete, error) {
	var data *model.TalkRecordsDelete
	cacheKey := c.GetTalkRecordsDeleteCacheKey(id)
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MultiSet multiple set cache
func (c *talkRecordsDeleteCache) MultiSet(ctx context.Context, data []*model.TalkRecordsDelete, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetTalkRecordsDeleteCacheKey(v.ID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}

	return nil
}

// MultiGet multiple get cache, return key in map is id value
func (c *talkRecordsDeleteCache) MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.TalkRecordsDelete, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetTalkRecordsDeleteCacheKey(v)
		keys = append(keys, cacheKey)
	}

	itemMap := make(map[string]*model.TalkRecordsDelete)
	err := c.cache.MultiGet(ctx, keys, itemMap)
	if err != nil {
		return nil, err
	}

	retMap := make(map[uint64]*model.TalkRecordsDelete)
	for _, id := range ids {
		val, ok := itemMap[c.GetTalkRecordsDeleteCacheKey(id)]
		if ok {
			retMap[id] = val
		}
	}

	return retMap, nil
}

// Del delete cache
func (c *talkRecordsDeleteCache) Del(ctx context.Context, id uint64) error {
	cacheKey := c.GetTalkRecordsDeleteCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetCacheWithNotFound set empty cache
func (c *talkRecordsDeleteCache) SetCacheWithNotFound(ctx context.Context, id uint64) error {
	cacheKey := c.GetTalkRecordsDeleteCacheKey(id)
	err := c.cache.SetCacheWithNotFound(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}
