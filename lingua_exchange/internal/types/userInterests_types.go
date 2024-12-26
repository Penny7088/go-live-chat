package types

import (
	"time"

	"github.com/zhufuyi/sponge/pkg/ggorm/query"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.

// CreateUserInterestsRequest request params
type CreateUserInterestsRequest struct {
	UserID int64 `json:"userID" binding:""`
	TagID  int64 `json:"tagID" binding:""`
}

// UpdateUserInterestsByIDRequest request params
type UpdateUserInterestsByIDRequest struct {
	UserID int64 `json:"userID" binding:""`
	TagID  int64 `json:"tagID" binding:""`
}

// UserInterestsObjDetail detail
type UserInterestsObjDetail struct {
	UserID    int64     `json:"userID"`
	TagID     int64     `json:"tagID"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// CreateUserInterestsReply only for api docs
type CreateUserInterestsReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ID uint64 `json:"id"` // id
	} `json:"data"` // return data
}

// UpdateUserInterestsByIDReply only for api docs
type UpdateUserInterestsByIDReply struct {
	Result
}

// GetUserInterestsByIDReply only for api docs
type GetUserInterestsByIDReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		UserInterests UserInterestsObjDetail `json:"userInterests"`
	} `json:"data"` // return data
}

// DeleteUserInterestsByIDReply only for api docs
type DeleteUserInterestsByIDReply struct {
	Result
}

// DeleteUserInterestssByIDsReply only for api docs
type DeleteUserInterestssByIDsReply struct {
	Result
}

// ListUserInterestssRequest request params
type ListUserInterestssRequest struct {
	query.Params
}

// ListUserInterestssReply only for api docs
type ListUserInterestssReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		UserInterestss []UserInterestsObjDetail `json:"userInterestss"`
	} `json:"data"` // return data
}

// DeleteUserInterestssByIDsRequest request params
type DeleteUserInterestssByIDsRequest struct {
	IDs []uint64 `json:"ids" binding:"min=1"` // id list
}

// GetUserInterestsByConditionRequest request params
type GetUserInterestsByConditionRequest struct {
	query.Conditions
}

// GetUserInterestsByConditionReply only for api docs
type GetUserInterestsByConditionReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		UserInterests UserInterestsObjDetail `json:"userInterests"`
	} `json:"data"` // return data
}

// ListUserInterestssByIDsRequest request params
type ListUserInterestssByIDsRequest struct {
	IDs []uint64 `json:"ids" binding:"min=1"` // id list
}

// ListUserInterestssByIDsReply only for api docs
type ListUserInterestssByIDsReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		UserInterestss []UserInterestsObjDetail `json:"userInterestss"`
	} `json:"data"` // return data
}
