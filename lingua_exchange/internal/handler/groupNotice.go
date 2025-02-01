package handler

import "github.com/gin-gonic/gin"

var _ GroupNoticeHandler = (*groupNoticeHandler)(nil)

type GroupNoticeHandler interface {
	List(ctx *gin.Context)
	CreateAndUpdate(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

type groupNoticeHandler struct {
}

func NewGroupNoticeHandler() GroupNoticeHandler {
	return &groupNoticeHandler{}
}

func (g groupNoticeHandler) List(ctx *gin.Context) {
	panic("implement me")
}

func (g groupNoticeHandler) CreateAndUpdate(ctx *gin.Context) {
	// TODO implement me
	panic("implement me")
}

func (g groupNoticeHandler) Delete(ctx *gin.Context) {
	// TODO implement me
	panic("implement me")
}
