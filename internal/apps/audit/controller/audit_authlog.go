package controller

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"valyria-backend/internal/apps/audit/repository"
	"valyria-backend/internal/apps/audit/service"
	pg "valyria-backend/internal/core/pagination"
)

// AuthLogController 处理集群相关请求的控制器结构体
type AuthLogController struct {
	authLog service.AuthLogService // 使用服务接口
}

func SetupAuthLogController(tx *gorm.DB) *AuthLogController {
	authLog := service.NewAuthLogService(tx, repository.NewAuthLogRepository(tx))
	return &AuthLogController{authLog: authLog}
}

func (c *AuthLogController) List(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.Query("page"))
	size, _ := strconv.Atoi(ctx.Query("size"))

	params := pg.QueryParams{
		Page:      page,
		Size:      size,
		SortBy:    ctx.Query("sortBy"),
		SortOrder: ctx.Query("sortOrder"),
	}
	// 默认排序
	if params.SortBy == "" && params.SortOrder == "" {
		params.SortBy = "created_at"
		params.SortOrder = "desc"
	}

	data, pagination, err := c.authLog.List(params)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data, "pagination": pagination})
}
