package cache

import (
	"context"
	"fmt"
	"strconv"

	"github.com/redis/go-redis/v9"
	"lingua_exchange/internal/model"
)

var _ UnreadCache = (*unreadCache)(nil)

type UnreadCache interface {
	Incr(ctx context.Context, mode, sender, receive int)
	PipeIncr(ctx context.Context, pipe redis.Pipeliner, mode, sender, receive int)
	Get(ctx context.Context, mode, sender, receive int) int
	Del(ctx context.Context, mode, sender, receive int)
	Reset(ctx context.Context, mode, sender, receive int)
	All(ctx context.Context, receive int) map[string]int
}

type unreadCache struct {
	redis *redis.Client
}

func NewUnreadCache() UnreadCache {

	return &unreadCache{redis: model.GetRedisCli()}
}

// Incr 消息未读数自增
// @params mode    对话模式 1私信 2群聊
// @params sender  发送者ID
// @params receive 接收者ID
func (u unreadCache) Incr(ctx context.Context, mode, sender, receive int) {
	u.redis.HIncrBy(ctx, u.name(receive), fmt.Sprintf("%d_%d", mode, sender), 1)
}

// PipeIncr 消息未读数自增
// @params mode    对话模式 1私信 2群聊
// @params sender  发送者ID
// @params receive 接收者ID
func (u unreadCache) PipeIncr(ctx context.Context, pipe redis.Pipeliner, mode, sender, receive int) {
	pipe.HIncrBy(ctx, u.name(receive), fmt.Sprintf("%d_%d", mode, sender), 1)
}

// Get 获取消息未读数
// @params mode    对话模式 1私信 2群聊
// @params sender  发送者ID
// @params receive 接收者ID
func (u unreadCache) Get(ctx context.Context, mode, sender, receive int) int {
	val, _ := u.redis.HGet(ctx, u.name(receive), fmt.Sprintf("%d_%d", mode, sender)).Int()
	return val
}

// Del 删除消息未读数
// @params mode    对话模式 1私信 2群聊
// @params sender  发送者ID
// @params receive 接收者ID
func (u unreadCache) Del(ctx context.Context, mode, sender, receive int) {
	u.redis.HDel(ctx, u.name(receive), fmt.Sprintf("%d_%d", mode, sender))
}

// Reset 消息未读数重置
// @params mode    对话模式 1私信 2群聊
// @params sender  发送者ID
// @params receive 接收者ID
func (u unreadCache) Reset(ctx context.Context, mode, sender, receive int) {
	u.redis.HSet(ctx, u.name(receive), fmt.Sprintf("%d_%d", mode, sender), 0)
}

// All 获取获取所有的未读数
// @params receive 接收者ID
func (u unreadCache) All(ctx context.Context, receive int) map[string]int {
	items := make(map[string]int)
	for k, v := range u.redis.HGetAll(ctx, u.name(receive)).Val() {
		items[k], _ = strconv.Atoi(v)
	}

	return items
}

func (u *unreadCache) name(receive int) string {
	return fmt.Sprintf("im:message:unread:uid_%d", receive)
}
