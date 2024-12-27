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

	g.GET("/", h.GlobalConfig)                                //  [get] api/v1/globalConfig/loginMethod
	g.POST("/sendSignUpVerifyCode", h.SendSignUpVerifyCode)   // [post] api/v1/globalConfig/sendSignUpVerifyCode
	g.POST("/sendResetPasswordCode", h.SendResetPasswordCode) // [post] api/v1/globalConfig/sendResetPasswordCode
}
