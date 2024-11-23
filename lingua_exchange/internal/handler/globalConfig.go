package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/zhufuyi/sponge/pkg/gin/middleware"
	"github.com/zhufuyi/sponge/pkg/gin/response"
	"github.com/zhufuyi/sponge/pkg/logger"
	"lingua_exchange/internal/ecode"
	"lingua_exchange/internal/types"
	"lingua_exchange/pkg/ip"
)

type GlobalConfigHandler interface {
	LoginMethod(c *gin.Context)
}

type globalConfigHandler struct {
}

func NewGlobalConfigHandler() GlobalConfigHandler {
	return &globalConfigHandler{}
}

// LoginMethod  obtain login method
// @Summary get user login method
// @Description  Get different login methods based on the user's IP
// @Tags    globalConfig
// @accept  json
// @Produce json
// @Success 200 {object} types.LoginMethodReply{}
// @Router /api/v1/globalConfig/LoginMethod [get]
func (g globalConfigHandler) LoginMethod(c *gin.Context) {
	clientIP := c.ClientIP()
	if clientIP == "" {
		logger.Warn("ip is nil  error: ", middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.ErrIpNotFound)
	}

	var methods []*types.LoginMethodDetailReply
	if ip.IsIpFromChina(clientIP) {
		methods = append(methods, queryLoginMethodFromCH())
	} else {
		methods = append(methods, queryLoginMethodFromOther())
	}

	response.Success(c, gin.H{
		"loginMethods": methods,
	})
}

// need query config
func queryLoginMethodFromCH() *types.LoginMethodDetailReply {
	data := &types.LoginMethodDetailReply{}
	data.Name = "email"
	return data
}

// need query config
func queryLoginMethodFromOther() *types.LoginMethodDetailReply {
	data := &types.LoginMethodDetailReply{}
	data.Name = "google"
	return data
}
