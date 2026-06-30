package dto

import "time"

type ProfileDTO struct {
	ID        int        `json:"id"`
	Username  string     `json:"username"`
	Nickname  string     `json:"nickname,omitempty"`
	Email     string     `json:"email"`
	Phone     string     `json:"phone,omitempty"`
	CreatedAt time.Time  `json:"createdAt"`
	LastLogin *time.Time `json:"lastLogin,omitempty"`
}

type LoginResponseDTO struct {
	UID          int    `json:"uid,omitempty"`
	Username     string `json:"username,omitempty"`
	AccessToken  string `json:"accessToken,omitempty"`
	RefreshToken string `json:"refreshToken,omitempty"`
	Expires      string `json:"expires,omitempty"`
}
