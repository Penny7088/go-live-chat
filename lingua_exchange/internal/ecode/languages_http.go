package ecode

import (
	"github.com/zhufuyi/sponge/pkg/errcode"
)

// languages business-level http error codes.
// the languagesNO value range is 1~100, if the same error code is used, it will cause panic.
var (
	languagesNO       = 2
	languagesName     = "languages"
	languagesBaseCode = errcode.HCode(languagesNO)

	ErrCreateLanguages     = errcode.NewError(languagesBaseCode+1, "failed to create "+languagesName)
	ErrDeleteByIDLanguages = errcode.NewError(languagesBaseCode+2, "failed to delete "+languagesName)
	ErrUpdateByIDLanguages = errcode.NewError(languagesBaseCode+3, "failed to update "+languagesName)
	ErrGetByIDLanguages    = errcode.NewError(languagesBaseCode+4, "failed to get "+languagesName+" details")
	ErrListLanguages       = errcode.NewError(languagesBaseCode+5, "failed to list of "+languagesName)

	ErrDeleteByIDsLanguages    = errcode.NewError(languagesBaseCode+6, "failed to delete by batch ids "+languagesName)
	ErrGetByConditionLanguages = errcode.NewError(languagesBaseCode+7, "failed to get "+languagesName+" details by conditions")
	ErrListByIDsLanguages      = errcode.NewError(languagesBaseCode+8, "failed to list by batch ids "+languagesName)
	ErrListByLastIDLanguages   = errcode.NewError(languagesBaseCode+9, "failed to list by last id "+languagesName)

	// error codes are globally unique, adding 1 to the previous error code
)
