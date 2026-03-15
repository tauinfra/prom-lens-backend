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

// EnvironmentController 定义控制器结构体
type EnvironmentController struct {
	manager service.EnvironmentManager // 使用服务接口
}

func NewEnvironmentController(manager service.EnvironmentManager) *EnvironmentController {
	return &EnvironmentController{manager: manager}
}

func (c *EnvironmentController) List(ctx *gin.Context) {
	projectID, ok := ginhelper.RequireUintParam(ctx, "projectID")
	if !ok {
		return
	}
	page, _ := strconv.Atoi(ctx.Query("page"))
	size, _ := strconv.Atoi(ctx.Query("size"))
	// 分页
	params := pg.QueryParams{
		Page:      page,
		Size:      size,
		SortBy:    ctx.Query("sortBy"),
		SortOrder: ctx.Query("sortOrder"),
		Keyword:   ctx.Query("keyword"),
		Preloads:  []string{"Project", "Harbor"},
		Filters:   map[string]interface{}{"project_id": projectID}, // EnvID 依赖 ProjID
	}
	userID := ctx.MustGet("uid").(int)
	isAdmin := ctx.MustGet("isSuperuser").(bool)
	data, pagination, err := c.manager.List(ctx, uint(userID), projectID, isAdmin, params)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data, "pagination": pagination})
}

func (c *EnvironmentController) Get(ctx *gin.Context) {
	projectID, ok := ginhelper.RequireUintParam(ctx, "projectID")
	if !ok {
		return
	}
	environmentID, ok := ginhelper.RequireUintParam(ctx, "environmentID")
	if !ok {
		return
	}
	data, err := c.manager.Get(ctx, projectID, environmentID)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

func (c *EnvironmentController) Create(ctx *gin.Context) {
	projectID, ok := ginhelper.RequireUintParam(ctx, "projectID")
	if !ok {
		return
	}
	var req request.CreateEnvironmentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	if err := c.manager.Create(ctx, projectID, &req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200})
}

func (c *EnvironmentController) Update(ctx *gin.Context) {
	projectID, ok := ginhelper.RequireUintParam(ctx, "projectID")
	if !ok {
		return
	}
	environmentID, ok := ginhelper.RequireUintParam(ctx, "environmentID")
	if !ok {
		return
	}
	var req request.UpdateEnvironmentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	if err := c.manager.Update(ctx, projectID, environmentID, &req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200})
}

func (c *EnvironmentController) Delete(ctx *gin.Context) {
	projectID, ok := ginhelper.RequireUintParam(ctx, "projectID")
	if !ok {
		return
	}
	environmentID, ok := ginhelper.RequireUintParam(ctx, "environmentID")
	if !ok {
		return
	}
	if err := c.manager.Delete(ctx, projectID, environmentID); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200})
}
