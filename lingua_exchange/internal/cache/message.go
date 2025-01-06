package cache

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/redis/go-redis/v9"
	basicCache "github.com/zhufuyi/sponge/pkg/cache"
	"github.com/zhufuyi/sponge/pkg/encoding"
	"lingua_exchange/internal/config"
	"lingua_exchange/internal/model"
)

type MessageCache interface {
	Set(ctx context.Context, channel string, fd string, uid int) error
	Del(ctx context.Context, channel, fd string) error
	IsOnline(ctx context.Context, channel, uid string) bool
	IsCurrentServerOnline(ctx context.Context, sid, channel, uid string) bool
	GetUidFromClientIds(ctx context.Context, sid, channel, uid string) []int64
	GetClientIdFromUid(ctx context.Context, sid, channel, cid string) (int64, error)
}

type messageCache struct {
	cache       basicCache.Cache
	redisClient *redis.Client
	config      *config.Config
	storage     ServerCache
}

// Set 设置客户端与用户绑定关系
// @params channel  渠道分组
// @params fd       客户端连接ID
// @params id       用户ID
func (m *messageCache) Set(ctx context.Context, channel string, fd string, uid int) error {
	sid := m.config.App.Sid
	_, err := m.redisClient.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.HSet(ctx, m.clientKey(sid, channel), fd, uid)
		pipe.SAdd(ctx, m.userKey(sid, channel, strconv.Itoa(uid)), fd)
		return nil
	})
	return err
}

// Del 删除客户端与用户绑定关系
// @params channel  渠道分组
// @params fd       客户端连接ID
func (m *messageCache) Del(ctx context.Context, channel, fd string) error {
	sid := m.config.App.Sid
	key := m.clientKey(sid, channel)
	uid, _ := m.redisClient.HGet(ctx, key, fd).Result()
	_, err := m.redisClient.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.HDel(ctx, key, fd)
		pipe.SRem(ctx, m.userKey(sid, channel, uid), fd)
		return nil
	})
	return err
}

// IsOnline 判断客户端是否在线[所有部署机器]
// @params channel  渠道分组
// @params uid      用户ID
func (m *messageCache) IsOnline(ctx context.Context, channel, uid string) bool {
	for _, sid := range m.storage.All(ctx, 1) {
		if m.IsCurrentServerOnline(ctx, sid, channel, uid) {
			return true
		}
	}

	return false
}

// IsCurrentServerOnline 判断当前节点是否在线
// @params sid      服务ID
// @params channel  渠道分组
// @params uid      用户ID
func (m *messageCache) IsCurrentServerOnline(ctx context.Context, sid, channel, uid string) bool {
	val, err := m.redisClient.SCard(ctx, m.userKey(sid, channel, uid)).Result()
	return err == nil && val > 0
}

// GetUidFromClientIds 获取当前节点用户ID关联的客户端ID
// @params sid      服务ID
// @params channel  渠道分组
// @params uid      用户ID
func (m *messageCache) GetUidFromClientIds(ctx context.Context, sid, channel, uid string) []int64 {
	cids := make([]int64, 0)

	items, err := m.redisClient.SMembers(ctx, m.userKey(sid, channel, uid)).Result()
	if err != nil {
		return cids
	}

	for _, cid := range items {
		if cid, err := strconv.ParseInt(cid, 10, 64); err == nil {
			cids = append(cids, cid)
		}
	}

	return cids
}

// GetClientIdFromUid 获取客户端ID关联的用户ID
// @params sid     服务节点ID
// @params channel 渠道分组
// @params cid     客户端ID
func (m *messageCache) GetClientIdFromUid(ctx context.Context, sid, channel, cid string) (int64, error) {
	uid, err := m.redisClient.HGet(ctx, m.clientKey(sid, channel), cid).Result()
	if err != nil {
		return 0, err
	}

	return strconv.ParseInt(uid, 10, 64)
}

func (m *messageCache) Bind(ctx context.Context, channel string, clientId int64, uid int) error {
	return m.Set(ctx, channel, strconv.FormatInt(clientId, 10), uid)
}

func (m *messageCache) UnBind(ctx context.Context, channel string, clientId int64) error {
	return m.Del(ctx, channel, strconv.FormatInt(clientId, 10))
}

func (m *messageCache) clientKey(sid, channel string) string {
	return fmt.Sprintf("ws:%s:%s:client", sid, channel)
}

func (m *messageCache) userKey(sid, channel, uid string) string {
	return fmt.Sprintf("ws:%s:%s:user:%s", sid, channel, uid)
}

type MessageModel struct {
}

func NewMessageCache(cacheType *model.CacheType) MessageCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := basicCache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &MessageModel{}
		})
		return &messageCache{cache: c, redisClient: model.GetRedisCli(), config: config.Get(), storage: NewServerCache(model.GetCacheType())}
	case "memory":
		c := basicCache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &MessageModel{}
		})
		return &messageCache{cache: c, redisClient: model.GetRedisCli(), config: config.Get(), storage: NewServerCache(model.GetCacheType())}
	}

	return nil // no cache
}
