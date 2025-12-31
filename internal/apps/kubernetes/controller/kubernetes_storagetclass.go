package controller

import (
	"net/http"
	"strconv"
	"valyria-backend/internal/apps/kubernetes/service"

	"github.com/gin-gonic/gin"
	storagev1 "k8s.io/api/storage/v1"
)

// StorageClassController 定义控制器结构体
type StorageClassController struct {
	storageClass service.StorageClassManager // 使用服务接口
}

// NewStorageClassController 创建新的 StorageClassController 实例
func NewStorageClassController(storageClass service.StorageClassManager) *StorageClassController {
	return &StorageClassController{storageClass: storageClass}
}

func (c *StorageClassController) List(ctx *gin.Context) {
	var (
		id, _ = strconv.Atoi(ctx.Param("id"))
	)
	data, err := c.storageClass.List(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

func (c *StorageClassController) Get(ctx *gin.Context) {
	var (
		id, _ = strconv.Atoi(ctx.Param("id"))
		name  = ctx.Param("name")
	)
	data, err := c.storageClass.Get(ctx, id, name)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

func (c *StorageClassController) Create(ctx *gin.Context) {
	var (
		id, _ = strconv.Atoi(ctx.Param("id"))
		body  = &storagev1.StorageClass{}
	)
	// 绑定数据
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
	}
	// 创建
	data, err := c.storageClass.Create(ctx, id, body)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

func (c *StorageClassController) Update(ctx *gin.Context) {
	var (
		id, _ = strconv.Atoi(ctx.Param("id"))
		name  = ctx.Param("name")
		body  = &storagev1.StorageClass{}
	)
	// 绑定数据
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	// 更新
	body.Name = name
	data, err := c.storageClass.Update(ctx, id, body)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

func (c *StorageClassController) Delete(ctx *gin.Context) {
	var (
		id, _ = strconv.Atoi(ctx.Param("id"))
		name  = ctx.Param("name")
	)
	err := c.storageClass.Delete(ctx, id, name)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200})
}
