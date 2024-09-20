package ecode

import (
	"github.com/zhufuyi/sponge/pkg/errcode"
)

// countries business-level http error codes.
// the countriesNO value range is 1~100, if the same error code is used, it will cause panic.
var (
	countriesNO       = 78
	countriesName     = "countries"
	countriesBaseCode = errcode.HCode(countriesNO)

	ErrCreateCountries     = errcode.NewError(countriesBaseCode+1, "failed to create "+countriesName)
	ErrDeleteByIDCountries = errcode.NewError(countriesBaseCode+2, "failed to delete "+countriesName)
	ErrUpdateByIDCountries = errcode.NewError(countriesBaseCode+3, "failed to update "+countriesName)
	ErrGetByIDCountries    = errcode.NewError(countriesBaseCode+4, "failed to get "+countriesName+" details")
	ErrListCountries       = errcode.NewError(countriesBaseCode+5, "failed to list of "+countriesName)

	ErrDeleteByIDsCountries    = errcode.NewError(countriesBaseCode+6, "failed to delete by batch ids "+countriesName)
	ErrGetByConditionCountries = errcode.NewError(countriesBaseCode+7, "failed to get "+countriesName+" details by conditions")
	ErrListByIDsCountries      = errcode.NewError(countriesBaseCode+8, "failed to list by batch ids "+countriesName)
	ErrListByLastIDCountries   = errcode.NewError(countriesBaseCode+9, "failed to list by last id "+countriesName)

	// error codes are globally unique, adding 1 to the previous error code
)
