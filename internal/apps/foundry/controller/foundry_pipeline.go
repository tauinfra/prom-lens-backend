package controller

import (
	"net/http"
	"strconv"
	"valyria-backend/internal/apps/foundry/model"
	"valyria-backend/internal/apps/foundry/service"
	"valyria-backend/internal/core/logger"
	pg "valyria-backend/internal/core/pagination"

	"github.com/gin-gonic/gin"
)

// PipelineController 定义控制器结构体
type PipelineController struct {
	pipeline service.PipelineManager // 使用服务接口
}

func NewPipelineController(pipeline service.PipelineManager) *PipelineController {
	return &PipelineController{pipeline: pipeline}
}

func (c *PipelineController) List(ctx *gin.Context) {
	envID, _ := strconv.Atoi(ctx.Param("envID"))
	page, _ := strconv.Atoi(ctx.Query("page"))
	size, _ := strconv.Atoi(ctx.Query("size"))
	// 分页
	params := pg.QueryParams{
		Page:      page,
		Size:      size,
		SortBy:    ctx.Query("sortBy"),
		SortOrder: ctx.Query("sortOrder"),
		Keyword:   ctx.Query("keyword"),
		Preloads:  []string{"Environment.Project", "Application"},
		Filters:   map[string]interface{}{"environment_id": envID},
	}

	data, pagination, err := c.pipeline.List(ctx, params)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data, "pagination": pagination})
}

func (c *PipelineController) Get(ctx *gin.Context) {
	pipelineID, _ := strconv.Atoi(ctx.Param("pipelineID"))
	data, err := c.pipeline.Get(ctx, pipelineID)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

func (c *PipelineController) Create(ctx *gin.Context) {
	var envID, _ = strconv.Atoi(ctx.Param("envID"))
	var data model.Pipeline
	data.EnvironmentID = envID
	if err := ctx.ShouldBindJSON(&data); err != nil {
		logger.Errorf("Data binding request failed, error: %v", err)
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	data.EnvironmentID = envID
	if err := c.pipeline.Create(ctx, &data); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

func (c *PipelineController) Update(ctx *gin.Context) {
	var data model.Pipeline
	var envID, _ = strconv.Atoi(ctx.Param("envID")) // Env ID
	var pipelineID, _ = strconv.Atoi(ctx.Param("pipelineID"))

	data.EnvironmentID = envID
	if err := ctx.ShouldBindJSON(&data); err != nil {
		logger.Errorf("Data binding request failed, error: %v", err)
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	if err := c.pipeline.Update(ctx, pipelineID, &data); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

func (c *PipelineController) Delete(ctx *gin.Context) {
	var pipelineID, _ = strconv.Atoi(ctx.Param("pipelineID"))

	if err := c.pipeline.Delete(ctx, pipelineID); err != nil {
		logger.Errorf("Data binding request failed, error: %v", err)
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200})
}
