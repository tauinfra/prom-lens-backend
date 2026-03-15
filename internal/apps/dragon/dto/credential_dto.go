package dto

import "time"

type CredentialDTO struct {
	ID         uint      `json:"id"`
	Name       string    `json:"name"`
	Type       string    `json:"type"`
	SecretType string    `json:"secretType"`
	BaseURL    string    `json:"baseURL"`
	Username   string    `json:"username"`
	Creator    string    `json:"creator"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}
