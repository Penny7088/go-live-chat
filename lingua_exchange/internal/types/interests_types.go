package types

import (
	"time"

	"github.com/zhufuyi/sponge/pkg/ggorm/query"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.

// CreateInterestsRequest request params
type CreateInterestsRequest struct {
	TagID   int64  `json:"tagID" binding:""`
	TagName string `json:"tagName" binding:""`
}

// UpdateInterestsByIDRequest request params
type UpdateInterestsByIDRequest struct {
	TagID   int64  `json:"tagID" binding:""`
	TagName string `json:"tagName" binding:""`
}

// InterestsObjDetail detail
type InterestsObjDetail struct {
	TagID     int64     `json:"tagID"`
	TagName   string    `json:"tagName"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type InterestTranslationDetail struct {
	TagID          int64  `json:"tagID"`
	TranslatedName string `json:"translatedName"`
}

// CreateInterestsReply only for api docs
type CreateInterestsReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ID uint64 `json:"id"` // id
	} `json:"data"` // return data
}

// UpdateInterestsByIDReply only for api docs
type UpdateInterestsByIDReply struct {
	Result
}

// GetInterestsByIDReply only for api docs
type GetInterestsByIDReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Interests InterestsObjDetail `json:"interests"`
	} `json:"data"` // return data
}

// DeleteInterestsByIDReply only for api docs
type DeleteInterestsByIDReply struct {
	Result
}

// DeleteInterestssByIDsReply only for api docs
type DeleteInterestssByIDsReply struct {
	Result
}

// ListInterestssRequest request params
type ListInterestssRequest struct {
	query.Params
}

// ListInterestssReply only for api docs
type ListInterestssReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Interestss []InterestsObjDetail `json:"interestss"`
	} `json:"data"` // return data
}

// DeleteInterestssByIDsRequest request params
type DeleteInterestssByIDsRequest struct {
	IDs []uint64 `json:"ids" binding:"min=1"` // id list
}

// GetInterestsByConditionRequest request params
type GetInterestsByConditionRequest struct {
	query.Conditions
}

// GetInterestsByConditionReply only for api docs
type GetInterestsByConditionReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Interests InterestsObjDetail `json:"interests"`
	} `json:"data"` // return data
}

// GetInterestsByLanguageReply only for api docs
type GetInterestsByLanguageReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Interests InterestTranslationDetail `json:"interests"`
	} `json:"data"` // return data
}

// ListInterestssByIDsRequest request params
type ListInterestssByIDsRequest struct {
	IDs []uint64 `json:"ids" binding:"min=1"` // id list
}

// ListInterestssByIDsReply only for api docs
type ListInterestssByIDsReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Interestss []InterestsObjDetail `json:"interestss"`
	} `json:"data"` // return data
}
