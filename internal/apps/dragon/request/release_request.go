package request

type CreateReleaseRequest struct {
	TaskID        string `json:"taskID,omitempty"`
	ProjectID     uint   `json:"projectID,omitempty"`                 // 项目
	EnvironmentID uint   `json:"environmentID,omitempty"`             // 环境
	PipelineID    uint   `json:"pipelineID,omitempty"`                // 流水线
	GitRef        string `json:"gitRef,omitempty" binding:"required"` // 只需要传递分支参数
	GitCommit     string `json:"gitCommit,omitempty"`
	GitRefType    string `json:"gitRefType,omitempty"`
	ImageRegistry string `json:"imageRegistry,omitempty"`
	ImageName     string `json:"imageName,omitempty"`
	ImageTag      string `json:"imageTag,omitempty"`
	Description   string `json:"description"`
	Creator       string `json:"creator,omitempty"`
}

// RollbackRequest 回滚请求：回滚到指定 release（该 release 必须为发布成功）
type RollbackRequest struct {
	TargetReleaseID uint `json:"targetReleaseID" binding:"required"` // 回滚目标 release id
}
