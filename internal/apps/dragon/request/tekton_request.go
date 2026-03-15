package request

type BatchDeleteRequest struct {
	Names []string `json:"names" binding:"required"`
}
