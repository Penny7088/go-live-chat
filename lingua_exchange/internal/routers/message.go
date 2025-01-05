package handler

import (
	"github.com/gin-gonic/gin"
	"lingua_exchange/pkg/jwt"
	"lingua_exchange/pkg/socket"
)

// todo 路由需要升级至websocket
func init() {
	apiWSRouterFns = append(apiWSRouterFns, func(r *gin.RouterGroup) {

	})
}

func messageRouter(group *gin.RouterGroup) {

	group.GET("/connect/detail", func(ctx *gin.Context) {
		ctx.JSON(200, map[string]any{
			"chat":    socket.Session.Chat.Count(),
			"example": socket.Session.Example.Count(),
		})
	})

	group.GET("/chat.io", jwt.AuthMiddleware())
}
