package routes

import (
	"valyria-backend/internal/pkg/di"

	"github.com/gin-gonic/gin"
)

func AuthnRouters(rg *gin.RouterGroup, provider *di.Provider) {
	authn := rg.Group("/authn")
	{
		//
		authn.POST("/login", provider.Authn.Auth.Controller.Login)
		authn.POST("/refresh-token", provider.Authn.Auth.Controller.RefreshToken)
		authn.POST("/me/change-password", provider.Authn.Auth.Controller.ChangePassword) // 需要登录
		authn.GET("/me/profile", provider.Authn.Auth.Controller.GetProfile)
		authn.GET("/me/routes", provider.Authn.Auth.Controller.GetUserRoutes)
		// Users
		authn.GET("/users", provider.Authn.User.Controller.List)
		authn.POST("/users", provider.Authn.User.Controller.Create)
		authn.PATCH("/users/:userID", provider.Authn.User.Controller.Update)
		authn.DELETE("/users/:userID", provider.Authn.User.Controller.Delete)
		authn.POST("/users/:userID/reset-password", provider.Authn.User.Controller.ResetPassword)
		authn.GET("/users/:userID/roles", provider.Authn.User.Controller.GetUserRoles)
		authn.PATCH("/users/:userID/roles", provider.Authn.User.Controller.UpdateRoleRoles)
		authn.GET("/users/:userID/menus", provider.Authn.User.Controller.GetUserMenus)
		authn.PATCH("/users/:userID/menus", provider.Authn.User.Controller.UpdateUserMenus)
		// Roles Route
		authn.GET("/roles", provider.Authn.Role.Controller.List)
		authn.POST("/roles", provider.Authn.Role.Controller.Create)
		authn.GET("/roles/:roleID", provider.Authn.Role.Controller.Get)
		authn.PATCH("/roles/:roleID", provider.Authn.Role.Controller.Update)
		authn.DELETE("/roles/:roleID", provider.Authn.Role.Controller.Delete)
		// Menus Route
		authn.GET("/menus", provider.Authn.Menu.Controller.List)
		authn.POST("/menus", provider.Authn.Menu.Controller.Create)
		authn.GET("/menus/:menuID", provider.Authn.Menu.Controller.Get)
		authn.PATCH("/menus/:menuID", provider.Authn.Menu.Controller.Update)
		authn.DELETE("/menus/:menuID", provider.Authn.Menu.Controller.Delete)
		//
		authn.GET("/roles/:roleID/permissions", provider.Authn.Role.Controller.GetRolePermissions)
		authn.PATCH("/roles/:roleID/permissions", provider.Authn.Role.Controller.UpdateRolePermissions)
		// Permissions Route
		authn.GET("/permissions", provider.Authn.Permission.Controller.List)
		authn.POST("/permissions", provider.Authn.Permission.Controller.Create)
		authn.GET("/permissions/:permissionID", provider.Authn.Permission.Controller.Get)
		authn.PATCH("/permissions/:permissionID", provider.Authn.Permission.Controller.Update)
		authn.DELETE("/permissions/:permissionID", provider.Authn.Permission.Controller.Delete)
	}
}
