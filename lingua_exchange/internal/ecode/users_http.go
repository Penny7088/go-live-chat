package ecode

import (
	"github.com/zhufuyi/sponge/pkg/errcode"
)

// users business-level http error codes.
// the usersNO value range is 1~100, if the same error code is used, it will cause panic.
var (
	usersNO       = 3
	usersName     = "users"
	usersBaseCode = errcode.HCode(usersNO)

	ErrCreateUsers     = errcode.NewError(usersBaseCode+1, "failed to create "+usersName)
	ErrDeleteByIDUsers = errcode.NewError(usersBaseCode+2, "failed to delete "+usersName)
	ErrUpdateByIDUsers = errcode.NewError(usersBaseCode+3, "failed to update "+usersName)
	ErrGetByIDUsers    = errcode.NewError(usersBaseCode+4, "failed to get "+usersName+" details")
	ErrListUsers       = errcode.NewError(usersBaseCode+5, "failed to list of "+usersName)

	ErrDeleteByIDsUsers     = errcode.NewError(usersBaseCode+6, "failed to delete by batch ids "+usersName)
	ErrGetByConditionUsers  = errcode.NewError(usersBaseCode+7, "failed to get "+usersName+" details by conditions")
	ErrListByIDsUsers       = errcode.NewError(usersBaseCode+8, "failed to list by batch ids "+usersName)
	ErrListByLastIDUsers    = errcode.NewError(usersBaseCode+9, "failed to list by last id "+usersName)
	ErrInvalidGoogleIdToken = errcode.NewError(usersBaseCode+10, "invalid Google ID Token")
	ErrUserAlreadyExists    = errcode.NewError(usersBaseCode+11, "User Already Exists")
	ErrUnsupportedPlatform  = errcode.NewError(usersBaseCode+12, "Unsupported platform")
	ErrToken                = errcode.NewError(usersBaseCode+13, "gen Token error")

	// error codes are globally unique, adding 1 to the previous error code
)
