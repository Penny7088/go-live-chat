package ecode

import (
	"github.com/zhufuyi/sponge/pkg/errcode"
)

// thirdPartyAuth business-level http error codes.
// the thirdPartyAuthNO value range is 1~100, if the same error code is used, it will cause panic.
var (
	thirdPartyAuthNO       = 4
	thirdPartyAuthName     = "thirdPartyAuth"
	thirdPartyAuthBaseCode = errcode.HCode(thirdPartyAuthNO)

	ErrCreateThirdPartyAuth     = errcode.NewError(thirdPartyAuthBaseCode+1, "failed to create "+thirdPartyAuthName)
	ErrDeleteByIDThirdPartyAuth = errcode.NewError(thirdPartyAuthBaseCode+2, "failed to delete "+thirdPartyAuthName)
	ErrUpdateByIDThirdPartyAuth = errcode.NewError(thirdPartyAuthBaseCode+3, "failed to update "+thirdPartyAuthName)
	ErrGetByIDThirdPartyAuth    = errcode.NewError(thirdPartyAuthBaseCode+4, "failed to get "+thirdPartyAuthName+" details")
	ErrListThirdPartyAuth       = errcode.NewError(thirdPartyAuthBaseCode+5, "failed to list of "+thirdPartyAuthName)

	ErrDeleteByIDsThirdPartyAuth    = errcode.NewError(thirdPartyAuthBaseCode+6, "failed to delete by batch ids "+thirdPartyAuthName)
	ErrGetByConditionThirdPartyAuth = errcode.NewError(thirdPartyAuthBaseCode+7, "failed to get "+thirdPartyAuthName+" details by conditions")
	ErrListByIDsThirdPartyAuth      = errcode.NewError(thirdPartyAuthBaseCode+8, "failed to list by batch ids "+thirdPartyAuthName)
	ErrListByLastIDThirdPartyAuth   = errcode.NewError(thirdPartyAuthBaseCode+9, "failed to list by last id "+thirdPartyAuthName)

	// error codes are globally unique, adding 1 to the previous error code
)
