package request

import "gorm.io/datatypes"

type CreateRuleRequest struct {
	Name        string         `json:"name" binding:"required"`
	Summary     string         `json:"summary" binding:"required"`
	Description string         `json:"description" binding:"required"`
	Expr        string         `json:"expr" binding:"required"`
	For         string         `json:"for" binding:"required"`
	Labels      datatypes.JSON `json:"labels" binding:"required"`
	Status      *bool          `json:"status"`
}

type UpdateRuleRequest struct {
	Name        *string         `json:"name"`
	Summary     *string         `json:"summary"`
	Description *string         `json:"description"`
	Expr        *string         `json:"expr"`
	For         *string         `json:"for"`
	Labels      *datatypes.JSON `json:"labels"`
	Status      *bool           `json:"status"`
}
