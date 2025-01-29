package handler

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/zhufuyi/sponge/pkg/gin/middleware"
	"github.com/zhufuyi/sponge/pkg/gin/response"
	"github.com/zhufuyi/sponge/pkg/logger"
	"lingua_exchange/internal/constant"
	"lingua_exchange/internal/ecode"
	"lingua_exchange/internal/imService"
	"lingua_exchange/internal/types"
	"lingua_exchange/pkg/jwt"
)

// 消息路由实现类
var _ IMMessageHandler = (*imMessageHandler)(nil)

var mapping map[string]func(ctx *context.Context, c *gin.Context) error

type IMMessageHandler interface {
	Publish(ctx *gin.Context)
}

type imMessageHandler struct {
	imAuthService  imService.IPermissionService
	messageService imService.MessageService
}

func NewIMMessageHandler() IMMessageHandler {
	return &imMessageHandler{}
}

func (i imMessageHandler) Publish(c *gin.Context) {
	params := &types.PublishBaseMessageRequest{}
	if err := c.ShouldBindJSON(params); err != nil {
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
	if err := i.imAuthService.IsAuth(ctx, &types.AuthOption{
		TalkType:          params.Receiver.TalkType,
		UserId:            uid,
		ReceiverId:        uint64(params.Receiver.ReceiverId),
		IsVerifyGroupMute: true,
	}); err != nil {
		logger.Warn("send message auth error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.ErrPublishPermissionError, err)
		return
	}

	err = i.transfer(ctx, params.Type, c)
	if err != nil {
		logger.Warn("send message auth error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.ErrPublishMessageError, err)
		return
	}
	response.Success(c)
}

func (i imMessageHandler) onSendText(ctx *context.Context, c *gin.Context) error {
	params := &types.TextMessageRequest{}
	err := c.ShouldBindBodyWith(params, binding.JSON)
	if err != nil {
		logger.Warn("ShouldBindBodyWith error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		return err
	}

	uid, err := jwt.HeaderObtainUID(c)
	if err != nil {
		logger.Warn("obtain uid  error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		return err
	}
	err = i.messageService.SendText(ctx, uid, params)
	if err != nil {
		logger.Warn("send message error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		return err
	}
	return nil
}

func (i imMessageHandler) onSendCode(ctx *context.Context, c *gin.Context) error {
	params := &types.CodeMessageRequest{}
	err := c.ShouldBindBodyWith(params, binding.JSON)
	if err != nil {
		logger.Warn("ShouldBindBodyWith error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		return err
	}
	uid, err := jwt.HeaderObtainUID(c)
	if err != nil {
		logger.Warn("obtain uid  error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		return err
	}
	err = i.messageService.SendCode(ctx, uid, params)
	if err != nil {
		logger.Warn("send code message error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		return err
	}
	return nil
}

func (i imMessageHandler) onSendLocation(ctx *context.Context, c *gin.Context) error {
	params := &types.LocationMessageRequest{}
	err := c.ShouldBindBodyWith(params, binding.JSON)
	if err != nil {
		logger.Warn("ShouldBindBodyWith error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		return err
	}
	uid, err := jwt.HeaderObtainUID(c)
	if err != nil {
		logger.Warn("obtain uid  error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		return err
	}
	err = i.messageService.SendLocation(ctx, uid, params)
	if err != nil {
		logger.Warn("send location message error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		return err
	}
	return nil
}

func (i imMessageHandler) onSendEmoticon(ctx *context.Context, c *gin.Context) error {
	params := &types.EmoticonMessageRequest{}
	if err := c.ShouldBindBodyWith(params, binding.JSON); err != nil {
		logger.Warn("ShouldBindBodyWith error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		return err
	}

	uid, err := jwt.HeaderObtainUID(c)
	if err != nil {
		logger.Warn("obtain uid  error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		return err
	}

	err = i.messageService.SendEmoticon(ctx, uid, params)
	if err != nil {
		logger.Warn("send emotion message error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		return err
	}

	return nil
}

func (i imMessageHandler) onSendVote(ctx *context.Context, c *gin.Context) error {
	params := &types.VoteMessageRequest{}
	if err := c.ShouldBindBodyWith(params, binding.JSON); err != nil {
		logger.Warn("ShouldBindBodyWith error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		return err
	}

	uid, err := jwt.HeaderObtainUID(c)
	if err != nil {
		logger.Warn("obtain uid  error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		return err
	}
	err = i.messageService.SendVote(ctx, uid, params)
	if err != nil {
		logger.Warn("send vote message error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		return err
	}
	return nil
}

func (i imMessageHandler) onSendImage(ctx *context.Context, c *gin.Context) error {
	params := &types.ImageMessageRequest{}
	if err := c.ShouldBindBodyWith(params, binding.JSON); err != nil {
		logger.Warn("ShouldBindBodyWith error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		return err
	}

	uid, err := jwt.HeaderObtainUID(c)
	if err != nil {
		logger.Warn("obtain uid  error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		return err
	}
	err = i.messageService.SendImage(ctx, uid, params)
	if err != nil {
		logger.Warn("send image message error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		return err
	}
	return nil
}

func (i imMessageHandler) onSendVoice(ctx *context.Context, c *gin.Context) error {
	params := &types.VoiceMessageRequest{}
	if err := c.ShouldBindBodyWith(params, binding.JSON); err != nil {
		logger.Warn("ShouldBindBodyWith error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		return err
	}

	uid, err := jwt.HeaderObtainUID(c)
	if err != nil {
		logger.Warn("obtain uid  error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		return err
	}
	err = i.messageService.SendVoice(ctx, uid, params)
	if err != nil {
		logger.Warn("send voice message error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		return err
	}

	return nil
}

func (i imMessageHandler) onSendVideo(ctx *context.Context, c *gin.Context) error {
	params := &types.VideoMessageRequest{}
	if err := c.ShouldBindBodyWith(params, binding.JSON); err != nil {
		logger.Warn("ShouldBindBodyWith error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		return err
	}

	uid, err := jwt.HeaderObtainUID(c)
	if err != nil {
		logger.Warn("obtain uid  error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		return err
	}
	err = i.messageService.SendVideo(ctx, uid, params)
	if err != nil {
		logger.Warn("send video message error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		return err
	}
	return nil
}

func (i imMessageHandler) onSendFile(ctx *context.Context, c *gin.Context) error {
	params := &types.FileMessageRequest{}
	if err := c.ShouldBindBodyWith(params, binding.JSON); err != nil {
		logger.Warn("ShouldBindBodyWith error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		return err
	}

	uid, err := jwt.HeaderObtainUID(c)
	if err != nil {
		logger.Warn("obtain uid  error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		return err
	}
	err = i.messageService.SendFile(ctx, uid, params)
	if err != nil {
		logger.Warn("send file message error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		return err
	}
	return nil
}

func (i imMessageHandler) onSendCard(ctx *context.Context, c *gin.Context) error {
	params := &types.CardMessageRequest{}
	if err := c.ShouldBindBodyWith(params, binding.JSON); err != nil {
		logger.Warn("ShouldBindBodyWith error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		return err
	}

	uid, err := jwt.HeaderObtainUID(c)
	if err != nil {
		logger.Warn("obtain uid  error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		return err
	}
	err = i.messageService.SendBusinessCard(ctx, uid, params)
	if err != nil {
		logger.Warn("send card message error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		return err
	}
	return nil
}

func (i imMessageHandler) onSendForward(ctx *context.Context, c *gin.Context) error {
	params := &types.ForwardMessageRequest{}
	if err := c.ShouldBindBodyWith(params, binding.JSON); err != nil {
		logger.Warn("ShouldBindBodyWith error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		return err
	}

	uid, err := jwt.HeaderObtainUID(c)
	if err != nil {
		logger.Warn("obtain uid  error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		return err
	}
	err = i.messageService.SendForward(ctx, uid, params)
	if err != nil {
		logger.Warn("send forward message error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
	}
	return nil
}

func (i imMessageHandler) onMixedMessage(ctx *context.Context, c *gin.Context) error {
	params := &types.MixedMessageRequest{}
	if err := c.ShouldBindBodyWith(params, binding.JSON); err != nil {
		logger.Warn("ShouldBindBodyWith error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		return err
	}

	uid, err := jwt.HeaderObtainUID(c)
	if err != nil {
		logger.Warn("obtain uid  error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		return err
	}
	err = i.messageService.SendMixedMessage(ctx, uid, params)
	if err != nil {
		logger.Warn("send mixed message error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		return err
	}
	return nil
}

func (i imMessageHandler) transfer(ctx context.Context, typeValue string, c *gin.Context) error {
	if mapping == nil {
		mapping = make(map[string]func(ctx *context.Context, c *gin.Context) error)
		mapping[constant.Text] = i.onSendText
		mapping[constant.Code] = i.onSendCode
		mapping[constant.Location] = i.onSendLocation
		mapping[constant.Emoticon] = i.onSendEmoticon
		mapping[constant.Vote] = i.onSendVote
		mapping[constant.Image] = i.onSendImage
		mapping[constant.Voice] = i.onSendVoice
		mapping[constant.Video] = i.onSendVideo
		mapping[constant.File] = i.onSendFile
		mapping[constant.Card] = i.onSendCard
		mapping[constant.Forward] = i.onSendForward
		mapping[constant.Mixed] = i.onMixedMessage
	}
	if call, ok := mapping[typeValue]; ok {
		return call(&ctx, c)
	}

	return nil
}
