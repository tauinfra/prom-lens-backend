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

// ReviewController 定义控制器结构体
type ReviewController struct {
	review service.ReviewManager // 使用服务接口
}

func NewReviewController(review service.ReviewManager) *ReviewController {
	return &ReviewController{review: review}
}

func (c *ReviewController) List(ctx *gin.Context) {
	releaseID, ok := ginhelper.RequireUintParam(ctx, "releaseID")
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
		Filters:   map[string]interface{}{"release_id": releaseID},
	}
	data, pagination, err := c.review.List(ctx, params)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data, "pagination": pagination})
}

func (c *ReviewController) Get(ctx *gin.Context) {
	reviewID, ok := ginhelper.RequireUintParam(ctx, "reviewID")
	if !ok {
		return
	}
	data, err := c.review.Get(ctx, reviewID)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

func (c *ReviewController) Create(ctx *gin.Context) {
	var req request.CreateReviewRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	userID := ctx.MustGet("uid").(int)
	if err := c.review.Create(ctx, &req, uint(userID)); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200})
}

func (c *ReviewController) Update(ctx *gin.Context) {
	reviewID, ok := ginhelper.RequireUintParam(ctx, "reviewID")
	if !ok {
		return
	}
	userID := ctx.MustGet("uid").(int)
	var req request.UpdateReviewRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	if err := c.review.Update(ctx, reviewID, uint(userID), &req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200})
}

func (c *ReviewController) Delete(ctx *gin.Context) {
	reviewID, ok := ginhelper.RequireUintParam(ctx, "reviewID")
	if !ok {
		return
	}
	if err := c.review.Delete(ctx, reviewID); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200})
}
