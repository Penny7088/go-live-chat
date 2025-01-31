package cache

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/zhufuyi/sponge/pkg/cache"
	"github.com/zhufuyi/sponge/pkg/encoding"
	"lingua_exchange/internal/model"
)

// ServerCache cache interface
type ServerCache interface {
	Set(ctx context.Context, server string, time int64) error
	GetExpireServerAll(ctx context.Context) []string
	Del(ctx context.Context, server string) error
	All(ctx context.Context, status int) []string
	SetExpireServer(ctx context.Context, server string) error
	DelExpireServer(ctx context.Context, server string) error
}

type ServerModel struct {
	Server string
}

type serverCache struct {
	cache cache.Cache
	redis *redis.Client
}

// Set 更新服务心跳时间
func (s *serverCache) Set(ctx context.Context, server string, time int64) error {
	_, err := s.redis.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.SRem(ctx, ServerKeyExpire, server)
		pipe.HSet(ctx, ServerKey, server, time)
		return nil
	})
	return err
}

// Del 删除指定 ServerStorage
func (s *serverCache) Del(ctx context.Context, server string) error {
	return s.redis.HDel(ctx, ServerKey, server).Err()
}

// All 获取指定状态的运行 ServerStorage
// status 状态[1:运行中;2:已超时;3:全部]
func (s *serverCache) All(ctx context.Context, status int) []string {

	var (
		unix  = time.Now().Unix()
		slice = make([]string, 0)
	)

	all, err := s.redis.HGetAll(ctx, ServerKey).Result()
	if err != nil {
		return slice
	}

	for key, val := range all {
		value, err := strconv.Atoi(val)
		if err != nil {
			continue
		}

		switch status {
		case 1:
			if unix-int64(value) >= ServerOverTime {
				continue
			}
		case 2:
			if unix-int64(value) < ServerOverTime {
				continue
			}
		}

		slice = append(slice, key)
	}

	return slice
}

func (s *serverCache) SetExpireServer(ctx context.Context, server string) error {
	return s.redis.SAdd(ctx, ServerKeyExpire, server).Err()
}

func (s *serverCache) DelExpireServer(ctx context.Context, server string) error {
	return s.redis.SRem(ctx, ServerKeyExpire, server).Err()
}

func (s *serverCache) GetExpireServerAll(ctx context.Context) []string {
	return s.redis.SMembers(ctx, ServerKeyExpire).Val()
}

func NewServerCache(cacheType *model.CacheType) ServerCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &ServerModel{}
		})
		return &serverCache{cache: c, redis: model.GetRedisCli()}
	case "memory":
		c := cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &ServerModel{}
		})
		return &serverCache{cache: c, redis: model.GetRedisCli()}
	}

	return nil // no
}
