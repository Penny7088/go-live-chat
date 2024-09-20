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
	countryLanguagesCachePrefixKey = "countryLanguages:"
	// CountryLanguagesExpireTime expire time
	CountryLanguagesExpireTime = 5 * time.Minute
)

var _ CountryLanguagesCache = (*countryLanguagesCache)(nil)

// CountryLanguagesCache cache interface
type CountryLanguagesCache interface {
	Set(ctx context.Context, id uint64, data *model.CountryLanguages, duration time.Duration) error
	Get(ctx context.Context, id uint64) (*model.CountryLanguages, error)
	MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.CountryLanguages, error)
	MultiSet(ctx context.Context, data []*model.CountryLanguages, duration time.Duration) error
	Del(ctx context.Context, id uint64) error
	SetCacheWithNotFound(ctx context.Context, id uint64) error
}

// countryLanguagesCache define a cache struct
type countryLanguagesCache struct {
	cache cache.Cache
}

// NewCountryLanguagesCache new a cache
func NewCountryLanguagesCache(cacheType *model.CacheType) CountryLanguagesCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.CountryLanguages{}
		})
		return &countryLanguagesCache{cache: c}
	case "memory":
		c := cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.CountryLanguages{}
		})
		return &countryLanguagesCache{cache: c}
	}

	return nil // no cache
}

// GetCountryLanguagesCacheKey cache key
func (c *countryLanguagesCache) GetCountryLanguagesCacheKey(id uint64) string {
	return countryLanguagesCachePrefixKey + utils.Uint64ToStr(id)
}

// Set write to cache
func (c *countryLanguagesCache) Set(ctx context.Context, id uint64, data *model.CountryLanguages, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetCountryLanguagesCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get cache value
func (c *countryLanguagesCache) Get(ctx context.Context, id uint64) (*model.CountryLanguages, error) {
	var data *model.CountryLanguages
	cacheKey := c.GetCountryLanguagesCacheKey(id)
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MultiSet multiple set cache
func (c *countryLanguagesCache) MultiSet(ctx context.Context, data []*model.CountryLanguages, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetCountryLanguagesCacheKey(v.ID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}

	return nil
}

// MultiGet multiple get cache, return key in map is id value
func (c *countryLanguagesCache) MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.CountryLanguages, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetCountryLanguagesCacheKey(v)
		keys = append(keys, cacheKey)
	}

	itemMap := make(map[string]*model.CountryLanguages)
	err := c.cache.MultiGet(ctx, keys, itemMap)
	if err != nil {
		return nil, err
	}

	retMap := make(map[uint64]*model.CountryLanguages)
	for _, id := range ids {
		val, ok := itemMap[c.GetCountryLanguagesCacheKey(id)]
		if ok {
			retMap[id] = val
		}
	}

	return retMap, nil
}

// Del delete cache
func (c *countryLanguagesCache) Del(ctx context.Context, id uint64) error {
	cacheKey := c.GetCountryLanguagesCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetCacheWithNotFound set empty cache
func (c *countryLanguagesCache) SetCacheWithNotFound(ctx context.Context, id uint64) error {
	cacheKey := c.GetCountryLanguagesCacheKey(id)
	err := c.cache.SetCacheWithNotFound(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}
