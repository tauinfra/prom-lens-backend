package dto

// DashboardCluster 仪表盘集群概览
type DashboardCluster struct {
	Clusters   int `json:"clusters"`
	Nodes      int `json:"nodes"`
	Namespaces int `json:"namespaces"`
	Pods       int `json:"pods"`
}

// DashboardResources 仪表盘资源使用（预留，当前返回空）
type DashboardResources struct {
	CPUUsage    int `json:"cpu_usage,omitempty"`
	MemoryUsage int `json:"memory_usage,omitempty"`
	PodCount    int `json:"pod_count,omitempty"`
	NodeErrors  int `json:"node_errors,omitempty"`
}

// DashboardPipelines 仪表盘流水线概览（今日 00:00–23:59 发布统计，不区分环境/项目）
type DashboardPipelines struct {
	TodayTotal int `json:"today_total"`
	Success    int `json:"success"`
	Failed     int `json:"failed"`
	Rollback   int `json:"rollback"`
}

// DashboardResponse 仪表盘接口返回
type DashboardResponse struct {
	Cluster   DashboardCluster   `json:"cluster"`
	Resources DashboardResources `json:"resources"`
	Pipelines DashboardPipelines `json:"pipelines"`
}

// PipelineTrendPoint 流水线趋势单点（按小时或按天）
type PipelineTrendPoint struct {
	Hour          string `json:"hour,omitempty"`
	Date          string `json:"date,omitempty"`
	SuccessCount  int    `json:"successCount"`
	FailedCount   int    `json:"failedCount"`
	RollbackCount int    `json:"rollbackCount"`
}

// PipelineTrendResponse 流水线趋势接口返回
type PipelineTrendResponse struct {
	Range  string               `json:"range"`
	Points []PipelineTrendPoint `json:"points"`
}

// PipelineRecentItem 最近发布单条（dashboard/pipelines/recent）
type PipelineRecentItem struct {
	ProjectName string `json:"projectName"`
	Pipeline    string `json:"pipeline"`
	Env         string `json:"env"`
	Status      string `json:"status"`
	Duration    int64  `json:"duration"` // 秒，未结束为 0
	Time        string `json:"time"`     // 显示时间，优先 FinishedAt 否则 CreatedAt
}

// PipelineProjectStatsItem 项目维度流水线统计（dashboard/pipelines/projects）
type PipelineProjectStatsItem struct {
	ProjectName   string `json:"projectName"`
	PipelineCount int    `json:"pipelineCount"`
	ReleaseCount  int    `json:"releaseCount"`
	SuccessCount  int    `json:"successCount"`
	FailedCount   int    `json:"failedCount"`
	RollbackCount int    `json:"rollbackCount"`
}
