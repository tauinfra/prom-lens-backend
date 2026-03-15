package service

import (
	"context"
	"time"
	"valyria-backend/internal/apps/dragon/dto"
	"valyria-backend/internal/apps/dragon/repository"
	"valyria-backend/internal/apps/dragon/request"
)

// ReleaseReportManager 发布报表
type ReleaseReportManager interface {
	Summary(ctx context.Context, req *request.ReleaseReportRequest) (dto.ReleaseReportResponse, error)
}

type releaseReportManager struct {
	repo repository.ReleaseReportRepository
}

func NewReleaseReportManager(repo repository.ReleaseReportRepository) ReleaseReportManager {
	return &releaseReportManager{repo: repo}
}

func (s *releaseReportManager) Summary(ctx context.Context, req *request.ReleaseReportRequest) (dto.ReleaseReportResponse, error) {
	from := req.From
	to := req.To
	if from == "" || to == "" {
		now := time.Now()
		to = now.Format("2006-01")
		from = now.AddDate(0, -11, 0).Format("2006-01") // 最近 12 个月
	}
	if from > to {
		from, to = to, from
	}
	summary, err := s.repo.GetSummary(ctx, from, to)
	if err != nil {
		return dto.ReleaseReportResponse{}, err
	}
	byMonth, err := s.repo.GetByMonth(ctx, from, to)
	if err != nil {
		return dto.ReleaseReportResponse{}, err
	}
	byProject, err := s.repo.GetByProject(ctx, from, to)
	if err != nil {
		return dto.ReleaseReportResponse{}, err
	}
	return dto.ReleaseReportResponse{
		Summary:      summary,
		ByMonth:      byMonth,
		ByProject:    byProject,
		ProjectCount: len(byProject),
	}, nil
}
