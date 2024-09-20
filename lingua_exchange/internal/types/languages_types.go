package types

import (
	"time"

	"github.com/zhufuyi/sponge/pkg/ggorm/query"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.

// CreateLanguagesRequest request params
type CreateLanguagesRequest struct {
	Name       string `json:"name" binding:""`
	NativeName string `json:"nativeName" binding:""`
	IsoCode    string `json:"isoCode" binding:""`
}

// UpdateLanguagesByIDRequest request params
type UpdateLanguagesByIDRequest struct {
	ID uint64 `json:"id" binding:""` // uint64 id

	Name       string `json:"name" binding:""`
	NativeName string `json:"nativeName" binding:""`
	IsoCode    string `json:"isoCode" binding:""`
}

// LanguagesObjDetail detail
type LanguagesObjDetail struct {
	ID uint64 `json:"id"` // convert to uint64 id

	Name       string `json:"name"`
	NativeName string `json:"nativeName"`
	IsoCode    string `json:"isoCode"`
}

// CreateLanguagesReply only for api docs
type CreateLanguagesReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ID uint64 `json:"id"` // id
	} `json:"data"` // return data
}

// UpdateLanguagesByIDReply only for api docs
type UpdateLanguagesByIDReply struct {
	Result
}

// GetLanguagesByIDReply only for api docs
type GetLanguagesByIDReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Languages LanguagesObjDetail `json:"languages"`
	} `json:"data"` // return data
}

// DeleteLanguagesByIDReply only for api docs
type DeleteLanguagesByIDReply struct {
	Result
}

// DeleteLanguagessByIDsReply only for api docs
type DeleteLanguagessByIDsReply struct {
	Result
}

// ListLanguagessRequest request params
type ListLanguagessRequest struct {
	query.Params
}

// ListLanguagessReply only for api docs
type ListLanguagessReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Languagess []LanguagesObjDetail `json:"languagess"`
	} `json:"data"` // return data
}

// DeleteLanguagessByIDsRequest request params
type DeleteLanguagessByIDsRequest struct {
	IDs []uint64 `json:"ids" binding:"min=1"` // id list
}

// GetLanguagesByConditionRequest request params
type GetLanguagesByConditionRequest struct {
	query.Conditions
}

// GetLanguagesByConditionReply only for api docs
type GetLanguagesByConditionReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Languages LanguagesObjDetail `json:"languages"`
	} `json:"data"` // return data
}

// ListLanguagessByIDsRequest request params
type ListLanguagessByIDsRequest struct {
	IDs []uint64 `json:"ids" binding:"min=1"` // id list
}

// ListLanguagessByIDsReply only for api docs
type ListLanguagessByIDsReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Languagess []LanguagesObjDetail `json:"languagess"`
	} `json:"data"` // return data
}
