package request

import "gorm.io/datatypes"

type CreateTargetRequest struct {
	IPAddress string         `json:"ipAddress" binding:"required,ipv4"`
	Port      int            `json:"port" binding:"required,min=1,max=65535"`
	Labels    datatypes.JSON `json:"labels" binding:"required"`
	Enabled   *bool          `json:"enabled"`
}

type UpdateTargetRequest struct {
	IPAddress *string         `json:"ipAddress" binding:"omitempty,ipv4"`
	Port      *int            `json:"port" binding:"omitempty,min=1,max=65535"`
	Labels    *datatypes.JSON `json:"labels"`
	Enabled   *bool           `json:"enabled"`
}
