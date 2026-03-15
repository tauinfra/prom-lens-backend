package dto

import "time"

type MenuDTO struct {
	ID         uint      `json:"id"`
	ParentID   *uint     `json:"parentId,omitempty"`
	Name       string    `json:"name"`
	Path       string    `json:"path"`
	Component  string    `json:"component,omitempty"`
	Title      string    `json:"title"`
	Icon       string    `json:"icon,omitempty"`
	Rank       int       `json:"rank"`
	ShowParent *bool     `json:"showParent,omitempty"`
	ShowLink   *bool     `json:"showLink,omitempty"`
	Creator    string    `json:"creator"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

// MenuTreeDTO 树形菜单（parentId 为空为一级，有 parentId 的挂在 children 下）
type MenuTreeDTO struct {
	ID         uint          `json:"id"`
	ParentID   *uint         `json:"parentId,omitempty"`
	Name       string        `json:"name"`
	Path       string        `json:"path"`
	Component  string        `json:"component,omitempty"`
	Title      string        `json:"title"`
	Icon       string        `json:"icon,omitempty"`
	Rank       int           `json:"rank"`
	ShowParent *bool        `json:"showParent,omitempty"`
	ShowLink   *bool        `json:"showLink,omitempty"`
	Creator    string       `json:"creator,omitempty"`
	CreatedAt  time.Time    `json:"createdAt,omitempty"`
	UpdatedAt  time.Time    `json:"updatedAt,omitempty"`
	Children   []MenuTreeDTO `json:"children,omitempty"`
}

type RouteMetaDTO struct {
	Title      string `json:"title"`
	Icon       string `json:"icon,omitempty"`
	Rank       int    `json:"rank"`
	ShowParent *bool  `json:"showParent,omitempty"`
	ShowLink   *bool  `json:"showLink,omitempty"`
}

type RouteDTO struct {
	Name      string      `json:"name"`
	Path      string      `json:"path"`
	Component string      `json:"component,omitempty"`
	Meta      RouteMetaDTO `json:"meta"`
	Children  []RouteDTO  `json:"children,omitempty"`
}
