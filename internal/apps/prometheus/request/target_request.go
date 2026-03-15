package request

import "gorm.io/datatypes"

type CreateTargetRequest struct {
	IPAddress string         `json:"ipAddress" binding:"required"`
	Port      int            `json:"port" binding:"required"`
	Labels    datatypes.JSON `json:"labels" binding:"required"`
	Enabled   *bool          `json:"enabled"`
}

type UpdateTargetRequest struct {
	IPAddress *string         `json:"ipAddress"`
	Port      *int            `json:"port"`
	Labels    *datatypes.JSON `json:"labels"`
	Enabled   *bool           `json:"enabled"`
}
