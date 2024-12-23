package ecode

import "github.com/zhufuyi/sponge/pkg/errcode"

var (
	globalConfig         = 81
	globalConfigName     = "globalConfig"
	globalConfigBaseCode = errcode.HCode(globalConfig)

	ErrIpNotFound                 = errcode.NewError(globalConfigBaseCode+3, "Client IP Not Found"+globalConfigName)
	ErrSendCode                   = errcode.NewError(globalConfigBaseCode+4, "Verification code send fail"+globalConfigName)
	ErrGetCode                    = errcode.NewError(globalConfigBaseCode+5, "get code "+globalConfigName)
	ErrVerificationSentRepeatedly = errcode.NewError(globalConfigBaseCode+6, "Verification sent repeatedly "+globalConfigName)
)
