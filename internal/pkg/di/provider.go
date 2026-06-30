package di

import (
	"prom-lens-backend/internal/core/config"
	coreMiddleware "prom-lens-backend/internal/core/middleware"
	"prom-lens-backend/internal/pkg/di/audit"
	"prom-lens-backend/internal/pkg/di/alerting"
	"prom-lens-backend/internal/pkg/di/authn"
	"prom-lens-backend/internal/pkg/di/prometheus"

	"gorm.io/gorm"
)

type Provider struct {
	JWTManager   *coreMiddleware.JWTManager
	AuditManager *coreMiddleware.AuditManager
	Audit        *audit.Provider
	Authn        *authn.AuthnProvider
	Prometheus   *prometheus.Provider
	Alerting     *alerting.Provider
}

func NewProvider(cfg *config.Config, db *gorm.DB) *Provider {
	jwtManager := coreMiddleware.NewJWTManager(cfg)
	authnProvider := authn.NewAuthnProvider(db, jwtManager)
	auditManager := coreMiddleware.NewAuditManager(db)
	prometheusProvider := prometheus.NewPrometheusProvider(db)
	alertingProvider := alerting.NewAlertingProvider(db, cfg)

	return &Provider{
		JWTManager:   jwtManager,
		AuditManager: auditManager,
		Audit:        audit.NewAuditProvider(db),
		Authn:        authnProvider,
		Prometheus:   prometheusProvider,
		Alerting:     alertingProvider,
	}
}
