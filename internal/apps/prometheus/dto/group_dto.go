package dto

import "time"

type GroupDTO struct {
	ID          int        `json:"id"`
	Name        string     `json:"name"`
	Type        string     `json:"type"`
	Description string     `json:"description"`
	RuleCount   int64      `json:"ruleCount"`
	RecordCount int64      `json:"recordCount"`
	CreatedAt   *time.Time `json:"createdAt"`
	UpdatedAt   *time.Time `json:"updatedAt"`
}
