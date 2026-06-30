package dto

import "time"

type UserDTO struct {
	ID          int        `json:"id"`
	Username    string     `json:"username"`
	Nickname    string     `json:"nickname"`
	Email       string     `json:"email"`
	Phone       string     `json:"phone,omitempty"`
	IsActive    *bool      `json:"isActive"`
	IsSuperuser *bool      `json:"isSuperuser"`
	IsLdap      *bool      `json:"isLdap"`
	DN          string     `json:"dn,omitempty"`
	Creator     string     `json:"creator"`
	LastLogin   *time.Time `json:"lastLogin,omitempty"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}
