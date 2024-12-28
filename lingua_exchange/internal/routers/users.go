package handler

import (
	"github.com/gin-gonic/gin"
	"lingua_exchange/internal/handler"
	"lingua_exchange/pkg/jwt"
)

func init() {
	apiV1RouterFns = append(apiV1RouterFns, func(group *gin.RouterGroup) {
		usersRouter(group, handler.NewUsersHandler())
	})
}

func usersRouter(group *gin.RouterGroup, h handler.UsersHandler) {
	g := group.Group("/users")

	// All the following routes use jwt authentication, you also can use middleware.Auth(middleware.WithVerify(fn))
	// g.Use(jwt.AuthMiddleware())

	// If jwt authentication is not required for all routes, authentication middleware can be added
	// separately for only certain routes. In this case, g.Use(jwt.AuthMiddleware()) above should not be used.
	g.POST("/auth", h.LoginOrRegister)                                       // [post] /api/v1/auth
	g.POST("/loginFromEmail", h.LoginFromEmail)                              // [post] /api/v1/LoginFromEmail
	g.POST("/signUpFromEmail", h.SignUpFromEmail)                            // [post] /api/v1/signUpFromEmail
	g.POST("/resetPassword", h.ResetPassword)                                // [post] /api/v1/resetPassword
	g.PUT("/updateUserInfo/:id", h.UpdateUserInfoByID, jwt.AuthMiddleware()) // [put] /api/v1/users/:id

	g.POST("/", h.Create, jwt.AuthMiddleware())          // [post] /api/v1/users
	g.DELETE("/:id", h.DeleteByID, jwt.AuthMiddleware()) // [delete] /api/v1/users/:id
	g.GET("/:id", h.GetByID, jwt.AuthMiddleware())       // [get] /api/v1/users/:id
	g.POST("/list", h.List, jwt.AuthMiddleware())        // [post] /api/v1/users/list

	g.POST("/delete/ids", h.DeleteByIDs, jwt.AuthMiddleware())   // [post] /api/v1/users/delete/ids
	g.POST("/condition", h.GetByCondition, jwt.AuthMiddleware()) // [post] /api/v1/users/condition
	g.POST("/list/ids", h.ListByIDs, jwt.AuthMiddleware())       // [post] /api/v1/users/list/ids
	g.GET("/list", h.ListByLastID, jwt.AuthMiddleware())         // [get] /api/v1/users/list
}
