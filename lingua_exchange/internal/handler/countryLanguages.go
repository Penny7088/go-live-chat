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

var _ CountryLanguagesHandler = (*countryLanguagesHandler)(nil)

// CountryLanguagesHandler defining the handler interface
type CountryLanguagesHandler interface {
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

type countryLanguagesHandler struct {
	iDao dao.CountryLanguagesDao
}

// NewCountryLanguagesHandler creating the handler interface
func NewCountryLanguagesHandler() CountryLanguagesHandler {
	return &countryLanguagesHandler{
		iDao: dao.NewCountryLanguagesDao(
			model.GetDB(),
			cache.NewCountryLanguagesCache(model.GetCacheType()),
		),
	}
}

// Create a record
// @Summary create countryLanguages
// @Description submit information to create countryLanguages
// @Tags countryLanguages
// @accept json
// @Produce json
// @Param data body types.CreateCountryLanguagesRequest true "countryLanguages information"
// @Success 200 {object} types.CreateCountryLanguagesReply{}
// @Router /api/v1/countryLanguages [post]
// @Security BearerAuth
func (h *countryLanguagesHandler) Create(c *gin.Context) {
	form := &types.CreateCountryLanguagesRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	countryLanguages := &model.CountryLanguages{}
	err = copier.Copy(countryLanguages, form)
	if err != nil {
		response.Error(c, ecode.ErrCreateCountryLanguages)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.Create(ctx, countryLanguages)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{"id": countryLanguages.ID})
}

// DeleteByID delete a record by id
// @Summary delete countryLanguages
// @Description delete countryLanguages by id
// @Tags countryLanguages
// @accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} types.DeleteCountryLanguagesByIDReply{}
// @Router /api/v1/countryLanguages/{id} [delete]
// @Security BearerAuth
func (h *countryLanguagesHandler) DeleteByID(c *gin.Context) {
	_, id, isAbort := getCountryLanguagesIDFromPath(c)
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
// @Summary update countryLanguages
// @Description update countryLanguages information by id
// @Tags countryLanguages
// @accept json
// @Produce json
// @Param id path string true "id"
// @Param data body types.UpdateCountryLanguagesByIDRequest true "countryLanguages information"
// @Success 200 {object} types.UpdateCountryLanguagesByIDReply{}
// @Router /api/v1/countryLanguages/{id} [put]
// @Security BearerAuth
func (h *countryLanguagesHandler) UpdateByID(c *gin.Context) {
	_, id, isAbort := getCountryLanguagesIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	form := &types.UpdateCountryLanguagesByIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	form.ID = id

	countryLanguages := &model.CountryLanguages{}
	err = copier.Copy(countryLanguages, form)
	if err != nil {
		response.Error(c, ecode.ErrUpdateByIDCountryLanguages)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.UpdateByID(ctx, countryLanguages)
	if err != nil {
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// GetByID get a record by id
// @Summary get countryLanguages detail
// @Description get countryLanguages detail by id
// @Tags countryLanguages
// @Param id path string true "id"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetCountryLanguagesByIDReply{}
// @Router /api/v1/countryLanguages/{id} [get]
// @Security BearerAuth
func (h *countryLanguagesHandler) GetByID(c *gin.Context) {
	_, id, isAbort := getCountryLanguagesIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	countryLanguages, err := h.iDao.GetByID(ctx, id)
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

	data := &types.CountryLanguagesObjDetail{}
	err = copier.Copy(data, countryLanguages)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDCountryLanguages)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"countryLanguages": data})
}

// List of records by query parameters
// @Summary list of countryLanguagess by query parameters
// @Description list of countryLanguagess by paging and conditions
// @Tags countryLanguages
// @accept json
// @Produce json
// @Param data body types.Params true "query parameters"
// @Success 200 {object} types.ListCountryLanguagessReply{}
// @Router /api/v1/countryLanguages/list [post]
// @Security BearerAuth
func (h *countryLanguagesHandler) List(c *gin.Context) {
	form := &types.ListCountryLanguagessRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	countryLanguagess, total, err := h.iDao.GetByColumns(ctx, &form.Params)
	if err != nil {
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertCountryLanguagess(countryLanguagess)
	if err != nil {
		response.Error(c, ecode.ErrListCountryLanguages)
		return
	}

	response.Success(c, gin.H{
		"countryLanguagess": data,
		"total":             total,
	})
}

// DeleteByIDs delete records by batch id
// @Summary delete countryLanguagess
// @Description delete countryLanguagess by batch id
// @Tags countryLanguages
// @Param data body types.DeleteCountryLanguagessByIDsRequest true "id array"
// @Accept json
// @Produce json
// @Success 200 {object} types.DeleteCountryLanguagessByIDsReply{}
// @Router /api/v1/countryLanguages/delete/ids [post]
// @Security BearerAuth
func (h *countryLanguagesHandler) DeleteByIDs(c *gin.Context) {
	form := &types.DeleteCountryLanguagessByIDsRequest{}
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
// @Summary get countryLanguages by condition
// @Description get countryLanguages by condition
// @Tags countryLanguages
// @Param data body types.Conditions true "query condition"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetCountryLanguagesByConditionReply{}
// @Router /api/v1/countryLanguages/condition [post]
// @Security BearerAuth
func (h *countryLanguagesHandler) GetByCondition(c *gin.Context) {
	form := &types.GetCountryLanguagesByConditionRequest{}
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
	countryLanguages, err := h.iDao.GetByCondition(ctx, &form.Conditions)
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

	data := &types.CountryLanguagesObjDetail{}
	err = copier.Copy(data, countryLanguages)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDCountryLanguages)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"countryLanguages": data})
}

// ListByIDs list of records by batch id
// @Summary list of countryLanguagess by batch id
// @Description list of countryLanguagess by batch id
// @Tags countryLanguages
// @Param data body types.ListCountryLanguagessByIDsRequest true "id array"
// @Accept json
// @Produce json
// @Success 200 {object} types.ListCountryLanguagessByIDsReply{}
// @Router /api/v1/countryLanguages/list/ids [post]
// @Security BearerAuth
func (h *countryLanguagesHandler) ListByIDs(c *gin.Context) {
	form := &types.ListCountryLanguagessByIDsRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	countryLanguagesMap, err := h.iDao.GetByIDs(ctx, form.IDs)
	if err != nil {
		logger.Error("GetByIDs error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	countryLanguagess := []*types.CountryLanguagesObjDetail{}
	for _, id := range form.IDs {
		if v, ok := countryLanguagesMap[id]; ok {
			record, err := convertCountryLanguages(v)
			if err != nil {
				response.Error(c, ecode.ErrListCountryLanguages)
				return
			}
			countryLanguagess = append(countryLanguagess, record)
		}
	}

	response.Success(c, gin.H{
		"countryLanguagess": countryLanguagess,
	})
}

// ListByLastID get records by last id and limit
// @Summary list of countryLanguagess by last id and limit
// @Description list of countryLanguagess by last id and limit
// @Tags countryLanguages
// @accept json
// @Produce json
// @Param lastID query int true "last id, default is MaxInt32" default(0)
// @Param limit query int false "number per page" default(10)
// @Param sort query string false "sort by column name of table, and the "-" sign before column name indicates reverse order" default(-id)
// @Success 200 {object} types.ListCountryLanguagessReply{}
// @Router /api/v1/countryLanguages/list [get]
// @Security BearerAuth
func (h *countryLanguagesHandler) ListByLastID(c *gin.Context) {
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
	countryLanguagess, err := h.iDao.GetByLastID(ctx, lastID, limit, sort)
	if err != nil {
		logger.Error("GetByLastID error", logger.Err(err), logger.Uint64("latsID", lastID), logger.Int("limit", limit), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertCountryLanguagess(countryLanguagess)
	if err != nil {
		response.Error(c, ecode.ErrListByLastIDCountryLanguages)
		return
	}

	response.Success(c, gin.H{
		"countryLanguagess": data,
	})
}

func getCountryLanguagesIDFromPath(c *gin.Context) (string, uint64, bool) {
	idStr := c.Param("id")
	id, err := utils.StrToUint64E(idStr)
	if err != nil || id == 0 {
		logger.Warn("StrToUint64E error: ", logger.String("idStr", idStr), middleware.GCtxRequestIDField(c))
		return "", 0, true
	}

	return idStr, id, false
}

func convertCountryLanguages(countryLanguages *model.CountryLanguages) (*types.CountryLanguagesObjDetail, error) {
	data := &types.CountryLanguagesObjDetail{}
	err := copier.Copy(data, countryLanguages)
	if err != nil {
		return nil, err
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	return data, nil
}

func convertCountryLanguagess(fromValues []*model.CountryLanguages) ([]*types.CountryLanguagesObjDetail, error) {
	toValues := []*types.CountryLanguagesObjDetail{}
	for _, v := range fromValues {
		data, err := convertCountryLanguages(v)
		if err != nil {
			return nil, err
		}
		toValues = append(toValues, data)
	}

	return toValues, nil
}
