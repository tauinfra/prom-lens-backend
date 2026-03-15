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

type ReviewConfigController struct {
	manager service.ReviewConfigManager
}

func NewReviewConfigController(manager service.ReviewConfigManager) *ReviewConfigController {
	return &ReviewConfigController{manager: manager}
}

func (c *ReviewConfigController) ListGlobal(ctx *gin.Context) {
	c.list(ctx)
}

func (c *ReviewConfigController) GetGlobal(ctx *gin.Context) {
	c.get(ctx)
}

func (c *ReviewConfigController) CreateGlobal(ctx *gin.Context) {
	c.create(ctx)
}

func (c *ReviewConfigController) UpdateGlobal(ctx *gin.Context) {
	c.update(ctx)
}

func (c *ReviewConfigController) DeleteGlobal(ctx *gin.Context) {
	c.delete(ctx)
}

func (c *ReviewConfigController) list(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.Query("page"))
	size, _ := strconv.Atoi(ctx.Query("size"))
	params := pg.QueryParams{
		Page:      page,
		Size:      size,
		SortBy:    ctx.Query("sortBy"),
		SortOrder: ctx.Query("sortOrder"),
		Keyword:   ctx.Query("keyword"),
	}
	data, pagination, err := c.manager.List(ctx, params)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data, "pagination": pagination})
}

func (c *ReviewConfigController) get(ctx *gin.Context) {
	stageID, ok := ginhelper.RequireUintParam(ctx, "stageID")
	if !ok {
		return
	}
	data, err := c.manager.Get(ctx, stageID)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

func (c *ReviewConfigController) create(ctx *gin.Context) {
	var req request.CreateReviewStageRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	creator := ctx.GetString("username")
	if err := c.manager.Create(ctx, &req, creator); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200})
}

func (c *ReviewConfigController) update(ctx *gin.Context) {
	stageID, ok := ginhelper.RequireUintParam(ctx, "stageID")
	if !ok {
		return
	}
	var req request.UpdateReviewStageRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	if err := c.manager.Update(ctx, stageID, &req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200})
}

func (c *ReviewConfigController) delete(ctx *gin.Context) {
	stageID, ok := ginhelper.RequireUintParam(ctx, "stageID")
	if !ok {
		return
	}
	if err := c.manager.Delete(ctx, stageID); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200})
}
