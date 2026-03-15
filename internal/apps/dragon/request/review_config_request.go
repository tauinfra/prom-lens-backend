package request

type CreateReviewStageRequest struct {
	Step   int  `json:"step" binding:"required"`
	RoleID uint `json:"roleID" binding:"required"`
}

type UpdateReviewStageRequest struct {
	Step   *int  `json:"step"`
	RoleID *uint `json:"roleID"`
}
