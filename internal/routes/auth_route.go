package routes

import (
	"valyria-backend/internal/pkg/di"

	"github.com/gin-gonic/gin"
)

func AuthRouters(rg *gin.RouterGroup, provider *di.Provider) {
	users := rg.Group("/auth")
	{
		// User Route
		users.GET("users", provider.Auth.User.Controller.List)
		users.GET("users/:id", provider.Auth.User.Controller.Get)
		users.POST("users", provider.Auth.User.Controller.Create)
		users.PATCH("users/:id", provider.Auth.User.Controller.Update)
		users.DELETE("users/:id", provider.Auth.User.Controller.Delete)
		users.POST("/login", provider.Auth.User.Controller.Login)
		// Role Route
		users.GET("roles", provider.Auth.Role.Controller.List)
		users.GET("roles/:id", provider.Auth.Role.Controller.Get)
		users.POST("roles", provider.Auth.Role.Controller.Create)
		users.PATCH("roles/:id", provider.Auth.Role.Controller.Update)
		users.DELETE("roles/:id", provider.Auth.Role.Controller.Delete)
		// Module Route
		users.GET("modules", provider.Auth.Module.Controller.List)
		users.GET("modules/:id", provider.Auth.Module.Controller.Get)
		users.POST("modules", provider.Auth.Module.Controller.Create)
		users.PATCH("modules/:id", provider.Auth.Module.Controller.Update)
		users.DELETE("modules/:id", provider.Auth.Module.Controller.Delete)
	}
}
