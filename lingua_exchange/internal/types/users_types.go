package types

import (
	"time"

	"github.com/zhufuyi/sponge/pkg/ggorm/query"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.

type LoginRequest struct {
	IdToken     string `json:"idToken" binding:"required"`
	Platform    string `json:"platform" binding:"required"`
	DeviceType  string `json:"deviceType" binding:"required"`
	DeviceToken string `json:"deviceToken" binding:"required"`
}

// CreateUsersRequest request params
type CreateUsersRequest struct {
	Email              string    `json:"email" binding:""`
	Username           string    `json:"username" binding:""`
	PasswordHash       string    `json:"passwordHash" binding:""`
	ProfilePicture     string    `json:"profilePicture" binding:""`
	NativeLanguageID   int64     `json:"nativeLanguageID" binding:""`
	LearningLanguageID int64     `json:"learningLanguageID" binding:""`
	LanguageLevel      string    `json:"languageLevel" binding:""`
	Age                int       `json:"age" binding:""`
	Gender             string    `json:"gender" binding:""`
	Interests          string    `json:"interests" binding:""`
	CountryID          int64     `json:"countryID" binding:""`
	RegistrationDate   time.Time `json:"registrationDate" binding:""`
	LastLogin          time.Time `json:"lastLogin" binding:""`
	Status             string    `json:"status" binding:""`
	EmailVerified      int       `json:"emailVerified" binding:""`
	VerificationToken  string    `json:"verificationToken" binding:""`
	TokenExpiration    time.Time `json:"tokenExpiration" binding:""`
}

// UpdateUsersByIDRequest request params
type UpdateUsersByIDRequest struct {
	ID uint64 `json:"id" binding:""` // uint64 id

	Email              string    `json:"email" binding:""`
	Username           string    `json:"username" binding:""`
	PasswordHash       string    `json:"passwordHash" binding:""`
	ProfilePicture     string    `json:"profilePicture" binding:""`
	NativeLanguageID   int64     `json:"nativeLanguageID" binding:""`
	LearningLanguageID int64     `json:"learningLanguageID" binding:""`
	LanguageLevel      string    `json:"languageLevel" binding:""`
	Age                int       `json:"age" binding:""`
	Gender             string    `json:"gender" binding:""`
	Interests          string    `json:"interests" binding:""`
	CountryID          int64     `json:"countryID" binding:""`
	RegistrationDate   time.Time `json:"registrationDate" binding:""`
	LastLogin          time.Time `json:"lastLogin" binding:""`
	Status             string    `json:"status" binding:""`
	EmailVerified      int       `json:"emailVerified" binding:""`
	VerificationToken  string    `json:"verificationToken" binding:""`
	TokenExpiration    time.Time `json:"tokenExpiration" binding:""`
}

// UsersObjDetail detail
type UsersObjDetail struct {
	ID uint64 `json:"id"` // convert to uint64 id

	Email              string    `json:"email"`
	Username           string    `json:"username"`
	PasswordHash       string    `json:"passwordHash"`
	ProfilePicture     string    `json:"profilePicture"`
	NativeLanguageID   int64     `json:"nativeLanguageID"`
	LearningLanguageID int64     `json:"learningLanguageID"`
	LanguageLevel      string    `json:"languageLevel"`
	Age                int       `json:"age"`
	Gender             string    `json:"gender"`
	Interests          string    `json:"interests"`
	CountryID          int64     `json:"countryID"`
	RegistrationDate   time.Time `json:"registrationDate"`
	LastLogin          time.Time `json:"lastLogin"`
	Status             string    `json:"status"`
	EmailVerified      int       `json:"emailVerified"`
	VerificationToken  string    `json:"verificationToken"`
	TokenExpiration    time.Time `json:"tokenExpiration"`
}

// LoginReply only for api docs
type LoginReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ID uint64 `json:"id"` // id
	} `json:"data"` // return data
}

// CreateUsersReply only for api docs
type CreateUsersReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ID uint64 `json:"id"` // id
	} `json:"data"` // return data
}

// UpdateUsersByIDReply only for api docs
type UpdateUsersByIDReply struct {
	Result
}

// GetUsersByIDReply only for api docs
type GetUsersByIDReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Users UsersObjDetail `json:"users"`
	} `json:"data"` // return data
}

// DeleteUsersByIDReply only for api docs
type DeleteUsersByIDReply struct {
	Result
}

// DeleteUserssByIDsReply only for api docs
type DeleteUserssByIDsReply struct {
	Result
}

// ListUserssRequest request params
type ListUserssRequest struct {
	query.Params
}

// ListUserssReply only for api docs
type ListUserssReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Userss []UsersObjDetail `json:"userss"`
	} `json:"data"` // return data
}

// DeleteUserssByIDsRequest request params
type DeleteUserssByIDsRequest struct {
	IDs []uint64 `json:"ids" binding:"min=1"` // id list
}

// GetUsersByConditionRequest request params
type GetUsersByConditionRequest struct {
	query.Conditions
}

// GetUsersByConditionReply only for api docs
type GetUsersByConditionReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Users UsersObjDetail `json:"users"`
	} `json:"data"` // return data
}

// ListUserssByIDsRequest request params
type ListUserssByIDsRequest struct {
	IDs []uint64 `json:"ids" binding:"min=1"` // id list
}

// ListUserssByIDsReply only for api docs
type ListUserssByIDsReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Userss []UsersObjDetail `json:"userss"`
	} `json:"data"` // return data
}
