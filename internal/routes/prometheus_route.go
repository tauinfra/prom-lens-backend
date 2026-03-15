package routes

import (
	"valyria-backend/internal/pkg/di"

	"github.com/gin-gonic/gin"
)

func PrometheusRouters(rg *gin.RouterGroup, provider *di.Provider) {
	prometheus := rg.Group("/prometheus")
	{
		// Group Routes
		prometheus.GET("/groups", provider.Prometheus.Group.Controller.List)
		prometheus.GET("/groups/:groupID", provider.Prometheus.Group.Controller.Get)
		prometheus.POST("/groups", provider.Prometheus.Group.Controller.Create)
		prometheus.PATCH("/groups/:groupID", provider.Prometheus.Group.Controller.Update)
		prometheus.DELETE("/groups/:groupID", provider.Prometheus.Group.Controller.Delete)

		// Target Group Routes
		prometheus.GET("/target-groups", provider.Prometheus.TargetGroup.Controller.List)
		prometheus.GET("/target-groups/:groupID", provider.Prometheus.TargetGroup.Controller.Get)
		prometheus.POST("/target-groups", provider.Prometheus.TargetGroup.Controller.Create)
		prometheus.PATCH("/target-groups/:groupID", provider.Prometheus.TargetGroup.Controller.Update)
		prometheus.DELETE("/target-groups/:groupID", provider.Prometheus.TargetGroup.Controller.Delete)

		// Target Routes
		prometheus.GET("/target-groups/:groupID/targets", provider.Prometheus.Target.Controller.List)
		prometheus.GET("/target-groups/:groupID/targets/:targetID", provider.Prometheus.Target.Controller.Get)
		prometheus.POST("/target-groups/:groupID/targets", provider.Prometheus.Target.Controller.Create)
		prometheus.PATCH("/target-groups/:groupID/targets/:targetID", provider.Prometheus.Target.Controller.Update)
		prometheus.DELETE("/target-groups/:groupID/targets/:targetID", provider.Prometheus.Target.Controller.Delete)

		// Record Routes
		prometheus.GET("/groups/:groupID/records", provider.Prometheus.Record.Controller.List)
		prometheus.GET("/groups/:groupID/records/:recordID", provider.Prometheus.Record.Controller.Get)
		prometheus.POST("/groups/:groupID/records", provider.Prometheus.Record.Controller.Create)
		prometheus.PATCH("/groups/:groupID/records/:recordID", provider.Prometheus.Record.Controller.Update)
		prometheus.DELETE("/groups/:groupID/records/:recordID", provider.Prometheus.Record.Controller.Delete)

		// Rule Routes
		prometheus.GET("/groups/:groupID/rules", provider.Prometheus.Rule.Controller.List)
		prometheus.GET("/groups/:groupID/rules/:ruleID", provider.Prometheus.Rule.Controller.Get)
		prometheus.POST("/groups/:groupID/rules", provider.Prometheus.Rule.Controller.Create)
		prometheus.PATCH("/groups/:groupID/rules/:ruleID", provider.Prometheus.Rule.Controller.Update)
		prometheus.DELETE("/groups/:groupID/rules/:ruleID", provider.Prometheus.Rule.Controller.Delete)
	}
}
