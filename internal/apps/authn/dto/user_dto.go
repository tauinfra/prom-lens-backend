package dto

import "time"

type UserDTO struct {
	ID          int       `json:"id"`
	Username    string    `json:"username"`
	Nickname    string    `json:"nickname"`
	Email       string    `json:"email"`
	Phone       string    `json:"phone"`
	IsActive    *bool     `json:"isActive"`
	IsSuperuser *bool     `json:"isSuperuser"`
	IsLdap      *bool     `json:"isLdap"`
	DN          string    `json:"dn,omitempty"`
	Creator     string    `json:"creator"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type UserSelfDTO struct {
	ID          int      `json:"id"`
	Username    string   `json:"username"`
	Nickname    string   `json:"nickname"`
	Email       string   `json:"email"`
	Phone       string   `json:"phone"`
	IsSuperuser bool     `json:"isSuperuser"`
	Roles       []string `json:"roles,omitempty"`       // 用户角色列表，用于前端菜单控制
	Permissions []string `json:"permissions,omitempty"` // 用户权限列表
}
// ProfileDTO 当前用户 profile 输出（GET /me/profile）
type ProfileDTO struct {
	ID          int        `json:"id"`
	Username    string     `json:"username"`
	Email       string     `json:"email"`
	Phone       string     `json:"phone"`
	CreatedAt   time.Time  `json:"createdAt"`
	LastLogin   *time.Time `json:"lastLogin,omitempty"`
	Roles       []string   `json:"roles"`
	Permissions []string   `json:"permissions"`
}

type LoginResponseDTO struct {
	UID          int      `json:"uid,omitempty"`
	Username     string   `json:"username,omitempty"`
	AccessToken  string   `json:"accessToken,omitempty"`
	RefreshToken string   `json:"refreshToken,omitempty"`
	Expires      string   `json:"expires,omitempty"`
	Roles        []string `json:"roles,omitempty"`
	Permissions  []string `json:"permissions,omitempty"`
}
