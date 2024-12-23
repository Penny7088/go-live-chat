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

	ErrDeleteByIDsUsers        = errcode.NewError(usersBaseCode+6, "failed to delete by batch ids "+usersName)
	ErrGetByConditionUsers     = errcode.NewError(usersBaseCode+7, "failed to get "+usersName+" details by conditions")
	ErrListByIDsUsers          = errcode.NewError(usersBaseCode+8, "failed to list by batch ids "+usersName)
	ErrListByLastIDUsers       = errcode.NewError(usersBaseCode+9, "failed to list by last id "+usersName)
	ErrInvalidGoogleIdToken    = errcode.NewError(usersBaseCode+10, "invalid Google ID Token")
	ErrUserAlreadyExists       = errcode.NewError(usersBaseCode+11, "User Already Exists")
	ErrUnsupportedPlatform     = errcode.NewError(usersBaseCode+12, "Unsupported platform")
	ErrToken                   = errcode.NewError(usersBaseCode+13, "gen Token error")
	ErrEmailNotFound           = errcode.NewError(usersBaseCode+14, "email not found")
	ErrPassword                = errcode.NewError(usersBaseCode+15, "password error")
	ErrVerificationCode        = errcode.NewError(usersBaseCode+16, "verification code invalidate")
	ErrVerificationCodeExpired = errcode.NewError(usersBaseCode+17, "verification code is Expired")
	ErrUserNotFound            = errcode.NewError(usersBaseCode+18, "user not found")
	ErrUpdateUsers             = errcode.NewError(usersBaseCode+19, "update user info is err")

	// error codes are globally unique, adding 1 to the previous error code
)
