package controller

import (
	"net/http"
	"strconv"
	"valyria-backend/internal/apps/dragon/model"
	"valyria-backend/internal/apps/dragon/request"
	"valyria-backend/internal/apps/dragon/service"
	pg "valyria-backend/internal/core/pagination"
	"valyria-backend/internal/pkg/ginhelper"

	"github.com/gin-gonic/gin"
)

// CredentialController 定义控制器结构体
type CredentialController struct {
	manager service.CredentialManager // 使用服务接口
}

func NewCredentialController(manager service.CredentialManager) *CredentialController {
	return &CredentialController{manager: manager}
}

func (c *CredentialController) List(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.Query("page"))
	size, _ := strconv.Atoi(ctx.Query("size"))
	typeParam := ctx.Query("type")
	// 分页
	params := pg.QueryParams{
		Page:      page,
		Size:      size,
		SortBy:    ctx.Query("sortBy"),
		SortOrder: ctx.Query("sortOrder"),
		Keyword:   ctx.Query("keyword"),
	}
	// 过滤凭证类型
	if typeParam != "" {
		params.Filters = map[string]interface{}{"type": typeParam}
	}
	// 查询结果
	data, pagination, err := c.manager.List(ctx, params)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data, "pagination": pagination})
}

func (c *CredentialController) Get(ctx *gin.Context) {
	id, ok := ginhelper.RequireUintParam(ctx, "credentialID")
	if !ok {
		return
	}
	data, err := c.manager.Get(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200, "data": data})
}

func (c *CredentialController) Create(ctx *gin.Context) {
	var req request.CreateCredentialRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": -1, "msg": err.Error()})
		return
	}
	if !model.IsValidCredentialType(req.Type) {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": "invalid credential type"})
		return
	}
	// 获取当前用户
	username, exists := ctx.Get("username")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"success": false, "code": http.StatusUnauthorized, "msg": "用户未登录"})
		return
	}
	// 绑定 creator
	data := model.Credential{
		Name:            req.Name,
		Type:            req.Type,
		SecretType:      req.SecretType,
		BaseURL:         req.BaseURL,
		Username:        req.Username,
		EncryptedSecret: req.EncryptedSecret,
		Creator:         username.(string),
	}
	if err := c.manager.Create(ctx, &data); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200})
}

func (c *CredentialController) Update(ctx *gin.Context) {
	id, ok := ginhelper.RequireUintParam(ctx, "credentialID")
	if !ok {
		return
	}
	var req request.UpdateCredentialRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	data := &model.Credential{}
	if req.Name != nil {
		data.Name = *req.Name
	}
	if req.SecretType != nil {
		data.SecretType = *req.SecretType
	}
	if req.BaseURL != nil {
		data.BaseURL = *req.BaseURL
	}
	if req.Username != nil {
		data.Username = *req.Username
	}
	if req.EncryptedSecret != nil {
		data.EncryptedSecret = *req.EncryptedSecret
	}
	if err := c.manager.Update(ctx, id, data); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200})
}

func (c *CredentialController) Delete(ctx *gin.Context) {
	id, ok := ginhelper.RequireUintParam(ctx, "credentialID")
	if !ok {
		return
	}
	if err := c.manager.Delete(ctx, id); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"success": false, "code": 20000, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true, "code": 200})
}
