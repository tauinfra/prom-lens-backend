package request

type BatchDeletePodRequest struct {
	Names []string `json:"names" binding:"required,min=1,dive,required"`
}
