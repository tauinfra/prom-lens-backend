package request

import "gorm.io/datatypes"

type CreateTargetGroupRequest struct {
	Name        string         `json:"name" binding:"required"`
	Description string         `json:"description"`
	Labels      datatypes.JSON `json:"labels" binding:"required"`
}

type UpdateTargetGroupRequest struct {
	Name        *string         `json:"name"`
	Description *string         `json:"description"`
	Labels      *datatypes.JSON `json:"labels"`
}
