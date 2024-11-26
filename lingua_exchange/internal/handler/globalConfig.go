package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/zhufuyi/sponge/pkg/errcode"
	"github.com/zhufuyi/sponge/pkg/gin/middleware"
	"github.com/zhufuyi/sponge/pkg/gin/response"
	"github.com/zhufuyi/sponge/pkg/logger"
	"lingua_exchange/internal/cache"
	"lingua_exchange/internal/ecode"
	"lingua_exchange/internal/model"
	"lingua_exchange/internal/types"
	"lingua_exchange/pkg/emailtool"
	"lingua_exchange/pkg/ip"
	"lingua_exchange/pkg/strutil"
)

type GlobalConfigHandler interface {
	LoginMethod(c *gin.Context)
	SendResetPasswordCode(c *gin.Context)
	SendSignUpVerifyCode(c *gin.Context)
}

type globalConfigHandler struct {
	cache cache.GlobalConfigCache
}

func NewGlobalConfigHandler() GlobalConfigHandler {
	return &globalConfigHandler{
		cache: cache.NewGlobalConfigCache(model.GetCacheType()),
	}
}

// SendSignUpVerifyCode
// @Summary  发送注册验证码
// @Description  发送验证码
// @Tags  验证码
// @accept      json
// @Param req body types.VerificationCodeReq true "Request payload containing email"
// @Success 200 {object} types.SignUpVerificationCodeRely
// @Router /api/v1/globalConfig/sendSignUpVerifyCode [post]
func (g globalConfigHandler) SendSignUpVerifyCode(c *gin.Context) {
	g.sendVerificationCode(c, cache.VCodeSignUpType, "register.html", "Your Sign Up code")
}

// SendResetPasswordCode
// @Summary  发送验证码，重置密码
// @Description  发送验证码，重置密码
// @Tags  验证码
// @accept      json
// @Param req body types.VerificationCodeReq true "Request payload containing email"
// @Success 200 {object} types.ResetVerificationCodeRely
// @Router /api/v1/globalConfig/sendResetPasswordCode [post]
func (g globalConfigHandler) SendResetPasswordCode(c *gin.Context) {
	g.sendVerificationCode(c, cache.VCodeForgetType, "reset_pwd.html", "Your Reset Password code")
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

// sendVerificationCode 发送验证代码
func (g globalConfigHandler) sendVerificationCode(c *gin.Context, codeType string, templatePath string, subject string) {
	req := &types.VerificationCodeReq{}
	if err := c.ShouldBindJSON(req); err != nil {
		g.handleValidationError(c, err, ecode.InvalidParams)
		return
	}

	validateCode := strutil.GenValidateCode(6)

	if err := g.sendEmail(req.Email, validateCode, subject, templatePath, c); err != nil {
		return
	}

	if err := g.storeVerificationCode(c, req.Email, validateCode, codeType); err != nil {
		return
	}

	response.Success(c)
}

// handleValidationError 处理请求验证错误
func (g globalConfigHandler) handleValidationError(c *gin.Context, err error, errorCode *errcode.Error) {
	logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
	response.Error(c, errorCode)
}

// sendEmail 发送电子邮件
func (g globalConfigHandler) sendEmail(email string, code string, subject string, templatePath string, c *gin.Context) error {
	if err := emailtool.SendEmail(email, code, subject, templatePath); err != nil {
		logger.Warn("send Code error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.ErrSendCode)
		return err
	}
	return nil
}

// storeVerificationCode 存储验证代码到缓存
func (g globalConfigHandler) storeVerificationCode(c *gin.Context, email string, code string, codeType string) error {
	if err := g.cache.SetVerificationCode(c, email, code, codeType); err != nil {
		logger.Warn("storage cache validate Code error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.ErrSendCode)
		return err
	}
	return nil
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
