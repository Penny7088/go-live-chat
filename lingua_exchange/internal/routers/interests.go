package handler

import (
	"github.com/gin-gonic/gin"
	"lingua_exchange/pkg/jwt"

	"lingua_exchange/internal/handler"
)

func init() {
	apiV1RouterFns = append(apiV1RouterFns, func(group *gin.RouterGroup) {
		interestsTranslationsRouter(group, handler.NewInterestsHandler())
	})
}

func interestsTranslationsRouter(group *gin.RouterGroup, h handler.InterestsHandler) {
	g := group.Group("/interests", jwt.AuthMiddleware())

	g.GET("/allList/:languageCode", h.AllList) // [get] /api/v1/interests/AllList:language_code
	// All the following routes use jwt authentication, you also can use middleware.Auth(middleware.WithVerify(fn))
	// g.Use(middleware.Auth())

	// If jwt authentication is not required for all routes, authentication middleware can be added
	// separately for only certain routes. In this case, g.Use(middleware.Auth()) above should not be used.

	// g.POST("/", h.Create) // [post] /api/v1/interestsTranslations
	// g.DELETE("/:id", h.DeleteByID) // [delete] /api/v1/interestsTranslations/:id
	// g.PUT("/:id", h.UpdateByID)    // [put] /api/v1/interestsTranslations/:id
	// g.GET("/:id", h.GetByID)       // [get] /api/v1/interestsTranslations/:id

	// g.POST("/delete/ids", h.DeleteByIDs)   // [post] /api/v1/interestsTranslations/delete/ids
	// g.POST("/condition", h.GetByCondition) // [post] /api/v1/interestsTranslations/condition
	// g.POST("/list/ids", h.ListByIDs)       // [post] /api/v1/interestsTranslations/list/ids
	// g.GET("/list", h.ListByLastID)         // [get] /api/v1/interestsTranslations/list
}
