package request

type CreateRecordRequest struct {
	Name string `json:"name" binding:"required"`
	Expr string `json:"expr" binding:"required"`
}

type UpdateRecordRequest struct {
	Name *string `json:"name"`
	Expr *string `json:"expr"`
}
