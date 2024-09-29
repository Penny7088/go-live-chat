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
	thirdPartyAuthCachePrefixKey = "thirdPartyAuth:"
	// ThirdPartyAuthExpireTime expire time
	ThirdPartyAuthExpireTime = 5 * time.Minute
)

var _ ThirdPartyAuthCache = (*thirdPartyAuthCache)(nil)

// ThirdPartyAuthCache cache interface
type ThirdPartyAuthCache interface {
	Set(ctx context.Context, id uint64, data *model.ThirdPartyAuth, duration time.Duration) error
	Get(ctx context.Context, id uint64) (*model.ThirdPartyAuth, error)
	MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.ThirdPartyAuth, error)
	MultiSet(ctx context.Context, data []*model.ThirdPartyAuth, duration time.Duration) error
	Del(ctx context.Context, id uint64) error
	SetCacheWithNotFound(ctx context.Context, id uint64) error
}

// thirdPartyAuthCache define a cache struct
type thirdPartyAuthCache struct {
	cache cache.Cache
}

// NewThirdPartyAuthCache new a cache
func NewThirdPartyAuthCache(cacheType *model.CacheType) ThirdPartyAuthCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.ThirdPartyAuth{}
		})
		return &thirdPartyAuthCache{cache: c}
	case "memory":
		c := cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.ThirdPartyAuth{}
		})
		return &thirdPartyAuthCache{cache: c}
	}

	return nil // no cache
}

// GetThirdPartyAuthCacheKey cache key
func (c *thirdPartyAuthCache) GetThirdPartyAuthCacheKey(id uint64) string {
	return thirdPartyAuthCachePrefixKey + utils.Uint64ToStr(id)
}

// Set write to cache
func (c *thirdPartyAuthCache) Set(ctx context.Context, id uint64, data *model.ThirdPartyAuth, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetThirdPartyAuthCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get cache value
func (c *thirdPartyAuthCache) Get(ctx context.Context, id uint64) (*model.ThirdPartyAuth, error) {
	var data *model.ThirdPartyAuth
	cacheKey := c.GetThirdPartyAuthCacheKey(id)
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MultiSet multiple set cache
func (c *thirdPartyAuthCache) MultiSet(ctx context.Context, data []*model.ThirdPartyAuth, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetThirdPartyAuthCacheKey(v.ID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}

	return nil
}

// MultiGet multiple get cache, return key in map is id value
func (c *thirdPartyAuthCache) MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.ThirdPartyAuth, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetThirdPartyAuthCacheKey(v)
		keys = append(keys, cacheKey)
	}

	itemMap := make(map[string]*model.ThirdPartyAuth)
	err := c.cache.MultiGet(ctx, keys, itemMap)
	if err != nil {
		return nil, err
	}

	retMap := make(map[uint64]*model.ThirdPartyAuth)
	for _, id := range ids {
		val, ok := itemMap[c.GetThirdPartyAuthCacheKey(id)]
		if ok {
			retMap[id] = val
		}
	}

	return retMap, nil
}

// Del delete cache
func (c *thirdPartyAuthCache) Del(ctx context.Context, id uint64) error {
	cacheKey := c.GetThirdPartyAuthCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetCacheWithNotFound set empty cache
func (c *thirdPartyAuthCache) SetCacheWithNotFound(ctx context.Context, id uint64) error {
	cacheKey := c.GetThirdPartyAuthCacheKey(id)
	err := c.cache.SetCacheWithNotFound(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}
