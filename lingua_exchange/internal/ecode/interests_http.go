package ecode

import (
	"github.com/zhufuyi/sponge/pkg/errcode"
)

// interests business-level http error codes.
// the interestsNO value range is 1~100, if the same error code is used, it will cause panic.
var (
	interestsNO       = 6
	interestsName     = "interests"
	interestsBaseCode = errcode.HCode(interestsNO)

	ErrCreateInterests     = errcode.NewError(interestsBaseCode+1, "failed to create "+interestsName)
	ErrDeleteByIDInterests = errcode.NewError(interestsBaseCode+2, "failed to delete "+interestsName)
	ErrUpdateByIDInterests = errcode.NewError(interestsBaseCode+3, "failed to update "+interestsName)
	ErrGetByIDInterests    = errcode.NewError(interestsBaseCode+4, "failed to get "+interestsName+" details")
	ErrListInterests       = errcode.NewError(interestsBaseCode+5, "failed to list of "+interestsName)

	ErrDeleteByIDsInterests    = errcode.NewError(interestsBaseCode+6, "failed to delete by batch ids "+interestsName)
	ErrGetByConditionInterests = errcode.NewError(interestsBaseCode+7, "failed to get "+interestsName+" details by conditions")
	ErrListByIDsInterests      = errcode.NewError(interestsBaseCode+8, "failed to list by batch ids "+interestsName)
	ErrListByLastIDInterests   = errcode.NewError(interestsBaseCode+9, "failed to list by last id "+interestsName)
	ErrLanguageCodeNull        = errcode.NewError(interestsBaseCode+10, "language code is nil "+interestsName)

	// error codes are globally unique, adding 1 to the previous error code
)
