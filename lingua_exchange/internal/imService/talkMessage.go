package imService

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/zhufuyi/sponge/pkg/ggorm/query"
	"github.com/zhufuyi/sponge/pkg/logger"
	"gorm.io/gorm"
	"lingua_exchange/internal/cache"
	"lingua_exchange/internal/constant"
	"lingua_exchange/internal/dao"
	"lingua_exchange/internal/model"
	"lingua_exchange/internal/types"
	"lingua_exchange/pkg/jsonutil"
	"lingua_exchange/pkg/strutil"
	"lingua_exchange/pkg/timeutil"
)

// 消息发送控制类
var _ IMessageService = (*MessageService)(nil)

type IMessageService interface {
	// SendSystemText 系统文本消息
	SendSystemText(ctx *context.Context, uid int, req *types.TextMessageRequest) error
	// SendText 文本消息
	SendText(ctx *context.Context, uid int, req *types.TextMessageRequest) error
	// SendImage 图片文件消息
	SendImage(ctx *context.Context, uid int, req *types.ImageMessageRequest) error
	// SendVoice 语音文件消息
	SendVoice(ctx *context.Context, uid int, req *types.VoiceMessageRequest) error
	// SendVideo 视频文件消息
	SendVideo(ctx *context.Context, uid int, req *types.VideoMessageRequest) error
	// SendFile 文件消息
	SendFile(ctx *context.Context, uid int, req *types.FileMessageRequest) error
	// SendCode 代码消息
	SendCode(ctx *context.Context, uid int, req *types.CodeMessageRequest) error
	// SendVote 投票消息
	SendVote(ctx *context.Context, uid int, req *types.VoteMessageRequest) error
	// SendEmoticon 表情消息
	SendEmoticon(ctx *context.Context, uid int, req *types.EmoticonMessageRequest) error
	// SendForward 转发消息
	SendForward(ctx *context.Context, uid int, req *types.ForwardMessageRequest) error
	// SendLocation 位置消息
	SendLocation(ctx *context.Context, uid int, req *types.LocationMessageRequest) error
	// SendBusinessCard 推送用户名片消息
	SendBusinessCard(ctx *context.Context, uid int, req *types.CardMessageRequest) error
	// SendLogin 推送用户登录消息
	SendLogin(ctx *context.Context, uid int, req *types.LoginMessageRequest) error
	// SendSysOther 推送其它消息
	SendSysOther(ctx *context.Context, data *model.TalkRecords) error
	// SendMixedMessage 图文消息
	SendMixedMessage(ctx *context.Context, uid int, req *types.MixedMessageRequest) error
	// Revoke 撤回消息
	Revoke(ctx *context.Context, uid int, msgId string) error
	// Vote 投票
	// Vote(ctx *context.Context, uid int, msgId string, optionsValue string) (*repo.VoteStatistics, error)
}

type MessageService struct {
	talkRecordsDao   dao.TalkRecordsDao
	talkRecordsCache cache.TalkRecordsCache
	usersDao         dao.UsersDao
	unreadCache      cache.UnreadCache
	redis            *redis.Client
	groupMemberDao   dao.GroupMemberDao
	messageCache     *cache.MessageCache
	serverCache      cache.ServerCache
	db               *gorm.DB
}

func NewMessageService() IMessageService {
	return &MessageService{
		talkRecordsDao:   dao.NewTalkRecordsDao(model.GetDB(), cache.NewTalkRecordsCache(model.GetCacheType())),
		usersDao:         dao.NewUsersDao(model.GetDB(), cache.NewUsersCache(model.GetCacheType())),
		talkRecordsCache: cache.NewTalkRecordsCache(model.GetCacheType()),
		unreadCache:      cache.NewUnreadCache(),
		redis:            model.GetRedisCli(),
		groupMemberDao:   dao.NewGroupMemberDao(model.GetDB(), cache.NewGroupMemberCache(model.GetCacheType())),
		messageCache:     cache.NewMessageCache(model.GetCacheType()),
		serverCache:      cache.NewServerCache(model.GetCacheType()),
		db:               model.GetDB(),
	}
}

func (m MessageService) SendSystemText(ctx *context.Context, uid int, req *types.TextMessageRequest) error {
	// TODO implement me
	panic("implement me")
}

func (m MessageService) SendText(ctx *context.Context, uid int, req *types.TextMessageRequest) error {
	data := &model.TalkRecords{
		TalkType:   uint(req.Receiver.TalkType),
		MsgType:    constant.ChatMsgTypeText,
		QuoteID:    req.QuoteID,
		UserID:     uid,
		ReceiverID: req.Receiver.ReceiverID,
		Extra: jsonutil.Encode(types.TalkRecordExtraText{
			Content:  strutil.EscapeHtml(req.Content),
			Mentions: req.Mentions,
		}),
	}
	return m.save(ctx, data)
}

func (m MessageService) SendImage(ctx *context.Context, uid int, req *types.ImageMessageRequest) error {
	data := &model.TalkRecords{
		TalkType:   uint(int(req.Receiver.TalkType)),
		MsgType:    constant.ChatMsgTypeImage,
		QuoteID:    req.QuoteID,
		UserID:     uid,
		ReceiverID: req.Receiver.ReceiverID,
		Extra: jsonutil.Encode(&types.TalkRecordExtraImage{
			Url:    req.URL,
			Width:  int(req.Width),
			Height: int(req.Height),
		}),
	}

	return m.save(ctx, data)
}

func (m MessageService) SendVoice(ctx *context.Context, uid int, req *types.VoiceMessageRequest) error {
	data := &model.TalkRecords{
		TalkType:   uint(int(req.Receiver.TalkType)),
		MsgType:    constant.ChatMsgTypeAudio,
		UserID:     uid,
		ReceiverID: req.Receiver.ReceiverID,
		Extra: jsonutil.Encode(&types.TalkRecordExtraAudio{
			Size:     int(req.Size),
			Url:      req.URL,
			Duration: 0,
		}),
	}

	return m.save(ctx, data)
}

func (m MessageService) SendVideo(ctx *context.Context, uid int, req *types.VideoMessageRequest) error {
	data := &model.TalkRecords{
		TalkType:   uint(int(req.Receiver.TalkType)),
		MsgType:    constant.ChatMsgTypeVideo,
		UserID:     uid,
		ReceiverID: req.Receiver.ReceiverID,
		Extra: jsonutil.Encode(&types.TalkRecordExtraVideo{
			Cover:    req.Cover,
			Size:     int(req.Size),
			Url:      req.URL,
			Duration: int(req.Duration),
		}),
	}

	return m.save(ctx, data)
}

func (m MessageService) SendFile(ctx *context.Context, uid int, req *types.FileMessageRequest) error {
	return errors.New("not implemented")
}

func (m MessageService) SendCode(ctx *context.Context, uid int, req *types.CodeMessageRequest) error {
	data := &model.TalkRecords{
		TalkType:   uint(int(req.Receiver.TalkType)),
		MsgType:    constant.ChatMsgTypeCode,
		UserID:     uid,
		ReceiverID: int(req.Receiver.ReceiverID),
		Extra: jsonutil.Encode(&types.TalkRecordExtraCode{
			Lang: req.Lang,
			Code: req.Code,
		}),
	}

	return m.save(ctx, data)
}

func (m MessageService) SendVote(ctx *context.Context, uid int, req *types.VoteMessageRequest) error {
	return errors.New("not implemented")
}

func (m MessageService) SendEmoticon(ctx *context.Context, uid int, req *types.EmoticonMessageRequest) error {
	return errors.New("not implemented")
}

func (m MessageService) SendForward(ctx *context.Context, uid int, req *types.ForwardMessageRequest) error {
	return errors.New("not implemented")
}

func (m MessageService) SendLocation(ctx *context.Context, uid int, req *types.LocationMessageRequest) error {
	data := &model.TalkRecords{
		TalkType:   req.Receiver.TalkType,
		MsgType:    constant.ChatMsgTypeLocation,
		UserID:     uid,
		ReceiverID: req.Receiver.ReceiverID,
		Extra: jsonutil.Encode(&types.TalkRecordExtraLocation{
			Longitude:   req.Longitude,
			Latitude:    req.Latitude,
			Description: req.Description,
		}),
	}

	return m.save(ctx, data)
}

func (m MessageService) SendBusinessCard(ctx *context.Context, uid int, req *types.CardMessageRequest) error {
	data := &model.TalkRecords{
		TalkType:   uint(int(req.Receiver.TalkType)),
		MsgType:    constant.ChatMsgTypeCard,
		UserID:     uid,
		ReceiverID: req.Receiver.ReceiverID,
		Extra: jsonutil.Encode(&types.TalkRecordExtraCard{
			UserId: req.UserID,
		}),
	}

	return m.save(ctx, data)
}

func (m MessageService) SendLogin(ctx *context.Context, uid int, req *types.LoginMessageRequest) error {
	return errors.New("not implemented")
}

func (m MessageService) SendSysOther(ctx *context.Context, data *model.TalkRecords) error {
	return m.save(ctx, data)
}

func (m MessageService) SendMixedMessage(ctx *context.Context, uid int, req *types.MixedMessageRequest) error {
	items := make([]*types.TalkRecordExtraMixedItem, 0)

	for _, item := range req.Items {
		items = append(items, &types.TalkRecordExtraMixedItem{
			Type:    int(item.Type),
			Content: item.Content,
		})
	}

	data := &model.TalkRecords{
		TalkType:   req.Receiver.TalkType,
		MsgType:    constant.ChatMsgTypeMixed,
		QuoteID:    req.QuoteID,
		UserID:     uid,
		ReceiverID: req.Receiver.ReceiverID,
		Extra:      jsonutil.Encode(types.TalkRecordExtraMixed{Items: items}),
	}

	return m.save(ctx, data)
}

func (m MessageService) Revoke(ctx *context.Context, uid int, msgId string) error {
	ctx2 := *ctx
	var record *model.TalkRecords
	if err := m.db.First(&record, "msg_id = ?", msgId).Error; err != nil {
		return err
	}

	if record.IsRevoke == 1 {
		return nil
	}

	if record.UserID != uid {
		return errors.New("无权撤回回消息")
	}

	if time.Now().Unix() > record.CreatedAt.Add(3*time.Minute).Unix() {
		return errors.New("超出有效撤回时间范围，无法进行撤销！")
	}
	record.IsRevoke = 1
	err := m.talkRecordsDao.UpdateByID(ctx2, record)
	if err != nil {
		return err
	}

	user, err := m.usersDao.GetByCondition(ctx2, &query.Conditions{
		Columns: []query.Column{
			{
				Name:  "id",
				Value: record.UserID,
			},
		},
	})

	if err != nil {
		return err
	}

	_ = m.messageCache.SetLastMessage(ctx2, int(record.TalkType), record.UserID, record.ReceiverID, &cache.LastCacheMessage{
		Content:  fmt.Sprintf("%s: 撤回了一条消息", user.Username),
		Datetime: timeutil.DateTime(),
	})

	body := map[string]any{
		"event": constant.SubEventImMessageRevoke,
		"data": jsonutil.Encode(map[string]any{
			"msg_id": record.MsgID,
		}),
	}

	m.redis.Publish(ctx2, constant.ImTopicChat, jsonutil.Encode(body))

	return nil
}

func (m MessageService) save(ctx *context.Context, data *model.TalkRecords) error {
	ctx2 := *ctx
	if data.MsgID == "" {
		data.MsgID = strutil.NewMsgId()
	}
	m.loadReply(ctx2, data)

	m.loadSequence(ctx2, data)

	err := m.talkRecordsDao.Create(ctx2, data)
	if err != nil {
		return err
	}

	lastMessage := types.TalkLastMessage{
		MsgId:      data.MsgID,
		Sequence:   int(data.Sequence),
		MsgType:    data.MsgType,
		UserId:     data.UserID,
		ReceiverId: data.ReceiverID,
		CreatedAt:  time.Now().Format(time.DateTime),
	}

	switch data.MsgType {
	case constant.ChatMsgTypeText:
		extra := types.TalkRecordExtraText{}
		if err := jsonutil.Decode(data.Extra, &extra); err != nil {
			logger.Errorf("MessageService Json Decode err: %s", err.Error())
			return err
		}
		lastMessage.Content = strutil.MtSubstr(strutil.ReplaceImgAll(extra.Content), 0, 300)
	default:
		if value, ok := constant.ChatMsgTypeMapping[data.MsgType]; ok {
			lastMessage.Content = value
		} else {
			lastMessage.Content = "[未知消息]"
		}
	}

	return nil

}

func (m *MessageService) loadReply(ctx context.Context, data *model.TalkRecords) {
	if data.QuoteID == "" {
		return
	}

	if data.Extra == "" {
		data.Extra = "{}"
	}

	extra := make(map[string]any)
	err := jsonutil.Decode(data.Extra, &extra)
	if err != nil {
		logger.Errorf("MessageService Json Decode err: %s", err.Error())
		return
	}
	record, err := m.talkRecordsDao.GetByID(ctx, data.QuoteID)
	if err != nil {
		return
	}

	users, err := m.usersDao.GetByCondition(ctx, &query.Conditions{
		Columns: []query.Column{
			{
				Name:  "Username",
				Value: record.UserID,
			},
		},
	})
	if err != nil {
		return
	}

	reply := types.Reply{
		UserId:   int(record.UserID),
		Nickname: users.Username,
		MsgType:  1,
		MsgId:    record.MsgID,
	}

	if record.MsgType != constant.ChatMsgTypeText {
		reply.Content = "[未知消息]"
		if value, ok := constant.ChatMsgTypeMapping[record.MsgType]; ok {
			reply.Content = value
		}
	} else {
		extra := types.TalkRecordExtraText{}
		if err := jsonutil.Decode(record.Extra, &extra); err != nil {
			logger.Errorf("loadReply Json Decode err: %s", err.Error())
			return
		}

		reply.Content = extra.Content
	}

	extra["reply"] = reply

	data.Extra = jsonutil.Encode(extra)
}

// loadSequence 加载时序
func (m *MessageService) loadSequence(ctx context.Context, data *model.TalkRecords) {
	if data.TalkType == constant.ChatGroupMode {
		data.Sequence = m.talkRecordsCache.GetSequence(ctx, 0, int(data.ReceiverID))
	} else {
		data.Sequence = m.talkRecordsCache.GetSequence(ctx, int(data.UserID), int(data.ReceiverID))
	}
}

// 消息发送完，后置处理
func (m *MessageService) afterHandler(ctx context.Context, record *model.TalkRecords, opt types.TalkLastMessage) {

	if record.TalkType == constant.ChatPrivateMode {
		m.unreadCache.Incr(ctx, constant.ChatPrivateMode, int(record.UserID), int(record.ReceiverID))
		if record.MsgType == constant.ChatMsgSysText {
			m.unreadCache.Incr(ctx, 1, int(record.ReceiverID), int(record.UserID))
		}
	} else if record.TalkType == constant.ChatGroupMode {
		pipe := m.redis.Pipeline()
		for _, uid := range m.groupMemberDao.GetMemberIds(ctx, int(record.ReceiverID)) {
			if uid != int(record.UserID) {
				m.unreadCache.PipeIncr(ctx, pipe, constant.ChatGroupMode, int(record.ReceiverID), uid)
			}
		}
		_, _ = pipe.Exec(ctx)
	}

	_ = m.messageCache.SetLastMessage(ctx, int(record.TalkType), int(record.UserID), int(record.ReceiverID), &cache.LastCacheMessage{
		Content:  opt.Content,
		Datetime: opt.CreatedAt,
	})

	content := jsonutil.Encode(map[string]any{
		"event": constant.SubEventImMessage,
		"data": jsonutil.Encode(map[string]any{
			"sender_id":   record.UserID,
			"receiver_id": record.ReceiverID,
			"talk_type":   record.TalkType,
			"msg_id":      record.MsgID,
		}),
	})

	if record.TalkType == constant.ChatPrivateMode {
		sids := m.serverCache.All(ctx, 1)

		if len(sids) > 3 {
			pipe := m.redis.Pipeline()

			for _, sid := range sids {
				for _, uid := range []int{record.UserID, record.ReceiverID} {
					if !m.messageCache.IsCurrentServerOnline(ctx, sid, constant.ImChannelChat, strconv.Itoa(uid)) {
						continue
					}

					pipe.Publish(ctx, fmt.Sprintf(constant.ImTopicChatPrivate, sid), content)
				}
			}

			if _, err := pipe.Exec(ctx); err == nil {
				return
			}
		}

	}

	if err := m.redis.Publish(ctx, constant.ImTopicChat, content).Err(); err != nil {
		logger.Errorf("[ALL]消息推送失败 %s", err.Error())
	}
}
