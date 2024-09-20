package types

import (
	"time"

	"github.com/zhufuyi/sponge/pkg/ggorm/query"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.

// CreateCountryLanguagesRequest request params
type CreateCountryLanguagesRequest struct {
	CountryID  int64 `json:"countryID" binding:""`
	LanguageID int64 `json:"languageID" binding:""`
}

// UpdateCountryLanguagesByIDRequest request params
type UpdateCountryLanguagesByIDRequest struct {
	ID uint64 `json:"id" binding:""` // uint64 id

	CountryID  int64 `json:"countryID" binding:""`
	LanguageID int64 `json:"languageID" binding:""`
}

// CountryLanguagesObjDetail detail
type CountryLanguagesObjDetail struct {
	ID uint64 `json:"id"` // convert to uint64 id

	CountryID  int64 `json:"countryID"`
	LanguageID int64 `json:"languageID"`
}

// CreateCountryLanguagesReply only for api docs
type CreateCountryLanguagesReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ID uint64 `json:"id"` // id
	} `json:"data"` // return data
}

// UpdateCountryLanguagesByIDReply only for api docs
type UpdateCountryLanguagesByIDReply struct {
	Result
}

// GetCountryLanguagesByIDReply only for api docs
type GetCountryLanguagesByIDReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		CountryLanguages CountryLanguagesObjDetail `json:"countryLanguages"`
	} `json:"data"` // return data
}

// DeleteCountryLanguagesByIDReply only for api docs
type DeleteCountryLanguagesByIDReply struct {
	Result
}

// DeleteCountryLanguagessByIDsReply only for api docs
type DeleteCountryLanguagessByIDsReply struct {
	Result
}

// ListCountryLanguagessRequest request params
type ListCountryLanguagessRequest struct {
	query.Params
}

// ListCountryLanguagessReply only for api docs
type ListCountryLanguagessReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		CountryLanguagess []CountryLanguagesObjDetail `json:"countryLanguagess"`
	} `json:"data"` // return data
}

// DeleteCountryLanguagessByIDsRequest request params
type DeleteCountryLanguagessByIDsRequest struct {
	IDs []uint64 `json:"ids" binding:"min=1"` // id list
}

// GetCountryLanguagesByConditionRequest request params
type GetCountryLanguagesByConditionRequest struct {
	query.Conditions
}

// GetCountryLanguagesByConditionReply only for api docs
type GetCountryLanguagesByConditionReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		CountryLanguages CountryLanguagesObjDetail `json:"countryLanguages"`
	} `json:"data"` // return data
}

// ListCountryLanguagessByIDsRequest request params
type ListCountryLanguagessByIDsRequest struct {
	IDs []uint64 `json:"ids" binding:"min=1"` // id list
}

// ListCountryLanguagessByIDsReply only for api docs
type ListCountryLanguagessByIDsReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		CountryLanguagess []CountryLanguagesObjDetail `json:"countryLanguagess"`
	} `json:"data"` // return data
}
