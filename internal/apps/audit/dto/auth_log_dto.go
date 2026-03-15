package dto

import "time"

type AuthLogDTO struct {
	ID         int       `json:"id"`
	Username   string    `json:"username"`
	IPAddress  string    `json:"ipAddress"`
	System     string    `json:"system"`
	Agent      string    `json:"agent"`
	StatusCode int       `json:"statusCode"`
	Success    bool      `json:"success"`
	ErrorCode  int       `json:"errorCode,omitempty"`
	ErrorMsg   string    `json:"errorMsg,omitempty"`
	CreatedAt  time.Time `json:"createdAt"`
}
