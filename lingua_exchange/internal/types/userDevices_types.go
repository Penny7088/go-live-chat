package types

import (
	"time"

	"github.com/zhufuyi/sponge/pkg/ggorm/query"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.

// CreateUserDevicesRequest request params
type CreateUserDevicesRequest struct {
	UserID      int64     `json:"userID" binding:""`
	DeviceToken string    `json:"deviceToken" binding:""`
	DeviceType  string    `json:"deviceType" binding:""`
	IPAddress   string    `json:"iPAddress" binding:""`
	LastActive  time.Time `json:"lastActive" binding:""`
}

// UpdateUserDevicesByIDRequest request params
type UpdateUserDevicesByIDRequest struct {
	ID uint64 `json:"id" binding:""` // uint64 id

	UserID      int64     `json:"userID" binding:""`
	DeviceToken string    `json:"deviceToken" binding:""`
	DeviceType  string    `json:"deviceType" binding:""`
	IPAddress   string    `json:"iPAddress" binding:""`
	LastActive  time.Time `json:"lastActive" binding:""`
}

// UserDevicesObjDetail detail
type UserDevicesObjDetail struct {
	ID uint64 `json:"id"` // convert to uint64 id

	UserID      int64     `json:"userID"`
	DeviceToken string    `json:"deviceToken"`
	DeviceType  string    `json:"deviceType"`
	IPAddress   string    `json:"iPAddress"`
	LastActive  time.Time `json:"lastActive"`
}

// CreateUserDevicesReply only for api docs
type CreateUserDevicesReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ID uint64 `json:"id"` // id
	} `json:"data"` // return data
}

// UpdateUserDevicesByIDReply only for api docs
type UpdateUserDevicesByIDReply struct {
	Result
}

// GetUserDevicesByIDReply only for api docs
type GetUserDevicesByIDReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		UserDevices UserDevicesObjDetail `json:"userDevices"`
	} `json:"data"` // return data
}

// DeleteUserDevicesByIDReply only for api docs
type DeleteUserDevicesByIDReply struct {
	Result
}

// DeleteUserDevicessByIDsReply only for api docs
type DeleteUserDevicessByIDsReply struct {
	Result
}

// ListUserDevicessRequest request params
type ListUserDevicessRequest struct {
	query.Params
}

// ListUserDevicessReply only for api docs
type ListUserDevicessReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		UserDevicess []UserDevicesObjDetail `json:"userDevicess"`
	} `json:"data"` // return data
}

// DeleteUserDevicessByIDsRequest request params
type DeleteUserDevicessByIDsRequest struct {
	IDs []uint64 `json:"ids" binding:"min=1"` // id list
}

// GetUserDevicesByConditionRequest request params
type GetUserDevicesByConditionRequest struct {
	query.Conditions
}

// GetUserDevicesByConditionReply only for api docs
type GetUserDevicesByConditionReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		UserDevices UserDevicesObjDetail `json:"userDevices"`
	} `json:"data"` // return data
}

// ListUserDevicessByIDsRequest request params
type ListUserDevicessByIDsRequest struct {
	IDs []uint64 `json:"ids" binding:"min=1"` // id list
}

// ListUserDevicessByIDsReply only for api docs
type ListUserDevicessByIDsReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		UserDevicess []UserDevicesObjDetail `json:"userDevicess"`
	} `json:"data"` // return data
}
