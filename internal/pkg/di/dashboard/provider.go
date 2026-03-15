package dashboard

import (
	k8ssvc 	"valyria-backend/internal/apps/kubernetes/service"
	"valyria-backend/internal/apps/dashboard/controller"
	dashboardsvc "valyria-backend/internal/apps/dashboard/service"
)

// Provider 仪表盘模块
type Provider struct {
	Service    dashboardsvc.DashboardManager
	Controller *controller.DashboardController
}

// NewDashboardProvider 创建仪表盘 Provider；pipelineStats 可为 nil
func NewDashboardProvider(
	cluster k8ssvc.ClusterManager,
	node k8ssvc.NodeManager,
	namespace k8ssvc.NamespaceManager,
	pod k8ssvc.PodManager,
	pipelineStats dashboardsvc.PipelineStatsGetter,
) *Provider {
	svc := dashboardsvc.NewDashboardManager(cluster, node, namespace, pod, pipelineStats)
	ctrl := controller.NewDashboardController(svc)
	return &Provider{
		Service:    svc,
		Controller: ctrl,
	}
}
