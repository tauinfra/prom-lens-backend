package routes

import (
	"valyria-backend/internal/pkg/di"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func FoundryRouters(rg *gin.RouterGroup, tx *gorm.DB, provider *di.Provider) {
	foundry := rg.Group("/foundry")
	{
		foundry.GET("/projects", provider.Foundry.Project.Controller.List)
		foundry.GET("/projects/:projectID", provider.Foundry.Project.Controller.Get)
		foundry.POST("/projects", provider.Foundry.Project.Controller.Create)
		foundry.PATCH("/projects/:projectID", provider.Foundry.Project.Controller.Update)
		foundry.DELETE("/projects/:projectID", provider.Foundry.Project.Controller.Delete)
		foundry.GET("/projects/:projectID/environments", provider.Foundry.Environment.Controller.List)
		foundry.GET("/projects/:projectID/environments/:envID", provider.Foundry.Environment.Controller.Get)
		foundry.POST("/projects/:projectID/environments", provider.Foundry.Environment.Controller.Create)
		foundry.PATCH("/projects/:projectID//environments/:envID", provider.Foundry.Environment.Controller.Update)
		foundry.DELETE("/projects/:projectID/environments/:envID", provider.Foundry.Environment.Controller.Delete)
		// Pipeline Routes
		foundry.GET("/projects/:projectID/environments/:envID/pipelines", provider.Foundry.Pipeline.Controller.List)
		foundry.GET("/projects/:projectID/environments/:envID/pipelines/:pipelineID", provider.Foundry.Pipeline.Controller.Get)
		foundry.POST("/projects/:projectID/environments/:envID/pipelines", provider.Foundry.Pipeline.Controller.Create)
		foundry.PATCH("/projects/:projectID/environments/:envID/pipelines/:pipelineID", provider.Foundry.Pipeline.Controller.Update)
		foundry.DELETE("/projects/:projectID/environments/:envID/pipelines/:pipelineID", provider.Foundry.Pipeline.Controller.Delete)
		// Release Routes
		foundry.GET("/projects/:projectID/environments/:envID/pipelines/:pipelineID/releases", provider.Foundry.Release.Controller.List)
		foundry.GET("/projects/:projectID/environments/:envID/pipelines/:pipelineID/releases/:releaseID", provider.Foundry.Release.Controller.Get)
		foundry.POST("/projects/:projectID/environments/:envID/pipelines/:pipelineID/releases", provider.Foundry.Release.Controller.Create)
		foundry.PATCH("/projects/:projectID/environments/:envID/pipelines/:pipelineID/releases/:releaseID", provider.Foundry.Release.Controller.Update)
		foundry.DELETE("/projects/:projectID/environments/:envID/pipelines/:pipelineID/releases/:releaseID", provider.Foundry.Release.Controller.Delete)

		foundry.GET("/pipelines/:pipelineID/releases", provider.Foundry.Release.Controller.List)
		foundry.GET("/pipelines/:pipelineID/releases/:releaseID", provider.Foundry.Release.Controller.Get)
		foundry.GET("/pipelines/:pipelineID/releases/:releaseID/logs", provider.Foundry.Release.Controller.GetLogs)
		// App Routes
		foundry.GET("/projects/:projectID/applications", provider.Foundry.Application.Controller.List)
		foundry.GET("/projects/:projectID/applications/:appID", provider.Foundry.Application.Controller.Get)
		foundry.POST("/projects/:projectID/applications", provider.Foundry.Application.Controller.Create)
		foundry.PATCH("/projects/:projectID/applications/:appID", provider.Foundry.Application.Controller.Update)
		foundry.DELETE("/projects/:projectID/applications/:appID", provider.Foundry.Application.Controller.Delete)
		// Gitlab Routes
		foundry.GET("/gitlab/:gitlabID/groups", provider.Foundry.Gitlab.Controller.ListGitlabGroups)
		foundry.GET("/gitlab/:gitlabID/groups/:groupID/projects", provider.Foundry.Gitlab.Controller.ListGroupProjects)
		foundry.GET("/gitlab/:gitlabID/projects/:projectID/branches", provider.Foundry.Gitlab.Controller.ListProjectBranches)
		foundry.GET("/gitlab/:gitlabID/projects/:projectID/branches/:branch", provider.Foundry.Gitlab.Controller.GetRefBranch)
		foundry.GET("/gitlab/:gitlabID/projects/:projectID/tags", provider.Foundry.Gitlab.Controller.ListProjectTags)
		foundry.GET("/gitlab/:gitlabID/projects/:projectID/tags/:tag", provider.Foundry.Gitlab.Controller.GetRefTag)
	}
	credentials := rg.Group("/foundry/credentials")
	{
		// Cred Gitlab Routes
		var resource = "credentials:gitlab"
		credentials.GET("/gitlab", provider.PermissionManager.FoundryPermission(tx, resource, "list"), provider.Foundry.CredGitlab.Controller.List)
		credentials.GET("/gitlab/:id", provider.PermissionManager.FoundryPermission(tx, resource, "get"), provider.Foundry.CredGitlab.Controller.Get)
		credentials.POST("/gitlab", provider.PermissionManager.FoundryPermission(tx, resource, "create"), provider.Foundry.CredGitlab.Controller.Create)
		credentials.PATCH("/gitlab/:id", provider.PermissionManager.FoundryPermission(tx, resource, "update"), provider.Foundry.CredGitlab.Controller.Update)
		credentials.DELETE("/gitlab/:id", provider.PermissionManager.FoundryPermission(tx, resource, "delete"), provider.Foundry.CredGitlab.Controller.Delete)
		// Cred Harbor Routes
		resource = "credentials:harbor"
		credentials.GET("/harbor", provider.PermissionManager.FoundryPermission(tx, resource, "list"), provider.Foundry.CredHarbor.Controller.List)
		credentials.GET("/harbor/:id", provider.PermissionManager.FoundryPermission(tx, resource, "get"), provider.Foundry.CredHarbor.Controller.Get)
		credentials.POST("/harbor", provider.PermissionManager.FoundryPermission(tx, resource, "create"), provider.Foundry.CredHarbor.Controller.Create)
		credentials.PATCH("/harbor/:id", provider.PermissionManager.FoundryPermission(tx, resource, "update"), provider.Foundry.CredHarbor.Controller.Update)
		credentials.DELETE("/harbor/:id", provider.PermissionManager.FoundryPermission(tx, resource, "delete"), provider.Foundry.CredHarbor.Controller.Delete)
		// Cred ArgoCD Routes
		resource = "credentials:argocd"
		credentials.GET("/argocd", provider.PermissionManager.FoundryPermission(tx, resource, "list"), provider.Foundry.CredArgoCD.Controller.List)
		credentials.GET("/argocd/:id", provider.PermissionManager.FoundryPermission(tx, resource, "get"), provider.Foundry.CredArgoCD.Controller.Get)
		credentials.POST("/argocd", provider.PermissionManager.FoundryPermission(tx, resource, "create"), provider.Foundry.CredArgoCD.Controller.Create)
		credentials.PATCH("/argocd/:id", provider.PermissionManager.FoundryPermission(tx, resource, "update"), provider.Foundry.CredArgoCD.Controller.Update)
		credentials.DELETE("/argocd/:id", provider.PermissionManager.FoundryPermission(tx, resource, "delete"), provider.Foundry.CredArgoCD.Controller.Delete)
	}
}
