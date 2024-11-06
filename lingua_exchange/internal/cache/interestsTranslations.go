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
	interestsTranslationsCachePrefixKey = "interestsTranslations:"
	// InterestsTranslationsExpireTime expire time
	InterestsTranslationsExpireTime = 5 * time.Minute
)

var _ InterestsTranslationsCache = (*interestsTranslationsCache)(nil)

// InterestsTranslationsCache cache interface
type InterestsTranslationsCache interface {
	Set(ctx context.Context, id uint64, data *model.InterestsTranslations, duration time.Duration) error
	Get(ctx context.Context, id uint64) (*model.InterestsTranslations, error)
	MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.InterestsTranslations, error)
	MultiSet(ctx context.Context, data []*model.InterestsTranslations, duration time.Duration) error
	Del(ctx context.Context, id uint64) error
	SetCacheWithNotFound(ctx context.Context, id uint64) error
}

// interestsTranslationsCache define a cache struct
type interestsTranslationsCache struct {
	cache cache.Cache
}

// NewInterestsTranslationsCache new a cache
func NewInterestsTranslationsCache(cacheType *model.CacheType) InterestsTranslationsCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.InterestsTranslations{}
		})
		return &interestsTranslationsCache{cache: c}
	case "memory":
		c := cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.InterestsTranslations{}
		})
		return &interestsTranslationsCache{cache: c}
	}

	return nil // no cache
}

// GetInterestsTranslationsCacheKey cache key
func (c *interestsTranslationsCache) GetInterestsTranslationsCacheKey(id uint64) string {
	return interestsTranslationsCachePrefixKey + utils.Uint64ToStr(id)
}

// Set write to cache
func (c *interestsTranslationsCache) Set(ctx context.Context, id uint64, data *model.InterestsTranslations, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetInterestsTranslationsCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get cache value
func (c *interestsTranslationsCache) Get(ctx context.Context, id uint64) (*model.InterestsTranslations, error) {
	var data *model.InterestsTranslations
	cacheKey := c.GetInterestsTranslationsCacheKey(id)
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MultiSet multiple set cache
func (c *interestsTranslationsCache) MultiSet(ctx context.Context, data []*model.InterestsTranslations, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetInterestsTranslationsCacheKey(v.ID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}

	return nil
}

// MultiGet multiple get cache, return key in map is id value
func (c *interestsTranslationsCache) MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.InterestsTranslations, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetInterestsTranslationsCacheKey(v)
		keys = append(keys, cacheKey)
	}

	itemMap := make(map[string]*model.InterestsTranslations)
	err := c.cache.MultiGet(ctx, keys, itemMap)
	if err != nil {
		return nil, err
	}

	retMap := make(map[uint64]*model.InterestsTranslations)
	for _, id := range ids {
		val, ok := itemMap[c.GetInterestsTranslationsCacheKey(id)]
		if ok {
			retMap[id] = val
		}
	}

	return retMap, nil
}

// Del delete cache
func (c *interestsTranslationsCache) Del(ctx context.Context, id uint64) error {
	cacheKey := c.GetInterestsTranslationsCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetCacheWithNotFound set empty cache
func (c *interestsTranslationsCache) SetCacheWithNotFound(ctx context.Context, id uint64) error {
	cacheKey := c.GetInterestsTranslationsCacheKey(id)
	err := c.cache.SetCacheWithNotFound(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}
