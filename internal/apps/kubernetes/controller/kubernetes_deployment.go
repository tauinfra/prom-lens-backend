package controller

import (
	"net/http"
	"strconv"
	"valyria-backend/internal/apps/kubernetes/service"

	"github.com/gin-gonic/gin"
	appsv1 "k8s.io/api/apps/v1"
)

// DeploymentController 定义控制器结构体
type DeploymentController struct {
	deployment service.DeploymentService // 使用服务接口
}

// NewDeploymentController 创建新的 DeploymentsController 实例
func NewDeploymentController(deployment service.DeploymentService) *DeploymentController {
	return &DeploymentController{deployment: deployment}
}

type Scale struct {
	Replicas int32 `json:"replicas" binding:"required"`
}

type ReplicaSet struct {
	Name string `json:"name" binding:"required"`
}

func (c *DeploymentController) List(ctx *gin.Context) {
	var (
		id, _ = strconv.Atoi(ctx.Param("id"))
		ns    = ctx.Param("namespace")
	)
	data, err := c.deployment.List(ctx, id, ns)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

func (c *DeploymentController) Get(ctx *gin.Context) {
	var (
		id, _ = strconv.Atoi(ctx.Param("id"))
		ns    = ctx.Param("namespace")
		name  = ctx.Param("name")
	)
	data, err := c.deployment.Get(ctx, id, ns, name)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

func (c *DeploymentController) GetDetail(ctx *gin.Context) {
	var (
		id, _ = strconv.Atoi(ctx.Param("id"))
		ns    = ctx.Param("namespace")
		name  = ctx.Param("name")
	)
	data, err := c.deployment.GetDetail(ctx, id, ns, name)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

func (c *DeploymentController) Create(ctx *gin.Context) {
	var (
		id, _ = strconv.Atoi(ctx.Param("id"))
		ns    = ctx.Param("namespace")
		body  = &appsv1.Deployment{}
	)
	// 绑定数据
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
	}
	// 创建
	data, err := c.deployment.Create(ctx, id, ns, body)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

func (c *DeploymentController) Update(ctx *gin.Context) {
	var (
		id, _ = strconv.Atoi(ctx.Param("id"))
		ns    = ctx.Param("namespace")
		name  = ctx.Param("name")
		body  = &appsv1.Deployment{}
	)
	// 绑定数据
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	// 更新
	body.Name = name
	body.Namespace = ns
	data, err := c.deployment.Update(ctx, id, ns, body)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

func (c *DeploymentController) Delete(ctx *gin.Context) {
	var (
		id, _ = strconv.Atoi(ctx.Param("id"))
		ns    = ctx.Param("namespace")
		name  = ctx.Param("name")
	)
	err := c.deployment.Delete(ctx, id, ns, name)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200})
}

func (c *DeploymentController) Scale(ctx *gin.Context) {
	var (
		id, _ = strconv.Atoi(ctx.Param("id"))
		ns    = ctx.Param("namespace")
		name  = ctx.Param("name")
		scale Scale
	)
	if err := ctx.ShouldBindJSON(&scale); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	// 更新副本
	data, err := c.deployment.Scale(ctx, id, ns, name, scale.Replicas)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

func (c *DeploymentController) Restart(ctx *gin.Context) {
	var (
		id, _ = strconv.Atoi(ctx.Param("id"))
		ns    = ctx.Param("namespace")
		name  = ctx.Param("name")
	)
	// 创建
	data, err := c.deployment.Restart(ctx, id, ns, name)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

func (c *DeploymentController) Rollout(ctx *gin.Context) {
	var (
		id, _      = strconv.Atoi(ctx.Param("id"))
		ns         = ctx.Param("namespace")
		name       = ctx.Param("name")
		replicaSet ReplicaSet
	)
	if err := ctx.ShouldBindJSON(&replicaSet); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	// 创建
	data, err := c.deployment.Rollout(ctx, id, ns, name, replicaSet.Name)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}
