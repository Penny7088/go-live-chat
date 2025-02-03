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
	"lingua_exchange/pkg/strutil"
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
	iDao             dao.GroupApplyDao
	iCache           cache.GroupApplyCache
	groupDao         dao.GroupDao
	groupMemberDao   dao.GroupMemberDao
	redis            *redis.Client
	db               *gorm.DB
	talkRecordsCache cache.TalkRecordsCache
}

func NewGroupApplyHandler() GroupApplyHandler {
	return &groupApplyHandler{
		iDao:             dao.NewGroupApplyDao(model.GetDB(), cache.NewGroupApplyCache(model.GetCacheType())),
		groupDao:         dao.NewGroupDao(model.GetDB(), cache.NewGroupCache(model.GetCacheType())),
		groupMemberDao:   dao.NewGroupMemberDao(model.GetDB(), cache.NewGroupMemberCache(model.GetCacheType())),
		redis:            model.GetRedisCli(),
		iCache:           cache.NewGroupApplyCache(model.GetCacheType()),
		db:               model.GetDB(),
		talkRecordsCache: cache.NewTalkRecordsCache(model.GetCacheType()),
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

func (g groupApplyHandler) Agree(c *gin.Context) {
	params := &types.GroupApplyAgreeRequest{}
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

	if !g.groupMemberDao.IsMember(ctx, int(apply.GroupID), int(apply.UserID), false) {
		err := g.invite(int(apply.GroupID), []int{int(apply.UserID)}, uid, ctx)
		if err != nil {
			logger.Warn("邀请失败 ", logger.Err(err2), middleware.GCtxRequestIDField(c))
			response.Error(c, ecode.ErrGroupApplyAlreadyHandler)
			return
		}
	}

	err = g.iDao.UpdateByID(ctx, &model.GroupApply{
		ID:     uint64(params.ApplyID),
		Status: constant.GroupApplyStatusPass,
	})

	if err != nil {
		logger.Warn("更新失败 ", logger.Err(err2), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.ErrGroupApplyAlreadyHandler)
		return
	}

	response.Success(c, "ok")

}

// invite 邀请用户加入群聊
func (g groupApplyHandler) invite(gid int, ids []int, uid int, ctx context.Context) error {
	var (
		err            error
		addMembers     []*model.GroupMember
		addTalkList    []*model.TalkSession
		updateTalkList []uint64
		talkList       []*model.TalkSession
		db             = g.db
	)
	m := make(map[int]struct{})
	for _, value := range g.groupMemberDao.GetMemberIds(ctx, gid) {
		m[value] = struct{}{}
	}

	listHash := make(map[int64]*model.TalkSession)
	db.Select("id", "user_id", "is_delete").Where("user_id in ? and receiver_id = ? and talk_type = 2", ids, gid).Find(&talkList)
	for _, item := range talkList {
		listHash[item.UserID] = item
	}
	mids := make([]int, 0)
	mids = append(mids, ids...)
	mids = append(mids, uid)

	memberItems := make([]*model.Users, 0)
	err = db.Table("users").Select("id,username").Where("id in ?", mids).Scan(&memberItems).Error
	if err != nil {
		return err
	}
	memberMaps := make(map[uint64]*model.Users)
	for _, item := range memberItems {
		memberMaps[item.ID] = item
	}
	members := make([]types.TalkRecordExtraGroupMembers, 0)

	for _, value := range ids {
		members = append(members, types.TalkRecordExtraGroupMembers{
			UserId:   value,
			Username: memberMaps[uint64(value)].Username,
		})

		if _, ok := m[value]; !ok {
			addMembers = append(addMembers, &model.GroupMember{
				GroupID:  uint(gid),
				UserID:   value,
				JoinTime: time.Now(),
			})
		}

		if item, ok := listHash[int64(value)]; !ok {
			addTalkList = append(addTalkList, &model.TalkSession{
				TalkType:   constant.ChatGroupMode,
				UserID:     int64(value),
				ReceiverID: uint64(gid),
			})
		} else if item.IsDelete == 1 {
			updateTalkList = append(
				updateTalkList,
				item.ID,
			)
		}

	}
	if len(addMembers) == 0 {
		return errors.New("邀请的好友，都已成为群成员")
	}

	record := &model.TalkRecords{
		MsgID:      strutil.NewMsgId(),
		TalkType:   constant.ChatGroupMode,
		ReceiverID: gid,
		MsgType:    constant.ChatMsgSysGroupMemberJoin,
		Sequence:   g.talkRecordsCache.GetSequence(ctx, 0, gid),
	}

	record.Extra = jsonutil.Encode(&types.TalkRecordExtraGroupJoin{
		OwnerId:   memberMaps[uint64(uid)].ID,
		OwnerName: memberMaps[uint64(uid)].Username,
		Members:   members,
	})

	err = db.Transaction(func(tx *gorm.DB) error {
		tx.Delete(&model.GroupMember{}, "group_id = ? and user_id in ? and is_quit = ?", gid, ids, constant.GroupMemberQuitStatusYes)

		if err = tx.Create(&addMembers).Error; err != nil {
			return err
		}
		// 添加用户的对话列表
		if len(addTalkList) > 0 {
			if err = tx.Create(&addTalkList).Error; err != nil {
				return err
			}
		}

		// 更新用户的对话列表
		if len(updateTalkList) > 0 {
			tx.Model(&model.TalkSession{}).Where("id in ?", updateTalkList).Updates(map[string]any{
				"is_delete":  0,
				"created_at": timeutil.DateTime(),
			})
		}

		if err = tx.Create(record).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	// 广播网关将在线的用户加入房间
	g.redis.Publish(ctx, constant.ImTopicChat, jsonutil.Encode(map[string]any{
		"event": constant.SubEventGroupJoin,
		"data": jsonutil.Encode(map[string]any{
			"type":     1,
			"group_id": gid,
			"uids":     ids,
		}),
	}))

	g.redis.Publish(ctx, constant.ImTopicChat, jsonutil.Encode(map[string]any{
		"event": constant.SubEventImMessage,
		"data": jsonutil.Encode(map[string]any{
			"sender_id":   record.UserID,
			"receiver_id": record.ReceiverID,
			"talk_type":   record.TalkType,
			"msg_id":      record.MsgID,
		}),
	}))

	return nil
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
