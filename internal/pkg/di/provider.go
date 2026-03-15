package di

import (
	"valyria-backend/internal/core/config"
	coreMiddleware "valyria-backend/internal/core/middleware"
	"valyria-backend/internal/pkg/di/audit"
	"valyria-backend/internal/pkg/di/authn"
	"valyria-backend/internal/pkg/di/dashboard"
	"valyria-backend/internal/pkg/di/dragon"
	"valyria-backend/internal/pkg/di/kubernetes"
	"valyria-backend/internal/pkg/di/prometheus"
	"valyria-backend/internal/pkg/encryption"

	"gorm.io/gorm"
)

type Provider struct {
	JWTManager        *coreMiddleware.JWTManager
	AuditManager      *coreMiddleware.AuditManager
	PermissionManager *coreMiddleware.PermissionManager
	Audit             *audit.Provider
	Authn             *authn.AuthnProvider
	Kubernetes        *kubernetes.KubernetesProvider
	Prometheus        *prometheus.Provider
	Dragon            *dragon.Provider
	Dashboard         *dashboard.Provider
}

func NewProvider(cfg *config.Config, db *gorm.DB, encryptor encryption.Encryptor) *Provider {
	// 创建 JWT 管理器
	jwtManager := coreMiddleware.NewJWTManager(cfg)
	authnProvider := authn.NewAuthnProvider(db, jwtManager)
	auditManager := coreMiddleware.NewAuditManager(db)
	// 权限中间件
	permissionManager := coreMiddleware.NewPermissionManager(cfg, *authnProvider.Auth.Service)
	// Kubernetes Provider
	kubernetesProvider := kubernetes.NewKubernetesProvider(db, encryptor, cfg)
	// 发布
	dragonProvider := dragon.NewDragonProvider(db, encryptor, cfg)
	// 仪表盘（独立模块，依赖 K8s 集群统计 + Dragon 今日发布统计）
	dashboardProvider := dashboard.NewDashboardProvider(
		*kubernetesProvider.Cluster.Service,
		*kubernetesProvider.Node.Service,
		*kubernetesProvider.Namespace.Service,
		*kubernetesProvider.Pod.Service,
		*dragonProvider.Release.Service,
	)
	prometheusProvider := prometheus.NewPrometheusProvider(db)

	return &Provider{
		JWTManager:        jwtManager,
		AuditManager:      auditManager,
		PermissionManager: permissionManager,
		Audit:             audit.NewAuditProvider(db),
		Authn:             authnProvider,
		Kubernetes:        kubernetesProvider,
		Prometheus:        prometheusProvider,
		Dragon:            dragonProvider,
		Dashboard:         dashboardProvider,
	}
}
