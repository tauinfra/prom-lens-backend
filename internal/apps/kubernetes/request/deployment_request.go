package request

type ScaleRequest struct {
	Replicas int32 `json:"replicas" binding:"required"`
}

type RolloutRequest struct {
	Name string `json:"name" binding:"required"`
}
