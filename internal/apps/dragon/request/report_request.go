package request

// ReleaseReportRequest 发布报表查询参数（query）
// from / to 格式：2006-01，不传则默认最近 12 个月
type ReleaseReportRequest struct {
	From string `form:"from"` // 起始月 2026-01
	To   string `form:"to"`   // 结束月 2026-02
}
