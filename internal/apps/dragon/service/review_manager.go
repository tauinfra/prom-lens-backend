package service

import (
	"context"
	"errors"
	authnmodel "valyria-backend/internal/apps/authn/model"
	"valyria-backend/internal/apps/dragon/dto"
	"valyria-backend/internal/apps/dragon/executor"
	"valyria-backend/internal/apps/dragon/model"
	"valyria-backend/internal/apps/dragon/repository"
	"valyria-backend/internal/apps/dragon/request"
	"valyria-backend/internal/core/config"
	pg "valyria-backend/internal/core/pagination"

	"gorm.io/gorm"
)

type ReviewManager interface {
	List(ctx context.Context, params pg.QueryParams) ([]dto.ReviewDTO, pg.Pagination, error)
	Get(ctx context.Context, id uint) (dto.ReviewDTO, error)
	Create(ctx context.Context, req *request.CreateReviewRequest, applicantID uint) error
	Update(ctx context.Context, id uint, reviewerID uint, req *request.UpdateReviewRequest) error
	Delete(ctx context.Context, id uint) error
}

// ReviewManager 实现 ReviewManager 接口
type reviewManager struct {
	repo          repository.ReviewRepository
	configRepo    repository.ReviewConfigRepository
	aclRepo       repository.PipelineACLRepository
	gitlabManager GitlabManager
	db            *gorm.DB
	cfg           *config.Config
}

// NewReviewManager 创建新的 ReviewManager 实例
func NewReviewManager(repo repository.ReviewRepository, configRepo repository.ReviewConfigRepository, aclRepo repository.PipelineACLRepository, gitlabManager GitlabManager, db *gorm.DB, cfg *config.Config) ReviewManager {
	// 使用 NewGitlabCredentialRepository 函数创建具体实现
	return &reviewManager{repo: repo, configRepo: configRepo, aclRepo: aclRepo, gitlabManager: gitlabManager, db: db, cfg: cfg}
}

// List 列表
func (s *reviewManager) List(ctx context.Context, params pg.QueryParams) ([]dto.ReviewDTO, pg.Pagination, error) {
	list, pagination, err := s.repo.List(ctx, params)
	if err != nil {
		return nil, pg.Pagination{}, err
	}
	reviewerIDs := make([]uint, 0, len(list))
	for _, r := range list {
		if r.ReviewerID != 0 {
			reviewerIDs = append(reviewerIDs, r.ReviewerID)
		}
	}
	names := s.getReviewerNamesByIDs(ctx, s.db, reviewerIDs)
	out := make([]dto.ReviewDTO, 0, len(list))
	for _, r := range list {
		out = append(out, reviewToDTO(r, names[r.ReviewerID]))
	}
	return out, pagination, nil
}

// Get 查询
func (s *reviewManager) Get(ctx context.Context, id uint) (dto.ReviewDTO, error) {
	r, err := s.repo.Get(ctx, id)
	if err != nil {
		return dto.ReviewDTO{}, err
	}
	var name string
	if r.ReviewerID != 0 {
		names := s.getReviewerNamesByIDs(ctx, s.db, []uint{r.ReviewerID})
		name = names[r.ReviewerID]
	}
	return reviewToDTO(r, name), nil
}

// Create 创建
func (s *reviewManager) Create(ctx context.Context, req *request.CreateReviewRequest, applicantID uint) error {
	data := &model.Review{
		ReleaseID:   req.ReleaseID,
		ApplicantID: applicantID,
	}
	return s.db.Transaction(func(tx *gorm.DB) error {
		if data.ReleaseID == 0 {
			return errors.New("releaseID is required")
		}
		release, err := s.getReleaseContext(ctx, tx, data.ReleaseID)
		if err != nil {
			return err
		}
		if data.ApplicantID == 0 {
			releaseCreatorID, err := s.getUserIDByUsername(ctx, tx, release.Creator)
			if err != nil {
				return err
			}
			data.ApplicantID = releaseCreatorID
		}
		stages, err := s.resolveReviewStages(ctx, tx, release)
		if err != nil {
			return err
		}
		if len(stages) == 0 {
			return errors.New("review stages not configured")
		}
		stepExists, err := s.repo.WithTx(tx).ExistsStepInRelease(ctx, release.ID, 1)
		if err != nil {
			return err
		}
		if stepExists {
			return errors.New("review step already exists")
		}
		data.Step = 1
		data.MaxStep = len(stages)
		data.ApprovalRole = stages[0].Role.Name
		data.Status = "pending"
		if err = s.repo.WithTx(tx).Create(ctx, data); err != nil {
			return err
		}
		if release.ApprovalStatus == "none" {
			if err = tx.WithContext(ctx).Model(&model.Release{}).
				Where("id = ?", release.ID).
				Update("approval_status", "pending").Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *reviewManager) Update(ctx context.Context, id uint, reviewerID uint, req *request.UpdateReviewRequest) error {
	var (
		shouldRun bool
		taskID    string
		repoURL   string
	)
	err := s.db.Transaction(func(tx *gorm.DB) error {
		current, err := s.repo.WithTx(tx).Get(ctx, id)
		if err != nil {
			return err
		}
		if current.Status != "pending" {
			return errors.New("review already processed")
		}
		if req.Status != "approved" && req.Status != "rejected" {
			return errors.New("invalid status")
		}
		if reviewerID == 0 {
			return errors.New("reviewer is required")
		}
		release, err := s.getReleaseContext(ctx, tx, current.ReleaseID)
		if err != nil {
			return err
		}
		stages, err := s.resolveReviewStages(ctx, tx, release)
		if err != nil {
			return err
		}
		if len(stages) == 0 {
			return errors.New("review stages not configured")
		}
		if current.Step <= 0 || current.Step > len(stages) {
			return errors.New("invalid review step")
		}
		expectedStage := stages[current.Step-1]
		if expectedStage.Role == nil || expectedStage.Role.Name == "" {
			return errors.New("approval role not configured")
		}
		if current.ApprovalRole != "" && current.ApprovalRole != expectedStage.Role.Name {
			// 兼容角色名称迁移/变更，自动同步当前审批角色名称
			if err := tx.WithContext(ctx).Model(&model.Review{}).
				Where("id = ?", current.ID).
				Update("approval_role", expectedStage.Role.Name).Error; err != nil {
				return err
			}
			current.ApprovalRole = expectedStage.Role.Name
		}
		if !IsSuperuser(ctx) {
			hasRole, err := s.configRepo.WithTx(tx).UserHasRole(ctx, reviewerID, expectedStage.RoleID)
			if err != nil {
				return err
			}
			if !hasRole {
				return errors.New("reviewer has no permission")
			}
			ok, err := s.aclRepo.HasUserAction(ctx, reviewerID, release.PipelineID, ActionApprove)
			if err != nil {
				return err
			}
			if !ok {
				return errors.New("reviewer has no permission")
			}
		}
		reviewerExists, err := s.repo.WithTx(tx).ExistsReviewerInRelease(ctx, current.ReleaseID, reviewerID, current.ID)
		if err != nil {
			return err
		}
		if reviewerExists {
			return errors.New("reviewer already approved another step")
		}

		// 1. 更新当前审批
		update := &model.Review{
			Status:     req.Status,
			ReviewerID: reviewerID,
			Comment:    req.Comment,
		}
		if err = s.repo.WithTx(tx).Update(ctx, id, update); err != nil {
			return err
		}
		// 2. 如果审批通过，生成下一步 Review
		if req.Status == "approved" {
			// 判断是否还有下一步
			if current.Step < len(stages) {
				stepExists, err := s.repo.WithTx(tx).ExistsStepInRelease(ctx, current.ReleaseID, current.Step+1)
				if err != nil {
					return err
				}
				if stepExists {
					return errors.New("review step already exists")
				}
				nextStage := stages[current.Step]
				if nextStage.Role == nil || nextStage.Role.Name == "" {
					return errors.New("approval role not configured")
				}
				nextReview := &model.Review{
					ReleaseID:    current.ReleaseID,
					ApplicantID:  current.ApplicantID,
					Status:       "pending",
					ApprovalRole: nextStage.Role.Name,
					Step:         current.Step + 1,
					MaxStep:      len(stages),
				}
				if err := s.repo.WithTx(tx).Create(ctx, nextReview); err != nil {
					return err
				}
			} else {
				// 最后一步审批完成，更新 Release 状态
				if err := tx.WithContext(ctx).Model(&model.Release{}).
					Where("id = ?", current.ReleaseID).
					Update("approval_status", "approved").Error; err != nil {
					return err
				}
				repoURL, err = s.gitlabManager.GetRepoURL(ctx, int64(release.Pipeline.GitlabID), int64(release.Pipeline.GitlabProjectID))
				if err != nil {
					return err
				}
				taskID = release.TaskID
				shouldRun = true

			}
		} else if req.Status == "rejected" {
			if err := tx.WithContext(ctx).Model(&model.Release{}).
				Where("id = ?", current.ReleaseID).
				Update("approval_status", "rejected").Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	if shouldRun {
		_, err = executor.Run(ctx, taskID, repoURL, s.db, s.cfg)
		return err
	}
	return nil
}

func (s *reviewManager) getReleaseContext(ctx context.Context, tx *gorm.DB, releaseID uint) (model.Release, error) {
	var release model.Release
	if err := tx.WithContext(ctx).
		Preload("Pipeline.Environment.Project").
		First(&release, releaseID).Error; err != nil {
		return release, err
	}
	return release, nil
}

func (s *reviewManager) getUserIDByUsername(ctx context.Context, tx *gorm.DB, username string) (uint, error) {
	if username == "" {
		return 0, gorm.ErrInvalidData
	}
	var user authnmodel.User
	if err := tx.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		return 0, err
	}
	return uint(user.ID), nil
}

// getReviewerNamesByIDs 根据用户 ID 批量查询审批人显示名（优先 nickname，否则 username）
func (s *reviewManager) getReviewerNamesByIDs(ctx context.Context, db *gorm.DB, ids []uint) map[uint]string {
	if len(ids) == 0 {
		return nil
	}
	var users []authnmodel.User
	if err := db.WithContext(ctx).Where("id IN ?", ids).Find(&users).Error; err != nil {
		return nil
	}
	out := make(map[uint]string, len(users))
	for _, u := range users {
		name := u.Nickname
		if name == "" {
			name = u.Username
		}
		out[uint(u.ID)] = name
	}
	return out
}

func reviewToDTO(r model.Review, reviewerName string) dto.ReviewDTO {
	return dto.ReviewDTO{
		ID:           r.ID,
		ApplicantID:  r.ApplicantID,
		ReviewerID:   r.ReviewerID,
		ReviewerName: reviewerName,
		ReleaseID:    r.ReleaseID,
		Status:       r.Status,
		ApprovalRole: r.ApprovalRole,
		Step:         r.Step,
		MaxStep:      r.MaxStep,
		Comment:      r.Comment,
		CreatedAt:    r.CreatedAt,
		UpdatedAt:    r.UpdatedAt,
	}
}

func (s *reviewManager) resolveReviewStages(ctx context.Context, tx *gorm.DB, release model.Release) ([]model.ReviewStageConfig, error) {
	configRepo := s.configRepo.WithTx(tx)
	return configRepo.ListAll(ctx)
}

// Delete 更新
func (s *reviewManager) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}
