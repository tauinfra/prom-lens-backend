package controller

import (
	"net/http"
	"strconv"
	"valyria-backend/internal/apps/kubernetes/service"

	"github.com/gin-gonic/gin"
	corev1 "k8s.io/api/core/v1"
)

// ServiceController 定义控制器结构体
type ServiceController struct {
	service service.ServiceManager // 使用服务接口
}

// NewServiceController 创建新的 ServiceController 实例
func NewServiceController(service service.ServiceManager) *ServiceController {
	return &ServiceController{service: service}
}

func (c *ServiceController) List(ctx *gin.Context) {
	var (
		id, _         = strconv.Atoi(ctx.Param("id"))
		ns            = ctx.Param("namespace")
		labelSelector = ctx.Query("labelSelector")
	)
	data, err := c.service.List(ctx, id, ns, labelSelector)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

func (c *ServiceController) Get(ctx *gin.Context) {
	var (
		id, _ = strconv.Atoi(ctx.Param("id"))
		ns    = ctx.Param("namespace")
		name  = ctx.Param("name")
	)
	data, err := c.service.Get(ctx, id, ns, name)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

func (c *ServiceController) Create(ctx *gin.Context) {
	var (
		id, _ = strconv.Atoi(ctx.Param("id"))
		ns    = ctx.Param("namespace")
		body  = &corev1.Service{}
	)
	// 绑定数据
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
	}
	// 创建
	data, err := c.service.Create(ctx, id, ns, body)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

func (c *ServiceController) Update(ctx *gin.Context) {
	var (
		id, _ = strconv.Atoi(ctx.Param("id"))
		ns    = ctx.Param("namespace")
		name  = ctx.Param("name")
		body  = &corev1.Service{}
	)
	// 绑定数据
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	// 更新
	body.Name = name
	body.Namespace = ns
	data, err := c.service.Update(ctx, id, ns, body)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

func (c *ServiceController) Delete(ctx *gin.Context) {
	var (
		id, _ = strconv.Atoi(ctx.Param("id"))
		ns    = ctx.Param("namespace")
		name  = ctx.Param("name")
	)
	err := c.service.Delete(ctx, id, ns, name)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200})
}
