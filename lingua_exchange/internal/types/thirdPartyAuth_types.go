package types

import (
	"time"

	"github.com/zhufuyi/sponge/pkg/ggorm/query"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.

// CreateThirdPartyAuthRequest request params
type CreateThirdPartyAuthRequest struct {
	UserID         int64  `json:"userID" binding:""`
	Provider       string `json:"provider" binding:""`
	ProviderUserID string `json:"providerUserID" binding:""`
}

// UpdateThirdPartyAuthByIDRequest request params
type UpdateThirdPartyAuthByIDRequest struct {
	ID uint64 `json:"id" binding:""` // uint64 id

	UserID         int64  `json:"userID" binding:""`
	Provider       string `json:"provider" binding:""`
	ProviderUserID string `json:"providerUserID" binding:""`
}

// ThirdPartyAuthObjDetail detail
type ThirdPartyAuthObjDetail struct {
	ID uint64 `json:"id"` // convert to uint64 id

	UserID         int64  `json:"userID"`
	Provider       string `json:"provider"`
	ProviderUserID string `json:"providerUserID"`
}

// CreateThirdPartyAuthReply only for api docs
type CreateThirdPartyAuthReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ID uint64 `json:"id"` // id
	} `json:"data"` // return data
}

// UpdateThirdPartyAuthByIDReply only for api docs
type UpdateThirdPartyAuthByIDReply struct {
	Result
}

// GetThirdPartyAuthByIDReply only for api docs
type GetThirdPartyAuthByIDReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ThirdPartyAuth ThirdPartyAuthObjDetail `json:"thirdPartyAuth"`
	} `json:"data"` // return data
}

// DeleteThirdPartyAuthByIDReply only for api docs
type DeleteThirdPartyAuthByIDReply struct {
	Result
}

// DeleteThirdPartyAuthsByIDsReply only for api docs
type DeleteThirdPartyAuthsByIDsReply struct {
	Result
}

// ListThirdPartyAuthsRequest request params
type ListThirdPartyAuthsRequest struct {
	query.Params
}

// ListThirdPartyAuthsReply only for api docs
type ListThirdPartyAuthsReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ThirdPartyAuths []ThirdPartyAuthObjDetail `json:"thirdPartyAuths"`
	} `json:"data"` // return data
}

// DeleteThirdPartyAuthsByIDsRequest request params
type DeleteThirdPartyAuthsByIDsRequest struct {
	IDs []uint64 `json:"ids" binding:"min=1"` // id list
}

// GetThirdPartyAuthByConditionRequest request params
type GetThirdPartyAuthByConditionRequest struct {
	query.Conditions
}

// GetThirdPartyAuthByConditionReply only for api docs
type GetThirdPartyAuthByConditionReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ThirdPartyAuth ThirdPartyAuthObjDetail `json:"thirdPartyAuth"`
	} `json:"data"` // return data
}

// ListThirdPartyAuthsByIDsRequest request params
type ListThirdPartyAuthsByIDsRequest struct {
	IDs []uint64 `json:"ids" binding:"min=1"` // id list
}

// ListThirdPartyAuthsByIDsReply only for api docs
type ListThirdPartyAuthsByIDsReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ThirdPartyAuths []ThirdPartyAuthObjDetail `json:"thirdPartyAuths"`
	} `json:"data"` // return data
}
