package model

import (
	"time"
)

type Credential struct {
	ID              uint      `gorm:"type:bigint;primaryKey;autoIncrement" json:"id,omitempty"`
	Name            string    `gorm:"type:varchar(255);not null;comment:'凭证名称'" json:"name,omitempty" binding:"required"`
	Type            string    `gorm:"type:enum('gitlab','harbor');not null;comment:'凭证类型'" json:"type,omitempty" binding:"required"`                             // 凭证类型: gitlab, harbor, argocd
	SecretType      string    `gorm:"type:enum('token','password');not null;default:'token';comment:'Secret 类型'" json:"secretType,omitempty" binding:"required"` // secret类型: token, password
	BaseURL         string    `gorm:"type:varchar(255);comment:'服务地址'" json:"baseURL,omitempty"`
	Username        string    `gorm:"type:varchar(64);comment:'用户名称'" json:"username,omitempty"`
	EncryptedSecret string    `gorm:"type:text;not null;comment:'Secret 加密'" json:"encryptedSecret,omitempty"` // 加密后的凭证
	Creator         string    `gorm:"type:varchar(64);not null;comment:创建人" json:"creator,omitempty"`
	CreatedAt       time.Time `gorm:"datetime(3);comment:'创建时间'"       json:"createdAt,omitempty"`
	UpdatedAt       time.Time `gorm:"datetime(3);comment:'更新时间'"       json:"updatedAt,omitempty"`
}

const (
	CredentialTypeGitlab = "gitlab"
	CredentialTypeHarbor = "harbor"
)

func IsValidCredentialType(value string) bool {
	switch value {
	case CredentialTypeGitlab, CredentialTypeHarbor:
		return true
	default:
		return false
	}
}

func (Credential) TableName() string {
	return "valyria_dragon_credential"
}
