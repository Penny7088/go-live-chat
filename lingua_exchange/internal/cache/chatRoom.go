package cache

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/zhufuyi/sponge/pkg/cache"
	"github.com/zhufuyi/sponge/pkg/encoding"
	"lingua_exchange/internal/model"
	"lingua_exchange/internal/types"
)

type ChatRoomCache interface {
	Add(ctx context.Context, opt *types.RoomOption) error
	BatchAdd(ctx context.Context, opts []*types.RoomOption) error
	Del(ctx context.Context, opt *types.RoomOption) error
	BatchDel(ctx context.Context, opts []*types.RoomOption) error
	All(ctx context.Context, opt *types.RoomOption) []int64
}

type chatRoomCache struct {
	cache cache.Cache
	redis *redis.Client
}

func NewChatRoomCache(cacheType *model.CacheType) ChatRoomCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &types.RoomOption{}
		})
		return &chatRoomCache{cache: c, redis: model.GetRedisCli()}
	case "memory":
		c := cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &types.RoomOption{}
		})
		return &chatRoomCache{cache: c, redis: model.GetRedisCli()}
	}

	return nil // no cache
}

func (c chatRoomCache) chatRoomKey(opt *types.RoomOption) string {
	return fmt.Sprintf("ws:%s:%s:%s", opt.Sid, opt.RoomType, opt.Number)
}

func (c chatRoomCache) Add(ctx context.Context, opt *types.RoomOption) error {
	key := c.chatRoomKey(opt)

	err := c.redis.SAdd(ctx, key, opt.Cid).Err()
	if err == nil {
		c.redis.Expire(ctx, key, time.Hour*24*7)
	}

	return err
}

func (c chatRoomCache) BatchAdd(ctx context.Context, opts []*types.RoomOption) error {
	pipeline := c.redis.Pipeline()
	for _, opt := range opts {
		key := c.name(opt)
		if err := pipeline.SAdd(ctx, key, opt.Cid).Err(); err == nil {
			pipeline.Expire(ctx, key, time.Hour*24*7)
		}
	}

	_, err := pipeline.Exec(ctx)
	return err
}

func (c chatRoomCache) Del(ctx context.Context, opt *types.RoomOption) error {
	return c.redis.SRem(ctx, c.name(opt), opt.Cid).Err()
}

func (c chatRoomCache) BatchDel(ctx context.Context, opts []*types.RoomOption) error {
	pipeline := c.redis.Pipeline()
	for _, opt := range opts {
		pipeline.SRem(ctx, c.name(opt), opt.Cid)
	}

	_, err := pipeline.Exec(ctx)
	return err
}

func (c chatRoomCache) All(ctx context.Context, opt *types.RoomOption) []int64 {
	arr := c.redis.SMembers(ctx, c.name(opt)).Val()

	cids := make([]int64, 0, len(arr))
	for _, val := range arr {
		if cid, err := strconv.ParseInt(val, 10, 64); err == nil {
			cids = append(cids, cid)
		}
	}

	return cids
}

func (c chatRoomCache) name(opt *types.RoomOption) string {
	return fmt.Sprintf("ws:%s:%s:%s", opt.Sid, opt.RoomType, opt.Number)
}
