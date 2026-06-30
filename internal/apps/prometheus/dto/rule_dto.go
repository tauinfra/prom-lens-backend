package dto

import (
	"time"

	"gorm.io/datatypes"
)

type RuleDTO struct {
	ID          uint           `json:"id"`
	Name        string         `json:"name"`
	GroupID     int            `json:"groupID"`
	Summary     string         `json:"summary"`
	Description string         `json:"description"`
	Expr        string         `json:"expr"`
	For         string         `json:"for"`
	Labels            datatypes.JSON `json:"labels"`
	ExtraAnnotations  datatypes.JSON `json:"extraAnnotations"`
	Status            *bool          `json:"status"`
	CreatedAt   *time.Time     `json:"createdAt"`
	UpdatedAt   *time.Time     `json:"updatedAt"`
}
