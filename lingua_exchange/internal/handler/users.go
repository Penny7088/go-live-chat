package handler

import (
	"context"
	"errors"
	"math"
	"time"

	"cloud.google.com/go/auth/credentials/idtoken"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/redis/go-redis/v9"
	"github.com/zhufuyi/sponge/pkg/errcode"
	"github.com/zhufuyi/sponge/pkg/gin/middleware"
	"github.com/zhufuyi/sponge/pkg/gin/response"
	"github.com/zhufuyi/sponge/pkg/logger"
	"github.com/zhufuyi/sponge/pkg/utils"
	"gorm.io/gorm"
	"lingua_exchange/internal/cache"
	"lingua_exchange/internal/dao"
	"lingua_exchange/internal/ecode"
	"lingua_exchange/internal/model"
	"lingua_exchange/internal/types"
	"lingua_exchange/pkg/encrypt"
	"lingua_exchange/pkg/jwt"
	"lingua_exchange/tools"
)

var _ UsersHandler = (*usersHandler)(nil)

const (
	GooglePlatform = "google"
)

// UsersHandler defining the handler interface
type UsersHandler interface {
	Create(c *gin.Context)
	DeleteByID(c *gin.Context)
	UpdateByID(c *gin.Context)
	GetByID(c *gin.Context)
	List(c *gin.Context)

	DeleteByIDs(c *gin.Context)
	GetByCondition(c *gin.Context)
	ListByIDs(c *gin.Context)
	ListByLastID(c *gin.Context)
	LoginOrRegister(c *gin.Context)
	LoginFromEmail(c *gin.Context)
	SignUpFromEmail(c *gin.Context)
	ResetPassword(c *gin.Context)
	UpdateUserInfoByID(c *gin.Context)
}

type usersHandler struct {
	iDao              dao.UsersDao
	thirdDao          dao.ThirdPartyAuthDao
	deviceDao         dao.UserDevicesDao
	userCache         cache.UsersCache
	globalConfigCache cache.GlobalConfigCache
	userInterestsDao  dao.UserInterestsDao
}

// NewUsersHandler creating the handler interface
func NewUsersHandler() UsersHandler {
	return &usersHandler{
		iDao: dao.NewUsersDao(
			model.GetDB(),
			cache.NewUsersCache(model.GetCacheType()),
		),
		thirdDao:          dao.NewThirdPartyAuthDao(model.GetDB(), cache.NewThirdPartyAuthCache(model.GetCacheType())),
		deviceDao:         dao.NewUserDevicesDao(model.GetDB(), cache.NewUserDevicesCache(model.GetCacheType())),
		userCache:         cache.NewUsersCache(model.GetCacheType()),
		globalConfigCache: cache.NewGlobalConfigCache(model.GetCacheType()),
		userInterestsDao:  dao.NewUserInterestsDao(model.GetDB(), cache.NewUserInterestsCache(model.GetCacheType())),
	}
}

// ResetPassword
// @Summary login from email
// @Description  used email login
// @Tags  Login
// @accept  json
// @Produce json
// @Param data body types.ResetPasswordReq true "Reset Password Information"
// @Success 200 {object} types.ResetPasswordReplay{}
// @Router  /api/v1/users/resetPassword [post]
func (h *usersHandler) ResetPassword(c *gin.Context) {
	form := &types.ResetPasswordReq{}
	if err := c.ShouldBindJSON(form); err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	// 验证验证码
	if err := h.verifyCode(c, form.Email, form.Code, cache.VCodeForgetType); err != nil {
		return
	}
	// 获取用户信息
	user, err := h.iDao.GetByEmail(c, form.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Warn("User not found", middleware.GCtxRequestIDField(c))
			response.Error(c, ecode.ErrUserNotFound)
			return
		}
		logger.Warn("failed to check user existence", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InternalServerError)
		return
	}

	// 更新用户密码
	user.PasswordHash = encrypt.HashPassword(form.NewPassword)
	if err := h.iDao.UpdateByID(c, user); err != nil {
		logger.Warn("failed to update user password", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.ErrUpdateUsers)
		return
	}

	response.Success(c)
}

// SignUpFromEmail
// @Summary login from email
// @Description  used email login
// @Tags  Login
// @accept  json
// @Produce json
// @Param req body types.SignUpFromEmailReq true "用户注册请求体"
// @Success 200 {object} types.LoginReply{}
// @Router  /api/v1/users/signUpFromEmail [post]
func (h *usersHandler) SignUpFromEmail(c *gin.Context) {
	form := &types.SignUpFromEmailReq{}
	if err := c.ShouldBind(form); err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	// 验证验证码
	if err := h.verifyCode(c, form.Email, form.Code, cache.VCodeSignUpType); err != nil {
		return
	}

	db := model.GetDB()
	user := &model.Users{
		Email:         form.Email,
		PasswordHash:  encrypt.HashPassword(form.Password),
		EmailVerified: 1,
	}

	err2 := db.Transaction(func(tx *gorm.DB) error {
		// 检查邮箱是否已注册
		_, err := h.iDao.GetByEmailTx(c, tx, user)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// 邮箱未注册，创建新用户
				if err := h.createUser(c, tx, user); err != nil {
					return err
				}
			} else {
				// 处理查询错误
				return err
			}
		} else {
			return ecode.ErrUserAlreadyExists.Err()
		}

		// 创建或更新设备信息
		device := &model.UserDevices{
			UserID:      int64(user.ID),
			DeviceToken: jwt.HeaderDeviceToken(c),
			DeviceType:  jwt.HeaderPlatform(c),
			IPAddress:   c.ClientIP(),
		}

		if err := h.createOrUpdateDevice(c, tx, device); err != nil {
			return err
		}

		return nil
	})

	if err2 != nil {
		parseError := errcode.ParseError(err2)
		h.handleError(c, err2, parseError)
		return
	}

	token, refreshToken, err2 := h.generateAndCacheToken(user)
	if err2 != nil {
		h.handleError(c, err2, ecode.ErrToken)
		return
	}

	// 成功响应
	data := h.buildUserDetailResponse(user, token, refreshToken, true)

	response.Success(c, gin.H{"user": data})

}

// 创建或更新设备的函数
func (h *usersHandler) createOrUpdateDevice(ctx *gin.Context, tx *gorm.DB, device *model.UserDevices) error {
	_, err := h.deviceDao.FirstOrCreateByTx(ctx, tx, device)
	if err != nil {
		h.handleError(ctx, err, ecode.ErrCreateUserDevices)
		return err
	}
	return nil
}

// 创建用户的函数
func (h *usersHandler) createUser(c *gin.Context, tx *gorm.DB, user *model.Users) error {
	_, err := h.iDao.CreateByTx(c, tx, user)
	if err != nil {
		if ecode.IsUniqueConstraintError(err) {
			h.handleError(c, err, ecode.ErrCreateUsers)
			return err // 返回错误而不是调用
		}
		return err
	}
	return nil // 成功创建用户时返回 nil
}

// verifyCode 验证验证码
func (h *usersHandler) verifyCode(c *gin.Context, email string, code string, codeType string) error {
	codeFromCache, err := h.globalConfigCache.GetVerificationCode(c, email, codeType)
	if err != nil || errors.Is(err, redis.Nil) {
		logger.Warn("verification code error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.ErrVerificationCodeExpired)
		return err
	}

	if codeFromCache != code {
		logger.Warn("invalid verification code", middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.ErrVerificationCode)
		return errors.New("invalid verification code")
	}
	return nil
}

// buildUserDetailResponse  组装用户数据
func (h *usersHandler) buildUserDetailResponse(users *model.Users, token string, refreshToken string, newUser bool) *types.UsersObjDetail {
	data := &types.UsersObjDetail{
		ID:                 users.ID,
		Email:              users.Email,
		ProfilePicture:     users.ProfilePicture,
		EmailVerified:      users.EmailVerified,
		Token:              token,
		RefreshToken:       refreshToken,
		Username:           users.Username,
		LanguageLevel:      users.LanguageLevel,
		CountryID:          users.CountryID,
		NativeLanguageID:   users.NativeLanguageID,
		LearningLanguageID: users.LearningLanguageID,
		Age:                users.Age,
		Gender:             users.Gender,
		RegistrationDate:   users.RegistrationDate,
		IsNewUser:          newUser, // 由于此时用户是新注册的，所以直接赋值为 true。
	}

	// 使用 copier 复制用户信息
	if err := copier.Copy(data, users); err != nil {
		logger.Warn("failed to copy user details", logger.Err(err))
		return nil
	}

	return data
}

// LoginFromEmail
// @Summary login from email
// @Description  used email login
// @Tags  login
// @accept  json
// @Produce json
// @Param req body types.LoginFromEmailReq true "用户登录请求体"
// @Success 200 {object} types.LoginReply{}
// @Router  /api/v1/users/loginFromEmail [post]
func (h *usersHandler) LoginFromEmail(c *gin.Context) {
	form := &types.LoginFromEmailReq{}
	if err := c.ShouldBind(form); err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	ctx := middleware.WrapCtx(c)
	users, err := h.iDao.GetByEmail(ctx, form.Email)
	if errors.Is(err, gorm.ErrRecordNotFound) || err != nil {
		logger.Warn("email error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.ErrEmailNotFound)
		return
	}

	if !encrypt.VerifyPassword(users.PasswordHash, form.Password) {
		logger.Warn("password error: ", middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.ErrPassword)
		return
	}
	token, refreshToken, err := h.generateAndCacheToken(users)
	if err != nil {
		logger.Warn("token gen error: ", middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.ErrToken)
		return
	}

	data := h.buildUserDetailResponse(users, token, refreshToken, false)
	response.Success(c, gin.H{"user": data})

}

// LoginOrRegister
// @Summary login users
// @Description submit information to create users
// @Tags users
// @accept json
// @Produce json
// @Param data body types.LoginRequest true "users information"
// @Success 200 {object} types.LoginReply{}
// @Router /api/v1/users/auth [post]
// @Security BearerAuth
func (h *usersHandler) LoginOrRegister(c *gin.Context) {
	form := &types.LoginRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	switch form.Platform {
	case GooglePlatform:
		if err := h.handleGoogleLogin(c, form); err != nil {
			h.handleError(c, err, ecode.ErrInvalidGoogleIdToken)
		}
	default:
		response.Error(c, ecode.ErrUnsupportedPlatform)
		return
	}

}

func (h *usersHandler) handleGoogleLogin(c *gin.Context, form *types.LoginRequest) error {
	// 验证 Google ID Token
	tokenInfo, err := idtoken.Validate(context.Background(), form.IdToken, "")
	if err != nil {
		return err
	}

	clientIP := c.ClientIP()
	name, _ := tokenInfo.Claims["name"].(string)
	email, _ := tokenInfo.Claims["email"].(string)
	picture, _ := tokenInfo.Claims["picture"].(string)
	emailVerified, _ := tokenInfo.Claims["email_verified"].(bool)

	emailStatus := 0
	if emailVerified {
		emailStatus = 1
	}

	user := &model.Users{
		Email:          email,
		Username:       name,
		ProfilePicture: picture,
		EmailVerified:  emailStatus,
	}

	ctx := middleware.WrapCtx(c)
	db := model.GetDB()

	newUser := false
	var device *model.UserDevices
	err = db.Transaction(func(tx *gorm.DB) error {
		user, err = h.iDao.GetByEmailTx(ctx, tx, user)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 用户不存在，创建新用户
			newUser = true
			_, err := h.iDao.CreateByTx(ctx, tx, user)
			if err != nil {
				if ecode.IsUniqueConstraintError(err) {
					response.Output(c, ecode.ErrUserAlreadyExists.ToHTTPCode())
					return nil
				}
				return err
			}
		} else if err != nil {
			return err
		}

		// 插入第三方认证信息
		if err := h.insertThirdPartyAuth(ctx, tx, user, form); err != nil {
			return err
		}

		// 插入设备信息
		_, err := h.insertUserDevice(ctx, tx, user, form, clientIP)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	if nil == device {
		response.Error(c, ecode.ErrDeviceNotFound)
		return err
	}

	// 生成 Token 和缓存
	token, refreshToken, err := h.generateAndCacheToken(user)
	if err != nil {
		return err
	}

	// 成功响应
	data := h.buildUserDetailResponse(user, token, refreshToken, newUser)
	response.Success(c, gin.H{"user": data})

	return nil
}

func (h *usersHandler) handleError(c *gin.Context, err error, errorCode *errcode.Error) {
	logger.Info("Error occurred", logger.Err(err), middleware.GCtxRequestIDField(c))
	response.Error(c, errorCode)
}

func (h *usersHandler) generateAndCacheToken(user *model.Users) (string, string, error) {

	token, refreshToken, err := jwt.GenerateTokens(user.ID)

	if err != nil {
		return "", "", err
	}

	return token, refreshToken, nil
}

func (h *usersHandler) insertThirdPartyAuth(ctx context.Context, tx *gorm.DB, user *model.Users, form *types.LoginRequest) error {
	thirdPartyAuth, err := h.thirdDao.GetByID(ctx, user.ID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		thirdPartyAuth = &model.ThirdPartyAuth{
			UserID:         int64(user.ID),
			ProviderUserID: form.IdToken,
			Provider:       form.Platform,
		}
		_, err := h.thirdDao.CreateByTx(ctx, tx, thirdPartyAuth)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	return nil
}

func (h *usersHandler) insertUserDevice(ctx context.Context, tx *gorm.DB, user *model.Users, form *types.LoginRequest, clientIP string) (*model.UserDevices, error) {
	token := form.DeviceToken
	deviceToken := tools.GenerateDeviceToken(token)

	device := &model.UserDevices{
		UserID:      int64(user.ID),
		DeviceToken: deviceToken,
		DeviceType:  form.Platform,
		IPAddress:   clientIP,
	}
	userDevices, err := h.deviceDao.FirstOrCreateByTx(ctx, tx, device)
	if err != nil {
		return userDevices, err
	}
	return userDevices, nil
}

// Create a record
// @Summary create users
// @Description submit information to create users
// @Tags users
// @accept json
// @Produce json
// @Param data body types.CreateUsersRequest true "users information"
// @Success 200 {object} types.CreateUsersReply{}
// @Router /api/v1/users [post]
// @Security BearerAuth
func (h *usersHandler) Create(c *gin.Context) {
	form := &types.CreateUsersRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	users := &model.Users{}
	err = copier.Copy(users, form)
	if err != nil {
		response.Error(c, ecode.ErrCreateUsers)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.Create(ctx, users)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{"id": users.ID})
}

// DeleteByID delete a record by id
// @Summary delete users
// @Description delete users by id
// @Tags users
// @accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} types.DeleteUsersByIDReply{}
// @Router /api/v1/users/{id} [delete]
// @Security BearerAuth
func (h *usersHandler) DeleteByID(c *gin.Context) {
	_, id, isAbort := getUsersIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	err := h.iDao.DeleteByID(ctx, id)
	if err != nil {
		logger.Error("DeleteByID error", logger.Err(err), logger.Any("id", id), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// UpdateUserInfoByID  update user info
// @Summary  update user info
// @Description  update user info
// @Tags users
// @accept json
// @Produce json
// @Param Authorization header string true "Authorization"
// @Param platform header string true "platform - ios/android"
// @Param deviceToken header string true "deviceToken device id"
// @Param id path string true "id"
// @Param data body types.UpdateUsersByIDRequest true "users information"
// @Success 200 {object} types.GetUsersByIDReply{}
// @Router /api/v1/users/updateUserInfo/{id} [put]
// @Security BearerAuth
func (h *usersHandler) UpdateUserInfoByID(c *gin.Context) {
	_, id, isAbort := getUsersIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}
	form := &types.UpdateUsersByIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	form.ID = id

	users := &model.Users{}
	err = copier.Copy(users, form)
	if err != nil {
		response.Error(c, ecode.ErrUpdateByIDUsers)
		return
	}
	users.RegistrationDate = time.Now()
	ctx := middleware.WrapCtx(c)
	db := model.GetDB()
	err = db.Transaction(func(db *gorm.DB) error {
		err2 := h.iDao.UpdateByTx(ctx, db, users)
		if err2 != nil {
			logger.Error("UpdateByID error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
			return err2
		}

		if len(form.Interests) > 0 {
			interests := form.Interests
			for _, tagId := range interests {
				userInterests := &model.UserInterests{
					UserID: int64(form.ID),
					TagID:  tagId,
				}

				if err := h.userInterestsDao.CreateByTx(ctx, db, userInterests); err != nil {
					if errors.Is(err, gorm.ErrDuplicatedKey) {
						continue
					} else {
						return err
					}
				}

			}
		}

		return nil
	})

	if err != nil {
		h.handleError(c, err, ecode.ErrUpdateByIDUsers)
	}

	response.Success(c, gin.H{"user": users})

}

// UpdateByID update information by id
// @Summary update users
// @Description update users information by id
// @Tags users
// @accept json
// @Produce json
// @Param id path string true "id"
// @Param data body types.UpdateUsersByIDRequest true "users information"
// @Success 200 {object} types.UpdateUsersByIDReply{}
// @Router /api/v1/users/{id} [put]
// @Security BearerAuth
func (h *usersHandler) UpdateByID(c *gin.Context) {
	_, id, isAbort := getUsersIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	form := &types.UpdateUsersByIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	form.ID = id

	users := &model.Users{}
	err = copier.Copy(users, form)
	if err != nil {
		response.Error(c, ecode.ErrUpdateByIDUsers)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.UpdateByID(ctx, users)
	if err != nil {
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// GetByID get a record by id
// @Summary get users detail
// @Description get users detail by id
// @Tags users
// @Param id path string true "id"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetUsersByIDReply{}
// @Router /api/v1/users/{id} [get]
// @Security BearerAuth
func (h *usersHandler) GetByID(c *gin.Context) {
	_, id, isAbort := getUsersIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	users, err := h.iDao.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, model.ErrRecordNotFound) {
			logger.Warn("GetByID not found", logger.Err(err), logger.Any("id", id), middleware.GCtxRequestIDField(c))
			response.Error(c, ecode.NotFound)
		} else {
			logger.Error("GetByID error", logger.Err(err), logger.Any("id", id), middleware.GCtxRequestIDField(c))
			response.Output(c, ecode.InternalServerError.ToHTTPCode())
		}
		return
	}

	data := &types.UsersObjDetail{}
	err = copier.Copy(data, users)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDUsers)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"users": data})
}

// List of records by query parameters
// @Summary list of userss by query parameters
// @Description list of userss by paging and conditions
// @Tags users
// @accept json
// @Produce json
// @Param data body types.Params true "query parameters"
// @Success 200 {object} types.ListUserssReply{}
// @Router /api/v1/users/list [post]
// @Security BearerAuth
func (h *usersHandler) List(c *gin.Context) {
	form := &types.ListUserssRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	userss, total, err := h.iDao.GetByColumns(ctx, &form.Params)
	if err != nil {
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertUserss(userss)
	if err != nil {
		response.Error(c, ecode.ErrListUsers)
		return
	}

	response.Success(c, gin.H{
		"userss": data,
		"total":  total,
	})
}

// DeleteByIDs delete records by batch id
// @Summary delete userss
// @Description delete userss by batch id
// @Tags users
// @Param data body types.DeleteUserssByIDsRequest true "id array"
// @Accept json
// @Produce json
// @Success 200 {object} types.DeleteUserssByIDsReply{}
// @Router /api/v1/users/delete/ids [post]
// @Security BearerAuth
func (h *usersHandler) DeleteByIDs(c *gin.Context) {
	form := &types.DeleteUserssByIDsRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	err = h.iDao.DeleteByIDs(ctx, form.IDs)
	if err != nil {
		logger.Error("GetByIDs error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// GetByCondition get a record by condition
// @Summary get users by condition
// @Description get users by condition
// @Tags users
// @Param data body types.Conditions true "query condition"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetUsersByConditionReply{}
// @Router /api/v1/users/condition [post]
// @Security BearerAuth
func (h *usersHandler) GetByCondition(c *gin.Context) {
	form := &types.GetUsersByConditionRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	err = form.Conditions.CheckValid()
	if err != nil {
		logger.Warn("Parameters error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	users, err := h.iDao.GetByCondition(ctx, &form.Conditions)
	if err != nil {
		if errors.Is(err, model.ErrRecordNotFound) {
			logger.Warn("GetByCondition not found", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
			response.Error(c, ecode.NotFound)
		} else {
			logger.Error("GetByCondition error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
			response.Output(c, ecode.InternalServerError.ToHTTPCode())
		}
		return
	}

	data := &types.UsersObjDetail{}
	err = copier.Copy(data, users)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDUsers)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"users": data})
}

// ListByIDs list of records by batch id
// @Summary list of userss by batch id
// @Description list of userss by batch id
// @Tags users
// @Param data body types.ListUserssByIDsRequest true "id array"
// @Accept json
// @Produce json
// @Success 200 {object} types.ListUserssByIDsReply{}
// @Router /api/v1/users/list/ids [post]
// @Security BearerAuth
func (h *usersHandler) ListByIDs(c *gin.Context) {
	form := &types.ListUserssByIDsRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	usersMap, err := h.iDao.GetByIDs(ctx, form.IDs)
	if err != nil {
		logger.Error("GetByIDs error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	userss := []*types.UsersObjDetail{}
	for _, id := range form.IDs {
		if v, ok := usersMap[id]; ok {
			record, err := convertUsers(v)
			if err != nil {
				response.Error(c, ecode.ErrListUsers)
				return
			}
			userss = append(userss, record)
		}
	}

	response.Success(c, gin.H{
		"userss": userss,
	})
}

// ListByLastID get records by last id and limit
// @Summary list of userss by last id and limit
// @Description list of userss by last id and limit
// @Tags users
// @accept json
// @Produce json
// @Param lastID query int true "last id, default is MaxInt32" default(0)
// @Param limit query int false "number per page" default(10)
// @Param sort query string false "sort by column name of table, and the "-" sign before column name indicates reverse order" default(-id)
// @Success 200 {object} types.ListUserssReply{}
// @Router /api/v1/users/list [get]
// @Security BearerAuth
func (h *usersHandler) ListByLastID(c *gin.Context) {
	lastID := utils.StrToUint64(c.Query("lastID"))
	if lastID == 0 {
		lastID = math.MaxInt32
	}
	limit := utils.StrToInt(c.Query("limit"))
	if limit == 0 {
		limit = 10
	}
	sort := c.Query("sort")

	ctx := middleware.WrapCtx(c)
	userss, err := h.iDao.GetByLastID(ctx, lastID, limit, sort)
	if err != nil {
		logger.Error("GetByLastID error", logger.Err(err), logger.Uint64("latsID", lastID), logger.Int("limit", limit), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertUserss(userss)
	if err != nil {
		response.Error(c, ecode.ErrListByLastIDUsers)
		return
	}

	response.Success(c, gin.H{
		"userss": data,
	})
}

func getUsersIDFromPath(c *gin.Context) (string, uint64, bool) {
	idStr := c.Param("id")
	id, err := utils.StrToUint64E(idStr)
	if err != nil || id == 0 {
		logger.Warn("StrToUint64E error: ", logger.String("idStr", idStr), middleware.GCtxRequestIDField(c))
		return "", 0, true
	}

	return idStr, id, false
}

func convertUsers(users *model.Users) (*types.UsersObjDetail, error) {
	data := &types.UsersObjDetail{}
	err := copier.Copy(data, users)
	if err != nil {
		return nil, err
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	return data, nil
}

func convertUserss(fromValues []*model.Users) ([]*types.UsersObjDetail, error) {
	toValues := []*types.UsersObjDetail{}
	for _, v := range fromValues {
		data, err := convertUsers(v)
		if err != nil {
			return nil, err
		}
		toValues = append(toValues, data)
	}

	return toValues, nil
}
