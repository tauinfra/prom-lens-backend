package controller

import (
	"net/http"
	"strconv"
	"valyria-backend/internal/apps/kubernetes/service"

	"github.com/gin-gonic/gin"
	corev1 "k8s.io/api/core/v1"
)

// ConfigmapController 定义控制器结构体
type ConfigmapController struct {
	configmap service.ConfigmapService // 使用服务接口
}

// NewConfigmapController 创建新的 ConfigmapController 实例
func NewConfigmapController(configmap service.ConfigmapService) *ConfigmapController {
	return &ConfigmapController{configmap: configmap}
}

func (c *ConfigmapController) List(ctx *gin.Context) {
	var (
		id, _ = strconv.Atoi(ctx.Param("id"))
		ns    = ctx.Param("namespace")
	)
	data, err := c.configmap.List(ctx, id, ns)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

func (c *ConfigmapController) Get(ctx *gin.Context) {
	var (
		id, _ = strconv.Atoi(ctx.Param("id"))
		ns    = ctx.Param("namespace")
		name  = ctx.Param("name")
	)
	data, err := c.configmap.Get(ctx, id, ns, name)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

func (c *ConfigmapController) Create(ctx *gin.Context) {
	var (
		id, _ = strconv.Atoi(ctx.Param("id"))
		ns    = ctx.Param("namespace")
		body  = &corev1.ConfigMap{}
	)
	// 绑定数据
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
	}
	// 创建
	data, err := c.configmap.Create(ctx, id, ns, body)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

func (c *ConfigmapController) Update(ctx *gin.Context) {
	var (
		id, _ = strconv.Atoi(ctx.Param("id"))
		ns    = ctx.Param("namespace")
		name  = ctx.Param("name")
		body  = &corev1.ConfigMap{}
	)
	// 绑定数据
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	// 更新
	body.Name = name
	body.Namespace = ns
	data, err := c.configmap.Update(ctx, id, ns, body)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

func (c *ConfigmapController) Delete(ctx *gin.Context) {
	var (
		id, _ = strconv.Atoi(ctx.Param("id"))
		ns    = ctx.Param("namespace")
		name  = ctx.Param("name")
	)
	err := c.configmap.Delete(ctx, id, ns, name)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200})
}
