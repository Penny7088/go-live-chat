package cache

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/zhufuyi/sponge/pkg/cache"
	"github.com/zhufuyi/sponge/pkg/encoding"
	"lingua_exchange/internal/model"
)

const (
	// cache prefix key, must end with a colon
	talkRecordsCachePrefixKey = "talkRecords:"
	// TalkRecordsExpireTime expire time
	TalkRecordsExpireTime = 5 * time.Minute
)

var _ TalkRecordsCache = (*talkRecordsCache)(nil)

// TalkRecordsCache cache interface
type TalkRecordsCache interface {
	Set(ctx context.Context, id string, data *model.TalkRecords, duration time.Duration) error
	Get(ctx context.Context, id string) (*model.TalkRecords, error)
	MultiGet(ctx context.Context, ids []string) (map[string]*model.TalkRecords, error)
	MultiSet(ctx context.Context, data []*model.TalkRecords, duration time.Duration) error
	Del(ctx context.Context, id string) error
	SetCacheWithNotFound(ctx context.Context, id string) error
	SetSequence(ctx context.Context, userId int, receiverId int, value int64) error
	GetSequence(ctx context.Context, userId int, receiverId int) int64
	BatchGetSequence(ctx context.Context, userId int, receiverId int, num int64) []int64
}

// talkRecordsCache define a cache struct
type talkRecordsCache struct {
	cache cache.Cache
	redis *redis.Client
}

func (s *talkRecordsCache) SequenceName(userId int, receiverId int) string {

	if userId == 0 {
		return fmt.Sprintf("im:sequence:chat:%d", receiverId)
	}

	if receiverId < userId {
		receiverId, userId = userId, receiverId
	}

	return fmt.Sprintf("im:sequence:chat:%d_%d", userId, receiverId)
}

// SetSequence 初始化发号器
func (c *talkRecordsCache) SetSequence(ctx context.Context, userId int, receiverId int, value int64) error {
	return c.redis.SetEx(ctx, c.SequenceName(userId, receiverId), value, 12*time.Hour).Err()
}

// GetSequence 获取消息时序ID
func (c *talkRecordsCache) GetSequence(ctx context.Context, userId int, receiverId int) int64 {
	return c.redis.Incr(ctx, c.SequenceName(userId, receiverId)).Val()
}

// BatchGetSequence 批量获取消息时序ID
func (c *talkRecordsCache) BatchGetSequence(ctx context.Context, userId int, receiverId int, num int64) []int64 {
	value := c.redis.IncrBy(ctx, c.SequenceName(userId, receiverId), num).Val()

	items := make([]int64, 0, num)
	for i := num; i > 0; i-- {
		items = append(items, value-i+1)
	}
	return items
}

// NewTalkRecordsCache new a cache
func NewTalkRecordsCache(cacheType *model.CacheType) TalkRecordsCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.TalkRecords{}
		})
		return &talkRecordsCache{cache: c, redis: model.GetRedisCli()}
	case "memory":
		c := cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.TalkRecords{}
		})
		return &talkRecordsCache{cache: c}
	}

	return nil // no cache
}

// GetTalkRecordsCacheKey cache key
func (c *talkRecordsCache) GetTalkRecordsCacheKey(id string) string {
	return talkRecordsCachePrefixKey + id
}

// Set write to cache
func (c *talkRecordsCache) Set(ctx context.Context, id string, data *model.TalkRecords, duration time.Duration) error {
	if data == nil || id == "" {
		return nil
	}
	cacheKey := c.GetTalkRecordsCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get cache value
func (c *talkRecordsCache) Get(ctx context.Context, id string) (*model.TalkRecords, error) {
	var data *model.TalkRecords
	cacheKey := c.GetTalkRecordsCacheKey(id)
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MultiSet multiple set cache
func (c *talkRecordsCache) MultiSet(ctx context.Context, data []*model.TalkRecords, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetTalkRecordsCacheKey(v.MsgID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}

	return nil
}

// MultiGet multiple get cache, return key in map is id value
func (c *talkRecordsCache) MultiGet(ctx context.Context, ids []string) (map[string]*model.TalkRecords, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetTalkRecordsCacheKey(v)
		keys = append(keys, cacheKey)
	}

	itemMap := make(map[string]*model.TalkRecords)
	err := c.cache.MultiGet(ctx, keys, itemMap)
	if err != nil {
		return nil, err
	}

	retMap := make(map[string]*model.TalkRecords)
	for _, id := range ids {
		val, ok := itemMap[c.GetTalkRecordsCacheKey(id)]
		if ok {
			retMap[id] = val
		}
	}

	return retMap, nil
}

// Del delete cache
func (c *talkRecordsCache) Del(ctx context.Context, id string) error {
	cacheKey := c.GetTalkRecordsCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetCacheWithNotFound set empty cache
func (c *talkRecordsCache) SetCacheWithNotFound(ctx context.Context, id string) error {
	cacheKey := c.GetTalkRecordsCacheKey(id)
	err := c.cache.SetCacheWithNotFound(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}
