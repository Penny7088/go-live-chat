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
	interestsCachePrefixKey = "interests:"
	// InterestsExpireTime expire time
	InterestsExpireTime = 5 * time.Minute
)

var _ InterestsCache = (*interestsCache)(nil)

// InterestsCache cache interface
type InterestsCache interface {
	Set(ctx context.Context, id uint64, data *model.Interests, duration time.Duration) error
	Get(ctx context.Context, id uint64) (*model.Interests, error)
	MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.Interests, error)
	MultiSet(ctx context.Context, data []*model.Interests, duration time.Duration) error
	Del(ctx context.Context, id uint64) error
	SetCacheWithNotFound(ctx context.Context, id uint64) error
	GetFromLanguageCodeAllInterestsCache(ctx context.Context, languageCode string) ([]*model.InterestsTranslations, error)
	SetFromLanguageCodeAllInterestsCache(ctx context.Context, languageCode string, data []*model.InterestsTranslations, duration time.Duration) error
}

// interestsCache define a cache struct
type interestsCache struct {
	cache cache.Cache
}

// NewInterestsCache new a cache
func NewInterestsCache(cacheType *model.CacheType) InterestsCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.Interests{}
		})
		return &interestsCache{cache: c}
	case "memory":
		c := cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.Interests{}
		})
		return &interestsCache{cache: c}
	}

	return nil // no cache
}

// GetInterestsCacheKey cache key
func (c *interestsCache) GetInterestsCacheKey(id uint64) string {
	return interestsCachePrefixKey + utils.Uint64ToStr(id)
}

func (c *interestsCache) GetFromLanguageInterestsCacheKey(languageCode string) string {
	return interestsCachePrefixKey + languageCode
}

func (c *interestsCache) GetFromLanguageCodeAllInterestsCache(ctx context.Context, languageCode string) ([]*model.InterestsTranslations, error) {
	cacheKey := c.GetFromLanguageInterestsCacheKey(languageCode)
	var interests []*model.InterestsTranslations
	err := c.cache.Get(ctx, cacheKey, &interests)
	if err != nil {
		return nil, err
	}
	return interests, nil
}

func (c *interestsCache) SetFromLanguageCodeAllInterestsCache(ctx context.Context, languageCode string, data []*model.InterestsTranslations, duration time.Duration) error {
	if data == nil || languageCode == "" {
		return nil
	}
	cacheKey := c.GetFromLanguageInterestsCacheKey(languageCode)

	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Set write to cache
func (c *interestsCache) Set(ctx context.Context, id uint64, data *model.Interests, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetInterestsCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get cache value
func (c *interestsCache) Get(ctx context.Context, id uint64) (*model.Interests, error) {
	var data *model.Interests
	cacheKey := c.GetInterestsCacheKey(id)
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MultiSet multiple set cache
func (c *interestsCache) MultiSet(ctx context.Context, data []*model.Interests, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetInterestsCacheKey(v.ID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}

	return nil
}

// MultiGet multiple get cache, return key in map is id value
func (c *interestsCache) MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.Interests, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetInterestsCacheKey(v)
		keys = append(keys, cacheKey)
	}

	itemMap := make(map[string]*model.Interests)
	err := c.cache.MultiGet(ctx, keys, itemMap)
	if err != nil {
		return nil, err
	}

	retMap := make(map[uint64]*model.Interests)
	for _, id := range ids {
		val, ok := itemMap[c.GetInterestsCacheKey(id)]
		if ok {
			retMap[id] = val
		}
	}

	return retMap, nil
}

// Del delete cache
func (c *interestsCache) Del(ctx context.Context, id uint64) error {
	cacheKey := c.GetInterestsCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetCacheWithNotFound set empty cache
func (c *interestsCache) SetCacheWithNotFound(ctx context.Context, id uint64) error {
	cacheKey := c.GetInterestsCacheKey(id)
	err := c.cache.SetCacheWithNotFound(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}
