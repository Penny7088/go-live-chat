package handler

import (
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

var _ GroupHandler = (*groupHandler)(nil)

type GroupHandler interface {
	Create(ctx *gin.Context)
	Dismiss(ctx *gin.Context)
	Invite(ctx *gin.Context)
	SignOut(ctx *gin.Context)
	Setting(ctx *gin.Context)
	RemoveMembers(ctx *gin.Context)
	Detail(ctx *gin.Context)
	UpdateMemberRemark(ctx *gin.Context)
	GetInviteFriends(ctx *gin.Context)
	GroupList(ctx *gin.Context)
	Members(ctx *gin.Context)
	OvertList(ctx *gin.Context)
	Overt(ctx *gin.Context)
	Handover(ctx *gin.Context)
	AssignAdmin(ctx *gin.Context)
	NoSpeak(ctx *gin.Context)
	Mute(ctx *gin.Context)
}

type groupHandler struct {
	groupDao        dao.GroupDao
	db              *gorm.DB
	talkRecordCache cache.TalkRecordsCache
	talkRecordsDao  dao.TalkRecordsDao
	redis           *redis.Client
	groupMemberDao  dao.GroupMemberDao
}

func (g groupHandler) Create(ctx *gin.Context) {
	var (
		members  []*model.GroupMember
		talkList []*model.TalkSession
	)
	params := &types.GroupCreateRequest{}
	if err := ctx.ShouldBindJSON(params); err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(ctx))
		response.Error(ctx, ecode.InvalidParams)
		return
	}

	uid, err2 := jwt.HeaderObtainUID(ctx)
	if err2 != nil {
		logger.Warn("uid obtain error: ", logger.Err(err2), middleware.GCtxRequestIDField(ctx))
		response.Error(ctx, ecode.InvalidParams)
		return
	}

	uids := sliceutil.Unique(append(sliceutil.ParseIds(params.IDs), uid))

	group := &model.Group{
		CreatorID: uid,
		Name:      params.Name,
		Avatar:    params.Avatar,
		MaxNum:    200,
	}

	joinTime := time.Now()

	err := g.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Create(group).Error
		if err != nil {
			return err
		}

		addMembers := make([]types.TalkRecordExtraGroupMembers, 0, len(uids))

		tx.Table("users").Select("id as user_id", "username").Where("id in ?", params.IDs).Scan(&addMembers)

		for _, val := range uids {
			leader := 0
			if val == uid {
				leader = 2
			}
			members = append(members, &model.GroupMember{
				GroupID:  uint(group.ID),
				UserID:   val,
				Leader:   leader,
				JoinTime: joinTime,
			})

			talkList = append(talkList, &model.TalkSession{
				TalkType:   2,
				UserID:     int64(val),
				ReceiverID: uint64(group.ID),
			})

		}

		if err = tx.Create(members).Error; err != nil {
			return err
		}

		if err = tx.Create(talkList).Error; err != nil {
			return err
		}

		var user model.Users
		err = tx.Table("users").Where("id = ?", uid).Scan(&user).Error
		if err != nil {
			return err
		}
		record := &model.TalkRecords{
			MsgID:      strutil.NewMsgId(),
			TalkType:   constant.ChatGroupMode,
			ReceiverID: int(group.ID),
			MsgType:    constant.ChatMsgSysGroupCreate,
			Sequence:   g.talkRecordCache.GetSequence(ctx, 0, int(group.ID)),
			Extra: jsonutil.Encode(types.TalkRecordExtraGroupCreate{
				OwnerId:   user.ID,
				OwnerName: user.Username,
				Members:   addMembers,
			}),
		}
		if err = tx.Create(record).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		logger.Warn("Transaction error: ", logger.Err(err2), middleware.GCtxRequestIDField(ctx))
		response.Error(ctx, ecode.ErrCreateGroup)
		return
	}

	// 广播网关将在线的用户加入房间
	body := map[string]any{
		"event": constant.SubEventGroupJoin,
		"data": jsonutil.Encode(map[string]any{
			"group_id": group.ID,
			"uids":     uids,
		}),
	}

	g.redis.Publish(ctx, constant.ImTopicChat, jsonutil.Encode(body))

	response.Success(ctx, gin.H{
		"group_id": &types.GroupCreateReply{
			GroupID: uint64(group.ID),
		},
	})
}

func (g groupHandler) Dismiss(c *gin.Context) {
	params := &types.GroupDismissRequest{}
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
	if !g.groupMemberDao.IsMaster(ctx, params.GroupID, uid) {
		logger.Warn("not permission dismiss ", logger.Err(err2), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.ErrGroupDismiss)
		return
	}

	err := g.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.Group{ID: params.GroupID, CreatorID: uid}).Updates(&model.Group{
			IsDismiss: 1,
		}).Error; err != nil {
			return err
		}

		if err := g.db.Model(&model.GroupMember{}).Where("group_id = ?", params.GroupID).Updates(&model.GroupMember{
			IsQuit: 1,
		}).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		logger.Warn("not permission dismiss ", logger.Err(err2), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.ErrGroupDismiss)
		return
	}

	response.Success(c, "ok")
}

func (g groupHandler) Invite(ctx *gin.Context) {
	// TODO implement me
	panic("implement me")
}

func (g groupHandler) SignOut(ctx *gin.Context) {
	// TODO implement me
	panic("implement me")
}

func (g groupHandler) Setting(ctx *gin.Context) {
	// TODO implement me
	panic("implement me")
}

func (g groupHandler) RemoveMembers(ctx *gin.Context) {
	// TODO implement me
	panic("implement me")
}

func (g groupHandler) Detail(ctx *gin.Context) {
	// TODO implement me
	panic("implement me")
}

func (g groupHandler) UpdateMemberRemark(ctx *gin.Context) {
	// TODO implement me
	panic("implement me")
}

func (g groupHandler) GetInviteFriends(ctx *gin.Context) {
	// TODO implement me
	panic("implement me")
}

func (g groupHandler) GroupList(ctx *gin.Context) {
	// TODO implement me
	panic("implement me")
}

func (g groupHandler) Members(ctx *gin.Context) {
	// TODO implement me
	panic("implement me")
}

func (g groupHandler) OvertList(ctx *gin.Context) {
	// TODO implement me
	panic("implement me")
}

func (g groupHandler) Overt(ctx *gin.Context) {
	// TODO implement me
	panic("implement me")
}

func (g groupHandler) Handover(ctx *gin.Context) {
	// TODO implement me
	panic("implement me")
}

func (g groupHandler) AssignAdmin(ctx *gin.Context) {
	// TODO implement me
	panic("implement me")
}

func (g groupHandler) NoSpeak(ctx *gin.Context) {
	// TODO implement me
	panic("implement me")
}

func (g groupHandler) Mute(ctx *gin.Context) {
	// TODO implement me
	panic("implement me")
}

func NewGroupHandler() GroupHandler {
	return &groupHandler{
		groupDao:        dao.NewGroupDao(model.GetDB(), cache.NewGroupCache(model.GetCacheType())),
		db:              model.GetDB(),
		talkRecordsDao:  dao.NewTalkRecordsDao(model.GetDB(), cache.NewTalkRecordsCache(model.GetCacheType())),
		talkRecordCache: cache.NewTalkRecordsCache(model.GetCacheType()),
		redis:           model.GetRedisCli(),
		groupMemberDao:  dao.NewGroupMemberDao(model.GetDB(), cache.NewGroupMemberCache(model.GetCacheType())),
	}
}
