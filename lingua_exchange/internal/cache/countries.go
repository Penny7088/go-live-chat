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
	countriesCachePrefixKey = "countries:"
	// CountriesExpireTime expire time
	CountriesExpireTime = 5 * time.Minute
)

var _ CountriesCache = (*countriesCache)(nil)

// CountriesCache cache interface
type CountriesCache interface {
	Set(ctx context.Context, id uint64, data *model.Countries, duration time.Duration) error
	Get(ctx context.Context, id uint64) (*model.Countries, error)
	MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.Countries, error)
	MultiSet(ctx context.Context, data []*model.Countries, duration time.Duration) error
	Del(ctx context.Context, id uint64) error
	SetCacheWithNotFound(ctx context.Context, id uint64) error
}

// countriesCache define a cache struct
type countriesCache struct {
	cache cache.Cache
}

// NewCountriesCache new a cache
func NewCountriesCache(cacheType *model.CacheType) CountriesCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.Countries{}
		})
		return &countriesCache{cache: c}
	case "memory":
		c := cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.Countries{}
		})
		return &countriesCache{cache: c}
	}

	return nil // no cache
}

// GetCountriesCacheKey cache key
func (c *countriesCache) GetCountriesCacheKey(id uint64) string {
	return countriesCachePrefixKey + utils.Uint64ToStr(id)
}

// Set write to cache
func (c *countriesCache) Set(ctx context.Context, id uint64, data *model.Countries, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetCountriesCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get cache value
func (c *countriesCache) Get(ctx context.Context, id uint64) (*model.Countries, error) {
	var data *model.Countries
	cacheKey := c.GetCountriesCacheKey(id)
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MultiSet multiple set cache
func (c *countriesCache) MultiSet(ctx context.Context, data []*model.Countries, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetCountriesCacheKey(v.ID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}

	return nil
}

// MultiGet multiple get cache, return key in map is id value
func (c *countriesCache) MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.Countries, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetCountriesCacheKey(v)
		keys = append(keys, cacheKey)
	}

	itemMap := make(map[string]*model.Countries)
	err := c.cache.MultiGet(ctx, keys, itemMap)
	if err != nil {
		return nil, err
	}

	retMap := make(map[uint64]*model.Countries)
	for _, id := range ids {
		val, ok := itemMap[c.GetCountriesCacheKey(id)]
		if ok {
			retMap[id] = val
		}
	}

	return retMap, nil
}

// Del delete cache
func (c *countriesCache) Del(ctx context.Context, id uint64) error {
	cacheKey := c.GetCountriesCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetCacheWithNotFound set empty cache
func (c *countriesCache) SetCacheWithNotFound(ctx context.Context, id uint64) error {
	cacheKey := c.GetCountriesCacheKey(id)
	err := c.cache.SetCacheWithNotFound(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}
