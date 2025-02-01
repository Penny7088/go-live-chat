package handler

import (
	"github.com/gin-gonic/gin"
	"lingua_exchange/internal/handler"
	"lingua_exchange/pkg/jwt"
)

func init() {
	apiV1RouterFns = append(apiV1RouterFns, func(r *gin.RouterGroup) {
		groupApplyRouter(r, handler.NewGroupApplyHandler())
	})
}

func groupApplyRouter(group *gin.RouterGroup, h handler.GroupApplyHandler) {
	routerGroup := group.Group("/groupApply")
	routerGroup.Use(jwt.AuthMiddleware())
	routerGroup.POST("/apply/create", h.Create)        // 提交入群申请
	routerGroup.POST("/apply/agree", h.Agree)          // 同意入群申请
	routerGroup.POST("/apply/decline", h.Decline)      // 拒绝入群申请
	routerGroup.GET("/apply/list", h.List)             // 入群申请列表
	routerGroup.GET("/apply/all", h.All)               // 入群申请列表
	routerGroup.GET("/apply/unread", h.ApplyUnreadNum) // 入群申请未读
}
