package handler

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/zhufuyi/sponge/pkg/gin/middleware"
	"github.com/zhufuyi/sponge/pkg/logger"
	"lingua_exchange/internal/cache"
	"lingua_exchange/internal/chat/event"
	"lingua_exchange/internal/config"
	"lingua_exchange/internal/dao"
	"lingua_exchange/internal/model"
	"lingua_exchange/pkg/jwt"
	"lingua_exchange/pkg/socket"
	"lingua_exchange/pkg/socket/adapter"
	"lingua_exchange/pkg/strutil"
)

var _ MessageHandler = (*messageHandler)(nil)

type MessageHandler interface {
	Connection(ctx *gin.Context)
}

type messageHandler struct {
	messageCache *cache.MessageCache
	event        *event.ChatEvent
}

func (m messageHandler) Connection(ctx *gin.Context) {
	err := m.conn(ctx)
	if err != nil {
		logger.Error("im connection error", logger.Err(err), middleware.GCtxRequestIDField(ctx))
	}
}

func (m messageHandler) conn(ctx *gin.Context) error {
	conn, err := adapter.NewWsAdapter(ctx.Writer, ctx.Request)
	if err != nil {
		log.Printf("websocket connect error: %s", err.Error())
		return err
	}
	uid := jwt.HeaderObtainUID(ctx)
	id, err := strutil.StringToInt(uid)
	if err != nil {
		return err
	}
	return m.newClient(id, conn)
}

func (m messageHandler) newClient(uid int, conn socket.IConn) error {
	return socket.NewClient(conn, &socket.ClientOption{
		Uid:     uid,
		Channel: socket.Session.Chat,
		Storage: m.messageCache,
		Buffer:  10,
	}, socket.NewEvent(
		// 连接成功回调
		socket.WithOpenEvent(m.event.OnOpen),
		// 接收消息回调
		socket.WithMessageEvent(m.event.OnMessage),
		// 关闭连接回调
		socket.WithCloseEvent(m.event.OnClose),
	))
}

func NewMessageHandler() MessageHandler {
	chatEvent := &event.ChatEvent{
		Redis:           model.GetRedisCli(),
		Config:          config.Get(),
		RoomStorage:     cache.NewChatRoomCache(model.GetCacheType()),
		GroupMemberRepo: dao.NewGroupMemberDao(model.GetDB(), cache.NewGroupMemberCache(model.GetCacheType()))}

	messageCache := cache.NewMessageCache(model.GetCacheType())

	return &messageHandler{messageCache: messageCache, event: chatEvent}
}
