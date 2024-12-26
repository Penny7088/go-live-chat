package ecode

import (
	"github.com/zhufuyi/sponge/pkg/errcode"
)

// userInterests business-level http error codes.
// the userInterestsNO value range is 1~100, if the same error code is used, it will cause panic.
var (
	userInterestsNO       = 1
	userInterestsName     = "userInterests"
	userInterestsBaseCode = errcode.HCode(userInterestsNO)

	ErrCreateUserInterests     = errcode.NewError(userInterestsBaseCode+1, "failed to create "+userInterestsName)
	ErrDeleteByIDUserInterests = errcode.NewError(userInterestsBaseCode+2, "failed to delete "+userInterestsName)
	ErrUpdateByIDUserInterests = errcode.NewError(userInterestsBaseCode+3, "failed to update "+userInterestsName)
	ErrGetByIDUserInterests    = errcode.NewError(userInterestsBaseCode+4, "failed to get "+userInterestsName+" details")
	ErrListUserInterests       = errcode.NewError(userInterestsBaseCode+5, "failed to list of "+userInterestsName)

	ErrDeleteByIDsUserInterests    = errcode.NewError(userInterestsBaseCode+6, "failed to delete by batch ids "+userInterestsName)
	ErrGetByConditionUserInterests = errcode.NewError(userInterestsBaseCode+7, "failed to get "+userInterestsName+" details by conditions")
	ErrListByIDsUserInterests      = errcode.NewError(userInterestsBaseCode+8, "failed to list by batch ids "+userInterestsName)
	ErrListByLastIDUserInterests   = errcode.NewError(userInterestsBaseCode+9, "failed to list by last id "+userInterestsName)

	// error codes are globally unique, adding 1 to the previous error code
)
