package cache

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/zhufuyi/sponge/pkg/cache"
	"github.com/zhufuyi/sponge/pkg/encoding"
	"github.com/zhufuyi/sponge/pkg/utils"
	"lingua_exchange/pkg/jsonutil"

	"lingua_exchange/internal/model"
)

const (
	// cache prefix key, must end with a colon
	talkRecordsVoteCachePrefixKey = "talkRecordsVote:"
	// TalkRecordsVoteExpireTime expire time
	TalkRecordsVoteExpireTime = 5 * time.Minute

	VoteUsersCache = "talk:vote:answer-users:%d"

	VoteStatisticCache = "talk:vote:statistic:%d"
)

var _ TalkRecordsVoteCache = (*talkRecordsVoteCache)(nil)

// TalkRecordsVoteCache cache interface
type TalkRecordsVoteCache interface {
	Set(ctx context.Context, id uint64, data *model.TalkRecordsVote, duration time.Duration) error
	Get(ctx context.Context, id uint64) (*model.TalkRecordsVote, error)
	MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.TalkRecordsVote, error)
	MultiSet(ctx context.Context, data []*model.TalkRecordsVote, duration time.Duration) error
	Del(ctx context.Context, id uint64) error
	SetCacheWithNotFound(ctx context.Context, id uint64) error
	GetVoteAnswerUser(ctx context.Context, voteId int) ([]int, error)
	SetVoteAnswerUser(ctx context.Context, vid int, uids []int) error
	GetVoteStatistics(ctx context.Context, vid int) (string, error)
	SetVoteStatistics(ctx context.Context, vid int, value string) error
}

// talkRecordsVoteCache define a cache struct
type talkRecordsVoteCache struct {
	cache cache.Cache
	redis *redis.Client
}

func (c *talkRecordsVoteCache) GetVoteAnswerUser(ctx context.Context, voteId int) ([]int, error) {
	val, err := c.redis.Get(ctx, fmt.Sprintf(VoteUsersCache, voteId)).Result()

	if err != nil {
		return nil, err
	}

	var ids []int
	if err := jsonutil.Decode(val, &ids); err != nil {
		return nil, err
	}

	return ids, nil
}

func (c *talkRecordsVoteCache) SetVoteAnswerUser(ctx context.Context, vid int, uids []int) error {
	return c.redis.Set(ctx, fmt.Sprintf(VoteUsersCache, vid), jsonutil.Encode(uids), time.Hour*24).Err()
}

func (c *talkRecordsVoteCache) GetVoteStatistics(ctx context.Context, vid int) (string, error) {
	return c.redis.Get(ctx, fmt.Sprintf(VoteStatisticCache, vid)).Result()
}

func (c *talkRecordsVoteCache) SetVoteStatistics(ctx context.Context, vid int, value string) error {
	return c.redis.Set(ctx, fmt.Sprintf(VoteStatisticCache, vid), value, time.Hour*24).Err()
}

// NewTalkRecordsVoteCache new a cache
func NewTalkRecordsVoteCache(cacheType *model.CacheType) TalkRecordsVoteCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.TalkRecordsVote{}
		})
		return &talkRecordsVoteCache{cache: c, redis: model.GetRedisCli()}
	case "memory":
		c := cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.TalkRecordsVote{}
		})
		return &talkRecordsVoteCache{cache: c, redis: model.GetRedisCli()}
	}

	return nil // no cache
}

// GetTalkRecordsVoteCacheKey cache key
func (c *talkRecordsVoteCache) GetTalkRecordsVoteCacheKey(id uint64) string {
	return talkRecordsVoteCachePrefixKey + utils.Uint64ToStr(id)
}

// Set write to cache
func (c *talkRecordsVoteCache) Set(ctx context.Context, id uint64, data *model.TalkRecordsVote, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetTalkRecordsVoteCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get cache value
func (c *talkRecordsVoteCache) Get(ctx context.Context, id uint64) (*model.TalkRecordsVote, error) {
	var data *model.TalkRecordsVote
	cacheKey := c.GetTalkRecordsVoteCacheKey(id)
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MultiSet multiple set cache
func (c *talkRecordsVoteCache) MultiSet(ctx context.Context, data []*model.TalkRecordsVote, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetTalkRecordsVoteCacheKey(v.ID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}

	return nil
}

// MultiGet multiple get cache, return key in map is id value
func (c *talkRecordsVoteCache) MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.TalkRecordsVote, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetTalkRecordsVoteCacheKey(v)
		keys = append(keys, cacheKey)
	}

	itemMap := make(map[string]*model.TalkRecordsVote)
	err := c.cache.MultiGet(ctx, keys, itemMap)
	if err != nil {
		return nil, err
	}

	retMap := make(map[uint64]*model.TalkRecordsVote)
	for _, id := range ids {
		val, ok := itemMap[c.GetTalkRecordsVoteCacheKey(id)]
		if ok {
			retMap[id] = val
		}
	}

	return retMap, nil
}

// Del delete cache
func (c *talkRecordsVoteCache) Del(ctx context.Context, id uint64) error {
	cacheKey := c.GetTalkRecordsVoteCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetCacheWithNotFound set empty cache
func (c *talkRecordsVoteCache) SetCacheWithNotFound(ctx context.Context, id uint64) error {
	cacheKey := c.GetTalkRecordsVoteCacheKey(id)
	err := c.cache.SetCacheWithNotFound(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}
