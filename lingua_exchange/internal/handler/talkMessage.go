package handler

import (
	"context"

	"lingua_exchange/internal/model"
	"lingua_exchange/internal/types"
)

var _ IMessageService = (*MessageService)(nil)

type IMessageService interface {
	// SendSystemText 系统文本消息
	SendSystemText(ctx context.Context, uid int, req *types.TextMessageRequest) error
	// SendText 文本消息
	SendText(ctx context.Context, uid int, req *types.TextMessageRequest) error
	// SendImage 图片文件消息
	SendImage(ctx context.Context, uid int, req *types.ImageMessageRequest) error
	// SendVoice 语音文件消息
	SendVoice(ctx context.Context, uid int, req *types.VoiceMessageRequest) error
	// SendVideo 视频文件消息
	SendVideo(ctx context.Context, uid int, req *types.VideoMessageRequest) error
	// SendFile 文件消息
	SendFile(ctx context.Context, uid int, req *types.FileMessageRequest) error
	// SendCode 代码消息
	SendCode(ctx context.Context, uid int, req *types.CodeMessageRequest) error
	// SendVote 投票消息
	SendVote(ctx context.Context, uid int, req *types.VoteMessageRequest) error
	// SendEmoticon 表情消息
	SendEmoticon(ctx context.Context, uid int, req *types.EmoticonMessageRequest) error
	// SendForward 转发消息
	SendForward(ctx context.Context, uid int, req *types.ForwardMessageRequest) error
	// SendLocation 位置消息
	SendLocation(ctx context.Context, uid int, req *types.LocationMessageRequest) error
	// SendBusinessCard 推送用户名片消息
	SendBusinessCard(ctx context.Context, uid int, req *types.CardMessageRequest) error
	// SendLogin 推送用户登录消息
	SendLogin(ctx context.Context, uid int, req *types.LoginMessageRequest) error
	// SendSysOther 推送其它消息
	SendSysOther(ctx context.Context, data *model.TalkRecords) error
	// SendMixedMessage 图文消息
	SendMixedMessage(ctx context.Context, uid int, req *types.MixedMessageRequest) error
	// Revoke 撤回消息
	Revoke(ctx context.Context, uid int, msgId string) error
	// Vote 投票
	// Vote(ctx context.Context, uid int, msgId string, optionsValue string) (*repo.VoteStatistics, error)
}

type MessageService struct {
}

func (m MessageService) SendSystemText(ctx context.Context, uid int, req *types.TextMessageRequest) error {
	// TODO implement me
	panic("implement me")
}

func (m MessageService) SendText(ctx context.Context, uid int, req *types.TextMessageRequest) error {
	// TODO implement me
	panic("implement me")
}

func (m MessageService) SendImage(ctx context.Context, uid int, req *types.ImageMessageRequest) error {
	// TODO implement me
	panic("implement me")
}

func (m MessageService) SendVoice(ctx context.Context, uid int, req *types.VoiceMessageRequest) error {
	// TODO implement me
	panic("implement me")
}

func (m MessageService) SendVideo(ctx context.Context, uid int, req *types.VideoMessageRequest) error {
	// TODO implement me
	panic("implement me")
}

func (m MessageService) SendFile(ctx context.Context, uid int, req *types.FileMessageRequest) error {
	// TODO implement me
	panic("implement me")
}

func (m MessageService) SendCode(ctx context.Context, uid int, req *types.CodeMessageRequest) error {
	// TODO implement me
	panic("implement me")
}

func (m MessageService) SendVote(ctx context.Context, uid int, req *types.VoteMessageRequest) error {
	// TODO implement me
	panic("implement me")
}

func (m MessageService) SendEmoticon(ctx context.Context, uid int, req *types.EmoticonMessageRequest) error {
	// TODO implement me
	panic("implement me")
}

func (m MessageService) SendForward(ctx context.Context, uid int, req *types.ForwardMessageRequest) error {
	// TODO implement me
	panic("implement me")
}

func (m MessageService) SendLocation(ctx context.Context, uid int, req *types.LocationMessageRequest) error {
	// TODO implement me
	panic("implement me")
}

func (m MessageService) SendBusinessCard(ctx context.Context, uid int, req *types.CardMessageRequest) error {
	// TODO implement me
	panic("implement me")
}

func (m MessageService) SendLogin(ctx context.Context, uid int, req *types.LoginMessageRequest) error {
	// TODO implement me
	panic("implement me")
}

func (m MessageService) SendSysOther(ctx context.Context, data *model.TalkRecords) error {
	// TODO implement me
	panic("implement me")
}

func (m MessageService) SendMixedMessage(ctx context.Context, uid int, req *types.MixedMessageRequest) error {
	// TODO implement me
	panic("implement me")
}

func (m MessageService) Revoke(ctx context.Context, uid int, msgId string) error {
	// TODO implement me
	panic("implement me")
}
