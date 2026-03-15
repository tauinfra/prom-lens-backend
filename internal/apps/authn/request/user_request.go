package request

type CreateUserRequest struct {
	Username    string `json:"username" binding:"required"`
	Password    string `json:"password"`
	Nickname    string `json:"nickname"`
	Email       string `json:"email" binding:"email"`
	Phone       string `json:"phone"`
	IsActive    *bool  `json:"isActive"`
	IsSuperuser *bool  `json:"isSuperuser"`
	IsLdap      *bool  `json:"isLdap"`
	DN          string `json:"dn"`
	Creator     string `json:"creator"`
}

type UpdateUserRequest struct {
	Username    *string `json:"username"`
	Nickname    *string `json:"nickname"`
	Email       *string `json:"email"`
	Phone       *string `json:"phone"`
	IsActive    *bool   `json:"isActive"`
	IsSuperuser *bool   `json:"isSuperuser"`
	IsLdap      *bool   `json:"isLdap"`
	DN          *string `json:"dn"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required"`
}

type ResetPasswordRequest struct {
	NewPassword string `json:"newPassword" binding:"required"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}
