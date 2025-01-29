package ecode

import (
	"github.com/zhufuyi/sponge/pkg/errcode"
)

// group business-level http error codes.
// the groupNO value range is 1~100, if the same error code is used, it will cause panic.
var (
	groupNO       = 60
	groupName     = "group"
	groupBaseCode = errcode.HCode(groupNO)

	ErrCreateGroup     = errcode.NewError(groupBaseCode+1, "failed to create "+groupName)
	ErrDeleteByIDGroup = errcode.NewError(groupBaseCode+2, "failed to delete "+groupName)
	ErrUpdateByIDGroup = errcode.NewError(groupBaseCode+3, "failed to update "+groupName)
	ErrGetByIDGroup    = errcode.NewError(groupBaseCode+4, "failed to get "+groupName+" details")
	ErrListGroup       = errcode.NewError(groupBaseCode+5, "failed to list of "+groupName)

	ErrDeleteByIDsGroup    = errcode.NewError(groupBaseCode+6, "failed to delete by batch ids "+groupName)
	ErrGetByConditionGroup = errcode.NewError(groupBaseCode+7, "failed to get "+groupName+" details by conditions")
	ErrListByIDsGroup      = errcode.NewError(groupBaseCode+8, "failed to list by batch ids "+groupName)
	ErrListByLastIDGroup   = errcode.NewError(groupBaseCode+9, "failed to list by last id "+groupName)
	ErrGroupDismiss        = errcode.NewError(groupBaseCode+10, "group is dismiss "+groupName)

	// error codes are globally unique, adding 1 to the previous error code
)
