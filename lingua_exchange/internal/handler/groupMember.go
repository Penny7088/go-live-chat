package handler

import (
	"context"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
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
)

var _ GroupMemberHandler = (*groupMemberHandler)(nil)

type GroupMemberHandler interface {
	Members(c *gin.Context)
	GetInviteFriends(c *gin.Context)
	RemoveMembers(c *gin.Context)
	UpdateMemberRemark(c *gin.Context)
}

type groupMemberHandler struct {
	groupMemberDao   dao.GroupMemberDao
	groupMemberCache cache.GroupMemberCache
	groupDao         dao.GroupDao
	db               *gorm.DB
	redis            *redis.Client
	talkRecordCache  cache.TalkRecordsCache
}

func NewGroupMemberHandler() GroupMemberHandler {
	return &groupMemberHandler{
		groupMemberDao:   dao.NewGroupMemberDao(model.GetDB(), cache.NewGroupMemberCache(model.GetCacheType())),
		groupMemberCache: cache.NewGroupMemberCache(model.GetCacheType()),
		groupDao:         dao.NewGroupDao(model.GetDB(), cache.NewGroupCache(model.GetCacheType())),
		db:               model.GetDB(),
		redis:            model.GetRedisCli(),
		talkRecordCache:  cache.NewTalkRecordsCache(model.GetCacheType()),
	}
}

func (g groupMemberHandler) Members(c *gin.Context) {
	params := &types.GroupMemberListRequest{}
	if err := c.ShouldBindJSON(params); err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	uid, err := jwt.HeaderObtainUID(c)
	if err != nil {
		logger.Warn("uid obtain error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	group, err := g.groupDao.GetByID(ctx, uint64(params.GroupID))
	if err != nil {
		logger.Warn("not found group  error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.ErrGroupNotExist)
		return
	}

	if group != nil && group.IsDismiss == 1 {
		logger.Warn("group is dismiss error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.ErrGroupDismiss)
		return
	}

	if !g.groupMemberDao.IsMember(ctx, params.GroupID, uid, false) {
		logger.Warn("not permission check group member: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.ErrGroupNotPermission)
		return
	}
	list := g.groupMemberDao.GetMembers(ctx, params.GroupID)
	items := make([]*types.GroupMemberItem, 0)
	for _, item := range list {
		items = append(items, &types.GroupMemberItem{
			UserID:   int32(item.UserId),
			Username: item.Nickname,
			Avatar:   item.Avatar,
			Gender:   int32(item.Gender),
			Leader:   int32(item.Leader),
			IsMute:   int32(item.IsMute),
			Remark:   item.UserCard,
		})
	}

	response.Success(c, items)
}

func (g groupMemberHandler) GetInviteFriends(c *gin.Context) {

}

func (g groupMemberHandler) RemoveMembers(c *gin.Context) {
	params := &types.GroupRemoveMemberRequest{}
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

	if !g.groupMemberDao.IsLeader(ctx, int(params.GroupID), uid) {
		logger.Warn("非群主，无权移除群成员", middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.ErrGroupNotPermission)
		return
	}

	g.removeMember(ctx, params, uid)

}

func (g groupMemberHandler) removeMember(ctx context.Context, params *types.GroupRemoveMemberRequest, uid int) error {
	var num int64
	membersIDs := sliceutil.ParseIds(params.MembersIDs)

	if err := g.db.Model(&model.GroupMember{}).Where("group_id =? and user_id =? and is_quit =0", params.GroupID, uid, constant.GroupMemberQuitStatusNo).Count(&num).Error; err != nil {
		return err
	}
	if int(num) != len(membersIDs) {
		return errors.New("移除成员失败")
	}

	mids := make([]int, 0)
	mids = append(mids, membersIDs...)
	mids = append(mids, uid)

	memberItems := make([]*model.Users, 0)
	err := g.db.Table("users").Select("id,username").Select("id in?", mids).Scan(&memberItems).Error
	if err != nil {
		return err
	}
	memberMaps := make(map[int]*model.Users)
	for _, item := range memberItems {
		memberMaps[int(item.ID)] = item
	}

	members := make([]types.TalkRecordExtraGroupMembers, 0)
	for _, value := range membersIDs {
		members = append(members, types.TalkRecordExtraGroupMembers{
			UserId:   value,
			Username: memberMaps[value].Username,
		})
	}

	record := &model.TalkRecords{
		MsgID:      strutil.NewMsgId(),
		Sequence:   g.talkRecordCache.GetSequence(ctx, 0, int(params.GroupID)),
		TalkType:   constant.ChatGroupMode,
		ReceiverID: int(params.GroupID),
		MsgType:    constant.ChatMsgSysGroupMemberKicked,
		Extra: jsonutil.Encode(&types.TalkRecordExtraGroupMemberKicked{
			OwnerId:   int(memberMaps[uid].ID),
			OwnerName: memberMaps[uid].Username,
			Members:   members,
		}),
	}

	err = g.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&model.GroupMember{}).Where("group_id = ? and user_id in ? and is_quit = 0", params.GroupID, membersIDs).Updates(map[string]any{
			"is_quit":    1,
			"updated_at": time.Now(),
		}).Error
		if err != nil {
			return err
		}

		return tx.Create(record).Error
	})

	if err != nil {
		return err
	}

	g.groupMemberCache.BatchDelGroupRelation(ctx, membersIDs, int(params.GroupID))

	// 广播网关将在线的用户加入房间
	_, _ = g.redis.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.Publish(ctx, constant.ImTopicChat, jsonutil.Encode(map[string]any{
			"event": constant.SubEventGroupJoin,
			"data": jsonutil.Encode(map[string]any{
				"type":     2,
				"group_id": params.GroupID,
				"uids":     membersIDs,
			}),
		}))

		pipe.Publish(ctx, constant.ImTopicChat, jsonutil.Encode(map[string]any{
			"event": constant.SubEventImMessage,
			"data": jsonutil.Encode(map[string]any{
				"sender_id":   int64(record.UserID),
				"receiver_id": int64(record.ReceiverID),
				"talk_type":   record.TalkType,
				"msg_id":      record.MsgID,
			}),
		}))
		return nil
	})

	return nil

}

func (g groupMemberHandler) UpdateMemberRemark(c *gin.Context) {
	params := &types.GroupRemarkUpdateRequest{}
	if err := c.ShouldBindJSON(params); err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	uid, err := jwt.HeaderObtainUID(c)
	if err != nil {
		logger.Warn("uid obtain error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	ctx := middleware.WrapCtx(c)
	_, err = g.groupMemberDao.UpdateWhere(ctx, map[string]any{
		"user_card": params.VisitCard,
	}, "group_id = ? and user_id = ?", params.GroupID, uid)
	if err != nil {
		logger.Warn("UpdateMemberRemark error: ", logger.Err(err))
		response.Error(c, ecode.ErrGroupUpdateGroupMemberRemark)
		return
	}

	response.Success(c, "ok")
}
