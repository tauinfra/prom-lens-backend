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

// ProjectController 定义控制器结构体
type ProjectController struct {
	cred service.ProjectManager // 使用服务接口
}

func NewProjectController(cred service.ProjectManager) *ProjectController {
	return &ProjectController{cred: cred}
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

	data, pagination, err := c.cred.List(ctx, params)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data, "pagination": pagination})
}

func (c *ProjectController) Get(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("projectID"))
	data, err := c.cred.Get(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

func (c *ProjectController) Create(ctx *gin.Context) {
	var data model.Project
	if err := ctx.ShouldBindJSON(&data); err != nil {
		logger.Errorf("Data binding request failed, error: %v", err)
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	if err := c.cred.Create(ctx, &data); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

func (c *ProjectController) Update(ctx *gin.Context) {
	var data model.Project
	var id, _ = strconv.Atoi(ctx.Param("id"))

	if err := ctx.ShouldBindJSON(&data); err != nil {
		logger.Errorf("Data binding request failed, error: %v", err)
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	if err := c.cred.Update(ctx, id, &data); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

func (c *ProjectController) Delete(ctx *gin.Context) {
	var id, _ = strconv.Atoi(ctx.Param("id"))

	if err := c.cred.Delete(ctx, id); err != nil {
		logger.Errorf("Data binding request failed, error: %v", err)
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200})
}
