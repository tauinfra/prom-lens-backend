package controller

import (
	"net/http"
	"strconv"
	"prom-lens-backend/internal/apps/audit/service"
	pg "prom-lens-backend/internal/core/pagination"

	"github.com/gin-gonic/gin"
)

// AuthLogController 处理登录日志请求的控制器结构体
type AuthLogController struct {
	authLog service.AuthLogService // 使用服务接口
}

func NewAuthLogController(authLog service.AuthLogService) *AuthLogController {
	return &AuthLogController{authLog: authLog}
}

func (c *AuthLogController) List(ctx *gin.Context) {
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

	data, pagination, err := c.authLog.List(ctx, params)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data, "pagination": pagination})
}
