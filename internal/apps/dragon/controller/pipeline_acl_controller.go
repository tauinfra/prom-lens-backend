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

type PipelineACLController struct {
	manager service.PipelineACLManager
}

func NewPipelineACLController(manager service.PipelineACLManager) *PipelineACLController {
	return &PipelineACLController{manager: manager}
}

func (c *PipelineACLController) List(ctx *gin.Context) {
	pipelineID, ok := ginhelper.RequireUintParam(ctx, "pipelineID")
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
	}
	data, pagination, err := c.manager.List(ctx, pipelineID, params)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data, "pagination": pagination})
}

func (c *PipelineACLController) Create(ctx *gin.Context) {
	pipelineID, ok := ginhelper.RequireUintParam(ctx, "pipelineID")
	if !ok {
		return
	}
	var req request.CreatePipelineACLRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	creator := ctx.GetString("username")
	if err := c.manager.Create(ctx, pipelineID, req.UserID, req.Action, creator); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200})
}

func (c *PipelineACLController) BatchCreate(ctx *gin.Context) {
	pipelineID, ok := ginhelper.RequireUintParam(ctx, "pipelineID")
	if !ok {
		return
	}
	var req request.BatchCreatePipelineACLRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	creator := ctx.GetString("username")
	if err := c.manager.BatchCreate(ctx, pipelineID, req.Users, req.Action, creator); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200})
}

func (c *PipelineACLController) Delete(ctx *gin.Context) {
	aclID, ok := ginhelper.RequireUintParam(ctx, "aclID")
	if !ok {
		return
	}
	if err := c.manager.Delete(ctx, aclID); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200})
}
