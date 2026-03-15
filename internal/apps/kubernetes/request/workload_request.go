package request

type BatchDeleteWorkloadRequest struct {
	Names []string `json:"names" binding:"required,min=1,dive,required"`
}

type BatchDeleteIDRequest struct {
	IDs []uint `json:"ids" binding:"required,min=1,dive,required"`
}
