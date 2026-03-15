package controller

import (
	"net/http"
	"strconv"
	"valyria-backend/internal/apps/dragon/request"
	"valyria-backend/internal/apps/dragon/service"
	pg "valyria-backend/internal/core/pagination"
	"valyria-backend/internal/pkg/ginhelper"

	"github.com/gin-gonic/gin"
)

// ProjectController 定义控制器结构体
type ProjectController struct {
	manager service.ProjectManager // 使用服务接口
}

func NewProjectController(manager service.ProjectManager) *ProjectController {
	return &ProjectController{manager: manager}
}

func (c *ProjectController) List(ctx *gin.Context) {
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
	userID := ctx.MustGet("uid").(int)
	isAdmin := ctx.MustGet("isSuperuser").(bool)
	data, pagination, err := c.manager.List(ctx, uint(userID), isAdmin, params)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data, "pagination": pagination})
}

func (c *ProjectController) Get(ctx *gin.Context) {
	id, ok := ginhelper.RequireUintParam(ctx, "projectID")
	if !ok {
		return
	}
	data, err := c.manager.Get(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

func (c *ProjectController) Create(ctx *gin.Context) {
	var req request.CreateProjectRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	username, exists := ctx.Get("username")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"success": false, "code": http.StatusUnauthorized, "msg": "用户未登录"})
		return
	}
	if err := c.manager.Create(ctx, &req, username.(string)); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200})
}

func (c *ProjectController) Update(ctx *gin.Context) {
	id, ok := ginhelper.RequireUintParam(ctx, "projectID")
	if !ok {
		return
	}
	var req request.UpdateProjectRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	if err := c.manager.Update(ctx, id, &req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200})
}

func (c *ProjectController) Delete(ctx *gin.Context) {
	id, ok := ginhelper.RequireUintParam(ctx, "projectID")
	if !ok {
		return
	}
	if err := c.manager.Delete(ctx, id); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200})
}
