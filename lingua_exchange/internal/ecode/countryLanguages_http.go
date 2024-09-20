package ecode

import (
	"github.com/zhufuyi/sponge/pkg/errcode"
)

// countryLanguages business-level http error codes.
// the countryLanguagesNO value range is 1~100, if the same error code is used, it will cause panic.
var (
	countryLanguagesNO       = 1
	countryLanguagesName     = "countryLanguages"
	countryLanguagesBaseCode = errcode.HCode(countryLanguagesNO)

	ErrCreateCountryLanguages     = errcode.NewError(countryLanguagesBaseCode+1, "failed to create "+countryLanguagesName)
	ErrDeleteByIDCountryLanguages = errcode.NewError(countryLanguagesBaseCode+2, "failed to delete "+countryLanguagesName)
	ErrUpdateByIDCountryLanguages = errcode.NewError(countryLanguagesBaseCode+3, "failed to update "+countryLanguagesName)
	ErrGetByIDCountryLanguages    = errcode.NewError(countryLanguagesBaseCode+4, "failed to get "+countryLanguagesName+" details")
	ErrListCountryLanguages       = errcode.NewError(countryLanguagesBaseCode+5, "failed to list of "+countryLanguagesName)

	ErrDeleteByIDsCountryLanguages    = errcode.NewError(countryLanguagesBaseCode+6, "failed to delete by batch ids "+countryLanguagesName)
	ErrGetByConditionCountryLanguages = errcode.NewError(countryLanguagesBaseCode+7, "failed to get "+countryLanguagesName+" details by conditions")
	ErrListByIDsCountryLanguages      = errcode.NewError(countryLanguagesBaseCode+8, "failed to list by batch ids "+countryLanguagesName)
	ErrListByLastIDCountryLanguages   = errcode.NewError(countryLanguagesBaseCode+9, "failed to list by last id "+countryLanguagesName)

	// error codes are globally unique, adding 1 to the previous error code
)
