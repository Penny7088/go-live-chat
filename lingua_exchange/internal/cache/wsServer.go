package cache

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"github.com/zhufuyi/sponge/pkg/cache"
	"github.com/zhufuyi/sponge/pkg/encoding"
	"lingua_exchange/internal/model"
	"strconv"
	"strings"
	"time"
)

const (
	// ServerKey 正在的运行服务
	ServerKey = "server_ids"

	// ServerKeyExpire 过期的运行服务
	ServerKeyExpire = "server_ids_expire"

	// ServerOverTime 运行检测超时时间（单位秒）
	ServerOverTime = 50
)

type IWsServerIDCache interface {
	Set(ctx context.Context, data string, duration int64) error
	Del(ctx context.Context, server string) error
	All(ctx context.Context, status int) []string
	SetExpireServer(ctx context.Context, server string) error
	DelExpireServer(ctx context.Context, server string) error
	GetExpireServerAll(ctx context.Context) []string
}

type WsServerIDCache struct {
	cache cache.Cache
	rds   *redis.Client
}

func (c *WsServerIDCache) Set(ctx context.Context, server string, duration int64) error {
	if server == "" {
		return nil
	}
	_, err := c.invalidCacheType()
	if err != nil {
		return err
	}
	_, err = c.rds.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.SRem(ctx, ServerKeyExpire, server)
		pipe.HSet(ctx, ServerKey, server, duration)
		return nil
	})
	return err
}

func (c *WsServerIDCache) invalidCacheType() (bool, error) {
	cacheType := model.GetCacheType()
	if strings.ToLower(cacheType.CType) == "redis" {
		return true, nil
	} else {
		return false, errors.New("invalid cache type")
	}
}

func (c *WsServerIDCache) Del(ctx context.Context, server string) error {
	_, err := c.invalidCacheType()
	if err != nil {
		return err
	}
	return c.rds.HDel(ctx, ServerKey, server).Err()
}

func (c *WsServerIDCache) All(ctx context.Context, status int) []string {

	var (
		unix  = time.Now().Unix()
		slice = make([]string, 0)
	)

	all, err := c.rds.HGetAll(ctx, ServerKey).Result()
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

func (c *WsServerIDCache) SetExpireServer(ctx context.Context, server string) error {
	return c.rds.SAdd(ctx, ServerKeyExpire, server).Err()
}

func (c *WsServerIDCache) DelExpireServer(ctx context.Context, server string) error {
	return c.rds.SRem(ctx, ServerKeyExpire, server).Err()
}

func (c *WsServerIDCache) GetExpireServerAll(ctx context.Context) []string {
	return c.rds.SMembers(ctx, ServerKeyExpire).Val()
}

func NewWsServerCache(cacheType *model.CacheType) *WsServerIDCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.WSServer{}
		})
		return &WsServerIDCache{cache: c, rds: model.GetRedisCli()}
	case "memory":
		c := cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.WSServer{}
		})
		return &WsServerIDCache{cache: c}
	}

	return nil // no cache
}
