package dto

import (
	"time"

	"gorm.io/datatypes"
)

type AuditLogDTO struct {
	ID         int            `json:"id"`
	Username   string         `json:"username"`
	UrlPath    string         `json:"urlPath"`
	Method     string         `json:"method"`
	IPAddress  string         `json:"ipAddress"`
	Agent      string         `json:"agent"`
	StatusCode int            `json:"statusCode"`
	Success    bool           `json:"success"`
	Params     datatypes.JSON `json:"params"`
	Response   datatypes.JSON `json:"response"`
	CreatedAt  time.Time      `json:"createdAt"`
}
