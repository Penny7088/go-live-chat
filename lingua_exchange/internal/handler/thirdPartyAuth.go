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

var _ ThirdPartyAuthHandler = (*thirdPartyAuthHandler)(nil)

// ThirdPartyAuthHandler defining the handler interface
type ThirdPartyAuthHandler interface {
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

type thirdPartyAuthHandler struct {
	iDao dao.ThirdPartyAuthDao
}

// NewThirdPartyAuthHandler creating the handler interface
func NewThirdPartyAuthHandler() ThirdPartyAuthHandler {
	return &thirdPartyAuthHandler{
		iDao: dao.NewThirdPartyAuthDao(
			model.GetDB(),
			cache.NewThirdPartyAuthCache(model.GetCacheType()),
		),
	}
}

// Create a record
// @Summary create thirdPartyAuth
// @Description submit information to create thirdPartyAuth
// @Tags thirdPartyAuth
// @accept json
// @Produce json
// @Param data body types.CreateThirdPartyAuthRequest true "thirdPartyAuth information"
// @Success 200 {object} types.CreateThirdPartyAuthReply{}
// @Router /api/v1/thirdPartyAuth [post]
// @Security BearerAuth
func (h *thirdPartyAuthHandler) Create(c *gin.Context) {
	form := &types.CreateThirdPartyAuthRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	thirdPartyAuth := &model.ThirdPartyAuth{}
	err = copier.Copy(thirdPartyAuth, form)
	if err != nil {
		response.Error(c, ecode.ErrCreateThirdPartyAuth)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.Create(ctx, thirdPartyAuth)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{"id": thirdPartyAuth.ID})
}

// DeleteByID delete a record by id
// @Summary delete thirdPartyAuth
// @Description delete thirdPartyAuth by id
// @Tags thirdPartyAuth
// @accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} types.DeleteThirdPartyAuthByIDReply{}
// @Router /api/v1/thirdPartyAuth/{id} [delete]
// @Security BearerAuth
func (h *thirdPartyAuthHandler) DeleteByID(c *gin.Context) {
	_, id, isAbort := getThirdPartyAuthIDFromPath(c)
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
// @Summary update thirdPartyAuth
// @Description update thirdPartyAuth information by id
// @Tags thirdPartyAuth
// @accept json
// @Produce json
// @Param id path string true "id"
// @Param data body types.UpdateThirdPartyAuthByIDRequest true "thirdPartyAuth information"
// @Success 200 {object} types.UpdateThirdPartyAuthByIDReply{}
// @Router /api/v1/thirdPartyAuth/{id} [put]
// @Security BearerAuth
func (h *thirdPartyAuthHandler) UpdateByID(c *gin.Context) {
	_, id, isAbort := getThirdPartyAuthIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	form := &types.UpdateThirdPartyAuthByIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	form.ID = id

	thirdPartyAuth := &model.ThirdPartyAuth{}
	err = copier.Copy(thirdPartyAuth, form)
	if err != nil {
		response.Error(c, ecode.ErrUpdateByIDThirdPartyAuth)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.UpdateByID(ctx, thirdPartyAuth)
	if err != nil {
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// GetByID get a record by id
// @Summary get thirdPartyAuth detail
// @Description get thirdPartyAuth detail by id
// @Tags thirdPartyAuth
// @Param id path string true "id"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetThirdPartyAuthByIDReply{}
// @Router /api/v1/thirdPartyAuth/{id} [get]
// @Security BearerAuth
func (h *thirdPartyAuthHandler) GetByID(c *gin.Context) {
	_, id, isAbort := getThirdPartyAuthIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	thirdPartyAuth, err := h.iDao.GetByID(ctx, id)
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

	data := &types.ThirdPartyAuthObjDetail{}
	err = copier.Copy(data, thirdPartyAuth)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDThirdPartyAuth)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"thirdPartyAuth": data})
}

// List of records by query parameters
// @Summary list of thirdPartyAuths by query parameters
// @Description list of thirdPartyAuths by paging and conditions
// @Tags thirdPartyAuth
// @accept json
// @Produce json
// @Param data body types.Params true "query parameters"
// @Success 200 {object} types.ListThirdPartyAuthsReply{}
// @Router /api/v1/thirdPartyAuth/list [post]
// @Security BearerAuth
func (h *thirdPartyAuthHandler) List(c *gin.Context) {
	form := &types.ListThirdPartyAuthsRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	thirdPartyAuths, total, err := h.iDao.GetByColumns(ctx, &form.Params)
	if err != nil {
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertThirdPartyAuths(thirdPartyAuths)
	if err != nil {
		response.Error(c, ecode.ErrListThirdPartyAuth)
		return
	}

	response.Success(c, gin.H{
		"thirdPartyAuths": data,
		"total":           total,
	})
}

// DeleteByIDs delete records by batch id
// @Summary delete thirdPartyAuths
// @Description delete thirdPartyAuths by batch id
// @Tags thirdPartyAuth
// @Param data body types.DeleteThirdPartyAuthsByIDsRequest true "id array"
// @Accept json
// @Produce json
// @Success 200 {object} types.DeleteThirdPartyAuthsByIDsReply{}
// @Router /api/v1/thirdPartyAuth/delete/ids [post]
// @Security BearerAuth
func (h *thirdPartyAuthHandler) DeleteByIDs(c *gin.Context) {
	form := &types.DeleteThirdPartyAuthsByIDsRequest{}
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
// @Summary get thirdPartyAuth by condition
// @Description get thirdPartyAuth by condition
// @Tags thirdPartyAuth
// @Param data body types.Conditions true "query condition"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetThirdPartyAuthByConditionReply{}
// @Router /api/v1/thirdPartyAuth/condition [post]
// @Security BearerAuth
func (h *thirdPartyAuthHandler) GetByCondition(c *gin.Context) {
	form := &types.GetThirdPartyAuthByConditionRequest{}
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
	thirdPartyAuth, err := h.iDao.GetByCondition(ctx, &form.Conditions)
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

	data := &types.ThirdPartyAuthObjDetail{}
	err = copier.Copy(data, thirdPartyAuth)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDThirdPartyAuth)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"thirdPartyAuth": data})
}

// ListByIDs list of records by batch id
// @Summary list of thirdPartyAuths by batch id
// @Description list of thirdPartyAuths by batch id
// @Tags thirdPartyAuth
// @Param data body types.ListThirdPartyAuthsByIDsRequest true "id array"
// @Accept json
// @Produce json
// @Success 200 {object} types.ListThirdPartyAuthsByIDsReply{}
// @Router /api/v1/thirdPartyAuth/list/ids [post]
// @Security BearerAuth
func (h *thirdPartyAuthHandler) ListByIDs(c *gin.Context) {
	form := &types.ListThirdPartyAuthsByIDsRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	thirdPartyAuthMap, err := h.iDao.GetByIDs(ctx, form.IDs)
	if err != nil {
		logger.Error("GetByIDs error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	thirdPartyAuths := []*types.ThirdPartyAuthObjDetail{}
	for _, id := range form.IDs {
		if v, ok := thirdPartyAuthMap[id]; ok {
			record, err := convertThirdPartyAuth(v)
			if err != nil {
				response.Error(c, ecode.ErrListThirdPartyAuth)
				return
			}
			thirdPartyAuths = append(thirdPartyAuths, record)
		}
	}

	response.Success(c, gin.H{
		"thirdPartyAuths": thirdPartyAuths,
	})
}

// ListByLastID get records by last id and limit
// @Summary list of thirdPartyAuths by last id and limit
// @Description list of thirdPartyAuths by last id and limit
// @Tags thirdPartyAuth
// @accept json
// @Produce json
// @Param lastID query int true "last id, default is MaxInt32" default(0)
// @Param limit query int false "number per page" default(10)
// @Param sort query string false "sort by column name of table, and the "-" sign before column name indicates reverse order" default(-id)
// @Success 200 {object} types.ListThirdPartyAuthsReply{}
// @Router /api/v1/thirdPartyAuth/list [get]
// @Security BearerAuth
func (h *thirdPartyAuthHandler) ListByLastID(c *gin.Context) {
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
	thirdPartyAuths, err := h.iDao.GetByLastID(ctx, lastID, limit, sort)
	if err != nil {
		logger.Error("GetByLastID error", logger.Err(err), logger.Uint64("latsID", lastID), logger.Int("limit", limit), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertThirdPartyAuths(thirdPartyAuths)
	if err != nil {
		response.Error(c, ecode.ErrListByLastIDThirdPartyAuth)
		return
	}

	response.Success(c, gin.H{
		"thirdPartyAuths": data,
	})
}

func getThirdPartyAuthIDFromPath(c *gin.Context) (string, uint64, bool) {
	idStr := c.Param("id")
	id, err := utils.StrToUint64E(idStr)
	if err != nil || id == 0 {
		logger.Warn("StrToUint64E error: ", logger.String("idStr", idStr), middleware.GCtxRequestIDField(c))
		return "", 0, true
	}

	return idStr, id, false
}

func convertThirdPartyAuth(thirdPartyAuth *model.ThirdPartyAuth) (*types.ThirdPartyAuthObjDetail, error) {
	data := &types.ThirdPartyAuthObjDetail{}
	err := copier.Copy(data, thirdPartyAuth)
	if err != nil {
		return nil, err
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	return data, nil
}

func convertThirdPartyAuths(fromValues []*model.ThirdPartyAuth) ([]*types.ThirdPartyAuthObjDetail, error) {
	toValues := []*types.ThirdPartyAuthObjDetail{}
	for _, v := range fromValues {
		data, err := convertThirdPartyAuth(v)
		if err != nil {
			return nil, err
		}
		toValues = append(toValues, data)
	}

	return toValues, nil
}
