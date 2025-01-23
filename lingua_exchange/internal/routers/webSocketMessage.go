package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"lingua_exchange/internal/handler"
	verify "lingua_exchange/pkg/jwt"
	"lingua_exchange/pkg/socket"
)

func NewWebSocketRouter() *gin.Engine {
	r := gin.New()

	r.Use(gin.RecoveryWithWriter(gin.DefaultWriter, func(c *gin.Context, err any) {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, map[string]any{"code": 500, "msg": "系统错误，请重试!!!"})
	}))
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, map[string]any{"msg": "请求地址不存在"})
	})

	messageRouter(r, handler.NewMessageHandler())

	return r
}

func messageRouter(group *gin.Engine, h handler.MessageHandler) {
	routerGroup := group.Group("ws")
	routerGroup.GET("/connect/detail", func(ctx *gin.Context) {
		ctx.JSON(200, map[string]any{
			"chat":    socket.Session.Chat.Count(),
			"example": socket.Session.Example.Count(),
		})
	})

	routerGroup.GET("/chat.io", verify.AuthWSMiddleware(), h.Connection)
}
