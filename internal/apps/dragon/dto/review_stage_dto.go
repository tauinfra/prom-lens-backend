package dto

import "time"

type ReviewStageDTO struct {
	ID        uint      `json:"id"`
	Step      int       `json:"step"`
	RoleID    uint      `json:"roleID"`
	RoleName  string    `json:"roleName"`
	Creator   string    `json:"creator"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
