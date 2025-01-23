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
	"lingua_exchange/pkg/jsonutil"
)

const lastMessageCacheKey = "redis:hash:last-message"

type MessageCache struct {
	cache       basicCache.Cache
	redisClient *redis.Client
	config      *config.Config
	storage     ServerCache
}

type MessageModel struct {
}
type LastCacheMessage struct {
	Content  string `json:"content"`
	Datetime string `json:"datetime"`
}

func (m *MessageCache) name(talkType int, sender int, receive int) string {
	if talkType == 2 {
		sender = 0
	}

	if sender > receive {
		sender, receive = receive, sender
	}

	return fmt.Sprintf("%d_%d_%d", talkType, sender, receive)
}

func (m *MessageCache) SetLastMessage(ctx context.Context, talkType int, sender int, receive int, message *LastCacheMessage) error {
	text := jsonutil.Encode(message)

	return m.redisClient.HSet(ctx, lastMessageCacheKey, m.name(talkType, sender, receive), text).Err()
}

func (m *MessageCache) GetLastMessage(ctx context.Context, talkType int, sender int, receive int) (*LastCacheMessage, error) {

	res, err := m.redisClient.HGet(ctx, lastMessageCacheKey, m.name(talkType, sender, receive)).Result()
	if err != nil {
		return nil, err
	}

	msg := &LastCacheMessage{}
	if err = jsonutil.Decode(res, msg); err != nil {
		return nil, err
	}

	return msg, nil
}

func (m *MessageCache) MGetLastMessage(ctx context.Context, fields []string) ([]*LastCacheMessage, error) {

	res := m.redisClient.HMGet(ctx, lastMessageCacheKey, fields...)

	items := make([]*LastCacheMessage, 0)
	for _, item := range res.Val() {
		if val, ok := item.(string); ok {
			msg := &LastCacheMessage{}
			if err := jsonutil.Decode(val, msg); err != nil {
				return nil, err
			}

			items = append(items, msg)
		}
	}

	return items, nil
}

// Set 设置客户端与用户绑定关系
// @params channel  渠道分组
// @params fd       客户端连接ID
// @params id       用户ID
func (m *MessageCache) Set(ctx context.Context, channel string, fd string, uid int) error {
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
func (m *MessageCache) Del(ctx context.Context, channel, fd string) error {
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
func (m *MessageCache) IsOnline(ctx context.Context, channel, uid string) bool {
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
func (m *MessageCache) IsCurrentServerOnline(ctx context.Context, sid, channel, uid string) bool {
	val, err := m.redisClient.SCard(ctx, m.userKey(sid, channel, uid)).Result()
	return err == nil && val > 0
}

// GetUidFromClientIds 获取当前节点用户ID关联的客户端ID
// @params sid      服务ID
// @params channel  渠道分组
// @params uid      用户ID
func (m *MessageCache) GetUidFromClientIds(ctx context.Context, sid, channel, uid string) []int64 {
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
func (m *MessageCache) GetClientIdFromUid(ctx context.Context, sid, channel, cid string) (int64, error) {
	uid, err := m.redisClient.HGet(ctx, m.clientKey(sid, channel), cid).Result()
	if err != nil {
		return 0, err
	}

	return strconv.ParseInt(uid, 10, 64)
}

func (m *MessageCache) Bind(ctx context.Context, channel string, clientId int64, uid int) error {
	return m.Set(ctx, channel, strconv.FormatInt(clientId, 10), uid)
}

func (m *MessageCache) UnBind(ctx context.Context, channel string, clientId int64) error {
	return m.Del(ctx, channel, strconv.FormatInt(clientId, 10))
}

func (m *MessageCache) clientKey(sid, channel string) string {
	return fmt.Sprintf("ws:%s:%s:client", sid, channel)
}

func (m *MessageCache) userKey(sid, channel, uid string) string {
	return fmt.Sprintf("ws:%s:%s:user:%s", sid, channel, uid)
}

func NewMessageCache(cacheType *model.CacheType) *MessageCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := basicCache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &MessageModel{}
		})
		return &MessageCache{cache: c, redisClient: model.GetRedisCli(), config: config.Get(), storage: NewServerCache(model.GetCacheType())}
	case "memory":
		c := basicCache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &MessageModel{}
		})
		return &MessageCache{cache: c, redisClient: model.GetRedisCli(), config: config.Get(), storage: NewServerCache(model.GetCacheType())}
	}

	return nil // no cache
}
