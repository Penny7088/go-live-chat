package handler

import (
	"context"
	"errors"
	"fmt"
	"time"

	"lingua_exchange/internal/cache"
	"lingua_exchange/internal/constant"
	"lingua_exchange/internal/dao"
	"lingua_exchange/internal/ecode"
	"lingua_exchange/internal/imService"
	"lingua_exchange/internal/model"
	"lingua_exchange/internal/types"
	"lingua_exchange/pkg/jsonutil"
	"lingua_exchange/pkg/jwt"
	"lingua_exchange/pkg/sliceutil"
	"lingua_exchange/pkg/strutil"
	"lingua_exchange/pkg/timeutil"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/zhufuyi/sponge/pkg/gin/middleware"
	"github.com/zhufuyi/sponge/pkg/gin/response"
	"github.com/zhufuyi/sponge/pkg/logger"
	"gorm.io/gorm"
)

var _ GroupHandler = (*groupHandler)(nil)

type GroupHandler interface {
	Create(ctx *gin.Context)
	Dismiss(ctx *gin.Context)
	Invite(ctx *gin.Context)
	SignOut(ctx *gin.Context)
	Setting(ctx *gin.Context)
	Detail(ctx *gin.Context)
	GroupList(ctx *gin.Context)
	Handover(ctx *gin.Context)
	AssignAdmin(ctx *gin.Context)
	NoSpeak(ctx *gin.Context)
	Mute(ctx *gin.Context)
}

type groupHandler struct {
	groupDao         dao.GroupDao
	db               *gorm.DB
	talkRecordCache  cache.TalkRecordsCache
	talkRecordsDao   dao.TalkRecordsDao
	redis            *redis.Client
	groupMemberDao   dao.GroupMemberDao
	groupMemberCache cache.GroupMemberCache
	redisLock        *cache.RedisLock
	talkSessionDao   dao.TalkSessionDao
	messageService   imService.IMessageService
	usersDao         dao.UsersDao
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

func (g groupHandler) Invite(c *gin.Context) {
	params := &types.GroupInviteRequest{}
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
	key := fmt.Sprintf("group-join:%d", params.GroupID)
	ctx := middleware.WrapCtx(c)
	if !g.redisLock.Lock(ctx, key, 20) {
		logger.Warn("uid obtain error: ", logger.Err(err2), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.Timeout)
		return
	}
	defer g.redisLock.UnLock(ctx, key)

	group, err := g.groupDao.GetByID(ctx, uint64(params.GroupID))
	if err != nil {
		logger.Warn("query group error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.Timeout)
		return
	}

	if group != nil && group.IsDismiss == 1 {
		logger.Warn("群组已解散: ", middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.ErrGroupAlreadyDismiss)
		return
	}

	uids := sliceutil.Unique(sliceutil.ParseIds(params.IDs))
	if len(uids) == 0 {
		logger.Warn("邀请好友列表为空: ", middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.ErrGroupInviteFriendsNil)
		return
	}

	if !g.groupMemberDao.IsMember(ctx, int(params.GroupID), uid, true) {
		logger.Warn("非群组成员，无权邀请好友！: ", middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.ErrGroupInviteNotPermission)
		return
	}

	err = g.invite(params, uids, uid, ctx)
	if err != nil {
		logger.Warn("邀请失败: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.ErrGroupInviteFailed)
		return
	}

	response.Success(c, "ok")
}

// invite 邀请用户加入群聊
func (g groupHandler) invite(prams *types.GroupInviteRequest, ids []int, uid int, ctx context.Context) error {
	var (
		err            error
		addMembers     []*model.GroupMember
		addTalkList    []*model.TalkSession
		updateTalkList []uint64
		talkList       []*model.TalkSession
		db             = g.db
	)
	m := make(map[int]struct{})
	for _, value := range g.groupMemberDao.GetMemberIds(ctx, int(prams.GroupID)) {
		m[value] = struct{}{}
	}

	listHash := make(map[int64]*model.TalkSession)
	db.Select("id", "user_id", "is_delete").Where("user_id in ? and receiver_id = ? and talk_type = 2", ids, prams.GroupID).Find(&talkList)
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
				GroupID:  uint(prams.GroupID),
				UserID:   value,
				JoinTime: time.Now(),
			})
		}

		if item, ok := listHash[int64(value)]; !ok {
			addTalkList = append(addTalkList, &model.TalkSession{
				TalkType:   constant.ChatGroupMode,
				UserID:     int64(value),
				ReceiverID: uint64(prams.GroupID),
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
		ReceiverID: int(prams.GroupID),
		MsgType:    constant.ChatMsgSysGroupMemberJoin,
		Sequence:   g.talkRecordCache.GetSequence(ctx, 0, int(prams.GroupID)),
	}

	record.Extra = jsonutil.Encode(&types.TalkRecordExtraGroupJoin{
		OwnerId:   memberMaps[uint64(uid)].ID,
		OwnerName: memberMaps[uint64(uid)].Username,
		Members:   members,
	})

	err = db.Transaction(func(tx *gorm.DB) error {
		tx.Delete(&model.GroupMember{}, "group_id = ? and user_id in ? and is_quit = ?", prams.GroupID, ids, constant.GroupMemberQuitStatusYes)

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
			"group_id": prams.GroupID,
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

func (g groupHandler) SignOut(c *gin.Context) {
	params := &types.GroupOutRequest{}
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
	sid := g.talkSessionDao.FindBySessionId(ctx, uid, int(params.GroupID), constant.ChatGroupMode)
	err := g.talkSessionDao.Delete(ctx, uid, sid)
	if err != nil {
		logger.Warn("delete  error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
	}

	response.Success(c, "ok")
}

func (g groupHandler) Setting(c *gin.Context) {
	params := &types.GroupSettingRequest{}
	if err := c.ShouldBindJSON(params); err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	ctx := middleware.WrapCtx(c)
	uid, err2 := jwt.HeaderObtainUID(c)
	if err2 != nil {
		logger.Warn("uid obtain error: ", logger.Err(err2), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	group, err := g.groupDao.GetByID(ctx, uint64(params.GroupID))
	if err != nil {
		logger.Warn("query group error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.ErrGroupNotExist)
		return
	}

	if group != nil && group.IsDismiss == 1 {
		logger.Warn("群组已解散: ", middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.ErrGroupAlreadyDismiss)
		return
	}

	if !g.groupMemberDao.IsLeader(ctx, int(params.GroupID), uid) {
		logger.Warn("非群主，无权修改群信息", middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.ErrGroupNotPermission)
		return
	}

	p := &model.Group{
		ID:      int(params.GroupID),
		Avatar:  params.Avatar,
		Name:    params.GroupName,
		Profile: params.Profile,
	}
	err2 = g.updateGroup(ctx, p)
	if err2 != nil {
		logger.Warn("update group error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.ErrGroupNotExist)
		return
	}

	g.messageService.SendSystemText(&ctx, uid, &types.TextMessageRequest{
		Content: "群主或管理员修改了群信息！",
		Receiver: types.MessageReceiver{
			TalkType:   uint(constant.ChatGroupMode),
			ReceiverID: int(params.GroupID),
		},
	})

	response.Success(c, "ok")
}

func (g groupHandler) updateGroup(ctx context.Context, params *model.Group) error {
	err := g.groupDao.UpdateByID(ctx, params)
	if err != nil {
		return err
	}
	return nil
}

func (g groupHandler) Detail(c *gin.Context) {
	params := &types.GroupDetailsRequest{}
	if err := c.ShouldBindQuery(params); err != nil {
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

	groupInfo, err := g.groupDao.GetByID(ctx, uint64(params.GroupID))
	if err != nil {
		logger.Warn("group details  obtain error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.ErrGroupDetailsFailed)
		return
	}

	if groupInfo.ID == 0 {
		logger.Warn("group details  obtain error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.ErrGroupNotExist)
		return
	}

	resp := &types.GroupDetailResponse{
		GroupID:   groupInfo.ID,
		GroupName: groupInfo.Name,
		Profile:   groupInfo.Profile,
		Avatar:    groupInfo.Avatar,
		CreatedAt: timeutil.FormatDatetime(groupInfo.CreatedAt),
		IsManager: uid == groupInfo.CreatorID,
		IsDisturb: 0,
		IsMute:    int32(groupInfo.IsMute),
		IsOvert:   int32(groupInfo.IsOvert),
		VisitCard: g.groupMemberDao.GetMemberRemark(ctx, int(params.GroupID), uid),
	}

	if g.talkSessionDao.IsDisturb(ctx, uid, groupInfo.ID, 2) {
		resp.IsDisturb = 1
	}

	response.Success(c, resp)
}

func (g groupHandler) GroupList(c *gin.Context) {
	uid, err := jwt.HeaderObtainUID(c)
	if err != nil {
		logger.Warn("uid obtain error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	ctx := middleware.WrapCtx(c)
	items, err := g.groupDao.GroupList(ctx, uid)

	resp := &types.GroupListResponse{
		Items: make([]types.GroupItem, 0, len(items)),
	}

	for _, item := range items {
		gi := types.GroupItem{
			ID:        item.ID,
			GroupName: item.GroupName,
			Avatar:    item.Avatar,
			Profile:   item.Profile,
			Leader:    item.Leader,
			IsDisturb: item.IsDisturb,
			CreatorID: item.CreatorID,
		}
		resp.Items = append(resp.Items, gi)
	}

	response.Success(c, resp)
}

func (g groupHandler) Handover(c *gin.Context) {
	params := &types.GroupHandoverRequest{}
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
		logger.Warn("not permission ", logger.Err(err2), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.ErrGroupNotPermission)
		return
	}

	if uid == params.UserID {
		logger.Warn("not permission ", logger.Err(err2), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.ErrGroupNotPermission)
		return
	}

	err := g.groupMemberDao.Handover(ctx, params.GroupID, uid, params.UserID)
	if err != nil {
		logger.Warn("转让失败", logger.Err(err2), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.ErrGroupHandoverFailed)
		return
	}

	members := make([]types.TalkRecordExtraGroupMembers, 0)
	g.db.Table("users").Select("id as user_id", "username").Where("id in ?", []int{uid, params.UserID}).Scan(&members)
	extra := types.TalkRecordExtraGroupTransfer{}
	for _, member := range members {
		if member.UserId == uid {
			extra.OldOwnerId = member.UserId
			extra.OldOwnerName = member.Username
		} else {
			extra.NewOwnerId = member.UserId
			extra.NewOwnerName = member.Username
		}
	}
	_ = g.messageService.SendSysOther(&ctx, &model.TalkRecords{
		MsgType:    constant.ChatMsgSysGroupTransfer,
		TalkType:   constant.TalkRecordTalkTypeGroup,
		UserID:     uid,
		ReceiverID: params.GroupID,
		Extra:      jsonutil.Encode(extra),
	})

	response.Success(c, "ok")
}

func (g groupHandler) AssignAdmin(c *gin.Context) {
	params := &types.GroupAssignAdminRequest{}
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
		logger.Warn("not permission ", logger.Err(err2), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.ErrGroupNotPermission)
		return
	}

	leader := 0
	if params.Mode == 1 {
		leader = 1
	}
	err := g.groupMemberDao.SetLeaderStatus(ctx, params.GroupID, params.UserID, leader)
	if err != nil {
		logger.Warn("[Group AssignAdmin] 设置管理员信息失败 err :%s", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.ErrGroupAssignAdmin)
		return
	}
	response.Success(c, "ok")
}

func (g groupHandler) NoSpeak(c *gin.Context) {
	params := &types.GroupNoSpeakRequest{}
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

	if !g.groupMemberDao.IsLeader(ctx, params.GroupID, uid) {
		logger.Warn("not permission ", logger.Err(err2), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.ErrGroupNotPermission)
		return
	}
	status := 1
	if params.Mode == 2 {
		status = 0
	}
	err := g.groupMemberDao.SetMuteStatus(ctx, params.GroupID, params.UserID, status)
	if err != nil {
		logger.Warn("设置群成员禁言失败 ", logger.Err(err2), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.ErrGroupBanSpeakFailed)
		return
	}

	data := &model.TalkRecords{
		TalkType:   constant.TalkRecordTalkTypeGroup,
		UserID:     uid,
		ReceiverID: params.GroupID,
	}
	members := make([]types.TalkRecordExtraGroupMembers, 0)
	g.db.Table("users").Select("id as user_id", "username").Where("id = ?", params.UserID).Scan(&members)

	user, err2 := g.usersDao.GetByID(ctx, uint64(uid))
	if err2 != nil {
		logger.Warn("获取当前用户信息失败 ", logger.Err(err2), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.ErrGroupBanSpeakFailed)
		return
	}

	if status == 1 {
		data.MsgType = constant.ChatMsgSysGroupMemberMuted
		data.Extra = jsonutil.Encode(types.TalkRecordExtraGroupMemberCancelMuted{
			OwnerId:   uid,
			OwnerName: user.Username,
			Members:   members,
		})
	} else {
		data.MsgType = constant.ChatMsgSysGroupMemberCancelMuted
		data.Extra = jsonutil.Encode(types.TalkRecordExtraGroupMemberCancelMuted{
			OwnerId:   uid,
			OwnerName: user.Username,
			Members:   members,
		})
	}
	_ = g.messageService.SendSysOther(&ctx, data)

	response.Success(c, "ok")
}

func (g groupHandler) Mute(c *gin.Context) {
	params := &types.GroupMuteRequest{}
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
	group, err := g.groupDao.GetByID(ctx, uint64(params.GroupID))
	if err != nil {
		logger.Warn("获取当前群失败: ", logger.Err(err2), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.ErrGroupMuteFailed)
		return
	}

	if group.IsDismiss == 1 {
		logger.Warn("当前群已解散: ", middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.ErrGroupAlreadyDismiss)
		return
	}

	if !g.groupMemberDao.IsLeader(ctx, params.GroupID, uid) {
		logger.Warn("not permission ", logger.Err(err2), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.ErrGroupNotPermission)
		return
	}

	data := &model.Group{}
	data.ID = params.GroupID
	if params.Mode == 1 {
		data.IsMute = 1
	} else {
		data.IsMute = 0
	}

	err = g.groupDao.UpdateByID(ctx, data)
	if err != nil {
		logger.Warn("更新群信息失败: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.ErrGroupMuteFailed)
		return
	}

	user, err := g.usersDao.GetByID(ctx, uint64(uid))
	if err != nil {
		logger.Warn("获取当前用户失败: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.ErrGroupMuteFailed)
		return
	}

	var (
		extra   any
		msgType int
	)
	if params.Mode == 1 {
		msgType = constant.ChatMsgSysGroupMuted
		extra = types.TalkRecordExtraGroupMuted{
			OwnerId:   int(user.ID),
			OwnerName: user.Username,
		}
	} else {
		msgType = constant.ChatMsgSysGroupCancelMuted
		extra = types.TalkRecordExtraGroupCancelMuted{
			OwnerId:   int(user.ID),
			OwnerName: user.Username,
		}
	}

	_ = g.messageService.SendSysOther(&ctx, &model.TalkRecords{
		MsgType:    uint(msgType),
		TalkType:   constant.TalkRecordTalkTypeGroup,
		UserID:     uid,
		ReceiverID: params.GroupID,
		Extra:      jsonutil.Encode(extra),
	})

	response.Success(c, "ok")
}

func NewGroupHandler() GroupHandler {
	return &groupHandler{
		groupDao:         dao.NewGroupDao(model.GetDB(), cache.NewGroupCache(model.GetCacheType())),
		db:               model.GetDB(),
		talkRecordsDao:   dao.NewTalkRecordsDao(model.GetDB(), cache.NewTalkRecordsCache(model.GetCacheType())),
		talkRecordCache:  cache.NewTalkRecordsCache(model.GetCacheType()),
		redis:            model.GetRedisCli(),
		groupMemberDao:   dao.NewGroupMemberDao(model.GetDB(), cache.NewGroupMemberCache(model.GetCacheType())),
		talkSessionDao:   dao.NewTalkSessionDao(model.GetDB()),
		groupMemberCache: cache.NewGroupMemberCache(model.GetCacheType()),
		usersDao:         dao.NewUsersDao(model.GetDB(), cache.NewUsersCache(model.GetCacheType())),
	}
}
