package model

import "time"

type CredGitlab struct {
	ID          int        `gorm:"type:bigint;primaryKey"                                 json:"id,omitempty"`
	Name        string     `gorm:"type:varchar(32);not null;unique;comment:'Gitlab 名称'"     json:"name,omitempty"`
	BaseURL     string     `gorm:"type:varchar(128);not null;comment:'Gitlab 地址'"    json:"baseURL,omitempty"`
	Token       string     `gorm:"type:text;comment:'Gitlab Token'"                          json:"token,omitempty"`
	Description string     `gorm:"type:text;comment:'Gitlab 描述'"                            json:"description,omitempty"`
	CreatedAt   *time.Time `gorm:"datetime(3);comment:'创建时间'"                          json:"createdAt,omitempty"`
	UpdatedAt   *time.Time `gorm:"datetime(3);comment:'更新时间'"                          json:"updatedAt,omitempty"`
}

func (CredGitlab) TableName() string {
	return "valyria_foundry_credential_gitlab"
}
