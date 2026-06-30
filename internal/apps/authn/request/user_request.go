package request

type CreateUserRequest struct {
	Username    string `json:"username" binding:"required"`
	Password    string `json:"password" binding:"required,min=6"`
	Nickname    string `json:"nickname"`
	Email       string `json:"email" binding:"required,email"`
	Phone       string `json:"phone"`
	IsActive    *bool  `json:"isActive"`
	IsSuperuser *bool  `json:"isSuperuser"`
	IsLdap      *bool  `json:"isLdap"`
	DN          string `json:"dn"`
	Creator     string `json:"-"`
}

type UpdateUserRequest struct {
	Username    *string `json:"username"`
	Nickname    *string `json:"nickname"`
	Email       *string `json:"email" binding:"omitempty,email"`
	Phone       *string `json:"phone"`
	IsActive    *bool   `json:"isActive"`
	IsSuperuser *bool   `json:"isSuperuser"`
	IsLdap      *bool   `json:"isLdap"`
	DN          *string `json:"dn"`
}

type ResetPasswordRequest struct {
	NewPassword string `json:"newPassword" binding:"required,min=6"`
}
