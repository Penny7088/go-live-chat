package handler

import (
	"github.com/gin-gonic/gin"
	"lingua_exchange/internal/handler"
)

func init() {
	apiV1RouterFns = append(apiV1RouterFns, func(group *gin.RouterGroup) {
		globalConfigRouter(group, handler.NewGlobalConfigHandler())
	})
}

func globalConfigRouter(group *gin.RouterGroup, h handler.GlobalConfigHandler) {
	g := group.Group("/globalConfig")

	g.GET("/LoginMethod", h.LoginMethod) // [get] /api/v1/countries/list
}
