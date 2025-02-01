package handler

import (
	"github.com/gin-gonic/gin"
	"lingua_exchange/internal/handler"
	"lingua_exchange/pkg/jwt"
)

func init() {
	apiV1RouterFns = append(apiV1RouterFns, func(r *gin.RouterGroup) {
		groupRouter(r, handler.NewGroupHandler())
	})

}

func groupRouter(group *gin.RouterGroup, h handler.GroupHandler) {
	routerGroup := group.Group("/group")
	routerGroup.Use(jwt.AuthMiddleware())
	routerGroup.GET("/list", h.GroupList)            // 群组列表
	routerGroup.GET("/overt/list", h.OvertList)      // 公开群组列表
	routerGroup.GET("/detail", h.Detail)             // 群组详情
	routerGroup.POST("/create", h.Create)            // 创建群组
	routerGroup.POST("/dismiss", h.Dismiss)          // 解散群组
	routerGroup.POST("/invite", h.Invite)            // 邀请加入群组
	routerGroup.POST("/secede", h.SignOut)           // 退出群组
	routerGroup.POST("/setting", h.Setting)          // 设置群组信息
	routerGroup.POST("/handover", h.Handover)        // 群主转让
	routerGroup.POST("/assign-admin", h.AssignAdmin) // 分配管理员
	routerGroup.POST("/no-speak", h.NoSpeak)         // 修改禁言状态
	routerGroup.POST("/mute", h.Mute)                // 修改禁言状态
	routerGroup.POST("/overt", h.Overt)              // 修改禁言状态
}
