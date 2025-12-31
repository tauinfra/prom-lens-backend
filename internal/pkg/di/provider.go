package di

import (
	"valyria-backend/internal/core/configs"
	"valyria-backend/internal/core/middleware"
	"valyria-backend/internal/pkg/di/auth"
	"valyria-backend/internal/pkg/di/foundry"
	"valyria-backend/internal/pkg/di/kingsguard"
	"valyria-backend/internal/pkg/di/kubernetes"
	"valyria-backend/internal/pkg/encryption"

	"gorm.io/gorm"
)

type Provider struct {
	DB                *gorm.DB
	Encryptor         *encryption.Encryptor
	JWTManager        *middleware.JWTManager
	AuditManager      *middleware.AuditManager
	PermissionManager *middleware.PermissionManager
	Auth              *auth.AuthProvider
	Kubernetes        *kubernetes.KubernetesProvider
	KingsGuard        *kingsguard.PolicyProvider
	Foundry           *foundry.FoundryProvider
}

func NewProvider(cfg *configs.Config, db *gorm.DB, encryptor encryption.Encryptor) *Provider {
	// 创建 JWT 管理器
	jwtManager := middleware.NewJWTManager(cfg)
	auditManager := middleware.NewAuditManager(db)
	permissionManager := middleware.NewPermissionManager(db, cfg)

	return &Provider{
		JWTManager:        jwtManager,
		AuditManager:      auditManager,
		PermissionManager: permissionManager,
		Auth:              auth.NewAuthProvider(db, jwtManager),
		Kubernetes:        kubernetes.NewKubernetesProvider(db, encryptor),
		KingsGuard:        kingsguard.NewPolicyProvider(db),
		Foundry:           foundry.NewFoundryProvider(db, encryptor),
	}
}
