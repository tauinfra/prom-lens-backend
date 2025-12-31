package controller

import (
	"net/http"
	"strconv"
	"valyria-backend/internal/apps/kubernetes/service"
	"valyria-backend/internal/core/logger"

	"github.com/gin-gonic/gin"
	corev1 "k8s.io/api/core/v1"
)

// NodeController 定义控制器结构体
type NodeController struct {
	node service.NodeService // 使用服务接口
}

// NewNodeController 创建新的 NodesController 实例
func NewNodeController(node service.NodeService) *NodeController {
	return &NodeController{
		node: node,
	}
}

type Node struct {
	Labels map[string]string `json:"labels,omitempty"`
	Taints []corev1.Taint    `json:"taints,omitempty"`
}

func (c *NodeController) List(ctx *gin.Context) {
	var id, _ = strconv.Atoi(ctx.Param("id"))
	data, err := c.node.List(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

func (c *NodeController) Get(ctx *gin.Context) {
	var (
		id, _ = strconv.Atoi(ctx.Param("id"))
		name  = ctx.Param("name")
	)
	data, err := c.node.Get(ctx, id, name)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

func (c *NodeController) GetDetail(ctx *gin.Context) {
	var (
		id, _ = strconv.Atoi(ctx.Param("id"))
		name  = ctx.Param("name")
	)
	data, err := c.node.GetDetail(ctx, id, name)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

// Labels 更新标签
func (c *NodeController) Labels(ctx *gin.Context) {
	var (
		id, _ = strconv.Atoi(ctx.Param("id"))
		name  = ctx.Param("name")
		req   Node
	)
	// 绑定标签数据
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Errorf("Data binding request failed, error: %v", err)
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	// 获取节点解析
	body, err := c.node.Get(ctx, id, name)
	if err != nil {
		logger.Errorf("Kubernetes Cluster '%v' node '%v' get failed, error: %v", id, name, err)
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	// 更新 Labels (完整替换 Labels)
	body.Labels = req.Labels

	data, err := c.node.Update(ctx, id, body)
	if err != nil {
		logger.Errorf("Kubernetes Cluster '%v' node '%v' labels update failed, error: %v", id, name, err)
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	logger.Infof("Kubernetes Cluster '%v' node '%v' labels has been updateed successfully.", id, name)
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

// Taints 更新污点
func (c *NodeController) Taints(ctx *gin.Context) {
	var (
		id, _ = strconv.Atoi(ctx.Param("id"))
		name  = ctx.Param("name")
		req   Node
	)
	// 绑定标签数据
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Errorf("Data binding request failed, error: %v", err)
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	// 获取节点解析
	body, err := c.node.Get(ctx, id, name)
	if err != nil {
		logger.Errorf("Kubernetes Cluster '%v' node '%v' get failed, error: %v", id, name, err)
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	// 更新 Taints (完整替换 Taints)
	body.Spec.Taints = req.Taints

	data, err := c.node.Update(ctx, id, body)
	if err != nil {
		logger.Errorf("Kubernetes Cluster '%v' node '%v' update failed, error: %v", id, name, err)
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	logger.Infof("Kubernetes Cluster '%v' node '%v' labels has been updateed successfully.", id, name)
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

func (c *NodeController) Cordon(ctx *gin.Context) {
	var (
		id, _ = strconv.Atoi(ctx.Param("id"))
		name  = ctx.Param("name")
	)
	data, err := c.node.Cordon(ctx, id, name)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}
