package routes

import (
	"valyria-backend/internal/pkg/di"

	"github.com/gin-gonic/gin"
)

func DashboardRouters(rg *gin.RouterGroup, provider *di.Provider) {
	rg.GET("/dashboard", provider.Dashboard.Controller.Get)
	rg.GET("/dashboard/pipelines/trend", provider.Dashboard.Controller.GetPipelinesTrend)
	rg.GET("/dashboard/pipelines/recent", provider.Dashboard.Controller.GetPipelinesRecent)
	rg.GET("/dashboard/pipelines/projects", provider.Dashboard.Controller.GetPipelinesProjects)
}
