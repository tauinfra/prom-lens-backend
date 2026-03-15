package audit

import "gorm.io/gorm"

type Provider struct {
	AuditLog *AuditLogProvider
	AuthLog  *AuthLogProvider
}

func NewAuditProvider(db *gorm.DB) *Provider {
	return &Provider{
		AuditLog: NewAuditLogProvider(db),
		AuthLog:  NewAuthLogProvider(db),
	}
}
