package handler

import (
	"github.com/gin-gonic/gin"
	"lingua_exchange/internal/handler"
	"lingua_exchange/pkg/jwt"
	"lingua_exchange/pkg/socket"
)

// todo 路由需要升级至websocket
func init() {
	apiWSRouterFns = append(apiWSRouterFns, func(r *gin.RouterGroup) {
		messageRouter(r, handler.NewMessageHandler())
	})
}

func messageRouter(group *gin.RouterGroup, h handler.MessageHandler) {
	routerGroup := group.Group("im")
	routerGroup.GET("/connect/detail", func(ctx *gin.Context) {
		ctx.JSON(200, map[string]any{
			"chat":    socket.Session.Chat.Count(),
			"example": socket.Session.Example.Count(),
		})
	})

	routerGroup.GET("/chat.io", jwt.AuthMiddleware(), h.Connection)
}
