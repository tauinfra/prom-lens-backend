package controller

import (
	"net/http"
	"valyria-backend/internal/apps/kubernetes/service"
	"valyria-backend/internal/pkg/ginhelper"

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
	idU, ok := ginhelper.RequireUintParam(ctx, "id")
	if !ok {
		return
	}
	id := idU
	data, err := c.persistentVolume.List(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

func (c *PersistentVolumeController) Get(ctx *gin.Context) {
	idU, ok := ginhelper.RequireUintParam(ctx, "id")
	if !ok {
		return
	}
	id := idU
	name := ctx.Param("name")
	data, err := c.persistentVolume.Get(ctx, id, name)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

func (c *PersistentVolumeController) Create(ctx *gin.Context) {
	idU, ok := ginhelper.RequireUintParam(ctx, "id")
	if !ok {
		return
	}
	id := idU
	body := &corev1.PersistentVolume{}
	// 绑定数据
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	// 创建
	_, err := c.persistentVolume.Create(ctx, id, body)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200})
}

func (c *PersistentVolumeController) Update(ctx *gin.Context) {
	idU, ok := ginhelper.RequireUintParam(ctx, "id")
	if !ok {
		return
	}
	id := idU
	name := ctx.Param("name")
	body := &corev1.PersistentVolume{}
	// 绑定数据
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	// 更新
	body.Name = name
	_, err := c.persistentVolume.Update(ctx, id, body)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200})
}

func (c *PersistentVolumeController) Delete(ctx *gin.Context) {
	idU, ok := ginhelper.RequireUintParam(ctx, "id")
	if !ok {
		return
	}
	id := idU
	name := ctx.Param("name")
	err := c.persistentVolume.Delete(ctx, id, name)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200})
}

func (c *PersistentVolumeController) DeleteBatch(ctx *gin.Context) {
	idU, ok := ginhelper.RequireUintParam(ctx, "id")
	if !ok {
		return
	}
	id := idU
	deleteBatchByNames(ctx, func(name string) error {
		return c.persistentVolume.Delete(ctx, id, name)
	})
}
