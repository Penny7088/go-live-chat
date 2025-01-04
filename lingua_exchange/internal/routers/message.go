package handler

import (
	"github.com/gin-gonic/gin"
)

// todo 路由需要升级至websocket
func init() {
	apiWSRouterFns = append(apiWSRouterFns, func(r *gin.RouterGroup) {

	})
}

func messageRouter(group *gin.RouterGroup) {

}
