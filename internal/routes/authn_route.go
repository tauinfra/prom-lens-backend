package routes

import (
	"prom-lens-backend/internal/pkg/di"

	"github.com/gin-gonic/gin"
)

func AuthnRouters(rg *gin.RouterGroup, provider *di.Provider) {
	authn := rg.Group("/authn")
	{
		authn.POST("/login", provider.Authn.Auth.Controller.Login)
		authn.POST("/refresh-token", provider.Authn.Auth.Controller.RefreshToken)
		authn.POST("/me/change-password", provider.Authn.Auth.Controller.ChangePassword)
		authn.GET("/me/profile", provider.Authn.Auth.Controller.GetProfile)

		users := authn.Group("/users")
		users.Use(provider.JWTManager.RequireSuperuser())
		{
			users.GET("", provider.Authn.User.Controller.List)
			users.GET("/:userID", provider.Authn.User.Controller.Get)
			users.POST("", provider.Authn.User.Controller.Create)
			users.PATCH("/:userID", provider.Authn.User.Controller.Update)
			users.DELETE("/:userID", provider.Authn.User.Controller.Delete)
			users.POST("/:userID/reset-password", provider.Authn.User.Controller.ResetPassword)
		}
	}
}
