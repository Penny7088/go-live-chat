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
	userDevicesCachePrefixKey = "userDevices:"
	// UserDevicesExpireTime expire time
	UserDevicesExpireTime = 5 * time.Minute
)

var _ UserDevicesCache = (*userDevicesCache)(nil)

// UserDevicesCache cache interface
type UserDevicesCache interface {
	Set(ctx context.Context, id uint64, data *model.UserDevices, duration time.Duration) error
	Get(ctx context.Context, id uint64) (*model.UserDevices, error)
	MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.UserDevices, error)
	MultiSet(ctx context.Context, data []*model.UserDevices, duration time.Duration) error
	Del(ctx context.Context, id uint64) error
	SetCacheWithNotFound(ctx context.Context, id uint64) error
}

// userDevicesCache define a cache struct
type userDevicesCache struct {
	cache cache.Cache
}

// NewUserDevicesCache new a cache
func NewUserDevicesCache(cacheType *model.CacheType) UserDevicesCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.UserDevices{}
		})
		return &userDevicesCache{cache: c}
	case "memory":
		c := cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.UserDevices{}
		})
		return &userDevicesCache{cache: c}
	}

	return nil // no cache
}

// GetUserDevicesCacheKey cache key
func (c *userDevicesCache) GetUserDevicesCacheKey(id uint64) string {
	return userDevicesCachePrefixKey + utils.Uint64ToStr(id)
}

// Set write to cache
func (c *userDevicesCache) Set(ctx context.Context, id uint64, data *model.UserDevices, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetUserDevicesCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get cache value
func (c *userDevicesCache) Get(ctx context.Context, id uint64) (*model.UserDevices, error) {
	var data *model.UserDevices
	cacheKey := c.GetUserDevicesCacheKey(id)
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MultiSet multiple set cache
func (c *userDevicesCache) MultiSet(ctx context.Context, data []*model.UserDevices, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetUserDevicesCacheKey(v.ID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}

	return nil
}

// MultiGet multiple get cache, return key in map is id value
func (c *userDevicesCache) MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.UserDevices, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetUserDevicesCacheKey(v)
		keys = append(keys, cacheKey)
	}

	itemMap := make(map[string]*model.UserDevices)
	err := c.cache.MultiGet(ctx, keys, itemMap)
	if err != nil {
		return nil, err
	}

	retMap := make(map[uint64]*model.UserDevices)
	for _, id := range ids {
		val, ok := itemMap[c.GetUserDevicesCacheKey(id)]
		if ok {
			retMap[id] = val
		}
	}

	return retMap, nil
}

// Del delete cache
func (c *userDevicesCache) Del(ctx context.Context, id uint64) error {
	cacheKey := c.GetUserDevicesCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetCacheWithNotFound set empty cache
func (c *userDevicesCache) SetCacheWithNotFound(ctx context.Context, id uint64) error {
	cacheKey := c.GetUserDevicesCacheKey(id)
	err := c.cache.SetCacheWithNotFound(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}
