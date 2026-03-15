package kubernetes

import (
	"valyria-backend/internal/apps/kubernetes/controller"
	"valyria-backend/internal/apps/kubernetes/repository"
	"valyria-backend/internal/apps/kubernetes/service"
	"valyria-backend/internal/pkg/encryption"

	"gorm.io/gorm"
)

type ClusterProvider struct {
	Repo       *repository.ClusterRepository
	Service    *service.ClusterManager
	Controller *controller.ClusterController
}

func NewClusterProvider(db *gorm.DB, encryptor encryption.Encryptor) *ClusterProvider {
	// 初始化用户模块的依赖
	repo := repository.NewClusterRepository(db, encryptor)
	svc := service.NewClusterManager(db, repo)
	ctrl := controller.NewClusterController(svc)

	return &ClusterProvider{
		Repo:       &repo,
		Service:    &svc,
		Controller: ctrl,
	}
}
