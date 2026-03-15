package request

type CreatePipelineRequest struct {
	Name            string `json:"name" binding:"required"`
	EnvironmentID   uint   `json:"environmentID"`
	GitlabID        uint   `json:"gitlabID" binding:"required"`
	GitlabGroupID   uint   `json:"gitlabGroupID" binding:"required"`
	GitlabProjectID uint   `json:"gitlabProjectID" binding:"required"`
	IsApproval      bool   `json:"isApproval"`
	Task            string `json:"task" binding:"required"`
	Script          string `json:"script" binding:"required"`
	Creator         string ` json:"creator"`
}

type UpdatePipelineRequest struct {
	Name            *string `json:"name"`
	GitlabID        *uint   `json:"gitlabID"`
	GitlabGroupID   *uint   `json:"gitlabGroupID"`
	GitlabProjectID *uint   `json:"gitlabProjectID"`
	IsApproval      *bool   `json:"isApproval"`
	Task            *string `json:"task"`
	Script          *string `json:"script"`
}
