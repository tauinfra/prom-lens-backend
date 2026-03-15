package controller

import (
	"net/http"
	"strconv"
	"valyria-backend/internal/apps/prometheus/request"
	"valyria-backend/internal/apps/prometheus/service"
	pg "valyria-backend/internal/core/pagination"
	"valyria-backend/internal/pkg/ginhelper"

	"github.com/gin-gonic/gin"
)

// RuleController 控制器结构体
type RuleController struct {
	rule service.RuleManager // 使用服务接口
}

func NewRuleController(rule service.RuleManager) *RuleController {
	return &RuleController{rule: rule}
}

func (c *RuleController) List(ctx *gin.Context) {
	groupIDU, ok := ginhelper.RequireUintParam(ctx, "groupID")
	if !ok {
		return
	}
	groupID := int(groupIDU)
	page, _ := strconv.Atoi(ctx.Query("page"))
	size, _ := strconv.Atoi(ctx.Query("size"))
	// 分页
	params := pg.QueryParams{
		Page:      page,
		Size:      size,
		SortBy:    ctx.Query("sortBy"),
		SortOrder: ctx.Query("sortOrder"),
		Keyword:   ctx.Query("keyword"),
		Filters:   map[string]interface{}{"group_id": groupID}, // Group 查询条件
	}
	data, pagination, err := c.rule.List(ctx, params)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data, "pagination": pagination})
}

func (c *RuleController) Get(ctx *gin.Context) {
	ruleIDU, ok := ginhelper.RequireUintParam(ctx, "ruleID")
	if !ok {
		return
	}
	ruleID := int(ruleIDU)
	data, err := c.rule.Get(ctx, ruleID)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

func (c *RuleController) Create(ctx *gin.Context) {
	groupIDU, ok := ginhelper.RequireUintParam(ctx, "groupID")
	if !ok {
		return
	}
	groupID := int(groupIDU)
	var req request.CreateRuleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	if err := c.rule.Create(ctx, groupID, &req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200})
}

func (c *RuleController) Update(ctx *gin.Context) {
	ruleIDU, ok := ginhelper.RequireUintParam(ctx, "ruleID")
	if !ok {
		return
	}
	var req request.UpdateRuleRequest
	ruleID := int(ruleIDU)
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	if err := c.rule.Update(ctx, ruleID, &req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200})
}

func (c *RuleController) Delete(ctx *gin.Context) {
	idU, ok := ginhelper.RequireUintParam(ctx, "ruleID")
	if !ok {
		return
	}
	id := int(idU)
	if err := c.rule.Delete(ctx, id); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200})
}
