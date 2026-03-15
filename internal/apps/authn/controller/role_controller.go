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

// RoleController 定义控制器结构体
type RoleController struct {
	Role service.RoleManager // 使用服务接口
}

func NewRoleController(Role service.RoleManager) *RoleController {
	return &RoleController{Role: Role}
}

func (c *RoleController) List(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.Query("page"))
	size, _ := strconv.Atoi(ctx.Query("size"))
	// 分页
	params := pg.QueryParams{
		Page:      page,
		Size:      size,
		SortBy:    ctx.Query("sortBy"),
		SortOrder: ctx.Query("sortOrder"),
		Keyword:   ctx.Query("keyword"),
	}

	data, pagination, err := c.Role.List(ctx, params)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data, "pagination": pagination})
}

func (c *RoleController) Get(ctx *gin.Context) {
	roleID, ok := ginhelper.RequireUintParam(ctx, "roleID")
	if !ok {
		return
	}
	data, err := c.Role.Get(ctx, roleID)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

func (c *RoleController) Create(ctx *gin.Context) {
	var req request.CreateRoleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	username := ctx.MustGet("username").(string)
	req.Creator = username
	if err := c.Role.Create(ctx, &req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200})
}

func (c *RoleController) Update(ctx *gin.Context) {
	var req request.UpdateRoleRequest
	roleID, ok := ginhelper.RequireUintParam(ctx, "roleID")
	if !ok {
		return
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	if err := c.Role.Update(ctx, roleID, &req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200})
}

func (c *RoleController) Delete(ctx *gin.Context) {
	roleID, ok := ginhelper.RequireUintParam(ctx, "roleID")
	if !ok {
		return
	}

	if err := c.Role.Delete(ctx, roleID); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200})
}

func (c *RoleController) GetRolePermissions(ctx *gin.Context) {
	roleID, ok := ginhelper.RequireUintParam(ctx, "roleID")
	if !ok {
		return
	}
	data, err := c.Role.GetRolePermissions(ctx, roleID)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

func (c *RoleController) UpdateRolePermissions(ctx *gin.Context) {
	var req request.UpdateRolePermissionsRequest
	roleID, ok := ginhelper.RequireUintParam(ctx, "roleID")
	if !ok {
		return
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	if err := c.Role.UpdateRolePermissions(ctx, roleID, req.Permissions); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200})
}
