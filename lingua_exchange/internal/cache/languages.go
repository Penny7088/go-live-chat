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
	languagesCachePrefixKey   = "languages:"
	allLanguageCachePrefixKey = "allLanguages:"
	// LanguagesExpireTime expire time
	LanguagesExpireTime = 5 * time.Minute
)

var _ LanguagesCache = (*languagesCache)(nil)

// LanguagesCache cache interface
type LanguagesCache interface {
	Set(ctx context.Context, id uint64, data *model.Languages, duration time.Duration) error
	Get(ctx context.Context, id uint64) (*model.Languages, error)
	MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.Languages, error)
	MultiSet(ctx context.Context, data []*model.Languages, duration time.Duration) error
	Del(ctx context.Context, id uint64) error
	SetCacheWithNotFound(ctx context.Context, id uint64) error
	SetAllLanguages(ctx context.Context, data []*model.Languages, duration time.Duration) error
	GetAllLanguages(ctx context.Context) ([]*model.Languages, error)
}

// languagesCache define a cache struct
type languagesCache struct {
	cache cache.Cache
}

// NewLanguagesCache new a cache
func NewLanguagesCache(cacheType *model.CacheType) LanguagesCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.Languages{}
		})
		return &languagesCache{cache: c}
	case "memory":
		c := cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.Languages{}
		})
		return &languagesCache{cache: c}
	}

	return nil // no cache
}

// GetAllLanguagesCacheKey cache key
func (c *languagesCache) GetAllLanguagesCacheKey() string {
	return allLanguageCachePrefixKey
}

// GetLanguagesCacheKey cache key
func (c *languagesCache) GetLanguagesCacheKey(id uint64) string {
	return languagesCachePrefixKey + utils.Uint64ToStr(id)
}

func (c *languagesCache) SetAllLanguages(ctx context.Context, data []*model.Languages, duration time.Duration) error {
	if data == nil || len(data) == 0 {
		return nil
	}
	key := c.GetAllLanguagesCacheKey()
	err := c.cache.Set(ctx, key, data, duration)
	if err != nil {
		return err
	}
	return nil
}

func (c *languagesCache) GetAllLanguages(ctx context.Context) ([]*model.Languages, error) {
	var data []*model.Languages
	key := c.GetAllLanguagesCacheKey()
	err := c.cache.Get(ctx, key, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// Set write to cache
func (c *languagesCache) Set(ctx context.Context, id uint64, data *model.Languages, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetLanguagesCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get cache value
func (c *languagesCache) Get(ctx context.Context, id uint64) (*model.Languages, error) {
	var data *model.Languages
	cacheKey := c.GetLanguagesCacheKey(id)
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MultiSet multiple set cache
func (c *languagesCache) MultiSet(ctx context.Context, data []*model.Languages, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetLanguagesCacheKey(v.ID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}

	return nil
}

// MultiGet multiple get cache, return key in map is id value
func (c *languagesCache) MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.Languages, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetLanguagesCacheKey(v)
		keys = append(keys, cacheKey)
	}

	itemMap := make(map[string]*model.Languages)
	err := c.cache.MultiGet(ctx, keys, itemMap)
	if err != nil {
		return nil, err
	}

	retMap := make(map[uint64]*model.Languages)
	for _, id := range ids {
		val, ok := itemMap[c.GetLanguagesCacheKey(id)]
		if ok {
			retMap[id] = val
		}
	}

	return retMap, nil
}

// Del delete cache
func (c *languagesCache) Del(ctx context.Context, id uint64) error {
	cacheKey := c.GetLanguagesCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetCacheWithNotFound set empty cache
func (c *languagesCache) SetCacheWithNotFound(ctx context.Context, id uint64) error {
	cacheKey := c.GetLanguagesCacheKey(id)
	err := c.cache.SetCacheWithNotFound(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}
