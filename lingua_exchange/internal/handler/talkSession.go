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
}

type sessionHandler struct {
	unreadCache    cache.UnreadCache
	talkSessionDao dao.TalkSessionDao
	messageCache   *cache.MessageCache
	lockCache      *cache.RedisLock
	userDao        dao.UsersDao
	groupDao       dao.GroupDao
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
	uid := jwt.HeaderObtainUID(c)
	ctx := middleware.WrapCtx(c)
	if err := c.ShouldBindJSON(body); err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
	}
	if s.verify(c, uid) {
		return
	}

	toInt, _, done := s.convertUID(c, uid)
	if done {
		return
	}
	err := s.talkSessionDao.Top(ctx, &model.TalkSessionTopOpt{
		UserId: toInt,
		Id:     int(body.SessionId),
		Type:   int(body.Type),
	})

	if err != nil {
		logger.Warn("top session error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
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
	uid := jwt.HeaderObtainUID(c)
	if err := c.ShouldBindJSON(body); err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
	}
	if s.verify(c, uid) {
		return
	}

	toInt, err, done := s.convertUID(c, uid)
	if done {
		return
	}
	ctx := middleware.WrapCtx(c)
	err = s.talkSessionDao.Delete(ctx, toInt, int(body.SessionId))
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
	uid := jwt.HeaderObtainUID(c)
	ctx := middleware.WrapCtx(c)
	if s.verify(c, uid) {
		return
	}

	toInt, err, done := s.convertUID(c, uid)
	if done {
		return
	}

	err2 := c.ShouldBindJSON(body)
	if err2 != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err2), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
	}

	if body.TalkType == constant.ChatPrivateMode && int(body.ReceiverID) == toInt {
		response.Error(c, ecode.ErrCreateSessionFailed)
		return
	}

	key := fmt.Sprintf("talk:list:%d-%d-%d", toInt, body.ReceiverID, body.TalkType)
	if !s.lockCache.Lock(ctx, key, 10) {
		response.Error(c, ecode.ErrCreateSessionFailed)
		return
	}

	talkSession, err := s.talkSessionDao.Create(ctx, &model.TalkSessionCreateOpt{
		UserId:     toInt,
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
		item.UnreadNum = int32(s.unreadCache.Get(ctx, 1, int(body.ReceiverID), toInt))
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
	if msg, err := s.messageCache.GetLastMessage(ctx, int(talkSession.TalkType), toInt, int(talkSession.ReceiverID)); err == nil {
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
	uid := jwt.HeaderObtainUID(c)
	if s.verify(c, uid) {
		return
	}
	ctx := middleware.WrapCtx(c)
	toInt, err, done := s.convertUID(c, uid)
	if done {
		return
	}

	unReads := s.unreadCache.All(ctx, toInt)
	if len(unReads) > 0 {
		s.talkSessionDao.BatchAddList(ctx, toInt, unReads)
	}

	talkSessions, err := s.talkSessionDao.List(ctx, toInt)
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
		if msg, err := s.messageCache.GetLastMessage(ctx, talkSession.TalkType, toInt, talkSession.ReceiverId); err == nil {
			value.MsgText = msg.Content
			value.UpdatedAt = msg.Datetime
		}
		items = append(items, value)
	}

	response.Success(c, gin.H{
		"data": items,
	})

}

func (s sessionHandler) convertUID(c *gin.Context, uid string) (int, error, bool) {
	toInt, err := strutil.StringToInt(uid)
	if err != nil {
		response.Error(c, ecode.ErrServerConvertID)
		return 0, nil, true
	}
	return toInt, err, false
}

func NewSessionHandler() SessionHandler {
	return &sessionHandler{
		unreadCache:    cache.NewUnreadCache(),
		talkSessionDao: dao.NewTalkSessionDao(model.GetDB()),
		messageCache:   cache.NewMessageCache(model.GetCacheType()),
		groupDao:       dao.NewGroupDao(model.GetDB(), cache.NewGroupCache(model.GetCacheType())),
	}
}
