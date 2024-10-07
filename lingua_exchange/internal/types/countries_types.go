package types

import (
	"time"

	"github.com/zhufuyi/sponge/pkg/ggorm/query"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.

// CreateCountriesRequest request params
type CreateCountriesRequest struct {
	Name    string `json:"name" binding:""`
	IsoCode string `json:"isoCode" binding:""`
}

// UpdateCountriesByIDRequest request params
type UpdateCountriesByIDRequest struct {
	ID uint64 `json:"id" binding:""` // uint64 id

	Name    string `json:"name" binding:""`
	IsoCode string `json:"isoCode" binding:""`
}

// CountriesObjDetail detail
type CountriesObjDetail struct {
	ID uint64 `json:"id"` // convert to uint64 id

	Name    string `json:"name"`
	IsoCode string `json:"isoCode"`
}

// CreateCountriesReply only for api docs
type CreateCountriesReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ID uint64 `json:"id"` // id
	} `json:"data"` // return data
}

// UpdateCountriesByIDReply only for api docs
type UpdateCountriesByIDReply struct {
	Result
}

// GetCountriesByIDReply only for api docs
type GetCountriesByIDReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Countries CountriesObjDetail `json:"countries"`
	} `json:"data"` // return data
}

// DeleteCountriesByIDReply only for api docs
type DeleteCountriesByIDReply struct {
	Result
}

// DeleteCountriessByIDsReply only for api docs
type DeleteCountriessByIDsReply struct {
	Result
}

// ListCountriessRequest request params
type ListCountriessRequest struct {
	query.Params
}

// ListCountriesReply only for api docs
type ListCountriesReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Countries []CountriesObjDetail `json:"countriess"`
	} `json:"data"` // return data
}

// DeleteCountriessByIDsRequest request params
type DeleteCountriessByIDsRequest struct {
	IDs []uint64 `json:"ids" binding:"min=1"` // id list
}

// GetCountriesByConditionRequest request params
type GetCountriesByConditionRequest struct {
	query.Conditions
}

// GetCountriesByConditionReply only for api docs
type GetCountriesByConditionReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Countries CountriesObjDetail `json:"countries"`
	} `json:"data"` // return data
}

// ListCountriessByIDsRequest request params
type ListCountriessByIDsRequest struct {
	IDs []uint64 `json:"ids" binding:"min=1"` // id list
}

// ListCountriessByIDsReply only for api docs
type ListCountriessByIDsReply struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Countriess []CountriesObjDetail `json:"countriess"`
	} `json:"data"` // return data
}
