package repository

import (
	"context"
	"time"
	"valyria-backend/internal/apps/dragon/model"
	pg "valyria-backend/internal/core/pagination"

	"gorm.io/gorm"
)

// TodayPipelineStats 今日发布统计（不区分环境、项目）
type TodayPipelineStats struct {
	TodayTotal int64
	Success    int64
	Failed     int64
	Rollback   int64
}

// TrendHourRow 按小时聚合的一行（用于 today）
type TrendHourRow struct {
	Hour         int   `gorm:"column:hour"`
	SuccessCount int64 `gorm:"column:success_count"`
	FailedCount  int64 `gorm:"column:failed_count"`
	RollbackCount int64 `gorm:"column:rollback_count"`
}

// TrendDayRow 按天聚合的一行（用于 7d/30d）
type TrendDayRow struct {
	Date         string `gorm:"column:date"`
	SuccessCount int64  `gorm:"column:success_count"`
	FailedCount  int64  `gorm:"column:failed_count"`
	RollbackCount int64 `gorm:"column:rollback_count"`
}

// ReleaseRepository 定义了数据访问层的接口
type ReleaseRepository interface {
	List(ctx context.Context, projectID, environmentID, pipelineID uint, params pg.QueryParams) ([]model.Release, pg.Pagination, error)
	Get(ctx context.Context, pipelineID, releaseID uint) (model.Release, error)
	Create(ctx context.Context, data *model.Release) error
	Update(ctx context.Context, pipelineID, releaseID uint, data *model.Release) error
	Delete(ctx context.Context, pipelineID, releaseID uint) error
	ExistsUnfinishedByPipelineID(ctx context.Context, pipelineID uint) (bool, error)
	GetTodayStats(ctx context.Context) (TodayPipelineStats, error)
	GetTrendByHour(ctx context.Context, dayStart, dayEnd time.Time) ([]TrendHourRow, error)
	GetTrendByDay(ctx context.Context, startDate, endDate time.Time) ([]TrendDayRow, error)
	ListRecentByEnvName(ctx context.Context, envName string, limit int) ([]model.Release, error)
	ListProjectPipelineStats(ctx context.Context) ([]ProjectPipelineStatsRow, error)
	WithTx(tx *gorm.DB) ReleaseRepository
}

// ProjectPipelineStatsRow 按项目聚合的流水线/发布统计
type ProjectPipelineStatsRow struct {
	ProjectID     uint   `gorm:"column:project_id"`
	ProjectName   string `gorm:"column:project_name"`
	PipelineCount int64  `gorm:"column:pipeline_count"`
	ReleaseCount  int64  `gorm:"column:release_count"`
	SuccessCount  int64  `gorm:"column:success_count"`
	FailedCount   int64  `gorm:"column:failed_count"`
	RollbackCount int64  `gorm:"column:rollback_count"`
}

// releaseRepository 实现了 ReleaseRepository 接口
type releaseRepository struct {
	db *gorm.DB
}

// NewReleaseRepository 创建新的 ReleaseRepository 实例
func NewReleaseRepository(db *gorm.DB) ReleaseRepository {
	return &releaseRepository{db: db}
}

// List 查询列表
func (r *releaseRepository) List(ctx context.Context, projectID, environmentID, pipelineID uint, params pg.QueryParams) (data []model.Release, pagination pg.Pagination, err error) {
	query := r.db.WithContext(ctx).
		Model(&model.Release{}).
		Joins("JOIN valyria_dragon_pipeline p ON p.id = valyria_dragon_release.pipeline_id").
		Joins("JOIN valyria_dragon_environment e ON e.id = p.environment_id").
		Where("valyria_dragon_release.pipeline_id = ? AND p.environment_id = ? AND e.project_id = ?", pipelineID, environmentID, projectID)
	// 分页查询
	if pagination, err = pg.Paginate(query, &data, params); err != nil {
		return nil, pg.Pagination{}, err
	}
	return
}

// Get 查询
func (r *releaseRepository) Get(ctx context.Context, pipelineID, releaseID uint) (model.Release, error) {
	var data model.Release
	if err := r.db.WithContext(ctx).
		Where("pipeline_id = ? AND id = ?", pipelineID, releaseID).
		Preload("Pipeline.Environment.Project").
		First(&data).Error; err != nil {
		return data, err
	}
	return data, nil
}

// Create 创建
func (r *releaseRepository) Create(ctx context.Context, data *model.Release) error {
	// 创建
	if err := r.db.WithContext(ctx).Create(data).Error; err != nil {
		return err
	}
	return nil
}

// Update 更新
func (r *releaseRepository) Update(ctx context.Context, pipelineID, releaseID uint, data *model.Release) error {
	if err := r.db.WithContext(ctx).Where("pipeline_id = ? AND id = ?", pipelineID, releaseID).Updates(data).Error; err != nil {
		return err
	}
	return nil
}

// Delete 删除
func (r *releaseRepository) Delete(ctx context.Context, pipelineID, releaseID uint) error {
	return r.db.WithContext(ctx).Where("pipeline_id = ? AND id = ?", pipelineID, releaseID).Delete(&model.Release{}).Error
}

// ExistsUnfinishedByPipelineID 是否存在未结束的发布（审批中或发布中）
func (r *releaseRepository) ExistsUnfinishedByPipelineID(ctx context.Context, pipelineID uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Release{}).
		Where("pipeline_id = ?", pipelineID).
		Where("release_status IN ? OR approval_status = ?", []string{"pending", "progressing"}, "pending").
		Count(&count).Error
	return count > 0, err
}

// GetTodayStats 今日 00:00–23:59 发布统计（不区分环境、项目）
func (r *releaseRepository) GetTodayStats(ctx context.Context) (TodayPipelineStats, error) {
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	todayEnd := todayStart.Add(24*time.Hour - time.Nanosecond)

	var stats TodayPipelineStats
	dayFilter := func(db *gorm.DB) *gorm.DB {
		return db.Where("created_at >= ? AND created_at <= ?", todayStart, todayEnd)
	}

	if err := dayFilter(r.db.WithContext(ctx).Model(&model.Release{})).Count(&stats.TodayTotal).Error; err != nil {
		return stats, err
	}
	if err := dayFilter(r.db.WithContext(ctx).Model(&model.Release{})).Where("release_status = ?", "success").Count(&stats.Success).Error; err != nil {
		return stats, err
	}
	if err := dayFilter(r.db.WithContext(ctx).Model(&model.Release{})).Where("release_status = ?", "failed").Count(&stats.Failed).Error; err != nil {
		return stats, err
	}
	if err := dayFilter(r.db.WithContext(ctx).Model(&model.Release{})).Where("target_release_id IS NOT NULL").Count(&stats.Rollback).Error; err != nil {
		return stats, err
	}
	return stats, nil
}

// GetTrendByHour 按小时聚合发布数（dayStart/dayEnd 为当日 00:00 至次日 00:00 前）
func (r *releaseRepository) GetTrendByHour(ctx context.Context, dayStart, dayEnd time.Time) ([]TrendHourRow, error) {
	var rows []TrendHourRow
	// MySQL: HOUR(created_at) 返回 0-23
	err := r.db.WithContext(ctx).Raw(`
		SELECT 
			HOUR(created_at) AS hour,
			SUM(CASE WHEN release_status = 'success' THEN 1 ELSE 0 END) AS success_count,
			SUM(CASE WHEN release_status = 'failed' THEN 1 ELSE 0 END) AS failed_count,
			SUM(CASE WHEN target_release_id IS NOT NULL THEN 1 ELSE 0 END) AS rollback_count
		FROM valyria_dragon_release
		WHERE created_at >= ? AND created_at < ?
		GROUP BY HOUR(created_at)
	`, dayStart, dayEnd).Scan(&rows).Error
	return rows, err
}

// GetTrendByDay 按天聚合发布数（startDate 含当日 00:00，endDate 为最后一日 00:00，查询区间 [startDate, endDate+1d)）
// SELECT 与 GROUP BY 使用同一表达式 DATE(created_at)，满足 ONLY_FULL_GROUP_BY；Date 为 "2006-01-02" 格式
func (r *releaseRepository) GetTrendByDay(ctx context.Context, startDate, endDate time.Time) ([]TrendDayRow, error) {
	var rows []TrendDayRow
	endExclusive := endDate.Add(24 * time.Hour)
	err := r.db.WithContext(ctx).Raw(`
		SELECT 
			DATE(created_at) AS date,
			SUM(CASE WHEN release_status = 'success' THEN 1 ELSE 0 END) AS success_count,
			SUM(CASE WHEN release_status = 'failed' THEN 1 ELSE 0 END) AS failed_count,
			SUM(CASE WHEN target_release_id IS NOT NULL THEN 1 ELSE 0 END) AS rollback_count
		FROM valyria_dragon_release
		WHERE created_at >= ? AND created_at < ?
		GROUP BY DATE(created_at)
		ORDER BY DATE(created_at)
	`, startDate, endExclusive).Scan(&rows).Error
	return rows, err
}

// ListRecentByEnvName 按环境名称取「每个 pipeline 最近一次发布」再按时间排序取前 limit 条
func (r *releaseRepository) ListRecentByEnvName(ctx context.Context, envName string, limit int) ([]model.Release, error) {
	if limit <= 0 {
		limit = 10
	}
	// 子查询：每个 pipeline 在指定环境下最近一次发布的 id（MAX(id) 即最新）
	// 再按 created_at 降序取前 limit 条
	var idRows []struct{ ID uint }
	err := r.db.WithContext(ctx).Raw(`
		SELECT valyria_dragon_release.id FROM valyria_dragon_release
		INNER JOIN (
			SELECT r2.pipeline_id, MAX(r2.id) AS max_id
			FROM valyria_dragon_release r2
			JOIN valyria_dragon_pipeline p2 ON p2.id = r2.pipeline_id
			JOIN valyria_dragon_environment e2 ON e2.id = p2.environment_id
			WHERE e2.name = ?
			GROUP BY r2.pipeline_id
		) latest ON latest.pipeline_id = valyria_dragon_release.pipeline_id AND latest.max_id = valyria_dragon_release.id
		JOIN valyria_dragon_pipeline p ON p.id = valyria_dragon_release.pipeline_id
		JOIN valyria_dragon_environment e ON e.id = p.environment_id
		WHERE e.name = ?
		ORDER BY valyria_dragon_release.created_at DESC
		LIMIT ?
	`, envName, envName, limit).Scan(&idRows).Error
	if err != nil {
		return nil, err
	}
	ids := make([]uint, 0, len(idRows))
	for _, row := range idRows {
		ids = append(ids, row.ID)
	}
	if len(ids) == 0 {
		return []model.Release{}, nil
	}
	var data []model.Release
	err = r.db.WithContext(ctx).
		Where("id IN ?", ids).
		Preload("Pipeline.Environment.Project").
		Order("created_at DESC").
		Find(&data).Error
	return data, err
}

// ListProjectPipelineStats 按项目统计：发布次数、pipeline 数、成功/失败/回滚次数
func (r *releaseRepository) ListProjectPipelineStats(ctx context.Context) ([]ProjectPipelineStatsRow, error) {
	var rows []ProjectPipelineStatsRow
	err := r.db.WithContext(ctx).Raw(`
		SELECT 
			pr.id AS project_id,
			pr.name AS project_name,
			COUNT(DISTINCT p.id) AS pipeline_count,
			COUNT(r.id) AS release_count,
			COALESCE(SUM(CASE WHEN r.release_status = 'success' THEN 1 ELSE 0 END), 0) AS success_count,
			COALESCE(SUM(CASE WHEN r.release_status = 'failed' THEN 1 ELSE 0 END), 0) AS failed_count,
			COALESCE(SUM(CASE WHEN r.target_release_id IS NOT NULL THEN 1 ELSE 0 END), 0) AS rollback_count
		FROM valyria_dragon_project pr
		LEFT JOIN valyria_dragon_environment e ON e.project_id = pr.id
		LEFT JOIN valyria_dragon_pipeline p ON p.environment_id = e.id
		LEFT JOIN valyria_dragon_release r ON r.pipeline_id = p.id
		GROUP BY pr.id, pr.name
		ORDER BY pr.id
	`).Scan(&rows).Error
	return rows, err
}

// WithTx 返回一个绑定事务的 Repository
func (r *releaseRepository) WithTx(db *gorm.DB) ReleaseRepository {
	return &releaseRepository{db: db}
}
