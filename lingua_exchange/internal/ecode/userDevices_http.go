package ecode

import (
	"github.com/zhufuyi/sponge/pkg/errcode"
)

// userDevices business-level http error codes.
// the userDevicesNO value range is 1~100, if the same error code is used, it will cause panic.
var (
	userDevicesNO       = 5
	userDevicesName     = "userDevices"
	userDevicesBaseCode = errcode.HCode(userDevicesNO)

	ErrCreateUserDevices     = errcode.NewError(userDevicesBaseCode+1, "failed to create "+userDevicesName)
	ErrDeleteByIDUserDevices = errcode.NewError(userDevicesBaseCode+2, "failed to delete "+userDevicesName)
	ErrUpdateByIDUserDevices = errcode.NewError(userDevicesBaseCode+3, "failed to update "+userDevicesName)
	ErrGetByIDUserDevices    = errcode.NewError(userDevicesBaseCode+4, "failed to get "+userDevicesName+" details")
	ErrListUserDevices       = errcode.NewError(userDevicesBaseCode+5, "failed to list of "+userDevicesName)

	ErrDeleteByIDsUserDevices    = errcode.NewError(userDevicesBaseCode+6, "failed to delete by batch ids "+userDevicesName)
	ErrGetByConditionUserDevices = errcode.NewError(userDevicesBaseCode+7, "failed to get "+userDevicesName+" details by conditions")
	ErrListByIDsUserDevices      = errcode.NewError(userDevicesBaseCode+8, "failed to list by batch ids "+userDevicesName)
	ErrListByLastIDUserDevices   = errcode.NewError(userDevicesBaseCode+9, "failed to list by last id "+userDevicesName)

	// error codes are globally unique, adding 1 to the previous error code
)
