package model

import (
	"time"
	"valyria-backend/internal/apps/authn/request"
)

// Menu 表示前端路由菜单结构
type Menu struct {
	ID         uint      `gorm:"type:bigint;primaryKey" json:"id,omitempty"`
	ParentID   *uint     `gorm:"type:bigint;default:null;comment:'父菜单ID'" json:"parentId,omitempty"`
	Name       string    `gorm:"type:varchar(64);not null;comment:'路由名称'" json:"name,omitempty"`
	Path       string    `gorm:"type:varchar(255);not null;comment:'路由路径'" json:"path,omitempty"`
	Component  string    `gorm:"type:varchar(255);comment:'组件路径'" json:"component,omitempty"`
	Title      string    `gorm:"type:varchar(64);not null;comment:'菜单标题'" json:"title,omitempty"`
	Icon       string    `gorm:"type:varchar(64);comment:'图标'" json:"icon,omitempty"`
	Rank       int       `gorm:"type:int;default:0;comment:'排序'" json:"rank,omitempty"`
	ShowParent *bool     `gorm:"type:tinyint(1);default:0;comment:'显示父级'" json:"showParent,omitempty"`
	ShowLink   *bool     `gorm:"type:tinyint(1);default:1;comment:'显示菜单'" json:"showLink,omitempty"`
	Children   []Menu    `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	Creator    string    `gorm:"type:varchar(64);not null;comment:创建人" json:"creator,omitempty"`
	CreatedAt  time.Time `gorm:"datetime(3);comment:'创建时间'" json:"createdAt,omitempty"`
	UpdatedAt  time.Time `gorm:"datetime(3);comment:'更新时间'" json:"updatedAt,omitempty"`
}

func (Menu) TableName() string {
	return "valyria_authn_menu"
}

func (m *Menu) ApplyRequest(req *request.UpdateMenuRequest) {
	if req.ParentID != nil {
		m.ParentID = req.ParentID
	}
	if req.Name != nil {
		m.Name = *req.Name
	}
	if req.Path != nil {
		m.Path = *req.Path
	}
	if req.Component != nil {
		m.Component = *req.Component
	}
	if req.Title != nil {
		m.Title = *req.Title
	}
	if req.Icon != nil {
		m.Icon = *req.Icon
	}
	if req.Rank != nil {
		m.Rank = *req.Rank
	}
	if req.ShowParent != nil {
		m.ShowParent = req.ShowParent
	}
	if req.ShowLink != nil {
		m.ShowLink = req.ShowLink
	}
}
