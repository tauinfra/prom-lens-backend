package service

import (
	"context"
	"errors"
	"valyria-backend/internal/apps/dashboard/dto"
	k8ssvc "valyria-backend/internal/apps/kubernetes/service"
	pg "valyria-backend/internal/core/pagination"
)

// PipelineStatsGetter 今日发布统计与趋势（由 dragon.ReleaseManager 实现）
type PipelineStatsGetter interface {
	GetTodayStats(ctx context.Context) (todayTotal, success, failed, rollback int, err error)
	GetPipelineTrend(ctx context.Context, rangeParam string) (*dto.PipelineTrendResponse, error)
	GetRecentReleases(ctx context.Context, envName string, limit int) ([]dto.PipelineRecentItem, error)
	GetProjectPipelineStats(ctx context.Context) ([]dto.PipelineProjectStatsItem, error)
}

// DashboardManager 仪表盘聚合数据
type DashboardManager interface {
	GetDashboard(ctx context.Context) (dto.DashboardResponse, error)
	GetPipelinesTrend(ctx context.Context, rangeParam string) (*dto.PipelineTrendResponse, error)
	GetPipelinesRecent(ctx context.Context, env string, limit int) ([]dto.PipelineRecentItem, error)
	GetPipelinesProjects(ctx context.Context) ([]dto.PipelineProjectStatsItem, error)
}

type dashboardManager struct {
	cluster       k8ssvc.ClusterManager
	node          k8ssvc.NodeManager
	namespace     k8ssvc.NamespaceManager
	pod           k8ssvc.PodManager
	pipelineStats PipelineStatsGetter
}

// NewDashboardManager 创建仪表盘服务；pipelineStats 可为 nil，为 nil 时 pipelines 返回零值
func NewDashboardManager(
	cluster k8ssvc.ClusterManager,
	node k8ssvc.NodeManager,
	namespace k8ssvc.NamespaceManager,
	pod k8ssvc.PodManager,
	pipelineStats PipelineStatsGetter,
) DashboardManager {
	return &dashboardManager{
		cluster:       cluster,
		node:          node,
		namespace:     namespace,
		pod:           pod,
		pipelineStats: pipelineStats,
	}
}

// GetDashboard 聚合集群/节点/命名空间/Pod 数量；resources 预留；pipelines 为今日发布统计
func (s *dashboardManager) GetDashboard(ctx context.Context) (dto.DashboardResponse, error) {
	clusters, _, err := s.cluster.List(ctx, pg.QueryParams{Page: 1, Size: 10000})
	if err != nil {
		return dto.DashboardResponse{}, err
	}

	clusterCount := len(clusters)
	var totalNodes, totalNamespaces, totalPods int

	for _, c := range clusters {
		clusterID := uint(c.ID)
		nodes, err := s.node.List(ctx, clusterID)
		if err != nil {
			continue
		}
		totalNodes += len(nodes)

		namespaces, err := s.namespace.List(ctx, clusterID)
		if err != nil {
			continue
		}
		totalNamespaces += len(namespaces)

		pods, err := s.pod.ListAll(ctx, clusterID, "", "")
		if err != nil {
			continue
		}
		totalPods += len(pods)
	}

	resp := dto.DashboardResponse{
		Cluster: dto.DashboardCluster{
			Clusters:   clusterCount,
			Nodes:      totalNodes,
			Namespaces: totalNamespaces,
			Pods:       totalPods,
		},
		Resources: dto.DashboardResources{},
		Pipelines: dto.DashboardPipelines{},
	}
	if s.pipelineStats != nil {
		todayTotal, success, failed, rollback, err := s.pipelineStats.GetTodayStats(ctx)
		if err == nil {
			resp.Pipelines = dto.DashboardPipelines{
				TodayTotal: todayTotal,
				Success:    success,
				Failed:     failed,
				Rollback:   rollback,
			}
		}
	}
	return resp, nil
}

// GetPipelinesTrend 流水线趋势，rangeParam=today|7d|30d
func (s *dashboardManager) GetPipelinesTrend(ctx context.Context, rangeParam string) (*dto.PipelineTrendResponse, error) {
	if s.pipelineStats == nil {
		return nil, errors.New("pipeline stats not available")
	}
	return s.pipelineStats.GetPipelineTrend(ctx, rangeParam)
}

// GetPipelinesRecent 指定环境最近发布记录
func (s *dashboardManager) GetPipelinesRecent(ctx context.Context, env string, limit int) ([]dto.PipelineRecentItem, error) {
	if s.pipelineStats == nil {
		return nil, errors.New("pipeline stats not available")
	}
	return s.pipelineStats.GetRecentReleases(ctx, env, limit)
}

// GetPipelinesProjects 按项目统计发布次数、pipeline 数、成功/失败/回滚次数
func (s *dashboardManager) GetPipelinesProjects(ctx context.Context) ([]dto.PipelineProjectStatsItem, error) {
	if s.pipelineStats == nil {
		return nil, errors.New("pipeline stats not available")
	}
	return s.pipelineStats.GetProjectPipelineStats(ctx)
}
