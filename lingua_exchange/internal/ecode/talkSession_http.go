package ecode

import "github.com/zhufuyi/sponge/pkg/errcode"

var (
	sessionNo          = 90
	sessionName        = "sessionName"
	sessionBaseCode    = errcode.HCode(sessionNo)
	ErrSessionIdNil    = errcode.NewError(sessionBaseCode+1, "session id is null "+sessionName)
	ErrServerConvertID = errcode.NewError(sessionBaseCode+2, "server convert id is  error "+sessionName)
	ErrServerQueryList = errcode.NewError(sessionBaseCode+3, "server query users group talk session  error "+sessionName)
)
