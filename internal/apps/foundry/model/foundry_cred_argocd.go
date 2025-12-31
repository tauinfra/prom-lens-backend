package model

import "time"

type CredArgoCD struct {
	ID          int        `gorm:"primaryKey" json:"id,omitempty"`
	Name        string     `gorm:"type:varchar(32);not null;unique;comment:'ArgoCD 名称'" json:"name,omitempty"`
	Server      string     `gorm:"type:varchar(128);not null;comment:'ArgoCD 地址'" json:"server,omitempty"`
	Username    string     `gorm:"type:varchar(128);comment:'ArgoCD 账号'" json:"username,omitempty"`
	Password    string     `gorm:"type:varchar(256);comment:'ArgoCD 密码，加密存储'" json:"password,omitempty"`
	Description string     `gorm:"type:text;comment:'ArgoCD 描述'" json:"description,omitempty"`
	CreatedAt   *time.Time `gorm:"datetime(3);comment:'创建时间'"                          json:"createdAt,omitempty"`
	UpdatedAt   *time.Time `gorm:"datetime(3);comment:'更新时间'"                          json:"updatedAt,omitempty"`
}

func (CredArgoCD) TableName() string {
	return "valyria_foundry_credential_argocd"
}
