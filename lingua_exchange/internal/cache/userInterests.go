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
	userInterestsCachePrefixKey = "userInterests:"
	// UserInterestsExpireTime expire time
	UserInterestsExpireTime = 5 * time.Minute
)

var _ UserInterestsCache = (*userInterestsCache)(nil)

// UserInterestsCache cache interface
type UserInterestsCache interface {
	Set(ctx context.Context, id uint64, data *model.UserInterests, duration time.Duration) error
	Get(ctx context.Context, id uint64) (*model.UserInterests, error)
	MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.UserInterests, error)
	MultiSet(ctx context.Context, data []*model.UserInterests, duration time.Duration) error
	Del(ctx context.Context, id uint64) error
	SetCacheWithNotFound(ctx context.Context, id uint64) error
}

// userInterestsCache define a cache struct
type userInterestsCache struct {
	cache cache.Cache
}

// NewUserInterestsCache new a cache
func NewUserInterestsCache(cacheType *model.CacheType) UserInterestsCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.UserInterests{}
		})
		return &userInterestsCache{cache: c}
	case "memory":
		c := cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.UserInterests{}
		})
		return &userInterestsCache{cache: c}
	}

	return nil // no cache
}

// GetUserInterestsCacheKey cache key
func (c *userInterestsCache) GetUserInterestsCacheKey(id uint64) string {
	return userInterestsCachePrefixKey + utils.Uint64ToStr(id)
}

// Set write to cache
func (c *userInterestsCache) Set(ctx context.Context, id uint64, data *model.UserInterests, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetUserInterestsCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get cache value
func (c *userInterestsCache) Get(ctx context.Context, id uint64) (*model.UserInterests, error) {
	var data *model.UserInterests
	cacheKey := c.GetUserInterestsCacheKey(id)
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MultiSet multiple set cache
func (c *userInterestsCache) MultiSet(ctx context.Context, data []*model.UserInterests, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetUserInterestsCacheKey(v.ID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}

	return nil
}

// MultiGet multiple get cache, return key in map is id value
func (c *userInterestsCache) MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.UserInterests, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetUserInterestsCacheKey(v)
		keys = append(keys, cacheKey)
	}

	itemMap := make(map[string]*model.UserInterests)
	err := c.cache.MultiGet(ctx, keys, itemMap)
	if err != nil {
		return nil, err
	}

	retMap := make(map[uint64]*model.UserInterests)
	for _, id := range ids {
		val, ok := itemMap[c.GetUserInterestsCacheKey(id)]
		if ok {
			retMap[id] = val
		}
	}

	return retMap, nil
}

// Del delete cache
func (c *userInterestsCache) Del(ctx context.Context, id uint64) error {
	cacheKey := c.GetUserInterestsCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetCacheWithNotFound set empty cache
func (c *userInterestsCache) SetCacheWithNotFound(ctx context.Context, id uint64) error {
	cacheKey := c.GetUserInterestsCacheKey(id)
	err := c.cache.SetCacheWithNotFound(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}
