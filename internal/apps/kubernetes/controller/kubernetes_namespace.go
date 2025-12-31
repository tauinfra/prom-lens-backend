package controller

import (
	"net/http"
	"strconv"
	"valyria-backend/internal/apps/kubernetes/service"
	"valyria-backend/internal/core/logger"

	"github.com/gin-gonic/gin"
	coreV1 "k8s.io/api/core/v1"
)

// NamespaceController 定义控制器结构体
type NamespaceController struct {
	namespace service.NamespaceService // 使用服务接口
}

// NewNamespaceController 创建新的 NamespacesController 实例
func NewNamespaceController(namespace service.NamespaceService) *NamespaceController {
	return &NamespaceController{
		namespace: namespace,
	}
}

func (c *NamespaceController) List(ctx *gin.Context) {
	var id, _ = strconv.Atoi(ctx.Param("id"))
	data, err := c.namespace.List(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

func (c *NamespaceController) Create(ctx *gin.Context) {
	var (
		id, _     = strconv.Atoi(ctx.Param("id"))
		body      = &coreV1.Namespace{}
		namespace struct{ name string }
	)
	if err := ctx.ShouldBindJSON(&namespace); err != nil {
		logger.Errorf("Data binding request failed, error: %v", err)
		ctx.JSON(http.StatusOK, gin.H{"code": -1, "msg": err.Error()})
		return
	}
	body.Name = namespace.name
	data, err := c.namespace.Create(ctx, id, body)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": -1, "msg": err.Error()})
		return
	}
	logger.Infof("Kubernetes cluster '%v' namespace '%v' has been created successfully.", id, namespace.name)
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}
