package model

type Policy struct {
	ID    int    `gorm:"type:bigint;primaryKey" json:"id"`
	Ptype string `json:"ptype"`                 // 策略类型
	V0    string `json:"v0" binding:"required"` // 用户
	V1    string `json:"v1" binding:"required"` // 集群
	V2    string `json:"v2" binding:"required"` // 命名空间
	V3    string `json:"v3" binding:"required"` // 资源
	V4    string `json:"v4" binding:"required"` // 动作
	V5    string `json:"v5"`                    // 可选
}

func (Policy) TableName() string {
	return "valyria_kubernetes_policy"
}
