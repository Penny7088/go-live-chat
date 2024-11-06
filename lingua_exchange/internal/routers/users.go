package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/zhufuyi/sponge/pkg/gin/middleware"
	"lingua_exchange/internal/handler"
)

func init() {
	apiV1RouterFns = append(apiV1RouterFns, func(group *gin.RouterGroup) {
		usersRouter(group, handler.NewUsersHandler())
	})
}

func usersRouter(group *gin.RouterGroup, h handler.UsersHandler) {
	g := group.Group("/users")

	// All the following routes use jwt authentication, you also can use middleware.Auth(middleware.WithVerify(fn))
	//g.Use(middleware.Auth())

	// If jwt authentication is not required for all routes, authentication middleware can be added
	// separately for only certain routes. In this case, g.Use(middleware.Auth()) above should not be used.
	g.POST("/auth", h.LoginOrRegister) // [post] /api/v1/auth

	g.POST("/", h.Create, middleware.Auth())          // [post] /api/v1/users
	g.DELETE("/:id", h.DeleteByID, middleware.Auth()) // [delete] /api/v1/users/:id
	g.PUT("/:id", h.UpdateByID, middleware.Auth())    // [put] /api/v1/users/:id
	g.GET("/:id", h.GetByID, middleware.Auth())       // [get] /api/v1/users/:id
	g.POST("/list", h.List, middleware.Auth())        // [post] /api/v1/users/list

	g.POST("/delete/ids", h.DeleteByIDs, middleware.Auth())   // [post] /api/v1/users/delete/ids
	g.POST("/condition", h.GetByCondition, middleware.Auth()) // [post] /api/v1/users/condition
	g.POST("/list/ids", h.ListByIDs, middleware.Auth())       // [post] /api/v1/users/list/ids
	g.GET("/list", h.ListByLastID, middleware.Auth())         // [get] /api/v1/users/list
}
