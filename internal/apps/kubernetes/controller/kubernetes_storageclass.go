package controller

import (
	"net/http"
	"valyria-backend/internal/apps/kubernetes/service"
	"valyria-backend/internal/pkg/ginhelper"

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
	idU, ok := ginhelper.RequireUintParam(ctx, "id")
	if !ok {
		return
	}
	id := idU
	data, err := c.storageClass.List(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

func (c *StorageClassController) Get(ctx *gin.Context) {
	idU, ok := ginhelper.RequireUintParam(ctx, "id")
	if !ok {
		return
	}
	id := idU
	name := ctx.Param("name")
	data, err := c.storageClass.Get(ctx, id, name)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

func (c *StorageClassController) Create(ctx *gin.Context) {
	idU, ok := ginhelper.RequireUintParam(ctx, "id")
	if !ok {
		return
	}
	id := idU
	body := &storagev1.StorageClass{}
	// 绑定数据
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	// 创建
	_, err := c.storageClass.Create(ctx, id, body)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200})
}

func (c *StorageClassController) Update(ctx *gin.Context) {
	idU, ok := ginhelper.RequireUintParam(ctx, "id")
	if !ok {
		return
	}
	id := idU
	name := ctx.Param("name")
	body := &storagev1.StorageClass{}
	// 绑定数据
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	// 更新
	body.Name = name
	_, err := c.storageClass.Update(ctx, id, body)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200})
}

func (c *StorageClassController) Delete(ctx *gin.Context) {
	idU, ok := ginhelper.RequireUintParam(ctx, "id")
	if !ok {
		return
	}
	id := idU
	name := ctx.Param("name")
	err := c.storageClass.Delete(ctx, id, name)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200})
}

func (c *StorageClassController) DeleteBatch(ctx *gin.Context) {
	idU, ok := ginhelper.RequireUintParam(ctx, "id")
	if !ok {
		return
	}
	id := idU
	deleteBatchByNames(ctx, func(name string) error {
		return c.storageClass.Delete(ctx, id, name)
	})
}
