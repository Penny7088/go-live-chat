package handler

import (
	"github.com/gin-gonic/gin"
	"lingua_exchange/internal/handler"
	"lingua_exchange/pkg/jwt"
)

func init() {
	apiV1RouterFns = append(apiV1RouterFns, func(r *gin.RouterGroup) {
		groupNoticeRouter(r, handler.NewGroupNoticeHandler())
	})
}

func groupNoticeRouter(group *gin.RouterGroup, h handler.GroupNoticeHandler) {
	routerGroup := group.Group("/groupNotice")
	routerGroup.Use(jwt.AuthMiddleware())
	routerGroup.GET("/notice/list", h.List)             // 群公告列表
	routerGroup.POST("/notice/edit", h.CreateAndUpdate) // 添加或编辑群公告
	routerGroup.POST("/notice/delete", h.Delete)        // 删除群公告
}
