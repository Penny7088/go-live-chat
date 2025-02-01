package handler

import (
	"github.com/gin-gonic/gin"
	"lingua_exchange/internal/handler"
	"lingua_exchange/pkg/jwt"
)

func init() {
	apiV1RouterFns = append(apiV1RouterFns, func(r *gin.RouterGroup) {
		groupMemberRouter(r, handler.NewGroupMemberHandler())
	})
}

func groupMemberRouter(group *gin.RouterGroup, h handler.GroupMemberHandler) {
	routerGroup := group.Group("/groupMember")
	routerGroup.Use(jwt.AuthMiddleware())
	routerGroup.GET("/list", h.Members)               // 群成员列表
	routerGroup.GET("/invites", h.GetInviteFriends)   // 群成员列表
	routerGroup.POST("/remove", h.RemoveMembers)      // 移出指定群成员
	routerGroup.POST("/remark", h.UpdateMemberRemark) // 设置群名片
}
