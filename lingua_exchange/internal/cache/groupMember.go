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

var _ GroupMemberCache = (*groupMemberCache)(nil)
var (
	KeyGroupRelation = "im:contact:relation:"
)

type GroupMemberCache interface {
	IsGroupRelation(ctx context.Context, uid, gid int) error
	SetGroupRelation(ctx context.Context, uid, gid int) error
	DelGroupRelation(ctx context.Context, uid, gid int) error
	BatchDelGroupRelation(ctx context.Context, uids []int, gid int) error
}

type groupMemberCache struct {
	cache cache.Cache
}

func (g groupMemberCache) SetGroupRelation(ctx context.Context, uid, gid int) error {
	relationKey := keyGroupRelation(uid, gid)
	err := g.cache.Set(ctx, relationKey, "1", time.Hour*1)
	if err != nil {
		return err
	}
	return nil
}

func (g groupMemberCache) DelGroupRelation(ctx context.Context, uid, gid int) error {
	key := keyGroupRelation(uid, gid)
	err := g.cache.Del(ctx, key)
	if err != nil {
		return err
	}
	return nil
}

func (g groupMemberCache) BatchDelGroupRelation(ctx context.Context, uids []int, gid int) error {
	for _, uid := range uids {
		err := g.DelGroupRelation(ctx, uid, gid)
		if err != nil {
			return err
		}
	}
	return nil
}

func (g groupMemberCache) IsGroupRelation(ctx context.Context, uid, gid int) error {
	relationKey := keyGroupRelation(uid, gid)
	return g.cache.Get(ctx, relationKey, gid)
}

func NewGroupMemberCache(cacheType *model.CacheType) GroupMemberCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.GroupMember{}
		})
		return &groupMemberCache{cache: c}
	case "memory":
		c := cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.GroupMember{}
		})
		return &groupMemberCache{cache: c}
	}

	return nil
}

func keyGroupRelation(uid, gid int) string {
	return KeyGroupRelation + utils.IntToStr(uid) + utils.IntToStr(gid)
}
