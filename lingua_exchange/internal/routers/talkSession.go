package handler

import (
	"github.com/gin-gonic/gin"
	"lingua_exchange/internal/handler"
	"lingua_exchange/pkg/jwt"
)

func init() {
	apiV1RouterFns = append(apiV1RouterFns, func(r *gin.RouterGroup) {
		sessionRouter(r, handler.NewSessionHandler())
	})
}

func sessionRouter(group *gin.RouterGroup, h handler.SessionHandler) {
	g := group.Group("/session")
	g.Use(jwt.AuthMiddleware())
	g.GET("/list", h.SessionList)
	g.POST("/create", h.Create)
	g.POST("/delete", h.Delete)
	g.POST("/topping", h.Top)
	g.POST("/disturb", h.Disturb)
	g.POST("/unread/clear", h.ClearUnreadMessage)
	g.POST("/records", h.GetRecords)
}
