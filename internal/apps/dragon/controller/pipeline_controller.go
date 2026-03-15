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

// PipelineController 定义控制器结构体
type PipelineController struct {
	manager service.PipelineManager // 使用服务接口
}

func NewPipelineController(manager service.PipelineManager) *PipelineController {
	return &PipelineController{manager: manager}
}

func (c *PipelineController) List(ctx *gin.Context) {
	projectID, ok := ginhelper.RequireUintParam(ctx, "projectID")
	if !ok {
		return
	}
	environmentID, ok := ginhelper.RequireUintParam(ctx, "environmentID")
	if !ok {
		return
	}
	page, _ := strconv.Atoi(ctx.Query("page"))
	size, _ := strconv.Atoi(ctx.Query("size"))
	params := pg.QueryParams{
		Page:      page,
		Size:      size,
		SortBy:    ctx.Query("sortBy"),
		SortOrder: ctx.Query("sortOrder"),
		Keyword:   ctx.Query("keyword"),
		Preloads:  []string{"Environment.Project", "Environment.Harbor", "Gitlab"},
	}
	data, pagination, err := c.manager.List(ctx, projectID, environmentID, params)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data, "pagination": pagination})
}

func (c *PipelineController) Get(ctx *gin.Context) {
	projectID, ok := ginhelper.RequireUintParam(ctx, "projectID")
	if !ok {
		return
	}
	environmentID, ok := ginhelper.RequireUintParam(ctx, "environmentID")
	if !ok {
		return
	}
	pipelineID, ok := ginhelper.RequireUintParam(ctx, "pipelineID")
	if !ok {
		return
	}
	data, err := c.manager.Get(ctx, projectID, environmentID, pipelineID)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

func (c *PipelineController) Create(ctx *gin.Context) {
	projectID, ok := ginhelper.RequireUintParam(ctx, "projectID")
	if !ok {
		return
	}
	environmentID, ok := ginhelper.RequireUintParam(ctx, "environmentID")
	if !ok {
		return
	}
	var req request.CreatePipelineRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	username := ctx.MustGet("username").(string)
	req.Creator = username
	if err := c.manager.Create(ctx, projectID, environmentID, &req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	} else {
		ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200})
	}
}

func (c *PipelineController) Update(ctx *gin.Context) {
	projectID, ok := ginhelper.RequireUintParam(ctx, "projectID")
	if !ok {
		return
	}
	environmentID, ok := ginhelper.RequireUintParam(ctx, "environmentID")
	if !ok {
		return
	}
	pipelineID, ok := ginhelper.RequireUintParam(ctx, "pipelineID")
	if !ok {
		return
	}
	var req request.UpdatePipelineRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	if err := c.manager.Update(ctx, projectID, environmentID, pipelineID, &req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "msg": "更新成功"})
}

func (c *PipelineController) Delete(ctx *gin.Context) {
	projectID, ok := ginhelper.RequireUintParam(ctx, "projectID")
	if !ok {
		return
	}
	environmentID, ok := ginhelper.RequireUintParam(ctx, "environmentID")
	if !ok {
		return
	}
	pipelineID, ok := ginhelper.RequireUintParam(ctx, "pipelineID")
	if !ok {
		return
	}
	if err := c.manager.Delete(ctx, projectID, environmentID, pipelineID); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200})
}
