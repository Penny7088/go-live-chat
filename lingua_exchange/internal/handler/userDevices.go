package handler

import (
	"errors"
	"math"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"

	"github.com/zhufuyi/sponge/pkg/gin/middleware"
	"github.com/zhufuyi/sponge/pkg/gin/response"
	"github.com/zhufuyi/sponge/pkg/logger"
	"github.com/zhufuyi/sponge/pkg/utils"

	"lingua_exchange/internal/cache"
	"lingua_exchange/internal/dao"
	"lingua_exchange/internal/ecode"
	"lingua_exchange/internal/model"
	"lingua_exchange/internal/types"
)

var _ UserDevicesHandler = (*userDevicesHandler)(nil)

// UserDevicesHandler defining the handler interface
type UserDevicesHandler interface {
	Create(c *gin.Context)
	DeleteByID(c *gin.Context)
	UpdateByID(c *gin.Context)
	GetByID(c *gin.Context)
	List(c *gin.Context)

	DeleteByIDs(c *gin.Context)
	GetByCondition(c *gin.Context)
	ListByIDs(c *gin.Context)
	ListByLastID(c *gin.Context)
}

type userDevicesHandler struct {
	iDao dao.UserDevicesDao
}

// NewUserDevicesHandler creating the handler interface
func NewUserDevicesHandler() UserDevicesHandler {
	return &userDevicesHandler{
		iDao: dao.NewUserDevicesDao(
			model.GetDB(),
			cache.NewUserDevicesCache(model.GetCacheType()),
		),
	}
}

// Create a record
// @Summary create userDevices
// @Description submit information to create userDevices
// @Tags userDevices
// @accept json
// @Produce json
// @Param data body types.CreateUserDevicesRequest true "userDevices information"
// @Success 200 {object} types.CreateUserDevicesReply{}
// @Router /api/v1/userDevices [post]
// @Security BearerAuth
func (h *userDevicesHandler) Create(c *gin.Context) {
	form := &types.CreateUserDevicesRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	userDevices := &model.UserDevices{}
	err = copier.Copy(userDevices, form)
	if err != nil {
		response.Error(c, ecode.ErrCreateUserDevices)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.Create(ctx, userDevices)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{"id": userDevices.ID})
}

// DeleteByID delete a record by id
// @Summary delete userDevices
// @Description delete userDevices by id
// @Tags userDevices
// @accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} types.DeleteUserDevicesByIDReply{}
// @Router /api/v1/userDevices/{id} [delete]
// @Security BearerAuth
func (h *userDevicesHandler) DeleteByID(c *gin.Context) {
	_, id, isAbort := getUserDevicesIDFromPath(c)
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
// @Summary update userDevices
// @Description update userDevices information by id
// @Tags userDevices
// @accept json
// @Produce json
// @Param id path string true "id"
// @Param data body types.UpdateUserDevicesByIDRequest true "userDevices information"
// @Success 200 {object} types.UpdateUserDevicesByIDReply{}
// @Router /api/v1/userDevices/{id} [put]
// @Security BearerAuth
func (h *userDevicesHandler) UpdateByID(c *gin.Context) {
	_, id, isAbort := getUserDevicesIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	form := &types.UpdateUserDevicesByIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	form.ID = id

	userDevices := &model.UserDevices{}
	err = copier.Copy(userDevices, form)
	if err != nil {
		response.Error(c, ecode.ErrUpdateByIDUserDevices)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.UpdateByID(ctx, userDevices)
	if err != nil {
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// GetByID get a record by id
// @Summary get userDevices detail
// @Description get userDevices detail by id
// @Tags userDevices
// @Param id path string true "id"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetUserDevicesByIDReply{}
// @Router /api/v1/userDevices/{id} [get]
// @Security BearerAuth
func (h *userDevicesHandler) GetByID(c *gin.Context) {
	_, id, isAbort := getUserDevicesIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	userDevices, err := h.iDao.GetByID(ctx, id)
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

	data := &types.UserDevicesObjDetail{}
	err = copier.Copy(data, userDevices)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDUserDevices)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"userDevices": data})
}

// List of records by query parameters
// @Summary list of userDevicess by query parameters
// @Description list of userDevicess by paging and conditions
// @Tags userDevices
// @accept json
// @Produce json
// @Param data body types.Params true "query parameters"
// @Success 200 {object} types.ListUserDevicessReply{}
// @Router /api/v1/userDevices/list [post]
// @Security BearerAuth
func (h *userDevicesHandler) List(c *gin.Context) {
	form := &types.ListUserDevicessRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	userDevicess, total, err := h.iDao.GetByColumns(ctx, &form.Params)
	if err != nil {
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertUserDevicess(userDevicess)
	if err != nil {
		response.Error(c, ecode.ErrListUserDevices)
		return
	}

	response.Success(c, gin.H{
		"userDevicess": data,
		"total":        total,
	})
}

// DeleteByIDs delete records by batch id
// @Summary delete userDevicess
// @Description delete userDevicess by batch id
// @Tags userDevices
// @Param data body types.DeleteUserDevicessByIDsRequest true "id array"
// @Accept json
// @Produce json
// @Success 200 {object} types.DeleteUserDevicessByIDsReply{}
// @Router /api/v1/userDevices/delete/ids [post]
// @Security BearerAuth
func (h *userDevicesHandler) DeleteByIDs(c *gin.Context) {
	form := &types.DeleteUserDevicessByIDsRequest{}
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
// @Summary get userDevices by condition
// @Description get userDevices by condition
// @Tags userDevices
// @Param data body types.Conditions true "query condition"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetUserDevicesByConditionReply{}
// @Router /api/v1/userDevices/condition [post]
// @Security BearerAuth
func (h *userDevicesHandler) GetByCondition(c *gin.Context) {
	form := &types.GetUserDevicesByConditionRequest{}
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
	userDevices, err := h.iDao.GetByCondition(ctx, &form.Conditions)
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

	data := &types.UserDevicesObjDetail{}
	err = copier.Copy(data, userDevices)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDUserDevices)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"userDevices": data})
}

// ListByIDs list of records by batch id
// @Summary list of userDevicess by batch id
// @Description list of userDevicess by batch id
// @Tags userDevices
// @Param data body types.ListUserDevicessByIDsRequest true "id array"
// @Accept json
// @Produce json
// @Success 200 {object} types.ListUserDevicessByIDsReply{}
// @Router /api/v1/userDevices/list/ids [post]
// @Security BearerAuth
func (h *userDevicesHandler) ListByIDs(c *gin.Context) {
	form := &types.ListUserDevicessByIDsRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	userDevicesMap, err := h.iDao.GetByIDs(ctx, form.IDs)
	if err != nil {
		logger.Error("GetByIDs error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	userDevicess := []*types.UserDevicesObjDetail{}
	for _, id := range form.IDs {
		if v, ok := userDevicesMap[id]; ok {
			record, err := convertUserDevices(v)
			if err != nil {
				response.Error(c, ecode.ErrListUserDevices)
				return
			}
			userDevicess = append(userDevicess, record)
		}
	}

	response.Success(c, gin.H{
		"userDevicess": userDevicess,
	})
}

// ListByLastID get records by last id and limit
// @Summary list of userDevicess by last id and limit
// @Description list of userDevicess by last id and limit
// @Tags userDevices
// @accept json
// @Produce json
// @Param lastID query int true "last id, default is MaxInt32" default(0)
// @Param limit query int false "number per page" default(10)
// @Param sort query string false "sort by column name of table, and the "-" sign before column name indicates reverse order" default(-id)
// @Success 200 {object} types.ListUserDevicessReply{}
// @Router /api/v1/userDevices/list [get]
// @Security BearerAuth
func (h *userDevicesHandler) ListByLastID(c *gin.Context) {
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
	userDevicess, err := h.iDao.GetByLastID(ctx, lastID, limit, sort)
	if err != nil {
		logger.Error("GetByLastID error", logger.Err(err), logger.Uint64("latsID", lastID), logger.Int("limit", limit), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertUserDevicess(userDevicess)
	if err != nil {
		response.Error(c, ecode.ErrListByLastIDUserDevices)
		return
	}

	response.Success(c, gin.H{
		"userDevicess": data,
	})
}

func getUserDevicesIDFromPath(c *gin.Context) (string, uint64, bool) {
	idStr := c.Param("id")
	id, err := utils.StrToUint64E(idStr)
	if err != nil || id == 0 {
		logger.Warn("StrToUint64E error: ", logger.String("idStr", idStr), middleware.GCtxRequestIDField(c))
		return "", 0, true
	}

	return idStr, id, false
}

func convertUserDevices(userDevices *model.UserDevices) (*types.UserDevicesObjDetail, error) {
	data := &types.UserDevicesObjDetail{}
	err := copier.Copy(data, userDevices)
	if err != nil {
		return nil, err
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	return data, nil
}

func convertUserDevicess(fromValues []*model.UserDevices) ([]*types.UserDevicesObjDetail, error) {
	toValues := []*types.UserDevicesObjDetail{}
	for _, v := range fromValues {
		data, err := convertUserDevices(v)
		if err != nil {
			return nil, err
		}
		toValues = append(toValues, data)
	}

	return toValues, nil
}
