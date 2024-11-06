package handler

import (
	"errors"
	"math"
	"time"

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

var _ InterestsHandler = (*interestsHandler)(nil)

// InterestsHandler defining the handler interface
type InterestsHandler interface {
	Create(c *gin.Context)
	DeleteByID(c *gin.Context)
	UpdateByID(c *gin.Context)
	GetByID(c *gin.Context)
	List(c *gin.Context)

	DeleteByIDs(c *gin.Context)
	GetByCondition(c *gin.Context)
	ListByIDs(c *gin.Context)
	ListByLastID(c *gin.Context)

	AllList(c *gin.Context)
}

type interestsHandler struct {
	iDao  dao.InterestsDao
	cache cache.InterestsCache
}

// NewInterestsHandler creating the handler interface
func NewInterestsHandler() InterestsHandler {
	return &interestsHandler{
		iDao: dao.NewInterestsDao(
			model.GetDB(),
			cache.NewInterestsCache(model.GetCacheType()),
		),
		cache: cache.NewInterestsCache(model.GetCacheType()),
	}
}

// AllList  obtain all interests tags
// @Summary create interestsTranslations
// @Description obtain all interests tags
// @Tags interests
// @accept json
// @Produce json
// @Param language_code path string true "language_code"
// @Success 200 {object} types.GetInterestsByLanguageReply{}
// @Router /api/v1/interests/allList/{language_code} [get]
// @Security BearerAuth
func (h *interestsHandler) AllList(c *gin.Context) {
	languageCode, isAbort := getInterestsLanguageFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	allInterestsCache, err2 := h.cache.GetFromLanguageCodeAllInterestsCache(c, languageCode)
	if err2 == nil && allInterestsCache != nil {
		cacheData, err := convertInterestsTranslationss(allInterestsCache)
		if err == nil && cacheData != nil {
			response.Success(c, gin.H{
				"interestss": cacheData,
			})
			return
		}
	}

	interests, err := h.iDao.GetByLanguage(c, languageCode)
	if err != nil {
		logger.Warn("GetByLanguage err: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.ErrListInterests)
		return
	}
	data, err := convertInterestsTranslationss(interests)
	if err != nil {
		response.Error(c, ecode.ErrListByLastIDInterests)
		return
	}

	err = h.cache.SetFromLanguageCodeAllInterestsCache(c, languageCode, interests, time.Hour*24*30)
	if err != nil {
		logger.Warn("GetByLanguage set cache err: ", logger.Err(err), middleware.GCtxRequestIDField(c))
	}

	response.Success(c, gin.H{
		"interestss": data,
	})

}

// Create a record
// @Summary create interests
// @Description submit information to create interests
// @Tags interests
// @accept json
// @Produce json
// @Param data body types.CreateInterestsRequest true "interests information"
// @Success 200 {object} types.CreateInterestsReply{}
// @Router /api/v1/interests [post]
// @Security BearerAuth
func (h *interestsHandler) Create(c *gin.Context) {
	form := &types.CreateInterestsRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	interests := &model.Interests{}
	err = copier.Copy(interests, form)
	if err != nil {
		response.Error(c, ecode.ErrCreateInterests)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.Create(ctx, interests)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{"id": interests.ID})
}

// DeleteByID delete a record by id
// @Summary delete interests
// @Description delete interests by id
// @Tags interests
// @accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} types.DeleteInterestsByIDReply{}
// @Router /api/v1/interests/{id} [delete]
// @Security BearerAuth
func (h *interestsHandler) DeleteByID(c *gin.Context) {
	_, id, isAbort := getInterestsIDFromPath(c)
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
// @Summary update interests
// @Description update interests information by id
// @Tags interests
// @accept json
// @Produce json
// @Param id path string true "id"
// @Param data body types.UpdateInterestsByIDRequest true "interests information"
// @Success 200 {object} types.UpdateInterestsByIDReply{}
// @Router /api/v1/interests/{id} [put]
// @Security BearerAuth
func (h *interestsHandler) UpdateByID(c *gin.Context) {
	_, id, isAbort := getInterestsIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	form := &types.UpdateInterestsByIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	form.TagID = int64(id)

	interests := &model.Interests{}
	err = copier.Copy(interests, form)
	if err != nil {
		response.Error(c, ecode.ErrUpdateByIDInterests)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.UpdateByID(ctx, interests)
	if err != nil {
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// GetByID get a record by id
// @Summary get interests detail
// @Description get interests detail by id
// @Tags interests
// @Param id path string true "id"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetInterestsByIDReply{}
// @Router /api/v1/interests/{id} [get]
// @Security BearerAuth
func (h *interestsHandler) GetByID(c *gin.Context) {
	_, id, isAbort := getInterestsIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	interests, err := h.iDao.GetByID(ctx, id)
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

	data := &types.InterestsObjDetail{}
	err = copier.Copy(data, interests)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDInterests)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"interests": data})
}

// List of records by query parameters
// @Summary list of interestss by query parameters
// @Description list of interestss by paging and conditions
// @Tags interests
// @accept json
// @Produce json
// @Param data body types.Params true "query parameters"
// @Success 200 {object} types.ListInterestssReply{}
// @Router /api/v1/interests/list [post]
// @Security BearerAuth
func (h *interestsHandler) List(c *gin.Context) {
	form := &types.ListInterestssRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	interestss, total, err := h.iDao.GetByColumns(ctx, &form.Params)
	if err != nil {
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertInterestss(interestss)
	if err != nil {
		response.Error(c, ecode.ErrListInterests)
		return
	}

	response.Success(c, gin.H{
		"interestss": data,
		"total":      total,
	})
}

// DeleteByIDs delete records by batch id
// @Summary delete interestss
// @Description delete interestss by batch id
// @Tags interests
// @Param data body types.DeleteInterestssByIDsRequest true "id array"
// @Accept json
// @Produce json
// @Success 200 {object} types.DeleteInterestssByIDsReply{}
// @Router /api/v1/interests/delete/ids [post]
// @Security BearerAuth
func (h *interestsHandler) DeleteByIDs(c *gin.Context) {
	form := &types.DeleteInterestssByIDsRequest{}
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
// @Summary get interests by condition
// @Description get interests by condition
// @Tags interests
// @Param data body types.Conditions true "query condition"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetInterestsByConditionReply{}
// @Router /api/v1/interests/condition [post]
// @Security BearerAuth
func (h *interestsHandler) GetByCondition(c *gin.Context) {
	form := &types.GetInterestsByConditionRequest{}
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
	interests, err := h.iDao.GetByCondition(ctx, &form.Conditions)
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

	data := &types.InterestsObjDetail{}
	err = copier.Copy(data, interests)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDInterests)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"interests": data})
}

// ListByIDs list of records by batch id
// @Summary list of interestss by batch id
// @Description list of interestss by batch id
// @Tags interests
// @Param data body types.ListInterestssByIDsRequest true "id array"
// @Accept json
// @Produce json
// @Success 200 {object} types.ListInterestssByIDsReply{}
// @Router /api/v1/interests/list/ids [post]
// @Security BearerAuth
func (h *interestsHandler) ListByIDs(c *gin.Context) {
	form := &types.ListInterestssByIDsRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	interestsMap, err := h.iDao.GetByIDs(ctx, form.IDs)
	if err != nil {
		logger.Error("GetByIDs error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	interestss := []*types.InterestsObjDetail{}
	for _, id := range form.IDs {
		if v, ok := interestsMap[id]; ok {
			record, err := convertInterests(v)
			if err != nil {
				response.Error(c, ecode.ErrListInterests)
				return
			}
			interestss = append(interestss, record)
		}
	}

	response.Success(c, gin.H{
		"interestss": interestss,
	})
}

// ListByLastID get records by last id and limit
// @Summary list of interestss by last id and limit
// @Description list of interestss by last id and limit
// @Tags interests
// @accept json
// @Produce json
// @Param lastID query int true "last id, default is MaxInt32" default(0)
// @Param limit query int false "number per page" default(10)
// @Param sort query string false "sort by column name of table, and the "-" sign before column name indicates reverse order" default(-id)
// @Success 200 {object} types.ListInterestssReply{}
// @Router /api/v1/interests/list [get]
// @Security BearerAuth
func (h *interestsHandler) ListByLastID(c *gin.Context) {
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
	interestss, err := h.iDao.GetByLastID(ctx, lastID, limit, sort)
	if err != nil {
		logger.Error("GetByLastID error", logger.Err(err), logger.Uint64("latsID", lastID), logger.Int("limit", limit), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertInterestss(interestss)
	if err != nil {
		response.Error(c, ecode.ErrListByLastIDInterests)
		return
	}

	response.Success(c, gin.H{
		"interestss": data,
	})
}

func getInterestsIDFromPath(c *gin.Context) (string, uint64, bool) {
	idStr := c.Param("id")
	id, err := utils.StrToUint64E(idStr)
	if err != nil || id == 0 {
		logger.Warn("StrToUint64E error: ", logger.String("idStr", idStr), middleware.GCtxRequestIDField(c))
		return "", 0, true
	}

	return idStr, id, false
}

func getInterestsLanguageFromPath(c *gin.Context) (string, bool) {
	languageCode := c.Param("language_code")
	if languageCode == "" {
		logger.Warn("languageCode error: ", logger.String("idStr", languageCode), middleware.GCtxRequestIDField(c))
		return "", true
	}

	return languageCode, false
}

func convertInterests(interests *model.Interests) (*types.InterestsObjDetail, error) {
	data := &types.InterestsObjDetail{}
	err := copier.Copy(data, interests)
	if err != nil {
		return nil, err
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	return data, nil
}

func convertInterestss(fromValues []*model.Interests) ([]*types.InterestsObjDetail, error) {
	toValues := []*types.InterestsObjDetail{}
	for _, v := range fromValues {
		data, err := convertInterests(v)
		if err != nil {
			return nil, err
		}
		toValues = append(toValues, data)
	}

	return toValues, nil
}

func convertInterestsTranslations(interestsTranslations *model.InterestsTranslations) (*types.InterestTranslationDetail, error) {
	data := &types.InterestTranslationDetail{}
	err := copier.Copy(data, interestsTranslations)
	if err != nil {
		return nil, err
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	return data, nil
}

func convertInterestsTranslationss(fromValues []*model.InterestsTranslations) ([]*types.InterestTranslationDetail, error) {
	var toValues []*types.InterestTranslationDetail
	for _, v := range fromValues {
		data, err := convertInterestsTranslations(v)
		if err != nil {
			return nil, err
		}
		toValues = append(toValues, data)
	}

	return toValues, nil
}
