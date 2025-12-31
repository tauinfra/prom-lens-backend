package model

type Policy struct {
	Username  string `json:"username"`
	Cluster   string `json:"cluster"`
	Namespace string `json:"namespace"`
	Resource  string `json:"resource"`
	Action    string `json:"action"` // 动作，如 GET、POST、PUT、DELETE、LIST
}

type CasbinPolicy struct {
	ID    int    `gorm:"type:bigint;primaryKey" json:"id"`
	Ptype string `json:"ptype"`                 // 策略类型
	V0    string `json:"v0" binding:"required"` // 用户
	V1    string `json:"v1" binding:"required"` // 集群
	V2    string `json:"v2" binding:"required"` // 命名空间
	V3    string `json:"v3" binding:"required"` // 资源
	V4    string `json:"v4" binding:"required"` // 资源
	V5    string `json:"v5"`                    // 可选
}

func (CasbinPolicy) TableName() string {
	return "valyria_kingsguard_casbin_policy"
}
