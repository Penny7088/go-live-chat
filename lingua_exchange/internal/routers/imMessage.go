package handler

import (
	"github.com/gin-gonic/gin"
	"lingua_exchange/internal/handler"
	"lingua_exchange/pkg/jwt"
)

func init() {
	apiV1RouterFns = append(apiV1RouterFns, func(r *gin.RouterGroup) {
		imMessageRouter(r, handler.NewIMMessageHandler())
	})
}

func imMessageRouter(group *gin.RouterGroup, h handler.IMMessageHandler) {
	routerGroup := group.Group("/message")
	routerGroup.Use(jwt.AuthMiddleware())
	routerGroup.POST("/publish", h.Publish)
}
