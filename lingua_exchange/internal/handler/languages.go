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

var _ LanguagesHandler = (*languagesHandler)(nil)

// LanguagesHandler defining the handler interface
type LanguagesHandler interface {
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

type languagesHandler struct {
	iDao dao.LanguagesDao
}

// NewLanguagesHandler creating the handler interface
func NewLanguagesHandler() LanguagesHandler {
	return &languagesHandler{
		iDao: dao.NewLanguagesDao(
			model.GetDB(),
			cache.NewLanguagesCache(model.GetCacheType()),
		),
	}
}

// Create a record
// @Summary create languages
// @Description submit information to create languages
// @Tags languages
// @accept json
// @Produce json
// @Param data body types.CreateLanguagesRequest true "languages information"
// @Success 200 {object} types.CreateLanguagesReply{}
// @Router /api/v1/languages [post]
// @Security BearerAuth
func (h *languagesHandler) Create(c *gin.Context) {
	form := &types.CreateLanguagesRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	languages := &model.Languages{}
	err = copier.Copy(languages, form)
	if err != nil {
		response.Error(c, ecode.ErrCreateLanguages)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.Create(ctx, languages)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{"id": languages.ID})
}

// DeleteByID delete a record by id
// @Summary delete languages
// @Description delete languages by id
// @Tags languages
// @accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} types.DeleteLanguagesByIDReply{}
// @Router /api/v1/languages/{id} [delete]
// @Security BearerAuth
func (h *languagesHandler) DeleteByID(c *gin.Context) {
	_, id, isAbort := getLanguagesIDFromPath(c)
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
// @Summary update languages
// @Description update languages information by id
// @Tags languages
// @accept json
// @Produce json
// @Param id path string true "id"
// @Param data body types.UpdateLanguagesByIDRequest true "languages information"
// @Success 200 {object} types.UpdateLanguagesByIDReply{}
// @Router /api/v1/languages/{id} [put]
// @Security BearerAuth
func (h *languagesHandler) UpdateByID(c *gin.Context) {
	_, id, isAbort := getLanguagesIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	form := &types.UpdateLanguagesByIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	form.ID = id

	languages := &model.Languages{}
	err = copier.Copy(languages, form)
	if err != nil {
		response.Error(c, ecode.ErrUpdateByIDLanguages)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	ctx := middleware.WrapCtx(c)
	err = h.iDao.UpdateByID(ctx, languages)
	if err != nil {
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// GetByID get a record by id
// @Summary get languages detail
// @Description get languages detail by id
// @Tags languages
// @Param id path string true "id"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetLanguagesByIDReply{}
// @Router /api/v1/languages/{id} [get]
// @Security BearerAuth
func (h *languagesHandler) GetByID(c *gin.Context) {
	_, id, isAbort := getLanguagesIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	languages, err := h.iDao.GetByID(ctx, id)
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

	data := &types.LanguagesObjDetail{}
	err = copier.Copy(data, languages)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDLanguages)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"languages": data})
}

// List of records by query parameters
// @Summary list of languagess by query parameters
// @Description list of languagess by paging and conditions
// @Tags languages
// @accept json
// @Produce json
// @Param data body types.Params true "query parameters"
// @Success 200 {object} types.ListLanguagessReply{}
// @Router /api/v1/languages/list [post]
// @Security BearerAuth
func (h *languagesHandler) List(c *gin.Context) {
	form := &types.ListLanguagessRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	languagess, total, err := h.iDao.GetByColumns(ctx, &form.Params)
	if err != nil {
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertLanguagess(languagess)
	if err != nil {
		response.Error(c, ecode.ErrListLanguages)
		return
	}

	response.Success(c, gin.H{
		"languagess": data,
		"total":      total,
	})
}

// DeleteByIDs delete records by batch id
// @Summary delete languagess
// @Description delete languagess by batch id
// @Tags languages
// @Param data body types.DeleteLanguagessByIDsRequest true "id array"
// @Accept json
// @Produce json
// @Success 200 {object} types.DeleteLanguagessByIDsReply{}
// @Router /api/v1/languages/delete/ids [post]
// @Security BearerAuth
func (h *languagesHandler) DeleteByIDs(c *gin.Context) {
	form := &types.DeleteLanguagessByIDsRequest{}
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
// @Summary get languages by condition
// @Description get languages by condition
// @Tags languages
// @Param data body types.Conditions true "query condition"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetLanguagesByConditionReply{}
// @Router /api/v1/languages/condition [post]
// @Security BearerAuth
func (h *languagesHandler) GetByCondition(c *gin.Context) {
	form := &types.GetLanguagesByConditionRequest{}
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
	languages, err := h.iDao.GetByCondition(ctx, &form.Conditions)
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

	data := &types.LanguagesObjDetail{}
	err = copier.Copy(data, languages)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDLanguages)
		return
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	response.Success(c, gin.H{"languages": data})
}

// ListByIDs list of records by batch id
// @Summary list of languagess by batch id
// @Description list of languagess by batch id
// @Tags languages
// @Param data body types.ListLanguagessByIDsRequest true "id array"
// @Accept json
// @Produce json
// @Success 200 {object} types.ListLanguagessByIDsReply{}
// @Router /api/v1/languages/list/ids [post]
// @Security BearerAuth
func (h *languagesHandler) ListByIDs(c *gin.Context) {
	form := &types.ListLanguagessByIDsRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	languagesMap, err := h.iDao.GetByIDs(ctx, form.IDs)
	if err != nil {
		logger.Error("GetByIDs error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	languagess := []*types.LanguagesObjDetail{}
	for _, id := range form.IDs {
		if v, ok := languagesMap[id]; ok {
			record, err := convertLanguages(v)
			if err != nil {
				response.Error(c, ecode.ErrListLanguages)
				return
			}
			languagess = append(languagess, record)
		}
	}

	response.Success(c, gin.H{
		"languagess": languagess,
	})
}

// ListByLastID get records by last id and limit
// @Summary list of languagess by last id and limit
// @Description list of languagess by last id and limit
// @Tags languages
// @accept json
// @Produce json
// @Param lastID query int true "last id, default is MaxInt32" default(0)
// @Param limit query int false "number per page" default(10)
// @Param sort query string false "sort by column name of table, and the "-" sign before column name indicates reverse order" default(-id)
// @Success 200 {object} types.ListLanguagessReply{}
// @Router /api/v1/languages/list [get]
// @Security BearerAuth
func (h *languagesHandler) ListByLastID(c *gin.Context) {
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
	languagess, err := h.iDao.GetByLastID(ctx, lastID, limit, sort)
	if err != nil {
		logger.Error("GetByLastID error", logger.Err(err), logger.Uint64("latsID", lastID), logger.Int("limit", limit), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertLanguagess(languagess)
	if err != nil {
		response.Error(c, ecode.ErrListByLastIDLanguages)
		return
	}

	response.Success(c, gin.H{
		"languagess": data,
	})
}

func getLanguagesIDFromPath(c *gin.Context) (string, uint64, bool) {
	idStr := c.Param("id")
	id, err := utils.StrToUint64E(idStr)
	if err != nil || id == 0 {
		logger.Warn("StrToUint64E error: ", logger.String("idStr", idStr), middleware.GCtxRequestIDField(c))
		return "", 0, true
	}

	return idStr, id, false
}

func convertLanguages(languages *model.Languages) (*types.LanguagesObjDetail, error) {
	data := &types.LanguagesObjDetail{}
	err := copier.Copy(data, languages)
	if err != nil {
		return nil, err
	}
	// Note: if copier.Copy cannot assign a value to a field, add it here

	return data, nil
}

func convertLanguagess(fromValues []*model.Languages) ([]*types.LanguagesObjDetail, error) {
	toValues := []*types.LanguagesObjDetail{}
	for _, v := range fromValues {
		data, err := convertLanguages(v)
		if err != nil {
			return nil, err
		}
		toValues = append(toValues, data)
	}

	return toValues, nil
}
