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
	userLanguagesCachePrefixKey = "userLanguages:"
	// UserLanguagesExpireTime expire time
	UserLanguagesExpireTime = 5 * time.Minute
)

var _ UserLanguagesCache = (*userLanguagesCache)(nil)

// UserLanguagesCache cache interface
type UserLanguagesCache interface {
	Set(ctx context.Context, id uint64, data *model.UserLanguages, duration time.Duration) error
	Get(ctx context.Context, id uint64) (*model.UserLanguages, error)
	MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.UserLanguages, error)
	MultiSet(ctx context.Context, data []*model.UserLanguages, duration time.Duration) error
	Del(ctx context.Context, id uint64) error
	SetCacheWithNotFound(ctx context.Context, id uint64) error
}

// userLanguagesCache define a cache struct
type userLanguagesCache struct {
	cache cache.Cache
}

// NewUserLanguagesCache new a cache
func NewUserLanguagesCache(cacheType *model.CacheType) UserLanguagesCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.UserLanguages{}
		})
		return &userLanguagesCache{cache: c}
	case "memory":
		c := cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.UserLanguages{}
		})
		return &userLanguagesCache{cache: c}
	}

	return nil // no cache
}

// GetUserLanguagesCacheKey cache key
func (c *userLanguagesCache) GetUserLanguagesCacheKey(id uint64) string {
	return userLanguagesCachePrefixKey + utils.Uint64ToStr(id)
}

// Set write to cache
func (c *userLanguagesCache) Set(ctx context.Context, id uint64, data *model.UserLanguages, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetUserLanguagesCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get cache value
func (c *userLanguagesCache) Get(ctx context.Context, id uint64) (*model.UserLanguages, error) {
	var data *model.UserLanguages
	cacheKey := c.GetUserLanguagesCacheKey(id)
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MultiSet multiple set cache
func (c *userLanguagesCache) MultiSet(ctx context.Context, data []*model.UserLanguages, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetUserLanguagesCacheKey(v.ID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}

	return nil
}

// MultiGet multiple get cache, return key in map is id value
func (c *userLanguagesCache) MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.UserLanguages, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetUserLanguagesCacheKey(v)
		keys = append(keys, cacheKey)
	}

	itemMap := make(map[string]*model.UserLanguages)
	err := c.cache.MultiGet(ctx, keys, itemMap)
	if err != nil {
		return nil, err
	}

	retMap := make(map[uint64]*model.UserLanguages)
	for _, id := range ids {
		val, ok := itemMap[c.GetUserLanguagesCacheKey(id)]
		if ok {
			retMap[id] = val
		}
	}

	return retMap, nil
}

// Del delete cache
func (c *userLanguagesCache) Del(ctx context.Context, id uint64) error {
	cacheKey := c.GetUserLanguagesCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetCacheWithNotFound set empty cache
func (c *userLanguagesCache) SetCacheWithNotFound(ctx context.Context, id uint64) error {
	cacheKey := c.GetUserLanguagesCacheKey(id)
	err := c.cache.SetCacheWithNotFound(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}
