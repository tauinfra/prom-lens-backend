package repository

import (
	"context"
	"time"
	"valyria-backend/internal/apps/dragon/dto"

	"gorm.io/gorm"
)

// ReleaseReportRepository 发布报表数据
type ReleaseReportRepository interface {
	GetSummary(ctx context.Context, from, to string) (dto.ReleaseReportSummaryDTO, error)
	GetByMonth(ctx context.Context, from, to string) ([]dto.ReleaseReportByMonthDTO, error)
	GetByProject(ctx context.Context, from, to string) ([]dto.ReleaseReportByProjectDTO, error)
}

type releaseReportRepository struct {
	db *gorm.DB
}

func NewReleaseReportRepository(db *gorm.DB) ReleaseReportRepository {
	return &releaseReportRepository{db: db}
}

func parseMonthRange(from, to string) (start, end time.Time, err error) {
	start, err = time.ParseInLocation("2006-01", from, time.Local)
	if err != nil {
		return
	}
	var endMonth time.Time
	endMonth, err = time.ParseInLocation("2006-01", to, time.Local)
	if err != nil {
		return
	}
	end = endMonth.AddDate(0, 1, 0) // 结束月下月 1 日 0 点（不包含）
	return
}

func (r *releaseReportRepository) GetSummary(ctx context.Context, from, to string) (dto.ReleaseReportSummaryDTO, error) {
	var res dto.ReleaseReportSummaryDTO
	start, end, err := parseMonthRange(from, to)
	if err != nil {
		return res, err
	}
	table := "valyria_dragon_release"
	q := r.db.WithContext(ctx).Table(table).Where("created_at >= ? AND created_at < ?", start, end)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return res, err
	}
	res.TotalReleases = int(total)
	if total == 0 {
		return res, nil
	}
	var successCount, failedCount, rollbackCount int64
	if err := r.db.WithContext(ctx).Table(table).Where("created_at >= ? AND created_at < ?", start, end).Where("release_status = ?", "success").Count(&successCount).Error; err != nil {
		return res, err
	}
	if err := r.db.WithContext(ctx).Table(table).Where("created_at >= ? AND created_at < ?", start, end).Where("release_status = ?", "failed").Count(&failedCount).Error; err != nil {
		return res, err
	}
	if err := r.db.WithContext(ctx).Table(table).Where("created_at >= ? AND created_at < ?", start, end).Where("target_release_id IS NOT NULL").Count(&rollbackCount).Error; err != nil {
		return res, err
	}
	res.SuccessCount = int(successCount)
	res.FailedCount = int(failedCount)
	res.RollbackCount = int(rollbackCount)
	res.FailureRate = float64(res.FailedCount) / float64(res.TotalReleases) * 100
	res.RollbackRate = float64(res.RollbackCount) / float64(res.TotalReleases) * 100
	return res, nil
}

func (r *releaseReportRepository) GetByMonth(ctx context.Context, from, to string) ([]dto.ReleaseReportByMonthDTO, error) {
	start, end, err := parseMonthRange(from, to)
	if err != nil {
		return nil, err
	}
	type row struct {
		Month         string
		Total         int
		SuccessCount  int
		FailedCount   int
		RollbackCount int
	}
	var rows []row
	// MySQL: DATE_FORMAT(created_at,'%Y-%m') as month, COUNT(*) as total, ...
	sql := `SELECT 
		DATE_FORMAT(created_at,'%Y-%m') AS month,
		COUNT(*) AS total,
		SUM(CASE WHEN release_status = 'success' THEN 1 ELSE 0 END) AS success_count,
		SUM(CASE WHEN release_status = 'failed' THEN 1 ELSE 0 END) AS failed_count,
		SUM(CASE WHEN target_release_id IS NOT NULL THEN 1 ELSE 0 END) AS rollback_count
		FROM valyria_dragon_release
		WHERE created_at >= ? AND created_at < ?
		GROUP BY DATE_FORMAT(created_at,'%Y-%m')
		ORDER BY month ASC`
	if err := r.db.WithContext(ctx).Raw(sql, start, end).Scan(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]dto.ReleaseReportByMonthDTO, 0, len(rows))
	for _, x := range rows {
		out = append(out, dto.ReleaseReportByMonthDTO{
			Month:         x.Month,
			Total:         x.Total,
			SuccessCount:  x.SuccessCount,
			FailedCount:   x.FailedCount,
			RollbackCount: x.RollbackCount,
		})
	}
	return out, nil
}

func (r *releaseReportRepository) GetByProject(ctx context.Context, from, to string) ([]dto.ReleaseReportByProjectDTO, error) {
	start, end, err := parseMonthRange(from, to)
	if err != nil {
		return nil, err
	}
	type row struct {
		ProjectID     uint
		ProjectName   string
		ReleaseCount  int
		SuccessCount  int
		FailedCount   int
		RollbackCount int
	}
	var rows []row
	sql := `SELECT 
		pr.id AS project_id,
		pr.name AS project_name,
		COUNT(*) AS release_count,
		SUM(CASE WHEN r.release_status = 'success' THEN 1 ELSE 0 END) AS success_count,
		SUM(CASE WHEN r.release_status = 'failed' THEN 1 ELSE 0 END) AS failed_count,
		SUM(CASE WHEN r.target_release_id IS NOT NULL THEN 1 ELSE 0 END) AS rollback_count
		FROM valyria_dragon_release r
		JOIN valyria_dragon_pipeline p ON p.id = r.pipeline_id
		JOIN valyria_dragon_environment e ON e.id = p.environment_id
		JOIN valyria_dragon_project pr ON pr.id = e.project_id
		WHERE r.created_at >= ? AND r.created_at < ?
		GROUP BY pr.id, pr.name
		ORDER BY release_count DESC`
	if err := r.db.WithContext(ctx).Raw(sql, start, end).Scan(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]dto.ReleaseReportByProjectDTO, 0, len(rows))
	for _, x := range rows {
		out = append(out, dto.ReleaseReportByProjectDTO{
			ProjectID:     x.ProjectID,
			ProjectName:   x.ProjectName,
			ReleaseCount:  x.ReleaseCount,
			SuccessCount:  x.SuccessCount,
			FailedCount:   x.FailedCount,
			RollbackCount: x.RollbackCount,
		})
	}
	return out, nil
}
