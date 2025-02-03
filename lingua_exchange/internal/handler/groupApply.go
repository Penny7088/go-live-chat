package handler

import (
	"context"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/zhufuyi/sponge/pkg/ggorm/query"
	"github.com/zhufuyi/sponge/pkg/gin/middleware"
	"github.com/zhufuyi/sponge/pkg/gin/response"
	"github.com/zhufuyi/sponge/pkg/logger"
	"gorm.io/gorm"
	"lingua_exchange/internal/cache"
	"lingua_exchange/internal/constant"
	"lingua_exchange/internal/dao"
	"lingua_exchange/internal/ecode"
	"lingua_exchange/internal/model"
	"lingua_exchange/internal/types"
	"lingua_exchange/pkg/jsonutil"
	"lingua_exchange/pkg/jwt"
	"lingua_exchange/pkg/sliceutil"
	"lingua_exchange/pkg/timeutil"
)

var _ GroupApplyHandler = (*groupApplyHandler)(nil)

type GroupApplyHandler interface {
	Create(ctx *gin.Context)
	Agree(ctx *gin.Context)
	Decline(ctx *gin.Context)
	List(ctx *gin.Context)
	All(ctx *gin.Context)
	ApplyUnreadNum(ctx *gin.Context)
}

type groupApplyHandler struct {
	iDao           dao.GroupApplyDao
	iCache         cache.GroupApplyCache
	groupDao       dao.GroupDao
	groupMemberDao dao.GroupMemberDao
	redis          *redis.Client
}

func NewGroupApplyHandler() GroupApplyHandler {
	return &groupApplyHandler{
		iDao:           dao.NewGroupApplyDao(model.GetDB(), cache.NewGroupApplyCache(model.GetCacheType())),
		groupDao:       dao.NewGroupDao(model.GetDB(), cache.NewGroupCache(model.GetCacheType())),
		groupMemberDao: dao.NewGroupMemberDao(model.GetDB(), cache.NewGroupMemberCache(model.GetCacheType())),
		redis:          model.GetRedisCli(),
		iCache:         cache.NewGroupApplyCache(model.GetCacheType()),
	}
}

func (g groupApplyHandler) Create(c *gin.Context) {
	params := &types.GroupApplyCreateRequest{}
	if err := c.ShouldBindJSON(params); err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	ctx := middleware.WrapCtx(c)
	args := &query.Conditions{
		Columns: []query.Column{
			{
				Name:  "group_id",
				Value: params.GroupID,
			},
			{
				Name:  "status",
				Value: constant.GroupApplyStatusWait,
			},
		},
	}
	groupApply, err := g.iDao.GetByCondition(ctx, args)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Warn("group apply  error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.ErrGroupApplyCreateFailed)
		return
	}

	uid, err2 := jwt.HeaderObtainUID(c)
	if err2 != nil {
		logger.Warn("uid obtain error: ", logger.Err(err2), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	err = g.createOrUpdate(groupApply, ctx, params, uid)
	if err != nil {
		logger.Warn("group apply  error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.ErrGroupApplyCreateFailed)
		return
	}

	groupMember, err := g.groupMemberDao.FindByWhere(ctx, "group_id = ? and leader = ?", params.GroupID, 2)
	if err == nil && groupMember != nil {
		g.iCache.Incr(ctx, uint64(groupMember.UserID))
	}

	g.redis.Publish(ctx, constant.ImTopicChat, jsonutil.Encode(map[string]any{
		"event": constant.SubEventGroupApply,
		"data": jsonutil.Encode(map[string]any{
			"group_id": params.GroupID,
			"user_id":  uid,
		}),
	}))

	response.Success(c, "ok")
}

func (g groupApplyHandler) createOrUpdate(groupApply *model.GroupApply, ctx context.Context, params *types.GroupApplyCreateRequest, uid int) error {
	if groupApply == nil {
		err := g.iDao.Create(ctx, &model.GroupApply{
			GroupID: uint(params.GroupID),
			UserID:  uint(uid),
			Status:  constant.GroupApplyStatusWait,
			Remark:  params.Remark,
		})
		if err != nil {
			return err
		}
		return nil
	} else {
		groupApply.Remark = params.Remark
		groupApply.UpdatedAt = time.Now()
		err := g.iDao.UpdateByID(ctx, groupApply)
		if err != nil {
			return err
		}
		return nil
	}
}

func (g groupApplyHandler) Agree(ctx *gin.Context) {
	// TODO implement me
	panic("implement me")
}

func (g groupApplyHandler) Decline(c *gin.Context) {
	params := &types.GroupApplyDeclineRequest{}
	if err := c.ShouldBindJSON(params); err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	uid, err2 := jwt.HeaderObtainUID(c)
	if err2 != nil {
		logger.Warn("uid obtain error: ", logger.Err(err2), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	ctx := middleware.WrapCtx(c)
	apply, err := g.iDao.GetByID(ctx, uint64(params.ApplyID))
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Warn("apply not found  error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.ErrGroupApplyCreateFailed)
		return
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Warn("apply not found  error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.ErrGroupApplyNotFound)
		return
	}

	isLeader := g.groupMemberDao.IsLeader(ctx, int(apply.GroupID), uid)
	if !isLeader {
		logger.Warn("not permission ", logger.Err(err2), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.ErrGroupNotPermission)
		return
	}

	if apply.Status != constant.GroupApplyStatusWait {
		logger.Warn("申请信息已被他(她)人处理 ", logger.Err(err2), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.ErrGroupApplyAlreadyHandler)
		return
	}

	apply = &model.GroupApply{
		ID:     uint64(params.ApplyID),
		Status: constant.GroupApplyStatusRefuse,
		Reason: params.Remark,
	}

	err = g.iDao.UpdateByID(ctx, apply)
	if err != nil {
		logger.Warn("更新失败 ", logger.Err(err2), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.ErrGroupApplyUpdate)
		return
	}

	response.Success(c, "ok")
}

func (g groupApplyHandler) List(c *gin.Context) {
	params := &types.GroupApplyListRequest{}
	if err := c.ShouldBindJSON(params); err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	uid, err2 := jwt.HeaderObtainUID(c)
	if err2 != nil {
		logger.Warn("uid obtain error: ", logger.Err(err2), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	ctx := middleware.WrapCtx(c)
	isLeader := g.groupMemberDao.IsLeader(ctx, params.GroupID, uid)
	if !isLeader {
		logger.Warn("not permission ", logger.Err(err2), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.ErrGroupNotPermission)
		return
	}
	list, err := g.iDao.List(ctx, []uint{uint(params.GroupID)})
	if err != nil {
		logger.Warn("list err ", logger.Err(err2), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.ErrGroupApplyCreateFailed)
		return
	}

	items := make([]*types.GroupApplyItem, 0)
	for _, item := range list {
		items = append(items, &types.GroupApplyItem{
			ID:        item.Id,
			UserID:    item.UserId,
			GroupID:   item.GroupId,
			Remark:    item.Remark,
			Avatar:    item.Avatar,
			Username:  item.Nickname,
			CreatedAt: timeutil.FormatDatetime(item.CreatedAt),
		})
	}

	response.Success(c, gin.H{
		"list": items,
	})
}

func (g groupApplyHandler) All(c *gin.Context) {
	uid, err2 := jwt.HeaderObtainUID(c)
	if err2 != nil {
		logger.Warn("uid obtain error: ", logger.Err(err2), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	ctx := middleware.WrapCtx(c)
	all, err := g.groupMemberDao.FindAll(ctx, func(db *gorm.DB) {
		db.Select("group_id")
		db.Where("user_id = ?", uid)
		db.Where("leader = ?", 2)
		db.Where("is_quit = ?", 0)
	})

	if err != nil {
		logger.Warn("查询 系统异常，请稍后再试: ", logger.Err(err2), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	groupIds := make([]uint, 0, len(all))
	for _, m := range all {
		groupIds = append(groupIds, m.GroupID)
	}

	resp := &types.GroupApplyListResponse{Items: make([]types.GroupApplyItem, 0)}
	if len(groupIds) == 0 {
		response.Success(c, gin.H{
			"items": resp,
		})
		return
	}

	list, err := g.iDao.List(ctx, groupIds)
	if err != nil {
		logger.Warn("查询 系统异常，请稍后再试: ", logger.Err(err2), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	groups, err := g.groupDao.FindAll(ctx, func(db *gorm.DB) {
		db.Select("id,name")
		db.Where("id in ?", groupIds)
	})
	if err != nil {
		logger.Warn("查询 群组 系统异常，请稍后再试: ", logger.Err(err2), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	groupMap := sliceutil.ToMap(groups, func(t *model.Group) int {
		return t.ID
	})

	for _, item := range list {
		resp.Items = append(resp.Items, types.GroupApplyItem{
			ID:        item.Id,
			UserID:    item.UserId,
			GroupName: groupMap[item.GroupId].Name,
			GroupID:   item.GroupId,
			Remark:    item.Remark,
			Avatar:    item.Avatar,
			Username:  item.Nickname,
			CreatedAt: timeutil.FormatDatetime(item.CreatedAt),
		})
	}

	g.iCache.Del(ctx, uint64(uid))
	response.Success(c, gin.H{
		"items": resp,
	})
}

func (g groupApplyHandler) ApplyUnreadNum(c *gin.Context) {
	uid, err2 := jwt.HeaderObtainUID(c)
	if err2 != nil {
		logger.Warn("uid obtain error: ", logger.Err(err2), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	ctx := middleware.WrapCtx(c)
	count := g.iCache.GetCount(ctx, uint64(uid))
	response.Success(c, gin.H{
		"unread_num": count,
	})
}
