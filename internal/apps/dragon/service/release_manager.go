package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
	"valyria-backend/internal/apps/dragon/dto"
	"valyria-backend/internal/apps/dragon/executor"
	"valyria-backend/internal/apps/dragon/model"
	"valyria-backend/internal/apps/dragon/repository"
	"valyria-backend/internal/apps/dragon/request"
	"valyria-backend/internal/core/config"
	dashdto "valyria-backend/internal/apps/dashboard/dto"
	pg "valyria-backend/internal/core/pagination"

	"gorm.io/gorm"
)

type ReleaseManager interface {
	List(ctx context.Context, projectID, environmentID, pipelineID uint, params pg.QueryParams) ([]dto.ReleaseDTO, pg.Pagination, error)
	Get(ctx context.Context, pipelineID, releaseID uint) (dto.ReleaseDTO, error)
	Create(ctx context.Context, req *request.CreateReleaseRequest) (dto.ReleaseDTO, error)
	Rollback(ctx context.Context, projectID, environmentID, pipelineID uint, req *request.RollbackRequest) error
	Delete(ctx context.Context, pipelineID, releaseID uint) error
	// GetTodayStats 今日 00:00–23:59 发布统计（不区分环境、项目）
	GetTodayStats(ctx context.Context) (todayTotal, success, failed, rollback int, err error)
	// GetPipelineTrend 流水线趋势，range=today|7d|30d
	GetPipelineTrend(ctx context.Context, rangeParam string) (*dashdto.PipelineTrendResponse, error)
	// GetRecentReleases 按环境名称取最近发布记录
	GetRecentReleases(ctx context.Context, envName string, limit int) ([]dashdto.PipelineRecentItem, error)
	// GetProjectPipelineStats 按项目统计发布次数、pipeline 数、成功/失败/回滚次数
	GetProjectPipelineStats(ctx context.Context) ([]dashdto.PipelineProjectStatsItem, error)
}

// releaseManager 实现 ReleaseManager 接口
type releaseManager struct {
	repo          repository.ReleaseRepository
	pipelineRepo  repository.PipelineRepository
	aclRepo       repository.PipelineACLRepository
	reviewManager ReviewManager
	gitlabManager GitlabManager
	db            *gorm.DB
	cfg           *config.Config
}

// NewReleaseManager 创建新的 ReleaseManager 实例
func NewReleaseManager(
	repo repository.ReleaseRepository,
	pipelineRepo repository.PipelineRepository,
	aclRepo repository.PipelineACLRepository,
	reviewManager ReviewManager,
	gitlabManager GitlabManager,
	db *gorm.DB,
	cfg *config.Config,
) ReleaseManager {
	return &releaseManager{repo: repo, pipelineRepo: pipelineRepo, aclRepo: aclRepo, reviewManager: reviewManager, gitlabManager: gitlabManager, db: db, cfg: cfg}
}

// List 列表
func (s *releaseManager) List(ctx context.Context, projectID, environmentID, pipelineID uint, params pg.QueryParams) (data []dto.ReleaseDTO, pagination pg.Pagination, err error) {
	if !IsSuperuser(ctx) {
		userID := GetUID(ctx)
		if userID == 0 {
			return data, pagination, errors.New("no permission")
		}
		ok, err := s.hasViewPermission(ctx, userID, pipelineID)
		if err != nil {
			return data, pagination, err
		}
		if !ok {
			return data, pagination, errors.New("no permission")
		}
	}
	releases, pagination, err := s.repo.List(ctx, projectID, environmentID, pipelineID, params)
	if err != nil {
		return data, pagination, err
	}
	for _, release := range releases {
		data = append(data, dto.ReleaseDTO{
			ID:              release.ID,
			TaskID:          release.TaskID,
			ProjectID:       release.Pipeline.Environment.Project.ID,
			ProjectName:     release.Pipeline.Environment.Project.Name,
			EnvironmentID:   release.Pipeline.Environment.ID,
			EnvironmentName: release.Pipeline.Environment.Name,
			PipelineID:      release.Pipeline.ID,
			PipelineName:    release.Pipeline.Name,
			GitRef:          release.GitRef,
			GitCommit:       release.GitCommit,
			GitRefType:      release.GitRefType,
			ImageRegistry:   release.ImageRegistry,
			ImageName:       release.ImageName,
			ImageTag:        release.ImageTag,
			ReleaseStatus:   release.ReleaseStatus,
			ApprovalStatus:  release.ApprovalStatus,
			Description:     release.Description,
			TargetReleaseID: release.TargetReleaseID,
			Creator:         release.Creator,
			StartedAt:       release.StartedAt,
			FinishedAt:      release.FinishedAt,
			CreatedAt:       release.CreatedAt,
			UpdatedAt:       release.UpdatedAt,
		})
	}
	return data, pagination, nil
}

// Get 查询
func (s *releaseManager) Get(ctx context.Context, pipelineID, releaseID uint) (data dto.ReleaseDTO, err error) {
	if !IsSuperuser(ctx) {
		userID := GetUID(ctx)
		if userID == 0 {
			return data, errors.New("no permission")
		}
		ok, err := s.hasViewPermission(ctx, userID, pipelineID)
		if err != nil {
			return data, err
		}
		if !ok {
			return data, errors.New("no permission")
		}
	}
	release, err := s.repo.Get(ctx, pipelineID, releaseID)
	if err != nil {
		return data, err
	}
	data = dto.ReleaseDTO{
		ID:              release.ID,
		TaskID:          release.TaskID,
		ProjectID:       release.Pipeline.Environment.Project.ID,
		ProjectName:     release.Pipeline.Environment.Project.Name,
		EnvironmentID:   release.Pipeline.Environment.ID,
		EnvironmentName: release.Pipeline.Environment.Name,
		PipelineID:      release.Pipeline.ID,
		PipelineName:    release.Pipeline.Name,
		GitRef:          release.GitRef,
		GitCommit:       release.GitCommit,
		GitRefType:      release.GitRefType,
		ImageRegistry:   release.ImageRegistry,
		ImageName:       release.ImageName,
		ImageTag:        release.ImageTag,
		ReleaseStatus:   release.ReleaseStatus,
		ApprovalStatus:  release.ApprovalStatus,
		Description:     release.Description,
		TargetReleaseID: release.TargetReleaseID,
		Creator:         release.Creator,
		StartedAt:       release.StartedAt,
		FinishedAt:      release.FinishedAt,
		CreatedAt:       release.CreatedAt,
		UpdatedAt:       release.UpdatedAt,
	}
	return data, nil
}

// Create 创建，返回当前创建的 release 信息（含 id），供前端查看进度
func (s *releaseManager) Create(ctx context.Context, req *request.CreateReleaseRequest) (dto.ReleaseDTO, error) {
	var zero dto.ReleaseDTO
	if !IsSuperuser(ctx) {
		userID := GetUID(ctx)
		if userID == 0 {
			return zero, errors.New("no permission")
		}
		ok, err := s.aclRepo.HasUserAction(ctx, userID, req.PipelineID, ActionDeploy)
		if err != nil {
			return zero, err
		}
		if !ok {
			return zero, errors.New("no permission")
		}
	}
	// 获取流水线
	pipeline, err := s.pipelineRepo.Get(ctx, req.ProjectID, req.EnvironmentID, req.PipelineID)
	if err != nil {
		return zero, err
	}
	// 同一流水线存在未结束的发布（审批中或发布中）时禁止再次创建
	exists, err := s.repo.ExistsUnfinishedByPipelineID(ctx, req.PipelineID)
	if err != nil {
		return zero, err
	}
	if exists {
		return zero, errors.New("该流水线当前存在进行中的发布或审批，请等待结束后再创建")
	}
	// 获取 Gitlab 分支信息
	git, err := s.gitlabManager.GetRef(ctx, int64(pipeline.GitlabID), int64(pipeline.GitlabProjectID), req.GitRef)
	if err != nil {
		return zero, err
	}
	// Gitlab 仓库地址
	repoURL, err := s.gitlabManager.GetRepoURL(ctx, int64(pipeline.GitlabID), int64(pipeline.GitlabProjectID))
	if err != nil {
		return zero, err
	}
	// 发布记录: 前端只接收分支信息（其他信息可能伪造，需要后端校验判断和生成）
	image := GenerateImage(&pipeline, git.RefName, git.CommitID)
	taskID := GenerateTaskID(&pipeline)
	// 构建 Model
	data := &model.Release{
		TaskID:        taskID,
		PipelineID:    req.PipelineID,
		ImageRegistry: image.Registry,
		ImageName:     image.Name,
		ImageTag:      image.Tag,
		GitRefType:    git.RefType,
		GitRef:        git.RefName,
		GitCommit:     git.CommitID,
		Creator:       req.Creator,
		Description:   req.Description,
	}
	if err = s.repo.Create(ctx, data); err != nil {
		return zero, err
	}
	if pipeline.IsApproval {
		applicantID := GetUID(ctx)
		if err = s.reviewManager.Create(ctx, &request.CreateReviewRequest{ReleaseID: data.ID}, applicantID); err != nil {
			_ = s.repo.Delete(ctx, req.PipelineID, data.ID)
			return zero, err
		}
	} else {
		// 执行构建发布（失败也返回已创建的 release，前端可据 id 查进度/状态）
		_, _ = executor.Run(ctx, taskID, repoURL, s.db, s.cfg)
	}
	// 返回当前创建的 release 完整信息，前端可用 data.id 查看进度
	return s.Get(ctx, req.PipelineID, data.ID)
}

// Rollback 回滚到指定 release（跳过审批，直接执行部署）
func (s *releaseManager) Rollback(ctx context.Context, projectID, environmentID, pipelineID uint, req *request.RollbackRequest) error {
	if !IsSuperuser(ctx) {
		userID := GetUID(ctx)
		if userID == 0 {
			return errors.New("no permission")
		}
		ok, err := s.aclRepo.HasUserAction(ctx, userID, pipelineID, ActionRollback)
		if err != nil {
			return err
		}
		if !ok {
			return errors.New("no permission")
		}
	}
	// 获取目标 release，且必须为同一 pipeline、发布状态为成功
	target, err := s.repo.Get(ctx, pipelineID, req.TargetReleaseID)
	if err != nil {
		return err
	}
	if target.ReleaseStatus != "success" {
		return errors.New("只能回滚到发布成功的记录")
	}
	// 同一流水线存在未结束的发布时禁止回滚
	exists, err := s.repo.ExistsUnfinishedByPipelineID(ctx, pipelineID)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("该流水线当前存在进行中的发布或审批，请等待结束后再回滚")
	}
	pipeline, err := s.pipelineRepo.Get(ctx, projectID, environmentID, pipelineID)
	if err != nil {
		return err
	}
	repoURL, err := s.gitlabManager.GetRepoURL(ctx, int64(pipeline.GitlabID), int64(pipeline.GitlabProjectID))
	if err != nil {
		return err
	}
	taskID := GenerateTaskID(&pipeline)
	targetID := req.TargetReleaseID
	data := &model.Release{
		TaskID:         taskID,
		PipelineID:     pipelineID,
		ImageRegistry:  target.ImageRegistry,
		ImageName:      target.ImageName,
		ImageTag:       target.ImageTag,
		GitRefType:     target.GitRefType,
		GitRef:         target.GitRef,
		GitCommit:      target.GitCommit,
		TargetReleaseID: &targetID,
		ApprovalStatus: "none",
		Creator:        GetUsername(ctx),
		Description:    "回滚至 release #" + fmt.Sprintf("%d", target.ID),
	}
	if err = s.repo.Create(ctx, data); err != nil {
		return err
	}
	_, err = executor.Run(ctx, taskID, repoURL, s.db, s.cfg)
	return err
}

// Delete 更新
func (s *releaseManager) Delete(ctx context.Context, pipelineID, releaseID uint) error {
	return s.repo.Delete(ctx, pipelineID, releaseID)
}

// GetTodayStats 今日 00:00–23:59 发布统计（不区分环境、项目）
func (s *releaseManager) GetTodayStats(ctx context.Context) (todayTotal, success, failed, rollback int, err error) {
	stats, err := s.repo.GetTodayStats(ctx)
	if err != nil {
		return 0, 0, 0, 0, err
	}
	return int(stats.TodayTotal), int(stats.Success), int(stats.Failed), int(stats.Rollback), nil
}

// GetPipelineTrend 流水线趋势，rangeParam=today|7d|30d
func (s *releaseManager) GetPipelineTrend(ctx context.Context, rangeParam string) (*dashdto.PipelineTrendResponse, error) {
	now := time.Now()
	loc := now.Location()

	switch rangeParam {
	case "today":
		dayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
		dayEnd := dayStart.Add(24 * time.Hour)
		rows, err := s.repo.GetTrendByHour(ctx, dayStart, dayEnd)
		if err != nil {
			return nil, err
		}
		byHour := make(map[int]struct{ Success, Failed, Rollback int64 })
		for _, r := range rows {
			byHour[r.Hour] = struct{ Success, Failed, Rollback int64 }{r.SuccessCount, r.FailedCount, r.RollbackCount}
		}
		points := make([]dashdto.PipelineTrendPoint, 0, 24)
		for h := 0; h < 24; h++ {
			c := byHour[h]
			points = append(points, dashdto.PipelineTrendPoint{
				Hour:          fmt.Sprintf("%02d:00", h),
				SuccessCount:  int(c.Success),
				FailedCount:   int(c.Failed),
				RollbackCount: int(c.Rollback),
			})
		}
		return &dashdto.PipelineTrendResponse{Range: "today", Points: points}, nil

	case "7d":
		endDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
		startDate := endDate.AddDate(0, 0, -6)
		rows, err := s.repo.GetTrendByDay(ctx, startDate, endDate)
		if err != nil {
			return nil, err
		}
		byDate := make(map[string]struct{ Success, Failed, Rollback int64 })
		for _, r := range rows {
			key := dateKey(r.Date)
			byDate[key] = struct{ Success, Failed, Rollback int64 }{r.SuccessCount, r.FailedCount, r.RollbackCount}
		}
		points := make([]dashdto.PipelineTrendPoint, 0, 7)
		for d := 0; d < 7; d++ {
			dayLocal := startDate.AddDate(0, 0, d)
			// 本地当日 00:00 与次日 00:00 转 UTC，可能跨两个 UTC 日期，需合并
			dayStartUTC := time.Date(dayLocal.Year(), dayLocal.Month(), dayLocal.Day(), 0, 0, 0, 0, loc).UTC()
			dayEndUTC := dayStartUTC.Add(24 * time.Hour)
			key1 := dayStartUTC.Format("2006-01-02")
			key2 := dayEndUTC.Add(-time.Second).Format("2006-01-02")
			var c struct{ Success, Failed, Rollback int64 }
			if v, ok := byDate[key1]; ok {
				c.Success += v.Success
				c.Failed += v.Failed
				c.Rollback += v.Rollback
			}
			if key2 != key1 {
				if v, ok := byDate[key2]; ok {
					c.Success += v.Success
					c.Failed += v.Failed
					c.Rollback += v.Rollback
				}
			}
			points = append(points, dashdto.PipelineTrendPoint{
				Date:          dayLocal.Format("1/2"),
				SuccessCount:  int(c.Success),
				FailedCount:   int(c.Failed),
				RollbackCount: int(c.Rollback),
			})
		}
		return &dashdto.PipelineTrendResponse{Range: "7d", Points: points}, nil

	case "30d":
		endDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
		startDate := endDate.AddDate(0, 0, -29)
		rows, err := s.repo.GetTrendByDay(ctx, startDate, endDate)
		if err != nil {
			return nil, err
		}
		byDate := make(map[string]struct{ Success, Failed, Rollback int64 })
		for _, r := range rows {
			key := dateKey(r.Date)
			byDate[key] = struct{ Success, Failed, Rollback int64 }{r.SuccessCount, r.FailedCount, r.RollbackCount}
		}
		points := make([]dashdto.PipelineTrendPoint, 0, 30)
		for d := 0; d < 30; d++ {
			dayLocal := startDate.AddDate(0, 0, d)
			dayStartUTC := time.Date(dayLocal.Year(), dayLocal.Month(), dayLocal.Day(), 0, 0, 0, 0, loc).UTC()
			dayEndUTC := dayStartUTC.Add(24 * time.Hour)
			key1 := dayStartUTC.Format("2006-01-02")
			key2 := dayEndUTC.Add(-time.Second).Format("2006-01-02")
			var c struct{ Success, Failed, Rollback int64 }
			if v, ok := byDate[key1]; ok {
				c.Success += v.Success
				c.Failed += v.Failed
				c.Rollback += v.Rollback
			}
			if key2 != key1 {
				if v, ok := byDate[key2]; ok {
					c.Success += v.Success
					c.Failed += v.Failed
					c.Rollback += v.Rollback
				}
			}
			points = append(points, dashdto.PipelineTrendPoint{
				Date:          dayLocal.Format("1/2"),
				SuccessCount:  int(c.Success),
				FailedCount:   int(c.Failed),
				RollbackCount: int(c.Rollback),
			})
		}
		return &dashdto.PipelineTrendResponse{Range: "30d", Points: points}, nil

		default:
		return nil, errors.New("invalid range: must be today, 7d or 30d")
	}
}

// GetRecentReleases 按环境名称取最近发布记录，返回 projectName / pipeline / env / status / duration(秒) / time
func (s *releaseManager) GetRecentReleases(ctx context.Context, envName string, limit int) ([]dashdto.PipelineRecentItem, error) {
	if envName == "" {
		return nil, errors.New("env is required")
	}
	if limit <= 0 {
		limit = 10
	}
	releases, err := s.repo.ListRecentByEnvName(ctx, envName, limit)
	if err != nil {
		return nil, err
	}
	out := make([]dashdto.PipelineRecentItem, 0, len(releases))
	for _, r := range releases {
		item := dashdto.PipelineRecentItem{
			Status: r.ReleaseStatus,
		}
		if r.Pipeline != nil {
			item.Pipeline = r.Pipeline.Name
			if r.Pipeline.Environment != nil {
				item.Env = r.Pipeline.Environment.Name
				if r.Pipeline.Environment.Project != nil {
					item.ProjectName = r.Pipeline.Environment.Project.Name
				}
			}
		}
		if r.StartedAt != nil && r.FinishedAt != nil {
			item.Duration = int64(r.FinishedAt.Sub(*r.StartedAt).Seconds())
		}
		if r.FinishedAt != nil {
			item.Time = r.FinishedAt.Format("2006-01-02 15:04:05")
		} else {
			item.Time = r.CreatedAt.Format("2006-01-02 15:04:05")
		}
		out = append(out, item)
	}
	return out, nil
}

// GetProjectPipelineStats 按项目统计发布次数、pipeline 数、成功/失败/回滚次数
func (s *releaseManager) GetProjectPipelineStats(ctx context.Context) ([]dashdto.PipelineProjectStatsItem, error) {
	rows, err := s.repo.ListProjectPipelineStats(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]dashdto.PipelineProjectStatsItem, 0, len(rows))
	for _, r := range rows {
		out = append(out, dashdto.PipelineProjectStatsItem{
			ProjectName:   r.ProjectName,
			PipelineCount: int(r.PipelineCount),
			ReleaseCount:  int(r.ReleaseCount),
			SuccessCount:  int(r.SuccessCount),
			FailedCount:   int(r.FailedCount),
			RollbackCount: int(r.RollbackCount),
		})
	}
	return out, nil
}

// dateKey 规范化 MySQL 返回的日期为 "2006-01-02"（可能带 " 00:00:00" 或 "T00:00:00Z"）
func dateKey(s string) string {
	s = strings.TrimSpace(s)
	if len(s) >= 10 {
		return s[:10]
	}
	return s
}

func (s *releaseManager) hasViewPermission(ctx context.Context, userID, pipelineID uint) (bool, error) {
	actions := []string{ActionView, ActionDeploy, ActionApprove, ActionRollback}
	for _, action := range actions {
		ok, err := s.aclRepo.HasUserAction(ctx, userID, pipelineID, action)
		if err != nil {
			return false, err
		}
		if ok {
			return true, nil
		}
	}
	return false, nil
}
