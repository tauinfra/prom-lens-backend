package service

import (
	"context"
	"fmt"
	"time"
	"valyria-backend/internal/apps/foundry/executor"
	"valyria-backend/internal/apps/foundry/model"
	"valyria-backend/internal/apps/foundry/repository"
	pg "valyria-backend/internal/core/pagination"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ReleaseManager interface {
	List(ctx context.Context, params pg.QueryParams) ([]model.Release, pg.Pagination, error)
	Get(ctx context.Context, id int) (model.Release, error)
	Create(ctx context.Context, data *model.Release) error
	Update(ctx context.Context, id int, data *model.Release) error
	Delete(ctx context.Context, id int) error
}

// releaseManager 实现 ReleaseManager 接口
type releaseManager struct {
	repo     repository.ReleaseRepository
	pipeline repository.PipelineRepository
	gitlab   GitlabManager
	db       *gorm.DB
}

type Image struct {
	Registry string
	Name     string
	Tag      string
}

type GitRefInfo struct {
	RefType  string
	RefName  string
	CommitID string
}

func GenerateImage(p *model.Pipeline, ref, commit string) Image {
	ts := time.Now().Format("20060102150405")
	return Image{
		Registry: fmt.Sprintf(
			"%s/%s",
			p.Environment.CredHarbor.Server,
			p.Environment.Project.Name,
		),
		Name: p.Application.Name,
		Tag:  fmt.Sprintf("%s-%s-%s", ts, ref, commit),
	}
}

func GenerateTaskID(p *model.Pipeline) string {
	return fmt.Sprintf(
		"%s-%s-%s-%s",
		p.Environment.Project.Name,
		p.Environment.Name,
		p.Application.Name,
		uuid.NewString()[:8],
	)
}

// NewReleaseManager 创建新的 ReleaseManager 实例
func NewReleaseManager(repo repository.ReleaseRepository, pipeline repository.PipelineRepository, gitlab GitlabManager, db *gorm.DB) ReleaseManager {
	// 使用 NewGitlabCredentialRepository 函数创建具体实现
	return &releaseManager{repo: repo, pipeline: pipeline, gitlab: gitlab, db: db}
}

// List 列表
func (s *releaseManager) List(ctx context.Context, params pg.QueryParams) ([]model.Release, pg.Pagination, error) {
	return s.repo.List(ctx, params)
}

// Get 查询
func (s *releaseManager) Get(ctx context.Context, id int) (model.Release, error) {
	return s.repo.Get(ctx, id)
}

// Create 创建
func (s *releaseManager) Create(ctx context.Context, data *model.Release) error {
	pipeline, err := s.pipeline.Get(ctx, data.PipelineID)
	if err != nil {
		return err
	}
	gitlab, err := s.gitlab.GetRef(ctx, pipeline.CredGitlabID, pipeline.GitlabProjectID, data.GitRef)
	if err != nil {
		return err
	}

	image := GenerateImage(&pipeline, gitlab.RefName, gitlab.CommitID)
	taskID := GenerateTaskID(&pipeline)
	data.TaskID = taskID
	data.ImageRegistry = image.Registry
	data.ImageName = image.Name
	data.ImageTag = image.Tag
	data.GitRefType = gitlab.RefType
	data.GitRef = gitlab.RefName
	data.GitCommit = gitlab.CommitID

	if err = s.repo.Create(ctx, data); err != nil {
		return err
	}
	// 执行发布
	if _, err = executor.Run(taskID, s.db); err != nil {
		return err
	}
	return nil
}

// Update 更新
func (s *releaseManager) Update(ctx context.Context, id int, data *model.Release) error {
	return s.repo.Update(ctx, id, data)
}

// Delete 更新
func (s *releaseManager) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
