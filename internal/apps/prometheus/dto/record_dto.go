package dto

import "time"

type RecordDTO struct {
	ID        int        `json:"id"`
	Name      string     `json:"name"`
	GroupID   int        `json:"groupID"`
	Expr      string     `json:"expr"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt"`
}
