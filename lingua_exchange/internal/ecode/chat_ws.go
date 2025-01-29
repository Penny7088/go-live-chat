package ecode

import "github.com/zhufuyi/sponge/pkg/errcode"

var (
	chatNo                      = 80
	chatName                    = "chatName"
	chatBaseCode                = errcode.HCode(chatNo)
	ErrServerClosed             = errcode.NewError(chatBaseCode+1, "Err Server Closed "+chatName)
	ErrServerNotCache           = errcode.NewError(chatBaseCode+2, "Err Server not cache,setup cache "+chatName)
	ErrWebsocketHealthSubscribe = errcode.NewError(chatBaseCode+3, "Websocket HealthSubscribe Report Err "+chatName)
	ErrPublishPermissionError   = errcode.NewError(chatBaseCode+4, "publish message auth error"+chatName)
	ErrPublishMessageError      = errcode.NewError(chatBaseCode+5, "publish message  error"+chatName)
	ErrPublishTextMessageError  = errcode.NewError(chatBaseCode+6, "publish message text error"+chatName)
	ErrPublishCodeMessageError  = errcode.NewError(chatBaseCode+7, "publish code message error"+chatName)
)
