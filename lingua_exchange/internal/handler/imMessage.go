package handler

import (
	"github.com/gin-gonic/gin"
)

// 消息路由实现类
var _ IMMessageHandler = (*imMessageHandler)(nil)

type IMMessageHandler interface {
	Publish(ctx *gin.Context)
}

type imMessageHandler struct {
}

func NewIMMessageHandler() IMMessageHandler {
	return &imMessageHandler{}
}

func (i imMessageHandler) Publish(ctx *gin.Context) {
	// TODO implement me
	panic("implement me")
}
