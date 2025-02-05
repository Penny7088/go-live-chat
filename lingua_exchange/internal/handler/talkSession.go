package handler

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zhufuyi/sponge/pkg/gin/middleware"
	"github.com/zhufuyi/sponge/pkg/gin/response"
	"github.com/zhufuyi/sponge/pkg/logger"
	"gorm.io/gorm"
	"lingua_exchange/internal/cache"
	"lingua_exchange/internal/constant"
	"lingua_exchange/internal/dao"
	"lingua_exchange/internal/ecode"
	"lingua_exchange/internal/imService"
	"lingua_exchange/internal/model"
	"lingua_exchange/internal/types"
	"lingua_exchange/pkg/jwt"
	"lingua_exchange/pkg/strutil"
	"lingua_exchange/pkg/timeutil"
)

var _ SessionHandler = (*sessionHandler)(nil)

type SessionHandler interface {
	SessionList(c *gin.Context)
	Create(c *gin.Context)
	Delete(c *gin.Context)
	Top(c *gin.Context)
	Disturb(c *gin.Context)
	ClearUnreadMessage(c *gin.Context)
	GetRecords(c *gin.Context)
}

type sessionHandler struct {
	unreadCache        cache.UnreadCache
	talkSessionDao     dao.TalkSessionDao
	messageCache       *cache.MessageCache
	lockCache          *cache.RedisLock
	userDao            dao.UsersDao
	groupDao           dao.GroupDao
	permissionsService imService.IPermissionService
	talkRecordsDao     dao.TalkRecordsDao
}

func (s sessionHandler) GetRecords(c *gin.Context) {
	params := &types.GetTalkRecordsRequest{}
	ctx := middleware.WrapCtx(c)
	if err := c.ShouldBind(params); err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	uid, err := jwt.HeaderObtainUID(c)
	if err != nil {
		logger.Warn("obtain uid  error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.Unauthorized, err)
		return
	}

	if params.TalkType == constant.ChatGroupMode {
		err := s.permissionsService.IsAuth(ctx, &types.AuthOption{
			TalkType:   params.TalkType,
			UserId:     uid,
			ReceiverId: uint64(params.ReceiverId),
		})

		if err != nil {
			items := make([]types.TalkRecordItem, 0)
			items = append(items, types.TalkRecordItem{
				ID:         1,
				MsgId:      strutil.NewMsgId(),
				Sequence:   1,
				TalkType:   params.TalkType,
				MsgType:    constant.ChatMsgSysText,
				ReceiverId: params.ReceiverId,
				Extra: types.TalkRecordExtraText{
					Content: "暂无权限查看群消息",
				},
				CreatedAt: timeutil.DateTime(),
			})
			response.Success(c, gin.H{
				"cursor": 1,
				"list":   items,
			})
			return
		}
	}

	records, err := s.talkRecordsDao.FindAllTalkRecords(ctx, &types.FindAllTalkRecordsOpt{
		TalkType:   params.TalkType,
		UserId:     uid,
		ReceiverId: params.ReceiverId,
		Cursor:     params.Cursor,
		Limit:      params.Limit,
	})

	if err != nil {
		logger.Warn("获取记录失败: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.ErrGetRecordsFailed, err)
		return

	}

	cursor := 0
	if length := len(records); length > 0 {
		cursor = records[length-1].Sequence
	}

	for i, record := range records {
		if record.IsRevoke == 1 {
			records[i].Extra = make(map[string]any)
		}
	}
	response.Success(c, gin.H{
		"cursor": cursor,
		"list":   records,
	})
}

// ClearUnreadMessage  清除未读消息
// @Summary 清除未读消息
// @Description  清除未读消息
// @Tags    消息
// @Param data body types.TalkSessionClearUnreadNumRequest true "request body"
// @accept  json
// @Produce json
// @Success 200 {object} types.Result
// @Router /api/v1/session/disturb [post]
func (s sessionHandler) ClearUnreadMessage(c *gin.Context) {
	params := &types.TalkSessionClearUnreadNumRequest{}
	ctx := middleware.WrapCtx(c)
	if err := c.ShouldBind(params); err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	uid, err := jwt.HeaderObtainUID(c)
	if err != nil {
		logger.Warn("obtain uid  error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.Unauthorized, err)
		return
	}

	s.unreadCache.Reset(ctx, int(params.TalkType), int(params.ReceiverId), uid)

	response.Success(c)
}

// Disturb  会话免打扰
// @Summary 会话免打扰
// @Description  会话免打扰
// @Tags    聊天列表
// @Param data body types.TalkSessionDisturbRequest true
// @accept  json
// @Produce json
// @Success 200 {object} types.TalkSessionDisturbReply{}
// @Router /api/v1/session/disturb [post]
func (s sessionHandler) Disturb(c *gin.Context) {
	body := &types.TalkSessionDisturbRequest{}
	ctx := middleware.WrapCtx(c)
	if err := c.ShouldBindJSON(body); err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	uid, err := jwt.HeaderObtainUID(c)
	if err != nil {
		logger.Warn("obtain uid  error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.Unauthorized, err)
		return
	}

	err2 := s.talkSessionDao.Disturb(ctx, &model.TalkSessionDisturbOpt{
		UserId:     uid,
		TalkType:   int(body.TalkType),
		ReceiverId: int(body.ReceiverID),
		IsDisturb:  int(body.IsDisturb),
	})

	if err2 != nil {
		logger.Warn("Disturb error: ", logger.Err(err2), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.ErrDisturbSessionFail)
		return
	}
	response.Success(c, "ok")
}

// Top  置顶会话
// @Summary 置顶会话
// @Description  置顶会话
// @Tags    聊天列表
// @Param data body types.TalkSessionTopRequest true
// @accept  json
// @Produce json
// @Success 200 {object} types.TalkSessionTopReply{}
// @Router /api/v1/session/topping [post]
func (s sessionHandler) Top(c *gin.Context) {
	body := &types.TalkSessionTopRequest{}
	ctx := middleware.WrapCtx(c)
	if err := c.ShouldBindJSON(body); err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	uid, err := jwt.HeaderObtainUID(c)
	if err != nil {
		logger.Warn("obtain uid  error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.Unauthorized, err)
		return
	}
	err2 := s.talkSessionDao.Top(ctx, &model.TalkSessionTopOpt{
		UserId: uid,
		Id:     int(body.SessionId),
		Type:   int(body.Type),
	})

	if err2 != nil {
		logger.Warn("top session error: ", logger.Err(err2), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.ErrTopSessionFail)
		return
	}

	response.Success(c, "ok")
}

// Delete  删除会话记录
// @Summary 删除会话记录
// @Description  删除会话记录
// @Tags    聊天列表
// @Param data body types.TalkSessionDeleteRequest true
// @accept  json
// @Produce json
// @Success 200 {object} types.TalkSessionDeleteReply{}
// @Router /api/v1/session/delete [post]
func (s sessionHandler) Delete(c *gin.Context) {
	body := &types.TalkSessionDeleteRequest{}
	if err := c.ShouldBindJSON(body); err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	uid, err := jwt.HeaderObtainUID(c)
	if err != nil {
		logger.Warn("obtain uid  error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.Unauthorized, err)
		return
	}
	ctx := middleware.WrapCtx(c)
	err = s.talkSessionDao.Delete(ctx, uid, int(body.SessionId))
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.ErrDeleteSessionFail)
		return
	}

	response.Success(c, "ok")

}

// Create  创建聊天记录
// @Summary 创建聊天记录
// @Description  创建聊天记录
// @Tags    聊天列表
// @Param data body types.TalkSessionCreateRequest true
// @accept  json
// @Produce json
// @Success 200 {object} types.TalkSessionCreateReply{}
// @Router /api/v1/session/create [post]
func (s sessionHandler) Create(c *gin.Context) {
	body := &types.TalkSessionCreateRequest{}
	err2 := c.ShouldBindJSON(body)
	if err2 != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err2), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	ctx := middleware.WrapCtx(c)

	uid, err := jwt.HeaderObtainUID(c)
	if err != nil {
		logger.Warn("obtain uid  error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.Unauthorized, err)
		return
	}

	if body.TalkType == constant.ChatPrivateMode && int(body.ReceiverID) == uid {
		response.Error(c, ecode.ErrCreateSessionFailed)
		return
	}

	key := fmt.Sprintf("talk:list:%d-%d-%d", uid, body.ReceiverID, body.TalkType)
	if !s.lockCache.Lock(ctx, key, 10) {
		response.Error(c, ecode.ErrCreateSessionFailed)
		return
	}

	talkSession, err := s.talkSessionDao.Create(ctx, &model.TalkSessionCreateOpt{
		UserId:     uid,
		TalkType:   int(body.TalkType),
		ReceiverId: int(body.ReceiverID),
	})
	if err != nil {
		response.Error(c, ecode.ErrCreateSessionFailed, err)
		return
	}

	item := &types.TalkSessionItem{
		ID:         int32(talkSession.ID),
		TalkType:   int32(talkSession.TalkType),
		ReceiverID: int32(talkSession.ReceiverID),
		IsRobot:    int32(talkSession.IsRobot),
		UpdatedAt:  timeutil.DateTime(),
	}

	if item.TalkType == constant.ChatPrivateMode {
		item.UnreadNum = int32(s.unreadCache.Get(ctx, 1, int(body.ReceiverID), uid))
		user, err := s.userDao.GetByID(ctx, talkSession.ReceiverID)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Error(c, ecode.ErrReceiverUserNotFound)
			return
		} else if err != nil {
			response.Error(c, ecode.ErrReceiverUserNotFound, err)
			return
		}
		item.Name = user.Username
		item.Avatar = user.ProfilePicture
	} else if item.TalkType == constant.ChatGroupMode {
		group, err := s.groupDao.GetByID(ctx, uint64(body.ReceiverID))
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Error(c, ecode.ErrReceiverGroupNotFound)
			return
		} else if err != nil {
			response.Error(c, ecode.ErrReceiverGroupNotFound, err)
			return
		}
		item.Name = group.Name
	}

	// 查询缓存消息
	if msg, err := s.messageCache.GetLastMessage(ctx, int(talkSession.TalkType), uid, int(talkSession.ReceiverID)); err == nil {
		item.MsgText = msg.Content
		item.UpdatedAt = msg.Datetime
	}

	response.Success(c, gin.H{
		"data": &types.TalkSessionCreateDetails{
			ID:         item.ID,
			TalkType:   item.TalkType,
			ReceiverId: item.ReceiverID,
			IsTop:      item.IsTop,
			IsDisturb:  item.IsDisturb,
			IsOnline:   item.IsOnline,
			IsRobot:    item.IsRobot,
			Name:       item.Name,
			Avatar:     item.Avatar,
			RemarkName: item.Remark,
			UnreadNum:  item.UnreadNum,
			MsgText:    item.MsgText,
			UpdatedAt:  item.UpdatedAt,
		},
	})

}

// verify 验证uid
func (s sessionHandler) verify(c *gin.Context, uid string) bool {
	if uid == "" {
		response.Error(c, ecode.ErrSessionIdNil)
		return true
	}
	return false
}

// SessionList  获取聊天列表
// @Summary 获取聊天列表
// @Description  获取用户的聊天列表
// @Tags    聊天列表
// @accept  json
// @Produce json
// @Success 200 {object} types.TalkSessionItemsReply{}
// @Router /api/v1/session/list [get]
func (s sessionHandler) SessionList(c *gin.Context) {
	ctx := middleware.WrapCtx(c)
	uid, err := jwt.HeaderObtainUID(c)
	if err != nil {
		logger.Warn("obtain uid  error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.Unauthorized, err)
		return
	}

	unReads := s.unreadCache.All(ctx, uid)
	if len(unReads) > 0 {
		s.talkSessionDao.BatchAddList(ctx, uid, unReads)
	}

	talkSessions, err := s.talkSessionDao.List(ctx, uid)
	if err != nil {
		response.Error(c, ecode.ErrServerQueryList)
	}
	items := make([]*types.TalkSessionItem, 0)
	for _, talkSession := range talkSessions {
		value := &types.TalkSessionItem{
			ID:         int32(talkSession.Id),
			TalkType:   int32(talkSession.TalkType),
			ReceiverID: int32(talkSession.ReceiverId),
			IsTop:      int32(talkSession.IsTop),
			IsDisturb:  int32(talkSession.IsDisturb),
			IsRobot:    int32(talkSession.IsRobot),
			Avatar:     talkSession.UserAvatar,
			MsgText:    "...",
			UpdatedAt:  timeutil.FormatDatetime(talkSession.UpdatedAt),
		}
		if num, ok := unReads[fmt.Sprintf("%d_%d", talkSession.TalkType, talkSession.ReceiverId)]; ok {
			value.UnreadNum = int32(num)
		}

		if talkSession.TalkType == 1 {
			value.Name = talkSession.Nickname
			value.Avatar = talkSession.UserAvatar
			value.IsOnline = int32(strutil.BoolToInt(s.messageCache.IsOnline(ctx, constant.ImChannelChat, strconv.Itoa(int(value.ReceiverID)))))
		} else {
			value.Name = talkSession.GroupName
			value.Avatar = talkSession.GroupAvatar
		}

		// 查询缓存消息
		if msg, err := s.messageCache.GetLastMessage(ctx, talkSession.TalkType, uid, talkSession.ReceiverId); err == nil {
			value.MsgText = msg.Content
			value.UpdatedAt = msg.Datetime
		}
		items = append(items, value)
	}

	response.Success(c, gin.H{
		"data": items,
	})

}

func NewSessionHandler() SessionHandler {
	return &sessionHandler{
		unreadCache:    cache.NewUnreadCache(),
		talkSessionDao: dao.NewTalkSessionDao(model.GetDB()),
		messageCache:   cache.NewMessageCache(model.GetCacheType()),
		groupDao:       dao.NewGroupDao(model.GetDB(), cache.NewGroupCache(model.GetCacheType())),
	}
}
