package dto

// ReleaseReportSummaryDTO 发布统计概览
type ReleaseReportSummaryDTO struct {
	TotalReleases int     `json:"totalReleases"` // 总发布次数
	SuccessCount  int     `json:"successCount"`  // 成功次数
	FailedCount   int     `json:"failedCount"`   // 失败次数
	FailureRate   float64 `json:"failureRate"`   // 失败率 0-100
	RollbackCount int     `json:"rollbackCount"` // 回滚次数
	RollbackRate  float64 `json:"rollbackRate"`  // 回滚率 0-100
}

// ReleaseReportByMonthDTO 按月的发布统计（趋势）
type ReleaseReportByMonthDTO struct {
	Month         string `json:"month"`         // 2026-01
	Total         int    `json:"total"`         // 发布次数
	SuccessCount  int    `json:"successCount"`  // 成功次数
	FailedCount   int    `json:"failedCount"`   // 失败次数
	RollbackCount int    `json:"rollbackCount"` // 回滚次数
}

// ReleaseReportByProjectDTO 按项目的发布统计
type ReleaseReportByProjectDTO struct {
	ProjectID     uint   `json:"projectID"`
	ProjectName   string `json:"projectName"`
	ReleaseCount  int    `json:"releaseCount"`
	SuccessCount  int    `json:"successCount"`
	FailedCount   int    `json:"failedCount"`
	RollbackCount int    `json:"rollbackCount"`
}

// ReleaseReportResponse 报表接口统一返回
type ReleaseReportResponse struct {
	Summary  ReleaseReportSummaryDTO   `json:"summary"`
	ByMonth  []ReleaseReportByMonthDTO `json:"byMonth"`
	ByProject []ReleaseReportByProjectDTO `json:"byProject"`
	ProjectCount int `json:"projectCount"` // 有发布记录的项目数量
}
