package controller

import (
	"net/http"
	"strconv"
	"valyria-backend/internal/apps/kubernetes/service"

	"github.com/gin-gonic/gin"
	corev1 "k8s.io/api/core/v1"
)

// PersistentVolumeController 定义控制器结构体
type PersistentVolumeController struct {
	persistentVolume service.PersistentVolumeManager // 使用服务接口
}

// NewPersistentVolumeController 创建新的 PersistentVolumeController 实例
func NewPersistentVolumeController(persistentVolume service.PersistentVolumeManager) *PersistentVolumeController {
	return &PersistentVolumeController{persistentVolume: persistentVolume}
}

func (c *PersistentVolumeController) List(ctx *gin.Context) {
	var (
		id, _ = strconv.Atoi(ctx.Param("id"))
	)
	data, err := c.persistentVolume.List(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

func (c *PersistentVolumeController) Get(ctx *gin.Context) {
	var (
		id, _ = strconv.Atoi(ctx.Param("id"))
		name  = ctx.Param("name")
	)
	data, err := c.persistentVolume.Get(ctx, id, name)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

func (c *PersistentVolumeController) Create(ctx *gin.Context) {
	var (
		id, _ = strconv.Atoi(ctx.Param("id"))
		body  = &corev1.PersistentVolume{}
	)
	// 绑定数据
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
	}
	// 创建
	data, err := c.persistentVolume.Create(ctx, id, body)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

func (c *PersistentVolumeController) Update(ctx *gin.Context) {
	var (
		id, _ = strconv.Atoi(ctx.Param("id"))
		name  = ctx.Param("name")
		body  = &corev1.PersistentVolume{}
	)
	// 绑定数据
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	// 更新
	body.Name = name
	data, err := c.persistentVolume.Update(ctx, id, body)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

func (c *PersistentVolumeController) Delete(ctx *gin.Context) {
	var (
		id, _ = strconv.Atoi(ctx.Param("id"))
		name  = ctx.Param("name")
	)
	err := c.persistentVolume.Delete(ctx, id, name)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200})
}
