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

var _ UserInterestsHandler = (*userInterestsHandler)(nil)

// UserInterestsHandler defining the handler interface
type UserInterestsHandler interface {
	Create(c *gin.Context)
	DeleteByID(c *gin.Context)
	GetByID(c *gin.Context)
	List(c *gin.Context)

	DeleteByIDs(c *gin.Context)
	GetByCondition(c *gin.Context)
	ListByIDs(c *gin.Context)
	ListByLastID(c *gin.Context)
}

type userInterestsHandler struct {
	iDao dao.UserInterestsDao
}

// NewUserInterestsHandler creating the handler interface
func NewUserInterestsHandler() UserInterestsHandler {
	return &userInterestsHandler{
		iDao: dao.NewUserInterestsDao(
			model.GetDB(),
			cache.NewUserInterestsCache(model.GetCacheType()),
		),
	}
}

// Create a record
// @Summary create userInterests
// @Description submit information to create userInterests
// @Tags userInterests
// @accept json
// @Produce json
// @Param data body types.CreateUserInterestsRequest true "userInterests information"
// @Success 200 {object} types.CreateUserInterestsReply{}
// @Router /api/v1/userInterests [post]
// @Security BearerAuth
func (h *userInterestsHandler) Create(c *gin.Context) {
	form := &types.CreateUserInterestsRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	userInterests := &model.UserInterests{}
	err = copier.Copy(userInterests, form)
	if err != nil {
		response.Error(c, ecode.ErrCreateUserInterests)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.Create(ctx, userInterests)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{"id": userInterests.ID})
}

// DeleteByID delete a record by id
// @Summary delete userInterests
// @Description delete userInterests by id
// @Tags userInterests
// @accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} types.DeleteUserInterestsByIDReply{}
// @Router /api/v1/userInterests/{id} [delete]
// @Security BearerAuth
func (h *userInterestsHandler) DeleteByID(c *gin.Context) {
	_, id, isAbort := getUserInterestsIDFromPath(c)
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

// GetByID get a record by id
// @Summary get userInterests detail
// @Description get userInterests detail by id
// @Tags userInterests
// @Param id path string true "id"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetUserInterestsByIDReply{}
// @Router /api/v1/userInterests/{id} [get]
// @Security BearerAuth
func (h *userInterestsHandler) GetByID(c *gin.Context) {
	_, id, isAbort := getUserInterestsIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	userInterests, err := h.iDao.GetByID(ctx, id)
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

	data := &types.UserInterestsObjDetail{}
	err = copier.Copy(data, userInterests)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDUserInterests)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"userInterests": data})
}

// List of records by query parameters
// @Summary list of userInterestss by query parameters
// @Description list of userInterestss by paging and conditions
// @Tags userInterests
// @accept json
// @Produce json
// @Param data body types.Params true "query parameters"
// @Success 200 {object} types.ListUserInterestssReply{}
// @Router /api/v1/userInterests/list [post]
// @Security BearerAuth
func (h *userInterestsHandler) List(c *gin.Context) {
	form := &types.ListUserInterestssRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	userInterestss, total, err := h.iDao.GetByColumns(ctx, &form.Params)
	if err != nil {
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertUserInterestss(userInterestss)
	if err != nil {
		response.Error(c, ecode.ErrListUserInterests)
		return
	}

	response.Success(c, gin.H{
		"userInterestss": data,
		"total":          total,
	})
}

// DeleteByIDs delete records by batch id
// @Summary delete userInterestss
// @Description delete userInterestss by batch id
// @Tags userInterests
// @Param data body types.DeleteUserInterestssByIDsRequest true "id array"
// @Accept json
// @Produce json
// @Success 200 {object} types.DeleteUserInterestssByIDsReply{}
// @Router /api/v1/userInterests/delete/ids [post]
// @Security BearerAuth
func (h *userInterestsHandler) DeleteByIDs(c *gin.Context) {
	form := &types.DeleteUserInterestssByIDsRequest{}
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
// @Summary get userInterests by condition
// @Description get userInterests by condition
// @Tags userInterests
// @Param data body types.Conditions true "query condition"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetUserInterestsByConditionReply{}
// @Router /api/v1/userInterests/condition [post]
// @Security BearerAuth
func (h *userInterestsHandler) GetByCondition(c *gin.Context) {
	form := &types.GetUserInterestsByConditionRequest{}
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
	userInterests, err := h.iDao.GetByCondition(ctx, &form.Conditions)
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

	data := &types.UserInterestsObjDetail{}
	err = copier.Copy(data, userInterests)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDUserInterests)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"userInterests": data})
}

// ListByIDs list of records by batch id
// @Summary list of userInterestss by batch id
// @Description list of userInterestss by batch id
// @Tags userInterests
// @Param data body types.ListUserInterestssByIDsRequest true "id array"
// @Accept json
// @Produce json
// @Success 200 {object} types.ListUserInterestssByIDsReply{}
// @Router /api/v1/userInterests/list/ids [post]
// @Security BearerAuth
func (h *userInterestsHandler) ListByIDs(c *gin.Context) {
	form := &types.ListUserInterestssByIDsRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	userInterestsMap, err := h.iDao.GetByIDs(ctx, form.IDs)
	if err != nil {
		logger.Error("GetByIDs error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	userInterestss := []*types.UserInterestsObjDetail{}
	for _, id := range form.IDs {
		if v, ok := userInterestsMap[id]; ok {
			record, err := convertUserInterests(v)
			if err != nil {
				response.Error(c, ecode.ErrListUserInterests)
				return
			}
			userInterestss = append(userInterestss, record)
		}
	}

	response.Success(c, gin.H{
		"userInterestss": userInterestss,
	})
}

// ListByLastID get records by last id and limit
// @Summary list of userInterestss by last id and limit
// @Description list of userInterestss by last id and limit
// @Tags userInterests
// @accept json
// @Produce json
// @Param lastID query int true "last id, default is MaxInt32" default(0)
// @Param limit query int false "number per page" default(10)
// @Param sort query string false "sort by column name of table, and the "-" sign before column name indicates reverse order" default(-id)
// @Success 200 {object} types.ListUserInterestssReply{}
// @Router /api/v1/userInterests/list [get]
// @Security BearerAuth
func (h *userInterestsHandler) ListByLastID(c *gin.Context) {
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
	userInterestss, err := h.iDao.GetByLastID(ctx, lastID, limit, sort)
	if err != nil {
		logger.Error("GetByLastID error", logger.Err(err), logger.Uint64("latsID", lastID), logger.Int("limit", limit), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertUserInterestss(userInterestss)
	if err != nil {
		response.Error(c, ecode.ErrListByLastIDUserInterests)
		return
	}

	response.Success(c, gin.H{
		"userInterestss": data,
	})
}

func getUserInterestsIDFromPath(c *gin.Context) (string, uint64, bool) {
	idStr := c.Param("id")
	id, err := utils.StrToUint64E(idStr)
	if err != nil || id == 0 {
		logger.Warn("StrToUint64E error: ", logger.String("idStr", idStr), middleware.GCtxRequestIDField(c))
		return "", 0, true
	}

	return idStr, id, false
}

func convertUserInterests(userInterests *model.UserInterests) (*types.UserInterestsObjDetail, error) {
	data := &types.UserInterestsObjDetail{}
	err := copier.Copy(data, userInterests)
	if err != nil {
		return nil, err
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	return data, nil
}

func convertUserInterestss(fromValues []*model.UserInterests) ([]*types.UserInterestsObjDetail, error) {
	toValues := []*types.UserInterestsObjDetail{}
	for _, v := range fromValues {
		data, err := convertUserInterests(v)
		if err != nil {
			return nil, err
		}
		toValues = append(toValues, data)
	}

	return toValues, nil
}
