package handler

import (
	"cloud.google.com/go/auth/credentials/idtoken"
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/zhufuyi/sponge/pkg/gin/middleware"
	"github.com/zhufuyi/sponge/pkg/gin/response"
	"github.com/zhufuyi/sponge/pkg/jwt"
	"github.com/zhufuyi/sponge/pkg/logger"
	"github.com/zhufuyi/sponge/pkg/utils"
	"gorm.io/gorm"
	"lingua_exchange/internal/cache"
	"lingua_exchange/internal/dao"
	"lingua_exchange/internal/ecode"
	"lingua_exchange/internal/model"
	"lingua_exchange/internal/types"
	"lingua_exchange/tools"
	"math"
	"time"
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
}

type usersHandler struct {
	iDao      dao.UsersDao
	thirdDao  dao.ThirdPartyAuthDao
	deviceDao dao.UserDevicesDao
	userCache cache.UsersCache
}

// NewUsersHandler creating the handler interface
func NewUsersHandler() UsersHandler {
	return &usersHandler{
		iDao: dao.NewUsersDao(
			model.GetDB(),
			cache.NewUsersCache(model.GetCacheType()),
		),
		thirdDao:  dao.NewThirdPartyAuthDao(model.GetDB(), cache.NewThirdPartyAuthCache(model.GetCacheType())),
		deviceDao: dao.NewUserDevicesDao(model.GetDB(), cache.NewUserDevicesCache(model.GetCacheType())),
		userCache: cache.NewUsersCache(model.GetCacheType()),
	}
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
			h.handleError(c, ecode.ErrInvalidGoogleIdToken.ToHTTPCode(), err)
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
	token, err := h.generateAndCacheToken(ctx, user)
	if err != nil {
		return err
	}

	// 成功响应
	data := &types.UsersObjDetail{
		ID:               user.ID,
		Email:            user.Email,
		ProfilePicture:   user.ProfilePicture,
		EmailVerified:    emailStatus,
		Platform:         form.Platform,
		DeviceToken:      device.DeviceToken,
		Token:            token,
		Username:         user.Username,
		RegistrationDate: time.Now(),
		IsNewUser:        newUser,
	}
	if err := copier.Copy(data, user); err != nil {
		response.Error(c, ecode.ErrGetByIDUsers)
		return nil
	}
	response.Success(c, gin.H{"userInfo": data})

	return nil
}

func (h *usersHandler) handleError(c *gin.Context, code int, err error) {
	logger.Error("Error occurred", logger.Err(err), middleware.GCtxRequestIDField(c))
	response.Output(c, code)
}

func (h *usersHandler) generateAndCacheToken(ctx context.Context, user *model.Users) (string, error) {

	token, err := jwt.GenerateToken(utils.Uint64ToStr(user.ID))
	if err != nil {
		return "", err
	}

	err = h.userCache.SetUserToken(ctx, user.ID, token, cache.UserTokenExpireTime)
	if err != nil {
		return "", err
	}

	return token, nil
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
