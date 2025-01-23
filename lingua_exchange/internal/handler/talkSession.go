package handler

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zhufuyi/sponge/pkg/gin/middleware"
	"github.com/zhufuyi/sponge/pkg/gin/response"
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
}

type sessionHandler struct {
	unreadCache    cache.UnreadCache
	talkSessionDao dao.TalkSessionDao
	messageCache   *cache.MessageCache
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
	if uid == "" {
		response.Error(c, ecode.ErrSessionIdNil)
		return
	}
	ctx := middleware.WrapCtx(c)
	toInt, err := strutil.StringToInt(uid)
	if err != nil {
		response.Error(c, ecode.ErrServerConvertID)
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

func NewSessionHandler() SessionHandler {
	return &sessionHandler{
		unreadCache:    cache.NewUnreadCache(),
		talkSessionDao: dao.NewTalkSessionDao(model.GetDB()),
		messageCache:   cache.NewMessageCache(model.GetCacheType()),
	}
}
