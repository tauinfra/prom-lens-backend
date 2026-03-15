package controller

import (
	"net/http"
	"strconv"
	"valyria-backend/internal/apps/authn/request"
	"valyria-backend/internal/apps/authn/service"
	pg "valyria-backend/internal/core/pagination"
	"valyria-backend/internal/pkg/ginhelper"

	"github.com/gin-gonic/gin"
)

// PermissionController 定义控制器结构体
type PermissionController struct {
	Permission service.PermissionManager // 使用服务接口
}

func NewPermissionController(Permission service.PermissionManager) *PermissionController {
	return &PermissionController{Permission: Permission}
}

func (c *PermissionController) List(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.Query("page"))
	size, _ := strconv.Atoi(ctx.Query("size"))
	// 分页
	params := pg.QueryParams{
		Page:      page,
		Size:      size,
		SortBy:    ctx.Query("sortBy"),
		SortOrder: ctx.Query("sortOrder"),
	}

	data, pagination, err := c.Permission.List(ctx, params)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data, "pagination": pagination})
}

func (c *PermissionController) Get(ctx *gin.Context) {
	permissionID, ok := ginhelper.RequireUintParam(ctx, "permissionID")
	if !ok {
		return
	}
	data, err := c.Permission.Get(ctx, permissionID)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

func (c *PermissionController) Create(ctx *gin.Context) {
	var req request.CreatePermissionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	username := ctx.MustGet("username").(string)
	req.Creator = username
	if err := c.Permission.Create(ctx, &req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200})
}

func (c *PermissionController) Update(ctx *gin.Context) {
	var req request.UpdatePermissionRequest
	permissionID, ok := ginhelper.RequireUintParam(ctx, "permissionID")
	if !ok {
		return
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	if err := c.Permission.Update(ctx, permissionID, &req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200})
}

func (c *PermissionController) Delete(ctx *gin.Context) {
	permissionID, ok := ginhelper.RequireUintParam(ctx, "permissionID")
	if !ok {
		return
	}
	if err := c.Permission.Delete(ctx, permissionID); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200})
}
