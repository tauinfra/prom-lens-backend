package controller

import (
	"net/http"
	"strconv"
	"valyria-backend/internal/apps/auth/model"
	"valyria-backend/internal/apps/auth/service"
	"valyria-backend/internal/core/logger"
	pg "valyria-backend/internal/core/pagination"

	"github.com/gin-gonic/gin"
)

// ModuleController 定义控制器结构体
type ModuleController struct {
	module service.ModuleManager // 使用服务接口
}

func NewModuleController(module service.ModuleManager) *ModuleController {
	return &ModuleController{module: module}
}

func (c *ModuleController) List(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.Query("page"))
	size, _ := strconv.Atoi(ctx.Query("size"))
	// 分页
	params := pg.QueryParams{
		Page:      page,
		Size:      size,
		SortBy:    ctx.Query("sortBy"),
		SortOrder: ctx.Query("sortOrder"),
	}

	data, pagination, err := c.module.List(ctx, params)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data, "pagination": pagination})
}

func (c *ModuleController) Get(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	data, err := c.module.Get(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

func (c *ModuleController) Create(ctx *gin.Context) {
	var data model.Module
	if err := ctx.ShouldBindJSON(&data); err != nil {
		logger.Errorf("Data binding request failed, error: %v", err)
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	if err := c.module.Create(ctx, &data); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	logger.Infof("Module '%v' has been created successfully.", data.Name)
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

func (c *ModuleController) Update(ctx *gin.Context) {
	var data model.Module
	var id, _ = strconv.Atoi(ctx.Param("id"))

	if err := ctx.ShouldBindJSON(&data); err != nil {
		logger.Errorf("Data binding request failed, error: %v", err)
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	if err := c.module.Update(ctx, id, &data); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	logger.Infof("Module '%v' has been updated successfully.", data.Name)
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

func (c *ModuleController) Delete(ctx *gin.Context) {
	var id, _ = strconv.Atoi(ctx.Param("id"))

	if err := c.module.Delete(ctx, id); err != nil {
		logger.Errorf("Data binding request failed, error: %v", err)
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	logger.Infof("ModuleID '%d' has been deleted successfully.", id)
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200})
}
