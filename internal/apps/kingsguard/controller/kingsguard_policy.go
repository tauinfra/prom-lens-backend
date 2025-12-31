package controller

import (
	"net/http"
	"strconv"
	"valyria-backend/internal/apps/kingsguard/model"
	"valyria-backend/internal/apps/kingsguard/service"
	"valyria-backend/internal/core/logger"
	pg "valyria-backend/internal/core/pagination"

	"github.com/gin-gonic/gin"
)

type PolicyController struct {
	policy service.PolicyService // 使用服务接口
}

func NewPolicyController(policy service.PolicyService) *PolicyController {
	return &PolicyController{policy: policy}
}

func (c *PolicyController) List(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.Query("page"))
	size, _ := strconv.Atoi(ctx.Query("size"))
	// 分页
	params := pg.QueryParams{
		Page:      page,
		Size:      size,
		SortBy:    ctx.Query("sortBy"),
		SortOrder: ctx.Query("sortOrder"),
		Keyword:   ctx.Query("keyword"), // 搜索关键字
	}
	data, pagination, err := c.policy.List(ctx, params)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data, "pagination": pagination})
}

func (c *PolicyController) Create(ctx *gin.Context) {
	var data model.CasbinPolicy
	if err := ctx.ShouldBindJSON(&data); err != nil {
		logger.Errorf("Data binding request failed, error: %v", err)
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	success, err := c.policy.Create(data)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	if !success {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": "user policies already exist."})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200})
}

func (c *PolicyController) Delete(ctx *gin.Context) {
	var id, _ = strconv.Atoi(ctx.Param("id"))
	success, err := c.policy.Delete(ctx, id)
	if err != nil {
		logger.Errorf("Data binding request failed, error: %v", err)
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	if !success {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": "KG k8s does not exist."})
		return
	}
	logger.Infof("KG Policy '%d' has been deleted successfully.", id)
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200})
}
