package request

type CreatePermissionRequest struct {
	ClusterID uint   `json:"clusterID" binding:"required"`
	UserID    uint   `json:"userID" binding:"required"`
	Namespace string `json:"namespace" binding:"required"`
	Role      string `json:"role" binding:"required,permission_role"`
}

// BatchCreatePermissionBinding 按 cluster 分组时，单集群下的一条绑定
type BatchCreatePermissionBinding struct {
	Namespace string `json:"namespace" binding:"required"`
	Role      string `json:"role" binding:"required,permission_role"`
}

// BatchCreatePermissionCluster 按 cluster 分组时，单集群及其 bindings
type BatchCreatePermissionCluster struct {
	ClusterID uint                           `json:"clusterID" binding:"required"`
	Bindings  []BatchCreatePermissionBinding `json:"bindings" binding:"required,min=1,dive"`
}

// BatchCreatePermissionRequest 批量同步权限：同一 userID，按集群分组
type BatchCreatePermissionRequest struct {
	UserID   uint                           `json:"userID" binding:"required"`
	Clusters []BatchCreatePermissionCluster `json:"clusters" binding:"required,min=1,dive"`
}
