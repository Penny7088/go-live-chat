package handler

import "github.com/gin-gonic/gin"

var _ GroupApplyHandler = (*groupApplyHandler)(nil)

type GroupApplyHandler interface {
	Create(ctx *gin.Context)
	Agree(ctx *gin.Context)
	Decline(ctx *gin.Context)
	List(ctx *gin.Context)
	All(ctx *gin.Context)
	ApplyUnreadNum(ctx *gin.Context)
}

type groupApplyHandler struct {
}

func NewGroupApplyHandler() GroupApplyHandler {
	return &groupApplyHandler{}
}

func (g groupApplyHandler) Create(ctx *gin.Context) {
	panic("implement me")
}

func (g groupApplyHandler) Agree(ctx *gin.Context) {
	// TODO implement me
	panic("implement me")
}

func (g groupApplyHandler) Decline(ctx *gin.Context) {
	// TODO implement me
	panic("implement me")
}

func (g groupApplyHandler) List(ctx *gin.Context) {
	// TODO implement me
	panic("implement me")
}

func (g groupApplyHandler) All(ctx *gin.Context) {
	// TODO implement me
	panic("implement me")
}

func (g groupApplyHandler) ApplyUnreadNum(ctx *gin.Context) {
	// TODO implement me
	panic("implement me")
}
