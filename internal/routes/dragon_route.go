package routes

import (
	"valyria-backend/internal/pkg/di"

	"github.com/gin-gonic/gin"
)

func DragonRouters(rg *gin.RouterGroup, provider *di.Provider) {
	dragon := rg.Group("/dragon")
	{
		// Projects Routes
		dragon.GET("/projects", provider.Dragon.Project.Controller.List)
		dragon.GET("/projects/:projectID", provider.Dragon.Project.Controller.Get)
		dragon.POST("/projects",
			provider.PermissionManager.RequirePermission("dragon:project:create"),
			provider.Dragon.Project.Controller.Create,
		)
		dragon.PATCH("/projects/:projectID",
			provider.PermissionManager.RequirePermission("dragon:project:update"),
			provider.Dragon.Project.Controller.Update,
		)
		dragon.DELETE("/projects/:projectID",
			provider.PermissionManager.RequirePermission("dragon:project:delete"),
			provider.Dragon.Project.Controller.Delete,
		)
		// Env Routes
		dragon.GET("/projects/:projectID/environments", provider.Dragon.Environment.Controller.List)
		dragon.GET("/projects/:projectID/environments/:environmentID", provider.Dragon.Environment.Controller.Get)
		dragon.POST("/projects/:projectID/environments",
			provider.PermissionManager.RequirePermission("dragon:environment:create"),
			provider.Dragon.Environment.Controller.Create,
		)
		dragon.PATCH("/projects/:projectID/environments/:environmentID",
			provider.PermissionManager.RequirePermission("dragon:environment:update"),
			provider.Dragon.Environment.Controller.Update,
		)
		dragon.DELETE("/projects/:projectID/environments/:environmentID",
			provider.PermissionManager.RequirePermission("dragon:environment:delete"),
			provider.Dragon.Environment.Controller.Delete,
		)
		// Pipeline Routes
		dragon.GET("/projects/:projectID/environments/:environmentID/pipelines", provider.Dragon.Pipeline.Controller.List)
		dragon.GET("/projects/:projectID/environments/:environmentID/pipelines/:pipelineID", provider.Dragon.Pipeline.Controller.Get)
		dragon.POST("/projects/:projectID/environments/:environmentID/pipelines",
			provider.PermissionManager.RequirePermission("dragon:pipeline:create"),
			provider.Dragon.Pipeline.Controller.Create,
		)
		dragon.PATCH("/projects/:projectID/environments/:environmentID/pipelines/:pipelineID",
			provider.PermissionManager.RequirePermission("dragon:pipeline:update"),
			provider.Dragon.Pipeline.Controller.Update,
		)
		dragon.DELETE("/projects/:projectID/environments/:environmentID/pipelines/:pipelineID",
			provider.PermissionManager.RequirePermission("dragon:pipeline:delete"),
			provider.Dragon.Pipeline.Controller.Delete,
		)
		dragon.GET("/projects/:projectID/environments/:environmentID/pipelines/:pipelineID/acls", provider.Dragon.PipelineACL.Controller.List)
		dragon.POST("/projects/:projectID/environments/:environmentID/pipelines/:pipelineID/acls",
			provider.PermissionManager.RequirePermission("dragon:acl:create"),
			provider.Dragon.PipelineACL.Controller.Create,
		)
		dragon.POST("/projects/:projectID/environments/:environmentID/pipelines/:pipelineID/acls/batch",
			provider.PermissionManager.RequirePermission("dragon:acl:create"),
			provider.Dragon.PipelineACL.Controller.BatchCreate,
		)
		dragon.DELETE("/projects/:projectID/environments/:environmentID/pipelines/:pipelineID/acls/:aclID",
			provider.PermissionManager.RequirePermission("dragon:acl:delete"),
			provider.Dragon.PipelineACL.Controller.Delete,
		)
		// Release Routes
		dragon.GET("/projects/:projectID/environments/:environmentID/pipelines/:pipelineID/releases", provider.Dragon.Release.Controller.List)
		dragon.GET("/projects/:projectID/environments/:environmentID/pipelines/:pipelineID/releases/:releaseID", provider.Dragon.Release.Controller.Get)
		dragon.POST("/projects/:projectID/environments/:environmentID/pipelines/:pipelineID/releases", provider.Dragon.Release.Controller.Create)
		dragon.POST("/projects/:projectID/environments/:environmentID/pipelines/:pipelineID/rollback", provider.Dragon.Release.Controller.Rollback)
		dragon.GET("/projects/:projectID/environments/:environmentID/pipelines/:pipelineID/releases/:releaseID/logs", provider.Dragon.Release.Controller.GetLogs)
		dragon.DELETE("/projects/:projectID/environments/:environmentID/pipelines/:pipelineID/releases/:releaseID",
			provider.PermissionManager.RequirePermission("dragon:release:delete"),
			provider.Dragon.Release.Controller.Delete,
		)
		// Report Routes
		dragon.GET("/reports/summary", provider.Dragon.Report.Controller.Summary)
		// Review Routes
		dragon.GET("/projects/:projectID/environments/:environmentID/pipelines/:pipelineID/releases/:releaseID/reviews", provider.Dragon.Review.Controller.List)
		dragon.POST("/projects/:projectID/environments/:environmentID/pipelines/:pipelineID/releases/:releaseID/reviews", provider.Dragon.Review.Controller.Create)
		dragon.GET("/projects/:projectID/environments/:environmentID/pipelines/:pipelineID/releases/:releaseID/reviews/:reviewID", provider.Dragon.Review.Controller.Get)
		dragon.PATCH("/projects/:projectID/environments/:environmentID/pipelines/:pipelineID/releases/:releaseID/reviews/:reviewID", provider.Dragon.Review.Controller.Update)
		dragon.DELETE("/projects/:projectID/environments/:environmentID/pipelines/:pipelineID/releases/:releaseID/reviews/:reviewID", provider.Dragon.Review.Controller.Delete)
		// Review Config Routes
		dragon.GET("/review-stages", provider.Dragon.ReviewConfig.Controller.ListGlobal)
		dragon.POST("/review-stages",
			provider.PermissionManager.RequirePermission("dragon:review-stage:create"),
			provider.Dragon.ReviewConfig.Controller.CreateGlobal,
		)
		dragon.GET("/review-stages/:stageID", provider.Dragon.ReviewConfig.Controller.GetGlobal)
		dragon.PATCH("/review-stages/:stageID",
			provider.PermissionManager.RequirePermission("dragon:review-stage:update"),
			provider.Dragon.ReviewConfig.Controller.UpdateGlobal,
		)
		dragon.DELETE("/review-stages/:stageID",
			provider.PermissionManager.RequirePermission("dragon:review-stage:delete"),
			provider.Dragon.ReviewConfig.Controller.DeleteGlobal,
		)
		// Tekton Pipeline Routes
		dragon.GET("/tekton/pipelines", provider.Dragon.TektonPipeline.Controller.List)
		dragon.GET("/tekton/pipelines/:name", provider.Dragon.TektonPipeline.Controller.Get)
		dragon.POST("/tekton/pipelines",
			provider.PermissionManager.RequirePermission("dragon:tekton-pipeline:create"),
			provider.Dragon.TektonPipeline.Controller.Create,
		)
		dragon.PATCH("/tekton/pipelines/:name",
			provider.PermissionManager.RequirePermission("dragon:tekton-pipeline:update"),
			provider.Dragon.TektonPipeline.Controller.Update,
		)
		dragon.DELETE("/tekton/pipelines/:name",
			provider.PermissionManager.RequirePermission("dragon:tekton-pipeline:delete"),
			provider.Dragon.TektonPipeline.Controller.Delete,
		)
		dragon.GET("/tekton/tasks", provider.Dragon.TektonTask.Controller.List)
		dragon.GET("/tekton/tasks/:name", provider.Dragon.TektonTask.Controller.Get)
		dragon.POST("/tekton/tasks",
			provider.PermissionManager.RequirePermission("dragon:tekton-task:create"),
			provider.Dragon.TektonTask.Controller.Create,
		)
		dragon.PATCH("/tekton/tasks/:name",
			provider.PermissionManager.RequirePermission("dragon:tekton-task:update"),
			provider.Dragon.TektonTask.Controller.Update,
		)
		dragon.DELETE("/tekton/tasks/:name",
			provider.PermissionManager.RequirePermission("dragon:tekton-task:delete"),
			provider.Dragon.TektonTask.Controller.Delete,
		)
		dragon.GET("/tekton/taskruns", provider.Dragon.TektonTaskRun.Controller.List)
		dragon.GET("/tekton/taskruns/:name", provider.Dragon.TektonTaskRun.Controller.Get)
		dragon.POST("/tekton/taskruns",
			provider.PermissionManager.RequirePermission("dragon:tekton-taskrun:create"),
			provider.Dragon.TektonTaskRun.Controller.Create,
		)
		dragon.PATCH("/tekton/taskruns/:name",
			provider.PermissionManager.RequirePermission("dragon:tekton-taskrun:update"),
			provider.Dragon.TektonTaskRun.Controller.Update,
		)
		dragon.DELETE("/tekton/taskruns/:name",
			provider.PermissionManager.RequirePermission("dragon:tekton-taskrun:delete"),
			provider.Dragon.TektonTaskRun.Controller.Delete,
		)
		dragon.POST("/tekton/taskruns/batch-delete",
			provider.PermissionManager.RequirePermission("dragon:tekton-taskrun:delete"),
			provider.Dragon.TektonTaskRun.Controller.BatchDelete,
		)
		dragon.GET("/tekton/pipelineruns", provider.Dragon.TektonPipelineRun.Controller.List)
		dragon.GET("/tekton/pipelineruns/:name", provider.Dragon.TektonPipelineRun.Controller.Get)
		dragon.POST("/tekton/pipelineruns",
			provider.PermissionManager.RequirePermission("dragon:tekton-pipelinerun:create"),
			provider.Dragon.TektonPipelineRun.Controller.Create,
		)
		dragon.PATCH("/tekton/pipelineruns/:name",
			provider.PermissionManager.RequirePermission("dragon:tekton-pipelinerun:update"),
			provider.Dragon.TektonPipelineRun.Controller.Update,
		)
		dragon.DELETE("/tekton/pipelineruns/:name",
			provider.PermissionManager.RequirePermission("dragon:tekton-pipelinerun:delete"),
			provider.Dragon.TektonPipelineRun.Controller.Delete,
		)
		dragon.POST("/tekton/pipelineruns/batch-delete",
			provider.PermissionManager.RequirePermission("dragon:tekton-pipelinerun:delete"),
			provider.Dragon.TektonPipelineRun.Controller.BatchDelete,
		)
		// Credentials Routes
		dragon.GET("/credentials", provider.Dragon.Credential.Controller.List)
		dragon.GET("/credentials/:credentialID", provider.Dragon.Credential.Controller.Get)
		dragon.POST("/credentials",
			provider.PermissionManager.RequirePermission("dragon:credential:create"),
			provider.Dragon.Credential.Controller.Create,
		)
		dragon.PATCH("/credentials/:credentialID",
			provider.PermissionManager.RequirePermission("dragon:credential:update"),
			provider.Dragon.Credential.Controller.Update,
		)
		dragon.DELETE("/credentials/:credentialID",
			provider.PermissionManager.RequirePermission("dragon:credential:delete"),
			provider.Dragon.Credential.Controller.Delete,
		)
		// Gitlab Routes
		dragon.GET("/gitlab/:gitlabID/groups", provider.Dragon.Gitlab.Controller.ListGitlabGroups)
		dragon.GET("/gitlab/:gitlabID/groups/:groupID/projects", provider.Dragon.Gitlab.Controller.ListGroupProjects)
		dragon.GET("/gitlab/:gitlabID/projects/:projectID/branches", provider.Dragon.Gitlab.Controller.ListProjectBranches)
		dragon.GET("/gitlab/:gitlabID/projects/:projectID/branches/:branch", provider.Dragon.Gitlab.Controller.GetRefBranch)
		dragon.GET("/gitlab/:gitlabID/projects/:projectID/tags", provider.Dragon.Gitlab.Controller.ListProjectTags)
		dragon.GET("/gitlab/:gitlabID/projects/:projectID/tags/:tag", provider.Dragon.Gitlab.Controller.GetRefTag)

	}
}
