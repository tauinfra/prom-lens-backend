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

// ApplicationController 定义控制器结构体
type ApplicationController struct {
	cred service.ApplicationManager // 使用服务接口
}

func NewApplicationController(cred service.ApplicationManager) *ApplicationController {
	return &ApplicationController{cred: cred}
}

func (c *ApplicationController) List(ctx *gin.Context) {
	projectID, _ := strconv.Atoi(ctx.Param("projectID"))
	page, _ := strconv.Atoi(ctx.Query("page"))
	size, _ := strconv.Atoi(ctx.Query("size"))
	// 分页
	params := pg.QueryParams{
		Page:      page,
		Size:      size,
		SortBy:    ctx.Query("sortBy"),
		SortOrder: ctx.Query("sortOrder"),
		Keyword:   ctx.Query("keyword"),
		Filters:   map[string]interface{}{"project_id": projectID}, // 查询指定项目关联的 App
	}

	data, pagination, err := c.cred.List(ctx, params)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data, "pagination": pagination})
}

func (c *ApplicationController) Get(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("appID"))
	data, err := c.cred.Get(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

func (c *ApplicationController) Create(ctx *gin.Context) {
	var projectID, _ = strconv.Atoi(ctx.Param("projectID"))
	var data model.Application
	if err := ctx.ShouldBindJSON(&data); err != nil {
		logger.Errorf("Data binding request failed, error: %v", err)
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	data.ProjectID = projectID
	if err := c.cred.Create(ctx, &data); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

func (c *ApplicationController) Update(ctx *gin.Context) {
	var data model.Application
	var id, _ = strconv.Atoi(ctx.Param("appID"))

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

func (c *ApplicationController) Delete(ctx *gin.Context) {
	var id, _ = strconv.Atoi(ctx.Param("appID"))

	if err := c.cred.Delete(ctx, id); err != nil {
		logger.Errorf("Data binding request failed, error: %v", err)
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200})
}
