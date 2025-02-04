package consume

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"lingua_exchange/internal/cache"
	"lingua_exchange/internal/config"
	"lingua_exchange/internal/constant"
	"lingua_exchange/internal/dao"
	"lingua_exchange/internal/model"
)

var handlers map[string]func(ctx context.Context, data []byte)

type IMHandler struct {
	chatRoom         cache.ChatRoomCache
	messageCache     *cache.MessageCache
	talkRecordsCache cache.TalkRecordsCache
	talkRecordsDao   dao.TalkRecordsDao
	redis            *redis.Client
	db               *gorm.DB
	config           *config.Config
}

func NewIMHandler() *IMHandler {
	return &IMHandler{
		chatRoom:         cache.NewChatRoomCache(model.GetCacheType()),
		messageCache:     cache.NewMessageCache(model.GetCacheType()),
		talkRecordsCache: cache.NewTalkRecordsCache(model.GetCacheType()),
		talkRecordsDao:   dao.NewTalkRecordsDao(model.GetDB(), cache.NewTalkRecordsCache(model.GetCacheType())),
		redis:            model.GetRedisCli(),
		db:               model.GetDB(),
		config:           config.Get(),
	}
}

func (h *IMHandler) init() {
	handlers = make(map[string]func(ctx context.Context, data []byte))

	handlers[constant.SubEventImMessage] = h.onConsumeTalk
	handlers[constant.SubEventImMessageKeyboard] = h.onConsumeTalkKeyboard
	handlers[constant.SubEventImMessageRead] = h.onConsumeTalkRead
	handlers[constant.SubEventImMessageRevoke] = h.onConsumeTalkRevoke
	handlers[constant.SubEventContactStatus] = h.onConsumeContactStatus
	handlers[constant.SubEventContactApply] = h.onConsumeContactApply
	handlers[constant.SubEventGroupJoin] = h.onConsumeGroupJoin
	handlers[constant.SubEventGroupApply] = h.onConsumeGroupApply
}

func (h *IMHandler) Call(ctx context.Context, event string, data []byte) {
	if handlers == nil {
		h.init()
	}

	if call, ok := handlers[event]; ok {
		call(ctx, data)
	} else {
		log.Printf("consume chat event: [%s]未注册回调事件\n", event)
	}
}
