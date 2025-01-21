package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"lingua_exchange/internal/handler"
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
	// r.Use(middleware.Cors())
	//
	// if config.Get().HTTP.Timeout > 0 {
	// 	// if you need more fine-grained control over your routes, set the timeout in your routes, unsetting the timeout globally here.
	// 	r.Use(middleware.Timeout(time.Second * time.Duration(config.Get().HTTP.Timeout)))
	// }

	// request id middleware
	// r.Use(middleware.RequestID())

	// logger middleware, to print simple messages, replace middleware.Logging with middleware.SimpleLog
	// r.Use(middleware.Logging(
	// 	middleware.WithLog(logger.Get()),
	// 	middleware.WithRequestIDFromContext(),
	// 	middleware.WithIgnoreRoutes("/metrics"), // ignore path
	// ))

	// init jwt middleware
	// jwt.Init(
	// 	jwt.WithExpire(verify.UserTokenExpireTime),
	// 	jwt.WithSigningKey("live_lingua:"),
	// 	jwt.WithSigningMethod(jwt.HS384),
	// )

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

	// routerGroup.GET("/chat.io", verify.AuthWSMiddleware(), h.Connection)
	routerGroup.GET("/chat.io", h.Connection)
}
