package controller

import (
	"net/http"
	"strconv"
	"valyria-backend/internal/apps/kubernetes/request"
	"valyria-backend/internal/apps/kubernetes/service"

	"github.com/gin-gonic/gin"
)

type PermissionController struct {
	svc service.PermissionManager
}

func NewPermissionController(svc service.PermissionManager) *PermissionController {
	return &PermissionController{svc: svc}
}

// List 获取权限列表 GET /kubernetes/rbac/permissions?clusterID=1（可选）
func (c *PermissionController) List(ctx *gin.Context) {
	clusterID, _ := strconv.Atoi(ctx.Query("clusterID"))
	data, err := c.svc.List(ctx, uint(clusterID))
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

// Create 添加权限 POST /kubernetes/rbac/permissions
func (c *PermissionController) Create(ctx *gin.Context) {
	var req request.CreatePermissionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	var creator string
	if v, exists := ctx.Get("username"); exists {
		if s, ok := v.(string); ok && s != "" {
			creator = s
		}
	}
	if err := c.svc.Create(ctx, &req, creator); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200})
}

// CreateBatch 批量添加权限（同一用户，不同集群/命名空间/角色）POST /kubernetes/permissions/batch
func (c *PermissionController) CreateBatch(ctx *gin.Context) {
	var req request.BatchCreatePermissionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	var creator string
	if v, exists := ctx.Get("username"); exists {
		if s, ok := v.(string); ok && s != "" {
			creator = s
		}
	}
	if err := c.svc.CreateBatch(ctx, &req, creator); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200})
}

// Delete 删除权限 DELETE /kubernetes/permissions/:id
func (c *PermissionController) Delete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": "invalid id"})
		return
	}
	if err := c.svc.Delete(ctx, uint(id)); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200})
}

func (c *PermissionController) DeleteBatch(ctx *gin.Context) {
	deleteBatchByIDs(ctx, func(id uint) error {
		return c.svc.Delete(ctx, id)
	})
}
